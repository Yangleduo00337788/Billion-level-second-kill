package com.seckill.usersystem.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.io.Serial;
import java.io.Serializable;
import java.time.LocalDate;

/**
 * 用户资料DTO
 */
@Data
public class UserProfileDTO implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private Long id;

    private Long userId;

    @Size(max = 50, message = "昵称长度不能超过50")
    private String nickname;

    @Size(max = 50, message = "真实姓名长度不能超过50")
    private String realName;

    private Integer gender;

    private LocalDate birthday;

    private String avatar;

    @Size(max = 255, message = "地址长度不能超过255")
    private String address;

    private String idCard;

    @Size(max = 500, message = "个人简介长度不能超过500")
    private String bio;
}
