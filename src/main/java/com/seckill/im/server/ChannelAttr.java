package com.seckill.im.server;

import io.netty.channel.Channel;
import io.netty.util.AttributeKey;

public final class ChannelAttr {

    public static final AttributeKey<Long> USER_ID = AttributeKey.valueOf("userId");
    public static final AttributeKey<String> DEVICE_ID = AttributeKey.valueOf("deviceId");
    public static final AttributeKey<String> DEVICE_TYPE = AttributeKey.valueOf("deviceType");
    public static final AttributeKey<Long> LOGIN_TIME = AttributeKey.valueOf("loginTime");

    private ChannelAttr() {
    }

    public static Long getUserId(Channel channel) {
        return channel.attr(USER_ID).get();
    }

    public static void setUserId(Channel channel, Long userId) {
        channel.attr(USER_ID).set(userId);
    }

    public static String getDeviceId(Channel channel) {
        return channel.attr(DEVICE_ID).get();
    }

    public static void setDeviceId(Channel channel, String deviceId) {
        channel.attr(DEVICE_ID).set(deviceId);
    }

    public static String getDeviceType(Channel channel) {
        return channel.attr(DEVICE_TYPE).get();
    }

    public static void setDeviceType(Channel channel, String deviceType) {
        channel.attr(DEVICE_TYPE).set(deviceType);
    }

    public static boolean isAuthenticated(Channel channel) {
        return channel.hasAttr(USER_ID) && channel.attr(USER_ID).get() != null;
    }
}