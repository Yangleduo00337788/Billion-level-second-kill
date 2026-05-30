package com.seckill.im.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.im.entity.ImSessionUser;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Update;

@Mapper
public interface ImSessionUserMapper extends BaseMapper<ImSessionUser> {

    @Update("UPDATE im_session_user SET last_read_msg_seq = #{lastReadSeq} " +
            "WHERE session_id = #{sessionId} AND user_id = #{userId} AND last_read_msg_seq < #{lastReadSeq}")
    int updateLastReadSeq(@Param("sessionId") Long sessionId,
                          @Param("userId") Long userId,
                          @Param("lastReadSeq") Long lastReadSeq);
}