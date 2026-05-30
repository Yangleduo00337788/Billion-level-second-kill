package com.seckill.usersystem.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@Table("sys_blacklist")
public class Blacklist implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long userId;

    private String ipAddress;

    private String deviceId;

    /** 黑名单类型：1-用户，2-IP，3-设备 */
    private Integer blacklistType;

    private String reason;

    private LocalDateTime expireTime;

    /** 状态：0-已解封，1-封禁中 */
    private Integer status;

    private Long operatorId;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}
