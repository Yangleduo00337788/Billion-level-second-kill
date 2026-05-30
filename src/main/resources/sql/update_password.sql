-- 更新admin用户密码为 admin123
-- 使用Hutool BCrypt生成的哈希
UPDATE sys_user SET password = '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWvnQlQ5aFvJH3Q5aFvJH3Q5aFvJH3' WHERE username = 'admin';
