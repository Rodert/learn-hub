package main

import (
	"log"

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

	// 自动迁移所有模型
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Permission{},
		&model.RolePermission{},
		&model.Menu{},
		&model.RoleMenu{},
		&model.Material{},
		&model.Question{},
		&model.Exam{},
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
func initDefaultData(db interface{}) error {
	// TODO: 实现默认数据初始化
	// 如：创建默认角色、权限等
	return nil
}
