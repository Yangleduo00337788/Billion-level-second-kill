package com.seckill.usersystem.service;

import com.seckill.usersystem.entity.Blacklist;

import java.util.List;

public interface IBlacklistService {

    void addToBlacklist(Blacklist blacklist);

    void removeUserFromBlacklist(Long blacklistId);

    boolean isUserInBlacklist(Long userId);

    boolean isIpInBlacklist(String ipAddress);

    boolean isDeviceInBlacklist(String deviceId);

    List<Blacklist> getActiveBlacklists();

    void checkAndExpireBlacklist();

    List<Blacklist> getBlacklistByUserId(Long userId);
}
