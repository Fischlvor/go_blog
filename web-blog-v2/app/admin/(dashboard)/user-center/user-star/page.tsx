'use client';

export const dynamic = 'force-dynamic';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { getUserLikedArticles } from '@/lib/api/user/article';
import type { Article } from '@/lib/api/types';
import { formatDate } from '@/lib/date';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';

export default function UserStarPage() {
  const [items, setItems] = useState<Article[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    getUserLikedArticles(page, pageSize)
      .then((res) => {
        setItems(res.list || []);
        setTotal(res.total_items || 0);
      })
      .catch(() => {
        setItems([]);
        setTotal(0);
      });
  }, [page, pageSize]);

  return (
    <Card>
      <CardHeader><CardTitle>我的收藏</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <Table>
          <TableHeader><TableRow><TableHead className="w-[120px]">封面</TableHead><TableHead>标题</TableHead><TableHead>类别</TableHead><TableHead>标签</TableHead><TableHead>简介</TableHead><TableHead>发布时间</TableHead><TableHead>文章id</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id} className="h-20"><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.featured_image} alt={it.title} /></div></TableCell><TableCell>{it.title}</TableCell><TableCell>{it.category?.name || '-'}</TableCell><TableCell>{(it.tags || []).map((t) => t.name).join('、') || '-'}</TableCell><TableCell className="max-w-[300px] truncate">{it.excerpt}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell><Link href={`/article/${it.slug}`} prefetch={false} className="text-primary hover:underline">{it.slug}</Link></TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
      </CardContent>
    </Card>
  );
}
