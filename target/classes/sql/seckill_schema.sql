-- ============================================================
-- 高并发秒杀电商系统数据库脚本
-- 数据库：MySQL 8.0+
-- [WARNING] 生产环境必须分库分表，此处仅展示逻辑表结构
-- ============================================================

CREATE DATABASE IF NOT EXISTS `seckill_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `seckill_system`;

-- ============================================================
-- 1. 商品表 (sk_product)
-- ============================================================
DROP TABLE IF EXISTS `sk_product`;
CREATE TABLE `sk_product` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `product_name` VARCHAR(200) NOT NULL COMMENT '商品名称',
    `title` VARCHAR(500) DEFAULT NULL COMMENT '商品标题',
    `description` TEXT COMMENT '商品描述',
    `main_image` VARCHAR(500) DEFAULT NULL COMMENT '商品主图URL',
    `detail_images` JSON DEFAULT NULL COMMENT '详情图片列表',
    `original_price` DECIMAL(12,2) NOT NULL COMMENT '原价',
    `price` DECIMAL(12,2) NOT NULL COMMENT '现价',
    `total_stock` INT NOT NULL DEFAULT 0 COMMENT '总库存',
    `sold_count` INT NOT NULL DEFAULT 0 COMMENT '已售数量',
    `category_id` BIGINT DEFAULT NULL COMMENT '分类ID',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-下架，1-上架',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- ============================================================
-- 2. 秒杀商品表 (sk_seckill_product)
-- ============================================================
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
    KEY `idx_start_time` (`start_time`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='秒杀商品表';

-- ============================================================
-- 3. 订单表 (sk_order)
-- ============================================================
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
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0-待支付，1-已支付，2-已取消，3-已退款，4-已完成',
    `pay_time` DATETIME DEFAULT NULL COMMENT '支付时间',
    `cancel_time` DATETIME DEFAULT NULL COMMENT '取消时间',
    `create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    UNIQUE KEY `uk_user_seckill` (`user_id`, `seckill_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- ============================================================
-- 4. 库存流水表 (sk_inventory_log) — 库存变更审计
-- ============================================================
DROP TABLE IF EXISTS `sk_inventory_log`;
CREATE TABLE `sk_inventory_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `seckill_id` BIGINT DEFAULT NULL COMMENT '秒杀ID',
    `change_type` TINYINT NOT NULL COMMENT '变更类型：1-秒杀扣减，2-退款回补，3-管理员调整',
    `change_amount` INT NOT NULL COMMENT '变更数量（负数表示扣减）',
    `before_stock` INT NOT NULL COMMENT '变更前库存',
    `after_stock` INT NOT NULL COMMENT '变更后库存',
    `order_id` BIGINT DEFAULT NULL COMMENT '关联订单ID',
    `user_id` BIGINT DEFAULT NULL COMMENT '操作人ID',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存流水表';