package com.seckill.seckill.limiter;

import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RRateLimiter;
import org.redisson.api.RateIntervalUnit;
import org.redisson.api.RateType;
import org.redisson.api.RedissonClient;
import org.springframework.stereotype.Component;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

@Slf4j
@Component
@RequiredArgsConstructor
public class TokenBucketLimiter {

    private final RedissonClient redissonClient;

    private final ConcurrentHashMap<String, RRateLimiter> limiterCache = new ConcurrentHashMap<>();

    public boolean tryAcquire(String key, long permitsPerSecond, long timeoutMillis) {
        RRateLimiter limiter = limiterCache.computeIfAbsent(key, k -> {
            RRateLimiter r = redissonClient.getRateLimiter("rate:limiter:" + k);
            r.trySetRate(RateType.OVERALL, permitsPerSecond, 1, RateIntervalUnit.SECONDS);
            return r;
        });

        try {
            return limiter.tryAcquire(timeoutMillis, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            log.warn("[TokenBucketLimiter] tryAcquire error, key:{}, cause:{}", key, e.getMessage());
            return true;
        }
    }

    public boolean tryAcquireByUser(Long userId, long permitsPerSecond) {
        return tryAcquire("user:" + userId, permitsPerSecond, 100);
    }

    public boolean tryAcquireByApi(String api, long permitsPerSecond, long timeoutMillis) {
        return tryAcquire("api:" + api, permitsPerSecond, timeoutMillis);
    }
}