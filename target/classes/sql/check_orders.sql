-- 查看所有订单
SELECT id, order_no, user_id, status, create_time FROM sk_order ORDER BY create_time DESC LIMIT 10;
