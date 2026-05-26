-- ============================================================
-- 企业级用户系统数据库脚本
-- 数据库：MySQL 8.0+
-- 字符集：utf8mb4
-- 排序规则：utf8mb4_unicode_ci
-- ============================================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `user_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `user_system`;

-- ============================================================
-- 1. 用户表 (sys_user)
-- ============================================================
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（BCrypt加密）',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用，2-锁定',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-未删除，1-已删除',
    `login_type` TINYINT NOT NULL DEFAULT 1 COMMENT '登录类型：1-密码，2-短信，3-OAuth',
    `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
    `login_count` INT NOT NULL DEFAULT 0 COMMENT '登录次数',
    `password_update_time` DATETIME DEFAULT NULL COMMENT '密码最后更新时间',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_phone` (`phone`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================================
-- 2. 角色表 (sys_role)
-- ============================================================
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `role_code` VARCHAR(50) NOT NULL COMMENT '角色编码',
    `role_name` VARCHAR(100) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(255) DEFAULT NULL COMMENT '角色描述',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-未删除，1-已删除',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_code` (`role_code`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================================
-- 3. 权限表 (sys_permission)
-- ============================================================
DROP TABLE IF EXISTS `sys_permission`;
CREATE TABLE `sys_permission` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `parent_id` BIGINT NOT NULL DEFAULT 0 COMMENT '父权限ID',
    `permission_code` VARCHAR(100) NOT NULL COMMENT '权限编码',
    `permission_name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `permission_type` TINYINT NOT NULL COMMENT '权限类型：1-菜单，2-按钮，3-接口',
    `path` VARCHAR(255) DEFAULT NULL COMMENT '路由路径',
    `method` VARCHAR(20) DEFAULT NULL COMMENT 'HTTP方法：GET/POST/PUT/DELETE',
    `icon` VARCHAR(100) DEFAULT NULL COMMENT '图标',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-未删除，1-已删除',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permission_code` (`permission_code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- ============================================================
-- 4. 用户角色关联表 (sys_user_role)
-- ============================================================
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================================
-- 5. 角色权限关联表 (sys_role_permission)
-- ============================================================
DROP TABLE IF EXISTS `sys_role_permission`;
CREATE TABLE `sys_role_permission` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT NOT NULL COMMENT '权限ID',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ============================================================
-- 6. 用户资料表 (sys_user_profile)
-- ============================================================
DROP TABLE IF EXISTS `sys_user_profile`;
CREATE TABLE `sys_user_profile` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `real_name` VARCHAR(50) DEFAULT NULL COMMENT '真实姓名',
    `gender` TINYINT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `address` VARCHAR(255) DEFAULT NULL COMMENT '地址',
    `id_card` VARCHAR(20) DEFAULT NULL COMMENT '身份证号',
    `bio` VARCHAR(500) DEFAULT NULL COMMENT '个人简介',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    KEY `idx_real_name` (`real_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户资料表';

-- ============================================================
-- 7. 黑名单表 (sys_blacklist)
-- ============================================================
DROP TABLE IF EXISTS `sys_blacklist`;
CREATE TABLE `sys_blacklist` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT DEFAULT NULL COMMENT '用户ID（封禁用户时）',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `device_id` VARCHAR(100) DEFAULT NULL COMMENT '设备ID',
    `blacklist_type` TINYINT NOT NULL COMMENT '黑名单类型：1-用户，2-IP，3-设备',
    `reason` VARCHAR(500) NOT NULL COMMENT '封禁原因',
    `expire_time` DATETIME DEFAULT NULL COMMENT '过期时间（NULL表示永久）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-已解封，1-封禁中',
    `operator_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_ip_address` (`ip_address`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_status` (`status`),
    KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='黑名单表';

-- ============================================================
-- 8. 设备表 (sys_device)
-- ============================================================
DROP TABLE IF EXISTS `sys_device`;
CREATE TABLE `sys_device` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '设备ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `device_id` VARCHAR(100) NOT NULL COMMENT '设备唯一标识',
    `device_type` VARCHAR(20) NOT NULL COMMENT '设备类型：WEB/IOS/ANDROID/PC/MAC',
    `device_name` VARCHAR(100) DEFAULT NULL COMMENT '设备名称',
    `os_name` VARCHAR(50) DEFAULT NULL COMMENT '操作系统',
    `browser` VARCHAR(50) DEFAULT NULL COMMENT '浏览器',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `location` VARCHAR(100) DEFAULT NULL COMMENT '登录地点',
    `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `last_active_time` DATETIME DEFAULT NULL COMMENT '最后活跃时间',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-正常',
    `is_trusted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否可信设备：0-否，1-是',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_device` (`user_id`, `device_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备表';

-- ============================================================
-- 9. 登录日志表 (sys_login_log)
-- ============================================================
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `login_type` TINYINT NOT NULL COMMENT '登录类型：1-密码，2-短信，3-OAuth',
    `device_type` VARCHAR(20) DEFAULT NULL COMMENT '设备类型',
    `device_id` VARCHAR(100) DEFAULT NULL COMMENT '设备ID',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `location` VARCHAR(100) DEFAULT NULL COMMENT '登录地点',
    `browser` VARCHAR(50) DEFAULT NULL COMMENT '浏览器',
    `os_name` VARCHAR(50) DEFAULT NULL COMMENT '操作系统',
    `status` TINYINT NOT NULL COMMENT '登录状态：0-失败，1-成功',
    `fail_reason` VARCHAR(255) DEFAULT NULL COMMENT '失败原因',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT 'User-Agent',
    `login_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_username` (`username`),
    KEY `idx_login_time` (`login_time`),
    KEY `idx_ip_address` (`ip_address`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ============================================================
-- 10. OAuth 授权表 (sys_user_oauth)
-- ============================================================
DROP TABLE IF EXISTS `sys_user_oauth`;
CREATE TABLE `sys_user_oauth` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `oauth_type` VARCHAR(20) NOT NULL COMMENT 'OAuth类型：WECHAT/ALIPAY/QQ/GOOGLE/GITHUB',
    `openid` VARCHAR(100) NOT NULL COMMENT '第三方用户唯一标识',
    `unionid` VARCHAR(100) DEFAULT NULL COMMENT '第三方用户统一标识',
    `access_token` VARCHAR(500) DEFAULT NULL COMMENT '访问令牌',
    `refresh_token` VARCHAR(500) DEFAULT NULL COMMENT '刷新令牌',
    `expires_in` INT DEFAULT NULL COMMENT '令牌有效期（秒）',
    `nickname` VARCHAR(100) DEFAULT NULL COMMENT '第三方昵称',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '第三方头像',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_oauth_type_openid` (`oauth_type`, `openid`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='OAuth授权表';

-- ============================================================
-- 初始数据
-- ============================================================

-- 插入默认角色
INSERT INTO `sys_role` (`role_code`, `role_name`, `description`, `sort`) VALUES
('ROLE_ADMIN', '超级管理员', '系统超级管理员，拥有所有权限', 1),
('ROLE_USER', '普通用户', '普通注册用户', 100);

-- 插入默认管理员用户（密码：admin123，BCrypt加密）
INSERT INTO `sys_user` (`username`, `password`, `status`, `login_count`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 1, 0);

-- 关联管理员角色
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (1, 1);

-- 插入管理员资料
INSERT INTO `sys_user_profile` (`user_id`, `nickname`, `real_name`) VALUES (1, '系统管理员', 'Admin');

-- 插入基础权限
INSERT INTO `sys_permission` (`parent_id`, `permission_code`, `permission_name`, `permission_type`, `sort`) VALUES
(0, 'system', '系统管理', 1, 1),
(0, 'user', '用户管理', 1, 2);

INSERT INTO `sys_permission` (`parent_id`, `permission_code`, `permission_name`, `permission_type`, `method`, `sort`) VALUES
(1, 'system:role:list', '角色列表', 2, 'GET', 1),
(1, 'system:role:add', '新增角色', 2, 'POST', 2),
(1, 'system:role:update', '修改角色', 2, 'PUT', 3),
(1, 'system:role:delete', '删除角色', 2, 'DELETE', 4),
(2, 'user:list', '用户列表', 2, 'GET', 1),
(2, 'user:add', '新增用户', 2, 'POST', 2),
(2, 'user:update', '修改用户', 2, 'PUT', 3),
(2, 'user:delete', '删除用户', 2, 'DELETE', 4);

-- 关联管理员权限（角色ID=1拥有所有权限）
INSERT INTO `sys_role_permission` (`role_id`, `permission_id`)
SELECT 1, id FROM `sys_permission`;
