package com.seckill.im.entity;

import com.mybatisflex.annotation.Column;
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

    @Column(ignore = true)
    private Integer maxMembers;

    @Column(ignore = true)
    private String description;

    @Column(ignore = true)
    private String notice;

    private Integer status;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}