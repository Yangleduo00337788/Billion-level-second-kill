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
@Table("im_message")
public class ImMessage {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private String clientMsgId;

    private Long sessionId;

    private Long senderId;

    private Integer msgType;

    private String content;

    private String mediaUrl;

    private String mediaThumbUrl;

    private Integer mediaDuration;

    private Long mediaSize;

    private Long seq;

    private Integer status;

    private LocalDateTime recallTime;

    private LocalDateTime createTime;
}