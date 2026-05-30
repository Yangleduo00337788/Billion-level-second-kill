package com.seckill.seckill.dto;

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
public class SeckillProductVO {

    private Long id;

    private Long productId;

    private String productName;

    private String mainImage;

    private BigDecimal seckillPrice;

    private BigDecimal originalPrice;

    private Integer totalStock;

    private Integer seckillStock;

    private Integer personLimit;

    private LocalDateTime startTime;

    private LocalDateTime endTime;

    private Integer status;
}
