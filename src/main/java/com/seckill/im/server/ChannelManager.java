package com.seckill.im.server;

import io.netty.channel.Channel;
import io.netty.channel.group.ChannelGroup;
import io.netty.channel.group.DefaultChannelGroup;
import io.netty.util.concurrent.GlobalEventExecutor;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArraySet;

@Slf4j
public class ChannelManager {

    private ChannelManager() {
    }

    private static final ChannelGroup ALL_CHANNELS = new DefaultChannelGroup(GlobalEventExecutor.INSTANCE);

    private static final Map<Long, Set<Channel>> USER_CHANNEL_MAP = new ConcurrentHashMap<>();

    private static final Map<String, Long> CHANNEL_USER_MAP = new ConcurrentHashMap<>();

    public static void addChannel(Channel channel) {
        ALL_CHANNELS.add(channel);
        channel.closeFuture().addListener(future -> removeChannel(channel));
    }

    public static void bindUser(Channel channel) {
        Long userId = ChannelAttr.getUserId(channel);
        String channelId = channel.id().asLongText();

        USER_CHANNEL_MAP.computeIfAbsent(userId, k -> new CopyOnWriteArraySet<>()).add(channel);
        CHANNEL_USER_MAP.put(channelId, userId);

        log.info("[ChannelManager] bind user:{} channel:{}, total users:{}, total channels:{}",
                userId, channel.id().asShortText(),
                USER_CHANNEL_MAP.size(), ALL_CHANNELS.size());
    }

    public static void removeChannel(Channel channel) {
        ALL_CHANNELS.remove(channel);

        String channelId = channel.id().asLongText();
        Long userId = CHANNEL_USER_MAP.remove(channelId);

        if (userId != null) {
            Set<Channel> channels = USER_CHANNEL_MAP.get(userId);
            if (channels != null) {
                channels.remove(channel);
                if (channels.isEmpty()) {
                    USER_CHANNEL_MAP.remove(userId);
                }
            }
        }

        log.info("[ChannelManager] remove channel:{}, userId:{}, remaining channels:{}",
                channel.id().asShortText(), userId, ALL_CHANNELS.size());
    }

    public static Set<Channel> getUserChannels(Long userId) {
        return USER_CHANNEL_MAP.getOrDefault(userId, Set.of());
    }

    public static boolean isUserOnline(Long userId) {
        Set<Channel> channels = USER_CHANNEL_MAP.get(userId);
        return channels != null && !channels.isEmpty();
    }

    public static int getOnlineUserCount() {
        return USER_CHANNEL_MAP.size();
    }

    public static int getTotalChannelCount() {
        return ALL_CHANNELS.size();
    }

    public static Long getUserIdByChannel(Channel channel) {
        return CHANNEL_USER_MAP.get(channel.id().asLongText());
    }

    public static void closeAll() {
        ALL_CHANNELS.close();
        USER_CHANNEL_MAP.clear();
        CHANNEL_USER_MAP.clear();
    }
}