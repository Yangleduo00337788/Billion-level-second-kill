package com.seckill.usersystem.controller;

import com.seckill.usersystem.entity.Role;
import com.seckill.usersystem.service.IRoleService;
import com.seckill.usersystem.vo.Result;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Tag(name = "角色管理")
@RestController
@RequestMapping("/api/roles")
@RequiredArgsConstructor
public class RoleController {

    private final IRoleService roleService;

    @Operation(summary = "获取角色列表")
    @GetMapping("/list")
    public Result<List<Role>> listRoles() {
        return Result.success(roleService.listRoles());
    }

    @Operation(summary = "新增角色")
    @PostMapping
    public Result<Void> saveRole(@RequestBody Role role) {
        roleService.saveRole(role);
        return Result.success();
    }

    @Operation(summary = "更新角色")
    @PutMapping("/{roleId}")
    public Result<Void> updateRole(@PathVariable Long roleId, @RequestBody Role role) {
        role.setId(roleId);
        roleService.updateRole(role);
        return Result.success();
    }

    @Operation(summary = "删除角色")
    @DeleteMapping("/{roleId}")
    public Result<Void> deleteRole(@PathVariable Long roleId) {
        roleService.deleteRole(roleId);
        return Result.success();
    }
}
