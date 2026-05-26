package com.seckill.seckill.entity;

import com.mybatisflex.annotation.Column;
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
@Table("sk_order")
public class SeckillOrder {

    @Id(keyType = KeyType.None)
    private Long id;

    private String orderNo;

    private Long userId;

    private Long productId;

    private Long seckillId;

    private String productName;

    private BigDecimal productPrice;

    private Integer quantity;

    private BigDecimal totalAmount;

    private Integer status;

    private LocalDateTime payTime;

    @Column(onInsertValue = "now()")
    private LocalDateTime createTime;

    @Column(onInsertValue = "now()", onUpdateValue = "now()")
    private LocalDateTime updateTime;
}