package com.seckill.usersystem.service;

import com.seckill.usersystem.entity.Device;
import com.seckill.usersystem.vo.DeviceVO;

import java.util.List;

public interface IDeviceService {

    Device getDevice(Long userId, String deviceId);

    Device saveDevice(Device device);

    void updateDeviceActivity(Long userId, String deviceId);

    List<DeviceVO> getUserDevices(Long userId, String currentDeviceId);

    void removeDevice(Long userId, String deviceId);

    void removeAllOtherDevices(Long userId, String currentDeviceId);

    void removeSameTypeDevices(Long userId, String deviceType, String currentDeviceId);

    void setDeviceTrusted(Long userId, String deviceId, boolean trusted);

    boolean isDeviceTrusted(Long userId, String deviceId);

    List<Device> getDevicesByUserId(Long userId);

    int getDeviceCount(Long userId);

    boolean hasSameTypeDevice(Long userId, String deviceType);
}
