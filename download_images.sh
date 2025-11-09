#!/bin/bash

# 在远程服务器上执行此脚本下载所有基础镜像
# 使用方法：ssh root@8.148.64.96 'bash -s' < download_images.sh

echo "开始下载基础镜像..."

# ==================== 构建镜像（Dockerfile中使用）====================
echo "1. 下载 Golang 构建镜像..."
docker pull golang:1.23.0-alpine

echo "2. 下载 Alpine 运行镜像..."
docker pull alpine:latest

echo "3. 下载 Node.js 构建镜像..."
docker pull node:20-alpine

echo "4. 下载 Nginx 镜像..."
docker pull nginx:alpine

echo ""
echo "=========================================="
echo "所有基础镜像下载完成！"
echo "=========================================="
echo ""
echo "查看已下载的镜像："
docker images

