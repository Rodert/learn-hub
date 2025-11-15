package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"learn-hub-backend/config"
	"learn-hub-backend/models"
	"learn-hub-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 构建 DSN，添加必要的连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	var sqlDB *sql.DB

	// 尝试连接，最多重试 10 次
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// 获取底层的 sql.DB 以配置连接池
			sqlDB, err = DB.DB()
			if err == nil {
				// 设置连接池参数
				sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中连接的最大数量
				sqlDB.SetMaxOpenConns(100)                 // 设置打开数据库连接的最大数量
				sqlDB.SetConnMaxLifetime(time.Hour)        // 设置了连接可复用的最大时间
				sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 设置空闲连接的最大存活时间

				// 测试连接
				if err = sqlDB.Ping(); err == nil {
					log.Println("数据库连接成功")
					break
				}
			}
		}

		if i < maxRetries-1 {
			log.Printf("数据库连接失败，正在重试 (%d/%d)...", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("数据库连接失败（已重试 %d 次）: %v", maxRetries, err)
	}

	// 自动迁移
	AutoMigrate()
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Rule{},
	)

	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	log.Println("数据库迁移完成")

	// 初始化基础数据
	InitData()
}

// InitData 初始化基础数据
func InitData() {
	// 检查是否已有数据
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	// 创建默认角色
	adminRole := models.Role{
		Code:        "admin",
		Name:        "管理员",
		Description: "系统管理员",
		Status:      1,
	}
	userRole := models.Role{
		Code:        "user",
		Name:        "普通用户",
		Description: "普通用户",
		Status:      1,
	}
	DB.Create(&adminRole)
	DB.Create(&userRole)

	// 创建默认管理员用户（密码: admin123）
	// 使用动态生成的密码哈希，确保正确性
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Println("生成密码哈希失败，使用默认值:", err)
		hashedPassword = "$2a$10$lZiuaxzQSb.5cuKSbV/F1.5dawptkLSm3p42zwEfY4wWuJDlm.qj2" // admin123
	}

	adminUser := models.User{
		Username: "admin",
		Password: hashedPassword,
		Name:     "管理员",
		Email:    "admin@example.com",
		UserID:   "00000001",
		Access:   "admin",
		Status:   1,
		Roles:    []models.Role{adminRole},
	}
	DB.Create(&adminUser)

	// 创建默认菜单
	welcomeMenu := models.Menu{
		ParentID: 0, Name: "welcome", Path: "/welcome", Component: "./Welcome", Icon: "smile", SortOrder: 1, Status: 1,
	}
	adminMenu := models.Menu{
		ParentID: 0, Name: "admin", Path: "/admin", Component: "", Icon: "crown", Access: "canAdmin", SortOrder: 2, Status: 1,
	}
	listMenu := models.Menu{
		ParentID: 0, Name: "list.table-list", Path: "/list", Component: "./table-list", Icon: "table", SortOrder: 3, Status: 1,
	}

	DB.Create(&welcomeMenu)
	DB.Create(&adminMenu)
	DB.Create(&listMenu)

	// 创建 admin 菜单的子菜单
	adminSubMenus := []models.Menu{
		{ParentID: adminMenu.ID, Path: "/admin", Redirect: "/admin/sub-page", SortOrder: 1, Status: 1},
		{ParentID: adminMenu.ID, Name: "sub-page", Path: "/admin/sub-page", Component: "./Admin", SortOrder: 2, Status: 1},
	}
	for _, subMenu := range adminSubMenus {
		DB.Create(&subMenu)
	}

	// 为管理员角色分配所有菜单
	allMenus := []models.Menu{welcomeMenu, adminMenu, listMenu}
	allMenus = append(allMenus, adminSubMenus...)
	DB.Model(&adminRole).Association("Menus").Append(&allMenus)

	log.Println("基础数据初始化完成")
}
