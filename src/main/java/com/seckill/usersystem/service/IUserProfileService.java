package com.seckill.usersystem.service;

import com.seckill.usersystem.dto.UserProfileDTO;
import com.seckill.usersystem.entity.UserProfile;

public interface IUserProfileService {

    UserProfile getUserProfileByUserId(Long userId);

    void saveProfile(UserProfileDTO profileDTO);

    void updateProfile(UserProfileDTO profileDTO);

    UserProfile getOrCreateProfile(Long userId);
}
