'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Sun, Moon, Monitor } from 'lucide-react';
import { useTheme } from 'next-themes';
import { Button } from '@/components/ui/button';
import { useState, useEffect } from 'react';

function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="w-9 h-9" />;
  const next = theme === 'light' ? 'dark' : theme === 'dark' ? 'system' : 'light';
  const Icon = theme === 'light' ? Sun : theme === 'dark' ? Moon : Monitor;
  return (
    <Button variant="ghost" size="icon" onClick={() => setTheme(next)}>
      <Icon className="h-4 w-4" />
    </Button>
  );
}

function getBreadcrumb(pathname: string): string {
  const map: Record<string, string> = {
    '/admin': '仪表盘',
    '/admin/articles': '文章列表',
    '/admin/articles/new': '发布文章',
    '/admin/comments': '评论管理',
    '/admin/users': '用户管理',
    '/admin/images': '图片管理',
    '/admin/resources': '资源管理',
    '/admin/emoji': '表情管理',
    '/admin/ai/models': '模型管理',
    '/admin/ai/sessions': '会话管理',
    '/admin/system/links': '友链管理',
    '/admin/system/ads': '广告管理',
    '/admin/system/feedback': '反馈管理',
    '/admin/system/login-logs': '登录日志',
    '/admin/system/config': '系统配置',
  };
  return map[pathname] || '管理后台';
}

export function AdminHeader() {
  const pathname = usePathname();

  return (
    <header className="h-16 border-b border-border bg-card flex items-center justify-between px-6 flex-shrink-0">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link href="/" className="hover:text-foreground transition-colors">
          <Home className="h-4 w-4" />
        </Link>
        <span>/</span>
        <span className="text-foreground font-medium">{getBreadcrumb(pathname)}</span>
      </div>
      <div className="flex items-center gap-2">
        <ThemeToggle />
      </div>
    </header>
  );
}
