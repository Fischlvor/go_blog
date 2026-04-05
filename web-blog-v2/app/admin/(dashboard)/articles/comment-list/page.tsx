'use client';

export const dynamic = 'force-dynamic';

import { Suspense } from 'react';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { toast } from 'sonner';
import { adminDeleteComment, adminListComments, type AdminCommentListQuery } from '@/lib/client-api/admin/comment';
import type { Comment } from '@/lib/client-api/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { UserCellPopover } from '@/components/admin/table-cells/UserCellPopover';
import { MarkdownContent } from '@/components/site/article/MarkdownContent';

function AdminCommentListContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [items, setItems] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());

  const [articleSlug, setArticleSlug] = useState(searchParams.get('article_slug') || '');
  const [userUUID, setUserUUID] = useState(searchParams.get('user_uuid') || '');
  const [contentKeyword, setContentKeyword] = useState(searchParams.get('content') || '');
  const [page, setPage] = useState(Number(searchParams.get('page') || 1));
  const [pageSize, setPageSize] = useState(Number(searchParams.get('page_size') || 10));
  const [total, setTotal] = useState(0);

  const query: AdminCommentListQuery = useMemo(() => ({
    page,
    page_size: pageSize,
    article_slug: articleSlug || null,
    user_uuid: userUUID || null,
    keyword: contentKeyword || null,
  }), [page, pageSize, articleSlug, userUUID, contentKeyword]);

  const syncQuery = () => {
    const params = new URLSearchParams();
    if (articleSlug) params.set('article_slug', articleSlug);
    if (userUUID) params.set('user_uuid', userUUID);
    if (contentKeyword) params.set('content', contentKeyword);
    params.set('page', String(page));
    params.set('page_size', String(pageSize));
    router.replace(`${pathname}?${params.toString()}`);
  };

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListComments(query);
      let list = res.list;
      if (contentKeyword) list = list.filter((item) => item.content?.includes(contentKeyword));
      setItems(list);
      setTotal(res.total_items);
      syncQuery();
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取评论列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchList();
  }, [query]);

  const toggleSelect = (id: number, checked: boolean) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (checked) next.add(id);
      else next.delete(id);
      return next;
    });
  };

  const deleteComments = async (ids: number[]) => {
    if (!ids.length) {
      toast.warning('请先选择要删除的评论');
      return;
    }
    try {
      await Promise.all(ids.map((id) => adminDeleteComment(id)));
      toast.success('删除成功');
      setSelectedIds(new Set());
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  return (
    <Card className="overflow-visible">
      <CardHeader>
        <CardTitle>评论列表</CardTitle>
      </CardHeader>

      <CardContent className="space-y-4 overflow-visible">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <Button variant="destructive" onClick={() => deleteComments(Array.from(selectedIds))}>批量删除</Button>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Input value={articleSlug} onChange={(e) => setArticleSlug(e.target.value)} placeholder="文章slug" className="w-52" />
            <Input value={userUUID} onChange={(e) => setUserUUID(e.target.value)} placeholder="用户uuid" className="w-52" />
            <Input value={contentKeyword} onChange={(e) => setContentKeyword(e.target.value)} placeholder="评论内容" className="w-52" />
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
                  <TableHead className="w-12" />
                  <TableHead>文章slug</TableHead>
                  <TableHead>用户</TableHead>
                  <TableHead>内容</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>
                      <input
                        type="checkbox"
                        checked={selectedIds.has(item.id)}
                        onChange={(e) => toggleSelect(item.id, e.target.checked)}
                      />
                    </TableCell>
                    <TableCell>
                      <Link href={`/article/${item.article_slug}`} className="text-primary hover:underline">
                        {item.article_slug}
                      </Link>
                    </TableCell>
                    <TableCell>
                      <UserCellPopover user={{ uuid: item.user?.uuid, nickname: item.user?.nickname, avatar: item.user?.avatar }} />
                    </TableCell>
                    <TableCell className="max-w-[520px] align-top">
                      <div className="max-h-40 overflow-auto rounded-md border border-border p-2">
                        <MarkdownContent content={item.content} className="prose-sm" />
                      </div>
                    </TableCell>
                    <TableCell>
                      <Button size="sm" variant="destructive" onClick={() => deleteComments([item.id])}>删除</Button>
                    </TableCell>
                  </TableRow>
                ))}
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

export default function AdminCommentListPage() {
  return (
    <Suspense fallback={<Card><CardContent className="pt-6"><p className="text-sm text-muted-foreground">加载中...</p></CardContent></Card>}>
      <AdminCommentListContent />
    </Suspense>
  );
}
