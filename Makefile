# ============================================================================
#  Go 跨平台编译 Makefile
#  用法:
#    make build          - 编译当前平台版本
#    make build-all      - 编译所有预定义平台
#    make build-linux    - 只编译 Linux amd64
#    make build-windows  - 只编译 Windows amd64
#    make build-darwin   - 只编译 macOS amd64
#    make clean          - 删除编译产物
#    make help           - 显示帮助信息
# ============================================================================

# 项目名称（可修改）
APP_NAME   := termhelper
# 主包路径（如果 main.go 在项目根目录，则留空；否则指定目录）
MAIN_PKG   := .
# 输出目录
BUILD_DIR  := ./bin
# 版本信息（可选，用于注入 ldflags）
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 静态编译核心设置（禁用 CGO，强制纯 Go 静态链接）
CGO_ENABLED := 0

# 编译参数：-s -w 去除调试信息，减小体积；-X 注入版本信息
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

# 支持的平台列表（GOOS/GOARCH）
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 windows/386 darwin/amd64 darwin/arm64

# ============================================================================
#  规则定义
# ============================================================================

.PHONY: all build build-all clean help $(PLATFORMS)

# 默认目标：构建当前平台的二进制
all: build

# 构建当前平台（自动检测 GOOS/GOARCH）
build:
	@echo "Building $(APP_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

# 构建所有预定义平台
build-all: $(PLATFORMS)

# 针对每个平台生成构建目标
$(PLATFORMS):
	@$(eval GOOS := $(word 1, $(subst /, ,$@)))
	@$(eval GOARCH := $(word 2, $(subst /, ,$@)))
	@$(eval OUTPUT := $(BUILD_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH))
	@if [ "$(GOOS)" = "windows" ]; then \
		OUTPUT="$(OUTPUT).exe"; \
	fi
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) $(MAIN_PKG)

# 单独构建各平台的快捷命令（方便调试）
build-linux: linux/amd64
build-windows: windows/amd64
build-darwin: darwin/amd64

# 清理构建产物
clean:
	@echo "Cleaning $(BUILD_DIR)..."
	@rm -rf $(BUILD_DIR)

# 显示帮助
help:
	@echo "可用目标："
	@echo "  build          构建当前平台版本"
	@echo "  build-all      构建所有预定义平台（见 PLATFORMS）"
	@echo "  build-linux    仅构建 Linux amd64"
	@echo "  build-windows  仅构建 Windows amd64"
	@echo "  build-darwin   仅构建 macOS amd64"
	@echo "  clean          删除 $(BUILD_DIR) 目录"
	@echo "  help           显示此帮助"
	@echo ""
	@echo "当前配置："
	@echo "  APP_NAME   = $(APP_NAME)"
	@echo "  VERSION    = $(VERSION)"
	@echo "  BUILD_TIME = $(BUILD_TIME)"
	@echo "  LDFLAGS    = $(LDFLAGS)"
