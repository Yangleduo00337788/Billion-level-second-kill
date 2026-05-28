package com.seckill.admin.controller;

import com.seckill.seckill.entity.*;
import com.seckill.seckill.mapper.*;
import com.seckill.usersystem.vo.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin/seckill")
@RequiredArgsConstructor
public class AdminSeckillController {

    private final ProductMapper productMapper;
    private final SeckillProductMapper seckillProductMapper;
    private final SeckillOrderMapper seckillOrderMapper;
    private final InventoryLogMapper inventoryLogMapper;

    @GetMapping("/products")
    public Result<List<Product>> listProducts() {
        return Result.success(productMapper.selectAll());
    }

    @PostMapping("/products")
    public Result<Void> saveProduct(@RequestBody Product product) {
        if (product.getId() == null) {
            productMapper.insert(product);
        } else {
            productMapper.update(product);
        }
        return Result.success();
    }

    @DeleteMapping("/products/{id}")
    public Result<Void> deleteProduct(@PathVariable Long id) {
        productMapper.deleteById(id);
        return Result.success();
    }

    @GetMapping("/activities")
    public Result<List<SeckillProduct>> listActivities() {
        return Result.success(seckillProductMapper.selectAll());
    }

    @PostMapping("/activities")
    public Result<Void> saveActivity(@RequestBody SeckillProduct activity) {
        if (activity.getId() == null) {
            seckillProductMapper.insert(activity);
        } else {
            seckillProductMapper.update(activity);
        }
        return Result.success();
    }

    @DeleteMapping("/activities/{id}")
    public Result<Void> deleteActivity(@PathVariable Long id) {
        seckillProductMapper.deleteById(id);
        return Result.success();
    }

    @GetMapping("/orders")
    public Result<List<SeckillOrder>> listOrders() {
        return Result.success(seckillOrderMapper.selectAll());
    }

    @GetMapping("/stock-logs")
    public Result<List<InventoryLog>> listStockLogs() {
        return Result.success(inventoryLogMapper.selectAll());
    }
}