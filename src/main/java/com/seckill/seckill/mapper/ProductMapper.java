package com.seckill.seckill.mapper;

import com.mybatisflex.core.BaseMapper;
import com.seckill.seckill.entity.Product;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface ProductMapper extends BaseMapper<Product> {
}