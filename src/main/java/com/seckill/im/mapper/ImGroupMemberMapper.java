package com.seckill.im.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.im.entity.ImGroupMember;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

import java.util.List;

@Mapper
public interface ImGroupMemberMapper extends BaseMapper<ImGroupMember> {

    @Select("SELECT user_id FROM im_group_member WHERE group_id = #{groupId}")
    List<Long> selectMemberUserIdsByGroupId(@Param("groupId") Long groupId);
}