-- 初始化数据库脚本
-- 如果数据库不存在则创建（MySQL 8.0 会自动创建，这里作为备用）

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS learn_hub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE learn_hub;

-- 注意：表结构会由 GORM 自动迁移创建，这里不需要手动创建表

