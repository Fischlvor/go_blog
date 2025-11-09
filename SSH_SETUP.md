# SSH 公钥配置说明

## 本机公钥信息

您的本机已有 SSH 公钥：
- 类型：ED25519
- 文件路径：`~/.ssh/id_ed25519.pub`

## 将公钥添加到远程服务器

### 方法一：使用 ssh-copy-id（推荐，最简单）

```bash
ssh-copy-id -p 22 root@8.148.64.96
```

如果命令不存在，可以使用以下命令：

```bash
cat ~/.ssh/id_ed25519.pub | ssh -p 22 root@8.148.64.96 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && chmod 700 ~/.ssh"
```

### 方法二：手动复制（如果方法一失败）

1. **复制公钥内容**：
   ```bash
   cat ~/.ssh/id_ed25519.pub
   ```
   复制输出的内容

2. **登录到远程服务器**：
   ```bash
   ssh -p 22 root@8.148.64.96
   ```

3. **在远程服务器上执行**：
   ```bash
   # 创建.ssh目录（如果不存在）
   mkdir -p ~/.ssh
   
   # 设置正确的权限
   chmod 700 ~/.ssh
   
   # 将公钥添加到authorized_keys文件
   echo "你的公钥内容" >> ~/.ssh/authorized_keys
   
   # 设置authorized_keys文件的权限
   chmod 600 ~/.ssh/authorized_keys
   ```

4. **退出远程服务器**：
   ```bash
   exit
   ```

### 方法三：使用脚本自动配置

运行以下命令（会自动处理）：

```bash
PUB_KEY=$(cat ~/.ssh/id_ed25519.pub)
ssh -p 22 root@8.148.64.96 "mkdir -p ~/.ssh && echo '$PUB_KEY' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && chmod 700 ~/.ssh && echo '公钥已添加成功'"
```

## 验证配置

配置完成后，测试免密登录：

```bash
ssh -p 22 root@8.148.64.96
```

如果不需要输入密码就能登录，说明配置成功！

## 注意事项

1. **权限很重要**：
   - `~/.ssh` 目录权限必须是 `700` (drwx------)
   - `~/.ssh/authorized_keys` 文件权限必须是 `600` (-rw-------)

2. **如果配置失败**，检查远程服务器的SSH配置：
   ```bash
   # 在远程服务器上检查
   cat /etc/ssh/sshd_config | grep -E "PubkeyAuthentication|AuthorizedKeysFile"
   ```
   确保：
   - `PubkeyAuthentication yes`
   - `AuthorizedKeysFile .ssh/authorized_keys`

3. **如果仍然需要密码**，可能需要检查：
   - SELinux 是否阻止（如果启用）
   - 防火墙设置
   - SSH 服务配置

## 公钥内容

您的公钥内容：
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEwAy9VIO2bxzM86vAKfh5vGVcska9MFdbpua3ZfzDFQ gitee_key
```

