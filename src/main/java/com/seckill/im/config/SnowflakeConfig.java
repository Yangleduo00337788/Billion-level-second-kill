package com.seckill.im.config;

import com.seckill.im.util.SnowflakeIdGenerator;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SnowflakeConfig {

    @Value("${im.snowflake.worker-id:1}")
    private long workerId;

    @Value("${im.snowflake.datacenter-id:1}")
    private long datacenterId;

    @Bean
    public SnowflakeIdGenerator snowflakeIdGenerator() {
        return new SnowflakeIdGenerator(workerId, datacenterId);
    }
}