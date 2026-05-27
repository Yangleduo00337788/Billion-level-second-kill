package com.seckill.usersystem.controller;

import com.seckill.usersystem.service.IUserService;
import com.seckill.usersystem.vo.Result;
import com.seckill.usersystem.vo.UserVO;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Tag(name = "用户管理")
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final IUserService userService;

    @Operation(summary = "获取用户列表")
    @GetMapping("/list")
    public Result<List<UserVO>> listUsers(
            @RequestParam(defaultValue = "1") int pageNum,
            @RequestParam(defaultValue = "10") int pageSize) {
        return Result.success(userService.listUsers(pageNum, pageSize));
    }

    @Operation(summary = "更新用户状态")
    @PutMapping("/{userId}/status")
    public Result<Void> updateStatus(@PathVariable Long userId, @RequestParam Integer status) {
        userService.updateUserStatus(userId, status);
        return Result.success();
    }

    @Operation(summary = "删除用户")
    @DeleteMapping("/{userId}")
    public Result<Void> deleteUser(@PathVariable Long userId) {
        userService.deleteUser(userId);
        return Result.success();
    }

    @Operation(summary = "锁定用户")
    @PostMapping("/{userId}/lock")
    public Result<Void> lockUser(@PathVariable Long userId, @RequestParam(defaultValue = "30") long lockMinutes) {
        userService.lockUser(userId, lockMinutes);
        return Result.success();
    }

    @Operation(summary = "重置密码")
    @PutMapping("/{userId}/password")
    public Result<Void> resetPassword(@PathVariable Long userId, @RequestParam String newPassword) {
        userService.resetPassword(userId, newPassword);
        return Result.success();
    }
}
