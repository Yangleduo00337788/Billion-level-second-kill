package com.seckill.seckill.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SeckillResult {

    private boolean success;
    private String orderNo;
    private Long orderId;
    private String message;
    private int code;

    public static SeckillResult success(String orderNo, Long orderId) {
        return SeckillResult.builder()
                .success(true)
                .orderNo(orderNo)
                .orderId(orderId)
                .message("秒杀成功，请尽快支付")
                .code(0)
                .build();
    }

    public static SeckillResult stockOut() {
        return SeckillResult.builder()
                .success(false)
                .message("商品已售罄")
                .code(-1)
                .build();
    }

    public static SeckillResult duplicate() {
        return SeckillResult.builder()
                .success(false)
                .message("您已参与过该秒杀")
                .code(-2)
                .build();
    }

    public static SeckillResult notStarted() {
        return SeckillResult.builder()
                .success(false)
                .message("秒杀尚未开始")
                .code(-3)
                .build();
    }

    public static SeckillResult ended() {
        return SeckillResult.builder()
                .success(false)
                .message("秒杀已结束")
                .code(-4)
                .build();
    }

    public static SeckillResult rateLimited() {
        return SeckillResult.builder()
                .success(false)
                .message("请求过于频繁，请稍后再试")
                .code(-5)
                .build();
    }
}