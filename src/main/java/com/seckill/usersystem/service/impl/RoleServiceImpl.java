package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.entity.Role;
import com.seckill.usersystem.mapper.RoleMapper;
import com.seckill.usersystem.service.IRoleService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class RoleServiceImpl implements IRoleService {

    private final RoleMapper roleMapper;
    private final RedisUtil redisUtil;

    private static final String USER_ROLES_KEY = "user:roles:";
    private static final String USER_ROLES_PERM_KEY = "user:roles:perm:";

    @Override
    public Role getRoleByCode(String roleCode) {
        return QueryChain.of(roleMapper)
                .where(Role::getRoleCode).eq(roleCode)
                .one();
    }

    @Override
    public List<Role> listRoles() {
        return roleMapper.selectAll();
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void saveRole(Role role) {
        roleMapper.insert(role);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void updateRole(Role role) {
        roleMapper.update(role);
        redisUtil.deleteByPattern(USER_ROLES_KEY + "*");
        redisUtil.deleteByPattern(USER_ROLES_PERM_KEY + "*");
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void deleteRole(Long roleId) {
        Role role = new Role();
        role.setId(roleId);
        role.setIsDeleted(1);
        roleMapper.update(role);
        redisUtil.deleteByPattern(USER_ROLES_KEY + "*");
        redisUtil.deleteByPattern(USER_ROLES_PERM_KEY + "*");
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void assignPermissions(Long roleId, List<Long> permissionIds) {
        roleMapper.deleteRolePermissions(roleId);
        if (permissionIds != null && !permissionIds.isEmpty()) {
            roleMapper.insertRolePermissions(roleId, permissionIds);
        }
        redisUtil.deleteByPattern(USER_ROLES_PERM_KEY + "*");
    }

    @Override
    public Set<Long> getRolePermissionIds(Long roleId) {
        List<Long> ids = roleMapper.selectPermissionIdsByRoleId(roleId);
        return ids != null ? Set.copyOf(ids) : Collections.emptySet();
    }

    @Override
    public List<Role> getRolesByUserId(Long userId) {
        String cacheKey = USER_ROLES_KEY + userId;
        String cached = redisUtil.get(cacheKey);
        if (cached != null) {
            try {
                String[] roleIdsArr = cached.split(",");
                List<Role> cachedRoles = new ArrayList<>();
                for (String idStr : roleIdsArr) {
                    if (!idStr.isEmpty()) {
                        Role r = roleMapper.selectOneById(Long.parseLong(idStr));
                        if (r != null) cachedRoles.add(r);
                    }
                }
                if (!cachedRoles.isEmpty()) return cachedRoles;
            } catch (Exception e) {
                log.warn("解析角色缓存失败: userId={}", userId);
            }
        }

        List<Role> roles = roleMapper.selectRolesByUserId(userId);
        if (roles != null && !roles.isEmpty()) {
            String idsStr = roles.stream().map(r -> String.valueOf(r.getId())).collect(Collectors.joining(","));
            redisUtil.set(cacheKey, idsStr, 3600);
        }
        return roles != null ? roles : Collections.emptyList();
    }
}
