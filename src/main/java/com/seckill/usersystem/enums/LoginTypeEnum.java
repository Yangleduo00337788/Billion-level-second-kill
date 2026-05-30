package com.seckill.usersystem.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * 登录类型枚举
 */
@Getter
@AllArgsConstructor
public enum LoginTypeEnum {

    PASSWORD(1, "密码登录"),
    SMS(2, "短信验证码登录"),
    OAUTH(3, "第三方OAuth登录");

    private final int code;
    private final String description;

    public static LoginTypeEnum fromCode(int code) {
        for (LoginTypeEnum type : values()) {
            if (type.code == code) {
                return type;
            }
        }
        return PASSWORD;
    }
}
