package com.seckill.usersystem.vo;

import com.seckill.usersystem.dto.TokenResponse;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;

/**
 * 登录结果视图对象
 */
@Data
public class LoginResultVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private Long userId;

    private String username;

    private String nickname;

    private String avatar;

    private TokenResponse token;

    private Boolean isFirstLogin;

    private Boolean needChangePassword;

    private Integer maxDeviceCount;

    private Integer currentDeviceCount;
}
