'use client';

export const dynamic = 'force-dynamic';

import { Suspense } from 'react';
import { useEffect, useRef, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import Image from 'next/image';
import { toast } from 'sonner';
import type { MDXEditorMethods } from '@mdxeditor/editor';
import { adminCreateArticle, adminUpdateArticle, adminSaveDraft, adminCreateCategory, adminCreateTag, adminGetArticle } from '@/lib/api/admin/article';
import { listCategories, listTags } from '@/lib/api/public/article';
import { adminUploadResourceChunk, adminCheckResource, adminInitResource, adminCompleteResource, adminCancelResource } from '@/lib/api/admin/resource';
import type { ArticleCategory, ArticleTag } from '@/lib/api/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { MDXEditorWrapper } from '@/components/admin/editor/MDXEditorWrapper';
import { MarkdownContent } from '@/components/site/article/MarkdownContent';
import { ImageThumbCell } from '@/components/admin/table-cells/ImageThumbCell';
import SparkMD5 from 'spark-md5';

function AdminArticlePublishContent() {
  const searchParams = useSearchParams();
  const editingSlug = searchParams.get('slug') || '';
  const editorRef = useRef<MDXEditorMethods>(null);

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [excerpt, setExcerpt] = useState('');
  const [featuredImage, setFeaturedImage] = useState('');
  const [status, setStatus] = useState<'draft' | 'published'>('published');
  const [visibility, setVisibility] = useState<'public' | 'private'>('public');
  const [categories, setCategories] = useState<ArticleCategory[]>([]);
  const [tags, setTags] = useState<ArticleTag[]>([]);
  const [categoryId, setCategoryId] = useState<number>(0);
  const [tagIds, setTagIds] = useState<number[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [autoSave, setAutoSave] = useState(true);
  const [publishOpen, setPublishOpen] = useState(false);
  const [draftSlug, setDraftSlug] = useState<string>('');
  const [lastSaveTime, setLastSaveTime] = useState<string>('');
  const [showNewCategoryDialog, setShowNewCategoryDialog] = useState(false);
  const [newCategoryForm, setNewCategoryForm] = useState({ name: '', slug: '' });
  const [showNewTagDialog, setShowNewTagDialog] = useState(false);
  const [newTagForm, setNewTagForm] = useState({ name: '', slug: '' });
  const autoSaveTimerRef = useRef<NodeJS.Timeout | null>(null);
  const performAutoSaveRef = useRef<(() => Promise<void>) | null>(null);
  const draftSlugRef = useRef<string>('');
  const lastSavedHashRef = useRef<string>('');
  const [uploadingImage, setUploadingImage] = useState(false);

  const HASH_CHUNK_SIZE = 2 * 1024 * 1024;

  const calcFileHash = async (file: File | Blob): Promise<string> => {
    return new Promise((resolve, reject) => {
      const chunkSize = HASH_CHUNK_SIZE;
      const chunks = Math.ceil(file.size / chunkSize);
      const spark = new SparkMD5.ArrayBuffer();
      const reader = new FileReader();
      let currentChunk = 0;

      reader.onload = (e) => {
        spark.append(e.target?.result as ArrayBuffer);
        currentChunk++;

        if (currentChunk < chunks) {
          loadNext();
        } else {
          resolve(spark.end());
        }
      };

      reader.onerror = () => {
        reject(new Error('文件读取失败'));
      };

      function loadNext() {
        const start = currentChunk * chunkSize;
        const end = Math.min(start + chunkSize, file.size);
        reader.readAsArrayBuffer(file.slice(start, end));
      }

      loadNext();
    });
  };

  const uploadImageChunks = async (
    taskId: string,
    file: File,
    chunkSize: number,
    missing: number[],
  ): Promise<void> => {
    let next = 0;

    const worker = async () => {
      while (next < missing.length) {
        const i = next;
        next += 1;

        const chunkNo = missing[i];
        const start = chunkNo * chunkSize;
        const end = Math.min(start + chunkSize, file.size);
        await adminUploadResourceChunk(taskId, chunkNo, file.slice(start, end));
      }
    };

    await Promise.all(Array.from({ length: Math.min(3, missing.length) }, () => worker()));
  };



  useEffect(() => {
    if (!editingSlug) return;
    adminGetArticle(editingSlug)
      .then((a) => {
        setTitle(a.title || '');
        setContent(a.content || '');
        setExcerpt(a.excerpt || '');
        setFeaturedImage(a.featured_image || '');
        setCategoryId(a.category?.id || 0);
        setTagIds((a.tags || []).map((t) => t.id));
        setStatus((a.status as 'draft' | 'published') || 'published');
        setVisibility((a.visibility as 'public' | 'private') || 'public');
        // 编辑器加载内容
        if (editorRef.current && a.content) {
          editorRef.current.setMarkdown(a.content);
        }
      })
      .catch(() => toast.error('加载文章详情失败'));
  }, [editingSlug]);

  // 自动保存到后端
  useEffect(() => {
    performAutoSaveRef.current = async () => {
      // draft 模式下只需要标题和内容，分类可选
      if (!title.trim() || !content.trim()) return;

      const currentData = {
        title: title.trim(),
        content,
        excerpt: excerpt.trim() || '',
        featuredImage: featuredImage.trim() || '',
        categoryId: categoryId || 0,
        tagIds,
      };

      // 计算当前数据的 hash
      const currentHash = await calcFileHash(new Blob([JSON.stringify(currentData)]));

      // 检查内容是否有变化
      if (lastSavedHashRef.current === currentHash) {
        return; // 内容没变，不发送请求
      }

      try {
        const payload = {
          title: currentData.title,
          content: currentData.content,
          excerpt: currentData.excerpt,
          featured_image: currentData.featuredImage,
          category_id: currentData.categoryId,
          tag_ids: currentData.tagIds,
        };

        if (draftSlugRef.current) {
          // 更新现有草稿
          await adminUpdateArticle({
            ...payload,
            slug: draftSlugRef.current,
            status: 'draft',
            visibility: 'private',
          });
        } else if (editingSlug) {
          // 编辑已发布的文章，保存为草稿
          await adminUpdateArticle({
            ...payload,
            slug: editingSlug,
            status: 'draft',
            visibility: 'private',
          });
        } else {
          // 创建新草稿
          const result = await adminSaveDraft(payload);
          draftSlugRef.current = result.slug;
          setDraftSlug(result.slug);
        }
        lastSavedHashRef.current = currentHash;
        setLastSaveTime(new Date().toLocaleTimeString('zh-CN'));
      } catch (error) {
        console.error('自动保存失败:', error);
      }
    };
  }, [title, content, excerpt, featuredImage, categoryId, tagIds]);

  // 打开发布配置对话框时加载分类和标签
  const handlePublishOpen = () => {
    if (categories.length === 0) {
      listCategories().then(setCategories).catch(() => setCategories([]));
    }
    if (tags.length === 0) {
      listTags().then(setTags).catch(() => setTags([]));
    }
    setPublishOpen(true);
  };

  // 创建新分类
  const handleCreateCategory = async () => {
    if (!newCategoryForm.name.trim() || !newCategoryForm.slug.trim()) {
      toast.error('请填写分类名称和标识');
      return;
    }
    try {
      const res = await adminCreateCategory(newCategoryForm);
      toast.success('分类创建成功');
      setShowNewCategoryDialog(false);
      setNewCategoryForm({ name: '', slug: '' });
      // 重新加载分类列表
      listCategories().then(setCategories).catch(() => setCategories([]));
    } catch {
      toast.error('创建分类失败');
    }
  };

  // 创建新标签
  const handleCreateTag = async () => {
    if (!newTagForm.name.trim() || !newTagForm.slug.trim()) {
      toast.error('请填写标签名称和标识');
      return;
    }
    try {
      const res = await adminCreateTag(newTagForm);
      toast.success('标签创建成功');
      setShowNewTagDialog(false);
      setNewTagForm({ name: '', slug: '' });
      // 重新加载标签列表
      listTags().then(setTags).catch(() => setTags([]));
    } catch {
      toast.error('创建标签失败');
    }
  };

  // 设置自动保存定时器
  useEffect(() => {
    if (!autoSave) {
      if (autoSaveTimerRef.current) clearInterval(autoSaveTimerRef.current);
      return;
    }

    if (autoSaveTimerRef.current) clearInterval(autoSaveTimerRef.current);
    autoSaveTimerRef.current = setInterval(() => {
      performAutoSaveRef.current?.();
    }, 30000); // 每 30 秒自动保存一次

    return () => {
      if (autoSaveTimerRef.current) clearInterval(autoSaveTimerRef.current);
    };
  }, [autoSave]);

  const tagContainerRef = useRef<HTMLDivElement>(null);
  const isScrollingRef = useRef(false);

  const handleTagWheel = (e: React.WheelEvent<HTMLDivElement>) => {
    if (tagContainerRef.current && !isScrollingRef.current) {
      e.preventDefault();
      isScrollingRef.current = true;
      const direction = e.deltaY > 0 ? 1 : -1;
      const container = tagContainerRef.current;
      container.scrollTo({
        top: container.scrollTop + direction * 41,
        behavior: 'smooth'
      });
      setTimeout(() => {
        isScrollingRef.current = false;
      }, 200);
    }
  };

  const onClear = () => {
    setTitle('');
    setContent('');
    setExcerpt('');
    setFeaturedImage('');
    setTagIds([]);
    setDraftSlug('');
    draftSlugRef.current = '';
    lastSavedHashRef.current = '';
    setLastSaveTime('');
  };

  const handleImageUpload = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      toast.error('请选择图片文件');
      return;
    }

    setUploadingImage(true);
    try {
      const hash = await calcFileHash(file);
      const checked = await adminCheckResource({
        file_hash: hash,
        file_size: file.size,
        file_name: file.name,
      });

      if (checked.exists) {
        setFeaturedImage(checked.file_url || '');
        toast.success('秒传成功');
        return;
      }

      let taskId = checked.task_id || '';
      let totalChunks = checked.total_chunks || 0;
      let chunkSize = 4 * 1024 * 1024;
      let missing = checked.missing_chunks || [];

      if (!taskId) {
        const init = await adminInitResource({
          file_hash: hash,
          file_size: file.size,
          file_name: file.name,
          mime_type: file.type || 'application/octet-stream',
        });
        taskId = init.task_id;
        totalChunks = init.total_chunks;
        chunkSize = init.chunk_size;
        missing = Array.from({ length: init.total_chunks }, (_, i) => i);
      }

      if (missing.length > 0) {
        await uploadImageChunks(taskId, file, chunkSize, missing);
      }

      const completed = await adminCompleteResource(taskId);
      setFeaturedImage(completed.file_url);
      toast.success('上传成功');
    } catch (error) {
      toast.error('上传失败');
      console.error('图片上传失败:', error);
    } finally {
      setUploadingImage(false);
    }
  };

  const openPublish = () => {
    if (!title.trim() || !content.trim()) {
      toast.error('请先填写标题和正文');
      return;
    }
    handlePublishOpen();
  };

  const submitPublish = async () => {
    if (!title.trim() || !content.trim()) {
      toast.error('请填写标题和正文');
      return;
    }
    // 发布时必须选择分类
    if (status === 'published' && !categoryId) {
      toast.error('发布文章必须选择分类');
      return;
    }

    setSubmitting(true);
    try {
      const basePayload = {
        title: title.trim(),
        content,
        excerpt: excerpt.trim() || '',
        featured_image: featuredImage.trim() || '',
        category_id: categoryId,
        tag_ids: tagIds,
        status,
        visibility,
      };
      if (editingSlug) {
        // 更新时需要 slug
        await adminUpdateArticle({ ...basePayload, slug: editingSlug });
        toast.success('更新成功');
      } else {
        // 创建时不传 slug，后端自动生成
        await adminCreateArticle(basePayload);
        toast.success('发布成功');
        onClear();
      }
      setPublishOpen(false);
    } catch {
      toast.error(editingSlug ? '更新失败' : '发布失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card>
      <CardHeader className="flex-row items-center">
        <CardTitle>{editingSlug ? '更新文章' : '发布文章'}</CardTitle>
        <div className="flex items-center gap-4 ml-auto">
          <div className="flex items-center gap-2 w-48">
            <span className="text-sm text-muted-foreground">自动保存</span>
            <input type="checkbox" checked={autoSave} onChange={(e) => setAutoSave(e.target.checked)} />
            {lastSaveTime && <span className="text-xs text-muted-foreground">最后保存: {lastSaveTime}</span>}
          </div>
          <div className="flex items-center gap-2">
            <Button variant="destructive" onClick={onClear}>清空文章</Button>
            <Button onClick={openPublish}>{editingSlug ? '更新' : '发布'}</Button>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="文章标题" />

        <div className="grid gap-3 lg:grid-cols-2 h-[700px]">
          <div className="space-y-2 flex flex-col">
            <p className="text-sm font-medium">Markdown 编辑器</p>
            <div className="rounded-md border border-border overflow-hidden flex-1">
              <MDXEditorWrapper
                ref={editorRef}
                markdown={content}
                onChange={(md) => setContent(md)}
                onImageUpload={async (file: File) => {
                  try {
                    const hash = await calcFileHash(file);
                    const checked = await adminCheckResource({
                      file_hash: hash,
                      file_size: file.size,
                      file_name: file.name,
                    });

                    if (checked.exists) {
                      return checked.file_url || '';
                    }

                    let taskId = checked.task_id || '';
                    let totalChunks = checked.total_chunks || 0;
                    let chunkSize = 4 * 1024 * 1024;
                    let missing = checked.missing_chunks || [];

                    if (!taskId) {
                      const init = await adminInitResource({
                        file_hash: hash,
                        file_size: file.size,
                        file_name: file.name,
                        mime_type: file.type || 'application/octet-stream',
                      });
                      taskId = init.task_id;
                      totalChunks = init.total_chunks;
                      chunkSize = init.chunk_size;
                      missing = Array.from({ length: init.total_chunks }, (_, i) => i);
                    }

                    if (missing.length > 0) {
                      await uploadImageChunks(taskId, file, chunkSize, missing);
                    }

                    const completed = await adminCompleteResource(taskId);
                    return completed.file_url;
                  } catch (error) {
                    console.error('编辑器图片上传失败:', error);
                    throw error;
                  }
                }}
              />
            </div>
          </div>

          <div className="space-y-2 flex flex-col">
            <p className="text-sm font-medium">渲染预览</p>
            <div className="overflow-auto rounded-md border border-border p-4 bg-muted/30 flex-1">
              <MarkdownContent content={content || '*(暂无内容)*'} />
            </div>
          </div>
        </div>
      </CardContent>

      <Dialog open={publishOpen} onOpenChange={setPublishOpen}>
        <DialogContent className="w-[500px] !max-w-none">
          <DialogHeader><DialogTitle>{editingSlug ? '更新文章配置' : '发布文章配置'}</DialogTitle></DialogHeader>
          <div className="space-y-3 max-h-[600px] overflow-y-auto pr-4">
            {/* 文章封面 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章封面</p>
              <div className="flex gap-4">
                <div 
                  className="border-2 border-dashed border-border rounded-lg p-1 text-center hover:bg-muted/50 transition h-28 flex flex-col items-center justify-center overflow-hidden flex-1 cursor-pointer"
                  onClick={() => document.getElementById('cover-upload')?.click()}
                >
                  {featuredImage ? (
                    <div className="relative flex h-full w-full max-h-full max-w-full items-center justify-center group">
                      <ImageThumbCell src={featuredImage} alt="cover" responsive={true} />
                      <p className="absolute text-xs text-muted-foreground bg-background/80 px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition">点击重新上传</p>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      <div className="text-2xl">📤</div>
                      <p className="text-sm">拖拽或点击上传</p>
                      <Button size="sm" variant="outline" disabled={uploadingImage}>{uploadingImage ? '上传中...' : '选择文件'}</Button>
                    </div>
                  )}
                  <Input 
                    type="file" 
                    accept="image/*" 
                    onChange={(e) => { 
                      const file = e.target.files?.[0]; 
                      if (file) handleImageUpload(file);
                      // 重置 input，允许重复选择同一文件
                      e.target.value = '';
                    }} 
                    className="hidden" 
                    id="cover-upload" 
                  />
                </div>
                <div className="flex flex-col justify-center w-1/4">
                  <p className="text-xs text-muted-foreground break-words">支持 jpg/ png/ jpeg/ gif/ webp，不超过 20MB</p>
                </div>
              </div>
              <Input value={featuredImage} disabled placeholder="文章封面 URL" className="text-xs" />
            </div>

            {/* 文章标题 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章标题</p>
              <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="请输入文章标题" />
            </div>

            {/* 分类 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章分类 {status === 'published' && <span className="text-red-500">*</span>}</p>
              <div className="flex gap-2">
                <select value={categoryId || ''} onChange={(e) => setCategoryId(Number(e.target.value))} className="flex-1 px-3 py-2 border border-border rounded-md text-sm bg-background hover:bg-muted/50 transition focus:outline-none focus:ring-2 focus:ring-primary/50">
                  <option value="">请选择分类</option>
                  {categories.map((c) => (<option key={c.id} value={c.id}>{c.name}</option>))}
                </select>
                <Button size="sm" variant="outline" onClick={() => setShowNewCategoryDialog(true)}>新建</Button>
              </div>
            </div>

            {/* 标签 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章标签</p>
              <div className="flex gap-2">
                <div ref={tagContainerRef} onWheel={handleTagWheel} className="flex-1 border border-border rounded-md h-10 overflow-hidden scrollbar-hide">
                  {tags.length > 0 ? (
                    <div className="flex flex-wrap content-start">
                      {tags.map((tag) => (
                        <label key={tag.id} className="flex items-center gap-2 text-sm cursor-pointer w-1/3 h-10 px-2.5">
                          <input type="checkbox" checked={tagIds.includes(tag.id)} onChange={() => { setTagIds((prev) => (prev.includes(tag.id) ? prev.filter((it) => it !== tag.id) : [...prev, tag.id])); }} className="rounded" />
                          <span className="truncate">{tag.name}</span>
                        </label>
                      ))}
                    </div>
                  ) : (
                    <p className="text-xs text-muted-foreground p-2">暂无标签</p>
                  )}
                </div>
                <Button size="sm" variant="outline" onClick={() => setShowNewTagDialog(true)}>新建</Button>
              </div>
            </div>

            {/* 文章简介 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章简介</p>
              <textarea value={excerpt} onChange={(e) => setExcerpt(e.target.value)} placeholder="请输入文章简介" className="w-full h-12 px-3 py-2 border border-border rounded-md text-sm bg-background focus:outline-none focus:ring-2 focus:ring-primary/50" />
            </div>

            {/* 可见性 */}
            <div className="space-y-2">
              <p className="text-sm font-medium">文章可见性</p>
              <div className="flex gap-6">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input type="radio" name="visibility" value="public" checked={visibility === 'public'} onChange={(e) => setVisibility(e.target.value as 'public' | 'private')} />
                  <span className="text-sm">公开</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input type="radio" name="visibility" value="private" checked={visibility === 'private'} onChange={(e) => setVisibility(e.target.value as 'public' | 'private')} />
                  <span className="text-sm">私密（仅自己可见）</span>
                </label>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPublishOpen(false)}>取消</Button>
            <Button onClick={submitPublish} disabled={submitting}>{submitting ? '提交中...' : (editingSlug ? '确认更新' : '确认')}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <Dialog open={showNewCategoryDialog} onOpenChange={setShowNewCategoryDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建分类</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">分类名称</label>
              <Input
                value={newCategoryForm.name}
                onChange={(e) => setNewCategoryForm({ ...newCategoryForm, name: e.target.value })}
                placeholder="请输入分类名称"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">分类标识</label>
              <Input
                value={newCategoryForm.slug}
                onChange={(e) => setNewCategoryForm({ ...newCategoryForm, slug: e.target.value })}
                placeholder="请输入分类标识（英文）"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowNewCategoryDialog(false)}>取消</Button>
            <Button onClick={handleCreateCategory}>确定</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 新建标签对话框 */}
      <Dialog open={showNewTagDialog} onOpenChange={setShowNewTagDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建标签</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">标签名称</label>
              <Input
                value={newTagForm.name}
                onChange={(e) => setNewTagForm({ ...newTagForm, name: e.target.value })}
                placeholder="请输入标签名称"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">标签标识</label>
              <Input
                value={newTagForm.slug}
                onChange={(e) => setNewTagForm({ ...newTagForm, slug: e.target.value })}
                placeholder="请输入标签标识（英文）"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowNewTagDialog(false)}>取消</Button>
            <Button onClick={handleCreateTag}>确定</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
}

export default function AdminArticlePublishPage() {
  return (
    <Suspense fallback={<Card><CardContent className="pt-6"><p className="text-sm text-muted-foreground">加载中...</p></CardContent></Card>}>
      <AdminArticlePublishContent />
    </Suspense>
  );
}
