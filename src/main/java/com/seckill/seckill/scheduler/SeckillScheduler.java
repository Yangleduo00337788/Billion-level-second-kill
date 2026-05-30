package com.seckill.seckill.scheduler;

import com.seckill.seckill.entity.SeckillOrder;
import com.seckill.seckill.entity.SeckillProduct;
import com.seckill.seckill.mapper.SeckillOrderMapper;
import com.seckill.seckill.mapper.SeckillProductMapper;
import com.seckill.seckill.service.SeckillStockService;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.util.List;

@Slf4j
@Component
@RequiredArgsConstructor
public class SeckillScheduler {

    private final SeckillProductMapper seckillProductMapper;
    private final SeckillStockService stockService;
    private final SeckillOrderMapper seckillOrderMapper;

    @PostConstruct
    public void init() {
        log.info("[SeckillScheduler] warming up stock for active seckill products...");
        List<SeckillProduct> activeProducts = seckillProductMapper.selectActiveProducts();
        for (SeckillProduct product : activeProducts) {
            try {
                stockService.warmUpStock(product.getId(), product.getSeckillStock());
                log.info("[SeckillScheduler] warmed up: id={}, stock={}", product.getId(), product.getSeckillStock());
            } catch (Exception e) {
                log.error("[SeckillScheduler] failed to warm up: id={}", product.getId(), e);
            }
        }
    }

    @Scheduled(cron = "*/5 * * * * ?")
    public void checkAndStartSeckill() {
        List<SeckillProduct> pendingProducts = seckillProductMapper.selectPendingStartProducts();
        for (SeckillProduct product : pendingProducts) {
            try {
                seckillProductMapper.startSeckill(product.getId());
                stockService.warmUpStock(product.getId(), product.getSeckillStock());
                log.info("[SeckillScheduler] started seckill: id={}, stock={}",
                        product.getId(), product.getSeckillStock());
            } catch (Exception e) {
                log.error("[SeckillScheduler] failed to start seckill: id={}", product.getId(), e);
            }
        }
    }

    @Scheduled(cron = "*/10 * * * * ?")
    public void checkAndEndSeckill() {
        List<SeckillProduct> expiredProducts = seckillProductMapper.selectExpiredProducts();
        for (SeckillProduct product : expiredProducts) {
            try {
                seckillProductMapper.endSeckill(product.getId());
                log.info("[SeckillScheduler] ended seckill: id={}", product.getId());
            } catch (Exception e) {
                log.error("[SeckillScheduler] failed to end seckill: id={}", product.getId(), e);
            }
        }
    }

    @Scheduled(cron = "*/60 * * * * ?")
    public void cancelUnpaidOrders() {
        List<SeckillOrder> unpaidOrders = seckillOrderMapper.selectUnpaidOrders();
        for (SeckillOrder order : unpaidOrders) {
            try {
                order.setStatus(2);
                seckillOrderMapper.update(order);
                log.info("[SeckillScheduler] auto cancelled unpaid order: orderNo={}", order.getOrderNo());
            } catch (Exception e) {
                log.error("[SeckillScheduler] failed to cancel order: orderNo={}", order.getOrderNo(), e);
            }
        }
    }
}