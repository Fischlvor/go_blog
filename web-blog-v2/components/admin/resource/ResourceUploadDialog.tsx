'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { toast } from 'sonner';
import {
  adminCancelResource,
  adminCheckResource,
  adminCompleteResource,
  adminGetResourceMaxSize,
  adminInitResource,
  adminUploadResourceChunk,
} from '@/lib/client-api/admin/resource';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';

const HASH_CHUNK_SIZE = 4 * 1024 * 1024; // 七牛云 qetag 使用 4MB 分块

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

function formatTime(seconds: number) {
  if (seconds < 60) return `${seconds}秒`;
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return secs > 0 ? `${minutes}分${secs}秒` : `${minutes}分`;
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

async function calcFileHash(file: File, signal: AbortSignal): Promise<string> {
  // 动态导入 qetag-js（避免 SSR 错误）
  const { File: QZFile, Normal: QETagNormal } = await import('qetag-js');

  return new Promise((resolve, reject) => {
    signal.addEventListener('abort', () => {
      reject(new Error('Hash calculation aborted'));
    });

    const qzFile = new QZFile({ file, blockSize: HASH_CHUNK_SIZE });
    const qetag = new QETagNormal(qzFile);

    qetag
      .get()
      .then((hash) => {
        if (signal.aborted) {
          reject(new Error('Hash calculation aborted'));
        } else {
          resolve(hash);
        }
      })
      .catch((err) => {
        reject(err);
      });
  });
}

function getMimeTypeFromExtension(filename: string): string {
  const ext = filename.toLowerCase().match(/\.([^.]+)$/)?.[1];
  const mimeMap: Record<string, string> = {
    '7z': 'application/x-7z-compressed',
    'rar': 'application/x-rar-compressed',
    'zip': 'application/zip',
    'tar': 'application/x-tar',
    'gz': 'application/gzip',
    'exe': 'application/x-msdownload',
    'pdf': 'application/pdf',
    'doc': 'application/msword',
    'docx': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'xls': 'application/vnd.ms-excel',
    'xlsx': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'txt': 'text/plain',
    'csv': 'text/csv',
  };
  return ext && mimeMap[ext] ? mimeMap[ext] : '';
}

function getFileMimeType(file: File): string {
  if (file.type && file.type !== 'application/octet-stream') {
    return file.type;
  }
  const mimeFromExt = getMimeTypeFromExtension(file.name);
  return mimeFromExt || file.type || 'application/octet-stream';
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

  // 新增：上传速度和时间统计
  const [uploadSpeed, setUploadSpeed] = useState(0);
  const [uploadedBytes, setUploadedBytes] = useState(0);
  const [estimatedTime, setEstimatedTime] = useState(0);
  const [activeTasks, setActiveTasks] = useState(0);
  const [concurrency, setConcurrency] = useState(3);
  const concurrencyRef = useRef(3);

  useEffect(() => {
    concurrencyRef.current = concurrency;
  }, [concurrency]);

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
    setUploadSpeed(0);
    setUploadedBytes(0);
    setEstimatedTime(0);
    setActiveTasks(0);
    setConcurrency(3);
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
    const state = {
      completed: 0,
      failed: false,
      error: null as Error | null,
      lastTime: Date.now(),
      lastBytes: 0,
      speedSamples: [] as number[],
    };

    const tasks = [...missing];
    const activeUploads = new Map<number, Promise<void>>();
    let nextTaskIndex = 0;

    const uploadChunk = async (chunkNo: number) => {
      const start = chunkNo * chunkSize;
      const end = Math.min(start + chunkSize, currentFile.size);

      await adminUploadResourceChunk(currentTaskId, chunkNo, currentFile.slice(start, end), signal);

      state.completed++;
      const uploaded = doneBefore + state.completed;
      const currentBytes = Math.min(uploaded * chunkSize, currentFile.size);

      setUploadedChunks(uploaded);
      setUploadedBytes(currentBytes);
      setProgress(Math.round((uploaded / totalCount) * 100));

      const now = Date.now();
      const elapsed = (now - state.lastTime) / 1000;
      if (elapsed >= 1) {
        const speed = (currentBytes - state.lastBytes) / elapsed;
        state.speedSamples.push(speed);
        if (state.speedSamples.length > 5) state.speedSamples.shift();

        const avgSpeed = state.speedSamples.reduce((a, b) => a + b, 0) / state.speedSamples.length;
        setUploadSpeed(Math.round(avgSpeed));

        const remaining = currentFile.size - currentBytes;
        if (avgSpeed > 0) setEstimatedTime(Math.round(remaining / avgSpeed));

        state.lastTime = now;
        state.lastBytes = currentBytes;
      }
    };

    const processNext = async (): Promise<void> => {
      while (!state.failed && nextTaskIndex < tasks.length) {
        while (activeUploads.size >= concurrencyRef.current) {
          await Promise.race(activeUploads.values());
        }

        if (state.failed) break;

        const chunkNo = tasks[nextTaskIndex++];
        const uploadPromise = uploadChunk(chunkNo)
          .catch(err => {
            if (!state.failed) {
              state.failed = true;
              state.error = err;
            }
          })
          .finally(() => {
            activeUploads.delete(chunkNo);
            setActiveTasks(activeUploads.size);
          });

        activeUploads.set(chunkNo, uploadPromise);
        setActiveTasks(activeUploads.size);
      }

      await Promise.all(activeUploads.values());
    };

    await processNext();

    if (state.failed && state.error) throw state.error;
  };

  const startUpload = async () => {
    if (!file) return;

    // 取消旧的上传（如果存在）
    if (abortController) {
      abortController.abort();
    }

    const controller = new AbortController();
    setAbortController(controller);
    setError('');
    setResultUrl('');

    try {
      // Hash 阶段
      setStatus('hashing');
      setProgress(0);
      const hash = await calcFileHash(file, controller.signal);

      // Checking 阶段
      setStatus('checking');
      setProgress(100);
      const mimeType = getFileMimeType(file);
      const checked = await adminCheckResource({
        file_hash: hash,
        file_size: file.size,
        file_name: file.name,
        mime_type: mimeType,
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
          mime_type: mimeType,
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
        setProgress(0); // 重置进度
        await uploadMissing(
          currentTaskId,
          file,
          chunkSize,
          missing,
          controller.signal,
          doneBefore,
          currentTotalChunks,
        );
      }

      setStatus('merging');
      setProgress(0);

      // 尝试完成上传，如果失败可能是有分片丢失
      let retryCount = 0;
      const maxRetries = 3;

      while (retryCount < maxRetries) {
        try {
          const completed = await adminCompleteResource(currentTaskId);
          setStatus('success');
          setProgress(100);
          setResultUrl(completed.file_url);
          toast.success('上传成功');
          await onSuccess();
          return;
        } catch (err) {
          const errorMessage = (err as Error).message || String(err);

          if (errorMessage.includes('任务已完成') || errorMessage.includes('already completed')) {
            setStatus('success');
            setProgress(100);
            toast.success('上传成功');
            await onSuccess();
            return;
          }

          const match = errorMessage.match(/\[([0-9\s]+)\]/);
          if (match) {
            const missingChunks = match[1].trim().split(/\s+/).map(Number);

            if (missingChunks.length > 0) {
              setStatus('uploading');

              await uploadMissing(
                currentTaskId,
                file,
                chunkSize,
                missingChunks,
                controller.signal,
                currentTotalChunks - missingChunks.length,
                currentTotalChunks,
              );

              retryCount++;
              setStatus('merging');
              continue;
            }
          }

          throw err;
        }
      }

      // 重试次数用尽
      throw new Error(`完成上传失败，已重试 ${maxRetries} 次`);
    } catch (e) {
      const errMsg = (e as Error).message;
      if ((e as Error).name === 'AbortError' || errMsg.includes('aborted')) {
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
    if (nextFile && nextFile.size > maxSize) {
      toast.error(`文件大小超限，最大 ${formatFileSize(maxSize)}，当前 ${formatFileSize(nextFile.size)}`);
      return;
    }
    if (nextFile) {
      const mimeType = getFileMimeType(nextFile);
      if (!mimeType) {
        toast.error(`无法识别文件类型：${mimeType}，请确保文件扩展名正确`);
        return;
      }
    }
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
                {/* 第1行：阶段名称 */}
                <div className="text-sm font-medium">
                  {statusLabel(status)}
                </div>

                {/* 第2行：进度条（hashing/checking/merging 不显示百分比） */}
                <div className="space-y-1">
                  <div className="h-2 overflow-hidden rounded-full bg-muted">
                    <div
                      className={`h-full bg-primary transition-all duration-300 ease-out ${
                        ['hashing', 'checking', 'merging'].includes(status) ? 'w-full animate-pulse' : ''
                      }`}
                      style={['hashing', 'checking', 'merging'].includes(status) ? undefined : { width: `${progress}%` }}
                    />
                  </div>
                  {!['hashing', 'checking', 'merging'].includes(status) && (
                    <div className="text-right text-xs text-muted-foreground">{progress}%</div>
                  )}
                </div>

                {/* 第3行：详细信息（仅上传阶段显示） */}
                {status === 'uploading' && uploadSpeed > 0 && (
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span>{formatFileSize(uploadedBytes)} / {formatFileSize(file.size)}</span>
                    <span className="flex items-center gap-2">
                      <span>{formatFileSize(uploadSpeed)}/s</span>
                      {estimatedTime > 0 && <span>• 剩余 {formatTime(estimatedTime)}</span>}
                    </span>
                  </div>
                )}

                {/* 并发控制（仅上传阶段显示） */}
                {status === 'uploading' && (
                  <div className="flex items-center gap-2 rounded border bg-muted/30 px-3 py-2 text-sm">
                    <span className="text-muted-foreground">并发数:</span>
                    <Button
                      size="sm"
                      variant="outline"
                      className="h-7 w-7 p-0"
                      onClick={() => setConcurrency(prev => Math.max(1, prev - 1))}
                      disabled={concurrency <= 1}
                    >
                      -
                    </Button>
                    <span className="w-8 text-center font-medium">{concurrency}</span>
                    <Button
                      size="sm"
                      variant="outline"
                      className="h-7 w-7 p-0"
                      onClick={() => setConcurrency(prev => Math.min(8, prev + 1))}
                      disabled={concurrency >= 8}
                    >
                      +
                    </Button>
                    <span className="ml-auto text-xs text-muted-foreground">
                      活跃: {activeTasks}
                    </span>
                  </div>
                )}
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
