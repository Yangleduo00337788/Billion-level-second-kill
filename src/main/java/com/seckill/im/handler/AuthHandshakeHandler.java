package com.seckill.im.handler;

import com.seckill.im.server.ChannelAttr;
import com.seckill.im.server.ChannelManager;
import com.seckill.usersystem.util.JwtUtil;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.handler.codec.http.FullHttpRequest;
import io.netty.handler.codec.http.websocketx.CloseWebSocketFrame;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.net.URI;
import java.util.List;
import java.util.Map;

@Slf4j
@Component
@RequiredArgsConstructor
@ChannelHandler.Sharable
public class AuthHandshakeHandler extends ChannelInboundHandlerAdapter {

    private static final String TOKEN_PARAM = "token";
    private static final String DEVICE_ID_PARAM = "deviceId";
    private static final String DEVICE_TYPE_PARAM = "deviceType";

    private final JwtUtil jwtUtil;

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (msg instanceof FullHttpRequest httpRequest) {
            handleHandshake(ctx, httpRequest);
        } else {
            super.channelRead(ctx, msg);
        }
    }

    private void handleHandshake(ChannelHandlerContext ctx, FullHttpRequest request) {
        String uri = request.uri();
        Map<String, List<String>> params = parseQueryParams(uri);

        String token = getParam(params, TOKEN_PARAM);
        String deviceId = getParam(params, DEVICE_ID_PARAM);
        String deviceType = getParam(params, DEVICE_TYPE_PARAM);

        if (!StringUtils.hasText(token)) {
            log.warn("[AuthHandshakeHandler] missing token, remote:{}", ctx.channel().remoteAddress());
            ctx.writeAndFlush(new CloseWebSocketFrame(4001, "Missing token"));
            ctx.close();
            return;
        }

        Long userId;
        try {
            userId = jwtUtil.getUserId(token);
        } catch (Exception e) {
            log.warn("[AuthHandshakeHandler] invalid token, remote:{}", ctx.channel().remoteAddress());
            ctx.writeAndFlush(new CloseWebSocketFrame(4002, "Invalid token"));
            ctx.close();
            return;
        }

        ChannelAttr.setUserId(ctx.channel(), userId);
        ChannelAttr.setDeviceId(ctx.channel(), deviceId);
        ChannelAttr.setDeviceType(ctx.channel(), deviceType);

        ChannelManager.addChannel(ctx.channel());
        ChannelManager.bindUser(ctx.channel());

        request.setUri("/ws");

        log.info("[AuthHandshakeHandler] user:{} device:{} type:{} authenticated",
                userId, deviceId, deviceType);
    }

    private Map<String, List<String>> parseQueryParams(String uri) {
        int queryIndex = uri.indexOf('?');
        if (queryIndex < 0) {
            return Map.of();
        }
        String query = uri.substring(queryIndex + 1);
        Map<String, List<String>> params = new java.util.HashMap<>();
        for (String pair : query.split("&")) {
            int eqIndex = pair.indexOf('=');
            if (eqIndex > 0) {
                String key = pair.substring(0, eqIndex);
                String value = pair.substring(eqIndex + 1);
                params.computeIfAbsent(key, k -> new java.util.ArrayList<>()).add(value);
            }
        }
        return params;
    }

    private String getParam(Map<String, List<String>> params, String key) {
        List<String> values = params.get(key);
        return (values != null && !values.isEmpty()) ? values.get(0) : null;
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
        log.error("[AuthHandshakeHandler] exception, channel:{}", ctx.channel().id().asShortText(), cause);
        ctx.close();
    }
}