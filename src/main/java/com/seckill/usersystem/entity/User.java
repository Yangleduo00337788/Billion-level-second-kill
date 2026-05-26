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
@Table("sys_user")
public class User implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private String username;

    private String password;

    private String phone;

    private String email;

    private Integer status;

    private Integer isDeleted;

    private Integer loginType;

    private LocalDateTime lastLoginTime;

    private String lastLoginIp;

    private Integer loginCount;

    private LocalDateTime passwordUpdateTime;

    @Column(onInsertValue = "now()")
    private LocalDateTime createTime;

    @Column(onInsertValue = "now()", onUpdateValue = "now()")
    private LocalDateTime updateTime;
}
