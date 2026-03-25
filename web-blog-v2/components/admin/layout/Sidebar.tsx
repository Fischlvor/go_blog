'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import {
  LayoutDashboard, FileText, MessageSquare, Users, Image,
  HardDrive, Smile, Settings, Bot, LogOut, ChevronLeft, ChevronRight,
  Link2, Megaphone, ClipboardList, ScrollText
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import { logout } from '@/lib/api/user/user';
import { useRouter } from 'next/navigation';

interface NavItem {
  href: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  exact?: boolean;
}

interface NavGroup {
  label: string;
  items: NavItem[];
}

const NAV_GROUPS: NavGroup[] = [
  {
    label: '概览',
    items: [
      { href: '/admin', label: '仪表盘', icon: LayoutDashboard, exact: true },
    ],
  },
  {
    label: '内容管理',
    items: [
      { href: '/admin/articles', label: '文章列表', icon: FileText },
      { href: '/admin/articles/new', label: '发布文章', icon: FileText },
      { href: '/admin/comments', label: '评论管理', icon: MessageSquare },
    ],
  },
  {
    label: '媒体管理',
    items: [
      { href: '/admin/images', label: '图片管理', icon: Image },
      { href: '/admin/resources', label: '资源管理', icon: HardDrive },
      { href: '/admin/emoji', label: '表情管理', icon: Smile },
    ],
  },
  {
    label: 'AI 管理',
    items: [
      { href: '/admin/ai/models', label: '模型管理', icon: Bot },
      { href: '/admin/ai/sessions', label: '会话管理', icon: Bot },
    ],
  },
  {
    label: '系统',
    items: [
      { href: '/admin/users', label: '用户管理', icon: Users },
      { href: '/admin/system/links', label: '友链管理', icon: Link2 },
      { href: '/admin/system/ads', label: '广告管理', icon: Megaphone },
      { href: '/admin/system/feedback', label: '反馈管理', icon: ClipboardList },
      { href: '/admin/system/login-logs', label: '登录日志', icon: ScrollText },
      { href: '/admin/system/config', label: '系统配置', icon: Settings },
    ],
  },
];

export function AdminSidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem('admin-sidebar-collapsed');
    if (saved) setCollapsed(saved === 'true');
  }, []);

  const toggleCollapse = () => {
    const next = !collapsed;
    setCollapsed(next);
    localStorage.setItem('admin-sidebar-collapsed', String(next));
  };

  const handleLogout = async () => {
    try {
      await logout();
      localStorage.removeItem('access_token');
      toast.success('已退出登录');
      router.push('/');
    } catch {
      localStorage.removeItem('access_token');
      router.push('/');
    }
  };

  const isActive = (href: string, exact?: boolean): boolean =>
    exact ? pathname === href : pathname.startsWith(href) && (href !== '/admin' || pathname === '/admin');

  return (
    <aside className={cn(
      'flex flex-col h-screen border-r border-border bg-card transition-all duration-300 flex-shrink-0',
      collapsed ? 'w-16' : 'w-60'
    )}>
      {/* Logo */}
      <div className={cn('flex items-center h-16 px-4 border-b border-border gap-3', collapsed && 'justify-center px-0')}  >
        <Link href="/" className="flex items-center gap-2 font-bold text-base">
          <span className="w-7 h-7 rounded-lg bg-gradient-to-br from-violet-600 to-cyan-500 flex-shrink-0" />
          {!collapsed && <span className="bg-gradient-to-r from-violet-600 to-cyan-500 bg-clip-text text-transparent truncate">管理后台</span>}
        </Link>
      </div>

      {/* Nav */}
      <ScrollArea className="flex-1 py-3">
        <nav className="px-2 space-y-4">
          {NAV_GROUPS.map((group) => (
            <div key={group.label}>
              {!collapsed && (
                <p className="px-3 mb-1 text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">
                  {group.label}
                </p>
              )}
              <ul className="space-y-0.5">
                {group.items.map(item => {
                  const Icon = item.icon as React.ComponentType<{ className?: string }>;
                  const active = isActive(item.href, item.exact);
                  return (
                    <li key={item.href}>
                      <Link
                        href={item.href}
                        title={collapsed ? item.label : undefined}
                        className={cn(
                          'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                          active ? 'bg-primary/10 text-primary' : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                          collapsed && 'justify-center px-0'
                        )}
                      >
                        <Icon className="h-4 w-4 flex-shrink-0" />
                        {!collapsed && <span className="truncate">{item.label}</span>}
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          ))}
        </nav>
      </ScrollArea>

      {/* Bottom */}
      <div className="border-t border-border p-2 space-y-1">
        <Button
          variant="ghost" size="sm"
          className={cn('w-full text-muted-foreground', collapsed ? 'justify-center px-0' : 'justify-start gap-2')}
          onClick={toggleCollapse}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <><ChevronLeft className="h-4 w-4" /> <span>收起菜单</span></>}
        </Button>
        <Button
          variant="ghost" size="sm"
          className={cn('w-full text-destructive hover:text-destructive', collapsed ? 'justify-center px-0' : 'justify-start gap-2')}
          onClick={handleLogout}
        >
          <LogOut className="h-4 w-4" />
          {!collapsed && <span>退出登录</span>}
        </Button>
      </div>
    </aside>
  );
}
