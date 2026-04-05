'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { adminDeleteFeedbacks, adminListFeedbacks, adminReplyFeedback } from '@/lib/client-api/admin/feedback';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Textarea } from '@/components/ui/textarea';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import type { Feedback } from '@/lib/client-api/types';

export default function AdminFeedbackListPage() {
  const [items, setItems] = useState<Feedback[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [replyOpen, setReplyOpen] = useState(false);
  const [replyText, setReplyText] = useState('');
  const [current, setCurrent] = useState<Feedback | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListFeedbacks({ page, page_size: pageSize });
      setItems(res.list);
      setTotal(res.total_items);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取反馈列表失败');
    } finally { setLoading(false); }
  };

  useEffect(() => { void fetchList(); }, [page, pageSize]);

  const onDelete = async (ids: number[]) => {
    if (!ids.length) { toast.warning('请先选择要删除的反馈'); return; }
    setSubmitting(true);
    try {
      await adminDeleteFeedbacks({ ids });
      toast.success('删除成功');
      setSelectedIds(new Set());
      await fetchList();
    } catch { toast.error('删除失败'); }
    finally { setSubmitting(false); }
  };

  const openReply = (item: Feedback) => { setCurrent(item); setReplyText(item.reply || ''); setReplyOpen(true); };

  const submitReply = async () => {
    if (!current) return;
    if (!replyText.trim()) { toast.error('回复内容不能为空'); return; }
    setSubmitting(true);
    try {
      await adminReplyFeedback({ id: current.id, reply: replyText.trim() });
      toast.success('回复成功');
      setReplyOpen(false);
      await fetchList();
    } catch { toast.error('回复失败'); }
    finally { setSubmitting(false); }
  };

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between"><CardTitle>反馈列表</CardTitle><Button variant="destructive" onClick={() => onDelete(Array.from(selectedIds))}>批量删除</Button></CardHeader>
      <CardContent className="space-y-4">
        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead className="w-12" /><TableHead>用户</TableHead><TableHead>时间</TableHead><TableHead>内容</TableHead><TableHead>回复</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => (<TableRow key={it.id}><TableCell><input type="checkbox" checked={selectedIds.has(it.id)} onChange={(e) => setSelectedIds((prev) => { const n = new Set(prev); if (e.target.checked) n.add(it.id); else n.delete(it.id); return n; })} /></TableCell><TableCell>{it.user_uuid || '-'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="max-w-[300px] truncate">{it.content}</TableCell><TableCell className="max-w-[260px] truncate">{it.reply || '-'}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => openReply(it)}>{it.reply ? '修改回复' : '回复'}</Button><Button size="sm" variant="destructive" onClick={() => onDelete([it.id])} disabled={submitting}>删除</Button></TableCell></TableRow>))}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
          </>
        )}

        <Dialog open={replyOpen} onOpenChange={setReplyOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>回复反馈</DialogTitle></DialogHeader>
            <Textarea value={replyText} onChange={(e) => setReplyText(e.target.value)} rows={6} placeholder="请输入回复内容" />
            <DialogFooter><Button variant="outline" onClick={() => setReplyOpen(false)}>取消</Button><Button onClick={submitReply} disabled={submitting}>提交回复</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
