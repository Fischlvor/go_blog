#!/bin/bash

# RSA密钥生成脚本

KEYS_DIR="../keys"

# 创建keys目录
mkdir -p "$KEYS_DIR"

echo "正在生成RSA密钥对..."

# 生成私钥（2048位）
openssl genrsa -out "$KEYS_DIR/private.pem" 2048

# 从私钥生成公钥
openssl rsa -in "$KEYS_DIR/private.pem" -pubout -out "$KEYS_DIR/public.pem"

# 设置权限（私钥只有owner可读写）
chmod 600 "$KEYS_DIR/private.pem"
chmod 644 "$KEYS_DIR/public.pem"

echo "✓ RSA密钥对生成成功！"
echo "  私钥: $KEYS_DIR/private.pem"
echo "  公钥: $KEYS_DIR/public.pem"
echo ""
echo "⚠️  请妥善保管私钥文件！"
echo "📋  公钥需要复制到博客服务(server/keys/public.pem)"

