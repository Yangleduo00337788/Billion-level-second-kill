package com.seckill.usersystem.service.impl;

import cn.hutool.crypto.digest.BCrypt;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.dto.LoginRequest;
import com.seckill.usersystem.dto.RegisterRequest;
import com.seckill.usersystem.dto.TokenResponse;
import com.seckill.usersystem.entity.Blacklist;
import com.seckill.usersystem.entity.Device;
import com.seckill.usersystem.entity.LoginLog;
import com.seckill.usersystem.entity.User;
import com.seckill.usersystem.entity.Role;
import com.seckill.usersystem.entity.UserProfile;
import com.seckill.usersystem.enums.DeviceTypeEnum;
import com.seckill.usersystem.enums.LoginTypeEnum;
import com.seckill.usersystem.enums.RiskLevelEnum;
import com.seckill.usersystem.enums.UserStatusEnum;
import com.seckill.usersystem.exception.BizException;
import com.seckill.usersystem.mapper.*;
import com.seckill.usersystem.service.*;
import com.seckill.usersystem.util.JwtUtil;
import com.seckill.usersystem.util.RedisUtil;
import com.seckill.usersystem.vo.LoginResultVO;
import com.seckill.usersystem.vo.UserVO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.*;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@RequiredArgsConstructor
public class UserServiceImpl implements IUserService {

    private final UserMapper userMapper;
    private final UserProfileMapper userProfileMapper;
    private final DeviceMapper deviceMapper;
    private final LoginLogMapper loginLogMapper;
    private final JwtUtil jwtUtil;
    private final RedisUtil redisUtil;
    private final IRiskControlService riskControlService;
    private final IDeviceService deviceService;
    private final IBlacklistService blacklistService;
    private final IRoleService roleService;
    private final IPermissionService permissionService;
    private final ObjectMapper objectMapper;
    private final ThreadPoolTaskExecutor loginLogExecutor;

    @Value("${user.max-device-count:5}")
    private int maxDeviceCount;

    @Value("${user.login-lock-minutes:30}")
    private long loginLockMinutes;

    private static final String TOKEN_BLACKLIST_PREFIX = "token:blacklist:";
    private static final String USER_LOCK_PREFIX = "user:lock:";
    private static final String USER_CACHE_PREFIX = "user:info:";
    private static final String USERNAME_ID_CACHE_PREFIX = "user:uname:id:";
    private static final String TOKEN_REFRESH_RATE_KEY = "rate:token:refresh:";

    @Override
    @Transactional(rollbackFor = Exception.class)
    public LoginResultVO login(LoginRequest request) {
        String username = request.getUsername();
        String ip = getClientIp();
        String deviceId = request.getDeviceId();

        log.info("用户登录尝试: username={}, ip={}", username, ip);

        if (blacklistService.isIpInBlacklist(ip) || riskControlService.isIpBlocked(ip)) {
            throw BizException.of(403, "您的IP已被封禁");
        }

        User user = getUserByUsername(username);
        if (user == null) {
            recordLoginLog(null, username, request.getLoginType(), ip, deviceId, 0, "用户不存在");
            riskControlService.recordLoginAttempt(null, ip, false);
            throw BizException.of("用户名或密码错误");
        }

        if (blacklistService.isUserInBlacklist(user.getId())) {
            throw BizException.of(403, "您的账号已被封禁");
        }

        if (riskControlService.isAccountLocked(user.getId())) {
            throw BizException.of(403, "账号已被锁定，请稍后重试");
        }

        if (user.getIsDeleted() == 1) {
            throw BizException.of("账号已注销");
        }

        if (user.getStatus() == UserStatusEnum.DISABLED.getCode()) {
            throw BizException.of("账号已被禁用");
        }

        if (request.getLoginType() == LoginTypeEnum.PASSWORD.getCode()) {
            if (!isPasswordCorrect(user, request.getPassword())) {
                riskControlService.recordLoginAttempt(user.getId(), ip, false);
                int failedCount = riskControlService.getFailedLoginCount(user.getId(), ip);
                if (failedCount >= 5) {
                    riskControlService.triggerRiskAction(user.getId(), RiskLevelEnum.HIGH);
                    lockUser(user.getId(), loginLockMinutes);
                    throw BizException.of("密码错误次数过多，账号已锁定30分钟");
                }
                recordLoginLog(user.getId(), username, request.getLoginType(), ip, deviceId, 0, "密码错误");
                throw BizException.of("用户名或密码错误");
            }
        }

        RiskLevelEnum riskLevel = riskControlService.assessLoginRisk(
                user.getId(), ip, deviceId, request.getDeviceType());

        if (!riskControlService.isLoginAllowed(riskLevel)) {
            throw BizException.of(403, "登录环境异常，请验证身份后重试");
        }

        if (riskControlService.isCaptchaRequired(riskLevel)) {
            log.warn("用户登录需要验证码: userId={}, riskLevel={}, ip={}", user.getId(), riskLevel, ip);
            riskControlService.requireCaptcha(ip, deviceId);
            throw BizException.of(429, "登录存在风险，请完成验证码验证");
        }

        if (riskControlService.isSmsCodeRequired(riskLevel)) {
            log.warn("用户登录需要短信验证: userId={}, riskLevel={}", user.getId(), riskLevel);
            throw BizException.of(429, "登录存在高风险，请完成短信验证");
        }

        TokenResponse tokenResponse = generateTokens(user, deviceId);

        String deviceType = request.getDeviceType() != null ? request.getDeviceType() : DeviceTypeEnum.WEB.getCode();
        if (deviceService.hasSameTypeDevice(user.getId(), deviceType)) {
            deviceService.removeSameTypeDevices(user.getId(), deviceType, deviceId);
            log.info("同类型设备互踢: userId={}, deviceType={}, newDeviceId={}", user.getId(), deviceType, deviceId);
        }

        Device device = deviceService.getDevice(user.getId(), deviceId);
        if (device == null) {
            if (deviceService.getDeviceCount(user.getId()) >= maxDeviceCount) {
                deviceService.removeAllOtherDevices(user.getId(), deviceId);
            }
            device = new Device();
            device.setUserId(user.getId());
            device.setDeviceId(deviceId);
            device.setDeviceType(deviceType);
            device.setDeviceName(parseDeviceName(request.getDeviceType()));
            device.setIpAddress(ip);
            device.setLastLoginTime(LocalDateTime.now());
            device.setLastActiveTime(LocalDateTime.now());
            device.setStatus(1);
            device.setIsTrusted(riskLevel == RiskLevelEnum.LOW ? 1 : 0);
            deviceMapper.insert(device);
        } else {
            device.setLastLoginTime(LocalDateTime.now());
            device.setLastActiveTime(LocalDateTime.now());
            device.setIpAddress(ip);
            device.setStatus(1);
            deviceMapper.update(device);
        }

        riskControlService.clearFailedLoginCount(user.getId(), ip);

        updateLastLoginInfo(user.getId(), ip);

        recordLoginLog(user.getId(), username, request.getLoginType(), ip, deviceId, 1, null);

        return buildLoginResult(user, tokenResponse, device);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public LoginResultVO register(RegisterRequest request) {
        if (!request.getPassword().equals(request.getConfirmPassword())) {
            throw BizException.of("两次密码输入不一致");
        }

        if (getUserByUsername(request.getUsername()) != null) {
            throw BizException.of("用户名已存在");
        }

        User user = new User();
        user.setUsername(request.getUsername());
        user.setPassword(BCrypt.hashpw(request.getPassword(), BCrypt.gensalt()));
        user.setPhone(request.getPhone());
        user.setEmail(request.getEmail());
        user.setStatus(UserStatusEnum.ENABLED.getCode());
        user.setIsDeleted(0);
        user.setLoginType(LoginTypeEnum.PASSWORD.getCode());
        user.setLoginCount(0);
        user.setPasswordUpdateTime(LocalDateTime.now());
        userMapper.insert(user);

        UserProfile profile = new UserProfile();
        profile.setUserId(user.getId());
        profile.setNickname(request.getUsername());
        userProfileMapper.insert(profile);

        LoginResultVO vo = new LoginResultVO();
        vo.setUserId(user.getId());
        vo.setUsername(user.getUsername());
        vo.setNickname(profile.getNickname());
        vo.setIsFirstLogin(true);
        return vo;
    }

    @Override
    public void logout(Long userId, String deviceId) {
        String tokenKey = "user:token:" + userId + ":" + deviceId;
        redisUtil.delete(tokenKey);

        String refreshTokenKey = "user:refresh:" + userId + ":" + deviceId;
        redisUtil.delete(refreshTokenKey);

        Device device = deviceService.getDevice(userId, deviceId);
        if (device != null) {
            device.setStatus(0);
            deviceMapper.update(device);
        }

        log.info("用户退出登录: userId={}, deviceId={}", userId, deviceId);
    }

    @Override
    public TokenResponse refreshToken(String refreshToken, String deviceId) {
        if (!jwtUtil.isTokenValid(refreshToken)) {
            throw BizException.of(401, "刷新令牌已失效");
        }

        Long userId = jwtUtil.getUserId(refreshToken);

        String rateKey = TOKEN_REFRESH_RATE_KEY + userId;
        Long refreshCount = redisUtil.increment(rateKey);
        redisUtil.expire(rateKey, 60);
        if (refreshCount > 50) {
            log.warn("Token刷新频率超限: userId={}, count={}", userId, refreshCount);
            throw BizException.of(429, "刷新过于频繁，请稍后重试");
        }

        String storedRefreshToken = redisUtil.get("user:refresh:" + userId + ":" + deviceId);

        if (storedRefreshToken == null || !storedRefreshToken.equals(refreshToken)) {
            redisUtil.deleteByPattern("user:token:" + userId + ":*");
            redisUtil.deleteByPattern("user:refresh:" + userId + ":*");
            throw BizException.of(401, "刷新令牌已被撤销");
        }

        String newAccessToken = jwtUtil.generateAccessToken(
                userId, jwtUtil.getUsername(refreshToken), deviceId, getUserRoles(userId),
                getUserPermissions(userId), jwtUtil.getSessionId(refreshToken), getClientIp());
        String newRefreshToken = jwtUtil.generateRefreshToken(userId, deviceId);

        redisUtil.set("user:token:" + userId + ":" + deviceId, newAccessToken,
                jwtUtil.getAccessTokenExpiration(), TimeUnit.SECONDS);
        redisUtil.set("user:refresh:" + userId + ":" + deviceId, newRefreshToken,
                jwtUtil.getRefreshTokenExpiration(), TimeUnit.SECONDS);

        deviceService.updateDeviceActivity(userId, deviceId);

        TokenResponse response = new TokenResponse();
        response.setAccessToken(newAccessToken);
        response.setRefreshToken(newRefreshToken);
        response.setExpiresIn(jwtUtil.getAccessTokenExpiration());
        return response;
    }

    @Override
    public User getUserByUsername(String username) {
        String userIdKey = USERNAME_ID_CACHE_PREFIX + username;
        String userIdStr = redisUtil.get(userIdKey);
        if (userIdStr != null) {
            return getUserById(Long.parseLong(userIdStr));
        }

        User user = QueryChain.of(userMapper)
                .where(User::getUsername).eq(username)
                .one();

        if (user != null) {
            redisUtil.set(userIdKey, String.valueOf(user.getId()), 3600, TimeUnit.SECONDS);
            cacheUser(user);
        }
        return user;
    }

    @Override
    public User getUserById(Long userId) {
        String cacheKey = USER_CACHE_PREFIX + userId;
        String cachedJson = redisUtil.get(cacheKey);
        if (cachedJson != null) {
            try {
                return objectMapper.readValue(cachedJson, User.class);
            } catch (Exception e) {
                log.warn("解析用户缓存失败: userId={}", userId);
            }
        }

        User user = userMapper.selectOneById(userId);
        if (user != null) {
            cacheUser(user);
        }
        return user;
    }

    private void cacheUser(User user) {
        try {
            String json = objectMapper.writeValueAsString(user);
            redisUtil.set(USER_CACHE_PREFIX + user.getId(), json, 3600, TimeUnit.SECONDS);
            redisUtil.set(USERNAME_ID_CACHE_PREFIX + user.getUsername(),
                    String.valueOf(user.getId()), 3600, TimeUnit.SECONDS);
        } catch (Exception e) {
            log.warn("缓存用户信息失败: userId={}", user.getId(), e);
        }
    }

    private void invalidateUserCache(Long userId, String username) {
        redisUtil.delete(USER_CACHE_PREFIX + userId);
        if (username != null) {
            redisUtil.delete(USERNAME_ID_CACHE_PREFIX + username);
        }
    }

    @Override
    public List<UserVO> listUsers(int pageNum, int pageSize) {
        return userMapper.selectAll().stream().map(this::convertToUserVO).toList();
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void updateUserStatus(Long userId, Integer status) {
        User user = getUserById(userId);
        if (user == null) {
            throw BizException.of("用户不存在");
        }
        user.setStatus(status);
        userMapper.update(user);
        invalidateUserCache(userId, user.getUsername());

        if (status == UserStatusEnum.DISABLED.getCode() || status == UserStatusEnum.LOCKED.getCode()) {
            redisUtil.deleteByPattern("user:token:" + userId + ":*");
            redisUtil.deleteByPattern("user:refresh:" + userId + ":*");
        }
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void deleteUser(Long userId) {
        User user = getUserById(userId);
        if (user == null) {
            throw BizException.of("用户不存在");
        }
        user.setIsDeleted(1);
        userMapper.update(user);
        invalidateUserCache(userId, user.getUsername());

        redisUtil.deleteByPattern("user:token:" + userId + ":*");
        redisUtil.deleteByPattern("user:refresh:" + userId + ":*");
    }

    @Override
    public void lockUser(Long userId, long lockMinutes) {
        String lockKey = USER_LOCK_PREFIX + userId;
        redisUtil.set(lockKey, "1", lockMinutes, TimeUnit.MINUTES);
        log.info("用户被锁定: userId={}, lockMinutes={}", userId, lockMinutes);
    }

    @Override
    public void resetPassword(Long userId, String newPassword) {
        User user = getUserById(userId);
        if (user == null) {
            throw BizException.of("用户不存在");
        }
        user.setPassword(BCrypt.hashpw(newPassword, BCrypt.gensalt()));
        user.setPasswordUpdateTime(LocalDateTime.now());
        userMapper.update(user);

        redisUtil.deleteByPattern("user:token:" + userId + ":*");
        redisUtil.deleteByPattern("user:refresh:" + userId + ":*");
    }

    @Override
    public Set<String> getUserRoles(Long userId) {
        return roleService.getRolesByUserId(userId).stream()
                .map(Role::getRoleCode)
                .collect(java.util.stream.Collectors.toSet());
    }

    @Override
    public Set<String> getUserPermissions(Long userId) {
        return new HashSet<>(permissionService.getPermissionCodesByUserId(userId));
    }

    @Override
    public boolean isPasswordCorrect(User user, String password) {
        return BCrypt.checkpw(password, user.getPassword());
    }

    @Override
    public void updateLastLoginInfo(Long userId, String ip) {
        User user = new User();
        user.setId(userId);
        user.setLastLoginTime(LocalDateTime.now());
        user.setLastLoginIp(ip);
        user.setLoginCount(getUserById(userId).getLoginCount() + 1);
        userMapper.update(user);
    }

    private TokenResponse generateTokens(User user, String deviceId) {
        Set<String> roles = getUserRoles(user.getId());
        Set<String> permissions = getUserPermissions(user.getId());
        String sessionId = UUID.randomUUID().toString();
        String ip = getClientIp();

        String accessToken = jwtUtil.generateAccessToken(
                user.getId(), user.getUsername(), deviceId, roles, permissions, sessionId, ip);
        String refreshToken = jwtUtil.generateRefreshToken(user.getId(), deviceId);

        redisUtil.set("user:token:" + user.getId() + ":" + deviceId, accessToken,
                jwtUtil.getAccessTokenExpiration(), TimeUnit.SECONDS);
        redisUtil.set("user:refresh:" + user.getId() + ":" + deviceId, refreshToken,
                jwtUtil.getRefreshTokenExpiration(), TimeUnit.SECONDS);

        saveSession(user.getId(), sessionId, deviceId, ip);

        TokenResponse response = new TokenResponse();
        response.setAccessToken(accessToken);
        response.setRefreshToken(refreshToken);
        response.setExpiresIn(jwtUtil.getAccessTokenExpiration());
        return response;
    }

    private void saveSession(Long userId, String sessionId, String deviceId, String ip) {
        String sessionKey = "user:session:" + sessionId;
        redisUtil.hashSet(sessionKey, "userId", String.valueOf(userId));
        redisUtil.hashSet(sessionKey, "deviceId", deviceId);
        redisUtil.hashSet(sessionKey, "ip", ip);
        redisUtil.hashSet(sessionKey, "loginTime", LocalDateTime.now().toString());
        redisUtil.hashSet(sessionKey, "lastActiveTime", LocalDateTime.now().toString());
        redisUtil.expire(sessionKey, 30 * 24 * 3600);
    }

    private LoginResultVO buildLoginResult(User user, TokenResponse token, Device device) {
        UserProfile profile = QueryChain.of(userProfileMapper)
                .where(UserProfile::getUserId).eq(user.getId())
                .one();

        LoginResultVO result = new LoginResultVO();
        result.setUserId(user.getId());
        result.setUsername(user.getUsername());
        result.setNickname(profile != null ? profile.getNickname() : user.getUsername());
        result.setAvatar(profile != null ? profile.getAvatar() : null);
        result.setToken(token);
        result.setIsFirstLogin(user.getLoginCount() <= 1);
        result.setNeedChangePassword(user.getPasswordUpdateTime() == null);
        result.setMaxDeviceCount(maxDeviceCount);
        result.setCurrentDeviceCount(deviceService.getDeviceCount(user.getId()));
        return result;
    }

    private UserVO convertToUserVO(User user) {
        UserVO vo = new UserVO();
        vo.setId(user.getId());
        vo.setUsername(user.getUsername());
        vo.setPhone(user.getPhone());
        vo.setEmail(user.getEmail());
        vo.setStatus(user.getStatus());
        vo.setStatusText(UserStatusEnum.fromCode(user.getStatus()).getDescription());
        vo.setLastLoginTime(user.getLastLoginTime());
        vo.setLastLoginIp(user.getLastLoginIp());
        vo.setLoginCount(user.getLoginCount());
        vo.setCreateTime(user.getCreateTime());
        vo.setRoles(new ArrayList<>(getUserRoles(user.getId())));
        vo.setPermissions(new ArrayList<>(getUserPermissions(user.getId())));
        return vo;
    }

    private void recordLoginLog(Long userId, String username, Integer loginType,
                                String ip, String deviceId, Integer status, String failReason) {
        loginLogExecutor.execute(() -> {
            try {
                LoginLog loginLog = new LoginLog();
                loginLog.setUserId(userId != null ? userId : 0L);
                loginLog.setUsername(username);
                loginLog.setLoginType(loginType);
                loginLog.setIpAddress(ip);
                loginLog.setDeviceId(deviceId);
                loginLog.setStatus(status);
                loginLog.setFailReason(failReason);
                loginLog.setLoginTime(LocalDateTime.now());
                loginLogMapper.insert(loginLog);
            } catch (Exception e) {
                log.error("异步写入登录日志失败: username={}", username, e);
            }
        });
    }

    private String getClientIp() {
        return "127.0.0.1";
    }

    private String parseDeviceName(String deviceType) {
        if (deviceType == null) return "未知设备";
        return switch (deviceType.toUpperCase()) {
            case "WEB" -> "网页浏览器";
            case "IOS" -> "iPhone/iPad";
            case "ANDROID" -> "Android设备";
            case "PC" -> "Windows电脑";
            case "MAC" -> "Mac电脑";
            default -> "未知设备";
        };
    }
}
