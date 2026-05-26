package com.seckill.usersystem.entity;

import com.mybatisflex.annotation.Id;
import com.mybatisflex.annotation.KeyType;
import com.mybatisflex.annotation.Table;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDateTime;

@Data
@Table("sys_user_oauth")
public class UserOauth implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    @Id(keyType = KeyType.Auto)
    private Long id;

    private Long userId;

    /** OAuth类型：WECHAT/ALIPAY/QQ/GOOGLE/GITHUB */
    private String oauthType;

    private String openid;

    private String unionid;

    private String accessToken;

    private String refreshToken;

    private Integer expiresIn;

    private String nickname;

    private String avatar;

    private LocalDateTime createTime;

    private LocalDateTime updateTime;
}
