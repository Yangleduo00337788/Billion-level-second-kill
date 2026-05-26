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
@Table("im_session")
public class ImSession {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Integer sessionType;

    private Long groupId;

    private Long lastMsgId;

    private String lastMsgContent;

    private LocalDateTime lastMsgTime;

    private Long lastMsgSenderId;

    private Integer status;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}