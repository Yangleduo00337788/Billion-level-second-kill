package com.seckill.usersystem.annotation;

import com.seckill.usersystem.util.JwtUtil;
import com.seckill.usersystem.vo.Result;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.stereotype.Component;

import java.lang.reflect.Method;
import java.util.List;

@Slf4j
@Aspect
@Component
@RequiredArgsConstructor
public class PermissionAspect {

    private final JwtUtil jwtUtil;

    @Around("@annotation(com.seckill.usersystem.annotation.RequirePermission) || " +
            "@within(com.seckill.usersystem.annotation.RequirePermission)")
    public Object around(ProceedingJoinPoint joinPoint) throws Throwable {
        HttpServletRequest request = getRequest(joinPoint);
        if (request == null) {
            return joinPoint.proceed();
        }

        String authHeader = request.getHeader("Authorization");
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            return Result.unauthorized("未提供认证令牌");
        }

        String token = authHeader.substring(7);
        List<String> userPermissions = jwtUtil.getPermissions(token);
        List<String> userRoles = jwtUtil.getRoles(token);

        RequirePermission annotation = getAnnotation(joinPoint);
        if (annotation == null) {
            return joinPoint.proceed();
        }

        boolean hasPermission = checkPermissions(annotation, userPermissions, userRoles);
        if (!hasPermission) {
            log.warn("权限不足: method={}, requiredPermissions={}, requiredRoles={}, userRoles={}",
                    joinPoint.getSignature().toShortString(),
                    annotation.value(), annotation.roles(), userRoles);
            return Result.forbidden("权限不足");
        }

        return joinPoint.proceed();
    }

    private boolean checkPermissions(RequirePermission annotation,
                                      List<String> userPermissions,
                                      List<String> userRoles) {
        boolean hasRoles = true;
        if (annotation.roles().length > 0) {
            hasRoles = checkValues(annotation.roles(), userRoles, annotation.logical());
        }

        boolean hasPerms = true;
        if (annotation.value().length > 0) {
            hasPerms = checkValues(annotation.value(), userPermissions, annotation.logical());
        }

        if (annotation.logical() == RequirePermission.Logical.AND) {
            return hasRoles && hasPerms;
        }
        return hasRoles || hasPerms;
    }

    private boolean checkValues(String[] required, List<String> userValues, RequirePermission.Logical logical) {
        if (logical == RequirePermission.Logical.AND) {
            for (String requiredValue : required) {
                if (!userValues.contains(requiredValue)) {
                    return false;
                }
            }
            return true;
        }
        for (String requiredValue : required) {
            if (userValues.contains(requiredValue)) {
                return true;
            }
        }
        return false;
    }

    private RequirePermission getAnnotation(ProceedingJoinPoint joinPoint) {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        Method method = signature.getMethod();
        RequirePermission methodAnnotation = method.getAnnotation(RequirePermission.class);
        if (methodAnnotation != null) {
            return methodAnnotation;
        }
        return method.getDeclaringClass().getAnnotation(RequirePermission.class);
    }

    private HttpServletRequest getRequest(ProceedingJoinPoint joinPoint) {
        Object[] args = joinPoint.getArgs();
        for (Object arg : args) {
            if (arg instanceof HttpServletRequest) {
                return (HttpServletRequest) arg;
            }
        }
        return null;
    }
}