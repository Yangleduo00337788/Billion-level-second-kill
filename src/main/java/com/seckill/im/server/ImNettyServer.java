package com.seckill.im.server;

import com.seckill.im.config.ImServerConfig;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.util.concurrent.DefaultThreadFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.stereotype.Component;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Slf4j
@Component
@RequiredArgsConstructor
public class ImNettyServer {

    private final ImServerConfig config;
    private final ObjectProvider<ImWebSocketServerInitializer> initializerProvider;

    private EventLoopGroup bossGroup;
    private EventLoopGroup workerGroup;
    private Channel serverChannel;

    @PostConstruct
    public void start() throws InterruptedException {
        bossGroup = new NioEventLoopGroup(config.getBossThreads(),
                new DefaultThreadFactory("im-boss"));
        workerGroup = new NioEventLoopGroup(config.getWorkerThreads(),
                new DefaultThreadFactory("im-worker"));

        ServerBootstrap bootstrap = new ServerBootstrap();
        bootstrap.group(bossGroup, workerGroup)
                .channel(NioServerSocketChannel.class)
                .handler(new LoggingHandler(LogLevel.INFO))
                .childHandler(initializerProvider.getObject())
                .option(ChannelOption.SO_BACKLOG, config.getSoBacklog())
                .childOption(ChannelOption.SO_KEEPALIVE, config.isSoKeepalive())
                .childOption(ChannelOption.TCP_NODELAY, config.isTcpNodelay())
                .childOption(ChannelOption.SO_SNDBUF, config.getSoSndbuf())
                .childOption(ChannelOption.SO_RCVBUF, config.getSoRcvbuf())
                .childOption(ChannelOption.WRITE_BUFFER_WATER_MARK,
                        new WriteBufferWaterMark(config.getWriteBufferLowWaterMark(),
                                config.getWriteBufferHighWaterMark()));

        ChannelFuture future = bootstrap.bind(config.getPort()).sync();
        serverChannel = future.channel();

        log.info("[ImNettyServer] started on port:{}, boss:{}, worker:{}",
                config.getPort(), config.getBossThreads(), config.getWorkerThreads());
    }

    @PreDestroy
    public void stop() {
        ChannelManager.closeAll();
        if (serverChannel != null) {
            serverChannel.close();
        }
        if (bossGroup != null) {
            bossGroup.shutdownGracefully();
        }
        if (workerGroup != null) {
            workerGroup.shutdownGracefully();
        }
        log.info("[ImNettyServer] stopped");
    }
}