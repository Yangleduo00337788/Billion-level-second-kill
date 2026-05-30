package com.seckill.usersystem.service.impl;

import com.seckill.usersystem.enums.RiskLevelEnum;
import com.seckill.usersystem.service.IBlacklistService;
import com.seckill.usersystem.service.IDeviceService;
import com.seckill.usersystem.service.IRiskControlService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.time.LocalDateTime;
import java.util.Map;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@RequiredArgsConstructor
public class RiskControlServiceImpl implements IRiskControlService {

    private final RedisUtil redisUtil;
    private final IBlacklistService blacklistService;
    private final IDeviceService deviceService;

    @Value("${risk-control.max-login-attempts:5}")
    private int maxLoginAttempts;

    @Value("${risk-control.attempt-window-minutes:15}")
    private int attemptWindowMinutes;

    @Value("${risk-control.max-device-changes-per-hour:3}")
    private int maxDeviceChangesPerHour;

    private static final String LOGIN_ATTEMPT_KEY = "risk:login:attempt:";
    private static final String LOGIN_TIME_KEY = "risk:login:time:";
    private static final String DEVICE_CHANGE_KEY = "risk:device:change:";
    private static final String ACCOUNT_LOCK_KEY = "risk:account:lock:";
    private static final String RISK_ACTION_KEY = "risk:action:";
    private static final String IP_BLOCK_KEY = "risk:ip:block:";
    private static final String CAPTCHA_REQUIRED_KEY = "risk:captcha:required:";
    private static final String LAST_LOGIN_LOCATION_KEY = "risk:location:";

    @Override
    public RiskLevelEnum assessLoginRisk(Long userId, String ip, String deviceId, String userAgent) {
        int riskScore = 0;

        if (userId == null) {
            riskScore += 10;
        }

        riskScore += checkIpRisk(ip);

        if (userId != null) {
            riskScore += checkDeviceRisk(userId, deviceId);
            riskScore += checkLoginTimeRisk(userId);
            riskScore += checkDeviceChangeRisk(userId, deviceId);
            riskScore += checkGeoLocationRisk(userId, ip);
        }

        if (riskScore >= 80) {
            return RiskLevelEnum.CRITICAL;
        } else if (riskScore >= 50) {
            return RiskLevelEnum.HIGH;
        } else if (riskScore >= 20) {
            return RiskLevelEnum.MEDIUM;
        }
        return RiskLevelEnum.LOW;
    }

    @Override
    public boolean isLoginAllowed(RiskLevelEnum riskLevel) {
        if (riskLevel == RiskLevelEnum.CRITICAL) {
            return false;
        }
        return true;
    }

    @Override
    public boolean isCaptchaRequired(RiskLevelEnum riskLevel) {
        return riskLevel == RiskLevelEnum.MEDIUM || riskLevel == RiskLevelEnum.HIGH;
    }

    @Override
    public boolean isSmsCodeRequired(RiskLevelEnum riskLevel) {
        return riskLevel == RiskLevelEnum.HIGH;
    }

    @Override
    public void requireCaptcha(String ip, String deviceId) {
        redisUtil.set(CAPTCHA_REQUIRED_KEY + "ip:" + ip, "1", 15, TimeUnit.MINUTES);
        if (deviceId != null) {
            redisUtil.set(CAPTCHA_REQUIRED_KEY + "device:" + deviceId, "1", 15, TimeUnit.MINUTES);
        }
    }

    @Override
    public boolean isCaptchaRequired(String ip, String deviceId) {
        return redisUtil.hasKey(CAPTCHA_REQUIRED_KEY + "ip:" + ip)
                || (deviceId != null && redisUtil.hasKey(CAPTCHA_REQUIRED_KEY + "device:" + deviceId));
    }

    @Override
    public void clearCaptchaRequired(String ip, String deviceId) {
        redisUtil.delete(CAPTCHA_REQUIRED_KEY + "ip:" + ip);
        if (deviceId != null) {
            redisUtil.delete(CAPTCHA_REQUIRED_KEY + "device:" + deviceId);
        }
    }

    @Override
    public void recordLoginAttempt(Long userId, String ip, boolean success) {
        String key;
        if (userId != null) {
            key = LOGIN_ATTEMPT_KEY + "user:" + userId;
        } else {
            key = LOGIN_ATTEMPT_KEY + "ip:" + ip;
        }

        if (!success) {
            long newCount = redisUtil.increment(key);
            redisUtil.expire(key, attemptWindowMinutes * 60L);
            log.info("记录失败登录尝试: key={}, count={}", key, newCount);

            String ipKey = LOGIN_ATTEMPT_KEY + "ip:" + ip;
            long ipCount = redisUtil.increment(ipKey);
            redisUtil.expire(ipKey, attemptWindowMinutes * 60L);

            if (ipCount > 10) {
                blockIp(ip, 24, TimeUnit.HOURS, "单IP登录失败超过10次");
                log.warn("IP已被封禁24小时: ip={}, count={}", ip, ipCount);
            } else if (ipCount > maxLoginAttempts) {
                blockIp(ip, 15, TimeUnit.MINUTES, "单IP登录失败超过5次");
                log.warn("IP已被封禁15分钟: ip={}, count={}", ip, ipCount);
            }
        } else {
            redisUtil.delete(key);
        }
    }

    @Override
    public int getFailedLoginCount(Long userId, String ip) {
        String key;
        if (userId != null) {
            key = LOGIN_ATTEMPT_KEY + "user:" + userId;
        } else {
            key = LOGIN_ATTEMPT_KEY + "ip:" + ip;
        }
        String count = redisUtil.get(key);
        return count != null ? Integer.parseInt(count) : 0;
    }

    @Override
    public void clearFailedLoginCount(Long userId, String ip) {
        if (userId != null) {
            redisUtil.delete(LOGIN_ATTEMPT_KEY + "user:" + userId);
        }
        redisUtil.delete(LOGIN_ATTEMPT_KEY + "ip:" + ip);
    }

    @Override
    public boolean isAccountLocked(Long userId) {
        return redisUtil.hasKey(ACCOUNT_LOCK_KEY + userId);
    }

    @Override
    public boolean isIpBlocked(String ip) {
        return redisUtil.hasKey(IP_BLOCK_KEY + ip);
    }

    @Override
    public void blockIp(String ip, long duration, TimeUnit unit, String reason) {
        redisUtil.set(IP_BLOCK_KEY + ip, reason, duration, unit);
        log.warn("IP被封禁: ip={}, duration={}{}, reason={}", ip, duration, unit, reason);
    }

    @Override
    public void triggerRiskAction(Long userId, RiskLevelEnum riskLevel) {
        String key = RISK_ACTION_KEY + userId;

        switch (riskLevel) {
            case MEDIUM:
                log.warn("中等风险，建议触发验证码: userId={}", userId);
                break;
            case HIGH:
                redisUtil.set(ACCOUNT_LOCK_KEY + userId, "1", 30, TimeUnit.MINUTES);
                redisUtil.set(key, "lock_30min");
                log.warn("高风险操作，锁定账号30分钟: userId={}", userId);
                break;
            case CRITICAL:
                redisUtil.set(ACCOUNT_LOCK_KEY + userId, "1", 24, TimeUnit.HOURS);
                redisUtil.set(key, "lock_24h");
                log.error("极高风险操作，锁定账号24小时: userId={}", userId);
                break;
            default:
                break;
        }
    }

    @Override
    public void recordDeviceChange(Long userId, String oldDeviceId, String newDeviceId, String ip) {
        String key = DEVICE_CHANGE_KEY + userId;
        redisUtil.increment(key);
        redisUtil.expire(key, 3600);

        log.info("记录设备变更: userId={}, oldDevice={}, newDevice={}, ip={}",
                userId, oldDeviceId, newDeviceId, ip);
    }

    @Override
    public void recordLoginLocation(Long userId, String ip, String location) {
        String key = LAST_LOGIN_LOCATION_KEY + userId;
        redisUtil.hashSet(key, "ip", ip);
        redisUtil.hashSet(key, "location", location);
        redisUtil.hashSet(key, "time", LocalDateTime.now().toString());
        redisUtil.expire(key, 86400);
    }

    @Override
    public String getLastLoginLocation(Long userId) {
        Map<Object, Object> data = redisUtil.hashGetAll(LAST_LOGIN_LOCATION_KEY + userId);
        if (data != null && !data.isEmpty()) {
            Object location = data.get("location");
            return location != null ? location.toString() : null;
        }
        return null;
    }

    private int checkIpRisk(String ip) {
        int score = 0;

        if (blacklistService.isIpInBlacklist(ip) || redisUtil.hasKey(IP_BLOCK_KEY + ip)) {
            return 100;
        }

        String loginCountKey = LOGIN_ATTEMPT_KEY + "ip:" + ip;
        String countStr = redisUtil.get(loginCountKey);
        if (countStr != null) {
            int count = Integer.parseInt(countStr);
            if (count > maxLoginAttempts) {
                score += 50;
            } else if (count > maxLoginAttempts / 2) {
                score += 20;
            }
        }

        if (isInternalIp(ip)) {
            score -= 10;
        }

        return Math.max(0, Math.min(100, score));
    }

    private int checkDeviceRisk(Long userId, String deviceId) {
        int score = 0;

        if (blacklistService.isDeviceInBlacklist(deviceId)) {
            return 100;
        }

        if (!deviceService.isDeviceTrusted(userId, deviceId)) {
            score += 10;
        }

        return Math.min(100, score);
    }

    private int checkLoginTimeRisk(Long userId) {
        int score = 0;

        LocalDateTime now = LocalDateTime.now();
        int hour = now.getHour();

        if (hour >= 0 && hour < 6) {
            score += 15;
        }

        String lastLoginTimeKey = LOGIN_TIME_KEY + userId;
        String lastTimeStr = redisUtil.get(lastLoginTimeKey);
        if (lastTimeStr != null) {
            try {
                LocalDateTime lastTime = LocalDateTime.parse(lastTimeStr);
                long minutesBetween = Duration.between(lastTime, now).toMinutes();
                if (minutesBetween < 5) {
                    score += 20;
                }
            } catch (Exception e) {
                log.warn("解析上次登录时间失败: {}", lastTimeStr);
            }
        }

        redisUtil.set(lastLoginTimeKey, now.toString(), 86400);

        return Math.min(100, score);
    }

    private int checkDeviceChangeRisk(Long userId, String deviceId) {
        String deviceChangeKey = DEVICE_CHANGE_KEY + userId;
        String countStr = redisUtil.get(deviceChangeKey);

        if (countStr != null) {
            int count = Integer.parseInt(countStr);
            if (count > maxDeviceChangesPerHour) {
                return 80;
            } else if (count > maxDeviceChangesPerHour / 2) {
                return 30;
            }
        }

        return 0;
    }

    private int checkGeoLocationRisk(Long userId, String ip) {
        if (userId == null) {
            return 0;
        }

        Map<Object, Object> lastLocation = redisUtil.hashGetAll(LAST_LOGIN_LOCATION_KEY + userId);
        if (lastLocation == null || lastLocation.isEmpty()) {
            return 0;
        }

        Object lastIp = lastLocation.get("ip");
        if (lastIp != null && !lastIp.toString().equals(ip)) {
            Object lastTime = lastLocation.get("time");
            if (lastTime != null) {
                try {
                    LocalDateTime lastLogin = LocalDateTime.parse(lastTime.toString());
                    long minutesBetween = Duration.between(lastLogin, LocalDateTime.now()).toMinutes();
                    if (minutesBetween < 30) {
                        return 50;
                    }
                } catch (Exception e) {
                    // ignore parse error
                }
            }
            return 10;
        }

        return 0;
    }

    private boolean isInternalIp(String ip) {
        return ip.startsWith("10.") || ip.startsWith("192.168.") ||
                ip.startsWith("172.") || ip.equals("127.0.0.1") || ip.equals("0:0:0:0:0:0:0:1");
    }
}
