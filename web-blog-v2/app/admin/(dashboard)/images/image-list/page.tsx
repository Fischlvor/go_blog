'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { adminDeleteImages, adminListImages } from '@/lib/client-api/admin/image';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';
import type { AdminImage } from '@/lib/client-api/admin/image';

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  const kb = bytes / 1024;
  if (kb < 1024) return `${kb.toFixed(kb < 10 ? 1 : 0)} KB`;
  const mb = kb / 1024;
  if (mb < 1024) return `${mb.toFixed(mb < 10 ? 1 : 0)} MB`;
  const gb = mb / 1024;
  return `${gb.toFixed(gb < 10 ? 1 : 0)} GB`;
}

export default function AdminImageListPage() {
  const [items, setItems] = useState<AdminImage[]>([]);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState('');
  const [mimeType, setMimeType] = useState('');
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListImages({ page, page_size: pageSize, name: name || null, mime_type: mimeType || null });
      setItems(res.list);
      setTotal(res.total_items);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('加载图片列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { void fetchList(); }, [page, pageSize]);

  const onDelete = async (ids: number[]) => {
    if (!ids.length) { toast.warning('请先选择要删除的图片'); return; }
    try {
      await adminDeleteImages(ids);
      toast.success('删除成功');
      setSelectedIds(new Set());
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  const allChecked = useMemo(() => items.length > 0 && items.every((it) => selectedIds.has(it.id)), [items, selectedIds]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>图片列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <Button variant="destructive" onClick={() => onDelete(Array.from(selectedIds))}>批量删除</Button>
          </div>
          <div className="flex items-center gap-2">
            <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="图片名称" className="w-56" />
            <Input value={mimeType} onChange={(e) => setMimeType(e.target.value)} placeholder="MIME类型" className="w-56" />
            <Button onClick={() => { setPage(1); void fetchList(); }}>查询</Button>
          </div>
        </div>

        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead className="w-12"><input type="checkbox" checked={allChecked} onChange={(e) => setSelectedIds(e.target.checked ? new Set(items.map((it) => it.id)) : new Set())} /></TableHead><TableHead>ID</TableHead><TableHead className="w-[120px]">缩略图</TableHead><TableHead>名称</TableHead><TableHead>MIME</TableHead><TableHead>大小</TableHead><TableHead>URL</TableHead><TableHead>创建时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => (<TableRow key={it.id} className="h-20"><TableCell><input type="checkbox" checked={selectedIds.has(it.id)} onChange={(e) => setSelectedIds((prev) => { const n = new Set(prev); if (e.target.checked) n.add(it.id); else n.delete(it.id); return n; })} /></TableCell><TableCell>{it.id}</TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.url} alt={it.name} boxWidth={96} boxHeight={64} /></div></TableCell><TableCell>{it.name}</TableCell><TableCell>{it.mime_type}</TableCell><TableCell>{formatFileSize(it.size)}</TableCell><TableCell className="max-w-[360px] truncate">{it.url}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => navigator.clipboard.writeText(it.url)}>复制链接</Button><Button size="sm" variant="destructive" onClick={() => onDelete([it.id])}>删除</Button></TableCell></TableRow>))}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
          </>
        )}
      </CardContent>
    </Card>
  );
}
