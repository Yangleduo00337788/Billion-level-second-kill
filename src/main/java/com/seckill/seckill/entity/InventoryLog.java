package com.seckill.seckill.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("sk_inventory_log")
public class InventoryLog {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long productId;

    private Long seckillId;

    private Integer changeType;

    private Integer changeAmount;

    private Integer beforeStock;

    private Integer afterStock;

    private Long orderId;

    private Long userId;

    private String remark;

    private LocalDateTime createTime;
}