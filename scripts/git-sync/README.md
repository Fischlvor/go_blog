# Git 双仓库同步工具

自动同步 Gitee 和 GitHub 仓库的脚本工具。

## 特性

- ✅ 只有 GitHub 走代理，Gitee 不受影响
- ✅ 代理失效时不影响 Gitee 推送
- ✅ 支持增量同步，自动检测缺失的分支和提交
- ✅ 超时保护，避免长时间等待

## 配置

### 1. Git 代理配置（已配置）

```bash
git config --global http.https://github.com.proxy http://10.21.71.52:7890
git config --global https.https://github.com.proxy http://10.21.71.52:7890
```

### 2. 远程仓库

- **origin** (Gitee): 主仓库
- **github** (GitHub): 镜像仓库

## 使用方法

### 日常推送（推荐）

```bash
# 1. 推送到 Gitee
git push

# 2. 同步当前分支到 GitHub
./scripts/git-sync/sync_repos.sh
```

### 推送所有分支

```bash
./scripts/git-sync/sync_repos.sh --all
```

### 推送指定分支

```bash
./scripts/git-sync/sync_repos.sh --branch master
```

### 增量同步（推荐）

自动检测并同步所有更新到 GitHub：
- 推送 GitHub 缺失的分支
- 更新所有分支的新提交（包括 master）

```bash
./scripts/git-sync/sync_repos.sh --sync
```

### 查看帮助

```bash
./scripts/git-sync/sync_repos.sh --help
```

## 常见场景

### 场景 1：开发新功能

```bash
git checkout -b feat-new-feature
git add .
git commit -m "feat: 新功能"
git push                                    # 推送到 Gitee
./scripts/git-sync/sync_repos.sh           # 同步到 GitHub
```

### 场景 2：合并分支后同步

```bash
# 在 Gitee 上合并 PR 后，拉取最新代码
git pull

# 增量同步（会自动同步 master 的新提交）
./scripts/git-sync/sync_repos.sh --sync
```

### 场景 3：定期全量同步

```bash
# 定期执行，保持两个仓库完全同步
./scripts/git-sync/sync_repos.sh --sync
```

## 注意事项

- Windows 代理关闭时，脚本会自动跳过 GitHub 推送
- GitHub 推送失败不会影响 Gitee
- `--sync` 模式会自动同步所有分支的更新（包括 master）
- 推荐使用 `--sync` 模式保持两个仓库完全同步
