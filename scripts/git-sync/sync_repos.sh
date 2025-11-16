#!/bin/bash

# Git 双仓库智能同步脚本
# Gitee (主仓库) <-> GitHub (镜像仓库)
# 特性：
# 1. 只有 GitHub 走代理（已在 git config 中配置）
# 2. GitHub 推送失败不影响 Gitee
# 3. 支持推送所有分支或指定分支

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
GITEE_REMOTE="origin"
GITHUB_REMOTE="github"
PROXY_HOST="10.21.71.52"
PROXY_PORT="7890"

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

# 检查代理是否可用
check_proxy() {
    log_info "检查代理连接 ${PROXY_HOST}:${PROXY_PORT}..."
    if timeout 3 nc -zv ${PROXY_HOST} ${PROXY_PORT} &>/dev/null; then
        log_success "代理可用"
        return 0
    else
        log_warning "代理不可用，GitHub 推送将跳过"
        return 1
    fi
}

# 推送到 Gitee
push_to_gitee() {
    local branch=$1
    log_info "推送到 Gitee (${GITEE_REMOTE})..."
    
    if [ -z "$branch" ]; then
        # 推送所有分支
        if git push ${GITEE_REMOTE} --all; then
            log_success "Gitee: 所有分支推送成功"
            git push ${GITEE_REMOTE} --tags 2>/dev/null || true
            return 0
        else
            log_error "Gitee: 推送失败"
            return 1
        fi
    else
        # 推送指定分支
        if git push ${GITEE_REMOTE} ${branch}; then
            log_success "Gitee: 分支 ${branch} 推送成功"
            return 0
        else
            log_error "Gitee: 分支 ${branch} 推送失败"
            return 1
        fi
    fi
}

# 推送到 GitHub
push_to_github() {
    local branch=$1
    
    # 检查代理
    if ! check_proxy; then
        log_warning "跳过 GitHub 推送（代理不可用）"
        return 0
    fi
    
    log_info "推送到 GitHub (${GITHUB_REMOTE})..."
    
    if [ -z "$branch" ]; then
        # 推送所有分支
        if timeout 30 git push ${GITHUB_REMOTE} --all 2>&1; then
            log_success "GitHub: 所有分支推送成功"
            timeout 30 git push ${GITHUB_REMOTE} --tags 2>/dev/null || true
            return 0
        else
            log_warning "GitHub: 推送失败（可能是网络问题）"
            return 1
        fi
    else
        # 推送指定分支
        if timeout 30 git push ${GITHUB_REMOTE} ${branch} 2>&1; then
            log_success "GitHub: 分支 ${branch} 推送成功"
            return 0
        else
            log_warning "GitHub: 分支 ${branch} 推送失败（可能是网络问题）"
            return 1
        fi
    fi
}

# 主函数
main() {
    echo "========================================"
    echo "  Git 双仓库同步脚本"
    echo "========================================"
    echo ""
    
    # 检查是否在 git 仓库中
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "当前目录不是 Git 仓库"
        exit 1
    fi
    
    # 检查 remote 是否存在
    if ! git remote | grep -q "^${GITEE_REMOTE}$"; then
        log_error "Gitee remote '${GITEE_REMOTE}' 不存在"
        exit 1
    fi
    
    if ! git remote | grep -q "^${GITHUB_REMOTE}$"; then
        log_error "GitHub remote '${GITHUB_REMOTE}' 不存在"
        exit 1
    fi
    
    # 获取当前分支
    CURRENT_BRANCH=$(git branch --show-current)
    log_info "当前分支: ${CURRENT_BRANCH}"
    
    # 解析参数
    BRANCH=""
    PUSH_ALL=false
    SYNC_MODE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --all|-a)
                PUSH_ALL=true
                shift
                ;;
            --branch|-b)
                BRANCH="$2"
                shift 2
                ;;
            --sync|-s)
                SYNC_MODE=true
                shift
                ;;
            --help|-h)
                echo "用法: $0 [选项]"
                echo ""
                echo "选项:"
                echo "  --all, -a           推送所有分支"
                echo "  --branch, -b BRANCH 推送指定分支"
                echo "  --sync, -s          增量同步（推送缺失分支和新提交）"
                echo "  --help, -h          显示帮助信息"
                echo ""
                echo "示例:"
                echo "  $0                  # 推送当前分支"
                echo "  $0 --all            # 推送所有分支"
                echo "  $0 --sync           # 增量同步缺失的分支"
                echo "  $0 --branch master  # 推送 master 分支"
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                echo "使用 --help 查看帮助"
                exit 1
                ;;
        esac
    done
    
    # 增量同步模式
    if [ "$SYNC_MODE" = true ]; then
        log_info "模式: 增量同步（推送缺失分支和新提交）"
        echo ""
        
        # 检查代理
        if ! check_proxy; then
            log_error "代理不可用，无法执行增量同步"
            exit 1
        fi
        
        # 获取 GitHub 远程分支
        log_info "获取 GitHub 远程分支信息..."
        timeout 30 git fetch ${GITHUB_REMOTE} &>/dev/null || {
            log_error "无法连接到 GitHub"
            exit 1
        }
        
        # 找出 Gitee 有但 GitHub 没有的分支
        MISSING_BRANCHES=$(comm -23 \
            <(git branch -r | grep "${GITEE_REMOTE}/" | sed "s/${GITEE_REMOTE}\///" | grep -v 'HEAD' | sort) \
            <(git branch -r | grep "${GITHUB_REMOTE}/" | sed "s/${GITHUB_REMOTE}\///" | sort))
        
        # 推送缺失的分支
        if [ -n "$MISSING_BRANCHES" ]; then
            log_info "发现 GitHub 缺失的分支:"
            echo "$MISSING_BRANCHES" | while read branch; do
                echo "  - $branch"
            done
            echo ""
            
            echo "$MISSING_BRANCHES" | while read branch; do
                log_info "推送分支: ${branch}"
                if timeout 30 git push ${GITHUB_REMOTE} refs/remotes/${GITEE_REMOTE}/${branch}:refs/heads/${branch} 2>&1; then
                    log_success "GitHub: 分支 ${branch} 推送成功"
                else
                    log_warning "GitHub: 分支 ${branch} 推送失败"
                fi
            done
            echo ""
        fi
        
        # 检查现有分支是否有新提交
        log_info "检查现有分支的新提交..."
        HAS_NEW_COMMITS=false
        for branch in $(git branch -r | grep "${GITEE_REMOTE}/" | sed "s/${GITEE_REMOTE}\///" | grep -v 'HEAD'); do
            # 检查 GitHub 是否有这个分支
            if git rev-parse --verify ${GITHUB_REMOTE}/${branch} &>/dev/null; then
                # 检查是否有新提交
                if git rev-list ${GITHUB_REMOTE}/${branch}..${GITEE_REMOTE}/${branch} 2>/dev/null | grep -q .; then
                    log_info "分支 ${branch} 有新提交需要同步"
                    HAS_NEW_COMMITS=true
                    if timeout 30 git push ${GITHUB_REMOTE} refs/remotes/${GITEE_REMOTE}/${branch}:refs/heads/${branch} 2>&1; then
                        log_success "GitHub: 分支 ${branch} 更新成功"
                    else
                        log_warning "GitHub: 分支 ${branch} 更新失败"
                    fi
                fi
            fi
        done
        
        if [ "$HAS_NEW_COMMITS" = false ] && [ -z "$MISSING_BRANCHES" ]; then
            log_success "所有分支都是最新的，无需同步"
        fi
        
        # 同步 tags
        log_info "同步 tags..."
        timeout 30 git push ${GITHUB_REMOTE} --tags 2>/dev/null || log_warning "Tags 同步失败"
        
    else
        # 原有的推送逻辑
        # 确定要推送的分支
        if [ "$PUSH_ALL" = true ]; then
            TARGET_BRANCH=""
            log_info "模式: 推送所有分支"
        elif [ -n "$BRANCH" ]; then
            TARGET_BRANCH="$BRANCH"
            log_info "模式: 推送指定分支 ${TARGET_BRANCH}"
        else
            TARGET_BRANCH="$CURRENT_BRANCH"
            log_info "模式: 推送当前分支 ${TARGET_BRANCH}"
        fi
        
        echo ""
        
        # 推送到 Gitee（主仓库，必须成功）
        if ! push_to_gitee "$TARGET_BRANCH"; then
            log_error "Gitee 推送失败，终止操作"
            exit 1
        fi
        
        echo ""
        
        # 推送到 GitHub（镜像仓库，失败不影响）
        push_to_github "$TARGET_BRANCH" || true
    fi
    
    echo ""
    echo "========================================"
    log_success "同步完成"
    echo "========================================"
}

# 执行主函数
main "$@"
