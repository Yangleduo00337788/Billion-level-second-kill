package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.entity.Role;
import com.seckill.usersystem.mapper.RoleMapper;
import com.seckill.usersystem.service.IRoleService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class RoleServiceImpl implements IRoleService {

    private final RoleMapper roleMapper;
    private final RedisUtil redisUtil;

    private static final String USER_ROLES_KEY = "user:roles:";

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
        redisUtil.delete(USER_ROLES_KEY + "*");
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void deleteRole(Long roleId) {
        Role role = new Role();
        role.setId(roleId);
        role.setIsDeleted(1);
        roleMapper.update(role);
        redisUtil.delete(USER_ROLES_KEY + "*");
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void assignPermissions(Long roleId, List<Long> permissionIds) {
    }

    @Override
    public Set<Long> getRolePermissionIds(Long roleId) {
        return Set.of();
    }

    @Override
    public List<Role> getRolesByUserId(Long userId) {
        return List.of();
    }
}
