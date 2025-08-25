#!/bin/bash

# CI/CD 部署脚本
# 用于自动构建前端代码并部署到服务器

set -e  # 遇到错误立即退出

# 配置变量
SERVER_HOST="8.148.64.96"
SERVER_USER="root"
SERVER_PORT="22"
SSH_KEY_PATH="$HOME/.ssh/id_rsa_go_blog"
REMOTE_DIR="/media/practice/onServer/go_blog/web"
LOCAL_DIST_DIR="web/dist"
PROJECT_NAME="go_blog"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查SSH密钥是否存在
check_ssh_key() {
    log_info "检查SSH密钥..."
    if [ ! -f "$SSH_KEY_PATH" ]; then
        log_error "SSH密钥不存在: $SSH_KEY_PATH"
        log_info "请先生成SSH密钥: ssh-keygen -t rsa -b 4096 -C 'your_email@example.com'"
        exit 1
    fi
    log_success "SSH密钥检查通过"
}

# 测试SSH连接
test_ssh_connection() {
    log_info "测试SSH连接..."
    if ssh -i "$SSH_KEY_PATH" -p "$SERVER_PORT" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$SERVER_USER@$SERVER_HOST" "echo 'SSH连接成功'" 2>/dev/null; then
        log_success "SSH连接测试成功"
    else
        log_error "SSH连接失败，请检查："
        log_error "1. 服务器地址: $SERVER_HOST"
        log_error "2. 用户名: $SERVER_USER"
        log_error "3. SSH密钥: $SSH_KEY_PATH"
        log_error "4. 网络连接"
        exit 1
    fi
}

# 检查Node.js和npm
check_node_environment() {
    log_info "检查Node.js环境..."
    
    # 检查nvm是否存在
    if ! command -v nvm &> /dev/null; then
        # 尝试加载nvm
        if [ -f "$HOME/.nvm/nvm.sh" ]; then
            source "$HOME/.nvm/nvm.sh"
        elif [ -f "/usr/local/nvm/nvm.sh" ]; then
            source "/usr/local/nvm/nvm.sh"
        else
            log_error "nvm未安装或未找到"
            log_info "请安装nvm: curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash"
            exit 1
        fi
    fi
    
    # 切换到指定Node.js版本
    log_info "切换到Node.js 22.15.0..."
    nvm use 22.15.0
    if [ $? -ne 0 ]; then
        log_error "Node.js 22.15.0未安装"
        log_info "请先安装: nvm install 22.15.0"
        exit 1
    fi
    
    # 重新检查Node.js和npm
    if ! command -v node &> /dev/null; then
        log_error "Node.js未安装"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        log_error "npm未安装"
        exit 1
    fi
    
    log_success "Node.js环境检查通过"
    log_info "Node.js版本: $(node --version)"
    log_info "npm版本: $(npm --version)"
}

# 构建前端项目
build_frontend() {
    log_info "开始构建前端项目..."
    
    # 进入web目录
    cd web
    
    # 安装依赖
    log_info "安装npm依赖..."
    npm install
    
    # 构建项目
    log_info "执行npm run build..."
    npm run build
    
    # 检查构建结果
    if [ ! -d "dist" ]; then
        log_error "构建失败，dist目录不存在"
        exit 1
    fi
    
    log_success "前端项目构建成功"
    
    # 返回项目根目录
    cd ..
}

# 备份远程文件
backup_remote_files() {
    log_info "备份远程文件..."
    local backup_dir="/tmp/${PROJECT_NAME}_backup_$(date +%Y%m%d_%H%M%S)"
    
    ssh -i "$SSH_KEY_PATH" -p "$SERVER_PORT" "$SERVER_USER@$SERVER_HOST" << EOF
        if [ -d "$REMOTE_DIR" ]; then
            mkdir -p $backup_dir
            cp -r $REMOTE_DIR/* $backup_dir/ 2>/dev/null || true
            echo "备份完成: $backup_dir"
        else
            echo "远程目录不存在，无需备份"
        fi
EOF
    log_success "远程文件备份完成"
}

# 部署到服务器
deploy_to_server() {
    log_info "开始部署到服务器..."
    
    # 创建远程目录（如果不存在）
    ssh -i "$SSH_KEY_PATH" -p "$SERVER_PORT" "$SERVER_USER@$SERVER_HOST" "mkdir -p $REMOTE_DIR"
    
    # 上传文件
    log_info "上传文件到服务器..."
    scp -i "$SSH_KEY_PATH" -P "$SERVER_PORT" -r "$LOCAL_DIST_DIR" "$SERVER_USER@$SERVER_HOST:$REMOTE_DIR/"
    
    if [ $? -eq 0 ]; then
        log_success "文件上传成功"
    else
        log_error "文件上传失败"
        exit 1
    fi
}

# 设置文件权限
set_permissions() {
    log_info "设置文件权限..."
    ssh -i "$SSH_KEY_PATH" -p "$SERVER_PORT" "$SERVER_USER@$SERVER_HOST" << EOF
        chmod -R 755 $REMOTE_DIR
        chown -R www-data:www-data $REMOTE_DIR 2>/dev/null || true
        echo "权限设置完成"
EOF
    log_success "文件权限设置完成"
}

# 验证部署
verify_deployment() {
    log_info "验证部署..."
    
    # 检查远程文件是否存在
    local file_count=$(ssh -i "$SSH_KEY_PATH" -p "$SERVER_PORT" "$SERVER_USER@$SERVER_HOST" "find $REMOTE_DIR -type f | wc -l")
    
    if [ "$file_count" -gt 0 ]; then
        log_success "部署验证成功，共部署 $file_count 个文件"
    else
        log_warning "部署验证失败，远程目录为空"
    fi
}

# 清理本地构建文件
cleanup_local() {
    log_info "清理本地构建文件..."
    if [ -d "$LOCAL_DIST_DIR" ]; then
        rm -rf "$LOCAL_DIST_DIR"
        log_success "本地构建文件清理完成"
    fi
}

# 主函数
main() {
    log_info "开始部署 $PROJECT_NAME 项目..."
    log_info "服务器: $SERVER_USER@$SERVER_HOST:$SERVER_PORT"
    log_info "远程目录: $REMOTE_DIR"
    
    # 执行部署步骤
    check_ssh_key
    test_ssh_connection
    check_node_environment
    build_frontend
    backup_remote_files
    deploy_to_server
    set_permissions
    verify_deployment
    cleanup_local
    
    log_success "部署完成！"
    log_info "访问地址: http://$SERVER_HOST"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -s, --server   指定服务器地址 (默认: $SERVER_HOST)"
    echo "  -u, --user     指定用户名 (默认: $SERVER_USER)"
    echo "  -p, --port     指定SSH端口 (默认: $SERVER_PORT)"
    echo "  -k, --key      指定SSH密钥路径 (默认: $SSH_KEY_PATH)"
    echo "  -d, --dir      指定远程目录 (默认: $REMOTE_DIR)"
    echo ""
    echo "示例:"
    echo "  $0"
    echo "  $0 -s 192.168.1.100 -u deploy -d /var/www/html"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -s|--server)
            SERVER_HOST="$2"
            shift 2
            ;;
        -u|--user)
            SERVER_USER="$2"
            shift 2
            ;;
        -p|--port)
            SERVER_PORT="$2"
            shift 2
            ;;
        -k|--key)
            SSH_KEY_PATH="$2"
            shift 2
            ;;
        -d|--dir)
            REMOTE_DIR="$2"
            shift 2
            ;;
        *)
            log_error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 执行主函数
main "$@" 