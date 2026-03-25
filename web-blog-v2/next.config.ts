import type { NextConfig } from 'next';

// 开发时允许的跨域来源（局域网访问用），生产环境忽略
const devOrigins = process.env.NEXT_DEV_ORIGINS
  ? process.env.NEXT_DEV_ORIGINS.split(',')
  : [];

const nextConfig: NextConfig = {
  ...(devOrigins.length > 0 ? { allowedDevOrigins: devOrigins } : {}),
  images: {
    remotePatterns: [
      { protocol: 'https', hostname: '**' },
      { protocol: 'http', hostname: '**' },
    ],
  },
  async rewrites() {
    const apiBase = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:8081';
    return [
      {
        source: '/api/v1/:path*',
        destination: `${apiBase}/api/v1/:path*`,
      },
      {
        source: '/api/v2/:path*',
        destination: `${apiBase}/api/v2/:path*`,
      },
    ];
  },
};

export default nextConfig;
