package com.seckill.seckill.limiter;

import java.lang.annotation.*;

@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface RateLimit {

    long permitsPerSecond() default 100;

    long timeoutMillis() default 500;
}