# 第一阶段：构建阶段
FROM node:22-alpine AS builder
# 写入前端构建所需环境变量（可被 --build-arg 覆盖）
ARG VITE_CDN_BASE=https://image.hsk423.cn
ENV VITE_CDN_BASE=$VITE_CDN_BASE

# 设置工作目录
WORKDIR /app

# 设置 npm 国内镜像源
RUN npm config set registry https://registry.npmmirror.com

# 复制package文件
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 第二阶段：运行阶段
FROM nginx:alpine

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置（如果需要自定义配置）
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]

