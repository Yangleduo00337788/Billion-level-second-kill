-- ============================================================
-- 企业级IM即时通讯系统数据库脚本
-- 数据库：MySQL 8.0+
-- 字符集：utf8mb4
-- ⚠ 注意：生产环境必须分库分表，此处仅展示逻辑表结构
-- ============================================================

CREATE DATABASE IF NOT EXISTS `im_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `im_system`;

-- ============================================================
-- 1. 会话表 (im_session)
-- ============================================================
DROP TABLE IF EXISTS `im_session`;
CREATE TABLE `im_session` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '会话ID',
    `session_type` TINYINT NOT NULL COMMENT '会话类型：1-单聊，2-群聊',
    `group_id` BIGINT DEFAULT NULL COMMENT '群ID（群聊时）',
    `last_msg_id` BIGINT NOT NULL DEFAULT 0 COMMENT '最后一条消息ID',
    `last_msg_content` VARCHAR(500) DEFAULT NULL COMMENT '最后一条消息内容摘要',
    `last_msg_time` DATETIME DEFAULT NULL COMMENT '最后消息时间',
    `last_msg_sender_id` BIGINT DEFAULT NULL COMMENT '最后消息发送者',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-关闭，1-活跃',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_last_msg_time` (`last_msg_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- ============================================================
-- 2. 会话-用户关联表 (im_session_user)
-- ============================================================
DROP TABLE IF EXISTS `im_session_user`;
CREATE TABLE `im_session_user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `last_read_msg_seq` BIGINT NOT NULL DEFAULT 0 COMMENT '最后已读的消息Seq（用于已读/未读判断）',
    `is_muted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否免打扰：0-否，1-是',
    `is_pinned` TINYINT NOT NULL DEFAULT 0 COMMENT '是否置顶：0-否，1-是',
    `joined_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_session_user` (`session_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_session_id` (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话-用户关联表';

-- ============================================================
-- 3. 消息表 (im_message) — 核心表，生产环境按月分表 (im_message_202501)
-- ============================================================
DROP TABLE IF EXISTS `im_message`;
CREATE TABLE `im_message` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '消息ID（Snowflake分布式ID）',
    `client_msg_id` VARCHAR(64) NOT NULL COMMENT '客户端消息ID（幂等去重key）',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `sender_id` BIGINT NOT NULL COMMENT '发送者用户ID',
    `msg_type` TINYINT NOT NULL COMMENT '消息类型：1-文本，2-图片，3-语音，4-视频，5-文件，6-系统通知',
    `content` TEXT NOT NULL COMMENT '消息内容（文本内容或多媒体URL JSON）',
    `media_url` VARCHAR(500) DEFAULT NULL COMMENT '多媒体文件URL',
    `media_thumb_url` VARCHAR(500) DEFAULT NULL COMMENT '缩略图URL',
    `media_duration` INT DEFAULT NULL COMMENT '语音/视频时长（秒）',
    `media_size` BIGINT DEFAULT NULL COMMENT '文件大小（字节）',
    `seq` BIGINT NOT NULL COMMENT '会话内消息序号（严格递增，用于顺序性保证）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-已撤回',
    `recall_time` DATETIME DEFAULT NULL COMMENT '撤回时间',
    `create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间（毫秒级）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_client_msg_id` (`client_msg_id`),
    UNIQUE KEY `uk_session_seq` (`session_id`, `seq`),
    KEY `idx_session_create_time` (`session_id`, `create_time`),
    KEY `idx_sender_id` (`sender_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表（生产环境按月分表）';

-- ============================================================
-- 4. 群组表 (im_group)
-- ============================================================
DROP TABLE IF EXISTS `im_group`;
CREATE TABLE `im_group` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '群ID',
    `group_name` VARCHAR(100) NOT NULL COMMENT '群名称',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '群头像URL',
    `owner_id` BIGINT NOT NULL COMMENT '群主用户ID',
    `member_count` INT NOT NULL DEFAULT 0 COMMENT '成员数量',
    `max_members` INT NOT NULL DEFAULT 200 COMMENT '最大成员数',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '群描述',
    `notice` VARCHAR(1000) DEFAULT NULL COMMENT '群公告',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-解散，1-正常',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群组表';

-- ============================================================
-- 5. 群成员表 (im_group_member)
-- ============================================================
DROP TABLE IF EXISTS `im_group_member`;
CREATE TABLE `im_group_member` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `group_id` BIGINT NOT NULL COMMENT '群ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role` TINYINT NOT NULL DEFAULT 2 COMMENT '角色：1-群主，2-普通成员，3-管理员',
    `nickname_in_group` VARCHAR(50) DEFAULT NULL COMMENT '群内昵称',
    `is_muted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否禁言：0-否，1-是',
    `muted_until` DATETIME DEFAULT NULL COMMENT '禁言截止时间',
    `joined_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    `last_read_seq` BIGINT NOT NULL DEFAULT 0 COMMENT '最后已读消息Seq',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_user` (`group_id`, `user_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群成员表';

-- ============================================================
-- 6. 消息已读回执表 (im_read_receipt) — 仅群聊，单聊读状态存Redis
-- ============================================================
DROP TABLE IF EXISTS `im_read_receipt`;
CREATE TABLE `im_read_receipt` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `msg_id` BIGINT NOT NULL COMMENT '消息ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `user_id` BIGINT NOT NULL COMMENT '阅读者用户ID',
    `read_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '阅读时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_msg_user` (`msg_id`, `user_id`),
    KEY `idx_session_user` (`session_id`, `user_id`),
    KEY `idx_msg_id` (`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息已读回执表';

-- ============================================================
-- 7. 离线消息推送轨迹表 (im_offline_push_log)
-- ============================================================
DROP TABLE IF EXISTS `im_offline_push_log`;
CREATE TABLE `im_offline_push_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `msg_id` BIGINT NOT NULL COMMENT '消息ID',
    `target_user_id` BIGINT NOT NULL COMMENT '目标用户ID',
    `push_type` TINYINT NOT NULL COMMENT '推送方式：1-AppPush，2-SMS，3-邮件',
    `push_status` TINYINT NOT NULL DEFAULT 0 COMMENT '推送状态：0-待推送，1-已推送，2-推送失败',
    `push_time` DATETIME DEFAULT NULL COMMENT '推送时间',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_msg_target_user` (`msg_id`, `target_user_id`),
    KEY `idx_target_user_status` (`target_user_id`, `push_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='离线消息推送轨迹表';