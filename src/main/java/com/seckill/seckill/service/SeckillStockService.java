package com.seckill.seckill.service;

import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.io.ClassPathResource;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.core.script.DefaultRedisScript;
import org.springframework.scripting.support.ResourceScriptSource;
import org.springframework.stereotype.Service;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@RequiredArgsConstructor
public class SeckillStockService {

    private final StringRedisTemplate redisTemplate;

    private DefaultRedisScript<Long> deductScript;
    private DefaultRedisScript<Long> rollbackScript;

    private static final String STOCK_KEY_PREFIX = "seckill:stock:";
    private static final String ORDER_FLAG_PREFIX = "seckill:order:user:";
    private static final String SUCCESS_USERS_PREFIX = "seckill:success:users:";

    @PostConstruct
    public void initScripts() {
        deductScript = new DefaultRedisScript<>();
        deductScript.setScriptSource(new ResourceScriptSource(
                new ClassPathResource("scripts/seckill_deduct.lua")));
        deductScript.setResultType(Long.class);

        rollbackScript = new DefaultRedisScript<>();
        rollbackScript.setScriptSource(new ResourceScriptSource(
                new ClassPathResource("scripts/seckill_rollback.lua")));
        rollbackScript.setResultType(Long.class);

        log.info("[SeckillStockService] Lua scripts loaded successfully");
    }

    public void warmUpStock(Long seckillId, int stock) {
        String stockKey = getStockKey(seckillId);
        redisTemplate.opsForValue().set(stockKey, String.valueOf(stock));
        log.info("[SeckillStockService] stock warmed up: seckillId={}, stock={}", seckillId, stock);
    }

    public Long deductStock(Long seckillId, Long userId, int personLimit) {
        String stockKey = getStockKey(seckillId);
        String orderFlagKey = getOrderFlagKey(userId, seckillId);

        List<String> keys = Arrays.asList(stockKey, orderFlagKey);
        Long result = redisTemplate.execute(
                deductScript,
                keys,
                String.valueOf(userId),
                String.valueOf(seckillId),
                String.valueOf(personLimit)
        );
        return result;
    }

    public Long rollbackStock(Long seckillId, Long userId) {
        String stockKey = getStockKey(seckillId);
        String orderFlagKey = getOrderFlagKey(userId, seckillId);

        List<String> keys = Arrays.asList(stockKey, orderFlagKey);
        return redisTemplate.execute(rollbackScript, keys);
    }

    public Integer getRemainingStock(Long seckillId) {
        String val = redisTemplate.opsForValue().get(getStockKey(seckillId));
        return val != null ? Integer.parseInt(val) : 0;
    }

    public boolean hasOrderFlag(Long userId, Long seckillId) {
        return Boolean.TRUE.equals(
                redisTemplate.hasKey(getOrderFlagKey(userId, seckillId)));
    }

    public String getStockKey(Long seckillId) {
        return STOCK_KEY_PREFIX + seckillId;
    }

    public String getOrderFlagKey(Long userId, Long seckillId) {
        return ORDER_FLAG_PREFIX + userId + ":" + seckillId;
    }

    public Long getSuccessUserCount(Long seckillId) {
        return redisTemplate.opsForSet().size(SUCCESS_USERS_PREFIX + seckillId);
    }
}