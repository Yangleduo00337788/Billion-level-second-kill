package com.seckill.usersystem.controller;

import com.seckill.usersystem.dto.LoginRequest;
import com.seckill.usersystem.dto.RefreshTokenRequest;
import com.seckill.usersystem.dto.RegisterRequest;
import com.seckill.usersystem.service.IUserService;
import com.seckill.usersystem.vo.LoginResultVO;
import com.seckill.usersystem.vo.Result;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

@Tag(name = "认证管理")
@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final IUserService userService;

    @Operation(summary = "用户登录")
    @PostMapping("/login")
    public Result<LoginResultVO> login(@Valid @RequestBody LoginRequest request) {
        return Result.success(userService.login(request));
    }

    @Operation(summary = "用户注册")
    @PostMapping("/register")
    public Result<LoginResultVO> register(@Valid @RequestBody RegisterRequest request) {
        return Result.success(userService.register(request));
    }

    @Operation(summary = "退出登录")
    @PostMapping("/logout")
    public Result<Void> logout(@RequestAttribute Long userId, @RequestAttribute String deviceId) {
        userService.logout(userId, deviceId);
        return Result.success();
    }

    @Operation(summary = "刷新令牌")
    @PostMapping("/refresh-token")
    public Result<?> refreshToken(@Valid @RequestBody RefreshTokenRequest request) {
        return Result.success(userService.refreshToken(request.getRefreshToken(), request.getDeviceId()));
    }
}
