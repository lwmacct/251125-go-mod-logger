package logger

import (
	"fmt"
	"log/slog"
)

// FormatBytes 格式化字节数为人类可读格式
//
// 用于日志中输出文件大小、传输速率等信息
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// LogError 记录错误日志并返回错误
//
// 这是一个便捷函数，用于在需要同时记录日志和返回错误的场景：
//   return logger.LogError(ctx, "操作失败", err, "user_id", userID)
func LogError(ctx interface{}, msg string, err error, attrs ...any) error {
	var logger *slog.Logger

	// 尝试从 context 获取 logger
	if c, ok := ctx.(interface{ Value(key any) any }); ok {
		if l, ok := c.Value(loggerKey).(*slog.Logger); ok {
			logger = l
		}
	}

	// 如果没有从 context 获取到，使用默认 logger
	if logger == nil {
		logger = slog.Default()
	}

	// 合并错误到属性中
	allAttrs := append([]any{"error", err}, attrs...)
	logger.Error(msg, allAttrs...)

	return err
}

// LogAndWrap 记录错误日志并包装错误信息
//
// 用于在错误传播链中添加上下文信息
func LogAndWrap(msg string, err error, attrs ...any) error {
	allAttrs := append([]any{"error", err}, attrs...)
	slog.Error(msg, allAttrs...)
	return fmt.Errorf("%s: %w", msg, err)
}

// Debug 结构化调试日志的快捷方式
func Debug(msg string, attrs ...any) {
	slog.Debug(msg, attrs...)
}

// Info 结构化信息日志的快捷方式
func Info(msg string, attrs ...any) {
	slog.Info(msg, attrs...)
}

// Warn 结构化警告日志的快捷方式
func Warn(msg string, attrs ...any) {
	slog.Warn(msg, attrs...)
}

// Error 结构化错误日志的快捷方式
func Error(msg string, attrs ...any) {
	slog.Error(msg, attrs...)
}
