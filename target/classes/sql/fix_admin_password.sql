-- 更新 admin 用户密码为 admin123
-- BCrypt 哈希值由 Hutool BCrypt.hashpw("admin123", BCrypt.gensalt()) 生成
UPDATE sys_user SET password = '$2a$10$EqKcp1WFKVQISheBxmXJfeQYHJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQJQ' WHERE username = 'admin';

-- 或者直接用你注册的账号密码覆盖
-- 先查看你注册的账号的密码哈希
-- SELECT username, password FROM sys_user WHERE username = 'test';
