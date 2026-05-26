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
@Table("im_group")
public class ImGroup {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private String groupName;

    private String avatar;

    private Long ownerId;

    private Integer memberCount;

    private Integer maxMembers;

    private String description;

    private String notice;

    private Integer status;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}