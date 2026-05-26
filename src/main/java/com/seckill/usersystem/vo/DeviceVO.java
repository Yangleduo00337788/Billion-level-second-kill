package com.seckill.usersystem.vo;

import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;
import java.util.List;

/**
 * 设备信息视图对象
 */
@Data
public class DeviceVO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private Long id;

    private String deviceId;

    private String deviceType;

    private String deviceName;

    private String osName;

    private String browser;

    private String ipAddress;

    private String location;

    private LocalDateTime lastLoginTime;

    private LocalDateTime lastActiveTime;

    private Integer status;

    private Boolean isTrusted;

    private Boolean isCurrentDevice;
}
