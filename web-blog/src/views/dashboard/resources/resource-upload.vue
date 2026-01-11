<template>
  <div class="resource-upload">
    <!-- 文件选择区域 -->
    <el-upload
      v-if="!uploadState.file"
      class="upload-dragger"
      drag
      :auto-upload="false"
      :show-file-list="false"
      :on-change="handleFileChange"
      accept="image/*,video/*,audio/*,.pdf,.zip,.rar,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
    >
      <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
      <div class="el-upload__text">
        将文件拖到此处，或<em>点击上传</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          支持图片、视频、音频、文档等格式，单文件最大 {{ maxFileSizeText }}
        </div>
      </template>
    </el-upload>

    <!-- 上传进度区域 -->
    <div v-else class="upload-progress">
      <div class="file-info">
        <el-icon :size="48" color="#409EFF">
          <Document />
        </el-icon>
        <div class="file-detail">
          <div class="file-name">{{ uploadState.file.name }}</div>
          <div class="file-size">{{ formatFileSize(uploadState.file.size) }}</div>
        </div>
      </div>

      <div class="progress-info">
        <div class="status-text">
          <span v-if="uploadState.status === 'hashing'">正在计算文件指纹...</span>
          <span v-else-if="uploadState.status === 'checking'">正在检查文件...</span>
          <span v-else-if="uploadState.status === 'instant'">秒传成功！</span>
          <span v-else-if="uploadState.status === 'uploading'">
            正在上传... {{ uploadState.uploadedChunks }}/{{ uploadState.totalChunks }} 块
          </span>
          <span v-else-if="uploadState.status === 'merging'">正在合并文件...</span>
          <span v-else-if="uploadState.status === 'success'">上传成功！</span>
          <span v-else-if="uploadState.status === 'error'" class="error-text">
            {{ uploadState.errorMsg }}
          </span>
        </div>
        <el-progress
          :percentage="uploadState.progress"
          :status="progressStatus"
          :stroke-width="10"
        />
      </div>

      <div v-if="uploadState.fileUrl" class="result-info">
        <el-input v-model="uploadState.fileUrl" readonly>
          <template #append>
            <el-button @click="copyUrl">复制链接</el-button>
          </template>
        </el-input>
      </div>

      <div class="action-buttons">
        <el-button
          v-if="uploadState.status === 'uploading'"
          type="danger"
          @click="handleCancel"
        >
          取消上传
        </el-button>
        <el-button
          v-if="['success', 'instant', 'error'].includes(uploadState.status)"
          type="primary"
          @click="handleReset"
        >
          继续上传
        </el-button>
        <el-button
          v-if="['success', 'instant'].includes(uploadState.status)"
          @click="emit('success')"
        >
          完成
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled, Document } from '@element-plus/icons-vue'
import type { UploadFile } from 'element-plus'
import {
  checkResource,
  initResource,
  uploadChunk,
  completeResource,
  cancelResource,
  calculateFileMD5,
  formatFileSize,
  getMaxFileSize
} from '@/api/resource'

const emit = defineEmits<{
  success: []
  cancel: []
}>()

interface UploadState {
  file: File | null
  fileHash: string
  taskId: string
  status: 'idle' | 'hashing' | 'checking' | 'instant' | 'uploading' | 'merging' | 'success' | 'error'
  progress: number
  totalChunks: number
  uploadedChunks: number
  fileUrl: string
  errorMsg: string
  abortController: AbortController | null
}

const uploadState = reactive<UploadState>({
  file: null,
  fileHash: '',
  taskId: '',
  status: 'idle',
  progress: 0,
  totalChunks: 0,
  uploadedChunks: 0,
  fileUrl: '',
  errorMsg: '',
  abortController: null
})

// 最大文件大小（字节）
let maxFileSize = 500 * 1024 * 1024 // 默认 500MB
let maxFileSizeText = '500MB'

// 初始化：获取最大文件大小
onMounted(async () => {
  try {
    const res = await getMaxFileSize()
    if (res.code === 0 && res.data?.max_size) {
      maxFileSize = res.data.max_size
      maxFileSizeText = formatFileSize(maxFileSize)
    }
  } catch (error) {
    console.warn('获取最大文件大小失败，使用默认值 500MB', error)
  }
})

const progressStatus = computed(() => {
  if (uploadState.status === 'error') return 'exception'
  if (['success', 'instant'].includes(uploadState.status)) return 'success'
  return undefined
})

// 并发上传配置
const CONCURRENT_LIMIT = 3   // 并发数
const MAX_RETRIES = 3        // 每块最大重试次数
const RETRY_DELAY = 1000     // 重试基础延迟(ms)

// 休眠函数
const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

// 并发上传（带自动重试）
const uploadChunksWithRetry = async (
  missingChunks: number[],
  taskId: string,
  file: File,
  chunkSize: number,
  abortSignal: AbortSignal
) => {
  interface ChunkTask {
    chunkNumber: number
    retries: number
    status: 'pending' | 'uploading' | 'success' | 'failed'
  }

  // 初始化任务队列
  const tasks: ChunkTask[] = missingChunks.map(n => ({
    chunkNumber: n,
    retries: 0,
    status: 'pending'
  }))

  let completedCount = uploadState.uploadedChunks

  // 更新进度
  const updateProgress = () => {
    uploadState.uploadedChunks = completedCount
    uploadState.progress = Math.round((completedCount / uploadState.totalChunks) * 90) + 10
  }

  // 并发执行器
  const runConcurrent = async () => {
    // 使用 Promise 池实现真正的并发
    const executeTask = async (task: ChunkTask): Promise<void> => {
      while (task.status !== 'success' && task.status !== 'failed') {
        if (abortSignal.aborted) return
        
        task.status = 'uploading'
        
        try {
          const start = task.chunkNumber * chunkSize
          const end = Math.min(start + chunkSize, file.size)
          const chunkData = file.slice(start, end)

          const res = await uploadChunk(taskId, task.chunkNumber, chunkData)

          if (res.code !== 0) {
            throw new Error(res.msg || `块 ${task.chunkNumber} 上传失败`)
          }

          task.status = 'success'
          completedCount++
          updateProgress()

        } catch (error: any) {
          task.retries++

          if (task.retries < MAX_RETRIES) {
            task.status = 'pending'
            console.warn(`块 ${task.chunkNumber} 失败，${RETRY_DELAY * task.retries}ms 后重试 (${task.retries}/${MAX_RETRIES})`)
            await sleep(RETRY_DELAY * task.retries)
          } else {
            task.status = 'failed'
            console.error(`块 ${task.chunkNumber} 上传失败，已重试 ${MAX_RETRIES} 次`)
          }
        }
      }
    }

    // 创建任务池，限制并发数
    const pool: Promise<void>[] = []
    let taskIndex = 0

    const runNext = async (): Promise<void> => {
      while (taskIndex < tasks.length) {
        if (abortSignal.aborted) return
        
        const task = tasks[taskIndex++]
        await executeTask(task)
      }
    }

    // 启动 CONCURRENT_LIMIT 个 worker
    for (let i = 0; i < Math.min(CONCURRENT_LIMIT, tasks.length); i++) {
      pool.push(runNext())
    }

    await Promise.all(pool)

    // 检查是否有失败的任务
    const failed = tasks.filter(t => t.status === 'failed')
    if (failed.length > 0) {
      throw new Error(`${failed.length} 个块上传失败，请重试`)
    }
  }

  await runConcurrent()
}

const handleFileChange = async (uploadFile: UploadFile) => {
  if (!uploadFile.raw) return

  // 前端初步检查：文件大小限制
  if (uploadFile.raw.size > maxFileSize) {
    const msg = `文件大小超过限制，最大允许 ${maxFileSizeText}，当前文件 ${formatFileSize(uploadFile.raw.size)}`
    ElMessage.error(msg)
    return
  }

  uploadState.file = uploadFile.raw
  uploadState.status = 'hashing'
  uploadState.progress = 0
  uploadState.errorMsg = ''

  try {
    // 1. 计算文件MD5
    uploadState.fileHash = await calculateFileMD5(uploadFile.raw)
    uploadState.progress = 10

    // 2. 检查文件（秒传/续传）
    uploadState.status = 'checking'
    const checkRes = await checkResource({
      file_hash: uploadState.fileHash,
      file_size: uploadFile.raw.size,
      file_name: uploadFile.raw.name
    })

    if (checkRes.code !== 0) {
      throw new Error(checkRes.msg || '检查文件失败')
    }

    // 秒传成功
    if (checkRes.data.exists) {
      uploadState.status = 'instant'
      uploadState.progress = 100
      uploadState.fileUrl = checkRes.data.file_url || ''
      ElMessage.success('秒传成功！')
      return
    }

    // 3. 初始化或续传
    let taskId = checkRes.data.task_id
    let totalChunks = checkRes.data.total_chunks || 0
    let missingChunks = checkRes.data.missing_chunks || []

    if (!taskId) {
      // 新上传，初始化任务
      const initRes = await initResource({
        file_hash: uploadState.fileHash,
        file_size: uploadFile.raw.size,
        file_name: uploadFile.raw.name,
        mime_type: uploadFile.raw.type || 'application/octet-stream'
      })

      if (initRes.code !== 0) {
        throw new Error(initRes.msg || '初始化任务失败')
      }

      taskId = initRes.data.task_id
      totalChunks = initRes.data.total_chunks
      const chunkSize = initRes.data.chunk_size

      // 生成所有块号
      missingChunks = Array.from({ length: totalChunks }, (_, i) => i)
    }

    uploadState.taskId = taskId
    uploadState.totalChunks = totalChunks
    uploadState.uploadedChunks = totalChunks - missingChunks.length
    uploadState.status = 'uploading'
    uploadState.progress = Math.round((uploadState.uploadedChunks / totalChunks) * 90) + 10

    // 4. 并发上传缺失的块（带自动重试）
    uploadState.abortController = new AbortController()
    const chunkSize = 4 * 1024 * 1024 // 4MB

    await uploadChunksWithRetry(
      missingChunks,
      taskId,
      uploadFile.raw,
      chunkSize,
      uploadState.abortController.signal
    )

    // 5. 合并文件
    uploadState.status = 'merging'
    const completeRes = await completeResource({ task_id: taskId })

    if (completeRes.code !== 0) {
      throw new Error(completeRes.msg || '合并文件失败')
    }

    uploadState.status = 'success'
    uploadState.progress = 100
    uploadState.fileUrl = completeRes.data.file_url
    ElMessage.success('上传成功！')
  } catch (error: any) {
    if (error.name === 'AbortError') {
      return
    }
    uploadState.status = 'error'
    uploadState.errorMsg = error.message || '上传失败'
    ElMessage.error(uploadState.errorMsg)
  }
}

const handleCancel = async () => {
  if (uploadState.abortController) {
    uploadState.abortController.abort()
  }

  if (uploadState.taskId) {
    try {
      await cancelResource({ task_id: uploadState.taskId })
    } catch {
      // 忽略取消错误
    }
  }

  handleReset()
  emit('cancel')
}

const handleReset = () => {
  uploadState.file = null
  uploadState.fileHash = ''
  uploadState.taskId = ''
  uploadState.status = 'idle'
  uploadState.progress = 0
  uploadState.totalChunks = 0
  uploadState.uploadedChunks = 0
  uploadState.fileUrl = ''
  uploadState.errorMsg = ''
  uploadState.abortController = null
}

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(uploadState.fileUrl)
    ElMessage.success('链接已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}
</script>

<style scoped lang="scss">
.resource-upload {
  .upload-dragger {
    width: 100%;

    :deep(.el-upload-dragger) {
      width: 100%;
      height: 200px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
    }
  }

  .upload-progress {
    .file-info {
      display: flex;
      align-items: center;
      padding: 20px;
      background-color: #f5f7fa;
      border-radius: 8px;
      margin-bottom: 20px;

      .file-detail {
        margin-left: 16px;

        .file-name {
          font-size: 16px;
          font-weight: 500;
          color: #303133;
          word-break: break-all;
        }

        .file-size {
          font-size: 14px;
          color: #909399;
          margin-top: 4px;
        }
      }
    }

    .progress-info {
      margin-bottom: 20px;

      .status-text {
        margin-bottom: 10px;
        font-size: 14px;
        color: #606266;

        .error-text {
          color: #f56c6c;
        }
      }
    }

    .result-info {
      margin-bottom: 20px;
    }

    .action-buttons {
      display: flex;
      justify-content: center;
      gap: 12px;
    }
  }
}
</style>
