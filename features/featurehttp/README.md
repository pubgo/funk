# Feature HTTP 模块

Feature HTTP 模块提供了一个简单的 Web 界面来查看和管理功能标志（Feature Flags）。

## 功能特性

- 📋 **功能列表展示**: 以美观的 HTML 界面展示所有注册的功能标志
- ✏️ **在线编辑**: 直接在 Web 界面中修改功能标志的值
- 🔍 **搜索功能**: 支持按名称或描述搜索功能标志
- 🔒 **安全控制**: 自动保护敏感字段，不允许修改
- 🎨 **现代化 UI**: 响应式设计，支持移动端访问

## 快速开始

### 基本使用

```go
package main

import (
	"log"

	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/features/featurehttp"
)

func main() {
	// 注册一些功能标志
	_ = features.String("app_status", "ok", "应用程序状态")
	_ = features.Int("replicas", 1, "副本数量")
	_ = features.Bool("debug", false, "启用调试模式")

	// 创建并启动 HTTP 服务器
	server := featurehttp.NewServer(":8181")
	log.Println("访问 http://localhost:8181/features 查看功能标志管理界面")
	
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
```

### 集成到现有 HTTP 服务器

```go
package main

import (
	"net/http"

	"github.com/pubgo/funk/v2/features/featurehttp"
)

func main() {
	mux := http.NewServeMux()
	
	// 创建 featurehttp 服务器
	server := featurehttp.NewServer("")
	
	// 挂载到现有路由
	mux.Handle("/features/", http.StripPrefix("/features", server.Handler()))
	
	// 其他路由...
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	
	http.ListenAndServe(":8080", mux)
}
```

### 自定义 URL 前缀

```go
server := featurehttp.NewServer(":8181")
server.WithPrefix("/admin/features")  // 默认是 /features
server.Start()
```

## API 端点

### GET `/features/`
返回 HTML 管理界面

### GET `/features/api/list`
获取所有功能标志的 JSON 列表

**响应示例:**
```json
[
  {
    "name": "app_status",
    "usage": "应用程序状态",
    "type": "string",
    "value": "ok",
    "valueString": "ok",
    "deprecated": false,
    "tags": {
      "group": "health",
      "mutable": true
    },
    "sensitive": false,
    "mutable": true
  }
]
```

### POST `/features/api/update`
更新功能标志的值

**请求体:**
```json
{
  "name": "app_status",
  "value": "degraded"
}
```

**响应示例:**
```json
{
  "success": true,
  "message": "Feature app_status updated successfully",
  "flag": {
    "name": "app_status",
    "usage": "应用程序状态",
    "type": "string",
    "value": "degraded",
    "valueString": "degraded",
    "deprecated": false,
    "tags": {},
    "sensitive": false,
    "mutable": true
  }
}
```

## 安全特性

### 敏感字段保护

标记为 `sensitive: true` 的功能标志：
- 在界面上显示为 `******`
- 不允许通过 API 修改
- 在 JSON 响应中值被隐藏

```go
_ = features.String("api_key", "sk-xxxx", "API 密钥",
    map[string]any{
        "sensitive": true,
        "group":     "security",
    })
```

### 不可修改字段

标记为 `mutable: false` 的功能标志：
- 在界面上显示"不可修改"标签
- 编辑按钮被禁用
- 不允许通过 API 修改

```go
_ = features.String("version", "1.0.0", "版本号",
    map[string]any{
        "mutable": false,
    })
```

## 界面功能

### 搜索
在搜索框中输入关键词，可以按功能名称或描述进行过滤。

### 编辑功能
1. 点击"编辑"按钮
2. 修改输入框中的值
3. 点击"保存"或按 Enter 键提交
4. 点击"取消"放弃修改

### 类型标识
每个功能标志都会显示其类型（string、int、bool、float、json）。

### 标签显示
如果功能标志包含标签（tags），会在卡片底部显示。

## 注意事项

1. **并发安全**: 功能标志的修改是线程安全的
2. **类型验证**: 更新时会自动验证值的类型是否正确
3. **错误处理**: 如果更新失败，会在界面上显示错误消息
4. **实时更新**: 修改后的值会立即反映在界面上

