'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { adminCreateEmojiGroup, adminDeleteEmojiGroup, adminListEmojiGroups, adminRegenerateEmojiSprites, adminUpdateEmojiGroup, type AdminEmojiGroup } from '@/lib/api/admin/emoji';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function AdminEmojiGroupsPage() {
  const [items, setItems] = useState<AdminEmojiGroup[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<AdminEmojiGroup | null>(null);
  const [groupName, setGroupName] = useState('');
  const [description, setDescription] = useState('');
  const [sortOrder, setSortOrder] = useState(0);

  const fetchList = async () => {
    try {
      const res = await adminListEmojiGroups();
      setItems(res || []);
    } catch {
      setItems([]);
      toast.error('获取表情组失败');
    }
  };

  useEffect(() => {
    adminListEmojiGroups()
      .then((res) => setItems(res || []))
      .catch(() => {
        setItems([]);
        toast.error('获取表情组失败');
      });
  }, []);

  const openCreate = () => {
    setEditing(null);
    setGroupName('');
    setDescription('');
    setSortOrder(0);
    setOpen(true);
  };

  const openEdit = (item: AdminEmojiGroup) => {
    setEditing(item);
    setGroupName(item.group_name || '');
    setDescription(item.description || '');
    setSortOrder(item.sort_order || 0);
    setOpen(true);
  };

  const onSubmit = async () => {
    if (!groupName.trim()) {
      toast.error('组名不能为空');
      return;
    }
    try {
      if (editing) {
        await adminUpdateEmojiGroup(editing.id, { group_name: groupName.trim(), description, sort_order: sortOrder });
        toast.success('更新成功');
      } else {
        await adminCreateEmojiGroup({ group_name: groupName.trim(), description, sort_order: sortOrder });
        toast.success('创建成功');
      }
      setOpen(false);
      await fetchList();
    } catch {
      toast.error('操作失败');
    }
  };

  const onRegenerate = async (groupKey?: string) => {
    try {
      await adminRegenerateEmojiSprites(groupKey ? [groupKey] : undefined);
      toast.success('已触发雪碧图重生成');
    } catch {
      toast.error('操作失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await adminDeleteEmojiGroup(id);
      toast.success('删除成功');
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between"><CardTitle>表情组管理</CardTitle><Button onClick={openCreate}>新建表情组</Button></CardHeader>
      <CardContent className="space-y-4">
        <div><Button variant="outline" onClick={() => onRegenerate()}>重新生成全部雪碧图</Button></div>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>组名</TableHead><TableHead>组Key</TableHead><TableHead>描述</TableHead><TableHead>排序</TableHead><TableHead>表情数量</TableHead><TableHead>状态</TableHead><TableHead>创建时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell>{it.id}</TableCell><TableCell>{it.group_name}</TableCell><TableCell>{it.group_key}</TableCell><TableCell>{it.description}</TableCell><TableCell>{it.sort_order}</TableCell><TableCell>{it.emoji_count}</TableCell><TableCell>{it.status === 1 ? '启用' : '禁用'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => openEdit(it)}>编辑</Button><Button size="sm" variant="outline" onClick={() => onRegenerate(it.group_key)} disabled={!it.emoji_count}>重生成</Button><Button size="sm" variant="destructive" onClick={() => onDelete(it.id)} disabled={it.emoji_count > 0}>删除</Button></TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? '编辑表情组' : '新建表情组'}</DialogTitle></DialogHeader>
            <div className="space-y-3"><Input value={groupName} onChange={(e) => setGroupName(e.target.value)} placeholder="组名" /><Input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="描述" /><Input value={String(sortOrder)} onChange={(e) => setSortOrder(Number(e.target.value) || 0)} placeholder="排序" /></div>
            <DialogFooter><Button variant="outline" onClick={() => setOpen(false)}>取消</Button><Button onClick={onSubmit}>{editing ? '更新' : '创建'}</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
