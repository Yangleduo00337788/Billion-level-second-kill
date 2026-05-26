package com.seckill.seckill.consumer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.seckill.seckill.config.SeckillKafkaConfig;
import com.seckill.seckill.dto.OrderMessage;
import com.seckill.seckill.entity.SeckillOrder;
import com.seckill.seckill.mapper.SeckillOrderMapper;
import com.seckill.seckill.service.SeckillStockService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.Acknowledgment;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.List;

@Slf4j
@Component
@RequiredArgsConstructor
public class SeckillOrderConsumer {

    private final SeckillOrderMapper seckillOrderMapper;
    private final SeckillStockService stockService;
    private final ObjectMapper objectMapper;

    @KafkaListener(topics = SeckillKafkaConfig.TOPIC_SECKILL_ORDER, groupId = "seckill-order-consumer", concurrency = "8")
    public void consumeOrder(List<byte[]> messages, Acknowledgment ack) {
        for (byte[] data : messages) {
            try {
                OrderMessage msg = objectMapper.readValue(data, OrderMessage.class);
                persistOrder(msg);
            } catch (Exception e) {
                log.error("[SeckillOrderConsumer] consume error", e);
            }
        }
        ack.acknowledge();
    }

    @Transactional(rollbackFor = Exception.class)
    public void persistOrder(OrderMessage msg) {
        SeckillOrder existing = seckillOrderMapper.selectByUserIdAndSeckillId(
                msg.getUserId(), msg.getSeckillId());
        if (existing != null) {
            log.warn("[SeckillOrderConsumer] duplicate order ignored: orderNo={}", msg.getOrderNo());
            return;
        }

        SeckillOrder order = SeckillOrder.builder()
                .id(msg.getOrderId())
                .orderNo(msg.getOrderNo())
                .userId(msg.getUserId())
                .productId(msg.getProductId())
                .seckillId(msg.getSeckillId())
                .productName(msg.getProductName())
                .productPrice(msg.getProductPrice())
                .quantity(msg.getQuantity())
                .totalAmount(msg.getTotalAmount())
                .status(0)
                .createTime(LocalDateTime.ofInstant(
                        Instant.ofEpochMilli(msg.getCreateTime()), ZoneId.systemDefault()))
                .build();

        seckillOrderMapper.insert(order);
        log.info("[SeckillOrderConsumer] order persisted: orderNo={}, userId={}",
                msg.getOrderNo(), msg.getUserId());
    }
}