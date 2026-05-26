package com.seckill.seckill.controller;

import com.seckill.seckill.dto.SeckillProductVO;
import com.seckill.seckill.dto.SeckillResult;
import com.seckill.seckill.entity.SeckillOrder;
import com.seckill.seckill.limiter.RateLimit;
import com.seckill.seckill.mapper.SeckillOrderMapper;
import com.seckill.seckill.mapper.SeckillProductMapper;
import com.seckill.seckill.entity.SeckillProduct;
import com.seckill.seckill.service.SeckillService;
import com.seckill.seckill.service.SeckillStockService;
import com.seckill.usersystem.util.JwtUtil;
import com.seckill.usersystem.vo.Result;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Slf4j
@RestController
@RequestMapping("/api/seckill")
@RequiredArgsConstructor
public class SeckillController {

    private final SeckillService seckillService;
    private final SeckillStockService stockService;
    private final SeckillProductMapper seckillProductMapper;
    private final SeckillOrderMapper seckillOrderMapper;
    private final JwtUtil jwtUtil;

    @RateLimit(permitsPerSecond = 10000, timeoutMillis = 200)
    @PostMapping("/execute/{seckillId}")
    public Result<SeckillResult> executeSeckill(@PathVariable Long seckillId,
                                                 HttpServletRequest request) {
        Long userId = getUserIdFromRequest(request);
        if (userId == null) {
            return Result.unauthorized("请先登录");
        }

        SeckillResult result = seckillService.executeSeckill(seckillId, userId);
        if (result.isSuccess()) {
            return Result.success(result);
        }
        return Result.fail(result.getCode(), result.getMessage());
    }

    @RateLimit(permitsPerSecond = 50000, timeoutMillis = 100)
    @GetMapping("/product/{seckillId}")
    public Result<Map<String, Object>> getSeckillProduct(@PathVariable Long seckillId) {
        SeckillProduct product = seckillProductMapper.selectOneById(seckillId);
        if (product == null) {
            return Result.fail("秒杀商品不存在");
        }

        int remainingStock = stockService.getRemainingStock(seckillId);

        Map<String, Object> data = new HashMap<>();
        data.put("seckillId", product.getId());
        data.put("productId", product.getProductId());
        data.put("seckillPrice", product.getSeckillPrice());
        data.put("remainingStock", remainingStock);
        data.put("personLimit", product.getPersonLimit());
        data.put("startTime", product.getStartTime());
        data.put("endTime", product.getEndTime());
        data.put("status", product.getStatus());

        return Result.success(data);
    }

    @RateLimit(permitsPerSecond = 50000, timeoutMillis = 100)
    @GetMapping("/products")
    public Result<List<SeckillProductVO>> getActiveProducts() {
        List<SeckillProductVO> products = seckillProductMapper.selectAllProductsWithName();
        return Result.success(products);
    }

    @GetMapping("/order/result/{seckillId}")
    public Result<SeckillResult> getOrderResult(@PathVariable Long seckillId,
                                                 HttpServletRequest request) {
        Long userId = getUserIdFromRequest(request);
        if (userId == null) {
            return Result.unauthorized("请先登录");
        }

        SeckillOrder order = seckillOrderMapper.selectByUserIdAndSeckillId(userId, seckillId);
        if (order != null) {
            return Result.success(SeckillResult.success(order.getOrderNo(), order.getId()));
        }

        if (stockService.getRemainingStock(seckillId) <= 0) {
            return Result.success(SeckillResult.stockOut());
        }
        return Result.success(SeckillResult.builder()
                .success(false)
                .message("排队中...")
                .code(100)
                .build());
    }

    @RateLimit(permitsPerSecond = 50000, timeoutMillis = 100)
    @GetMapping("/stock/{seckillId}")
    public Result<Integer> getStock(@PathVariable Long seckillId) {
        return Result.success(stockService.getRemainingStock(seckillId));
    }

    @GetMapping("/my-orders")
    public Result<List<SeckillOrder>> getMyOrders(HttpServletRequest request) {
        Long userId = getUserIdFromRequest(request);
        if (userId == null) {
            return Result.unauthorized("请先登录");
        }
        List<SeckillOrder> orders = seckillOrderMapper.selectByUserId(userId);
        return Result.success(orders);
    }

    @PostMapping("/pay/{orderNo}")
    public Result<String> payOrder(@PathVariable String orderNo,
                                   HttpServletRequest request) {
        Long userId = getUserIdFromRequest(request);
        if (userId == null) {
            return Result.unauthorized("请先登录");
        }

        SeckillOrder order = seckillOrderMapper.selectByOrderNo(orderNo);
        if (order == null) {
            return Result.fail("订单不存在");
        }
        if (!order.getUserId().equals(userId)) {
            return Result.fail("无权操作此订单");
        }
        if (order.getStatus() != 0) {
            return Result.fail("订单状态不正确");
        }

        SeckillProduct product = seckillProductMapper.selectOneById(order.getSeckillId());
        if (product == null || product.getStatus() == 2) {
            order.setStatus(2);
            seckillOrderMapper.update(order);
            return Result.fail("秒杀活动已结束，订单已自动取消");
        }

        order.setStatus(1);
        order.setPayTime(java.time.LocalDateTime.now());
        seckillOrderMapper.update(order);

        log.info("[SeckillController] order paid: orderNo={}, userId={}", orderNo, userId);
        return Result.success("支付成功");
    }

    @PostMapping("/cancel/{orderNo}")
    public Result<String> cancelOrder(@PathVariable String orderNo,
                                      HttpServletRequest request) {
        Long userId = getUserIdFromRequest(request);
        if (userId == null) {
            return Result.unauthorized("请先登录");
        }

        SeckillOrder order = seckillOrderMapper.selectByOrderNo(orderNo);
        if (order == null) {
            return Result.fail("订单不存在");
        }
        if (!order.getUserId().equals(userId)) {
            return Result.fail("无权操作此订单");
        }
        if (order.getStatus() != 0) {
            return Result.fail("订单状态不正确");
        }

        order.setStatus(2);
        seckillOrderMapper.update(order);

        log.info("[SeckillController] order cancelled: orderNo={}, userId={}", orderNo, userId);
        return Result.success("取消成功");
    }

    private Long getUserIdFromRequest(HttpServletRequest request) {
        String authHeader = request.getHeader("Authorization");
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            return null;
        }
        try {
            String token = authHeader.substring(7);
            return jwtUtil.getUserId(token);
        } catch (Exception e) {
            return null;
        }
    }
}