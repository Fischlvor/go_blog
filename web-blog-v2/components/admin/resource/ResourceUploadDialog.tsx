'use client';

import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import SparkMD5 from 'spark-md5';
import {
  adminCancelResource,
  adminCheckResource,
  adminCompleteResource,
  adminGetResourceMaxSize,
  adminInitResource,
  adminUploadResourceChunk,
} from '@/lib/api/admin/resource';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';

const HASH_CHUNK_SIZE = 2 * 1024 * 1024;

type UploadStatus = 'idle' | 'hashing' | 'checking' | 'instant' | 'uploading' | 'merging' | 'success' | 'error' | 'cancelled';

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  const kb = bytes / 1024;
  if (kb < 1024) return `${kb.toFixed(kb < 10 ? 1 : 0)} KB`;
  const mb = kb / 1024;
  if (mb < 1024) return `${mb.toFixed(mb < 10 ? 1 : 0)} MB`;
  const gb = mb / 1024;
  return `${gb.toFixed(gb < 10 ? 1 : 0)} GB`;
}

function statusLabel(status: UploadStatus) {
  if (status === 'hashing') return '正在计算文件指纹...';
  if (status === 'checking') return '正在检查文件（秒传/续传）...';
  if (status === 'instant') return '秒传成功';
  if (status === 'uploading') return '正在分片上传...';
  if (status === 'merging') return '正在合并文件...';
  if (status === 'success') return '上传成功';
  if (status === 'cancelled') return '上传已取消';
  if (status === 'error') return '上传失败';
  return '等待上传';
}

async function calcFileHash(file: File): Promise<string> {
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
}

interface ResourceUploadDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => Promise<void> | void;
}

export function ResourceUploadDialog({ open, onOpenChange, onSuccess }: ResourceUploadDialogProps) {
  const [file, setFile] = useState<File | null>(null);
  const [maxSize, setMaxSize] = useState(500 * 1024 * 1024);
  const [status, setStatus] = useState<UploadStatus>('idle');
  const [progress, setProgress] = useState(0);
  const [uploadedChunks, setUploadedChunks] = useState(0);
  const [totalChunks, setTotalChunks] = useState(0);
  const [taskId, setTaskId] = useState('');
  const [resultUrl, setResultUrl] = useState('');
  const [error, setError] = useState('');
  const [abortController, setAbortController] = useState<AbortController | null>(null);

  const isUploading = useMemo(() => ['hashing', 'checking', 'uploading', 'merging'].includes(status), [status]);

  useEffect(() => {
    adminGetResourceMaxSize()
      .then((res) => setMaxSize(res.max_size || 500 * 1024 * 1024))
      .catch(() => setMaxSize(500 * 1024 * 1024));
  }, []);

  const canStart = useMemo(() => !!file && ['idle', 'error', 'cancelled'].includes(status), [file, status]);
  const canReset = useMemo(() => ['success', 'instant', 'error', 'cancelled'].includes(status), [status]);

  const reset = () => {
    setFile(null);
    setStatus('idle');
    setProgress(0);
    setUploadedChunks(0);
    setTotalChunks(0);
    setTaskId('');
    setResultUrl('');
    setError('');
    setAbortController(null);
  };

  const uploadMissing = async (
    currentTaskId: string,
    currentFile: File,
    chunkSize: number,
    missing: number[],
    signal: AbortSignal,
    doneBefore: number,
    totalCount: number,
  ) => {
    let next = 0;
    let done = 0;

    const worker = async () => {
      while (next < missing.length) {
        const i = next;
        next += 1;

        const chunkNo = missing[i];
        const start = chunkNo * chunkSize;
        const end = Math.min(start + chunkSize, currentFile.size);
        await adminUploadResourceChunk(currentTaskId, chunkNo, currentFile.slice(start, end), signal);

        done += 1;
        const uploaded = doneBefore + done;
        setUploadedChunks(uploaded);
        setProgress(Math.min(90, 10 + Math.round((uploaded / totalCount) * 80)));
      }
    };

    await Promise.all(Array.from({ length: Math.min(3, missing.length) }, () => worker()));
  };

  const startUpload = async () => {
    if (!file) return;
    if (file.size > maxSize) {
      toast.error(`文件大小超限，最大 ${formatFileSize(maxSize)}，当前 ${formatFileSize(file.size)}`);
      return;
    }

    const controller = new AbortController();
    setAbortController(controller);
    setError('');
    setResultUrl('');

    try {
      setStatus('hashing');
      setProgress(5);
      const hash = await calcFileHash(file);

      setStatus('checking');
      setProgress(10);
      const checked = await adminCheckResource({
        file_hash: hash,
        file_size: file.size,
        file_name: file.name,
      });

      if (checked.exists) {
        setStatus('instant');
        setProgress(100);
        setResultUrl(checked.file_url || '');
        toast.success('秒传成功');
        await onSuccess();
        return;
      }

      let currentTaskId = checked.task_id || '';
      let currentTotalChunks = checked.total_chunks || 0;
      let chunkSize = 4 * 1024 * 1024;
      let missing = checked.missing_chunks || [];

      if (!currentTaskId) {
        const init = await adminInitResource({
          file_hash: hash,
          file_size: file.size,
          file_name: file.name,
          mime_type: file.type || 'application/octet-stream',
        });
        currentTaskId = init.task_id;
        currentTotalChunks = init.total_chunks;
        chunkSize = init.chunk_size;
        missing = Array.from({ length: init.total_chunks }, (_, i) => i);
      }

      setTaskId(currentTaskId);
      setTotalChunks(currentTotalChunks);
      const doneBefore = Math.max(0, currentTotalChunks - missing.length);
      setUploadedChunks(doneBefore);

      if (missing.length > 0) {
        setStatus('uploading');
        await uploadMissing(currentTaskId, file, chunkSize, missing, controller.signal, doneBefore, currentTotalChunks);
      }

      setStatus('merging');
      setProgress(95);
      const completed = await adminCompleteResource(currentTaskId);
      setStatus('success');
      setProgress(100);
      setResultUrl(completed.file_url);
      toast.success('上传成功');
      await onSuccess();
    } catch (e) {
      if ((e as Error).name === 'AbortError') {
        setStatus('cancelled');
        setError('上传已取消');
        return;
      }
      setStatus('error');
      setError((e as Error).message || '上传失败');
    } finally {
      setAbortController(null);
    }
  };

  const cancelUpload = async () => {
    abortController?.abort();
    if (taskId) {
      try {
        await adminCancelResource(taskId);
      } catch {
        // ignore
      }
    }
    setStatus('cancelled');
    setError('上传已取消');
  };

  const handleFileChange = (nextFile: File | null) => {
    setFile(nextFile);
    setStatus('idle');
    setProgress(0);
    setUploadedChunks(0);
    setTotalChunks(0);
    setTaskId('');
    setResultUrl('');
    setError('');
  };

  const handleDialogOpenChange = async (next: boolean) => {
    if (next) {
      onOpenChange(true);
      return;
    }

    if (isUploading) {
      const confirmed = window.confirm('关闭将取消当前上传任务，是否继续？');
      if (!confirmed) return;
      await cancelUpload();
    }

    onOpenChange(false);
    reset();
  };

  return (
    <Dialog open={open} onOpenChange={(next) => { void handleDialogOpenChange(next); }}>
      <DialogContent>
        <DialogHeader><DialogTitle>上传资源</DialogTitle></DialogHeader>

        <div className="space-y-3">
          {!file ? (
            <>
              <Input type="file" onChange={(e) => handleFileChange(e.target.files?.[0] || null)} />
              <p className="text-xs text-muted-foreground">支持图片、视频、音频、文档等格式，单文件最大 {formatFileSize(maxSize)}</p>
            </>
          ) : (
            <>
              <div className="rounded border p-3 text-sm">
                <div className="font-medium">{file.name}</div>
                <div className="text-muted-foreground">{formatFileSize(file.size)}</div>
              </div>

              <div className="space-y-2">
                <div className="text-sm">
                  {statusLabel(status)} {status === 'uploading' ? `${uploadedChunks}/${totalChunks}` : ''}
                </div>
                <div className="h-2 overflow-hidden rounded bg-muted">
                  <div className="h-full bg-primary transition-all" style={{ width: `${progress}%` }} />
                </div>
                <div className="text-xs text-muted-foreground">{progress}%</div>
              </div>

              {error ? <p className="text-sm text-destructive">{error}</p> : null}

              {resultUrl ? (
                <div className="space-y-2">
                  <Input value={resultUrl} readOnly />
                  <Button size="sm" variant="outline" onClick={() => navigator.clipboard.writeText(resultUrl)}>复制链接</Button>
                </div>
              ) : null}
            </>
          )}
        </div>

        <DialogFooter>
          {isUploading ? <Button variant="destructive" onClick={cancelUpload}>取消上传</Button> : null}
          {file ? <Button variant="outline" onClick={() => handleFileChange(null)} disabled={isUploading}>重新选择文件</Button> : null}
          {!file ? <Button variant="outline" onClick={() => { void handleDialogOpenChange(false); }}>关闭</Button> : null}
          {canStart ? <Button onClick={startUpload}>开始上传</Button> : null}
          {['success', 'instant'].includes(status) ? <Button onClick={() => { void handleDialogOpenChange(false); }}>完成</Button> : null}
          {canReset ? <Button variant="outline" onClick={reset}>继续上传</Button> : null}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
