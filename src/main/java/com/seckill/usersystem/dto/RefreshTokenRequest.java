package com.seckill.usersystem.dto;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 刷新令牌请求DTO
 */
@Data
public class RefreshTokenRequest implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private String refreshToken;

    private String deviceId;
}
