package com.seckill.usersystem.util;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.*;
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

    private SecretKey getSigningKey() {
        byte[] keyBytes = secret.getBytes(StandardCharsets.UTF_8);
        return Keys.hmacShaKeyFor(keyBytes);
    }

    public String generateAccessToken(Long userId, String username, String deviceId) {
        return generateAccessToken(userId, username, deviceId, null, null, null, null);
    }

    public String generateAccessToken(Long userId, String username, String deviceId,
                                       Set<String> roles, Set<String> permissions,
                                       String sessionId, String ip) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", userId);
        claims.put("username", username);
        claims.put("deviceId", deviceId);
        claims.put("type", "access");
        claims.put("sessionId", sessionId != null ? sessionId : UUID.randomUUID().toString());
        claims.put("ip", ip != null ? ip : "127.0.0.1");
        claims.put("jti", UUID.randomUUID().toString());
        claims.put("loginTime", System.currentTimeMillis());
        if (roles != null && !roles.isEmpty()) {
            claims.put("roles", new ArrayList<>(roles));
        }
        if (permissions != null && !permissions.isEmpty()) {
            claims.put("permissions", new ArrayList<>(permissions));
        }
        return generateToken(claims, accessTokenExpiration);
    }

    public String generateRefreshToken(Long userId, String deviceId) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", userId);
        claims.put("deviceId", deviceId);
        claims.put("type", "refresh");
        claims.put("jti", UUID.randomUUID().toString());
        return generateToken(claims, refreshTokenExpiration);
    }

    private String generateToken(Map<String, Object> claims, Long expiration) {
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + expiration * 1000);

        return Jwts.builder()
                .claims(claims)
                .subject(String.valueOf(claims.get("username") != null ? claims.get("username") : claims.get("userId")))
                .issuedAt(now)
                .expiration(expiryDate)
                .id(claims.get("jti") != null ? claims.get("jti").toString() : null)
                .signWith(getSigningKey())
                .compact();
    }

    public Long getUserId(String token) {
        return getClaim(token, claims -> claims.get("userId", Long.class));
    }

    public String getUsername(String token) {
        String username = getClaim(token, claims -> claims.get("username", String.class));
        return username != null ? username : getClaim(token, Claims::getSubject);
    }

    public String getDeviceId(String token) {
        return getClaim(token, claims -> claims.get("deviceId", String.class));
    }

    public String getTokenType(String token) {
        return getClaim(token, claims -> claims.get("type", String.class));
    }

    public String getSessionId(String token) {
        return getClaim(token, claims -> claims.get("sessionId", String.class));
    }

    public String getJti(String token) {
        return getClaim(token, claims -> claims.get("jti", String.class));
    }

    public String getIp(String token) {
        return getClaim(token, claims -> claims.get("ip", String.class));
    }

    @SuppressWarnings("unchecked")
    public List<String> getRoles(String token) {
        return getClaim(token, claims -> {
            Object rolesObj = claims.get("roles");
            if (rolesObj instanceof List) {
                return (List<String>) rolesObj;
            }
            return Collections.emptyList();
        });
    }

    @SuppressWarnings("unchecked")
    public List<String> getPermissions(String token) {
        return getClaim(token, claims -> {
            Object permsObj = claims.get("permissions");
            if (permsObj instanceof List) {
                return (List<String>) permsObj;
            }
            return Collections.emptyList();
        });
    }

    public Date getExpiration(String token) {
        return getClaim(token, Claims::getExpiration);
    }

    public <T> T getClaim(String token, Function<Claims, T> claimsResolver) {
        Claims claims = getAllClaims(token);
        return claimsResolver.apply(claims);
    }

    private Claims getAllClaims(String token) {
        return Jwts.parser()
                .verifyWith(getSigningKey())
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }

    public boolean isTokenValid(String token) {
        try {
            return !isTokenExpired(token);
        } catch (Exception e) {
            log.error("JWT令牌验证失败: {}", e.getMessage());
            return false;
        }
    }

    private boolean isTokenExpired(String token) {
        return getExpiration(token).before(new Date());
    }

    public Long getAccessTokenExpiration() {
        return accessTokenExpiration;
    }

    public Long getRefreshTokenExpiration() {
        return refreshTokenExpiration;
    }
}
