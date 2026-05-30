package com.seckill.im.handler;

import com.seckill.im.server.ChannelAttr;
import com.seckill.im.server.ChannelManager;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.handler.timeout.IdleState;
import io.netty.handler.timeout.IdleStateEvent;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@ChannelHandler.Sharable
public class HeartbeatHandler extends ChannelInboundHandlerAdapter {

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        if (evt instanceof IdleStateEvent event) {
            if (event.state() == IdleState.READER_IDLE) {
                handleReaderIdle(ctx);
            } else if (event.state() == IdleState.WRITER_IDLE) {
                handleWriterIdle(ctx);
            }
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    private void handleReaderIdle(ChannelHandlerContext ctx) {
        Long userId = ChannelAttr.getUserId(ctx.channel());
        log.info("[Heartbeat] reader idle, closing channel:{} userId:{}",
                ctx.channel().id().asShortText(), userId);
        ctx.close();
    }

    private void handleWriterIdle(ChannelHandlerContext ctx) {
        ctx.channel().writeAndFlush(
                com.seckill.im.proto.ImMessageProto.WrapperMessage.newBuilder()
                        .setType(com.seckill.im.proto.ImMessageProto.WrapperMessage.MessageType.HEARTBEAT_ACK)
                        .setHeartbeatAck(com.seckill.im.proto.ImMessageProto.HeartbeatAck.newBuilder()
                                .setPongTime(System.currentTimeMillis())
                                .build())
                        .build()
        );
    }
}