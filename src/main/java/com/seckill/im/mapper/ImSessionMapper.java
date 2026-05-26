package com.seckill.im.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.im.entity.ImSession;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

@Mapper
public interface ImSessionMapper extends BaseMapper<ImSession> {

    @Select("SELECT s.* FROM im_session s " +
            "INNER JOIN im_session_user su ON s.id = su.session_id " +
            "WHERE su.user_id = #{userId} AND s.status = 1 " +
            "ORDER BY s.last_msg_time DESC")
    java.util.List<ImSession> selectByUserId(@Param("userId") Long userId);

    @Select("SELECT s.* FROM im_session s " +
            "INNER JOIN im_session_user su1 ON s.id = su1.session_id " +
            "INNER JOIN im_session_user su2 ON s.id = su2.session_id " +
            "WHERE s.session_type = 1 AND su1.user_id = #{userId1} AND su2.user_id = #{userId2}")
    ImSession selectSingleSessionByUserIds(@Param("userId1") Long userId1, @Param("userId2") Long userId2);
}