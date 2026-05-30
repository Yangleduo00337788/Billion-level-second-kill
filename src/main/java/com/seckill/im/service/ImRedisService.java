package com.seckill.im.service;

import com.google.protobuf.InvalidProtocolBufferException;
import com.seckill.im.proto.ImMessageProto;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.concurrent.TimeUnit;

@Slf4j
@Service
@RequiredArgsConstructor
public class ImRedisService {

    private final StringRedisTemplate redisTemplate;

    private static final String REDIS_KEY_MSG_IDEMPOTENT = "im:msg:idempotent:%s";
    private static final String REDIS_KEY_SESSION_SEQ = "im:session:seq:%d";
    private static final String REDIS_KEY_USER_ONLINE = "im:user:online:%d";
    private static final String REDIS_KEY_USER_OFFLINE_QUEUE = "im:user:offline:queue:%d";
    private static final String REDIS_KEY_UNREAD_COUNT = "im:user:unread:%d:%d";
    private static final String REDIS_KEY_LAST_READ_SEQ = "im:user:lastread:%d:%d";

    private static final int IDEMPOTENT_TTL_HOURS = 24;
    private static final int ONLINE_TTL_MINUTES = 5;
    private static final int OFFLINE_QUEUE_TTL_DAYS = 7;

    public boolean checkAndMarkIdempotent(String clientMsgId) {
        String key = String.format(REDIS_KEY_MSG_IDEMPOTENT, clientMsgId);
        Boolean success = redisTemplate.opsForValue()
                .setIfAbsent(key, "1", IDEMPOTENT_TTL_HOURS, TimeUnit.HOURS);
        return Boolean.TRUE.equals(success);
    }

    public boolean hasProcessed(String clientMsgId) {
        String key = String.format(REDIS_KEY_MSG_IDEMPOTENT, clientMsgId);
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }

    public Long generateSessionSeq(Long sessionId) {
        String key = String.format(REDIS_KEY_SESSION_SEQ, sessionId);
        return redisTemplate.opsForValue().increment(key);
    }

    public Long getCurrentSeq(Long sessionId) {
        String key = String.format(REDIS_KEY_SESSION_SEQ, sessionId);
        String val = redisTemplate.opsForValue().get(key);
        return val != null ? Long.parseLong(val) : 0L;
    }

    public void markUserOnline(Long userId, String deviceId, String deviceType) {
        String key = String.format(REDIS_KEY_USER_ONLINE, userId);
        redisTemplate.opsForHash().put(key, deviceId, deviceType);
        redisTemplate.expire(key, ONLINE_TTL_MINUTES, TimeUnit.MINUTES);
    }

    public void refreshUserOnline(Long userId) {
        String key = String.format(REDIS_KEY_USER_ONLINE, userId);
        if (Boolean.TRUE.equals(redisTemplate.hasKey(key))) {
            redisTemplate.expire(key, ONLINE_TTL_MINUTES, TimeUnit.MINUTES);
        }
    }

    public void markUserOffline(Long userId) {
        String key = String.format(REDIS_KEY_USER_ONLINE, userId);
        redisTemplate.delete(key);
    }

    public boolean isUserOnline(Long userId) {
        String key = String.format(REDIS_KEY_USER_ONLINE, userId);
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }

    public void addOfflineMessage(Long userId, ImMessageProto.DownstreamMessage msg) {
        String key = String.format(REDIS_KEY_USER_OFFLINE_QUEUE, userId);
        double score = (double) msg.getServerTime();
        redisTemplate.opsForZSet().add(key, msg.toByteArray().toString(), score);
        redisTemplate.expire(key, OFFLINE_QUEUE_TTL_DAYS, TimeUnit.DAYS);
    }

    public List<ImMessageProto.DownstreamMessage> getAndRemoveOfflineMessages(Long userId, int limit) {
        String key = String.format(REDIS_KEY_USER_OFFLINE_QUEUE, userId);
        Set<String> rawSet = redisTemplate.opsForZSet().range(key, 0, limit - 1);
        if (rawSet == null || rawSet.isEmpty()) {
            return List.of();
        }
        List<ImMessageProto.DownstreamMessage> messages = new ArrayList<>();
        for (String raw : rawSet) {
            try {
                ImMessageProto.DownstreamMessage msg = ImMessageProto.DownstreamMessage.parseFrom(raw.getBytes());
                messages.add(msg);
            } catch (Exception e) {
                log.warn("[ImRedisService] failed to parse offline message for user:{}", userId);
            }
        }
        redisTemplate.opsForZSet().removeRange(key, 0, limit - 1);
        return messages;
    }

    public int getUnreadCount(Long userId) {
        int count = 0;
        Set<String> keys = redisTemplate.keys(String.format(REDIS_KEY_UNREAD_COUNT, userId, 0).replace("0", "*"));
        if (keys != null) {
            for (String key : keys) {
                String val = redisTemplate.opsForValue().get(key);
                if (val != null) {
                    count += Integer.parseInt(val);
                }
            }
        }
        return count;
    }

    public void incrementUnreadCount(Long userId, Long sessionId) {
        String key = String.format(REDIS_KEY_UNREAD_COUNT, userId, sessionId);
        redisTemplate.opsForValue().increment(key);
    }

    public void resetUnreadCount(Long userId, Long sessionId) {
        String key = String.format(REDIS_KEY_UNREAD_COUNT, userId, sessionId);
        redisTemplate.delete(key);
    }

    public void updateLastReadSeq(Long userId, Long sessionId, Long seq) {
        String key = String.format(REDIS_KEY_LAST_READ_SEQ, userId, sessionId);
        redisTemplate.opsForValue().set(key, String.valueOf(seq));
    }

    public Long getLastReadSeq(Long userId, Long sessionId) {
        String key = String.format(REDIS_KEY_LAST_READ_SEQ, userId, sessionId);
        String val = redisTemplate.opsForValue().get(key);
        return val != null ? Long.parseLong(val) : 0L;
    }
}