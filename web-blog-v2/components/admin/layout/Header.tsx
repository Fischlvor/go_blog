'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Sun, Moon } from 'lucide-react';
import { useTheme } from 'next-themes';
import { Button } from '@/components/ui/button';
import { useState, useEffect } from 'react';
import { UserMenu } from '@/components/site/auth/UserMenu';

function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="w-9 h-9" />;
  const isDark = theme === 'dark';

  return (
    <Button variant="ghost" size="icon" onClick={() => setTheme(isDark ? 'light' : 'dark')} aria-label="切换主题">
      {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
    </Button>
  );
}

function getBreadcrumb(pathname: string): string {
  const map: Record<string, string> = {
    '/admin': '主页',
    '/admin/user-center/user-info': '我的信息',
    '/admin/user-center/user-star': '我的收藏',
    '/admin/user-center/user-comment': '我的评论',
    '/admin/user-center/user-feedback': '我的反馈',
    '/admin/users/user-list': '用户列表',
    '/admin/articles/article-publish': '发布文章',
    '/admin/articles/comment-list': '评论列表',
    '/admin/articles/article-list': '文章列表',
    '/admin/images/image-list': '图片列表',
    '/admin/resources/resource-list': '资源列表',
    '/admin/ai-management/models': '模型管理',
    '/admin/ai-management/sessions': '会话管理',
    '/admin/ai-management/messages': '消息管理',
    '/admin/emoji/emoji-list': '表情列表',
    '/admin/emoji/emoji-groups': '表情组管理',
    '/admin/emoji/emoji-sprites': '雪碧图管理',
    '/admin/system/friend-link-list': '友链列表',
    '/admin/system/advertisement-list': '广告列表',
    '/admin/system/feedback-list': '反馈列表',
    '/admin/system/login-logs': '登录日志',
    '/admin/system/app-config': '应用配置',
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
      <div className="flex items-center gap-1">
        <ThemeToggle />
        <UserMenu />
      </div>
    </header>
  );
}
