package com.seckill.usersystem.util;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

/**
 * JWT 工具类
 * 基于 jjwt 0.12.x 实现，支持 HS256 签名
 */
@Slf4j
@Component
public class JwtUtil {

    @Value("${jwt.secret:user-system-secret-key-must-be-at-least-256-bits-long-for-hs256-algorithm}")
    private String secret;

    @Value("${jwt.access-token-expiration:7200}")
    private Long accessTokenExpiration;

    @Value("${jwt.refresh-token-expiration:604800}")
    private Long refreshTokenExpiration;

    /**
     * 获取签名密钥
     */
    private SecretKey getSigningKey() {
        byte[] keyBytes = secret.getBytes(StandardCharsets.UTF_8);
        return Keys.hmacShaKeyFor(keyBytes);
    }

    /**
     * 生成访问令牌
     */
    public String generateAccessToken(Long userId, String username, String deviceId) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", userId);
        claims.put("username", username);
        claims.put("deviceId", deviceId);
        claims.put("type", "access");
        return generateToken(claims, accessTokenExpiration);
    }

    /**
     * 生成刷新令牌
     */
    public String generateRefreshToken(Long userId, String deviceId) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", userId);
        claims.put("deviceId", deviceId);
        claims.put("type", "refresh");
        return generateToken(claims, refreshTokenExpiration);
    }

    /**
     * 生成令牌
     */
    private String generateToken(Map<String, Object> claims, Long expiration) {
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + expiration * 1000);

        return Jwts.builder()
                .claims(claims)
                .issuedAt(now)
                .expiration(expiryDate)
                .signWith(getSigningKey())
                .compact();
    }

    /**
     * 从令牌中获取用户ID
     */
    public Long getUserId(String token) {
        return getClaim(token, claims -> claims.get("userId", Long.class));
    }

    /**
     * 从令牌中获取用户名
     */
    public String getUsername(String token) {
        return getClaim(token, Claims::getSubject);
    }

    /**
     * 从令牌中获取设备ID
     */
    public String getDeviceId(String token) {
        return getClaim(token, claims -> claims.get("deviceId", String.class));
    }

    /**
     * 获取令牌类型
     */
    public String getTokenType(String token) {
        return getClaim(token, claims -> claims.get("type", String.class));
    }

    /**
     * 获取过期时间
     */
    public Date getExpiration(String token) {
        return getClaim(token, Claims::getExpiration);
    }

    /**
     * 获取指定声明
     */
    public <T> T getClaim(String token, Function<Claims, T> claimsResolver) {
        Claims claims = getAllClaims(token);
        return claimsResolver.apply(claims);
    }

    /**
     * 获取所有声明
     */
    private Claims getAllClaims(String token) {
        return Jwts.parser()
                .verifyWith(getSigningKey())
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }

    /**
     * 验证令牌是否有效
     */
    public boolean isTokenValid(String token) {
        try {
            return !isTokenExpired(token);
        } catch (Exception e) {
            log.error("JWT令牌验证失败: {}", e.getMessage());
            return false;
        }
    }

    /**
     * 检查令牌是否过期
     */
    private boolean isTokenExpired(String token) {
        return getExpiration(token).before(new Date());
    }

    /**
     * 获取访问令牌过期时间（秒）
     */
    public Long getAccessTokenExpiration() {
        return accessTokenExpiration;
    }

    /**
     * 获取刷新令牌过期时间（秒）
     */
    public Long getRefreshTokenExpiration() {
        return refreshTokenExpiration;
    }
}
