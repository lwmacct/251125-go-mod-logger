# Logger 包

基于 Go 1.21+ `log/slog` 的统一结构化日志系统。

## 特性

- 支持多种输出格式：JSON、Text、Colored（彩色终端）
- 灵活的日志级别控制（DEBUG、INFO、WARN、ERROR）
- 多种时间格式配置
- 环境自动检测（开发/生产环境智能默认值）
- Context 集成，支持请求链路追踪
- 彩色 Handler 支持 JSON/map/struct 自动平铺
- WithGroup 分组支持（嵌套分组自动添加前缀）
- 配置验证（无效格式/级别会返回错误）
- 文件输出支持，带关闭机制

## 快速开始

### 初始化方式

```go
// 方式一：自动检测环境（推荐）
// IS_SANDBOX=1 时使用开发配置，否则使用生产配置
logger.InitAuto()

// 方式二：从环境变量初始化（固定默认值）
logger.InitEnv()

// 方式三：手动配置
logger.InitCfg(&logger.Config{
    Level:  "DEBUG",
    Format: "color",
})
```

### 基础用法

```go
package main

import (
    "log/slog"
    "github.com/lwmacct/251125-go-mod-logger/pkg/logger"
)

func main() {
    // 自动检测环境初始化
    if err := logger.InitAuto(); err != nil {
        panic(err)
    }
    defer logger.Close()

    // 使用全局 logger
    slog.Info("server started", "port", 8080)
    slog.Debug("debug info", "key", "value")
    slog.Warn("warning message")
    slog.Error("error occurred", "error", err)
}
```

## 初始化 API

### InitAuto - 自动检测环境

根据 `IS_SANDBOX` 环境变量自动选择开发或生产配置：

| 配置项         | 开发环境 (IS_SANDBOX=1) | 生产环境       |
|----------------|-------------------------|----------------|
| LOG_LEVEL      | DEBUG                   | INFO           |
| LOG_FORMAT     | color                   | json           |
| LOG_ADD_SOURCE | true                    | false          |
| LOG_TIME_FORMAT| time (15:04:05)         | datetime       |

```go
logger.InitAuto()
```

### InitEnv - 从环境变量初始化

使用固定的默认值，适合需要明确控制的场景：

```go
// 支持的环境变量：
// - LOG_LEVEL: DEBUG, INFO, WARN, ERROR（默认 INFO）
// - LOG_FORMAT: json, text, color（默认 color）
// - LOG_OUTPUT: stdout, stderr, 或文件路径（默认 stdout）
// - LOG_ADD_SOURCE: true, false（默认 true）
// - LOG_TIME_FORMAT: datetime, time, rfc3339, rfc3339ms（默认 datetime）

logger.InitEnv()
```

### InitCfg - 手动配置

完全控制所有配置项：

```go
cfg := &logger.Config{
    Level:      "DEBUG",
    Format:     "json",        // json, text, color
    Output:     "stdout",      // stdout, stderr, /path/to/file.log
    AddSource:  true,          // 添加源码位置
    TimeFormat: "rfc3339ms",   // 时间格式
    Timezone:   "Asia/Shanghai",
}
logger.InitCfg(cfg)
```

## 输出格式

### JSON 格式

```json
{"time":"2024-01-15T10:30:00.123+08:00","level":"INFO","msg":"request received","method":"GET","path":"/api"}
```

### Text 格式

```
time="2024-01-15 10:30:00" level=INFO msg="request received" method=GET path=/api
```

### Colored 格式（终端）

```json
{"time":"10:30:00","level":"INFO","msg":"request received","method":"GET","path":"/api"}
```

带有颜色高亮：
- DEBUG: 蓝色
- INFO: 绿色
- WARN: 黄色
- ERROR: 红色

## 时间格式

| 格式        | 示例                            |
| ----------- | ------------------------------- |
| `time`      | `10:30:00`                      |
| `timems`    | `10:30:00.123`                  |
| `datetime`  | `2024-01-15 10:30:00`           |
| `rfc3339`   | `2024-01-15T10:30:00+08:00`     |
| `rfc3339ms` | `2024-01-15T10:30:00.123+08:00` |
| 自定义      | 使用 Go 时间格式字符串           |

## 高级功能

### 自动平铺（JSON/map/struct）

彩色 Handler 会自动将复杂类型平铺为 `key.subkey` 格式。

```go
// JSON 字符串
slog.Info("request", "body", `{"user":"alice","age":30}`)
// 输出: {"msg":"request","body.user":"alice","body.age":"30"}

// map[string]any
slog.Info("request", "data", map[string]any{"user": "bob", "active": true})
// 输出: {"msg":"request","data.user":"bob","data.active":"true"}

// struct
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
slog.Info("user", "info", User{Name: "charlie", Age: 30})
// 输出: {"msg":"user","info.name":"charlie","info.age":"30"}
```

### Context 集成

```go
// 将 logger 存入 context
ctx := logger.WithLogger(ctx, customLogger)

// 从 context 获取 logger
log := logger.FromContext(ctx)

// 便捷方法：添加 request_id
ctx = logger.WithRequestID(ctx, "req-123")
```

### 创建独立 Logger

```go
// 为特定模块创建独立配置的 logger
moduleLogger, err := logger.New(&logger.Config{
    Level:  "DEBUG",
    Format: "json",
    Output: "/var/log/module.log",
})

// 如需手动关闭文件
moduleLogger, closer, err := logger.NewWithCloser(cfg)
defer closer.Close()
```

### 带属性的 Logger

```go
// 添加固定属性
log := logger.WithAttrs("service", "api", "version", "1.0")
log.Info("started")  // 每条日志都会包含 service 和 version

// 分组
log := logger.WithGroup("request")
log.Info("received", "method", "GET")  // method 在 request 分组下
```

### 辅助函数

```go
// 格式化字节数
logger.FormatBytes(1536 * 1024)  // "1.5 MB"

// 记录错误并返回
return logger.LogError(ctx, "operation failed", err, "user_id", userID)

// 记录并包装错误
return logger.LogAndWrap("fetch failed", err, "url", url)
```

## 资源管理

输出到文件时，应在程序退出时关闭：

```go
func main() {
    logger.InitCfg(&logger.Config{
        Output: "/var/log/app.log",
    })
    defer logger.Close()  // 确保文件正确关闭

    // ...
}
```

## 配置验证

Config 提供 `Validate()` 方法，在 `InitCfg()` 和 `New()` 时自动调用：

```go
cfg := &logger.Config{
    Level:  "TRACE",        // 无效级别
    Format: "yaml",         // 无效格式
}

// 验证会返回详细错误
err := cfg.Validate()
// err: invalid log format: "yaml", valid options: json, text, color

// InitCfg/New 也会自动验证
err := logger.InitCfg(cfg)
// err: invalid log level: "TRACE", valid options: DEBUG, INFO, WARN, ERROR
```

有效的配置选项：
- **Level**: `DEBUG`, `INFO`, `WARN`, `WARNING`, `ERROR`（大小写不敏感）
- **Format**: `json`, `text`, `color`, `colored`

## 最佳实践

1. **应用启动时初始化一次**

   ```go
   func main() {
       logger.InitAuto()  // 或 InitEnv()
       defer logger.Close()
   }
   ```

2. **使用结构化日志**

   ```go
   // 好
   slog.Info("user login", "user_id", userID, "ip", ip)

   // 避免
   slog.Info(fmt.Sprintf("user %s login from %s", userID, ip))
   ```

3. **传递 Context**

   ```go
   func HandleRequest(ctx context.Context) {
       log := logger.FromContext(ctx)
       log.Info("processing request")
   }
   ```

4. **使用适当的日志级别**
   - DEBUG: 开发调试信息
   - INFO: 重要业务事件
   - WARN: 警告，但不影响运行
   - ERROR: 错误，需要关注

## 测试

```bash
go test ./pkg/logger/... -v
```
