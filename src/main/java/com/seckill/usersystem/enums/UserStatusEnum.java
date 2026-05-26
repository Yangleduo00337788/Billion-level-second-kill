package com.seckill.usersystem.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * 用户状态枚举
 */
@Getter
@AllArgsConstructor
public enum UserStatusEnum {

    DISABLED(0, "禁用"),
    ENABLED(1, "启用"),
    LOCKED(2, "锁定");

    private final int code;
    private final String description;

    public static UserStatusEnum fromCode(int code) {
        for (UserStatusEnum status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return DISABLED;
    }
}
