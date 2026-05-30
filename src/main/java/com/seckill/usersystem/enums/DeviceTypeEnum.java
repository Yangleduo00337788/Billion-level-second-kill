package com.seckill.usersystem.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * 设备类型枚举
 */
@Getter
@AllArgsConstructor
public enum DeviceTypeEnum {

    WEB("WEB", "网页端"),
    IOS("IOS", "苹果移动端"),
    ANDROID("ANDROID", "安卓移动端"),
    PC("PC", "Windows客户端"),
    MAC("MAC", "Mac客户端"),
    MINI_PROGRAM("MINI_PROGRAM", "小程序");

    private final String code;
    private final String description;

    public static DeviceTypeEnum fromCode(String code) {
        for (DeviceTypeEnum type : values()) {
            if (type.code.equals(code)) {
                return type;
            }
        }
        return WEB;
    }
}
