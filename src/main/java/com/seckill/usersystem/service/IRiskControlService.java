package com.seckill.usersystem.service;

import com.seckill.usersystem.enums.RiskLevelEnum;
import com.seckill.usersystem.vo.LoginResultVO;

/**
 * 风控服务接口
 */
public interface IRiskControlService {

    RiskLevelEnum assessLoginRisk(Long userId, String ip, String deviceId, String userAgent);

    boolean isLoginAllowed(RiskLevelEnum riskLevel);

    void recordLoginAttempt(Long userId, String ip, boolean success);

    int getFailedLoginCount(Long userId, String ip);

    void clearFailedLoginCount(Long userId, String ip);

    boolean isAccountLocked(Long userId);

    void triggerRiskAction(Long userId, RiskLevelEnum riskLevel);

    void recordDeviceChange(Long userId, String oldDeviceId, String newDeviceId, String ip);
}
