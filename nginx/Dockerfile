FROM nginx:alpine

# 复制nginx配置文件
COPY go_blog.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80 443

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]

