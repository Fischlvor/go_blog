'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { adminCreateFriendLink, adminDeleteFriendLinks, adminListFriendLinks, adminUpdateFriendLink } from '@/lib/client-api/admin/friendLink';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';
import type { FriendLink } from '@/lib/client-api/types';

interface FriendLinkForm {
  logo: string;
  link: string;
  name: string;
  description: string;
}

const DEFAULT_FORM: FriendLinkForm = { logo: '', link: '', name: '', description: '' };

export default function AdminFriendLinkListPage() {
  const [items, setItems] = useState<FriendLink[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<FriendLink | null>(null);
  const [form, setForm] = useState<FriendLinkForm>(DEFAULT_FORM);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListFriendLinks({ page, page_size: pageSize, name: name || null, description: description || null });
      setItems(res.list);
      setTotal(res.total_items);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取友链列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { void fetchList(); }, [page, pageSize]);

  const resetForm = () => { setEditing(null); setForm(DEFAULT_FORM); };

  const openCreate = () => { resetForm(); setOpen(true); };

  const openEdit = (item: FriendLink) => {
    setEditing(item);
    setForm({ logo: item.logo || '', link: item.link || item.url || '', name: item.name || '', description: item.description || '' });
    setOpen(true);
  };

  const submit = async () => {
    if (!form.name.trim() || !form.link.trim()) { toast.error('名称和链接不能为空'); return; }
    setSubmitting(true);
    try {
      if (editing) { await adminUpdateFriendLink({ id: editing.id, ...form }); toast.success('更新成功'); }
      else { await adminCreateFriendLink(form); toast.success('创建成功'); }
      setOpen(false);
      resetForm();
      await fetchList();
    } catch { toast.error('操作失败'); }
    finally { setSubmitting(false); }
  };

  const onDelete = async (ids: number[]) => {
    if (!ids.length) { toast.warning('请先选择要删除的友链'); return; }
    setSubmitting(true);
    try {
      await adminDeleteFriendLinks({ ids });
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
        <CardTitle>友链列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2"><div className="space-x-2"><Button onClick={openCreate}>新建友链</Button><Button variant="destructive" onClick={() => onDelete(Array.from(selectedIds))}>批量删除</Button></div><div className="flex items-center gap-2"><Input value={name} onChange={(e) => setName(e.target.value)} placeholder="友链名称" className="w-52" /><Input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="友链描述" className="w-52" /><Button onClick={() => { setPage(1); void fetchList(); }}>查询</Button></div></div>

        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead className="w-12"><input type="checkbox" checked={allChecked} onChange={(e) => setSelectedIds(e.target.checked ? new Set(items.map((it) => it.id)) : new Set())} /></TableHead><TableHead className="w-[120px]">Logo</TableHead><TableHead>名称</TableHead><TableHead>链接</TableHead><TableHead>描述</TableHead><TableHead>创建时间</TableHead><TableHead className="text-right">操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => (<TableRow key={it.id} className="h-20"><TableCell><input type="checkbox" checked={selectedIds.has(it.id)} onChange={(e) => setSelectedIds((prev) => { const n = new Set(prev); if (e.target.checked) n.add(it.id); else n.delete(it.id); return n; })} /></TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.logo} alt={it.name} /></div></TableCell><TableCell>{it.name}</TableCell><TableCell className="max-w-[260px] truncate">{it.link || it.url || '-'}</TableCell><TableCell className="max-w-[320px] truncate">{it.description}</TableCell><TableCell>{it.created_at ? formatDate(it.created_at) : '-'}</TableCell><TableCell className="text-right space-x-2"><Button size="sm" variant="outline" onClick={() => openEdit(it)}>更新</Button><Button size="sm" variant="destructive" onClick={() => onDelete([it.id])} disabled={submitting}>删除</Button></TableCell></TableRow>))}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
          </>
        )}

        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? '更新友链' : '新建友链'}</DialogTitle></DialogHeader>
            <div className="space-y-3"><Input placeholder="名称" value={form.name} onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} /><Input placeholder="链接" value={form.link} onChange={(e) => setForm((p) => ({ ...p, link: e.target.value }))} /><Input placeholder="Logo URL" value={form.logo} onChange={(e) => setForm((p) => ({ ...p, logo: e.target.value }))} /><Input placeholder="描述" value={form.description} onChange={(e) => setForm((p) => ({ ...p, description: e.target.value }))} /></div>
            <DialogFooter><Button variant="outline" onClick={() => setOpen(false)}>取消</Button><Button onClick={submit} disabled={submitting}>{editing ? '更新' : '创建'}</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
