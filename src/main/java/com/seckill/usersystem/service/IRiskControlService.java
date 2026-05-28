package com.seckill.usersystem.service;

import com.seckill.usersystem.enums.RiskLevelEnum;
import com.seckill.usersystem.vo.LoginResultVO;

import java.util.concurrent.TimeUnit;

/**
 * 风控服务接口
 */
public interface IRiskControlService {

    RiskLevelEnum assessLoginRisk(Long userId, String ip, String deviceId, String userAgent);

    boolean isLoginAllowed(RiskLevelEnum riskLevel);

    boolean isCaptchaRequired(RiskLevelEnum riskLevel);

    boolean isSmsCodeRequired(RiskLevelEnum riskLevel);

    void requireCaptcha(String ip, String deviceId);

    boolean isCaptchaRequired(String ip, String deviceId);

    void clearCaptchaRequired(String ip, String deviceId);

    void recordLoginAttempt(Long userId, String ip, boolean success);

    int getFailedLoginCount(Long userId, String ip);

    void clearFailedLoginCount(Long userId, String ip);

    boolean isAccountLocked(Long userId);

    boolean isIpBlocked(String ip);

    void blockIp(String ip, long duration, TimeUnit unit, String reason);

    void triggerRiskAction(Long userId, RiskLevelEnum riskLevel);

    void recordDeviceChange(Long userId, String oldDeviceId, String newDeviceId, String ip);

    void recordLoginLocation(Long userId, String ip, String location);

    String getLastLoginLocation(Long userId);
}
