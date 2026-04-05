'use client';

import { useMemo, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { ArticleCategory } from '@/lib/client-api/types';

interface ArticlesFiltersProps {
  categories: ArticleCategory[];
  keyword: string;
  categoryId?: number;
}

export function ArticlesFilters({ categories, keyword, categoryId }: ArticlesFiltersProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [inputValue, setInputValue] = useState(keyword);

  const createQueryString = (updates: Record<string, string | undefined>) => {
    const params = new URLSearchParams(searchParams.toString());
    Object.entries(updates).forEach(([key, value]) => {
      if (value) {
        params.set(key, value);
      } else {
        params.delete(key);
      }
    });
    params.delete('page');
    return params.toString();
  };

  const applyFilters = (updates: Record<string, string | undefined>) => {
    const nextQuery = createQueryString(updates);
    router.push(nextQuery ? `${pathname}?${nextQuery}` : pathname);
  };

  const hasFilters = useMemo(() => Boolean(keyword || categoryId), [keyword, categoryId]);

  return (
    <>
      <div className="flex flex-col sm:flex-row gap-4 mb-8">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder="搜索文章..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && applyFilters({ q: inputValue || undefined })}
          />
        </div>
        <Button variant="outline" onClick={() => applyFilters({ q: inputValue || undefined })}>
          搜索
        </Button>
        {hasFilters && (
          <Button
            variant="ghost"
            onClick={() => {
              setInputValue('');
              router.push(pathname);
            }}
          >
            清除筛选
          </Button>
        )}
      </div>

      <div className="flex flex-wrap gap-2 mb-8">
        <Badge
          variant={!categoryId ? 'default' : 'outline'}
          className="cursor-pointer"
          onClick={() => applyFilters({ category: undefined })}
        >
          全部
        </Badge>
        {categories.map((cat) => (
          <Badge
            key={cat.id}
            variant={categoryId === cat.id ? 'default' : 'outline'}
            className="cursor-pointer"
            onClick={() => applyFilters({ category: String(cat.id) })}
          >
            {cat.name}
            {cat.article_count != null && <span className="ml-1 opacity-60">{cat.article_count}</span>}
          </Badge>
        ))}
      </div>
    </>
  );
}
