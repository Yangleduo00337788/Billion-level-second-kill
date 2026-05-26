package com.seckill.usersystem.dto;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * OAuth登录请求DTO
 */
@Data
public class OAuthLoginRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private String oauthType;

    private String code;

    private String state;

    private String redirectUri;

    private String deviceId;

    private String deviceType;
}
