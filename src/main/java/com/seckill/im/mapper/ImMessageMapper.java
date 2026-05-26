package com.seckill.im.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.im.entity.ImMessage;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;

import java.util.List;

@Mapper
public interface ImMessageMapper extends BaseMapper<ImMessage> {

    @Select("SELECT * FROM im_message WHERE session_id = #{sessionId} AND seq > #{lastSeq} ORDER BY seq ASC LIMIT #{limit}")
    List<ImMessage> selectUnreadMessages(@Param("sessionId") Long sessionId,
                                         @Param("lastSeq") Long lastSeq,
                                         @Param("limit") int limit);

    @Select("SELECT MAX(seq) FROM im_message WHERE session_id = #{sessionId}")
    Long selectMaxSeqBySessionId(@Param("sessionId") Long sessionId);

    @Select("SELECT * FROM im_message WHERE client_msg_id = #{clientMsgId}")
    ImMessage selectByClientMsgId(@Param("clientMsgId") String clientMsgId);

    @Update("UPDATE im_message SET status = 2, recall_time = NOW() WHERE id = #{msgId} AND sender_id = #{userId} AND status = 1")
    int recallMessage(@Param("msgId") Long msgId, @Param("userId") Long userId);

    @Select("SELECT * FROM im_message WHERE session_id = #{sessionId} ORDER BY seq DESC LIMIT #{limit}")
    List<ImMessage> selectLatestMessages(@Param("sessionId") Long sessionId, @Param("limit") int limit);
}