import Link from 'next/link';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import type { Article } from '@/lib/client-api/types';
import { Eye, Heart } from 'lucide-react';
import { cn } from '@/lib/utils';

interface ArticleCardProps {
  article: Article;
  className?: string;
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
}

export function ArticleCard({ article, className }: ArticleCardProps) {
  const coverUrl = article.featured_image || '';

  return (
    <Link href={`/article/${article.slug}`} prefetch={false} className="group block h-full cursor-pointer">
      <Card className={cn(
        'h-full overflow-hidden transition-all duration-300 border-border/60 flex flex-col',
        'hover:shadow-xl hover:shadow-black/5 dark:hover:shadow-black/30',
        'hover:-translate-y-1 hover:border-primary/30',
        'py-0 gap-0',
        className
      )}>
        <div className="relative aspect-video overflow-hidden bg-muted">
          {coverUrl ? (
            <img
              src={coverUrl}
              alt={article.title}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              loading="lazy"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <span className="font-mono text-xs text-muted-foreground/40">// no cover</span>
            </div>
          )}
          {article.category && (
            <div className="absolute top-3 left-3">
              <span className="inline-flex items-center rounded-full bg-background/85 backdrop-blur-sm text-xs border-0 text-foreground shadow-sm px-2 h-6">
                {article.category.name}
              </span>
            </div>
          )}
          {article.is_featured && (
            <div className="absolute top-3 right-3">
              <span className="inline-block w-2 h-2 rounded-full bg-amber-400 shadow-lg shadow-amber-400/50" title="精选" />
            </div>
          )}
        </div>

        <CardContent className="p-3 flex-1 flex flex-col gap-2">
          {/* {article.category && (
            <span className="inline-flex items-center max-w-fit rounded-full bg-primary/20 text-primary px-2 h-6 text-xs font-normal whitespace-nowrap">
              {article.category.name}
            </span>
          )} */}
          <h3 className="text-lg font-semibold leading-7 line-clamp-2 min-h-14 group-hover:text-primary transition-colors">
            {article.title}
          </h3>
          <span className="text-sm text-muted-foreground line-clamp-2 leading-6 min-h-12">
            {article.excerpt || ''}
          </span>
          <div className="flex flex-wrap content-start gap-1 mt-1 min-h-12 max-h-12 overflow-hidden">
            {article.tags?.slice(0, 6).map((tag) => (
              <span key={tag.id} className="inline-flex items-center max-w-fit rounded-full bg-[#ebebeb] text-neutral-600 dark:bg-neutral-700 dark:text-neutral-300 px-2 h-6 text-xs font-normal whitespace-nowrap">
                {tag.name}
              </span>
            ))}
          </div>
        </CardContent>

        <CardFooter className="px-3 py-2 flex items-center justify-between text-xs text-muted-foreground border-0 bg-transparent rounded-none">
          <span>{formatDate(article.published_at)}</span>
          <div className="flex items-center gap-3">
            <span className="flex items-center gap-1"><Eye className="h-3 w-3" />{article.views}</span>
            <span className="flex items-center gap-1"><Heart className="h-3 w-3" />{article.like?.likes ?? 0}</span>
          </div>
        </CardFooter>
      </Card>
    </Link>
  );
}
