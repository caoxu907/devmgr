ARG VERSION=dev
ARG BUILD_TIME
ARG BINARY_NAME=devmgr

# =================================================================
# 第一阶段：构建阶段 (builder)
# =================================================================

FROM 10.17.196.52:12888/go-base-alpine:v1 AS builder

# 继承构建参数
#ARG VERSION
#ARG BUILD_TIME
#ARG BINARY_NAME

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl -ldflags="-s -w" -o main .

# =================================================================
# 第二阶段：运行阶段 (runtime)
# =================================================================

FROM 10.17.196.52:12888/alpine:3.8

# 继承构建参数，用于设置环境变量或标签
ARG VERSION
ARG BUILD_TIME
ARG BINARY_NAME

WORKDIR /app
COPY --from=builder /app/main .
COPY resource ./resource
COPY manifest ./manifest
COPY manifest/config/config_1.yaml ./manifest/config/config.yaml
CMD ["./main"]
