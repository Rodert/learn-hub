package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"learn-hub/config"
	"learn-hub/internal/model"
	"learn-hub/pkg/database"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 执行迁移
	log.Println("Starting database migration...")

	// 自动迁移所有模型（注意顺序：先创建被引用的表）
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Permission{},
		&model.RolePermission{},
		&model.Menu{},
		&model.RoleMenu{},
		&model.Material{},
		&model.Exam{},
		&model.Question{},
		&model.ExamRecord{},
		&model.CourseRecord{},
		&model.Topic{},
		&model.TopicMaterial{},
		&model.TopicExam{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		log.Fatalf("Failed to initialize default data: %v", err)
	}

	log.Println("Default data initialized successfully")
}

// 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	// 1. 创建默认角色
	roles := []model.Role{
		{Name: "user", Description: "普通用户"},
		{Name: "admin", Description: "管理员"},
		{Name: "system_admin", Description: "系统管理员"},
	}

	for _, role := range roles {
		var count int64
		db.Model(&model.Role{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created role: %s", role.Name)
		}
	}

	// 2. 创建默认权限
	permissions := []model.Permission{
		// 资料权限
		{Name: "materials:view", Description: "浏览学习资料", Resource: "materials", Action: "read"},
		{Name: "materials:create", Description: "上传学习资料", Resource: "materials", Action: "create"},
		{Name: "materials:update", Description: "编辑学习资料", Resource: "materials", Action: "update"},
		{Name: "materials:delete", Description: "删除学习资料", Resource: "materials", Action: "delete"},
		// 考试权限
		{Name: "exams:view", Description: "查看试卷", Resource: "exams", Action: "read"},
		{Name: "exams:submit", Description: "提交答卷", Resource: "exams", Action: "create"},
		{Name: "exams:manage", Description: "管理试卷", Resource: "exams", Action: "update"},
		// 题库权限
		{Name: "questions:manage", Description: "管理题库", Resource: "questions", Action: "update"},
		// 用户权限
		{Name: "users:view", Description: "查看用户数据", Resource: "users", Action: "read"},
		{Name: "users:manage", Description: "管理用户", Resource: "users", Action: "update"},
		// 角色权限
		{Name: "roles:manage", Description: "管理角色", Resource: "roles", Action: "update"},
	}

	for _, perm := range permissions {
		var count int64
		db.Model(&model.Permission{}).Where("name = ?", perm.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&perm).Error; err != nil {
				return err
			}
			log.Printf("Created permission: %s", perm.Name)
		}
	}

	// 3. 分配权限给角色
	rolePermissions := map[string][]string{
		"user": {
			"materials:view",
			"exams:view",
			"exams:submit",
		},
		"admin": {
			"materials:view",
			"materials:create",
			"materials:update",
			"materials:delete",
			"exams:view",
			"exams:manage",
			"questions:manage",
			"users:view",
		},
		"system_admin": {
			"materials:view",
			"materials:create",
			"materials:update",
			"materials:delete",
			"exams:view",
			"exams:submit",
			"exams:manage",
			"questions:manage",
			"users:view",
			"users:manage",
			"roles:manage",
		},
	}

	for roleName, permNames := range rolePermissions {
		var role model.Role
		if err := db.First(&role, "name = ?", roleName).Error; err != nil {
			continue
		}

		for _, permName := range permNames {
			var perm model.Permission
			if err := db.First(&perm, "name = ?", permName).Error; err != nil {
				continue
			}

			var count int64
			db.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, perm.ID).Count(&count)
			if count == 0 {
				if err := db.Create(&model.RolePermission{
					RoleID:       role.ID,
					PermissionID: perm.ID,
				}).Error; err != nil {
					return err
				}
			}
		}
		log.Printf("Assigned permissions to role: %s", roleName)
	}

	// 4. 创建默认菜单
	menus := []model.Menu{
		{Name: "仪表盘", Path: "/dashboard", Icon: "dashboard", Component: "Dashboard", OrderNum: 1, Visible: 1, Type: "menu", Permission: ""},
		{Name: "学习资料", Path: "/materials", Icon: "file-text", Component: "Materials", OrderNum: 2, Visible: 1, Type: "menu", Permission: "materials:view"},
		{Name: "题库管理", Path: "/questions", Icon: "file-text", Component: "Questions", OrderNum: 3, Visible: 1, Type: "menu", Permission: "questions:manage"},
		{Name: "考试管理", Path: "/exams", Icon: "file-text", Component: "Exams", OrderNum: 4, Visible: 1, Type: "menu", Permission: "exams:manage"},
		{Name: "用户管理", Path: "/users", Icon: "users", Component: "Users", OrderNum: 5, Visible: 1, Type: "menu", Permission: "users:view"},
		{Name: "角色权限", Path: "/roles", Icon: "lock", Component: "Roles", OrderNum: 6, Visible: 1, Type: "menu", Permission: "roles:manage"},
	}

	for _, menu := range menus {
		var count int64
		db.Model(&model.Menu{}).Where("path = ?", menu.Path).Count(&count)
		if count == 0 {
			if err := db.Create(&menu).Error; err != nil {
				return err
			}
			log.Printf("Created menu: %s", menu.Name)
		}
	}

	// 5. 创建默认管理员账户
	var adminCount int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "系统管理员",
			Status:   "active",
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}

		// 分配 system_admin 角色
		var systemAdminRole model.Role
		if err := db.First(&systemAdminRole, "name = ?", "system_admin").Error; err == nil {
			db.Create(&model.UserRole{
				UserID: admin.ID,
				RoleID: systemAdminRole.ID,
			})
		}

		log.Printf("Created default admin user: admin (password: admin123)")
	}

	return nil
}
