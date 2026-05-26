package com.seckill.im.server;

import com.seckill.im.config.ImServerConfig;
import com.seckill.im.handler.AuthHandshakeHandler;
import com.seckill.im.handler.HeartbeatHandler;
import com.seckill.im.handler.ImMessageHandler;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpServerCodec;
import io.netty.handler.codec.http.websocketx.WebSocketServerProtocolHandler;
import io.netty.handler.codec.http.websocketx.extensions.compression.WebSocketServerCompressionHandler;
import io.netty.handler.codec.protobuf.ProtobufDecoder;
import io.netty.handler.codec.protobuf.ProtobufEncoder;
import io.netty.handler.codec.protobuf.ProtobufVarint32FrameDecoder;
import io.netty.handler.codec.protobuf.ProtobufVarint32LengthFieldPrepender;
import io.netty.handler.stream.ChunkedWriteHandler;
import io.netty.handler.timeout.IdleStateHandler;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Scope;
import org.springframework.stereotype.Component;

import java.util.concurrent.TimeUnit;

@Component
@Scope("prototype")
@RequiredArgsConstructor
public class ImWebSocketServerInitializer extends ChannelInitializer<SocketChannel> {

    private final ImServerConfig config;
    private final AuthHandshakeHandler authHandshakeHandler;
    private final ImMessageHandler imMessageHandler;
    private final com.seckill.im.proto.ImMessageProto.WrapperMessage defaultInstance =
            com.seckill.im.proto.ImMessageProto.WrapperMessage.getDefaultInstance();

    @Override
    protected void initChannel(SocketChannel ch) {
        ChannelPipeline pipeline = ch.pipeline();

        pipeline.addLast(new HttpServerCodec());
        pipeline.addLast(new ChunkedWriteHandler());
        pipeline.addLast(new HttpObjectAggregator(config.getMaxContentLength()));

        pipeline.addLast(new WebSocketServerCompressionHandler());

        pipeline.addLast(authHandshakeHandler);

        pipeline.addLast(new WebSocketServerProtocolHandler("/ws", null, true,
                config.getMaxFrameSize()));

        pipeline.addLast(new IdleStateHandler(
                config.getReaderIdleTimeSeconds(),
                config.getWriterIdleTimeSeconds(),
                config.getAllIdleTimeSeconds(),
                TimeUnit.SECONDS));
        pipeline.addLast(new HeartbeatHandler());

        pipeline.addLast(new ProtobufVarint32FrameDecoder());
        pipeline.addLast(new ProtobufDecoder(defaultInstance));
        pipeline.addLast(new ProtobufVarint32LengthFieldPrepender());
        pipeline.addLast(new ProtobufEncoder());

        pipeline.addLast(imMessageHandler);
    }
}