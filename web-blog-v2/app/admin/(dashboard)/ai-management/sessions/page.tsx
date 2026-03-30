'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { adminDeleteAISession, adminGetAISession, adminListAISessions, type AIMessage, type AISession } from '@/lib/api/admin/ai-management';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';

export default function AdminAISessionsPage() {
  const [items, setItems] = useState<AISession[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [userUUID, setUserUUID] = useState('');
  const [model, setModel] = useState('');
  const [detailOpen, setDetailOpen] = useState(false);
  const [detail, setDetail] = useState<(AISession & { messages?: AIMessage[] }) | null>(null);

  const fetchList = async () => {
    try {
      const res = await adminListAISessions({ page, page_size: pageSize, user_uuid: userUUID || undefined, keyword: model || undefined });
      setItems(res.list || []);
      setTotal(res.total_items || 0);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取AI会话列表失败');
    }
  };

  useEffect(() => {
    adminListAISessions({ page, page_size: pageSize, user_uuid: userUUID || undefined, keyword: model || undefined })
      .then((res) => {
        setItems(res.list || []);
        setTotal(res.total_items || 0);
      })
      .catch(() => {
        setItems([]);
        setTotal(0);
        toast.error('获取AI会话列表失败');
      });
  }, [page, pageSize, userUUID, model]);

  const onSearch = () => {
    setPage(1);
    void fetchList();
  };

  const onDetail = async (id: number) => {
    try {
      const res = await adminGetAISession(id);
      setDetail(res);
      setDetailOpen(true);
    } catch {
      toast.error('获取会话详情失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await adminDeleteAISession(id);
      toast.success('删除成功');
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>AI会话管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2"><div /><div className="flex items-center gap-2"><Input value={userUUID} onChange={(e) => setUserUUID(e.target.value)} placeholder="用户UUID" className="w-56" /><Input value={model} onChange={(e) => setModel(e.target.value)} placeholder="模型" className="w-48" /><Button onClick={onSearch}>搜索</Button></div></div>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>标题</TableHead><TableHead>用户UUID</TableHead><TableHead>模型</TableHead><TableHead>创建时间</TableHead><TableHead>更新时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell>{it.id}</TableCell><TableCell>{it.title || '-'}</TableCell><TableCell>{it.user_uuid || '-'}</TableCell><TableCell>{it.model || '-'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell>{formatDate(it.updated_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => onDetail(it.id)}>查看详情</Button><Button size="sm" variant="destructive" onClick={() => onDelete(it.id)}>删除</Button></TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />

        <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>会话详情</DialogTitle></DialogHeader>
            <pre className="max-h-[420px] overflow-auto rounded border p-3 text-xs">{JSON.stringify(detail, null, 2)}</pre>
            <DialogFooter><Button variant="outline" onClick={() => setDetailOpen(false)}>关闭</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
