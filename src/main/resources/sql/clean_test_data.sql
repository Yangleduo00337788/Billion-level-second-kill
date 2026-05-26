-- 清理 Redis 中的秒杀标记（在 Redis CLI 中执行）
-- KEYS seckill:order:user:*
-- DEL seckill:order:user:1:1

-- 或者直接在 MySQL 中删除订单（如果需要重新测试）
-- DELETE FROM sk_order WHERE user_id = 1;
