package com.seckill.usersystem.service.impl;

import com.seckill.usersystem.entity.Permission;
import com.seckill.usersystem.mapper.PermissionMapper;
import com.seckill.usersystem.service.IPermissionService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class PermissionServiceImpl implements IPermissionService {

    private final PermissionMapper permissionMapper;
    private final RedisUtil redisUtil;

    private static final String USER_PERMISSIONS_KEY = "user:permissions:";
    private static final String USER_PERM_CODES_KEY = "user:perm:codes:";

    @Override
    public List<Permission> listPermissions() {
        return permissionMapper.selectAll();
    }

    @Override
    public List<Permission> getPermissionsByUserId(Long userId) {
        String cacheKey = USER_PERMISSIONS_KEY + userId;
        String cachedIds = redisUtil.get(cacheKey);
        if (cachedIds != null && !cachedIds.isEmpty()) {
            try {
                String[] ids = cachedIds.split(",");
                List<Permission> cachedPerms = new java.util.ArrayList<>();
                for (String idStr : ids) {
                    if (!idStr.isEmpty()) {
                        Permission p = permissionMapper.selectOneById(Long.parseLong(idStr));
                        if (p != null) cachedPerms.add(p);
                    }
                }
                if (!cachedPerms.isEmpty()) return cachedPerms;
            } catch (Exception e) {
                log.warn("解析权限缓存失败: userId={}", userId);
            }
        }

        List<Permission> permissions = permissionMapper.selectPermissionsByUserId(userId);
        if (permissions != null && !permissions.isEmpty()) {
            String idsStr = permissions.stream().map(p -> String.valueOf(p.getId())).collect(Collectors.joining(","));
            redisUtil.set(cacheKey, idsStr, 1800, TimeUnit.SECONDS);
        }
        return permissions != null ? permissions : Collections.emptyList();
    }

    @Override
    public List<Permission> getPermissionsByRoleId(Long roleId) {
        return permissionMapper.selectPermissionsByRoleId(roleId);
    }

    @Override
    public void savePermission(Permission permission) {
        permissionMapper.insert(permission);
        invalidateAllPermissionCache();
    }

    @Override
    public void updatePermission(Permission permission) {
        permissionMapper.update(permission);
        invalidateAllPermissionCache();
    }

    @Override
    public void deletePermission(Long permissionId) {
        Permission permission = new Permission();
        permission.setId(permissionId);
        permission.setIsDeleted(1);
        permissionMapper.update(permission);
        invalidateAllPermissionCache();
    }

    @Override
    public List<String> getPermissionCodesByUserId(Long userId) {
        String cacheKey = USER_PERM_CODES_KEY + userId;
        String cached = redisUtil.get(cacheKey);
        if (cached != null && !cached.isEmpty()) {
            try {
                return List.of(cached.split(","));
            } catch (Exception e) {
                log.warn("解析权限码缓存失败: userId={}", userId);
            }
        }

        List<String> codes = permissionMapper.selectPermissionCodesByUserId(userId);
        if (codes != null && !codes.isEmpty()) {
            redisUtil.set(cacheKey, String.join(",", codes), 1800, TimeUnit.SECONDS);
        }
        return codes != null ? codes : Collections.emptyList();
    }

    private void invalidateAllPermissionCache() {
        redisUtil.deleteByPattern(USER_PERMISSIONS_KEY + "*");
        redisUtil.deleteByPattern(USER_PERM_CODES_KEY + "*");
    }
}
