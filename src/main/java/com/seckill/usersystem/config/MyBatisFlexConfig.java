package com.seckill.usersystem.config;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@MapperScan({"com.seckill.usersystem.mapper", "com.seckill.im.mapper", "com.seckill.seckill.mapper"})
public class MyBatisFlexConfig {
}
