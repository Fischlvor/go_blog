#!/bin/bash

# CI/CD 部署脚本
# 在远程服务器上执行：从Git仓库拉取代码，构建Docker镜像并部署

set -e  # 遇到错误立即退出

# ==================== 配置区域 ====================
# Git仓库配置
GIT_REPO="git@gitee.com:qiyana423/go_blog.git"
GIT_BRANCH="master"  # 或 "main"，根据实际情况修改

# 远程服务器路径
BASE_DIR="/media/practice/onServer/go_blog"
PROJECT_DIR="${BASE_DIR}/source"
COMPOSE_DIR="${BASE_DIR}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 未安装，请先安装"
        exit 1
    fi
}

# ==================== 主流程 ====================

# 1. 检查必要的命令
log_info "检查必要的命令..."
check_command docker
check_command docker-compose
check_command git

# 2. 创建必要的目录
log_info "创建必要的目录..."
mkdir -p ${PROJECT_DIR}
mkdir -p ${COMPOSE_DIR}
mkdir -p ${COMPOSE_DIR}/nginx
mkdir -p ${BASE_DIR}/server/configs
mkdir -p ${BASE_DIR}/server/uploads
mkdir -p ${BASE_DIR}/server/log
mkdir -p ${BASE_DIR}/auth-server/configs
mkdir -p ${BASE_DIR}/auth-server/log
mkdir -p ${BASE_DIR}/auth-server/keys

# 3. 克隆或更新代码
log_info "从Git仓库获取代码..."
cd ${PROJECT_DIR}

if [ -d ".git" ]; then
    log_info "代码仓库已存在，更新代码..."
    git fetch origin
    git reset --hard origin/${GIT_BRANCH}
    git clean -fd
else
    log_info "克隆代码仓库..."
    git clone -b ${GIT_BRANCH} ${GIT_REPO} .
fi

log_info "代码更新完成！当前分支: $(git rev-parse --abbrev-ref HEAD)"
log_info "最新提交: $(git log -1 --oneline)"

# 4. 检查配置文件是否存在
log_info "检查配置文件..."

if [ ! -f "${BASE_DIR}/server/configs/config.yaml" ]; then
    log_warn "server-blog 配置文件不存在，请确保已配置: ${BASE_DIR}/server/configs/config.yaml"
fi

if [ ! -f "${BASE_DIR}/auth-server/configs/config.yaml" ]; then
    log_warn "server-auth-service 配置文件不存在，请确保已配置: ${BASE_DIR}/auth-server/configs/config.yaml"
fi

# 5. 构建Docker镜像
log_info "=========================================="
log_info "开始构建Docker镜像..."
log_info "=========================================="

# 构建 server-blog 镜像
log_info "构建 server-blog 镜像..."
cd ${PROJECT_DIR}/server-blog
docker build -t go-blog-server:latest .
cd ${PROJECT_DIR}

# 构建 server-auth-service 镜像
log_info "构建 server-auth-service 镜像..."
cd ${PROJECT_DIR}/server-auth-service
docker build -t auth-service-server:latest .
cd ${PROJECT_DIR}

# 构建 web-blog 镜像
log_info "构建 web-blog 镜像..."
cd ${PROJECT_DIR}/web-blog
docker build -t web-blog:latest .
cd ${PROJECT_DIR}

# 构建 web-auth-service 镜像
log_info "构建 web-auth-service 镜像..."
cd ${PROJECT_DIR}/web-auth-service
docker build -t web-auth-service:latest .
cd ${PROJECT_DIR}

# 构建 nginx 镜像
log_info "构建 nginx 镜像..."
cd ${PROJECT_DIR}/nginx
docker build -t nginx-proxy:latest .
cd ${PROJECT_DIR}

log_info "=========================================="
log_info "所有镜像构建完成！"
log_info "=========================================="

# 6. 修改docker-compose.yml中的build context路径
log_info "准备docker-compose文件..."
cd ${PROJECT_DIR}

# 创建部署用的docker-compose文件，修改build context为绝对路径
sed "s|context: ./server-blog|context: ${PROJECT_DIR}/server-blog|g" docker-compose.yml > ${COMPOSE_DIR}/docker-compose.yml
sed -i "s|context: ./server-auth-service|context: ${PROJECT_DIR}/server-auth-service|g" ${COMPOSE_DIR}/docker-compose.yml
sed -i "s|context: ./web-blog|context: ${PROJECT_DIR}/web-blog|g" ${COMPOSE_DIR}/docker-compose.yml
sed -i "s|context: ./web-auth-service|context: ${PROJECT_DIR}/web-auth-service|g" ${COMPOSE_DIR}/docker-compose.yml
sed -i "s|context: ./nginx|context: ${PROJECT_DIR}/nginx|g" ${COMPOSE_DIR}/docker-compose.yml

# 复制基础服务compose文件
cp ${PROJECT_DIR}/docker-compose.base.yml ${COMPOSE_DIR}/

# 7. 复制nginx配置文件
log_info "复制nginx配置文件..."
cp ${PROJECT_DIR}/nginx/go_blog.conf ${COMPOSE_DIR}/nginx/

# 8. 停止旧容器
log_info "停止旧容器..."
cd ${COMPOSE_DIR}
docker-compose -f docker-compose.yml down 2>/dev/null || true
docker-compose -f docker-compose.base.yml down 2>/dev/null || true

# 9. 启动基础服务
log_info "启动基础服务..."
docker-compose -f docker-compose.base.yml up -d

# 10. 等待基础服务就绪
log_info "等待基础服务就绪..."
sleep 10

# 11. 启动业务服务
log_info "启动业务服务..."
docker-compose -f docker-compose.yml up -d

# 12. 显示服务状态
log_info "=========================================="
log_info "部署完成！"
log_info "=========================================="
log_info ""
log_info "服务状态："
docker-compose -f docker-compose.yml ps

log_info ""
log_info "查看日志:"
log_info "  cd ${COMPOSE_DIR} && docker-compose logs -f"
log_info ""
log_info "查看所有镜像:"
log_info "  docker images | grep -E \"go-blog|auth-service|web-blog|web-auth|nginx-proxy\""
