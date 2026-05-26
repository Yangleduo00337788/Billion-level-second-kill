package com.seckill.usersystem.config;

import com.seckill.usersystem.interceptor.JwtAuthInterceptor;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

/**
 * Web MVC 配置
 */
@Configuration
@RequiredArgsConstructor
public class WebConfig implements WebMvcConfigurer {

    private final JwtAuthInterceptor jwtAuthInterceptor;

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(jwtAuthInterceptor)
                .addPathPatterns("/api/**")
                .excludePathPatterns(
                        "/api/auth/login",
                        "/api/auth/register",
                        "/api/auth/refresh-token",
                        "/api/auth/oauth/**",
                        "/api/auth/captcha",
                        "/api/auth/sms-code",
                        "/doc.html",
                        "/swagger-resources/**",
                        "/webjars/**",
                        "/v3/api-docs/**",
                        "/error"
                );
    }
}
