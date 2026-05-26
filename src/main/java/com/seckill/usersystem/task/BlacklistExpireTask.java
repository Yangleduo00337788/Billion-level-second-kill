package com.seckill.usersystem.task;

import com.seckill.usersystem.service.IBlacklistService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
public class BlacklistExpireTask {

    private final IBlacklistService blacklistService;
    private final RedisUtil redisUtil;

    @Scheduled(cron = "0 0 * * * ?")
    public void expireBlacklist() {
        log.info("开始执行黑名单过期检查任务");
        blacklistService.checkAndExpireBlacklist();
    }

    @Scheduled(cron = "0 0 */6 * * ?")
    public void cleanExpiredTokens() {
        log.info("开始清理过期令牌缓存");
        redisUtil.deleteByPattern("token:blacklist:*");
    }
}
