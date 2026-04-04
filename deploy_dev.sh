#!/bin/bash

# 本地开发部署脚本
# 功能：构建Docker镜像、启动/停止/重启容器、查看日志
# 支持分步执行和单服务操作

set -e

# ==================== 配置区域 ====================
LOCAL_PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${LOCAL_PROJECT_DIR}/docker-compose.dev.yml"

# 服务目录名
SERVER_BLOG_DIR="server-blog-v2"
SERVER_AUTH_DIR="server-auth-service"
WEB_BLOG_DIR="web-blog-v2"
WEB_AUTH_DIR="web-auth-service"
NGINX_DIR="nginx"

# Docker镜像名称
IMAGE_TAG="latest"
SERVER_BLOG_IMAGE="server-blog-dev:${IMAGE_TAG}"
SERVER_AUTH_IMAGE="server-auth-dev:${IMAGE_TAG}"
WEB_BLOG_IMAGE="web-blog-dev:${IMAGE_TAG}"
WEB_AUTH_IMAGE="web-auth-dev:${IMAGE_TAG}"
NGINX_IMAGE="nginx-proxy-dev:${IMAGE_TAG}"

# 服务列表
SERVICES=("server-blog" "server-auth" "web-blog" "web-auth" "nginx")

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ==================== 辅助函数 ====================
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 未安装，请先安装"
        exit 1
    fi
}

# 检查服务名是否有效
is_valid_service() {
    local service="$1"
    for s in "${SERVICES[@]}"; do
        if [ "$s" = "$service" ]; then
            return 0
        fi
    done
    return 1
}

# 获取容器名称
get_container_name() {
    local service="$1"
    echo "${service}-dev"
}

# ==================== 步骤函数 ====================

# 构建单个服务镜像（目录逻辑）
_build_service() {
    local service="$1"
    case "$service" in
        server-blog)
            log_info "构建 server-blog 镜像..."
            cd "${LOCAL_PROJECT_DIR}/${SERVER_BLOG_DIR}"
            docker build -t ${SERVER_BLOG_IMAGE} .
            ;;
        server-auth)
            log_info "构建 server-auth 镜像..."
            cd "${LOCAL_PROJECT_DIR}/${SERVER_AUTH_DIR}"
            docker build -t ${SERVER_AUTH_IMAGE} .
            ;;
        web-blog)
            log_info "构建 web-blog 镜像..."
            cd "${LOCAL_PROJECT_DIR}/${WEB_BLOG_DIR}"
            log_info "构建 Docker 镜像..."
            docker build -t ${WEB_BLOG_IMAGE} .
            ;;
        web-auth)
            log_info "构建 web-auth 镜像..."
            cd "${LOCAL_PROJECT_DIR}/${WEB_AUTH_DIR}"
            docker build -t ${WEB_AUTH_IMAGE} .
            ;;
        nginx)
            log_info "构建 nginx 镜像..."
            cd "${LOCAL_PROJECT_DIR}/${NGINX_DIR}"
            docker build -t ${NGINX_IMAGE} .
            ;;
    esac
}

# 构建镜像
step_build() {
    local service="$1"

    if [ -z "$service" ]; then
        log_info "=========================================="
        log_info "构建所有服务镜像"
        log_info "=========================================="
        for s in "${SERVICES[@]}"; do
            _build_service "$s"
        done
        log_info "所有服务镜像构建完成！"
    else
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        log_info "=========================================="
        log_info "构建服务镜像: $service"
        log_info "=========================================="
        _build_service "$service"
        log_info "服务镜像构建完成: $service"
    fi
}

# 启动容器
step_up() {
    local service="$1"
    
    if [ -z "$service" ]; then
        # 启动所有服务
        log_info "=========================================="
        log_info "启动所有服务"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" up -d
        log_info "所有服务启动完成！"
        log_info ""
        log_info "服务状态:"
        docker-compose -f "${COMPOSE_FILE}" ps
    else
        # 启动单个服务
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        
        log_info "=========================================="
        log_info "启动服务: $service"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" up -d "$service"
        log_info "服务启动完成: $service"
        log_info ""
        log_info "服务状态:"
        docker-compose -f "${COMPOSE_FILE}" ps "$service"
    fi
}

# 停止容器
step_down() {
    local service="$1"
    
    if [ -z "$service" ]; then
        # 停止所有服务
        log_info "=========================================="
        log_info "停止所有服务"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" down
        log_info "所有服务已停止！"
    else
        # 停止单个服务
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        
        log_info "=========================================="
        log_info "停止服务: $service"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" stop "$service"
        log_info "服务已停止: $service"
    fi
}

# 重启容器
step_restart() {
    local service="$1"
    
    if [ -z "$service" ]; then
        # 重启所有服务
        log_info "=========================================="
        log_info "重启所有服务"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" restart
        log_info "所有服务已重启！"
        log_info ""
        log_info "服务状态:"
        docker-compose -f "${COMPOSE_FILE}" ps
    else
        # 重启单个服务
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        
        log_info "=========================================="
        log_info "重启服务: $service"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" restart "$service"
        log_info "服务已重启: $service"
        log_info ""
        log_info "服务状态:"
        docker-compose -f "${COMPOSE_FILE}" ps "$service"
    fi
}

# 查看日志
step_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        # 查看所有服务日志
        log_info "=========================================="
        log_info "查看所有服务日志（按 Ctrl+C 退出）"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" logs -f
    else
        # 查看单个服务日志
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        
        log_info "=========================================="
        log_info "查看服务日志: $service（按 Ctrl+C 退出）"
        log_info "=========================================="
        cd "${LOCAL_PROJECT_DIR}"
        docker-compose -f "${COMPOSE_FILE}" logs -f "$service"
    fi
}

# 查看服务状态
step_ps() {
    local service="$1"
    
    log_info "=========================================="
    log_info "服务状态"
    log_info "=========================================="
    cd "${LOCAL_PROJECT_DIR}"
    if [ -z "$service" ]; then
        docker-compose -f "${COMPOSE_FILE}" ps
    else
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        docker-compose -f "${COMPOSE_FILE}" ps "$service"
    fi
}

# 完整发布（构建 + 停止旧容器 + 启动新容器）
step_release() {
    local service="$1"
    
    if [ -z "$service" ]; then
        log_info "=========================================="
        log_info "执行完整发布流程（所有服务）"
        log_info "=========================================="
        step_build
        step_down
        step_up
        log_info "=========================================="
        log_info "完整发布流程执行完成！"
        log_info "=========================================="
    else
        # 发布单个服务
        if ! is_valid_service "$service"; then
            log_error "未知服务: $service"
            log_error "支持的服务: ${SERVICES[*]}"
            exit 1
        fi
        
        log_info "=========================================="
        log_info "发布服务: $service（build + down + up）"
        log_info "=========================================="
        step_build "$service"
        step_down "$service"
        step_up "$service"
        log_info "=========================================="
        log_info "服务发布完成: $service"
        log_info "=========================================="
    fi
}

# 显示帮助信息
show_help() {
    cat << EOF
${GREEN}本地开发部署脚本${NC}

${YELLOW}用法:${NC}
  $0 <命令> [服务]

${YELLOW}命令:${NC}
  build [service]         - 构建镜像（不指定服务则构建所有）
  up [service]            - 启动容器（不指定服务则启动所有）
  down [service]          - 停止容器（不指定服务则停止所有）
  restart [service]       - 重启容器（不指定服务则重启所有）
  logs [service]          - 查看日志（不指定服务则查看所有）
  ps [service]            - 查看服务状态（不指定服务则查看所有）
  release [service]       - 发布服务（停止 + 构建 + 启动，不指定服务则发布所有）

${YELLOW}支持的服务:${NC}
  ${SERVICES[*]}

${YELLOW}示例:${NC}
  $0 build                        # 构建所有服务镜像
  $0 build server-blog            # 仅构建 server-blog 镜像
  $0 up                           # 启动所有服务
  $0 up web-blog                  # 仅启动 web-blog 服务
  $0 down                         # 停止所有服务
  $0 down server-auth             # 仅停止 server-auth 服务
  $0 restart                      # 重启所有服务
  $0 restart nginx                # 仅重启 nginx 服务
  $0 logs                         # 查看所有服务日志
  $0 logs server-blog             # 查看 server-blog 日志
  $0 ps                           # 查看所有服务状态
  $0 ps web-auth                  # 查看 web-auth 服务状态
  $0 release                      # 发布所有服务
  $0 release server-blog          # 发布 server-blog 服务

EOF
}

# ==================== 主流程 ====================

# 检查必要的命令
check_command docker

# 检查参数
if [ $# -eq 0 ]; then
    log_error "必须指定执行命令！"
    echo ""
    show_help
    exit 1
fi

# 解析参数
COMMAND=$1
SERVICE=$2

case ${COMMAND} in
    build)
        step_build "$SERVICE"
        ;;
    up)
        step_up "$SERVICE"
        ;;
    down)
        step_down "$SERVICE"
        ;;
    restart)
        step_restart "$SERVICE"
        ;;
    logs)
        step_logs "$SERVICE"
        ;;
    ps)
        step_ps "$SERVICE"
        ;;
    release)
        step_release "$SERVICE"
        ;;
    help|-h|--help)
        show_help
        ;;
    *)
        log_error "未知命令: ${COMMAND}"
        echo ""
        show_help
        exit 1
        ;;
esac
