package com.seckill.im.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Data
@Component
@ConfigurationProperties(prefix = "im.netty")
public class ImServerConfig {

    private int port = 9999;

    private int bossThreads = 1;

    private int workerThreads = Runtime.getRuntime().availableProcessors() * 2;

    private int soBacklog = 1024;

    private boolean soKeepalive = true;

    private boolean tcpNodelay = true;

    private int soSndbuf = 65536;

    private int soRcvbuf = 65536;

    private int writeBufferLowWaterMark = 32768;

    private int writeBufferHighWaterMark = 65536;

    private int maxContentLength = 65536;

    private int maxFrameSize = 65536;

    private int readerIdleTimeSeconds = 120;

    private int writerIdleTimeSeconds = 60;

    private int allIdleTimeSeconds = 0;
}