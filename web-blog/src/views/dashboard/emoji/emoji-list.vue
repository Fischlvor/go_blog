<template>
  <div class="emoji-list">
    <div class="page-header">
      <h2>表情列表</h2>
      <div class="header-actions">
        <el-button type="primary" @click="showUploadDialog = true" :icon="Plus">
          批量上传表情
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filters">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索表情键名或文件名"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="表情组">
          <el-select 
            v-model="searchForm.groupName" 
            placeholder="选择表情组" 
            clearable 
            @change="handleSearch"
            style="width: 200px"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="group in emojiGroups"
              :key="group.id"
              :label="group.group_name"
              :value="group.group_key"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :icon="Search">搜索</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表情表格 -->
    <el-table :data="emojiList" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="缩略图" width="100" align="center">
        <template #default="{ row }">
          <div class="emoji-thumbnail">
            <img 
              :src="row.cdn_url" 
              :alt="row.key"
              class="emoji-image"
              @error="handleImageError"
            />
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="key" label="表情键" width="120">
        <template #default="{ row }">
          <el-tag type="info">{{ row.key }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="filename" label="文件名" min-width="180" />
      <el-table-column prop="group_name" label="表情组" width="130">
        <template #default="{ row }">
          <el-tag>{{ row.group_name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="雪碧图信息" width="160">
        <template #default="{ row }">
          <div class="sprite-info">
            <div>组号: {{ row.sprite_group }}</div>
            <div>位置: ({{ row.sprite_position_x }}, {{ row.sprite_position_y }})</div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="file_size" label="文件大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="upload_time" label="上传时间" min-width="150">
        <template #default="{ row }">
          {{ formatTime(row.upload_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '正常' : '已删除' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 1"
            type="danger"
            size="small"
            @click="deleteEmoji(row)"
            :icon="Delete"
          >
            删除
          </el-button>
          <el-button
            v-else
            type="success"
            size="small"
            @click="restoreEmoji(row)"
            :icon="RefreshRight"
          >
            恢复
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 上传对话框 -->
    <el-dialog 
      v-model="showUploadDialog" 
      title="批量上传表情" 
      width="600px"
      class="upload-dialog"
    >
      <el-form :model="uploadForm" label-width="100px">
        <el-form-item label="表情组" required>
          <el-select 
            v-model="uploadForm.groupKey" 
            placeholder="选择表情组"
            style="width: 100%"
          >
            <el-option
              v-for="group in emojiGroups"
              :key="group.id"
              :label="group.group_name"
              :value="group.group_key"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择文件" required>
          <!-- 无文件时：显示上传按钮 -->
          <div v-if="uploadForm.files.length === 0" class="upload-empty-state">
            <el-upload
              ref="uploadRef"
              :auto-upload="false"
              :multiple="true"
              :file-list="[]"
              :on-change="handleFileChange"
              :before-upload="beforeUpload"
              accept=".png,.jpg,.jpeg"
              drag
              class="upload-dragger-full"
              :show-file-list="false"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                拖拽文件到此处，或<em>点击选择文件</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                    只能上传 PNG/JPG 文件，建议尺寸 64x64px
                </div>
              </template>
            </el-upload>
          </div>
          
          <!-- 有文件时：显示文件列表 -->
          <div v-else class="upload-with-files">
            <div class="file-list-container-full">
              <div class="file-list-header">
                <span>已选择文件 ({{ uploadForm.files.length }})</span>
                <div class="header-actions">
                  <el-button 
                    size="small" 
                    type="primary" 
                    text 
                    @click="addMoreFiles"
                    :icon="Plus"
                    :disabled="uploading"
                  >
                    继续添加
                  </el-button>
                  <el-button 
                    size="small" 
                    type="danger" 
                    text 
                    @click="clearAllFiles"
                    :disabled="uploading"
                  >
                    清空
                  </el-button>
                </div>
              </div>
              <div class="custom-file-list">
                <div 
                  v-for="(file, index) in uploadForm.files" 
                  :key="file.uid" 
                  class="file-item"
                >
                  <div class="file-info">
                    <el-icon class="file-icon"><document /></el-icon>
                    <span class="file-name" :title="file.name">{{ file.name }}</span>
                    <span class="file-size">{{ formatFileSize(file.size || 0) }}</span>
                  </div>
                  <el-button 
                    size="small" 
                    type="danger" 
                    text 
                    @click="removeFile(index)"
                    class="remove-btn"
                    :disabled="uploading"
                  >
                    <el-icon><close /></el-icon>
                  </el-button>
                </div>
              </div>
            </div>
            
            <!-- 隐藏的上传组件，用于继续添加文件 -->
            <el-upload
              ref="addMoreUploadRef"
              :auto-upload="false"
              :multiple="true"
              :file-list="[]"
              :on-change="handleFileChange"
              :before-upload="beforeUpload"
              accept=".png,.jpg,.jpeg"
              :show-file-list="false"
              style="display: none;"
            />
          </div>
        </el-form-item>
        
        <!-- 上传进度展示 -->
        <el-form-item :label="uploading ? '上传进度' : ' '">
          <div class="progress-container">
            <transition name="progress-slide" mode="out-in">
              <div v-if="uploading" class="simple-progress">
                <div class="progress-bar">
                  <div 
                    class="progress-fill" 
                    :style="{ width: uploadProgress.percentage + '%' }"
                  ></div>
                </div>
                <div class="progress-info">
                  <span class="progress-text">{{ uploadProgress.percentage }}%</span>
                  <span class="current-file" v-if="uploadProgress.currentFile">
                    {{ uploadProgress.currentFile }}
                  </span>
                </div>
              </div>
            </transition>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button type="primary" @click="handleUpload" :loading="uploading">
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Delete, RefreshRight, UploadFilled, Document, Close } from '@element-plus/icons-vue'
import type { UploadFile } from 'element-plus'
import { 
  getEmojiList, 
  getEmojiGroups, 
  uploadEmoji as apiUploadEmoji,
  uploadEmojiStream as apiUploadEmojiStream, 
  deleteEmoji as apiDeleteEmoji, 
  restoreEmoji as apiRestoreEmoji
} from '@/api/emoji'

// 接口类型定义
interface Emoji {
  id: number
  key: string
  filename: string
  group_name: string
  sprite_group: number
  sprite_position_x: number
  sprite_position_y: number
  file_size: number
  cdn_url: string
  upload_time: string
  status: number
  created_by: string
}

interface EmojiGroup {
  id: number
  group_key: string
  group_name: string
  description: string
  emoji_count: number
  status: number
}

interface EmojiListRequest {
  page: number
  page_size: number
  group_key?: string
  keyword?: string
}

// 响应式数据
const loading = ref(false)
const uploading = ref(false)
const showUploadDialog = ref(false)
const emojiList = ref<Emoji[]>([])
const emojiGroups = ref<EmojiGroup[]>([])

// 上传进度相关
const uploadProgress = ref({
  current: 0,
  total: 0,
  percentage: 0,
  currentFile: '',
  completedFiles: [] as string[],
  failedFiles: [] as string[]
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  groupName: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 上传表单
const uploadForm = reactive({
  groupKey: '',
  files: [] as UploadFile[]
})

// 获取表情列表
const fetchEmojiList = async () => {
  loading.value = true
  try {
    const params: EmojiListRequest = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...(searchForm.keyword && { keyword: searchForm.keyword }),
      ...(searchForm.groupName && { group_key: searchForm.groupName })
    }

    const result = await getEmojiList(params)
    if (result.code === 0) {
      emojiList.value = result.data.list || []
      pagination.total = result.data.total || 0
    } else {
      ElMessage.error(result.msg || '获取表情列表失败')
    }
  } catch (error) {
    console.error('获取表情列表失败:', error)
    ElMessage.error('获取表情列表失败')
  } finally {
    loading.value = false
  }
}

// 获取表情组列表
const fetchEmojiGroups = async () => {
  try {
    const result = await getEmojiGroups()
    if (result.code === 0) {
      emojiGroups.value = result.data || []
    }
  } catch (error) {
    console.error('获取表情组失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchEmojiList()
}

// 分页处理
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchEmojiList()
}

const handleCurrentChange = (val: number) => {
  pagination.page = val
  fetchEmojiList()
}

// 删除表情
const deleteEmoji = async (emoji: Emoji) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除表情 "${emoji.key}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const result = await apiDeleteEmoji(emoji.id)
    if (result.code === 0) {
      ElMessage.success('删除成功')
      fetchEmojiList()
    } else {
      ElMessage.error(result.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除表情失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 恢复表情
const restoreEmoji = async (emoji: Emoji) => {
  try {
    const result = await apiRestoreEmoji(emoji.id)
    if (result.code === 0) {
      ElMessage.success('恢复成功')
      fetchEmojiList()
    } else {
      ElMessage.error(result.msg || '恢复失败')
    }
  } catch (error) {
    console.error('恢复表情失败:', error)
    ElMessage.error('恢复失败')
  }
}

// 文件上传处理
const handleFileChange = (file: UploadFile, fileList: UploadFile[]) => {
  uploadForm.files = fileList
  console.log('文件列表更新:', fileList.length, '个文件')
}

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 清空所有文件
const clearAllFiles = () => {
  uploadForm.files = []
}

// 移除单个文件
const removeFile = (index: number) => {
  uploadForm.files.splice(index, 1)
}

// 继续添加文件
const addMoreFiles = () => {
  const addMoreUploadRef = document.querySelector('input[type="file"]') as HTMLInputElement
  if (addMoreUploadRef) {
    addMoreUploadRef.click()
  }
}

// 处理图片加载错误
const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjQiIGhlaWdodD0iNjQiIHZpZXdCb3g9IjAgMCA2NCA2NCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjY0IiBoZWlnaHQ9IjY0IiBmaWxsPSIjRjVGN0ZBIi8+CjxwYXRoIGQ9Ik0yMCAyMEg0NFY0NEgyMFYyMFoiIHN0cm9rZT0iI0MxQzRDRCIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPHBhdGggZD0iTTI4IDI4TDM2IDM2TDQwIDMyIiBzdHJva2U9IiNDMUM0Q0QiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+Cjwvc3ZnPgo='
  img.alt = '图片加载失败'
}

// 处理上传
const handleUpload = async () => {
  if (!uploadForm.groupKey) {
    ElMessage.error('请选择表情组')
    return
  }

  if (!uploadForm.files || uploadForm.files.length === 0) {
    ElMessage.error('请选择要上传的文件')
    return
  }

  // 过滤出有效的文件（状态为ready的文件）
  const validFiles = uploadForm.files.filter(file => file.status !== 'fail')
  if (validFiles.length === 0) {
    ElMessage.error('没有有效的文件可以上传')
    return
  }

  uploading.value = true
  uploadProgress.value = {
    current: 0,
    total: validFiles.length,
    percentage: 0,
    currentFile: '',
    completedFiles: [],
    failedFiles: []
  }

  try {
    // 使用SSE流式上传
    await handleStreamUpload(validFiles)
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
    
    // 错误时也要重置状态
    setTimeout(() => {
      showUploadDialog.value = false
      uploadForm.files = []
      uploadForm.groupKey = ''
      uploadProgress.value = {
        current: 0,
        total: 0,
        percentage: 0,
        currentFile: '',
        completedFiles: [],
        failedFiles: []
      }
    }, 1000)
  } finally {
    uploading.value = false
  }
}

// SSE流式上传处理
const handleStreamUpload = (validFiles: UploadFile[]) => {
  // 创建FormData
  const formData = new FormData()
  formData.append('group_key', uploadForm.groupKey)
  
  // 添加所有文件到FormData
  validFiles.forEach((file) => {
    formData.append('files', file.raw!)
  })

  // 使用统一的API管理系统处理流式上传
  return apiUploadEmojiStream(formData, handleSSEEvent)
}

// 处理SSE事件
const handleSSEEvent = (event: string, data: any) => {
  switch (event) {
    case 'start':
      uploadProgress.value.total = data.total
      uploadProgress.value.current = 0
      uploadProgress.value.percentage = 0
      uploadProgress.value.completedFiles = []
      uploadProgress.value.failedFiles = []
      console.log('开始上传:', data.message)
      break
      
    case 'progress':
      uploadProgress.value.current = data.current
      uploadProgress.value.total = data.total
      uploadProgress.value.percentage = Math.round((data.current / data.total) * 100)
      uploadProgress.value.currentFile = data.filename
      console.log(`进度: ${data.current}/${data.total} - ${data.filename}`)
      break
      
    case 'file_complete':
      uploadProgress.value.completedFiles.push(data.filename)
      // 更新进度
      uploadProgress.value.current = data.current
      uploadProgress.value.percentage = Math.round((data.current / data.total) * 100)
      
      // 延迟移除已完成的文件，让用户看到完成状态
      setTimeout(() => {
        const fileIndex = uploadForm.files.findIndex(f => f.name === data.filename)
        if (fileIndex > -1) {
          uploadForm.files.splice(fileIndex, 1)
        }
      }, 500)
      console.log(`文件完成: ${data.filename}`)
      break
      
    case 'file_error':
      uploadProgress.value.failedFiles.push(data.filename)
      ElMessage.error(`文件 ${data.filename} 上传失败: ${data.error}`)
      break
      
    case 'complete':
      uploadProgress.value.percentage = 100
      uploadProgress.value.current = uploadProgress.value.total
      uploadProgress.value.currentFile = '' // 清空当前文件名
      console.log('上传完成事件:', data)
      
      // 安全地访问数据字段，提供默认值
      const successCount = data?.success ?? 0
      const failedCount = data?.failed ?? 0
      const totalCount = data?.total ?? 0
      
      // 显示完成消息
      if (failedCount > 0) {
        ElMessage.warning(`上传完成！成功: ${successCount}, 失败: ${failedCount}`)
      } else {
        ElMessage.success(`上传完成！成功上传 ${successCount} 个文件`)
      }
      
      // 延迟关闭对话框，让用户看到完成状态
      setTimeout(() => {
      showUploadDialog.value = false
      uploadForm.files = []
      uploadForm.groupKey = ''
      // 重置进度状态
      uploadProgress.value = {
        current: 0,
        total: 0,
        percentage: 0,
        currentFile: '',
        completedFiles: [],
        failedFiles: []
      }
      fetchEmojiList()
      }, 200)
      break
      
    case 'error':
      ElMessage.error(data.message)
      break
      
    case 'warning':
      ElMessage.warning(data.message)
      break
  }
}

// 工具函数
const formatFileSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(1)} MB`
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString()
}

// 初始化
onMounted(() => {
  fetchEmojiGroups()
  fetchEmojiList()
})
</script>

<style scoped lang="scss">
// 进度条容器 - 防抖动
.progress-container {
  width: 100%;
  min-height: 50px; // 固定最小高度防止抖动
  transition: all 0.3s ease;
  position: relative;
}

// 简洁进度条样式
.simple-progress {
  width: 100%;
  
  .progress-bar {
    width: 100%;
    height: 8px;
    background-color: #f0f2f5;
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 8px;
    
    .progress-fill {
      height: 100%;
      background: linear-gradient(90deg, #409eff 0%, #67c23a 100%);
      border-radius: 4px;
      transition: width 0.3s ease;
    }
  }
  
  .progress-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    min-height: 20px; // 固定文字区域高度
    
    .progress-text {
      font-weight: 600;
      color: #409eff;
    }
    
    .current-file {
      color: #606266;
      font-size: 13px;
      max-width: 200px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

// 进度条过渡动画
.progress-slide-enter-active,
.progress-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.progress-slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.progress-slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}



.upload-dragger {
  width: 100%;
}

.emoji-list {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      color: #333;
    }

    .header-actions {
      display: flex;
      gap: 10px;
    }
  }

  .filters {
    background: #f8f9fa;
    padding: 15px;
    border-radius: 8px;
    margin-bottom: 20px;

    .search-form {
      margin: 0;
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

// 上传弹窗样式
:deep(.upload-dialog) {
  .el-dialog__body {
    max-height: 70vh;
    overflow-y: auto;
    padding: 20px 24px;
  }
  
  .el-dialog__header {
    padding: 20px 24px 10px;
  }
  
  .el-dialog__footer {
    padding: 10px 24px 20px;
  }
}

// 无文件时的上传区域样式
.upload-empty-state {
  width: 100%;
  
  .upload-dragger-full {
    width: 100%;
    height: 240px;
    display: flex;
    flex-direction: column;
    
    :deep(.el-upload-dragger) {
      height: 200px;
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      border: 2px dashed #d9d9d9;
      border-radius: 8px;
      transition: all 0.3s ease;
      margin-bottom: 8px;
      
      &:hover {
        border-color: #409eff;
        background-color: #f5f7fa;
      }
    }
    
    :deep(.el-upload__tip) {
      margin-top: 0 !important;
      padding: 6px 8px;
      text-align: center;
      color: #909399;
      font-size: 12px;
      background: transparent;
      border: none;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      white-space: nowrap;
      overflow: hidden;
    }
    
    :deep(.el-icon--upload) {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 16px;
    }
    
    :deep(.el-upload__text) {
      color: #606266;
      font-size: 14px;
      text-align: center;
      white-space: nowrap;
      margin-top: 8px;
      
      em {
        color: #409eff;
        font-style: normal;
      }
    }
  }
}

// 有文件时的文件列表样式
.upload-with-files {
  width: 100%;
}

.file-list-container-full {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 240px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  
  .file-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 16px;
    background: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    font-size: 14px;
    font-weight: 500;
    color: #303133;
    
    .header-actions {
      display: flex;
      gap: 8px;
    }
  }
  
  .custom-file-list {
    flex: 1;
    overflow-y: auto;
    background: #fff;
    
    // 自定义滚动条
    scrollbar-width: thin;
    scrollbar-color: #c1c4cd transparent;
    
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: transparent;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #c1c4cd;
      border-radius: 3px;
      
      &:hover {
        background: #a4a7b0;
      }
    }
    
    .empty-state {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #909399;
      font-size: 14px;
      min-height: 100px;
    }
    
    .file-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 12px;
      border-bottom: 1px solid #f5f7fa;
      transition: all 0.3s ease;
      
      &:hover {
        background: #f5f7fa;
      }
      
      &:last-child {
        border-bottom: none;
      }
      
      .file-info {
        display: flex;
        align-items: center;
        flex: 1;
        min-width: 0;
        
        .file-icon {
          color: #409eff;
          margin-right: 8px;
          flex-shrink: 0;
        }
        
        .file-name {
          flex: 1;
          min-width: 0;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          color: #303133;
          font-size: 14px;
          margin-right: 8px;
        }
        
        .file-size {
          color: #909399;
          font-size: 12px;
          flex-shrink: 0;
        }
      }
      
      .remove-btn {
        margin-left: 8px;
        flex-shrink: 0;
      }
    }
  }
}

// 表情缩略图样式
.emoji-thumbnail {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 64px;
  height: 64px;
  margin: 0 auto;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  background: #fafafa;
  overflow: hidden;
  
  .emoji-image {
    width: 48px;
    height: 48px;
    object-fit: contain;
    border-radius: 4px;
    transition: transform 0.2s ease;
    
    &:hover {
      transform: scale(1.1);
    }
  }
}

// 表格行高调整
:deep(.el-table .el-table__row) {
  height: 80px;
  
  .el-table__cell {
    padding: 8px 0;
  }
}

// 雪碧图信息样式
.sprite-info {
  font-size: 12px;
  line-height: 1.4;
  color: #606266;
  
  div {
    margin-bottom: 2px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
