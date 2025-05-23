# go_blog

gin + vue3 开发的个人博客项目

## 本地测试

### 开发工具及版本

golang: 1.23.0

node: v22.13.0

docker: 26.1.3

编译器：IDEA

### 启动容器

```bash
docker run -itd --name mysql -p 13306:3306 -e  MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=blog_db -d mysql:8.0.20

docker run --name es -p 127.0.0.1:9200:9200 -e "discovery.type=single-node" -e "xpack.security.http.ssl.enabled=false" -e "xpack.license.self_generated.type=trial" -e "xpack.security.enabled=false" -e ES_JAVA_OPTS="-Xms84m -Xmx512m" -d elasticsearch:8.17.0

docker run --name redis -p 16379:6379 -d redis:6.2
```

### 启动服务

```bash
# 进入 server 文件夹，修改配置文件 config.yaml
go mod tidy
# 初始化 mysql
go run main.go -sql
# 初始化 es
go run main.go -es
# 创建管理员
go run main.go -admin
# 运行后端
go run main.go

# 进入 web 文件夹
npm install
# 运行前端
npm run dev
```

## 项目部署

### 部署环境

Alibaba Cloud Linux release 3 (OpenAnolis Edition)(CentOS 8)

### 环境准备

```bash
# 安装 docker
yum install -y docker-ce
systemctl start docker
systemctl enable docker

# 安装 supervisor
yum install -y supervisor
systemctl start supervisord
systemctl enable supervisord

# 安装 nginx
阿里云，去宝塔面板安装 nginx:1.26.×
```

### 准备工作

编译后端，得到 main 文件

```bash
# windows环境下，打开项目所在目录，进入 server 文件夹，打开 cmd
set GOOS=linux
set GOARCH=amd64
go mod tidy
go build main.go
```

编译前端，得到 dist 文件夹

```bash
# windows环境下，打开项目所在目录，进入 web 文件夹，打开 cmd
npm install
# 将 index.html 中 http://127.0.0.1:8080/api/website/logo 替换为域名 https://www.your_domain/api/website/logo
npm run build
```

### 服务端目录

将文件按照下述目录上传

```bash
# /opt/go_blog
├── go_blog
    ├── data
    │   ├── es
    │   └── mysql
    ├── server
    │   ├── main
    │   └── config.yaml
    └── web
        └── dist
```

### 容器配置

#### redis

```
docker run --name redis6 --restart=always  -p 16379:6379 -d redis:6.2
```

#### mysql

```
docker run -itd --name mysql8 --restart=always -p 13306:3306 -v /media/practice/onServer/go_blog/data/mysql/conf:/etc/mysql/conf.d -v /media/practice/onServer/go_blog/data/mysql/datadir:/var/lib/mysql -v /media/practice/onServer/go_blog/data/mysql/go_blog.sql:/opt/go_blog.sql -e  MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=blog_db mysql:8.0.20
```

#### elasticsearch 

elasticsearch 无法直接数据卷挂载本地，需要先启动一个不挂载数据卷的容器，将文件复制到本地，再进行挂载

```bash
# 复制文件
docker run --name es8 --restart=always -p 127.0.0.1:9200:9200 -e "discovery.type=single-node" -e "xpack.security.http.ssl.enabled=false" -e "xpack.license.self_generated.type=trial" -e "xpack.security.enabled=false" -e ES_JAVA_OPTS="-Xms84m -Xmx512m"  -d elasticsearch:8.17.0

docker cp es8:/usr/share/elasticsearch/config /media/practice/onServer/go_blog/data/es/config
docker cp es8:/usr/share/elasticsearch/data /media/practice/onServer/go_blog/data/es/data 
docker cp es8:/usr/share/elasticsearch/plugins /media/practice/onServer/go_blog/data/es/plugins

docker rm -f es8
```

```
docker run --name es8 --restart=always -p 127.0.0.1:9200:9200 -e "discovery.type=single-node" -e "xpack.security.http.ssl.enabled=false" -e "xpack.license.self_generated.type=trial" -e "xpack.security.enabled=false" -e ES_JAVA_OPTS="-Xms84m -Xmx512m" -v /media/practice/onServer/go_blog/data/es/config:/usr/share/elasticsearch/config -v /media/practice/onServer/go_blog/data/es/data:/usr/share/elasticsearch/data -v /media/practice/onServer/go_blog/data/es/plugins:/usr/share/elasticsearch/plugins  -d elasticsearch:8.17.0
```

### nginx配置

修改 /www/server/nginx/conf/nginx.conf

```nginx

user  www www;
worker_processes auto;
error_log  /www/wwwlogs/nginx_error.log  crit;  # 改这里
pid        /www/server/nginx/logs/nginx.pid;    # 改这里
worker_rlimit_nofile 51200;

stream {
    log_format tcp_format '$time_local|$remote_addr|$protocol|$status|$bytes_sent|$bytes_received|$session_time|$upstream_addr|$upstream_bytes_sent|$upstream_bytes_received|$upstream_connect_time';
  
    access_log /www/wwwlogs/tcp-access.log tcp_format;   # 改这里
    error_log /www/wwwlogs/tcp-error.log;                # 改这里
    include /www/server/panel/vhost/nginx/tcp/*.conf;    # 改这里，配置文件放这里
}

events
    {
        use epoll;
        worker_connections 51200;
        multi_accept on;
    }

http
    {
        include       mime.types;
		#include luawaf.conf;

		include proxy.conf;
        lua_package_path "/www/server/nginx/lib/lua/?.lua;;";

        default_type  application/octet-stream;

        server_names_hash_bucket_size 512;
        client_header_buffer_size 32k;
        large_client_header_buffers 4 32k;
        client_max_body_size 50m;

        sendfile   on;
        tcp_nopush on;

        keepalive_timeout 60;

        tcp_nodelay on;

        fastcgi_connect_timeout 300;
        fastcgi_send_timeout 300;
        fastcgi_read_timeout 300;
        fastcgi_buffer_size 64k;
        fastcgi_buffers 4 64k;
        fastcgi_busy_buffers_size 128k;
        fastcgi_temp_file_write_size 256k;
		fastcgi_intercept_errors on;

        gzip on;
        gzip_min_length  1k;
        gzip_buffers     4 16k;
        gzip_http_version 1.1;
        gzip_comp_level 2;
        gzip_types     text/plain application/javascript application/x-javascript text/javascript text/css application/xml application/json image/jpeg image/gif image/png font/ttf font/otf image/svg+xml application/xml+rss text/x-js;
        gzip_vary on;
        gzip_proxied   expired no-cache no-store private auth;
        gzip_disable   "MSIE [1-6]\.";

        limit_conn_zone $binary_remote_addr zone=perip:10m;
		limit_conn_zone $server_name zone=perserver:10m;

        server_tokens off;
        access_log off;

server
    {
        listen 888;
        server_name phpmyadmin;
        index index.html index.htm index.php;
        root  /www/server/phpmyadmin;

        #error_page   404   /404.html;
        include enable-php.conf;

        location ~ .*\.(gif|jpg|jpeg|png|bmp|swf)$
        {
            expires      30d;
        }

        location ~ .*\.(js|css)?$
        {
            expires      12h;
        }

        location ~ /\.
        {
            deny all;
        }

        access_log  /www/wwwlogs/access.log;
    }
include /www/server/panel/vhost/nginx/*.conf;  # 改这里，配置文件放这里
}



```

创建 /www/server/panel/vhost/nginx/go_blog.conf

#### 没有域名
```nginx
server {
    gzip on;
    gzip_vary on;
    gzip_disable "MSIE [1-6]\.";
    gzip_static on;
    gzip_min_length 256;
    gzip_buffers 32 8k;
    gzip_http_version 1.1;
    gzip_comp_level 5;
    gzip_proxied any;
    gzip_types text/plain text/css text/xml application/javascript application/x-javascript application/xml application/xml+rss application/emacscript application/json image/svg+xml;

    listen 80 default_server;
    #server_name _; # 匹配所有域名

    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";

    location / {
        try_files $uri $uri/ /index.html;
        root   /media/practice/onServer/go_blog/web/dist;
        index  index.html index.htm;
    }

    location /api/ {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header REMOTE-HOST $remote_addr;
        proxy_pass http://127.0.0.1:8081/api/;
    }

    location /image {
        alias /media/practice/onServer/go_blog/web/dist/image;
    }

    location /emoji {
        alias /media/practice/onServer/go_blog/web/dist/emoji;
    }

    location /uploads/ {
        alias /media/practice/onServer/go_blog/server/uploads/;
    }
}
    
```
#### 有域名
将 your_domain 替换为你的域名，请自行获取 ssl 证书，上传证书文件至 /etc/nginx/cert/

```nginx
server {
    listen 80;
    server_name your_domain www.your_domain;
    return 301 https://www.your_domain$request_uri;
}

server { 
    listen 443 ssl; 
    server_name your_domain;  # 仅匹配非 www 的域名
    ssl_certificate /etc/nginx/cert/your_domain.crt; # 证书公钥
    ssl_certificate_key /etc/nginx/cert/your_domain.key; # 证书私钥
    return 301 https://www.your_domain$request_uri;  # 强制跳转到 www.your_domain
}

server {
    gzip on;
    gzip_vary on;
    gzip_disable "MSIE [1-6]\.";
    gzip_static on;
    gzip_min_length 256;
    gzip_buffers 32 8k;
    gzip_http_version 1.1;
    gzip_comp_level 5;
    gzip_proxied any;
    gzip_types text/plain text/css text/xml application/javascript application/x-javascript application/xml application/xml+rss application/emacscript application/json image/svg+xml;

    listen 443 ssl;
    server_name www.your_domain; # 多个域名⽤空格分开 
    ssl_certificate /etc/nginx/cert/your_domain.crt; # 证书公钥
    ssl_certificate_key /etc/nginx/cert/your_domain.key; # 证书私钥
    ssl_session_timeout 5m; 
    ssl_session_cache shared:MozSSL:10m;  # 设置会话缓存以提⾼性能 
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;  # 配置加密算法 
    ssl_protocols TLSv1.2 TLSv1.3;  # 配置加密协议 
    ssl_prefer_server_ciphers on;

    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always; #可选配置，开启HSTS 
    add_header X-Frame-Options DENY; # 可选配置，防⽌点击劫持 
    add_header X-Content-Type-Options nosniff; # 可选配置，防⽌MIME类型嗅探 
    add_header X-XSS-Protection "1; mode=block"; # 可选配置，防⽌XSS攻击

    location / {
        try_files $uri $uri/ /index.html;
        root   /opt/go_blog/web/dist;
        index  index.html index.htm;
    }

    location /api/ {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header REMOTE-HOST $remote_addr;
        proxy_pass http://127.0.0.1:8080/api/;
    }

    location /image {
        alias /opt/go_blog/web/dist/image;
    }

    location /emoji {
        alias /opt/go_blog/web/dist/emoji;
    }

    location /uploads/ {
        alias /opt/go_blog/server/uploads/;
    }
}
```

重新启动，宝塔面板直接**重载配置**

### supervisor配置

给 main 执行权限，并初始化项目

```bash
# 进入 /media/practice/onServer/go_blog/server
chmod +x ./main

./main -sql
./main -es
./main -admin
```

创建 /etc/supervisord.d/go_blog.ini

```ini
[program: go_blog]
command=/media/practice/onServer/go_blog/server/main
directory=/media/practice/onServer/go_blog/server/
autorestart=true ; 程序意外退出是否自动重启
autostart=true ; 是否自动启动
user=root ; 进程执行的用户身份
stopsignal=INT
startsecs=1 ; 自动重启间隔
stopasgroup=true ;默认为false,进程被杀死时，是否向这个进程组发送stop信号，包括子进程
killasgroup=true ;默认为false，向进程组发送kill信号，包括子进程
```

重新启动

```bash
systemctl restart supervisord
```
