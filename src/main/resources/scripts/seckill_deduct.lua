-- ============================================================
-- 秒杀扣减库存 Lua 脚本（核心）
-- KEYS[1]: seckill:stock:{seckillId} — Redis库存Key
-- KEYS[2]: seckill:order:user:{userId}:{seckillId} — 用户秒杀标记Key
-- ARGV[1]: userId
-- ARGV[2]: seckillId
-- ARGV[3]: personLimit — 单用户限购数
--
-- 返回值：
--   0: 秒杀成功
--  -1: 库存不足
--  -2: 重复秒杀
--  -3: 用户限购超限
-- ============================================================

local stockKey = KEYS[1]
local userOrderKey = KEYS[2]
local userId = ARGV[1]
local seckillId = ARGV[2]
local personLimit = tonumber(ARGV[3])

-- 1. 检查是否重复下单
local exists = redis.call('EXISTS', userOrderKey)
if exists == 1 then
    return -2
end

-- 2. 检查库存
local stock = tonumber(redis.call('GET', stockKey))
if stock == nil then
    return -1
end
if stock <= 0 then
    return -1
end

-- 3. 扣减库存
local afterStock = redis.call('DECR', stockKey)

-- 4. 库存扣减后<0说明并发超卖，回滚
if afterStock < 0 then
    redis.call('INCR', stockKey)
    return -1
end

-- 5. 标记用户已下单（去重）
redis.call('SETEX', userOrderKey, 3600, userId)

-- 6. 记录秒杀成功用户（用于后续统计）
redis.call('SADD', 'seckill:success:users:' .. seckillId, userId)

return 0