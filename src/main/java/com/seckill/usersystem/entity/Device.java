package com.seckill.usersystem.entity;

import com.mybatisflex.annotation.Column;
import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@Table("sys_device")
public class Device implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long userId;

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

    private Integer isTrusted;

    @Column(onInsertValue = "now()")
    private LocalDateTime createTime;

    @Column(onInsertValue = "now()", onUpdateValue = "now()")
    private LocalDateTime updateTime;
}
