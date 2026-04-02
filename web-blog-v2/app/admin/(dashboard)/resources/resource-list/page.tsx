'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useState } from 'react';
import { FileArchive, FileAudio2, FileCode2, FileText, FileVideo2, FileWarning, Image as ImageIcon } from 'lucide-react';
import { toast } from 'sonner';
import { adminDeleteResources, adminListResources } from '@/lib/api/admin/resource';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ResourceUploadDialog } from '@/components/admin/resource/ResourceUploadDialog';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';
import type { AdminResource } from '@/lib/api/admin/resource';

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  const kb = bytes / 1024;
  if (kb < 1024) return `${kb.toFixed(kb < 10 ? 1 : 0)} KB`;
  const mb = kb / 1024;
  if (mb < 1024) return `${mb.toFixed(mb < 10 ? 1 : 0)} MB`;
  const gb = mb / 1024;
  return `${gb.toFixed(gb < 10 ? 1 : 0)} GB`;
}

function ResourceTypeIcon({ mime }: { mime: string }) {
  if (mime.startsWith('image/')) return <ImageIcon className="h-full w-full text-emerald-600" />;
  if (mime.startsWith('video/')) return <FileVideo2 className="h-full w-full text-violet-600" />;
  if (mime.startsWith('audio/')) return <FileAudio2 className="h-full w-full text-orange-600" />;
  if (mime.includes('pdf') || mime.includes('word') || mime.includes('text/')) return <FileText className="h-full w-full text-blue-600" />;
  if (mime.includes('zip') || mime.includes('rar') || mime.includes('tar') || mime.includes('7z')) return <FileArchive className="h-full w-full text-amber-600" />;
  if (mime.includes('json') || mime.includes('javascript') || mime.includes('typescript') || mime.includes('xml')) return <FileCode2 className="h-full w-full text-cyan-600" />;
  return <FileWarning className="h-full w-full text-muted-foreground" />;
}

function ResourcePreview({ item }: { item: AdminResource }) {
  if (item.mime_type.startsWith('image/')) {
    return <ImageThumbCell src={item.file_url} alt={item.file_name} boxWidth={96} boxHeight={64} />;
  }

  if (item.mime_type.startsWith('video/') && item.transcode_status === 2 && item.thumbnail_url) {
    return <ImageThumbCell src={item.thumbnail_url} alt={item.file_name} boxWidth={96} boxHeight={64} />;
  }

  return (
    <div className="inline-flex items-center justify-center overflow-hidden rounded-lg border border-border bg-muted/20" style={{ width: 96, height: 64 }}>
      <div className="flex h-full w-full items-center justify-center p-3">
        <ResourceTypeIcon mime={item.mime_type} />
      </div>
    </div>
  );
}

export default function AdminResourceListPage() {
  const [items, setItems] = useState<AdminResource[]>([]);
  const [loading, setLoading] = useState(true);
  const [fileName, setFileName] = useState('');
  const [mimeType, setMimeType] = useState('');
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [uploadOpen, setUploadOpen] = useState(false);

  const fetchList = async () => {
    setLoading(true);
    setSelectedIds(new Set());
    try {
      const res = await adminListResources({ page, page_size: pageSize, file_name: fileName || null, mime_type: mimeType || null });
      setItems(res.list);
      setTotal(res.total_items);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('加载资源列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { void fetchList(); }, [page, pageSize]);

  const onDelete = async (ids: number[]) => {
    if (!ids.length) { toast.warning('请先选择要删除的资源'); return; }
    const isBatch = ids.length > 1;
    const confirmed = window.confirm(isBatch ? `确认删除选中的 ${ids.length} 个资源吗？` : '确认删除该资源吗？');
    if (!confirmed) return;

    try {
      await adminDeleteResources(ids);
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
        <CardTitle>资源列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <Button onClick={() => setUploadOpen(true)}>上传资源</Button>
            <Button variant="destructive" disabled={selectedIds.size === 0} onClick={() => onDelete(Array.from(selectedIds))}>批量删除</Button>
          </div>
          <div className="flex items-center gap-2">
            <Input value={fileName} onChange={(e) => setFileName(e.target.value)} placeholder="资源名称" className="w-56" />
            <Input value={mimeType} onChange={(e) => setMimeType(e.target.value)} placeholder="MIME类型" className="w-56" />
            <Button onClick={() => { setPage(1); setSelectedIds(new Set()); void fetchList(); }}>查询</Button>
          </div>
        </div>

        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead className="w-12"><input type="checkbox" checked={allChecked} onChange={(e) => setSelectedIds(e.target.checked ? new Set(items.map((it) => it.id)) : new Set())} /></TableHead><TableHead>ID</TableHead><TableHead className="w-[120px]">预览</TableHead><TableHead>文件名</TableHead><TableHead>MIME</TableHead><TableHead>状态</TableHead><TableHead>大小</TableHead><TableHead>URL</TableHead><TableHead>创建时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => { const status = it.transcode_status; const copyUrl = it.mime_type.startsWith('video/') && status === 2 ? (it.transcode_url || it.file_url) : it.file_url; return (<TableRow key={it.id} className="h-20"><TableCell><input type="checkbox" checked={selectedIds.has(it.id)} onChange={(e) => setSelectedIds((prev) => { const n = new Set(prev); if (e.target.checked) n.add(it.id); else n.delete(it.id); return n; })} /></TableCell><TableCell>{it.id}</TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ResourcePreview item={it} /></div></TableCell><TableCell className="max-w-[220px] truncate">{it.file_name}</TableCell><TableCell>{it.mime_type}</TableCell><TableCell>{status === 2 ? <Badge>已转码</Badge> : status === 1 ? <Badge variant="secondary">转码中</Badge> : status === 3 ? <Badge variant="destructive">转码失败</Badge> : <Badge variant="outline">-</Badge>}</TableCell><TableCell>{formatFileSize(it.file_size)}</TableCell><TableCell className="max-w-[320px] truncate">{it.file_url}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" disabled={it.mime_type.startsWith('video/') && status !== 2} onClick={() => navigator.clipboard.writeText(copyUrl)}>复制链接</Button><Button size="sm" variant="destructive" onClick={() => onDelete([it.id])}>删除</Button></TableCell></TableRow>); })}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={(nextPage) => { setSelectedIds(new Set()); setPage(nextPage); }} onPageSizeChange={(size) => { setSelectedIds(new Set()); setPageSize(size); setPage(1); }} />
          </>
        )}

        <ResourceUploadDialog open={uploadOpen} onOpenChange={setUploadOpen} onSuccess={fetchList} />
      </CardContent>
    </Card>
  );
}
