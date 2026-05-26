package com.seckill.usersystem.service;

import com.seckill.usersystem.dto.LoginRequest;
import com.seckill.usersystem.dto.RegisterRequest;
import com.seckill.usersystem.dto.TokenResponse;
import com.seckill.usersystem.entity.User;
import com.seckill.usersystem.vo.LoginResultVO;
import com.seckill.usersystem.vo.UserVO;

import java.util.List;
import java.util.Set;

public interface IUserService {

    LoginResultVO login(LoginRequest request);

    LoginResultVO register(RegisterRequest request);

    void logout(Long userId, String deviceId);

    TokenResponse refreshToken(String refreshToken, String deviceId);

    User getUserByUsername(String username);

    User getUserById(Long userId);

    List<UserVO> listUsers(int pageNum, int pageSize);

    void updateUserStatus(Long userId, Integer status);

    void deleteUser(Long userId);

    void lockUser(Long userId, long lockMinutes);

    void resetPassword(Long userId, String newPassword);

    Set<String> getUserRoles(Long userId);

    Set<String> getUserPermissions(Long userId);

    boolean isPasswordCorrect(User user, String password);

    void updateLastLoginInfo(Long userId, String ip);
}
