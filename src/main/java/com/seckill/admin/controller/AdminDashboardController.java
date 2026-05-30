package com.seckill.admin.controller;

import com.seckill.im.mapper.*;
import com.seckill.seckill.mapper.*;
import com.seckill.usersystem.mapper.*;
import com.seckill.usersystem.vo.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.*;

@RestController
@RequestMapping("/api/admin")
@RequiredArgsConstructor
public class AdminDashboardController {

    private final UserMapper userMapper;
    private final RoleMapper roleMapper;
    private final BlacklistMapper blacklistMapper;
    private final ProductMapper productMapper;
    private final SeckillProductMapper seckillProductMapper;
    private final SeckillOrderMapper seckillOrderMapper;
    private final ImSessionMapper imSessionMapper;
    private final ImMessageMapper imMessageMapper;
    private final ImGroupMapper imGroupMapper;

    @GetMapping("/dashboard")
    public Result<Map<String, Object>> dashboard() {
        Map<String, Object> data = new HashMap<>();

        Map<String, Object> userStats = new HashMap<>();
        userStats.put("userCount", userMapper.selectAll().size());
        userStats.put("roleCount", roleMapper.selectAll().size());
        userStats.put("blacklistCount", blacklistMapper.selectAll().size());
        data.put("user", userStats);

        Map<String, Object> seckillStats = new HashMap<>();
        seckillStats.put("productCount", productMapper.selectAll().size());
        seckillStats.put("activityCount", seckillProductMapper.selectAll().size());
        seckillStats.put("todayOrders", seckillOrderMapper.selectAll().size());
        seckillStats.put("qpsLimit", 10000);
        data.put("seckill", seckillStats);

        Map<String, Object> imStats = new HashMap<>();
        imStats.put("onlineUsers", 0);
        imStats.put("activeSessions", imSessionMapper.selectAll().size());
        imStats.put("todayMessages", imMessageMapper.selectAll().size());
        imStats.put("groupCount", imGroupMapper.selectAll().size());
        data.put("im", imStats);

        return Result.success(data);
    }
}