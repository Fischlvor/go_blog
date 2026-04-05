'use client';

export const dynamic = 'force-dynamic';

import { Suspense } from 'react';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { toast } from 'sonner';
import { adminDeleteArticles, adminListArticles, type AdminArticleListQuery } from '@/lib/client-api/admin/article';
import type { Article } from '@/lib/client-api/types';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';

function AdminArticleListContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [items, setItems] = useState<Article[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

  const [titleKeyword, setTitleKeyword] = useState(searchParams.get('keyword') || '');
  const [categoryKeyword, setCategoryKeyword] = useState(searchParams.get('category') || '');
  const [abstractKeyword, setAbstractKeyword] = useState(searchParams.get('abstract') || '');
  const [page, setPage] = useState(Number(searchParams.get('page') || 1));
  const [pageSize, setPageSize] = useState(Number(searchParams.get('page_size') || 10));
  const [total, setTotal] = useState(0);

  const query: AdminArticleListQuery = useMemo(() => ({
    keyword: titleKeyword || null,
    page,
    page_size: pageSize,
  }), [titleKeyword, page, pageSize]);

  const syncQuery = () => {
    const params = new URLSearchParams();
    if (titleKeyword) params.set('keyword', titleKeyword);
    if (categoryKeyword) params.set('category', categoryKeyword);
    if (abstractKeyword) params.set('abstract', abstractKeyword);
    params.set('page', String(page));
    params.set('page_size', String(pageSize));
    router.replace(`${pathname}?${params.toString()}`);
  };

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListArticles(query);
      let list = res.list;
      if (categoryKeyword) list = list.filter((item) => item.category?.name?.includes(categoryKeyword));
      if (abstractKeyword) list = list.filter((item) => item.excerpt?.includes(abstractKeyword));

      setItems(list);
      setTotal(res.total_items);
      syncQuery();
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取文章列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchList();
  }, [query, categoryKeyword, abstractKeyword]);

  const toggleSelect = (slug: string, checked: boolean) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (checked) next.add(slug);
      else next.delete(slug);
      return next;
    });
  };

  const toggleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedIds(new Set(items.map((item) => item.slug)));
      return;
    }
    setSelectedIds(new Set());
  };

  const handleDelete = async (ids: string[]) => {
    if (!ids.length) {
      toast.warning('请先选择要删除的文章');
      return;
    }
    try {
      await adminDeleteArticles(ids);
      toast.success('删除成功');
      setSelectedIds(new Set());
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  const allChecked = items.length > 0 && items.every((item) => selectedIds.has(String(item.id)));

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between">
        <CardTitle>文章列表</CardTitle>
      </CardHeader>

      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={() => router.push('/admin/articles/article-publish')}>新建文章</Button>
            <Button variant="destructive" onClick={() => handleDelete(Array.from(selectedIds))}>批量删除</Button>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Input value={titleKeyword} onChange={(e) => setTitleKeyword(e.target.value)} placeholder="文章标题" className="w-52" />
            <Input value={categoryKeyword} onChange={(e) => setCategoryKeyword(e.target.value)} placeholder="文章类别" className="w-44" />
            <Input value={abstractKeyword} onChange={(e) => setAbstractKeyword(e.target.value)} placeholder="文章简介" className="w-52" />
            <Button onClick={() => setPage(1)}>查询</Button>
          </div>
        </div>

        {loading ? (
          <p className="text-sm text-muted-foreground">加载中...</p>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12">
                    <input type="checkbox" checked={allChecked} onChange={(e) => toggleSelectAll(e.target.checked)} />
                  </TableHead>
                  <TableHead className="w-[120px]">封面</TableHead>
                  <TableHead>标题</TableHead>
                  <TableHead>类别</TableHead>
                  <TableHead>标签</TableHead>
                  <TableHead>简介</TableHead>
                  <TableHead>发布时间</TableHead>
                  <TableHead>可见性</TableHead>
                  <TableHead>文章id</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map((a) => {
                  const slug = a.slug;
                  return (
                    <TableRow key={a.id} className="h-20">
                      <TableCell>
                        <input
                          type="checkbox"
                          checked={selectedIds.has(slug)}
                          onChange={(e) => toggleSelect(slug, e.target.checked)}
                        />
                      </TableCell>
                      <TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={a.featured_image} alt={a.title} boxWidth={96} boxHeight={64} /></div></TableCell>
                      <TableCell className="max-w-[180px] truncate">{a.title}</TableCell>
                      <TableCell>{a.category?.name || '-'}</TableCell>
                      <TableCell className="max-w-[180px]">
                        <div className="flex flex-nowrap gap-1 overflow-hidden">
                          {(a.tags || []).map((tag) => <Badge key={tag.id} variant="outline">{tag.name}</Badge>)}
                        </div>
                      </TableCell>
                      <TableCell className="max-w-[260px] truncate">{a.excerpt || '-'}</TableCell>
                      <TableCell>{formatDate(a.created_at)}</TableCell>
                      <TableCell>
                        {a.status === 'draft' ? (
                          <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
                            草稿
                          </Badge>
                        ) : (
                          <Badge variant={a.visibility === 'public' ? 'default' : a.visibility === 'private' ? 'secondary' : 'outline'}>
                            {a.visibility === 'public' ? '公开' : a.visibility === 'private' ? '私密' : '-'}
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="max-w-[220px] truncate">
                        <Link href={`/article/${a.slug}`} prefetch={false} className="text-primary hover:underline">{a.slug}</Link>
                      </TableCell>
                      <TableCell className="space-x-2">
                        <Button size="sm" variant="outline" onClick={() => router.push(`/admin/articles/article-publish?slug=${a.slug}`)}>更新</Button>
                        <Button size="sm" variant="destructive" onClick={() => handleDelete([slug])}>删除</Button>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>

            <AdminTablePagination
              page={page}
              pageSize={pageSize}
              total={total}
              onPageChange={setPage}
              onPageSizeChange={(size) => {
                setPageSize(size);
                setPage(1);
              }}
            />
          </>
        )}
      </CardContent>
    </Card>
  );
}

export default function AdminArticleListPage() {
  return (
    <Suspense fallback={<Card><CardContent className="pt-6"><p className="text-sm text-muted-foreground">加载中...</p></CardContent></Card>}>
      <AdminArticleListContent />
    </Suspense>
  );
}
