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
@Table("sk_product")
public class Product {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private String productName;

    private String title;

    private String description;

    private String mainImage;

    private String detailImages;

    private BigDecimal originalPrice;

    private BigDecimal price;

    private Integer totalStock;

    private Integer soldCount;

    private Long categoryId;

    private Integer status;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}