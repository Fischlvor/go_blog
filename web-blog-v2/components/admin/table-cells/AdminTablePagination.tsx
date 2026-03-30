'use client';

import { Button } from '@/components/ui/button';

interface AdminTablePaginationProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  pageSizeOptions?: number[];
}

export function AdminTablePagination({
  page,
  pageSize,
  total,
  onPageChange,
  onPageSizeChange,
  pageSizeOptions = [10, 30, 50, 100],
}: AdminTablePaginationProps) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const canPrev = page > 1;
  const canNext = page < totalPages;

  return (
    <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
      <div className="text-sm text-muted-foreground">共 {total} 条，第 {page}/{totalPages} 页</div>

      <div className="flex items-center gap-2">
        <select
          className="h-8 rounded-md border border-border bg-background px-2 text-sm"
          value={pageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
        >
          {pageSizeOptions.map((size) => (
            <option key={size} value={size}>
              {size} / 页
            </option>
          ))}
        </select>

        <Button size="sm" variant="outline" disabled={!canPrev} onClick={() => onPageChange(page - 1)}>
          上一页
        </Button>
        <Button size="sm" variant="outline" disabled={!canNext} onClick={() => onPageChange(page + 1)}>
          下一页
        </Button>
      </div>
    </div>
  );
}
