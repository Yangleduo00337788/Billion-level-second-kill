package com.seckill.im.handler;

import com.seckill.im.proto.ImMessageProto;
import com.seckill.im.server.ChannelAttr;
import com.seckill.im.server.ChannelManager;
import com.seckill.im.service.ImMessageService;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
@ChannelHandler.Sharable
public class ImMessageHandler extends SimpleChannelInboundHandler<ImMessageProto.WrapperMessage> {

    private final ImMessageService imMessageService;

    @Override
    public void channelActive(ChannelHandlerContext ctx) {
        log.info("[ImMessageHandler] channel active: {}", ctx.channel().id().asShortText());
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) {
        Long userId = ChannelAttr.getUserId(ctx.channel());
        ChannelManager.removeChannel(ctx.channel());

        if (userId != null && !ChannelManager.isUserOnline(userId)) {
            imMessageService.handleUserOffline(userId);
        }

        log.info("[ImMessageHandler] channel inactive: {}, userId:{}",
                ctx.channel().id().asShortText(), userId);
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, ImMessageProto.WrapperMessage wrapper) {
        Channel channel = ctx.channel();

        if (!ChannelAttr.isAuthenticated(channel)) {
            log.warn("[ImMessageHandler] unauthenticated message from channel:{}",
                    channel.id().asShortText());
            return;
        }

        Long userId = ChannelAttr.getUserId(channel);

        switch (wrapper.getType()) {
            case UPSTREAM_MESSAGE -> handleUpstreamMessage(channel, userId, wrapper.getUpstreamMessage());
            case HEARTBEAT -> handleHeartbeat(channel, userId, wrapper.getHeartbeat());
            case READ_RECEIPT -> handleReadReceipt(channel, userId, wrapper.getReadReceipt());
            case RECALL_MESSAGE -> handleRecallMessage(channel, userId, wrapper.getRecallMessage());
            case SYNC_OFFLINE_REQUEST -> handleSyncOffline(channel, userId, wrapper.getSyncOfflineRequest());
            default -> log.warn("[ImMessageHandler] unknown message type: {} from user:{}",
                    wrapper.getType(), userId);
        }
    }

    private void handleUpstreamMessage(Channel channel, Long userId, ImMessageProto.UpstreamMessage msg) {
        imMessageService.processUpstreamMessage(channel, userId, msg);
    }

    private void handleHeartbeat(Channel channel, Long userId, ImMessageProto.Heartbeat heartbeat) {
        ImMessageProto.HeartbeatAck ack = ImMessageProto.HeartbeatAck.newBuilder()
                .setPongTime(System.currentTimeMillis())
                .setUnreadCount(imMessageService.getUnreadCount(userId))
                .build();

        ImMessageProto.WrapperMessage response = ImMessageProto.WrapperMessage.newBuilder()
                .setType(ImMessageProto.WrapperMessage.MessageType.HEARTBEAT_ACK)
                .setHeartbeatAck(ack)
                .build();

        channel.writeAndFlush(response);
    }

    private void handleReadReceipt(Channel channel, Long userId, ImMessageProto.ReadReceipt receipt) {
        imMessageService.processReadReceipt(userId, receipt);
    }

    private void handleRecallMessage(Channel channel, Long userId, ImMessageProto.RecallMessage recall) {
        imMessageService.processRecall(userId, recall);
    }

    private void handleSyncOffline(Channel channel, Long userId, ImMessageProto.SyncOfflineRequest request) {
        imMessageService.processSyncOffline(channel, userId, request);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
        log.error("[ImMessageHandler] exception, channel:{}",
                ctx.channel().id().asShortText(), cause);
        ctx.close();
    }
}