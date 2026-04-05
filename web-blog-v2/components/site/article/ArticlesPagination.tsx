import { Button } from '@/components/ui/button';

interface ArticlesPaginationProps {
  currentPage: number;
  totalPages: number;
  basePathname: string;
  queryString: string;
}

function buildHref(basePathname: string, queryString: string, page: number): string {
  const params = new URLSearchParams(queryString);
  if (page <= 1) {
    params.delete('page');
  } else {
    params.set('page', String(page));
  }
  const nextQuery = params.toString();
  return nextQuery ? `${basePathname}?${nextQuery}` : basePathname;
}

export function ArticlesPagination({
  currentPage,
  totalPages,
  basePathname,
  queryString,
}: ArticlesPaginationProps) {
  if (totalPages <= 1) {
    return null;
  }

  return (
    <div className="flex justify-center gap-2 mt-12">
      <Button variant="outline" size="sm" disabled={currentPage <= 1} render={<a href={buildHref(basePathname, queryString, currentPage - 1)} />}>
        上一页
      </Button>
      <span className="flex items-center px-4 text-sm text-muted-foreground">
        {currentPage} / {totalPages}
      </span>
      <Button variant="outline" size="sm" disabled={currentPage >= totalPages} render={<a href={buildHref(basePathname, queryString, currentPage + 1)} />}>
        下一页
      </Button>
    </div>
  );
}
