package com.seckill.usersystem.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@Table("sys_role")
public class Role implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private String roleCode;

    private String roleName;

    private String description;

    private Integer status;

    private Integer isDeleted;

    private Integer sort;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}
