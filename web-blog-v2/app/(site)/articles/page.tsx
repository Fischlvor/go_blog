'use client';

import { Suspense } from 'react';
import { useEffect, useState, useCallback } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { ArticleCard } from '@/components/site/article/ArticleCard';
import { listArticles, listCategories } from '@/lib/api/public/article';
import type { Article, ArticleCategory } from '@/lib/api/types';

const PAGE_SIZE = 9;

function ArticlesContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [articles, setArticles] = useState<Article[]>([]);
  const [categories, setCategories] = useState<ArticleCategory[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);

  const page = Number(searchParams.get('page') || 1);
  const keyword = searchParams.get('q') || '';
  const categoryId = searchParams.get('category') ? Number(searchParams.get('category')) : undefined;

  const [inputValue, setInputValue] = useState(keyword);

  const updateParams = (updates: Record<string, string | undefined>) => {
    const params = new URLSearchParams(searchParams.toString());
    Object.entries(updates).forEach(([k, v]) => {
      if (v) params.set(k, v); else params.delete(k);
    });
    params.delete('page');
    router.push(`?${params.toString()}`);
  };

  const fetchArticles = useCallback(async () => {
    setLoading(true);
    try {
      const res = await listArticles({
        page, page_size: PAGE_SIZE,
        keyword: keyword || undefined,
        'filter.category_id': categoryId,
        order: 'desc',
      });
      setArticles(res.list);
      setTotal(res.total_items);
    } catch {
      setArticles([]);
    } finally {
      setLoading(false);
    }
  }, [page, keyword, categoryId]);

  useEffect(() => { fetchArticles(); }, [fetchArticles]);
  useEffect(() => {
    listCategories().then(setCategories).catch(() => {});
  }, []);

  const totalPages = Math.ceil(total / PAGE_SIZE);

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">
      <div className="mb-10 space-y-1">
        <h1 className="text-3xl font-bold tracking-tight">文章</h1>
        <p className="text-muted-foreground">共 {total} 篇文章</p>
      </div>

      {/* Search + Filter */}
      <div className="flex flex-col sm:flex-row gap-4 mb-8">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder="搜索文章..."
            value={inputValue}
            onChange={e => setInputValue(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && updateParams({ q: inputValue || undefined })}
          />
        </div>
      </div>

      {/* Categories */}
      <div className="flex flex-wrap gap-2 mb-8">
        <Badge
          variant={!categoryId ? 'default' : 'outline'}
          className="cursor-pointer"
          onClick={() => updateParams({ category: undefined })}
        >全部</Badge>
        {categories.map(cat => (
          <Badge
            key={cat.id}
            variant={categoryId === cat.id ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => updateParams({ category: String(cat.id) })}
          >
            {cat.name}
            {cat.article_count != null && <span className="ml-1 opacity-60">{cat.article_count}</span>}
          </Badge>
        ))}
      </div>

      {/* Grid */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-72 rounded-xl" />)}
        </div>
      ) : articles.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {articles.map(article => <ArticleCard key={article.id} article={article} />)}
        </div>
      ) : (
        <div className="text-center py-20 text-muted-foreground">
          <p className="text-lg mb-4">没有找到相关文章</p>
          <Button variant="outline" onClick={() => { setInputValue(''); router.push('/articles'); }}>
            清除筛选
          </Button>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-2 mt-12">
          <Button variant="outline" size="sm" disabled={page <= 1}
            onClick={() => { const p = new URLSearchParams(searchParams.toString()); p.set('page', String(page - 1)); router.push(`?${p.toString()}`); }}
          >上一页</Button>
          <span className="flex items-center px-4 text-sm text-muted-foreground">{page} / {totalPages}</span>
          <Button variant="outline" size="sm" disabled={page >= totalPages}
            onClick={() => { const p = new URLSearchParams(searchParams.toString()); p.set('page', String(page + 1)); router.push(`?${p.toString()}`); }}
          >下一页</Button>
        </div>
      )}
    </div>
  );
}

export default function ArticlesPage() {
  return (
    <Suspense fallback={
      <div className="max-w-6xl mx-auto px-4 py-12">
        <div className="mb-10 space-y-1">
          <Skeleton className="h-10 w-20" />
          <Skeleton className="h-4 w-32" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-72 rounded-xl" />)}
        </div>
      </div>
    }>
      <ArticlesContent />
    </Suspense>
  );
}
