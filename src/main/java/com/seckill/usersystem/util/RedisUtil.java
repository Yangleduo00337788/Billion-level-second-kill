package com.seckill.usersystem.util;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.redisson.api.RBucket;
import org.redisson.api.RLock;
import org.redisson.api.RSet;
import org.redisson.api.RedissonClient;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;

/**
 * Redis 工具类
 * 基于 Redisson + StringRedisTemplate 实现
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class RedisUtil {

    private final StringRedisTemplate stringRedisTemplate;

    private final RedissonClient redissonClient;

    /**
     * 设置字符串值
     */
    public void set(String key, String value) {
        stringRedisTemplate.opsForValue().set(key, value);
    }

    /**
     * 设置字符串值并指定过期时间（秒）
     */
    public void set(String key, String value, long timeout) {
        stringRedisTemplate.opsForValue().set(key, value, timeout, TimeUnit.SECONDS);
    }

    /**
     * 设置字符串值并指定过期时间和时间单位
     */
    public void set(String key, String value, long timeout, TimeUnit unit) {
        stringRedisTemplate.opsForValue().set(key, value, timeout, unit);
    }

    /**
     * 获取字符串值
     */
    public String get(String key) {
        return stringRedisTemplate.opsForValue().get(key);
    }

    /**
     * 删除键
     */
    public Boolean delete(String key) {
        return stringRedisTemplate.delete(key);
    }

    /**
     * 批量删除键
     */
    public Long delete(String... keys) {
        List<String> keyList = Arrays.asList(keys);
        return stringRedisTemplate.delete(keyList);
    }

    /**
     * 删除匹配模式的键
     */
    public Long deleteByPattern(String pattern) {
        Set<String> keys = stringRedisTemplate.keys(pattern);
        if (keys == null || keys.isEmpty()) {
            return 0L;
        }
        return stringRedisTemplate.delete(keys);
    }

    /**
     * 判断键是否存在
     */
    public Boolean hasKey(String key) {
        return stringRedisTemplate.hasKey(key);
    }

    /**
     * 设置过期时间（秒）
     */
    public Boolean expire(String key, long timeout) {
        return stringRedisTemplate.expire(key, timeout, TimeUnit.SECONDS);
    }

    /**
     * 获取过期时间（秒）
     */
    public Long getExpire(String key) {
        return stringRedisTemplate.getExpire(key, TimeUnit.SECONDS);
    }

    /**
     * 自增
     */
    public Long increment(String key) {
        return stringRedisTemplate.opsForValue().increment(key);
    }

    /**
     * 自增指定值
     */
    public Long increment(String key, long delta) {
        return stringRedisTemplate.opsForValue().increment(key, delta);
    }

    /**
     * 设置哈希值
     */
    public void hashSet(String key, String field, String value) {
        stringRedisTemplate.opsForHash().put(key, field, value);
    }

    /**
     * 获取哈希值
     */
    public String hashGet(String key, String field) {
        Object value = stringRedisTemplate.opsForHash().get(key, field);
        return value != null ? value.toString() : null;
    }

    /**
     * 获取整个哈希
     */
    public Map<Object, Object> hashGetAll(String key) {
        return stringRedisTemplate.opsForHash().entries(key);
    }

    /**
     * 向集合添加元素
     */
    public Long setAdd(String key, String... values) {
        return stringRedisTemplate.opsForSet().add(key, values);
    }

    /**
     * 获取集合所有元素
     */
    public Set<String> setMembers(String key) {
        return stringRedisTemplate.opsForSet().members(key);
    }

    /**
     * 判断元素是否在集合中
     */
    public Boolean setIsMember(String key, String value) {
        return stringRedisTemplate.opsForSet().isMember(key, value);
    }

    /**
     * 从集合移除元素
     */
    public Long setRemove(String key, Object... values) {
        return stringRedisTemplate.opsForSet().remove(key, values);
    }

    /**
     * 向列表右侧推入元素
     */
    public Long rightPush(String key, String value) {
        return stringRedisTemplate.opsForList().rightPush(key, value);
    }

    /**
     * 从列表左侧弹出元素
     */
    public String leftPop(String key) {
        return stringRedisTemplate.opsForList().leftPop(key);
    }

    /**
     * 获取列表长度
     */
    public Long listSize(String key) {
        return stringRedisTemplate.opsForList().size(key);
    }

    /**
     * 获取列表指定范围元素
     */
    public List<String> listRange(String key, long start, long end) {
        return stringRedisTemplate.opsForList().range(key, start, end);
    }

    /**
     * 获取分布式锁
     */
    public RLock getLock(String lockKey) {
        return redissonClient.getLock(lockKey);
    }

    /**
     * 尝试获取锁
     */
    public boolean tryLock(String lockKey, long waitTime, long leaseTime) {
        RLock lock = redissonClient.getLock(lockKey);
        try {
            return lock.tryLock(waitTime, leaseTime, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("获取分布式锁被中断: {}", lockKey);
            return false;
        }
    }

    /**
     * 释放锁
     */
    public void unlock(String lockKey) {
        RLock lock = redissonClient.getLock(lockKey);
        if (lock.isHeldByCurrentThread()) {
            lock.unlock();
        }
    }

    /**
     * 获取RBucket对象（用于更细粒度的操作）
     */
    public <T> RBucket<T> getBucket(String key) {
        return redissonClient.getBucket(key);
    }

    /**
     * 获取RSet对象
     */
    public <T> RSet<T> getSet(String key) {
        return redissonClient.getSet(key);
    }

    /**
     * 检查键是否存在且值匹配
     */
    public boolean matchValue(String key, String expectedValue) {
        String actualValue = get(key);
        return expectedValue.equals(actualValue);
    }
}
