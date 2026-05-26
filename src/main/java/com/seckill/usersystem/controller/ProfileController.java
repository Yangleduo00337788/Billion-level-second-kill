package com.seckill.usersystem.controller;

import com.seckill.usersystem.dto.UserProfileDTO;
import com.seckill.usersystem.entity.UserProfile;
import com.seckill.usersystem.service.IUserProfileService;
import com.seckill.usersystem.vo.Result;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

@Tag(name = "用户资料管理")
@RestController
@RequestMapping("/api/profile")
@RequiredArgsConstructor
public class ProfileController {

    private final IUserProfileService userProfileService;

    @Operation(summary = "获取用户资料")
    @GetMapping("/me")
    public Result<UserProfile> getMyProfile(@RequestAttribute Long userId) {
        return Result.success(userProfileService.getOrCreateProfile(userId));
    }

    @Operation(summary = "更新用户资料")
    @PutMapping("/me")
    public Result<Void> updateMyProfile(@Valid @RequestBody UserProfileDTO profileDTO, @RequestAttribute Long userId) {
        profileDTO.setUserId(userId);
        userProfileService.updateProfile(profileDTO);
        return Result.success();
    }
}
