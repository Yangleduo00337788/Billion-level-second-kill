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
@Table("im_group_member")
public class ImGroupMember {

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long groupId;

    private Long userId;

    private Integer role;

    private String nicknameInGroup;

    private Integer isMuted;

    private LocalDateTime mutedUntil;

    private LocalDateTime joinedTime;

    private Long lastReadSeq;
}