package com.seckill.usersystem.dto;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 令牌响应DTO
 */
@Data
public class TokenResponse implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private String accessToken;

    private String refreshToken;

    private Long expiresIn;

    private String tokenType = "Bearer";
}
