'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { adminDeleteEmoji, adminListEmojiGroups, adminListEmojis, adminRestoreEmoji, adminUploadEmojis, type AdminEmoji, type AdminEmojiGroup } from '@/lib/api/admin/emoji';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';

export default function AdminEmojiListPage() {
  const [items, setItems] = useState<AdminEmoji[]>([]);
  const [groups, setGroups] = useState<AdminEmojiGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploadOpen, setUploadOpen] = useState(false);
  const [files, setFiles] = useState<File[]>([]);
  const [groupKey, setGroupKey] = useState('');
  const [keyword, setKeyword] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const fetchList = async () => {
    setLoading(true);
    try {
      const res = await adminListEmojis({ page, page_size: pageSize, keyword: keyword || undefined, group_key: groupKey || undefined });
      setItems(res.list || []);
      setTotal(res.total_items || 0);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取表情列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    adminListEmojiGroups().then(setGroups).catch(() => setGroups([]));
  }, []);

  useEffect(() => {
    void fetchList();
  }, [page, pageSize]);

  const onUpload = async () => {
    if (!groupKey || files.length === 0) {
      toast.error('请选择表情组并添加文件');
      return;
    }

    const formData = new FormData();
    formData.append('group_key', groupKey);
    files.forEach((f) => formData.append('files', f));

    try {
      await adminUploadEmojis(formData);
      toast.success('上传成功');
      setUploadOpen(false);
      setFiles([]);
      await fetchList();
    } catch {
      toast.error('上传失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await adminDeleteEmoji(id);
      toast.success('删除成功');
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  const onRestore = async (id: number) => {
    try {
      await adminRestoreEmoji(id);
      toast.success('恢复成功');
      await fetchList();
    } catch {
      toast.error('恢复失败');
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>表情列表</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <Button onClick={() => setUploadOpen(true)}>批量上传表情</Button>
          </div>
          <div className="flex items-center gap-2">
            <Input value={keyword} onChange={(e) => setKeyword(e.target.value)} placeholder="搜索表情键名或文件名" className="w-64" />
            <select className="h-10 rounded-md border border-border bg-background px-3" value={groupKey} onChange={(e) => setGroupKey(e.target.value)}>
              <option value="">全部表情组</option>
              {groups.map((g) => <option key={g.id} value={g.group_key}>{g.group_name}</option>)}
            </select>
            <Button onClick={() => { setPage(1); void fetchList(); }}>搜索</Button>
          </div>
        </div>

        {loading ? <p className="text-sm text-muted-foreground">加载中...</p> : (
          <>
            <Table>
              <TableHeader><TableRow><TableHead>ID</TableHead><TableHead className="w-[120px]">缩略图</TableHead><TableHead>表情键</TableHead><TableHead>文件名</TableHead><TableHead>表情组</TableHead><TableHead>雪碧图信息</TableHead><TableHead>文件大小</TableHead><TableHead>上传时间</TableHead><TableHead>状态</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
              <TableBody>{items.map((it) => (
                <TableRow key={it.id} className="h-20"><TableCell>{it.id}</TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.cdn_url} alt={it.key} boxWidth={96} boxHeight={64} /></div></TableCell><TableCell>{it.key}</TableCell><TableCell>{it.filename}</TableCell><TableCell>{it.group_name}</TableCell><TableCell>{`组号:${it.sprite_group} 位置:(${it.sprite_position_x},${it.sprite_position_y})`}</TableCell><TableCell>{it.file_size}</TableCell><TableCell>{formatDate(it.upload_time || it.created_at)}</TableCell><TableCell>{it.status === 1 ? '正常' : '已删除'}</TableCell><TableCell>{it.status === 1 ? <Button size="sm" variant="destructive" onClick={() => onDelete(it.id)}>删除</Button> : <Button size="sm" onClick={() => onRestore(it.id)}>恢复</Button>}</TableCell></TableRow>
              ))}</TableBody>
            </Table>

            <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />
          </>
        )}

        <Dialog open={uploadOpen} onOpenChange={setUploadOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>批量上传表情</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <select className="h-10 w-full rounded-md border border-border bg-background px-3" value={groupKey} onChange={(e) => setGroupKey(e.target.value)}>
                <option value="">选择表情组</option>
                {groups.map((g) => <option key={g.id} value={g.group_key}>{g.group_name}</option>)}
              </select>
              <Input type="file" multiple accept=".png,.jpg,.jpeg,.gif,.webp" onChange={(e) => setFiles(Array.from(e.target.files || []))} />
              <p className="text-xs text-muted-foreground">已选择 {files.length} 个文件</p>
            </div>
            <DialogFooter><Button variant="outline" onClick={() => setUploadOpen(false)}>取消</Button><Button onClick={onUpload}>上传</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
