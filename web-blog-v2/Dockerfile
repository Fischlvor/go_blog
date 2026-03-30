# 运行阶段 - 使用本地预构建的产物
FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0

RUN npm config set registry https://registry.npmmirror.com && \
    npm install -g pnpm

COPY package*.json pnpm-lock.yaml ./

# 只安装生产依赖，然后删除缓存
RUN pnpm install --prod --frozen-lockfile && \
    pnpm store prune

# 复制预构建的产物
COPY .next/standalone ./
COPY .next/static ./.next/static
COPY public ./public

EXPOSE 3000

CMD ["node", "server.js"]
