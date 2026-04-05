import { ArticleCard } from '@/components/site/article/ArticleCard';
import { ArticlesFilters } from '@/components/site/article/ArticlesFilters';
import { ArticlesPagination } from '@/components/site/article/ArticlesPagination';
import { Button } from '@/components/ui/button';
import { listArticlesServer, listCategoriesServer } from '@/lib/server-api/article';

const PAGE_SIZE = 9;

interface ArticlesPageProps {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}

function getSingleParam(value: string | string[] | undefined): string | undefined {
  return Array.isArray(value) ? value[0] : value;
}

export default async function ArticlesPage({ searchParams }: ArticlesPageProps) {
  const params = await searchParams;
  const pageParam = Number(getSingleParam(params.page) || 1);
  const page = Number.isFinite(pageParam) && pageParam > 0 ? pageParam : 1;
  const keyword = getSingleParam(params.q) || '';
  const categoryParam = getSingleParam(params.category);
  const categoryId = categoryParam ? Number(categoryParam) : undefined;

  const [categories, articleResult] = await Promise.all([
    listCategoriesServer(),
    listArticlesServer({
      page,
      page_size: PAGE_SIZE,
      keyword: keyword || undefined,
      'filter.category_id': typeof categoryId === 'number' && !Number.isNaN(categoryId) ? categoryId : undefined,
      order: 'desc',
    }),
  ]);

  const totalPages = Math.ceil(articleResult.total_items / PAGE_SIZE);
  const currentQuery = new URLSearchParams();
  if (keyword) currentQuery.set('q', keyword);
  if (typeof categoryId === 'number' && !Number.isNaN(categoryId)) currentQuery.set('category', String(categoryId));

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">
      <div className="mb-10 space-y-1">
        <h1 className="text-3xl font-bold tracking-tight">文章</h1>
        <p className="text-muted-foreground">共 {articleResult.total_items} 篇文章</p>
      </div>

      <ArticlesFilters categories={categories} keyword={keyword} categoryId={categoryId} />

      {articleResult.list.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {articleResult.list.map((article) => (
            <ArticleCard key={article.id} article={article} />
          ))}
        </div>
      ) : (
        <div className="text-center py-20 text-muted-foreground">
          <p className="text-lg mb-4">没有找到相关文章</p>
          <Button variant="outline" render={<a href="/articles" />}>
            清除筛选
          </Button>
        </div>
      )}

      <ArticlesPagination
        currentPage={page}
        totalPages={totalPages}
        basePathname="/articles"
        queryString={currentQuery.toString()}
      />
    </div>
  );
}
