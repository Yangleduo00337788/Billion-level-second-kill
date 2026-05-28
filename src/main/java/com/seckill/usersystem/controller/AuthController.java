package com.seckill.usersystem.controller;

import com.seckill.seckill.limiter.TokenBucketLimiter;
import com.seckill.usersystem.dto.LoginRequest;
import com.seckill.usersystem.dto.RefreshTokenRequest;
import com.seckill.usersystem.dto.RegisterRequest;
import com.seckill.usersystem.service.IRiskControlService;
import com.seckill.usersystem.service.IUserService;
import com.seckill.usersystem.vo.LoginResultVO;
import com.seckill.usersystem.vo.Result;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

@Tag(name = "认证管理")
@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final IUserService userService;
    private final IRiskControlService riskControlService;
    private final TokenBucketLimiter tokenBucketLimiter;

    @Operation(summary = "用户登录")
    @PostMapping("/login")
    public Result<LoginResultVO> login(@Valid @RequestBody LoginRequest request,
                                        HttpServletRequest httpRequest) {
        String ip = httpRequest.getRemoteAddr();
        if (!tokenBucketLimiter.tryAcquire("login:ip:" + ip, 100, 200)) {
            return Result.fail(429, "请求过于频繁，请稍后重试");
        }
        return Result.success(userService.login(request));
    }

    @Operation(summary = "用户注册")
    @PostMapping("/register")
    public Result<LoginResultVO> register(@Valid @RequestBody RegisterRequest request,
                                           HttpServletRequest httpRequest) {
        String ip = httpRequest.getRemoteAddr();
        if (!tokenBucketLimiter.tryAcquire("register:ip:" + ip, 10, 500)) {
            return Result.fail(429, "注册请求过于频繁，请稍后重试");
        }
        return Result.success(userService.register(request));
    }

    @Operation(summary = "验证码登录")
    @PostMapping("/login-with-captcha")
    public Result<LoginResultVO> loginWithCaptcha(@Valid @RequestBody LoginRequest request,
                                                   @RequestParam String captcha,
                                                   HttpServletRequest httpRequest) {
        String ip = httpRequest.getRemoteAddr();
        String deviceId = request.getDeviceId();
        if (!riskControlService.isCaptchaRequired(ip, deviceId)) {
            return Result.success(userService.login(request));
        }
        riskControlService.clearCaptchaRequired(ip, deviceId);
        return Result.success(userService.login(request));
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
