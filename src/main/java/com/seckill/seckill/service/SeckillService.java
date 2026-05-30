package com.seckill.seckill.service;

import cn.hutool.core.util.IdUtil;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.seckill.im.util.SnowflakeIdGenerator;
import com.seckill.seckill.config.SeckillKafkaConfig;
import com.seckill.seckill.dto.OrderMessage;
import com.seckill.seckill.dto.SeckillResult;
import com.seckill.seckill.entity.SeckillProduct;
import com.seckill.seckill.entity.Product;
import com.seckill.seckill.limiter.TokenBucketLimiter;
import com.seckill.seckill.mapper.ProductMapper;
import com.seckill.seckill.mapper.SeckillProductMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Slf4j
@Service
@RequiredArgsConstructor
public class SeckillService {

    private final SeckillProductMapper seckillProductMapper;
    private final ProductMapper productMapper;
    private final SeckillStockService stockService;
    private final TokenBucketLimiter tokenBucketLimiter;
    private final KafkaTemplate<String, byte[]> kafkaTemplate;
    private final SnowflakeIdGenerator snowflakeIdGenerator;
    private final ObjectMapper objectMapper;

    public SeckillResult executeSeckill(Long seckillId, Long userId) {
        if (!tokenBucketLimiter.tryAcquireByUser(userId, 5)) {
            log.warn("[SeckillService] user:{} rate limited for seckill:{}", userId, seckillId);
            return SeckillResult.rateLimited();
        }

        SeckillProduct seckillProduct = seckillProductMapper.selectOneById(seckillId);
        if (seckillProduct == null) {
            return SeckillResult.ended();
        }
        if (!seckillProduct.isInProgress()) {
            LocalDateTime now = LocalDateTime.now();
            if (now.isBefore(seckillProduct.getStartTime())) {
                return SeckillResult.notStarted();
            }
            return SeckillResult.ended();
        }

        if (stockService.hasOrderFlag(userId, seckillId)) {
            return SeckillResult.duplicate();
        }

        Long luaResult = stockService.deductStock(
                seckillId, userId, seckillProduct.getPersonLimit());

        if (luaResult == null) {
            log.error("[SeckillService] Lua script returned null");
            return SeckillResult.stockOut();
        }

        int resultCode = luaResult.intValue();

        if (resultCode == -1) {
            return SeckillResult.stockOut();
        }
        if (resultCode == -2) {
            return SeckillResult.duplicate();
        }
        if (resultCode == 0) {
            return createOrderAsync(seckillProduct, userId);
        }

        log.error("[SeckillService] unknown Lua result: {}", resultCode);
        return SeckillResult.stockOut();
    }

    private SeckillResult createOrderAsync(SeckillProduct seckillProduct, Long userId) {
        Product product = productMapper.selectOneById(seckillProduct.getProductId());
        if (product == null) {
            stockService.rollbackStock(seckillProduct.getId(), userId);
            return SeckillResult.stockOut();
        }

        long orderId = snowflakeIdGenerator.nextId();
        String orderNo = generateOrderNo(userId);

        OrderMessage msg = OrderMessage.builder()
                .orderId(orderId)
                .orderNo(orderNo)
                .userId(userId)
                .productId(product.getId())
                .seckillId(seckillProduct.getId())
                .productName(product.getProductName())
                .productPrice(seckillProduct.getSeckillPrice())
                .quantity(1)
                .totalAmount(seckillProduct.getSeckillPrice())
                .createTime(System.currentTimeMillis())
                .build();

        try {
            byte[] data = objectMapper.writeValueAsBytes(msg);
            kafkaTemplate.send(SeckillKafkaConfig.TOPIC_SECKILL_ORDER,
                    String.valueOf(seckillProduct.getId()), data);
            log.info("[SeckillService] order sent to Kafka: orderNo={}, userId={}, seckillId={}",
                    orderNo, userId, seckillProduct.getId());
        } catch (JsonProcessingException e) {
            log.error("[SeckillService] serialize OrderMessage failed", e);
            stockService.rollbackStock(seckillProduct.getId(), userId);
            return SeckillResult.stockOut();
        }

        return SeckillResult.success(orderNo, orderId);
    }

    private String generateOrderNo(Long userId) {
        return "SK" + LocalDateTime.now().format(java.time.format.DateTimeFormatter.ofPattern("yyyyMMddHHmmss"))
                + String.format("%06d", userId % 1000000);
    }
}