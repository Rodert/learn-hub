-- 初始化用户数据脚本
-- 使用方法: docker compose exec -T mysql mysql -u root -proot123456 learn_hub < scripts/init-user.sql

-- 检查并创建角色
INSERT IGNORE INTO sys_role (id, code, name, description, status, created_at, updated_at)
VALUES 
  (1, 'admin', '管理员', '系统管理员', 1, NOW(), NOW()),
  (2, 'user', '普通用户', '普通用户', 1, NOW(), NOW());

-- 检查并创建管理员用户（密码: admin123）
-- 密码哈希: $2a$10$lZiuaxzQSb.5cuKSbV/F1.5dawptkLSm3p42zwEfY4wWuJDlm.qj2
INSERT IGNORE INTO sys_user (id, username, password, name, email, userid, access, status, created_at, updated_at)
VALUES (1, 'admin', '$2a$10$lZiuaxzQSb.5cuKSbV/F1.5dawptkLSm3p42zwEfY4wWuJDlm.qj2', '管理员', 'admin@example.com', '00000001', 'admin', 1, NOW(), NOW());

-- 关联用户和角色
INSERT IGNORE INTO sys_user_role (user_id, role_id, created_at)
VALUES (1, 1, NOW());

-- 创建默认菜单
INSERT IGNORE INTO sys_menu (id, parent_id, name, path, component, icon, sort_order, access, status, created_at, updated_at)
VALUES 
  (1, 0, 'welcome', '/welcome', './Welcome', 'smile', 1, NULL, 1, NOW(), NOW()),
  (2, 0, 'admin', '/admin', '', 'crown', 2, 'canAdmin', 1, NOW(), NOW()),
  (3, 0, 'list.table-list', '/list', './table-list', 'table', 3, NULL, 1, NOW(), NOW()),
  (4, 2, '', '/admin', '/admin/sub-page', '', 1, NULL, 1, NOW(), NOW()),
  (5, 2, 'sub-page', '/admin/sub-page', './Admin', '', 2, NULL, 1, NOW(), NOW());

-- 为管理员角色分配菜单
INSERT IGNORE INTO sys_role_menu (role_id, menu_id, created_at)
VALUES 
  (1, 1, NOW()),
  (1, 2, NOW()),
  (1, 3, NOW()),
  (1, 4, NOW()),
  (1, 5, NOW());

