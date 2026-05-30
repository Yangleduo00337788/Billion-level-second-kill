package com.seckill.usersystem.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.usersystem.entity.UserProfile;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface UserProfileMapper extends BaseMapper<UserProfile> {
}
