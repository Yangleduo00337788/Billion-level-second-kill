-- ============================================================
-- 完整数据库初始化脚本
-- 数据库：MySQL 8.0+
-- 所有表创建在 user_system 数据库中
-- ============================================================

CREATE DATABASE IF NOT EXISTS `user_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `user_system`;

-- ============================================================
-- 用户系统表
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
-- 秒杀系统表
-- ============================================================

DROP TABLE IF EXISTS `sk_product`;
CREATE TABLE `sk_product` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `product_name` VARCHAR(200) NOT NULL COMMENT '商品名称',
    `title` VARCHAR(500) DEFAULT NULL COMMENT '商品标题',
    `description` TEXT COMMENT '商品描述',
    `main_image` VARCHAR(500) DEFAULT NULL COMMENT '商品主图URL',
    `original_price` DECIMAL(12,2) NOT NULL COMMENT '原价',
    `price` DECIMAL(12,2) NOT NULL COMMENT '现价',
    `total_stock` INT NOT NULL DEFAULT 0 COMMENT '总库存',
    `sold_count` INT NOT NULL DEFAULT 0 COMMENT '已售数量',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-下架，1-上架',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

DROP TABLE IF EXISTS `sk_seckill_product`;
CREATE TABLE `sk_seckill_product` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '秒杀ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `seckill_price` DECIMAL(12,2) NOT NULL COMMENT '秒杀价格',
    `seckill_stock` INT NOT NULL COMMENT '秒杀总库存',
    `redis_stock_key` VARCHAR(100) NOT NULL COMMENT 'Redis库存Key',
    `person_limit` INT NOT NULL DEFAULT 1 COMMENT '单用户限购数量',
    `start_time` DATETIME NOT NULL COMMENT '秒杀开始时间',
    `end_time` DATETIME NOT NULL COMMENT '秒杀结束时间',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0-未开始，1-进行中，2-已结束',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='秒杀商品表';

DROP TABLE IF EXISTS `sk_order`;
CREATE TABLE `sk_order` (
    `id` BIGINT NOT NULL COMMENT '订单ID（Snowflake）',
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单号',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `seckill_id` BIGINT DEFAULT NULL COMMENT '秒杀ID',
    `product_name` VARCHAR(200) NOT NULL COMMENT '商品名称',
    `product_price` DECIMAL(12,2) NOT NULL COMMENT '商品价格',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '购买数量',
    `total_amount` DECIMAL(12,2) NOT NULL COMMENT '订单金额',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0-待支付，1-已支付，2-已取消',
    `pay_time` DATETIME DEFAULT NULL COMMENT '支付时间',
    `create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    UNIQUE KEY `uk_user_seckill` (`user_id`, `seckill_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

DROP TABLE IF EXISTS `sk_inventory_log`;
CREATE TABLE `sk_inventory_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `seckill_id` BIGINT DEFAULT NULL COMMENT '秒杀ID',
    `change_type` TINYINT NOT NULL COMMENT '变更类型：1-秒杀扣减，2-退款回补',
    `change_amount` INT NOT NULL COMMENT '变更数量',
    `before_stock` INT NOT NULL COMMENT '变更前库存',
    `after_stock` INT NOT NULL COMMENT '变更后库存',
    `order_id` BIGINT DEFAULT NULL COMMENT '关联订单ID',
    `user_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存流水表';

-- ============================================================
-- IM系统表
-- ============================================================

DROP TABLE IF EXISTS `im_session`;
CREATE TABLE `im_session` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '会话ID',
    `session_type` TINYINT NOT NULL COMMENT '会话类型：1-单聊，2-群聊',
    `group_id` BIGINT DEFAULT NULL COMMENT '群ID',
    `last_msg_id` BIGINT NOT NULL DEFAULT 0 COMMENT '最后一条消息ID',
    `last_msg_content` VARCHAR(500) DEFAULT NULL COMMENT '最后一条消息内容',
    `last_msg_time` DATETIME DEFAULT NULL COMMENT '最后消息时间',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-关闭，1-活跃',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

DROP TABLE IF EXISTS `im_session_user`;
CREATE TABLE `im_session_user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `last_read_msg_seq` BIGINT NOT NULL DEFAULT 0 COMMENT '最后已读消息Seq',
    `joined_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_session_user` (`session_id`, `user_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话用户表';

DROP TABLE IF EXISTS `im_message`;
CREATE TABLE `im_message` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `client_msg_id` VARCHAR(64) NOT NULL COMMENT '客户端消息ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `sender_id` BIGINT NOT NULL COMMENT '发送者用户ID',
    `msg_type` TINYINT NOT NULL COMMENT '消息类型：1-文本，2-图片，3-语音',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `seq` BIGINT NOT NULL COMMENT '会话内消息序号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-已撤回',
    `recall_time` DATETIME DEFAULT NULL COMMENT '撤回时间',
    `create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_client_msg_id` (`client_msg_id`),
    KEY `idx_session_create_time` (`session_id`, `create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

DROP TABLE IF EXISTS `im_group`;
CREATE TABLE `im_group` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '群ID',
    `group_name` VARCHAR(100) NOT NULL COMMENT '群名称',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '群头像',
    `owner_id` BIGINT NOT NULL COMMENT '群主用户ID',
    `member_count` INT NOT NULL DEFAULT 0 COMMENT '成员数量',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-解散，1-正常',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群组表';

DROP TABLE IF EXISTS `im_group_member`;
CREATE TABLE `im_group_member` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `group_id` BIGINT NOT NULL COMMENT '群ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role` TINYINT NOT NULL DEFAULT 2 COMMENT '角色：1-群主，2-普通成员',
    `joined_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_user` (`group_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群成员表';

DROP TABLE IF EXISTS `im_read_receipt`;
CREATE TABLE `im_read_receipt` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `msg_id` BIGINT NOT NULL COMMENT '消息ID',
    `session_id` BIGINT NOT NULL COMMENT '会话ID',
    `user_id` BIGINT NOT NULL COMMENT '阅读者用户ID',
    `read_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_msg_user` (`msg_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='已读回执表';

-- ============================================================
-- 初始数据
-- ============================================================

INSERT INTO `sys_role` (`role_code`, `role_name`, `description`, `sort`) VALUES
('ROLE_ADMIN', '超级管理员', '系统超级管理员', 1),
('ROLE_USER', '普通用户', '普通注册用户', 100);

INSERT INTO `sys_user` (`username`, `password`, `status`, `login_count`) VALUES
('admin', '$2a$10$EqKcp1WFKVQISheBxmXJfeQYHJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQ', 1, 0);

INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (1, 1);

INSERT INTO `sys_user_profile` (`user_id`, `nickname`, `real_name`) VALUES (1, '系统管理员', 'Admin');

INSERT INTO `sys_permission` (`parent_id`, `permission_code`, `permission_name`, `permission_type`, `sort`) VALUES
(0, 'system', '系统管理', 1, 1),
(0, 'user', '用户管理', 1, 2);

INSERT INTO `sys_permission` (`parent_id`, `permission_code`, `permission_name`, `permission_type`, `method`, `sort`) VALUES
(1, 'system:role:list', '角色列表', 2, 'GET', 1),
(1, 'system:role:add', '新增角色', 2, 'POST', 2),
(2, 'user:list', '用户列表', 2, 'GET', 1),
(2, 'user:add', '新增用户', 2, 'POST', 2);

INSERT INTO `sys_role_permission` (`role_id`, `permission_id`)
SELECT 1, id FROM `sys_permission`;

-- 秒杀商品测试数据
INSERT INTO `sk_product` (`product_name`, `original_price`, `price`, `total_stock`, `status`) VALUES
('iPhone 15 Pro Max', 9999.00, 8999.00, 100, 1),
('MacBook Pro 14', 14999.00, 12999.00, 50, 1),
('AirPods Pro 2', 1899.00, 1499.00, 200, 1);

INSERT INTO `sk_seckill_product` (`product_id`, `seckill_price`, `seckill_stock`, `redis_stock_key`, `person_limit`, `start_time`, `end_time`, `status`) VALUES
(1, 7999.00, 10, 'seckill:stock:1', 1, DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_ADD(NOW(), INTERVAL 1 HOUR), 1),
(2, 10999.00, 5, 'seckill:stock:2', 1, DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_ADD(NOW(), INTERVAL 2 HOUR), 1),
(3, 999.00, 50, 'seckill:stock:3', 2, DATE_SUB(NOW(), INTERVAL 30 MINUTE), DATE_ADD(NOW(), INTERVAL 3 HOUR), 1);
