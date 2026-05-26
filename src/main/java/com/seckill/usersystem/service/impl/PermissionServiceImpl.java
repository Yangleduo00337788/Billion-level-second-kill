package com.seckill.usersystem.service.impl;

import com.seckill.usersystem.entity.Permission;
import com.seckill.usersystem.mapper.PermissionMapper;
import com.seckill.usersystem.service.IPermissionService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class PermissionServiceImpl implements IPermissionService {

    private final PermissionMapper permissionMapper;

    @Override
    public List<Permission> listPermissions() {
        return permissionMapper.selectAll();
    }

    @Override
    public List<Permission> getPermissionsByUserId(Long userId) {
        return List.of();
    }

    @Override
    public List<Permission> getPermissionsByRoleId(Long roleId) {
        return List.of();
    }

    @Override
    public void savePermission(Permission permission) {
        permissionMapper.insert(permission);
    }

    @Override
    public void updatePermission(Permission permission) {
        permissionMapper.update(permission);
    }

    @Override
    public void deletePermission(Long permissionId) {
        Permission permission = new Permission();
        permission.setId(permissionId);
        permission.setIsDeleted(1);
        permissionMapper.update(permission);
    }

    @Override
    public List<String> getPermissionCodesByUserId(Long userId) {
        return List.of();
    }
}
