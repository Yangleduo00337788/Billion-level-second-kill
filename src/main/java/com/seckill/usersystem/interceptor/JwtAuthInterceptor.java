package com.seckill.usersystem.interceptor;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.seckill.usersystem.util.JwtUtil;
import com.seckill.usersystem.util.RedisUtil;
import com.seckill.usersystem.vo.Result;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

/**
 * JWT 认证拦截器
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class JwtAuthInterceptor implements HandlerInterceptor {

    private final JwtUtil jwtUtil;
    private final RedisUtil redisUtil;
    private final ObjectMapper objectMapper;

    private static final String AUTH_HEADER = "Authorization";
    private static final String BEARER_PREFIX = "Bearer ";

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
        if ("OPTIONS".equalsIgnoreCase(request.getMethod())) {
            return true;
        }

        String authHeader = request.getHeader(AUTH_HEADER);
        if (authHeader == null || !authHeader.startsWith(BEARER_PREFIX)) {
            sendUnauthorized(response, "未提供认证令牌");
            return false;
        }

        String token = authHeader.substring(BEARER_PREFIX.length());

        if (!jwtUtil.isTokenValid(token)) {
            sendUnauthorized(response, "认证令牌已过期");
            return false;
        }

        String tokenType = jwtUtil.getTokenType(token);
        if (!"access".equals(tokenType)) {
            sendUnauthorized(response, "无效的令牌类型");
            return false;
        }

        Long userId = jwtUtil.getUserId(token);
        String deviceId = jwtUtil.getDeviceId(token);

        String storedToken = redisUtil.get("user:token:" + userId + ":" + deviceId);
        if (storedToken == null || !storedToken.equals(token)) {
            sendUnauthorized(response, "认证令牌已失效，请重新登录");
            return false;
        }

        String blacklistKey = "token:blacklist:" + token;
        if (redisUtil.hasKey(blacklistKey)) {
            sendUnauthorized(response, "认证令牌已被撤销");
            return false;
        }

        request.setAttribute("userId", userId);
        request.setAttribute("username", jwtUtil.getUsername(token));
        request.setAttribute("deviceId", deviceId);

        return true;
    }

    private void sendUnauthorized(HttpServletResponse response, String message) throws Exception {
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        response.setCharacterEncoding("UTF-8");
        response.getWriter().write(
                objectMapper.writeValueAsString(Result.unauthorized(message))
        );
    }
}
