package com.seckill.usersystem.service;

import com.seckill.usersystem.dto.OAuthLoginRequest;
import com.seckill.usersystem.vo.LoginResultVO;

/**
 * OAuth登录服务接口
 */
public interface IOAuthService {

    LoginResultVO oauthLogin(OAuthLoginRequest request);

    String getOAuthAuthorizeUrl(String oauthType, String redirectUri);

    void bindOAuth(Long userId, String oauthType, String code);

    void unbindOAuth(Long userId, String oauthType);
}
