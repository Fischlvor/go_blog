'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import {
  adminCreateAIModel,
  adminDeleteAIModel,
  adminListAIModels,
  adminUpdateAIModel,
  type AIModel,
} from '@/lib/api/admin/ai-management';
import { formatDate } from '@/lib/date';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminTablePagination } from '@/components/admin/table-cells/AdminTablePagination';

const DEFAULT_FORM = { name: '', display_name: '', provider: '', endpoint: '', api_key: '', max_tokens: 4096, temperature: 0.7, is_active: true };

export default function AdminAIModelsPage() {
  const [items, setItems] = useState<AIModel[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [name, setName] = useState('');
  const [provider, setProvider] = useState('');
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<AIModel | null>(null);
  const [form, setForm] = useState(DEFAULT_FORM);

  const fetchList = async () => {
    try {
      const res = await adminListAIModels({ page, page_size: pageSize, name: name || undefined, provider: provider || undefined });
      setItems(res.list || []);
      setTotal(res.total_items || 0);
    } catch {
      setItems([]);
      setTotal(0);
      toast.error('获取AI模型列表失败');
    }
  };

  useEffect(() => {
    adminListAIModels({ page, page_size: pageSize, name: name || undefined, provider: provider || undefined })
      .then((res) => {
        setItems(res.list || []);
        setTotal(res.total_items || 0);
      })
      .catch(() => {
        setItems([]);
        setTotal(0);
        toast.error('获取AI模型列表失败');
      });
  }, [page, pageSize, name, provider]);

  const onSearch = () => {
    setPage(1);
    void fetchList();
  };

  const onAdd = () => {
    setEditing(null);
    setForm(DEFAULT_FORM);
    setOpen(true);
  };

  const onEdit = (item: AIModel) => {
    setEditing(item);
    setForm({
      name: item.name || '',
      display_name: item.display_name || '',
      provider: item.provider || '',
      endpoint: item.endpoint || '',
      api_key: item.api_key || '',
      max_tokens: item.max_tokens || 4096,
      temperature: item.temperature || 0.7,
      is_active: item.is_active ?? true,
    });
    setOpen(true);
  };

  const onSubmit = async () => {
    if (!form.name.trim() || !form.display_name.trim() || !form.provider.trim()) {
      toast.error('请填写模型名称、显示名称和提供商');
      return;
    }

    try {
      if (editing) {
        await adminUpdateAIModel({ id: editing.id, ...form });
        toast.success('更新成功');
      } else {
        await adminCreateAIModel(form);
        toast.success('创建成功');
      }
      setOpen(false);
      await fetchList();
    } catch {
      toast.error('操作失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await adminDeleteAIModel(id);
      toast.success('删除成功');
      await fetchList();
    } catch {
      toast.error('删除失败');
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>AI模型管理</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2"><div><Button onClick={onAdd}>新增模型</Button></div><div className="flex items-center gap-2"><Input value={name} onChange={(e) => setName(e.target.value)} placeholder="模型名称" className="w-52" /><Input value={provider} onChange={(e) => setProvider(e.target.value)} placeholder="提供商" className="w-52" /><Button onClick={onSearch}>搜索</Button></div></div>

        <Table>
          <TableHeader><TableRow><TableHead>ID</TableHead><TableHead>模型名称</TableHead><TableHead>显示名称</TableHead><TableHead>提供商</TableHead><TableHead>端点</TableHead><TableHead>最大Token数</TableHead><TableHead>温度参数</TableHead><TableHead>状态</TableHead><TableHead>创建时间</TableHead><TableHead>操作</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell>{it.id}</TableCell><TableCell>{it.name}</TableCell><TableCell>{it.display_name || '-'}</TableCell><TableCell>{it.provider || '-'}</TableCell><TableCell className="max-w-[220px] truncate">{it.endpoint || '-'}</TableCell><TableCell>{it.max_tokens || '-'}</TableCell><TableCell>{it.temperature ?? '-'}</TableCell><TableCell>{it.is_active ? '启用' : '禁用'}</TableCell><TableCell>{formatDate(it.created_at)}</TableCell><TableCell className="space-x-2"><Button size="sm" variant="outline" onClick={() => onEdit(it)}>编辑</Button><Button size="sm" variant="destructive" onClick={() => onDelete(it.id)}>删除</Button></TableCell></TableRow>
          ))}</TableBody>
        </Table>

        <AdminTablePagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPageSize(size); setPage(1); }} />

        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? '编辑模型' : '新增模型'}</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <Input placeholder="模型名称" value={form.name} onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} />
              <Input placeholder="显示名称" value={form.display_name} onChange={(e) => setForm((p) => ({ ...p, display_name: e.target.value }))} />
              <Input placeholder="提供商" value={form.provider} onChange={(e) => setForm((p) => ({ ...p, provider: e.target.value }))} />
              <Input placeholder="API端点" value={form.endpoint} onChange={(e) => setForm((p) => ({ ...p, endpoint: e.target.value }))} />
              <Input placeholder="API密钥" value={form.api_key} onChange={(e) => setForm((p) => ({ ...p, api_key: e.target.value }))} />
              <Input placeholder="最大Token数" value={String(form.max_tokens)} onChange={(e) => setForm((p) => ({ ...p, max_tokens: Number(e.target.value) || 0 }))} />
              <Input placeholder="温度参数" value={String(form.temperature)} onChange={(e) => setForm((p) => ({ ...p, temperature: Number(e.target.value) || 0 }))} />
              <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={form.is_active} onChange={(e) => setForm((p) => ({ ...p, is_active: e.target.checked }))} /> 启用</label>
            </div>
            <DialogFooter><Button variant="outline" onClick={() => setOpen(false)}>取消</Button><Button onClick={onSubmit}>确定</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  );
}
