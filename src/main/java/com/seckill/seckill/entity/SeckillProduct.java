package com.seckill.seckill.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("sk_seckill_product")
public class SeckillProduct {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long productId;

    private BigDecimal seckillPrice;

    private Integer seckillStock;

    private String redisStockKey;

    private Integer personLimit;

    private LocalDateTime startTime;

    private LocalDateTime endTime;

    private Integer status;

    private Integer version;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;

    public boolean isInProgress() {
        LocalDateTime now = LocalDateTime.now();
        return status == 1 && now.isAfter(startTime) && now.isBefore(endTime);
    }
}