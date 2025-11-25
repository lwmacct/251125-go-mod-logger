package logger

import (
	"os"
	"strings"
)

// InitFromEnv 从环境变量初始化日志系统
//
// 支持的环境变量：
//   - LOG_LEVEL: 日志级别 (DEBUG, INFO, WARN, ERROR)
//   - LOG_FORMAT: 输出格式 (json, text, color/colored)
//   - LOG_OUTPUT: 输出目标 (stdout, stderr, 或文件路径)
//   - LOG_ADD_SOURCE: 是否添加源代码位置 (true, false)
//   - LOG_TIME_FORMAT: 时间格式 (rfc3339, rfc3339ms, unix, unixms, unixfloat, datetime)

func InitFromEnv() error {
	cfg := &Config{
		Level:      getEnv("LOG_LEVEL", "INFO"),
		Format:     getEnv("LOG_FORMAT", "color"), // 默认使用彩色输出
		Output:     getEnv("LOG_OUTPUT", "stdout"),
		AddSource:  getEnvBool("LOG_ADD_SOURCE", true), // 默认启用源代码位置
		TimeFormat: getEnv("LOG_TIME_FORMAT", "rfc3339ms"),
	}

	return Init(cfg)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔类型的环境变量
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}
