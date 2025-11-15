package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	JWTExpire  int // 小时
	Port       string
}

var AppConfig *Config

func LoadConfig() {
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load()

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "learn_hub"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpire:  24, // 24小时
		Port:       getEnv("PORT", "8080"),
	}

	// 调试输出（生产环境可删除）
	log.Printf("数据库配置: host=%s, port=%s, user=%s, dbname=%s",
		AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBUser, AppConfig.DBName)
	log.Println("配置加载完成")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
