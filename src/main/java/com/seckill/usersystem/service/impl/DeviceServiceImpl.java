package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.entity.Device;
import com.seckill.usersystem.enums.DeviceTypeEnum;
import com.seckill.usersystem.mapper.DeviceMapper;
import com.seckill.usersystem.service.IDeviceService;
import com.seckill.usersystem.util.RedisUtil;
import com.seckill.usersystem.vo.DeviceVO;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class DeviceServiceImpl implements IDeviceService {

    private final DeviceMapper deviceMapper;
    private final RedisUtil redisUtil;

    private static final String DEVICE_TRUSTED_KEY = "device:trusted:";

    @Override
    public Device getDevice(Long userId, String deviceId) {
        return QueryChain.of(deviceMapper)
                .where(Device::getUserId).eq(userId)
                .where(Device::getDeviceId).eq(deviceId)
                .one();
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public Device saveDevice(Device device) {
        deviceMapper.insert(device);
        return device;
    }

    @Override
    public void updateDeviceActivity(Long userId, String deviceId) {
        Device device = getDevice(userId, deviceId);
        if (device != null) {
            device.setLastActiveTime(LocalDateTime.now());
            deviceMapper.update(device);
        }
    }

    @Override
    public List<DeviceVO> getUserDevices(Long userId, String currentDeviceId) {
        return getDevicesByUserId(userId).stream()
                .map(d -> convertToDeviceVO(d, currentDeviceId))
                .sorted(Comparator.comparing(DeviceVO::getLastActiveTime).reversed())
                .collect(Collectors.toList());
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void removeDevice(Long userId, String deviceId) {
        Device device = getDevice(userId, deviceId);
        if (device != null) {
            deviceMapper.deleteById(device.getId());
            redisUtil.delete("user:token:" + userId + ":" + deviceId);
            redisUtil.delete("user:refresh:" + userId + ":" + deviceId);
            log.info("移除设备: userId={}, deviceId={}", userId, deviceId);
        }
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void removeAllOtherDevices(Long userId, String currentDeviceId) {
        List<Device> devices = getDevicesByUserId(userId);
        for (Device device : devices) {
            if (!device.getDeviceId().equals(currentDeviceId)) {
                deviceMapper.deleteById(device.getId());
                redisUtil.delete("user:token:" + userId + ":" + device.getDeviceId());
                redisUtil.delete("user:refresh:" + userId + ":" + device.getDeviceId());
            }
        }
    }

    @Override
    public void setDeviceTrusted(Long userId, String deviceId, boolean trusted) {
        String key = DEVICE_TRUSTED_KEY + userId + ":" + deviceId;
        if (trusted) {
            redisUtil.set(key, "1");
        } else {
            redisUtil.delete(key);
        }

        Device device = getDevice(userId, deviceId);
        if (device != null) {
            device.setIsTrusted(trusted ? 1 : 0);
            deviceMapper.update(device);
        }
    }

    @Override
    public boolean isDeviceTrusted(Long userId, String deviceId) {
        Device device = getDevice(userId, deviceId);
        return device != null && device.getIsTrusted() == 1;
    }

    @Override
    public List<Device> getDevicesByUserId(Long userId) {
        return QueryChain.of(deviceMapper)
                .where(Device::getUserId).eq(userId)
                .where(Device::getStatus).eq(1)
                .list();
    }

    @Override
    public int getDeviceCount(Long userId) {
        return getDevicesByUserId(userId).size();
    }

    private DeviceVO convertToDeviceVO(Device device, String currentDeviceId) {
        DeviceVO vo = new DeviceVO();
        vo.setId(device.getId());
        vo.setDeviceId(device.getDeviceId());
        vo.setDeviceType(device.getDeviceType());
        vo.setDeviceName(device.getDeviceName());
        vo.setOsName(device.getOsName());
        vo.setBrowser(device.getBrowser());
        vo.setIpAddress(device.getIpAddress());
        vo.setLocation(device.getLocation());
        vo.setLastLoginTime(device.getLastLoginTime());
        vo.setLastActiveTime(device.getLastActiveTime());
        vo.setStatus(device.getStatus());
        vo.setIsTrusted(device.getIsTrusted() == 1);
        vo.setIsCurrentDevice(device.getDeviceId().equals(currentDeviceId));
        return vo;
    }
}
