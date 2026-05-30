package com.seckill.usersystem.service;

import com.seckill.usersystem.entity.Permission;

import java.util.List;

public interface IPermissionService {

    List<Permission> listPermissions();

    List<Permission> getPermissionsByUserId(Long userId);

    List<Permission> getPermissionsByRoleId(Long roleId);

    void savePermission(Permission permission);

    void updatePermission(Permission permission);

    void deletePermission(Long permissionId);

    List<String> getPermissionCodesByUserId(Long userId);
}
