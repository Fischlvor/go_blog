import service, { type ApiResponse } from '@/utils/request'
import type { PageInfo, PageResult } from './common'

// ==================== 请求类型 ====================

export interface ResourceCheckRequest {
  file_hash: string
  file_size: number
  file_name: string
}

export interface ResourceInitRequest {
  file_hash: string
  file_size: number
  file_name: string
  mime_type: string
}

export interface ResourceCompleteRequest {
  task_id: string
}

export interface ResourceCancelRequest {
  task_id: string
}

export interface ResourceProgressRequest {
  task_id: string
}

export interface ResourceListRequest extends PageInfo {
  file_name?: string
  mime_type?: string
}

export interface ResourceDeleteRequest {
  ids: number[]
}

// ==================== 响应类型 ====================

export interface ResourceCheckResponse {
  exists: boolean
  file_url?: string
  task_id?: string
  total_chunks?: number
  uploaded_chunks?: number[]
  missing_chunks?: number[]
}

export interface ResourceInitResponse {
  task_id: string
  total_chunks: number
  chunk_size: number
}

export interface ResourceUploadChunkResponse {
  success: boolean
  chunk_number: number
}

export interface ResourceCompleteResponse {
  file_url: string
  file_key: string
}

export interface ResourceProgressResponse {
  task_id: string
  total_chunks: number
  uploaded_chunks: number[]
  missing_chunks: number[]
  progress: number
}

export interface ResourceItem {
  id: number
  file_key: string
  file_name: string
  file_url: string
  file_size: number
  mime_type: string
  transcode_status: number  // 转码状态：0=无需转码, 1=转码中, 2=成功, 3=失败
  transcode_url?: string    // 转码后视频URL
  thumbnail_url?: string    // 缩略图URL
}

// ==================== API 方法 ====================

/**
 * 检查文件（秒传/续传检测）
 */
export function checkResource(data: ResourceCheckRequest): Promise<ApiResponse<ResourceCheckResponse>> {
  return service.post('/resources/check', data)
}

/**
 * 初始化上传任务
 */
export function initResource(data: ResourceInitRequest): Promise<ApiResponse<ResourceInitResponse>> {
  return service.post('/resources/init', data)
}

/**
 * 上传分片
 */
export function uploadChunk(taskId: string, chunkNumber: number, chunkData: Blob): Promise<ApiResponse<ResourceUploadChunkResponse>> {
  const formData = new FormData()
  formData.append('task_id', taskId)
  formData.append('chunk_number', chunkNumber.toString())
  formData.append('chunk_data', chunkData)
  
  return service.post('/resources/upload-chunk', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 完成上传
 */
export function completeResource(data: ResourceCompleteRequest): Promise<ApiResponse<ResourceCompleteResponse>> {
  return service.post('/resources/complete', data)
}

/**
 * 取消上传
 */
export function cancelResource(data: ResourceCancelRequest): Promise<ApiResponse<null>> {
  return service.post('/resources/cancel', data)
}

/**
 * 查询上传进度
 */
export function getResourceProgress(taskId: string): Promise<ApiResponse<ResourceProgressResponse>> {
  return service.get('/resources/progress', {
    params: { task_id: taskId }
  })
}

/**
 * 资源列表
 */
export function getResourceList(params: ResourceListRequest): Promise<ApiResponse<PageResult<ResourceItem>>> {
  return service.get('/resources/list', { params })
}

/**
 * 删除资源
 */
export function deleteResource(data: ResourceDeleteRequest): Promise<ApiResponse<null>> {
  return service.post('/resources/delete', data)
}

// ==================== 工具函数 ====================

/**
 * 计算文件MD5
 */
export async function calculateFileMD5(file: File): Promise<string> {
  const SparkMD5 = (await import('spark-md5')).default
  
  return new Promise((resolve, reject) => {
    const chunkSize = 2 * 1024 * 1024 // 2MB chunks for MD5 calculation
    const chunks = Math.ceil(file.size / chunkSize)
    const spark = new SparkMD5.ArrayBuffer()
    const reader = new FileReader()
    let currentChunk = 0

    reader.onload = (e) => {
      spark.append(e.target?.result as ArrayBuffer)
      currentChunk++

      if (currentChunk < chunks) {
        loadNext()
      } else {
        resolve(spark.end())
      }
    }

    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }

    function loadNext() {
      const start = currentChunk * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      reader.readAsArrayBuffer(file.slice(start, end))
    }

    loadNext()
  })
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 根据文件URL生成缩略图URL
 * 视频文件：将 xxx.mp4 转换为 xxx_thumb.jpg
 * 图片文件：直接返回原URL
 * 其他文件：返回空字符串
 */
export function getThumbnailUrl(fileUrl: string, mimeType: string): string {
  if (!fileUrl) return ''
  
  // 图片直接返回原URL
  if (mimeType.startsWith('image/')) {
    return fileUrl
  }
  
  // 视频生成缩略图URL
  if (mimeType.startsWith('video/')) {
    // 移除文件扩展名，添加 _thumb.jpg
    const lastDotIndex = fileUrl.lastIndexOf('.')
    if (lastDotIndex === -1) return ''
    return fileUrl.substring(0, lastDotIndex) + '_thumb.jpg'
  }
  
  return ''
}

/**
 * 判断是否为视频文件
 */
export function isVideoFile(mimeType: string): boolean {
  return mimeType.startsWith('video/')
}

/**
 * 判断是否为图片文件
 */
export function isImageFile(mimeType: string): boolean {
  return mimeType.startsWith('image/')
}

/**
 * 判断是否为音频文件
 */
export function isAudioFile(mimeType: string): boolean {
  return mimeType.startsWith('audio/')
}

/**
 * 获取文件类型分类
 */
export function getFileCategory(mimeType: string): 'video' | 'image' | 'audio' | 'document' | 'other' {
  if (mimeType.startsWith('video/')) return 'video'
  if (mimeType.startsWith('image/')) return 'image'
  if (mimeType.startsWith('audio/')) return 'audio'
  if (mimeType.includes('pdf') || mimeType.includes('document') || mimeType.includes('sheet') || mimeType.includes('presentation')) return 'document'
  return 'other'
}
