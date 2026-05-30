package com.seckill.seckill.limiter;

import com.seckill.seckill.dto.SeckillResult;
import com.seckill.usersystem.vo.Result;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.stereotype.Component;

@Slf4j
@Aspect
@Component
@RequiredArgsConstructor
public class RateLimitAspect {

    private final TokenBucketLimiter tokenBucketLimiter;

    @Around("@annotation(rateLimit)")
    public Object around(ProceedingJoinPoint joinPoint, RateLimit rateLimit) throws Throwable {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        String key = signature.getDeclaringTypeName() + "." + signature.getName();

        boolean acquired = tokenBucketLimiter.tryAcquireByApi(
                key, rateLimit.permitsPerSecond(), rateLimit.timeoutMillis());

        if (!acquired) {
            log.warn("[RateLimit] API:{} blocked", key);
            return Result.fail(429, "请求过于频繁，请稍后再试");
        }

        return joinPoint.proceed();
    }
}