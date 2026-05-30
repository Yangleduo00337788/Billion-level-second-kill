package com.seckill.admin.controller;

import com.seckill.im.entity.*;
import com.seckill.im.mapper.*;
import com.seckill.usersystem.vo.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin/im")
@RequiredArgsConstructor
public class AdminImController {

    private final ImSessionMapper imSessionMapper;
    private final ImMessageMapper imMessageMapper;
    private final ImGroupMapper imGroupMapper;
    private final ImGroupMemberMapper imGroupMemberMapper;
    private final ImReadReceiptMapper imReadReceiptMapper;

    @GetMapping("/sessions")
    public Result<List<ImSession>> listSessions() {
        return Result.success(imSessionMapper.selectAll());
    }

    @GetMapping("/messages")
    public Result<List<ImMessage>> listMessages() {
        return Result.success(imMessageMapper.selectAll());
    }

    @GetMapping("/groups")
    public Result<List<ImGroup>> listGroups() {
        return Result.success(imGroupMapper.selectAll());
    }

    @GetMapping("/groups/{groupId}/members")
    public Result<List<ImGroupMember>> getGroupMembers(@PathVariable Long groupId) {
        return Result.success(imGroupMemberMapper.selectAll());
    }

    @GetMapping("/push-logs")
    public Result<List<ImReadReceipt>> listPushLogs() {
        return Result.success(imReadReceiptMapper.selectAll());
    }
}