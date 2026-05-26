package com.seckill.usersystem.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.usersystem.entity.User;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface UserMapper extends BaseMapper<User> {
}
