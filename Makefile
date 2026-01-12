# Stock Query System Makefile
# 股票查询系统构建脚本

# 变量定义
BACKEND_IMAGE = stock-api
BACKEND_TAG = latest
BACKEND_CONTAINER = stock-backend
FRONTEND_CONTAINER = stock-frontend

BACKEND_DIR = backend
FRONTEND_DIR = frontend/stock-frontend

# 默认目标
.PHONY: help
help:
	@echo "Stock Query System - Available Commands:"
	@echo "  make build    - 构建后端Docker镜像"
	@echo "  make stop     - 停止前端和后端容器"
	@echo "  make clean    - 删除Docker镜像和容器"
	@echo "  make run      - 启动前端和后端服务"
	@echo "  make logs     - 查看后端容器日志"
	@echo "  make status   - 查看容器状态"

# 构建后端Docker镜像
.PHONY: build
build:
	@echo "🔨 构建后端Docker镜像..."
	cd $(BACKEND_DIR) && docker build -t $(BACKEND_IMAGE):$(BACKEND_TAG) .
	@echo "✅ 后端镜像构建完成"

# 停止所有容器
.PHONY: stop
stop:
	@echo "🛑 停止所有容器..."
	-docker stop $(BACKEND_CONTAINER) 2>/dev/null || true
	-docker stop $(FRONTEND_CONTAINER) 2>/dev/null || true
	-pkill -f "npm run dev" 2>/dev/null || true
	-pkill -f "vite" 2>/dev/null || true
	@echo "✅ 所有服务已停止"

# 清理Docker镜像和容器
.PHONY: clean
clean: stop
	@echo "🧹 清理Docker资源..."
	-docker rm $(BACKEND_CONTAINER) 2>/dev/null || true
	-docker rm $(FRONTEND_CONTAINER) 2>/dev/null || true
	-docker rmi $(BACKEND_IMAGE):$(BACKEND_TAG) 2>/dev/null || true
	@echo "✅ 清理完成"

# 启动前端和后端服务
.PHONY: run
run: stop build
	@echo "🚀 启动服务..."
	@echo "📡 启动后端服务..."
	docker run -d -p 8080:8080 --name $(BACKEND_CONTAINER) $(BACKEND_IMAGE):$(BACKEND_TAG)
	@echo "⏳ 等待后端服务启动..."
	sleep 5
	@echo "🌐 启动前端开发服务器..."
	cd $(FRONTEND_DIR) && nohup npm run dev > /dev/null 2>&1 &
	@echo "⏳ 等待前端服务启动..."
	sleep 3
	@echo ""
	@echo "✅ 服务启动完成！"
	@echo "🔗 后端API: http://localhost:8080"
	@echo "🔗 前端界面: http://localhost:5173"
	@echo "📊 健康检查: curl http://localhost:8080/health"

# 查看后端日志
.PHONY: logs
logs:
	@echo "📋 后端容器日志:"
	docker logs -f $(BACKEND_CONTAINER)

# 查看容器状态
.PHONY: status
status:
	@echo "📊 容器状态:"
	@echo "后端容器:"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" --filter name=$(BACKEND_CONTAINER) || echo "  未运行"
	@echo ""
	@echo "前端进程:"
	@ps aux | grep -E "(npm run dev|vite)" | grep -v grep || echo "  未运行"

# 重启服务
.PHONY: restart
restart: stop run

# 开发模式（仅启动后端Docker，前端本地开发）
.PHONY: dev
dev: stop
	@echo "🛠️  开发模式启动..."
	@echo "📡 启动后端Docker容器..."
	docker run -d -p 8080:8080 --name $(BACKEND_CONTAINER) $(BACKEND_IMAGE):$(BACKEND_TAG)
	@echo "✅ 后端服务已启动"
	@echo "💡 请手动启动前端: cd $(FRONTEND_DIR) && npm run dev"

# 生产构建
.PHONY: prod-build
prod-build: build
	@echo "🏗️  生产环境构建..."
	cd $(FRONTEND_DIR) && npm run build
	@echo "✅ 前端生产版本构建完成"