package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.dto.UserProfileDTO;
import com.seckill.usersystem.entity.UserProfile;
import com.seckill.usersystem.mapper.UserProfileMapper;
import com.seckill.usersystem.service.IUserProfileService;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
public class UserProfileServiceImpl implements IUserProfileService {

    private final UserProfileMapper userProfileMapper;

    @Override
    public UserProfile getUserProfileByUserId(Long userId) {
        return QueryChain.of(userProfileMapper)
                .where(UserProfile::getUserId).eq(userId)
                .one();
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void saveProfile(UserProfileDTO profileDTO) {
        UserProfile profile = new UserProfile();
        BeanUtils.copyProperties(profileDTO, profile);
        userProfileMapper.insert(profile);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void updateProfile(UserProfileDTO profileDTO) {
        UserProfile profile = getUserProfileByUserId(profileDTO.getUserId());
        if (profile != null) {
            profileDTO.setId(profile.getId());
        }
        UserProfile entity = new UserProfile();
        BeanUtils.copyProperties(profileDTO, entity);
        if (profile != null) {
            userProfileMapper.update(entity);
        } else {
            userProfileMapper.insert(entity);
        }
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public UserProfile getOrCreateProfile(Long userId) {
        UserProfile profile = getUserProfileByUserId(userId);
        if (profile == null) {
            profile = new UserProfile();
            profile.setUserId(userId);
            userProfileMapper.insert(profile);
        }
        return profile;
    }
}
