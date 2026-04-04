package featurehttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/pubgo/funk/v2/features"
)

// Server 提供 HTTP 服务器来展示和修改 features
type Server struct {
	mux    *http.ServeMux
	addr   string
	prefix string
}

// NewServer 创建一个新的 feature HTTP 服务器
func NewServer(addr string) *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		addr:   addr,
		prefix: "/features",
	}
	s.setupRoutes()
	return s
}

// WithPrefix 设置 URL 前缀
func (s *Server) WithPrefix(prefix string) *Server {
	s.prefix = prefix
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux = http.NewServeMux()

	// HTML 页面
	s.mux.HandleFunc(s.prefix+"/", s.handleIndex)

	// API 端点
	s.mux.HandleFunc(s.prefix+"/api/list", s.handleList)
	s.mux.HandleFunc(s.prefix+"/api/update", s.handleUpdate)
}

// Start 启动 HTTP 服务器
func (s *Server) Start() error {
	log.Printf("Feature HTTP server starting on %s", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}

// Handler 返回 HTTP handler，可以集成到现有的 HTTP 服务器中
func (s *Server) Handler() http.Handler {
	return s.mux
}

// handleIndex 处理首页请求，返回 HTML 页面
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	flags := s.getAllFlags()

	data := map[string]any{
		"Flags":  flags,
		"Prefix": s.prefix,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleList 处理获取所有 features 的 API 请求
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flags := s.getAllFlags()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flags); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleUpdate 处理更新 feature 值的 API 请求
func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	flag := features.Lookup(req.Name)
	if flag == nil {
		http.Error(w, fmt.Sprintf("Feature not found: %s", req.Name), http.StatusNotFound)
		return
	}

	// 检查是否标记为敏感字段，如果是则不允许修改
	if sensitive, ok := flag.Tags["sensitive"].(bool); ok && sensitive {
		http.Error(w, fmt.Sprintf("Feature %s is sensitive and cannot be modified", req.Name), http.StatusForbidden)
		return
	}

	// 尝试设置新值
	if err := flag.Value.Set(req.Value); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set value: %v", err), http.StatusBadRequest)
		return
	}

	// 返回更新后的值
	response := UpdateResponse{
		Success: true,
		Message: fmt.Sprintf("Feature %s updated successfully", req.Name),
		Flag:    s.flagToJSON(flag),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// getAllFlags 获取所有 flags 并转换为 JSON 格式
func (s *Server) getAllFlags() []FlagJSON {
	var flags []FlagJSON
	features.VisitAll(func(flag *features.Flag) {
		flags = append(flags, s.flagToJSON(flag))
	})
	return flags
}

// flagToJSON 将 Flag 转换为 JSON 格式
func (s *Server) flagToJSON(flag *features.Flag) FlagJSON {
	// 检查是否敏感
	isSensitive := false
	if sensitive, ok := flag.Tags["sensitive"].(bool); ok && sensitive {
		isSensitive = true
	}

	// 检查是否可修改
	isMutable := true
	if mutable, ok := flag.Tags["mutable"].(bool); ok && !mutable {
		isMutable = false
	}

	// 获取值
	var value any
	var valueStr string
	if isSensitive {
		value = "******"
		valueStr = "******"
	} else {
		value = flag.Value.Value()
		valueStr = flag.Value.String()
	}

	return FlagJSON{
		Name:        flag.Name,
		Usage:       flag.Usage,
		Type:        flag.Value.Type(),
		Value:       value,
		ValueString: valueStr,
		Deprecated:  flag.Deprecated,
		Tags:        flag.Tags,
		Sensitive:   isSensitive,
		Mutable:     isMutable,
	}
}

// FlagJSON 用于 JSON 序列化的 Flag 结构
type FlagJSON struct {
	Name        string         `json:"name"`
	Usage       string         `json:"usage"`
	Type        string         `json:"type"`
	Value       any            `json:"value"`
	ValueString string         `json:"valueString"`
	Deprecated  bool           `json:"deprecated"`
	Tags        map[string]any `json:"tags"`
	Sensitive   bool           `json:"sensitive"`
	Mutable     bool           `json:"mutable"`
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UpdateResponse 更新响应
type UpdateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Flag    FlagJSON `json:"flag"`
}
