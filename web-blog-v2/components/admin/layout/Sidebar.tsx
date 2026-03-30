'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import {
  LayoutDashboard, FileText, MessageSquare, Users, Image,
  HardDrive, Smile, Settings, Bot, ChevronLeft, ChevronRight,
  Link2, Megaphone, ClipboardList, ScrollText, ChevronDown, ChevronUp
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { cn } from '@/lib/utils';

interface NavItem {
  href: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  exact?: boolean;
}

interface NavGroup {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  items: NavItem[];
}

const NAV_GROUPS: NavGroup[] = [
  {
    label: '概览',
    icon: LayoutDashboard,
    items: [
      { href: '/admin', label: '主页', icon: LayoutDashboard, exact: true },
    ],
  },
  {
    label: '个人中心',
    icon: Users,
    items: [
      { href: '/admin/user-center/user-info', label: '我的信息', icon: Users },
      { href: '/admin/user-center/user-star', label: '我的收藏', icon: FileText },
      { href: '/admin/user-center/user-comment', label: '我的评论', icon: MessageSquare },
      { href: '/admin/user-center/user-feedback', label: '我的反馈', icon: ClipboardList },
    ],
  },
  {
    label: '用户管理',
    icon: Users,
    items: [
      { href: '/admin/users/user-list', label: '用户列表', icon: Users },
    ],
  },
  {
    label: '文章管理',
    icon: ScrollText,
    items: [
      { href: '/admin/articles/article-publish', label: '发布文章', icon: FileText },
      { href: '/admin/articles/comment-list', label: '评论列表', icon: MessageSquare },
      { href: '/admin/articles/article-list', label: '文章列表', icon: ScrollText },
    ],
  },
  {
    label: '图片与资源',
    icon: Image,
    items: [
      { href: '/admin/images/image-list', label: '图片列表', icon: Image },
      { href: '/admin/resources/resource-list', label: '资源列表', icon: HardDrive },
    ],
  },
  {
    label: 'AI 对话管理',
    icon: Bot,
    items: [
      { href: '/admin/ai-management/models', label: '模型管理', icon: Bot },
      { href: '/admin/ai-management/sessions', label: '会话管理', icon: Bot },
      { href: '/admin/ai-management/messages', label: '消息管理', icon: MessageSquare },
    ],
  },
  {
    label: '表情管理',
    icon: Smile,
    items: [
      { href: '/admin/emoji/emoji-list', label: '表情列表', icon: Smile },
      { href: '/admin/emoji/emoji-groups', label: '表情组管理', icon: Link2 },
      { href: '/admin/emoji/emoji-sprites', label: '雪碧图管理', icon: Megaphone },
    ],
  },
  {
    label: '系统管理',
    icon: Settings,
    items: [
      { href: '/admin/system/feedback-list', label: '反馈列表', icon: ClipboardList },
      { href: '/admin/system/advertisement-list', label: '广告列表', icon: Megaphone },
      { href: '/admin/system/friend-link-list', label: '友链列表', icon: Link2 },
      { href: '/admin/system/login-logs', label: '登录日志', icon: ScrollText },
      { href: '/admin/system/app-config', label: '应用配置', icon: Settings },
    ],
  },
];

export function AdminSidebar() {
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);
  const [expandedGroups, setExpandedGroups] = useState<Record<string, boolean>>({});

  useEffect(() => {
    const saved = localStorage.getItem('admin-sidebar-collapsed');
    if (saved) setCollapsed(saved === 'true');

    const storedGroups = localStorage.getItem('admin-sidebar-expanded-groups');
    if (storedGroups) {
      try {
        setExpandedGroups(JSON.parse(storedGroups) as Record<string, boolean>);
        return;
      } catch {
        // ignore broken local storage
      }
    }

    const defaults = Object.fromEntries(NAV_GROUPS.map((group) => [group.label, true]));
    setExpandedGroups(defaults);
  }, []);

  const toggleCollapse = () => {
    const next = !collapsed;
    setCollapsed(next);
    localStorage.setItem('admin-sidebar-collapsed', String(next));
  };

  const toggleGroup = (label: string) => {
    setExpandedGroups((prev) => {
      const next = { ...prev, [label]: !prev[label] };
      localStorage.setItem('admin-sidebar-expanded-groups', JSON.stringify(next));
      return next;
    });
  };

  const isActive = (href: string, exact?: boolean): boolean =>
    exact ? pathname === href : pathname.startsWith(href) && (href !== '/admin' || pathname === '/admin');

  const isGroupActive = (group: NavGroup): boolean => group.items.some((item) => isActive(item.href, item.exact));

  return (
    <aside className={cn(
      'flex flex-col h-screen border-r border-border bg-sidebar text-sidebar-foreground transition-all duration-300 flex-shrink-0',
      collapsed ? 'w-16' : 'w-60'
    )}>
      <div className={cn('flex items-center h-16 px-4 border-b border-border gap-3', collapsed && 'px-2 justify-center')}>
        <Link href="/" className={cn('flex items-center gap-2 font-bold text-base min-w-0', collapsed && 'hidden')}>
          <span className="w-7 h-7 rounded-lg bg-gradient-to-br from-violet-600 to-cyan-500 flex-shrink-0" />
          <span className="bg-gradient-to-r from-violet-600 to-cyan-500 bg-clip-text text-transparent truncate">管理后台</span>
        </Link>

        <Button
          variant="ghost"
          size="icon"
          className={cn('ml-auto text-muted-foreground hover:text-foreground', collapsed && 'ml-0')}
          onClick={toggleCollapse}
          title={collapsed ? '展开菜单' : '收起菜单'}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </Button>
      </div>

      <ScrollArea className="flex-1 min-h-0">
        <nav className="px-2 py-3 space-y-3">
          {collapsed ? (
            <ul className="space-y-1">
              {NAV_GROUPS.map((group) => {
                const GroupIcon = group.icon;
                const active = isGroupActive(group);
                const target = group.items[0]?.href || '/admin';
                return (
                  <li key={group.label}>
                    <Link
                      href={target}
                      title={group.label}
                      className={cn(
                        'group flex items-center justify-center rounded-lg py-2.5 transition-all outline-none focus-visible:ring-2 focus-visible:ring-primary/40',
                        active
                          ? 'bg-primary/14 text-primary ring-1 ring-primary/30'
                          : 'text-muted-foreground hover:bg-accent/75 hover:text-foreground'
                      )}
                    >
                      <GroupIcon className="h-4.5 w-4.5" />
                    </Link>
                  </li>
                );
              })}
            </ul>
          ) : (
            NAV_GROUPS.map((group) => {
              const groupExpanded = expandedGroups[group.label] ?? true;
              return (
                <div key={group.label} className="rounded-xl border border-border/70 bg-background/70 dark:bg-muted/15 backdrop-blur-[1px]">
                  <button
                    type="button"
                    onClick={() => toggleGroup(group.label)}
                    className="w-full px-3 py-2.5 flex items-center justify-between text-xs font-semibold tracking-[0.08em] text-muted-foreground/90 hover:text-foreground transition-colors border-b border-border/60"
                  >
                    <span className="inline-flex items-center gap-2"><span className="h-1.5 w-1.5 rounded-full bg-primary/65" />{group.label}</span>
                    {groupExpanded ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
                  </button>

                  {groupExpanded && (
                    <ul className="space-y-1 p-2 ml-2 my-1 pl-2 border-l border-border/70">
                      {group.items.map(item => {
                        const Icon = item.icon as React.ComponentType<{ className?: string }>;
                        const active = isActive(item.href, item.exact);
                        return (
                          <li key={item.href}>
                            <Link
                              href={item.href}
                              className={cn(
                                'group flex items-center gap-3 rounded-lg px-3 py-2.5 text-[13px] font-medium transition-all outline-none focus-visible:ring-2 focus-visible:ring-primary/40 focus-visible:ring-offset-1 focus-visible:ring-offset-background',
                                active
                                  ? 'bg-primary/14 text-primary ring-1 ring-primary/30 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.06)] dark:shadow-none'
                                  : 'text-muted-foreground hover:bg-accent/75 hover:text-foreground'
                              )}
                            >
                              <Icon className={cn('h-4 w-4 flex-shrink-0 transition-transform', !active && 'group-hover:scale-105')} />
                              <span className="truncate">{item.label}</span>
                            </Link>
                          </li>
                        );
                      })}
                    </ul>
                  )}
                </div>
              );
            })
          )}
        </nav>
      </ScrollArea>
    </aside>
  );
}
