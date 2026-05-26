package com.seckill.seckill.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.seckill.dto.SeckillProductVO;
import com.seckill.seckill.entity.SeckillProduct;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

import java.util.List;

@Mapper
public interface SeckillProductMapper extends BaseMapper<SeckillProduct> {

    @Select("SELECT * FROM sk_seckill_product WHERE status = 1 AND start_time <= NOW() AND end_time > NOW() ORDER BY start_time ASC")
    List<SeckillProduct> selectActiveProducts();

    @Select("SELECT * FROM sk_seckill_product WHERE status = 1 AND end_time <= NOW() ORDER BY start_time ASC")
    List<SeckillProduct> selectExpiredProducts();

    @Select("SELECT * FROM sk_seckill_product ORDER BY start_time ASC")
    List<SeckillProduct> selectAllProducts();

    @Select("SELECT * FROM sk_seckill_product WHERE status = 0 AND start_time <= NOW() ORDER BY start_time ASC")
    List<SeckillProduct> selectPendingStartProducts();

    @org.apache.ibatis.annotations.Update("UPDATE sk_seckill_product SET status = 1 WHERE id = #{id} AND status = 0")
    int startSeckill(@Param("id") Long id);

    @org.apache.ibatis.annotations.Update("UPDATE sk_seckill_product SET status = 2 WHERE id = #{id} AND status = 1")
    int endSeckill(@Param("id") Long id);

    @Select("SELECT sp.id, sp.product_id, p.product_name, p.main_image, sp.seckill_price, p.original_price, p.total_stock, sp.seckill_stock, sp.person_limit, sp.start_time, sp.end_time, sp.status " +
            "FROM sk_seckill_product sp " +
            "LEFT JOIN sk_product p ON sp.product_id = p.id " +
            "ORDER BY sp.start_time ASC")
    List<SeckillProductVO> selectAllProductsWithName();
}