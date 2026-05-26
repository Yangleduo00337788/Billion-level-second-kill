package com.seckill.usersystem.service;

import com.seckill.usersystem.entity.Role;

import java.util.List;
import java.util.Set;

public interface IRoleService {

    Role getRoleByCode(String roleCode);

    List<Role> listRoles();

    void saveRole(Role role);

    void updateRole(Role role);

    void deleteRole(Long roleId);

    void assignPermissions(Long roleId, List<Long> permissionIds);

    Set<Long> getRolePermissionIds(Long roleId);

    List<Role> getRolesByUserId(Long userId);
}
