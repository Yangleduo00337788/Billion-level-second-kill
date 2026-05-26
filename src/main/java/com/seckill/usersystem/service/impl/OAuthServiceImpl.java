package com.seckill.usersystem.service.impl;

import cn.hutool.http.HttpUtil;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.dto.OAuthLoginRequest;
import com.seckill.usersystem.entity.Device;
import com.seckill.usersystem.entity.User;
import com.seckill.usersystem.entity.UserOauth;
import com.seckill.usersystem.entity.UserProfile;
import com.seckill.usersystem.enums.DeviceTypeEnum;
import com.seckill.usersystem.enums.LoginTypeEnum;
import com.seckill.usersystem.exception.BizException;
import com.seckill.usersystem.mapper.DeviceMapper;
import com.seckill.usersystem.mapper.UserMapper;
import com.seckill.usersystem.mapper.UserOauthMapper;
import com.seckill.usersystem.mapper.UserProfileMapper;
import com.seckill.usersystem.service.IOAuthService;
import com.seckill.usersystem.util.JwtUtil;
import com.seckill.usersystem.util.RedisUtil;
import com.seckill.usersystem.dto.TokenResponse;
import com.seckill.usersystem.vo.LoginResultVO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@RequiredArgsConstructor
public class OAuthServiceImpl implements IOAuthService {

    private final UserOauthMapper userOauthMapper;
    private final UserMapper userMapper;
    private final UserProfileMapper userProfileMapper;
    private final DeviceMapper deviceMapper;
    private final JwtUtil jwtUtil;
    private final RedisUtil redisUtil;

    @Value("${oauth.wechat.app-id:}")
    private String wechatAppId;

    @Value("${oauth.wechat.app-secret:}")
    private String wechatAppSecret;

    @Value("${oauth.wechat.authorize-url:https://open.weixin.qq.com/connect/oauth2/authorize}")
    private String wechatAuthorizeUrl;

    @Value("${oauth.wechat.token-url:https://api.weixin.qq.com/sns/oauth2/access_token}")
    private String wechatTokenUrl;

    @Value("${oauth.wechat.userinfo-url:https://api.weixin.qq.com/sns/userinfo}")
    private String wechatUserInfoUrl;

    @Override
    @Transactional(rollbackFor = Exception.class)
    public LoginResultVO oauthLogin(OAuthLoginRequest request) {
        String oauthType = request.getOauthType().toUpperCase();

        cn.hutool.json.JSONObject oauthUserInfo = getOAuthUserInfo(oauthType, request.getCode());
        String openid = oauthUserInfo.getStr("openid");
        String unionid = oauthUserInfo.getStr("unionid");

        UserOauth oauth = QueryChain.of(userOauthMapper)
                .where(UserOauth::getOauthType).eq(oauthType)
                .where(UserOauth::getOpenid).eq(openid)
                .one();

        Long userId;
        UserProfile profile;

        if (oauth == null) {
            User user = new User();
            user.setUsername("oauth_" + oauthType.toLowerCase() + "_" + openid.substring(0, 8));
            user.setPassword("");
            user.setStatus(1);
            user.setIsDeleted(0);
            user.setLoginType(LoginTypeEnum.OAUTH.getCode());
            user.setLoginCount(0);
            user.setPasswordUpdateTime(LocalDateTime.now());
            userMapper.insert(user);

            profile = new UserProfile();
            profile.setUserId(user.getId());
            profile.setNickname(oauthUserInfo.getStr("nickname"));
            profile.setAvatar(oauthUserInfo.getStr("headimgurl"));
            userProfileMapper.insert(profile);

            UserOauth newOauth = new UserOauth();
            newOauth.setUserId(user.getId());
            newOauth.setOauthType(oauthType);
            newOauth.setOpenid(openid);
            newOauth.setUnionid(unionid);
            newOauth.setAccessToken(oauthUserInfo.getStr("access_token"));
            newOauth.setRefreshToken(oauthUserInfo.getStr("refresh_token"));
            newOauth.setExpiresIn(oauthUserInfo.getInt("expires_in"));
            newOauth.setNickname(oauthUserInfo.getStr("nickname"));
            newOauth.setAvatar(oauthUserInfo.getStr("headimgurl"));
            userOauthMapper.insert(newOauth);

            userId = user.getId();
        } else {
            oauth.setAccessToken(oauthUserInfo.getStr("access_token"));
            oauth.setRefreshToken(oauthUserInfo.getStr("refresh_token"));
            oauth.setExpiresIn(oauthUserInfo.getInt("expires_in"));
            oauth.setUpdateTime(LocalDateTime.now());
            userOauthMapper.update(oauth);

            userId = oauth.getUserId();
            profile = QueryChain.of(userProfileMapper)
                    .where(UserProfile::getUserId).eq(userId)
                    .one();
        }

        String deviceId = request.getDeviceId() != null ? request.getDeviceId() : "oauth_" + openid;
        TokenResponse tokenResponse = generateTokens(userId, "oauth_" + openid, deviceId);

        saveOrUpdateDevice(userId, deviceId, request.getDeviceType());

        User user = userMapper.selectOneById(userId);
        user.setLastLoginTime(LocalDateTime.now());
        user.setLoginCount(user.getLoginCount() + 1);
        userMapper.update(user);

        return buildLoginResult(user, profile, tokenResponse);
    }

    @Override
    public String getOAuthAuthorizeUrl(String oauthType, String redirectUri) {
        return switch (oauthType.toUpperCase()) {
            case "WECHAT" -> wechatAuthorizeUrl +
                    "?appid=" + wechatAppId +
                    "&redirect_uri=" + redirectUri +
                    "&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect";
            case "GITHUB" -> "https://github.com/login/oauth/authorize" +
                    "?client_id=" + wechatAppId +
                    "&redirect_uri=" + redirectUri +
                    "&scope=user:email";
            case "GOOGLE" -> "https://accounts.google.com/o/oauth2/v2/auth" +
                    "?client_id=" + wechatAppId +
                    "&redirect_uri=" + redirectUri +
                    "&response_type=code&scope=openid%20email%20profile";
            default -> throw BizException.of("不支持的OAuth类型: " + oauthType);
        };
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void bindOAuth(Long userId, String oauthType, String code) {
        cn.hutool.json.JSONObject oauthUserInfo = getOAuthUserInfo(oauthType.toUpperCase(), code);
        String openid = oauthUserInfo.getStr("openid");

        UserOauth existing = QueryChain.of(userOauthMapper)
                .where(UserOauth::getOauthType).eq(oauthType.toUpperCase())
                .where(UserOauth::getOpenid).eq(openid)
                .one();

        if (existing != null) {
            throw BizException.of("该第三方账号已绑定其他用户");
        }

        UserOauth oauth = new UserOauth();
        oauth.setUserId(userId);
        oauth.setOauthType(oauthType.toUpperCase());
        oauth.setOpenid(openid);
        oauth.setUnionid(oauthUserInfo.getStr("unionid"));
        oauth.setAccessToken(oauthUserInfo.getStr("access_token"));
        oauth.setRefreshToken(oauthUserInfo.getStr("refresh_token"));
        oauth.setExpiresIn(oauthUserInfo.getInt("expires_in"));
        userOauthMapper.insert(oauth);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void unbindOAuth(Long userId, String oauthType) {
        UserOauth oauth = QueryChain.of(userOauthMapper)
                .where(UserOauth::getUserId).eq(userId)
                .where(UserOauth::getOauthType).eq(oauthType.toUpperCase())
                .one();

        if (oauth != null) {
            userOauthMapper.deleteById(oauth.getId());
        }
    }

    private cn.hutool.json.JSONObject getOAuthUserInfo(String oauthType, String code) {
        return switch (oauthType) {
            case "WECHAT" -> getWechatUserInfo(code);
            case "GITHUB" -> {
                JSONObject githubObj = new JSONObject();
                githubObj.set("openid", "github_" + code);
                githubObj.set("nickname", "GitHub User");
                yield githubObj;
            }
            case "GOOGLE" -> {
                JSONObject googleObj = new JSONObject();
                googleObj.set("openid", "google_" + code);
                googleObj.set("nickname", "Google User");
                yield googleObj;
            }
            default -> throw BizException.of("不支持的OAuth类型: " + oauthType);
        };
    }

    private cn.hutool.json.JSONObject getWechatUserInfo(String code) {
        Map<String, Object> tokenParams = new HashMap<>();
        tokenParams.put("appid", wechatAppId);
        tokenParams.put("secret", wechatAppSecret);
        tokenParams.put("code", code);
        tokenParams.put("grant_type", "authorization_code");

        String tokenResult = HttpUtil.get(wechatTokenUrl, tokenParams);
        JSONObject tokenJson = JSONUtil.parseObj(tokenResult);

        if (tokenJson.containsKey("errcode")) {
            throw BizException.of("微信授权失败: " + tokenJson.getStr("errmsg"));
        }

        String accessToken = tokenJson.getStr("access_token");
        String openid = tokenJson.getStr("openid");

        String userInfoUrl = wechatUserInfoUrl +
                "?access_token=" + accessToken +
                "&openid=" + openid +
                "&lang=zh_CN";

        String userInfo = HttpUtil.get(userInfoUrl);
        return JSONUtil.parseObj(userInfo);
    }

    private TokenResponse generateTokens(Long userId, String username, String deviceId) {
        String accessToken = jwtUtil.generateAccessToken(userId, username, deviceId);
        String refreshToken = jwtUtil.generateRefreshToken(userId, deviceId);

        redisUtil.set("user:token:" + userId + ":" + deviceId, accessToken,
                jwtUtil.getAccessTokenExpiration(), TimeUnit.SECONDS);
        redisUtil.set("user:refresh:" + userId + ":" + deviceId, refreshToken,
                jwtUtil.getRefreshTokenExpiration(), TimeUnit.SECONDS);

        TokenResponse response = new TokenResponse();
        response.setAccessToken(accessToken);
        response.setRefreshToken(refreshToken);
        response.setExpiresIn(jwtUtil.getAccessTokenExpiration());
        return response;
    }

    private void saveOrUpdateDevice(Long userId, String deviceId, String deviceType) {
        Device device = QueryChain.of(deviceMapper)
                .where(Device::getUserId).eq(userId)
                .where(Device::getDeviceId).eq(deviceId)
                .one();

        if (device == null) {
            device = new Device();
            device.setUserId(userId);
            device.setDeviceId(deviceId);
            device.setDeviceType(deviceType != null ? deviceType : DeviceTypeEnum.WEB.getCode());
            device.setDeviceName("OAuth设备");
            device.setLastLoginTime(LocalDateTime.now());
            device.setLastActiveTime(LocalDateTime.now());
            device.setStatus(1);
            deviceMapper.insert(device);
        } else {
            device.setLastLoginTime(LocalDateTime.now());
            device.setLastActiveTime(LocalDateTime.now());
            deviceMapper.update(device);
        }
    }

    private LoginResultVO buildLoginResult(User user, UserProfile profile, TokenResponse token) {
        LoginResultVO result = new LoginResultVO();
        result.setUserId(user.getId());
        result.setUsername(user.getUsername());
        result.setNickname(profile != null ? profile.getNickname() : user.getUsername());
        result.setAvatar(profile != null ? profile.getAvatar() : null);
        result.setToken(token);
        result.setIsFirstLogin(user.getLoginCount() <= 1);
        result.setNeedChangePassword(false);
        return result;
    }
}
