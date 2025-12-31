package featurewatcher

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/pubgo/funk/v2/closer"
	"github.com/pubgo/funk/v2/features"
	"github.com/pubgo/funk/v2/log"
)

// Watcher 监控指定目录的文件变化，并自动加载到 features
type Watcher struct {
	dir      string
	feature  *features.Feature
	watcher  *fsnotify.Watcher
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	debounce time.Duration
	mu       sync.RWMutex
	fileMap  map[string]time.Time // 用于防抖，记录文件最后修改时间
}

// Option 配置选项
type Option func(*Watcher)

// WithFeature 指定要使用的 Feature 实例，默认使用全局的 defaultFeature
func WithFeature(f *features.Feature) Option {
	return func(w *Watcher) {
		w.feature = f
	}
}

// WithDebounce 设置防抖时间，避免频繁触发
func WithDebounce(d time.Duration) Option {
	return func(w *Watcher) {
		w.debounce = d
	}
}

// NewWatcher 创建一个新的文件监控器
func NewWatcher(dir string, opts ...Option) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := &Watcher{
		dir:      dir,
		feature:  nil, // 默认使用全局的
		watcher:  watcher,
		ctx:      ctx,
		cancel:   cancel,
		debounce: 500 * time.Millisecond, // 默认 500ms 防抖
		fileMap:  make(map[string]time.Time),
	}

	// 应用选项
	for _, opt := range opts {
		opt(w)
	}

	return w, nil
}

// Start 开始监控目录
func (w *Watcher) Start() error {
	// 检查目录是否存在
	if _, err := os.Stat(w.dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", w.dir)
	}

	// 添加目录到监控
	if err := w.watcher.Add(w.dir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	// 初始加载所有文件
	if err := w.loadAllFiles(); err != nil {
		log.Err(err).Str("dir", w.dir).Msg("failed to load initial files")
	}

	// 启动监控 goroutine
	w.wg.Add(1)
	go w.watchLoop()

	log.Info().Str("dir", w.dir).Msg("Feature watcher started")

	return nil
}

// Stop 停止监控
func (w *Watcher) Stop() error {
	w.cancel()
	closer.SafeClose(w.watcher)
	w.wg.Wait()
	log.Info().Str("dir", w.dir).Msg("Feature watcher stopped")
	return nil
}

// watchLoop 监控循环
func (w *Watcher) watchLoop() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Err(err).Msg("file watcher error")
		}
	}
}

// handleEvent 处理文件事件
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// 只处理文件写入和创建事件
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return
	}

	// skip ~ suffix
	if strings.HasSuffix(event.Name, "~") {
		return
	}

	// 跳过目录
	if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
		return
	}

	// 防抖处理
	w.mu.Lock()
	lastTime, exists := w.fileMap[event.Name]
	now := time.Now()
	if exists && now.Sub(lastTime) < w.debounce {
		w.mu.Unlock()
		return
	}
	w.fileMap[event.Name] = now
	w.mu.Unlock()

	// 延迟加载，避免文件正在写入
	time.Sleep(100 * time.Millisecond)

	// 加载文件
	if err := w.loadFile(event.Name); err != nil {
		log.Err(err).Str("file", event.Name).Msg("failed to load file")
	} else {
		log.Info().Str("file", event.Name).Msg("file reloaded")
	}
}

// loadAllFiles 加载目录中的所有文件
func (w *Watcher) loadAllFiles() error {
	return filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return w.loadFile(path)
	})
}

// loadFile 加载单个文件并解析
func (w *Watcher) loadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closer.SafeClose(file)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	loaded := 0
	errors := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// 解析 name=value 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Warn().
				Str("file", filePath).
				Int("line", lineNum).
				Str("content", line).
				Msg("invalid format, expected 'name=value'")
			errors++
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if name == "" {
			log.Warn().
				Str("file", filePath).
				Int("line", lineNum).
				Msg("empty feature name")
			errors++
			continue
		}

		// 查找对应的 feature
		var flag *features.Flag
		if w.feature != nil {
			flag = w.feature.Lookup(name)
		} else {
			flag = features.Lookup(name)
		}

		if flag == nil {
			log.Debug().
				Str("file", filePath).
				Str("name", name).
				Msg("feature not found, skipping")
			continue
		}

		// 检查是否可修改
		if mutable, ok := flag.Tags["mutable"].(bool); ok && !mutable {
			log.Debug().
				Str("file", filePath).
				Str("name", name).
				Msg("feature is not mutable, skipping")
			continue
		}

		// 设置值
		if err := flag.Value.Set(value); err != nil {
			log.Err(err).
				Str("file", filePath).
				Str("name", name).
				Str("value", value).
				Int("line", lineNum).
				Msg("failed to set feature value")
			errors++
			continue
		}

		loaded++
		log.Debug().
			Str("file", filePath).
			Str("name", name).
			Str("value", value).
			Msg("feature updated")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	log.Info().
		Str("file", filePath).
		Int("loaded", loaded).
		Int("errors", errors).
		Int("total_lines", lineNum).
		Msg("file loaded")

	return nil
}
