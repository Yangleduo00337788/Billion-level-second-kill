package com.seckill.seckill.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.seckill.entity.SeckillOrder;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;

@Mapper
public interface SeckillOrderMapper extends BaseMapper<SeckillOrder> {

    @Select("SELECT * FROM sk_order WHERE user_id = #{userId} AND seckill_id = #{seckillId}")
    SeckillOrder selectByUserIdAndSeckillId(@Param("userId") Long userId, @Param("seckillId") Long seckillId);

    @Select("SELECT * FROM sk_order WHERE user_id = #{userId} ORDER BY create_time DESC")
    java.util.List<SeckillOrder> selectByUserId(@Param("userId") Long userId);

    @Select("SELECT * FROM sk_order WHERE order_no = #{orderNo}")
    SeckillOrder selectByOrderNo(@Param("orderNo") String orderNo);

    @Update("UPDATE sk_order SET status = 2 WHERE id = #{orderId} AND status = 0")
    int cancelOrder(@Param("orderId") Long orderId);

    @Select("SELECT * FROM sk_order WHERE status = 0 AND create_time < DATE_SUB(NOW(), INTERVAL 15 MINUTE)")
    java.util.List<SeckillOrder> selectUnpaidOrders();
}