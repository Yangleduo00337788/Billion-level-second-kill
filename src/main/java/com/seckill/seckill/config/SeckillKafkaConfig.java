package com.seckill.seckill.config;

import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Slf4j
@Configuration
public class SeckillKafkaConfig {

    public static final String TOPIC_SECKILL_ORDER = "seckill-order-topic";
    public static final String TOPIC_SECKILL_STOCK_SYNC = "seckill-stock-sync-topic";
    public static final String TOPIC_SECKILL_CANCEL = "seckill-cancel-topic";

    @Bean
    public NewTopic seckillOrderTopic() {
        return new NewTopic(TOPIC_SECKILL_ORDER, 16, (short) 1);
    }

    @Bean
    public NewTopic seckillStockSyncTopic() {
        return new NewTopic(TOPIC_SECKILL_STOCK_SYNC, 8, (short) 1);
    }

    @Bean
    public NewTopic seckillCancelTopic() {
        return new NewTopic(TOPIC_SECKILL_CANCEL, 4, (short) 1);
    }
}