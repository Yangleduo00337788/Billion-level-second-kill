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
@Table("sys_login_log")
public class LoginLog implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long userId;

    private String username;

    private Integer loginType;

    private String deviceType;

    private String deviceId;

    private String ipAddress;

    private String location;

    private String browser;

    private String osName;

    private Integer status;

    private String failReason;

    private String userAgent;

    @Column(onInsertValue = "now()")
    private LocalDateTime loginTime;
}
