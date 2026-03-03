package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	LogLevel string
	JWT      JWTConfig
	Database DatabaseConfig
	Redis    RedisConfig
	OSS      OSSConfig
	WeChat   WeChatConfig
}

type JWTConfig struct {
	Secret     string
	ExpireTime int // hours
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type OSSConfig struct {
	Region    string
	SecretID  string
	SecretKey string
	Bucket    string
}

type WeChatConfig struct {
	AppID      string
	AppSecret  string
	TemplateID string // 订阅消息模板ID
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: 24 * 7, // 7 days
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "elderly_care"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		OSS: OSSConfig{
			Region:    getEnv("OSS_REGION", ""),
			SecretID:  getEnv("OSS_SECRET_ID", ""),
			SecretKey: getEnv("OSS_SECRET_KEY", ""),
			Bucket:    getEnv("OSS_BUCKET", ""),
		},
		WeChat: WeChatConfig{
			AppID:      getEnv("WECHAT_APP_ID", ""),
			AppSecret:  getEnv("WECHAT_APP_SECRET", ""),
			TemplateID: getEnv("WECHAT_TEMPLATE_ID", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
