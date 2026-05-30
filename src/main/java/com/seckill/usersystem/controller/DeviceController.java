package com.seckill.usersystem.controller;

import com.seckill.usersystem.service.IBlacklistService;
import com.seckill.usersystem.service.IDeviceService;
import com.seckill.usersystem.vo.DeviceVO;
import com.seckill.usersystem.vo.Result;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Tag(name = "设备管理")
@RestController
@RequestMapping("/api/devices")
@RequiredArgsConstructor
public class DeviceController {

    private final IDeviceService deviceService;
    private final IBlacklistService blacklistService;

    @Operation(summary = "获取我的设备列表")
    @GetMapping("/my")
    public Result<List<DeviceVO>> getMyDevices(@RequestAttribute Long userId, @RequestAttribute String deviceId) {
        return Result.success(deviceService.getUserDevices(userId, deviceId));
    }

    @Operation(summary = "移除设备")
    @DeleteMapping("/{deviceId}")
    public Result<Void> removeDevice(@RequestAttribute Long userId, @PathVariable String deviceId) {
        deviceService.removeDevice(userId, deviceId);
        return Result.success();
    }

    @Operation(summary = "移除所有其他设备")
    @PostMapping("/remove-others")
    public Result<Void> removeAllOthers(@RequestAttribute Long userId, @RequestAttribute String deviceId) {
        deviceService.removeAllOtherDevices(userId, deviceId);
        return Result.success();
    }

    @Operation(summary = "设置可信设备")
    @PutMapping("/{deviceId}/trust")
    public Result<Void> setTrusted(@RequestAttribute Long userId, @PathVariable String deviceId, @RequestParam boolean trusted) {
        deviceService.setDeviceTrusted(userId, deviceId, trusted);
        return Result.success();
    }
}
