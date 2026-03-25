'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Calendar } from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import { listArticles } from '@/lib/api/public/article';
import type { Article } from '@/lib/api/types';

interface GroupedArticles {
  year: number;
  months: {
    month: number;
    articles: Article[];
  }[];
}

function groupByYearMonth(articles: Article[]): GroupedArticles[] {
  const map = new Map<number, Map<number, Article[]>>();
  for (const a of articles) {
    const d = new Date(a.published_at || a.created_at);
    const y = d.getFullYear();
    const m = d.getMonth() + 1;
    if (!map.has(y)) map.set(y, new Map());
    const yMap = map.get(y)!;
    if (!yMap.has(m)) yMap.set(m, []);
    yMap.get(m)!.push(a);
  }
  return Array.from(map.entries())
    .sort(([a], [b]) => b - a)
    .map(([year, months]) => ({
      year,
      months: Array.from(months.entries())
        .sort(([a], [b]) => b - a)
        .map(([month, articles]) => ({ month, articles })),
    }));
}

export default function ArchivePage() {
  const [groups, setGroups] = useState<GroupedArticles[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listArticles({ page: 1, page_size: 999, order: 'desc' })
      .then(res => {
        setTotal(res.total_items);
        setGroups(groupByYearMonth(res.list));
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <div className="mb-10">
        <h1 className="text-3xl font-bold tracking-tight">归档</h1>
        <p className="text-muted-foreground mt-1">共 {total} 篇文章</p>
      </div>

      {loading ? (
        <div className="space-y-8">
          {[0,1,2].map(i => <Skeleton key={i} className="h-40 rounded-xl" />)}
        </div>
      ) : (
        <div className="space-y-12">
          {groups.map(group => (
            <div key={group.year}>
              <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
                <Calendar className="h-5 w-5 text-primary" />
                {group.year}
              </h2>
              <div className="space-y-8">
                {group.months.map(({ month, articles }) => (
                  <div key={month}>
                    <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-widest mb-3">
                      {group.year} 年 {month} 月
                    </h3>
                    <ul className="space-y-3 border-l border-border pl-4">
                      {articles.map(article => (
                        <li key={article.id} className="relative">
                          <div className="absolute -left-[1.3rem] top-2 w-2 h-2 rounded-full bg-primary/40" />
                          <div className="flex items-start justify-between gap-4">
                            <Link
                              href={`/article/${article.slug}`}
                              className="text-sm font-medium hover:text-primary transition-colors line-clamp-1 flex-1"
                            >
                              {article.title}
                            </Link>
                            <div className="flex items-center gap-2 flex-shrink-0">
                              {article.category && (
                                <Badge variant="outline" className="text-xs py-0">
                                  {article.category.name}
                                </Badge>
                              )}
                              <span className="text-xs text-muted-foreground whitespace-nowrap">
                                {new Date(article.published_at || article.created_at).toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })}
                              </span>
                            </div>
                          </div>
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
