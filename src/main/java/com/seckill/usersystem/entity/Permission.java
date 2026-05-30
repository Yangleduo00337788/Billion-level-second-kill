package com.seckill.usersystem.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@Table("sys_permission")
public class Permission implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long parentId;

    private String permissionCode;

    private String permissionName;

    /** 权限类型：1-菜单，2-按钮，3-接口 */
    private Integer permissionType;

    private String path;

    private String method;

    private String icon;

    private Integer sort;

    private Integer status;

    private Integer isDeleted;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}
