package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.entity.Blacklist;
import com.seckill.usersystem.mapper.BlacklistMapper;
import com.seckill.usersystem.service.IBlacklistService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class BlacklistServiceImpl implements IBlacklistService {

    private final BlacklistMapper blacklistMapper;
    private final RedisUtil redisUtil;

    private static final String BLACKLIST_USER_PREFIX = "blacklist:user:";
    private static final String BLACKLIST_IP_PREFIX = "blacklist:ip:";
    private static final String BLACKLIST_DEVICE_PREFIX = "blacklist:device:";

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void addToBlacklist(Blacklist blacklist) {
        blacklist.setStatus(1);
        blacklist.setCreateTime(LocalDateTime.now());
        blacklistMapper.insert(blacklist);

        cacheBlacklist(blacklist);
        log.info("添加到黑名单: type={}, target={}", blacklist.getBlacklistType(), getTargetInfo(blacklist));
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void removeUserFromBlacklist(Long blacklistId) {
        Blacklist blacklist = blacklistMapper.selectOneById(blacklistId);
        if (blacklist != null) {
            blacklist.setStatus(0);
            blacklist.setUpdateTime(LocalDateTime.now());
            blacklistMapper.update(blacklist);

            removeCacheBlacklist(blacklist);
            log.info("从黑名单移除: blacklistId={}", blacklistId);
        }
    }

    @Override
    public boolean isUserInBlacklist(Long userId) {
        String key = BLACKLIST_USER_PREFIX + userId;
        if (redisUtil.hasKey(key)) {
            return true;
        }

        Blacklist blacklist = QueryChain.of(blacklistMapper)
                .where(Blacklist::getUserId).eq(userId)
                .where(Blacklist::getStatus).eq(1)
                .one();

        if (blacklist != null) {
            if (blacklist.getExpireTime() == null || blacklist.getExpireTime().isAfter(LocalDateTime.now())) {
                cacheBlacklist(blacklist);
                return true;
            }
        }
        return false;
    }

    @Override
    public boolean isIpInBlacklist(String ipAddress) {
        String key = BLACKLIST_IP_PREFIX + ipAddress;
        if (redisUtil.hasKey(key)) {
            return true;
        }

        Blacklist blacklist = QueryChain.of(blacklistMapper)
                .where(Blacklist::getIpAddress).eq(ipAddress)
                .where(Blacklist::getStatus).eq(1)
                .one();

        if (blacklist != null) {
            if (blacklist.getExpireTime() == null || blacklist.getExpireTime().isAfter(LocalDateTime.now())) {
                cacheBlacklist(blacklist);
                return true;
            }
        }
        return false;
    }

    @Override
    public boolean isDeviceInBlacklist(String deviceId) {
        String key = BLACKLIST_DEVICE_PREFIX + deviceId;
        if (redisUtil.hasKey(key)) {
            return true;
        }

        Blacklist blacklist = QueryChain.of(blacklistMapper)
                .where(Blacklist::getDeviceId).eq(deviceId)
                .where(Blacklist::getStatus).eq(1)
                .one();

        if (blacklist != null) {
            if (blacklist.getExpireTime() == null || blacklist.getExpireTime().isAfter(LocalDateTime.now())) {
                cacheBlacklist(blacklist);
                return true;
            }
        }
        return false;
    }

    @Override
    public List<Blacklist> getActiveBlacklists() {
        return QueryChain.of(blacklistMapper)
                .where(Blacklist::getStatus).eq(1)
                .list();
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void checkAndExpireBlacklist() {
        LocalDateTime now = LocalDateTime.now();
        List<Blacklist> expiredBlacklists = QueryChain.of(blacklistMapper)
                .where(Blacklist::getStatus).eq(1)
                .where(Blacklist::getExpireTime).le(now)
                .list();

        for (Blacklist blacklist : expiredBlacklists) {
            blacklist.setStatus(0);
            blacklist.setUpdateTime(now);
            blacklistMapper.update(blacklist);
            removeCacheBlacklist(blacklist);
        }

        if (!expiredBlacklists.isEmpty()) {
            log.info("过期黑名单清理: count={}", expiredBlacklists.size());
        }
    }

    @Override
    public List<Blacklist> getBlacklistByUserId(Long userId) {
        return QueryChain.of(blacklistMapper)
                .where(Blacklist::getUserId).eq(userId)
                .list();
    }

    private void cacheBlacklist(Blacklist blacklist) {
        if (blacklist.getUserId() != null) {
            String key = BLACKLIST_USER_PREFIX + blacklist.getUserId();
            if (blacklist.getExpireTime() != null) {
                long seconds = java.time.Duration.between(LocalDateTime.now(), blacklist.getExpireTime()).getSeconds();
                if (seconds > 0) {
                    redisUtil.set(key, "1", seconds);
                }
            } else {
                redisUtil.set(key, "1");
            }
        }
        if (blacklist.getIpAddress() != null) {
            String key = BLACKLIST_IP_PREFIX + blacklist.getIpAddress();
            redisUtil.set(key, "1");
        }
        if (blacklist.getDeviceId() != null) {
            String key = BLACKLIST_DEVICE_PREFIX + blacklist.getDeviceId();
            redisUtil.set(key, "1");
        }
    }

    private void removeCacheBlacklist(Blacklist blacklist) {
        if (blacklist.getUserId() != null) {
            redisUtil.delete(BLACKLIST_USER_PREFIX + blacklist.getUserId());
        }
        if (blacklist.getIpAddress() != null) {
            redisUtil.delete(BLACKLIST_IP_PREFIX + blacklist.getIpAddress());
        }
        if (blacklist.getDeviceId() != null) {
            redisUtil.delete(BLACKLIST_DEVICE_PREFIX + blacklist.getDeviceId());
        }
    }

    private String getTargetInfo(Blacklist blacklist) {
        return switch (blacklist.getBlacklistType()) {
            case 1 -> "userId=" + blacklist.getUserId();
            case 2 -> "ip=" + blacklist.getIpAddress();
            case 3 -> "device=" + blacklist.getDeviceId();
            default -> "unknown";
        };
    }
}
