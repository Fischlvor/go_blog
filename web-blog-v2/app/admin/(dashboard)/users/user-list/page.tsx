'use client';

export const dynamic = 'force-dynamic';

import { Suspense } from 'react';
import { useEffect, useMemo, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { toast } from 'sonner';
import { adminFreezeUser, adminListUsers, adminUnfreezeUser, type AdminUserListQuery } from '@/lib/client-api/admin/user';
import type { User } from '@/lib/client-api/types';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { UserCellPopover } from '@/components/admin/table-cells/UserCellPopover';

function AdminUserListContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const initialKeyword = searchParams.get('keyword') || '';
  const initialPage = Number(searchParams.get('page') || 1);
  const initialPageSize = Number(searchParams.get('page_size') || 10);

  const [items, setItems] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [keyword, setKeyword] = useState(initialKeyword);
  const [page, setPage] = useState(initialPage);
  const [pageSize, setPageSize] = useState(initialPageSize);
  const [total, setTotal] = useState(0);

  const query: AdminUserListQuery = useMemo(() => ({
    keyword: keyword || null,
    page,
    page_size: pageSize,
  }), [keyword, page, pageSize]);

  const syncQuery = (next: { keyword?: string; page?: number; page_size?: number }) => {
    const params = new URLSearchParams(searchParams.toString());
    if (next.keyword) params.set('keyword', next.keyword);
    else params.delete('keyword');
    params.set('page', String(next.page || 1));
    params.set('page_size', String(next.page_size || pageSize));
    router.replace(`${pathname}?${params.toString()}`);
  };

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListUsers(query);
      setItems(res.list);
      setTotal(res.total_items);
      syncQuery({ keyword, page, page_size: pageSize });
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchList();
  }, [query]);

  const onToggleFreeze = async (item: User) => {
    try {
      if (item.freeze) {
        await adminUnfreezeUser(item.id);
        toast.success('解冻成功');
      } else {
        await adminFreezeUser(item.id);
        toast.success('冻结成功');
      }
      await fetchList();
    } catch {
      toast.error('操作失败');
    }
  };

  return (
    <Card className="overflow-visible">
      <CardHeader>
        <CardTitle>用户列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 overflow-visible">
        <div className="flex items-center justify-end gap-2">
          <Input
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="请输入用户名或邮箱"
            className="w-72"
          />
          <Button onClick={() => setPage(1)}>查询</Button>
        </div>

        {loading ? (
          <p className="text-sm text-muted-foreground">加载中...</p>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>头像</TableHead>
                  <TableHead>用户名</TableHead>
                  <TableHead>UUID</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>注册时间</TableHead>
                  <TableHead>角色</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell>
                      <UserCellPopover user={{ ...u, username: u.nickname }} fallbackText="无用户" />
                    </TableCell>
                    <TableCell>{u.nickname}</TableCell>
                    <TableCell className="max-w-[260px] truncate">{u.uuid}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{formatDate(u.created_at)}</TableCell>
                    <TableCell>{u.role_id === 2 ? '管理员' : '普通用户'}</TableCell>
                    <TableCell>
                      {u.role_id === 1 ? (
                        <Button size="sm" variant="outline" onClick={() => onToggleFreeze(u)}>
                          {u.freeze ? '解冻' : '冻结'}
                        </Button>
                      ) : (
                        '-'
                      )}
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

export default function AdminUserListPage() {
  return (
    <Suspense fallback={<Card><CardContent className="pt-6"><p className="text-sm text-muted-foreground">加载中...</p></CardContent></Card>}>
      <AdminUserListContent />
    </Suspense>
  );
}
