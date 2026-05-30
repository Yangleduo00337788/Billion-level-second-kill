package com.seckill.admin.controller;

import com.seckill.usersystem.entity.Blacklist;
import com.seckill.usersystem.entity.Device;
import com.seckill.usersystem.entity.Permission;
import com.seckill.usersystem.entity.Role;
import com.seckill.usersystem.mapper.*;
import com.seckill.usersystem.service.IBlacklistService;
import com.seckill.usersystem.service.IPermissionService;
import com.seckill.usersystem.vo.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin")
@RequiredArgsConstructor
public class AdminUserController {

    private final BlacklistMapper blacklistMapper;
    private final IBlacklistService blacklistService;
    private final DeviceMapper deviceMapper;
    private final PermissionMapper permissionMapper;
    private final IPermissionService permissionService;
    private final RoleMapper roleMapper;

    @GetMapping("/blacklist")
    public Result<List<Blacklist>> listBlacklist() {
        return Result.success(blacklistMapper.selectAll());
    }

    @PostMapping("/blacklist")
    public Result<Void> addBlacklist(@RequestBody Blacklist blacklist) {
        blacklistService.addToBlacklist(blacklist);
        return Result.success();
    }

    @DeleteMapping("/blacklist/{id}")
    public Result<Void> removeBlacklist(@PathVariable Long id) {
        blacklistService.removeUserFromBlacklist(id);
        return Result.success();
    }

    @GetMapping("/devices")
    public Result<List<Device>> listAllDevices() {
        return Result.success(deviceMapper.selectAll());
    }

    @GetMapping("/permissions")
    public Result<List<Permission>> listPermissions() {
        return Result.success(permissionMapper.selectAll());
    }

    @PostMapping("/permissions")
    public Result<Void> savePermission(@RequestBody Permission permission) {
        permissionService.savePermission(permission);
        return Result.success();
    }

    @PutMapping("/permissions/{id}")
    public Result<Void> updatePermission(@PathVariable Long id, @RequestBody Permission permission) {
        permission.setId(id);
        permissionService.updatePermission(permission);
        return Result.success();
    }

    @DeleteMapping("/permissions/{id}")
    public Result<Void> deletePermission(@PathVariable Long id) {
        permissionService.deletePermission(id);
        return Result.success();
    }

    @GetMapping("/roles/permissions/{roleId}")
    public Result<List<Permission>> getRolePermissions(@PathVariable Long roleId) {
        return Result.success(permissionMapper.selectPermissionsByRoleId(roleId));
    }

    @PutMapping("/roles/{roleId}/permissions")
    public Result<Void> assignPermissions(@PathVariable Long roleId, @RequestBody List<Long> permissionIds) {
        roleMapper.deleteRolePermissions(roleId);
        if (permissionIds != null && !permissionIds.isEmpty()) {
            roleMapper.insertRolePermissions(roleId, permissionIds);
        }
        return Result.success();
    }
}