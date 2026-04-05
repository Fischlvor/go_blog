'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { adminListLoginRecords, type LoginRecord } from '@/lib/client-api/admin/user';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { UserCellPopover } from '@/components/admin/table-cells/UserCellPopover';

export default function AdminLoginLogsPage() {
  const [items, setItems] = useState<LoginRecord[]>([]);
  const [uuid, setUUID] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    try {
      const res = await adminListLoginRecords({ page, page_size: pageSize, uuid: uuid || null });
      setItems(res.list || []);
      setTotal(res.total_items || 0);
    } catch {
      setItems([]);
      setTotal(0);
    }
  };

  useEffect(() => {
    adminListLoginRecords({ page, page_size: pageSize, uuid: uuid || null })
      .then((res) => {
        setItems(res.list || []);
        setTotal(res.total_items || 0);
      })
      .catch(() => {
        setItems([]);
        setTotal(0);
      });
  }, [page, pageSize, uuid]);

  return (
    <Card>
      <CardHeader><CardTitle>登录日志</CardTitle></CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center justify-end gap-2">
          <Input value={uuid} onChange={(e) => setUUID(e.target.value)} placeholder="请输入用户UUID" className="w-72" />
          <Button onClick={() => { setPage(1); void fetchList(); }}>查询</Button>
        </div>

        <Table>
          <TableHeader><TableRow><TableHead>用户</TableHead><TableHead>用户名</TableHead><TableHead>登录时间</TableHead><TableHead>登录方式</TableHead><TableHead>IP</TableHead><TableHead>登录地址</TableHead><TableHead>操作系统</TableHead><TableHead>设备信息</TableHead><TableHead>浏览器信息</TableHead><TableHead>登录状态</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell><UserCellPopover user={it.user} /></TableCell><TableCell>{it.user?.nickname || it.user?.username || '-'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell>{it.login_method}</TableCell><TableCell>{it.ip}</TableCell><TableCell>{it.address}</TableCell><TableCell>{it.os}</TableCell><TableCell className="max-w-[180px] truncate">{it.device_info}</TableCell><TableCell className="max-w-[180px] truncate">{it.browser_info}</TableCell><TableCell>{it.status}</TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
      </CardContent>
    </Card>
  );
}
