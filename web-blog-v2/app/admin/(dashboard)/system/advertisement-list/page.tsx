'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import {
  adminCreateAdvertisement,
  adminDeleteAdvertisements,
  adminListAdvertisements,
  adminUpdateAdvertisement,
  type Advertisement,
} from '@/lib/api/admin/advertisement';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';

interface AdvertisementForm { ad_image: string; link: string; title: string; content: string; }
const DEFAULT_FORM: AdvertisementForm = { ad_image: '', link: '', title: '', content: '' };

export default function AdminAdvertisementListPage() {
  const [items, setItems] = useState<Advertisement[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Advertisement | null>(null);
  const [form, setForm] = useState<AdvertisementForm>(DEFAULT_FORM);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListAdvertisements({ page, page_size: pageSize, title: title || null, content: content || null });
      setItems(res.list);
      setTotal(res.total_items);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取广告列表失败');
    } finally { setLoading(false); }
  };

  useEffect(() => { void fetchList(); }, [page, pageSize]);

  const openCreate = () => { setEditing(null); setForm(DEFAULT_FORM); setOpen(true); };
  const openEdit = (item: Advertisement) => { setEditing(item); setForm({ ad_image: item.ad_image || '', link: item.link || '', title: item.title || '', content: item.content || '' }); setOpen(true); };

  const submit = async () => {
    if (!form.title.trim() || !form.link.trim()) { toast.error('标题和链接不能为空'); return; }
    setSubmitting(true);
    try {
      if (editing) { await adminUpdateAdvertisement({ id: editing.id, ...form }); toast.success('更新成功'); }
      else { await adminCreateAdvertisement(form); toast.success('创建成功'); }
      setOpen(false);
      await fetchList();
    } catch { toast.error('操作失败'); }
    finally { setSubmitting(false); }
  };

  const onDelete = async (ids: number[]) => {
    if (!ids.length) { toast.warning('请先选择要删除的广告'); return; }
    setSubmitting(true);
    try {
      await adminDeleteAdvertisements({ ids });
      toast.success('删除成功');
      setSelectedIds(new Set());
      await fetchList();
    } catch { toast.error('删除失败'); }
    finally { setSubmitting(false); }
  };

  const allChecked = useMemo(() => items.length > 0 && items.every((it) => selectedIds.has(it.id)), [items, selectedIds]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>广告列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2"><div className="space-x-2"><Button onClick={openCreate}>新建广告</Button><Button variant="destructive" onClick={() => onDelete(Array.from(selectedIds))}>批量删除</Button></div><div className="flex items-center gap-2"><Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="广告标题" className="w-52" /><Input value={content} onChange={(e) => setContent(e.target.value)} placeholder="广告内容" className="w-52" /><Button onClick={() => { setPage(1); void fetchList(); }}>查询</Button></div></div>

        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead className="w-12"><input type="checkbox" checked={allChecked} onChange={(e) => setSelectedIds(e.target.checked ? new Set(items.map((it) => it.id)) : new Set())} /></TableHead><TableHead className="w-[120px]">图片</TableHead><TableHead>标题</TableHead><TableHead>链接</TableHead><TableHead>内容</TableHead><TableHead>创建时间</TableHead><TableHead className="text-right">操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => (<TableRow key={it.id} className="h-20"><TableCell><input type="checkbox" checked={selectedIds.has(it.id)} onChange={(e) => setSelectedIds((prev) => { const n = new Set(prev); if (e.target.checked) n.add(it.id); else n.delete(it.id); return n; })} /></TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.ad_image} alt={it.title} /></div></TableCell><TableCell>{it.title}</TableCell><TableCell className="max-w-[260px] truncate">{it.link}</TableCell><TableCell className="max-w-[320px] truncate">{it.content}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="text-right space-x-2"><Button size="sm" variant="outline" onClick={() => openEdit(it)}>更新</Button><Button size="sm" variant="destructive" onClick={() => onDelete([it.id])} disabled={submitting}>删除</Button></TableCell></TableRow>))}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
          </>
        )}

        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? '更新广告' : '新建广告'}</DialogTitle></DialogHeader>
            <div className="space-y-3"><Input placeholder="标题" value={form.title} onChange={(e) => setForm((p) => ({ ...p, title: e.target.value }))} /><Input placeholder="链接" value={form.link} onChange={(e) => setForm((p) => ({ ...p, link: e.target.value }))} /><Input placeholder="图片 URL" value={form.ad_image} onChange={(e) => setForm((p) => ({ ...p, ad_image: e.target.value }))} /><Input placeholder="内容" value={form.content} onChange={(e) => setForm((p) => ({ ...p, content: e.target.value }))} /></div>
            <DialogFooter><Button variant="outline" onClick={() => setOpen(false)}>取消</Button><Button onClick={submit} disabled={submitting}>{editing ? '更新' : '创建'}</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
