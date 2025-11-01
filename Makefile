DOCKER_NAME ?= "devmgr"
DOCKER_TAG  ?= "latest"
VERSION     ?= "dev"

# 二进制文件名称
BINARY_NAME ?= "main"

GOBIN_DIR   = $(shell go env GOPATH)/bin
PATH_WITH_GOBIN = $(GOBIN_DIR):$(PATH)

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image: $(DOCKER_NAME):$(DOCKER_TAG) with version $(VERSION)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		--build-arg BINARY_NAME=$(BINARY_NAME) \
		-t $(DOCKER_NAME):$(DOCKER_TAG) .
	@echo "Docker image built successfully!"
	@echo "  Image: $(DOCKER_NAME):$(DOCKER_TAG)"
	@echo "  Version: $(VERSION)"
	@echo "  Binary: $(BINARY_NAME)"

.PHONY: clean-images
clean-images: ## 清理无用的Docker镜像
	@echo "清理无用的Docker镜像..."
	@echo "清理悬空镜像 (dangling images)..."
	@docker image prune -f
	@echo "清理项目相关的未使用镜像..."
	#@docker images ${DOCKER_NAME} --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}" | grep -v "$(VERSION)" | awk 'NR>1 {print $$3}' | xargs -r docker rmi --force|| true
	@docker images ${DOCKER_NAME} -aq | xargs -r docker rmi --force|| true
	@echo "Docker镜像清理完成"
	

