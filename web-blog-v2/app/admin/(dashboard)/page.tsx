'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { FileText, MessageSquare, Users, Eye } from 'lucide-react';
import { listArticles } from '@/lib/client-api/public/article';
import { adminListUsers } from '@/lib/client-api/admin/user';

export default function AdminDashboardPage() {
  const [stats, setStats] = useState({ articles: 0, users: 0 });

  useEffect(() => {
    Promise.all([
      listArticles({ page: 1, page_size: 1 }).catch(() => ({ total_items: 0, list: [] as never[], current_page: 1, page_size: 1, total_pages: 0 })),
      adminListUsers({ page: 1, page_size: 1 }).catch(() => ({ total_items: 0, list: [] as never[], current_page: 1, page_size: 1, total_pages: 0 })),
    ]).then(([arts, users]) => {
      setStats({ articles: arts.total_items ?? 0, users: users.total_items ?? 0 });
    });
  }, []);

  const cards = [
    { label: '文章总数', value: stats.articles, icon: FileText, color: 'text-blue-500' },
    { label: '注册用户', value: stats.users, icon: Users, color: 'text-green-500' },
    { label: '今日访问', value: '—', icon: Eye, color: 'text-orange-500' },
    { label: '待审评论', value: '—', icon: MessageSquare, color: 'text-purple-500' },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">仪表盘</h1>
        <p className="text-muted-foreground">欢迎回来</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map(card => {
          const Icon = card.icon;
          return (
            <Card key={card.label}>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">{card.label}</CardTitle>
                <Icon className={`h-4 w-4 ${card.color}`} />
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{card.value}</p>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
