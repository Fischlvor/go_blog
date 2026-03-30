'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { adminListEmojiSprites, adminRegenerateEmojiSprites, type AdminEmojiSprite } from '@/lib/api/admin/emoji';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';

export default function AdminEmojiSpritesPage() {
  const [items, setItems] = useState<AdminEmojiSprite[]>([]);
  const [detailOpen, setDetailOpen] = useState(false);
  const [current, setCurrent] = useState<AdminEmojiSprite | null>(null);

  const fetchList = async () => {
    try {
      const res = await adminListEmojiSprites();
      setItems(res || []);
    } catch {
      setItems([]);
      toast.error('获取雪碧图列表失败');
    }
  };

  useEffect(() => {
    adminListEmojiSprites()
      .then((res) => setItems(res || []))
      .catch(() => {
        setItems([]);
        toast.error('获取雪碧图列表失败');
      });
  }, []);

  const totalEmojis = useMemo(() => items.reduce((sum, it) => sum + (it.emoji_count || 0), 0), [items]);
  const totalSize = useMemo(() => items.reduce((sum, it) => sum + (it.file_size || 0), 0), [items]);

  const onRegenerate = async () => {
    try {
      await adminRegenerateEmojiSprites();
      toast.success('已触发雪碧图重生成');
      await fetchList();
    } catch {
      toast.error('操作失败');
    }
  };

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between"><CardTitle>雪碧图管理</CardTitle><div className="space-x-2"><Button onClick={onRegenerate}>重新生成全部雪碧图</Button><Button variant="outline" onClick={() => fetchList()}>刷新列表</Button></div></CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-3 md:grid-cols-3 text-sm"><div className="rounded border p-3">雪碧图总数：{items.length}</div><div className="rounded border p-3">表情总数：{totalEmojis}</div><div className="rounded border p-3">总文件大小：{totalSize}</div></div>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>组号</TableHead><TableHead>文件名</TableHead><TableHead className="w-[120px]">预览</TableHead><TableHead>尺寸</TableHead><TableHead>表情数量</TableHead><TableHead>文件大小</TableHead><TableHead>状态</TableHead><TableHead>生成时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (<TableRow key={it.id} className="h-20"><TableCell>{it.id}</TableCell><TableCell>{it.sprite_group}</TableCell><TableCell>{it.filename}</TableCell><TableCell className="w-[120px] p-0 align-middle"><div className="h-16 w-24"><ImageThumbCell src={it.cdn_url} alt={it.filename} boxWidth={96} boxHeight={64} /></div></TableCell><TableCell>{it.width} × {it.height}</TableCell><TableCell>{it.emoji_count}</TableCell><TableCell>{it.file_size}</TableCell><TableCell>{it.status === 1 ? '正常' : '异常'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => { setCurrent(it); setDetailOpen(true); }}>查看详情</Button>{it.cdn_url ? <a href={it.cdn_url} target="_blank" rel="noreferrer"><Button size="sm">下载</Button></a> : null}</TableCell></TableRow>))}</TableBody>
        </Table>

        <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>雪碧图详情</DialogTitle></DialogHeader>
            {current ? <div className="space-y-2 text-sm"><p>组号：{current.sprite_group}</p><p>文件名：{current.filename}</p><p>尺寸：{current.width} × {current.height}</p><p>表情数量：{current.emoji_count}</p><p>文件大小：{current.file_size}</p><p>生成时间：{formatDate(current.created_at)}</p><ImageThumbCell src={current.cdn_url} alt={current.filename} /></div> : null}
            <DialogFooter><Button variant="outline" onClick={() => setDetailOpen(false)}>关闭</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
