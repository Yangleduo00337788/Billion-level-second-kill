-- ============================================================
-- 库存回补 Lua 脚本（订单取消/超时未支付时回退库存）
-- KEYS[1]: seckill:stock:{seckillId} — Redis库存Key
-- KEYS[2]: seckill:order:user:{userId}:{seckillId} — 用户秒杀标记Key
--
-- 返回值：
--   1: 回补成功
--   0: 无需回补（Key不存在）
-- ============================================================

local stockKey = KEYS[1]
local userOrderKey = KEYS[2]

local stock = redis.call('GET', stockKey)
if stock == nil then
    return 0
end

redis.call('INCR', stockKey)
redis.call('DEL', userOrderKey)

return 1