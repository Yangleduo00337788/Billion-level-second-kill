-- 查看秒杀商品状态
SELECT id, status, seckill_stock, start_time, end_time,
       CASE status
           WHEN 0 THEN '未开始'
           WHEN 1 THEN '进行中'
           WHEN 2 THEN '已结束'
       END as status_text
FROM sk_seckill_product;
