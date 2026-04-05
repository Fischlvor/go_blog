'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { adminDeleteAIMessage, adminGetAIMessage, adminListAIMessages, type AIMessage } from '@/lib/client-api/admin/ai-management';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';

export default function AdminAIMessagesPage() {
  const [items, setItems] = useState<AIMessage[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [sessionId, setSessionId] = useState('');
  const [role, setRole] = useState('');
  const [detailOpen, setDetailOpen] = useState(false);
  const [detail, setDetail] = useState<AIMessage | null>(null);

  const fetchList = async () => {
    try {
      const res = await adminListAIMessages({
        page,
        page_size: pageSize,
        session_id: sessionId ? Number(sessionId) : undefined,
        role: role || undefined,
      });
      setItems(res.list || []);
      setTotal(res.total_items || 0);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取AI消息列表失败');
    }
  };

  useEffect(() => {
    adminListAIMessages({
      page,
      page_size: pageSize,
      session_id: sessionId ? Number(sessionId) : undefined,
      role: role || undefined,
    })
      .then((res) => {
        setItems(res.list || []);
        setTotal(res.total_items || 0);
      })
      .catch(() => {
        setItems([]);
        setTotal(0);
        toast.error('获取AI消息列表失败');
      });
  }, [page, pageSize, sessionId, role]);

  const onSearch = () => {
    setPage(1);
    void fetchList();
  };

  const onDetail = async (id: number) => {
    try {
      const res = await adminGetAIMessage(id);
      setDetail(res);
      setDetailOpen(true);
    } catch {
      toast.error('获取消息详情失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await adminDeleteAIMessage(id);
      toast.success('删除成功');
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>AI消息管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2"><div /><div className="flex items-center gap-2"><Input value={sessionId} onChange={(e) => setSessionId(e.target.value)} placeholder="会话ID" className="w-40" /><Input value={role} onChange={(e) => setRole(e.target.value)} placeholder="角色(user/assistant)" className="w-52" /><Button onClick={onSearch}>搜索</Button></div></div>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>会话ID</TableHead><TableHead>角色</TableHead><TableHead>内容</TableHead><TableHead>Token数</TableHead><TableHead>创建时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell>{it.id}</TableCell><TableCell>{it.session_id}</TableCell><TableCell>{it.role}</TableCell><TableCell className="max-w-[360px] truncate">{it.content}</TableCell><TableCell>{it.tokens ?? '-'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => onDetail(it.id)}>查看详情</Button><Button size="sm" variant="destructive" onClick={() => onDelete(it.id)}>删除</Button></TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />

        <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>消息详情</DialogTitle></DialogHeader>
            <pre className="max-h-[420px] overflow-auto rounded border p-3 text-xs">{JSON.stringify(detail, null, 2)}</pre>
            <DialogFooter><Button variant="outline" onClick={() => setDetailOpen(false)}>关闭</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
