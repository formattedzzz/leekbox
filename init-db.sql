USE mysql; ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'leekbox'; flush privileges;
-- 开放连接授权

CREATE DATABASE IF NOT EXISTS leekbox; USE leekbox;
-- 建库