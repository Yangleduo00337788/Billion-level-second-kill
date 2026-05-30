package com.seckill.usersystem.exception;

import lombok.Getter;

/**
 * 业务异常类
 */
@Getter
public class BizException extends RuntimeException {

    private final Integer code;

    public BizException(String message) {
        super(message);
        this.code = 500;
    }

    public BizException(Integer code, String message) {
        super(message);
        this.code = code;
    }

    public static BizException of(String message) {
        return new BizException(message);
    }

    public static BizException of(Integer code, String message) {
        return new BizException(code, message);
    }
}
