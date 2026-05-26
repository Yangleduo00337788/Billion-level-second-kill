package com.seckill.im.entity;

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
@Table("im_read_receipt")
public class ImReadReceipt {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long msgId;

    private Long sessionId;

    private Long userId;

    private LocalDateTime readTime;
}