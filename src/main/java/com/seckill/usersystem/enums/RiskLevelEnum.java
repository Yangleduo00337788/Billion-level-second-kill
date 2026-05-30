package com.seckill.usersystem.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * 风控等级枚举
 */
@Getter
@AllArgsConstructor
public enum RiskLevelEnum {

    LOW(1, "低风险"),
    MEDIUM(2, "中风险"),
    HIGH(3, "高风险"),
    CRITICAL(4, "极高风险");

    private final int code;
    private final String description;

    public static RiskLevelEnum fromCode(int code) {
        for (RiskLevelEnum level : values()) {
            if (level.code == code) {
                return level;
            }
        }
        return LOW;
    }
}
