package com.seckill.usersystem.service.impl;

import com.mybatisflex.core.query.QueryChain;
import com.seckill.usersystem.dto.UserProfileDTO;
import com.seckill.usersystem.entity.User;
import com.seckill.usersystem.entity.UserProfile;
import com.seckill.usersystem.mapper.UserMapper;
import com.seckill.usersystem.mapper.UserProfileMapper;
import com.seckill.usersystem.service.IUserProfileService;
import com.seckill.usersystem.util.RedisUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashMap;
import java.util.Map;

@Slf4j
@Service
@RequiredArgsConstructor
public class UserProfileServiceImpl implements IUserProfileService {

    private final UserProfileMapper userProfileMapper;
    private final UserMapper userMapper;
    private final RedisUtil redisUtil;

    private static final String PROFILE_CACHE_KEY = "user:profile:";
    private static final long PROFILE_CACHE_TTL = 7 * 24 * 3600;

    @Override
    public UserProfile getUserProfileByUserId(Long userId) {
        UserProfile cached = getCachedProfile(userId);
        if (cached != null) {
            return cached;
        }

        UserProfile profile = QueryChain.of(userProfileMapper)
                .where(UserProfile::getUserId).eq(userId)
                .one();

        if (profile != null) {
            cacheProfile(profile);
        }
        return profile;
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void saveProfile(UserProfileDTO profileDTO) {
        UserProfile profile = new UserProfile();
        BeanUtils.copyProperties(profileDTO, profile);
        userProfileMapper.insert(profile);
        cacheProfile(profile);
    }

    @Override
    @Transactional(rollbackFor = Exception.class)
    public void updateProfile(UserProfileDTO profileDTO) {
        UserProfile existing = QueryChain.of(userProfileMapper)
                .where(UserProfile::getUserId).eq(profileDTO.getUserId())
                .one();
        if (existing != null) {
            profileDTO.setId(existing.getId());
        }
        UserProfile entity = new UserProfile();
        BeanUtils.copyProperties(profileDTO, entity);
        if (existing != null) {
            userProfileMapper.update(entity);
        } else {
            userProfileMapper.insert(entity);
        }
        if (profileDTO.getEmail() != null || profileDTO.getPhone() != null) {
            User user = QueryChain.of(userMapper)
                    .where(User::getId).eq(profileDTO.getUserId())
                    .one();
            if (user != null) {
                if (profileDTO.getEmail() != null) {
                    user.setEmail(profileDTO.getEmail());
                }
                if (profileDTO.getPhone() != null) {
                    user.setPhone(profileDTO.getPhone());
                }
                userMapper.update(user);
            }
        }
        invalidateProfileCache(profileDTO.getUserId());
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
        User user = QueryChain.of(userMapper)
                .where(User::getId).eq(userId)
                .one();
        if (user != null) {
            profile.setEmail(user.getEmail());
            profile.setPhone(user.getPhone());
        }
        return profile;
    }

    private void cacheProfile(UserProfile profile) {
        if (profile == null || profile.getUserId() == null) {
            return;
        }
        String key = PROFILE_CACHE_KEY + profile.getUserId();
        try {
            Map<String, String> hashData = new HashMap<>();
            if (profile.getId() != null) hashData.put("id", String.valueOf(profile.getId()));
            if (profile.getNickname() != null) hashData.put("nickname", profile.getNickname());
            if (profile.getRealName() != null) hashData.put("realName", profile.getRealName());
            if (profile.getGender() != null) hashData.put("gender", String.valueOf(profile.getGender()));
            if (profile.getAvatar() != null) hashData.put("avatar", profile.getAvatar());
            if (profile.getAddress() != null) hashData.put("address", profile.getAddress());
            if (profile.getBio() != null) hashData.put("bio", profile.getBio());
            if (profile.getBirthday() != null) hashData.put("birthday", profile.getBirthday().toString());
            if (profile.getIdCard() != null) hashData.put("idCard", profile.getIdCard());
            if (profile.getEmail() != null) hashData.put("email", profile.getEmail());
            if (profile.getPhone() != null) hashData.put("phone", profile.getPhone());

            for (Map.Entry<String, String> entry : hashData.entrySet()) {
                redisUtil.hashSet(key, entry.getKey(), entry.getValue());
            }
            redisUtil.expire(key, PROFILE_CACHE_TTL);
        } catch (Exception e) {
            log.warn("缓存用户画像失败: userId={}", profile.getUserId(), e);
        }
    }

    private UserProfile getCachedProfile(Long userId) {
        String key = PROFILE_CACHE_KEY + userId;
        Map<Object, Object> hashData = redisUtil.hashGetAll(key);
        if (hashData == null || hashData.isEmpty()) {
            return null;
        }
        UserProfile profile = new UserProfile();
        profile.setId(parseLong(getHashValue(hashData, "id")));
        profile.setUserId(userId);
        profile.setNickname(getHashValue(hashData, "nickname"));
        profile.setRealName(getHashValue(hashData, "realName"));
        profile.setGender(parseInt(getHashValue(hashData, "gender")));
        profile.setAvatar(getHashValue(hashData, "avatar"));
        profile.setAddress(getHashValue(hashData, "address"));
        profile.setBio(getHashValue(hashData, "bio"));
        profile.setIdCard(getHashValue(hashData, "idCard"));
        profile.setEmail(getHashValue(hashData, "email"));
        profile.setPhone(getHashValue(hashData, "phone"));
        String birthday = getHashValue(hashData, "birthday");
        if (birthday != null) {
            try {
                profile.setBirthday(java.time.LocalDate.parse(birthday));
            } catch (Exception e) {
                // ignore
            }
        }
        return profile;
    }

    private void invalidateProfileCache(Long userId) {
        redisUtil.delete(PROFILE_CACHE_KEY + userId);
    }

    private String getHashValue(Map<Object, Object> hash, String key) {
        Object value = hash.get(key);
        return value != null ? value.toString() : null;
    }

    private Integer parseInt(String value) {
        if (value == null) return null;
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            return null;
        }
    }

    private Long parseLong(String value) {
        if (value == null) return null;
        try {
            return Long.parseLong(value);
        } catch (NumberFormatException e) {
            return null;
        }
    }
}
