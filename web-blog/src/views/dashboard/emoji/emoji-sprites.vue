<template>
  <div class="emoji-sprites">
    <div class="page-header">
      <h2>雪碧图管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="regenerateSprites" :loading="regenerating" :icon="Refresh">
          重新生成全部雪碧图
        </el-button>
        <el-button @click="fetchSpriteList" :icon="RefreshRight">
          刷新列表
        </el-button>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats-cards">
      <el-card class="stat-card">
        <div class="stat-item">
          <div class="stat-number">{{ spriteList.length }}</div>
          <div class="stat-label">雪碧图总数</div>
        </div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-item">
          <div class="stat-number">{{ totalEmojis }}</div>
          <div class="stat-label">表情总数</div>
        </div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-item">
          <div class="stat-number">{{ formatFileSize(totalSize) }}</div>
          <div class="stat-label">总文件大小</div>
        </div>
      </el-card>
    </div>

    <!-- 雪碧图列表 -->
    <el-table :data="spriteList" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="sprite_group" label="组号" width="100">
        <template #default="{ row }">
          <el-tag type="info">{{ row.sprite_group }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="filename" label="文件名" min-width="200" />
      <el-table-column label="预览" width="120">
        <template #default="{ row }">
          <div class="preview-cell">
            <el-popover
              v-if="row.cdn_url"
              placement="right"
              :width="520"
              trigger="hover"
              :show-arrow="true"
            >
              <template #reference>
                <el-image
                  :src="row.cdn_url"
                  fit="cover"
                  class="sprite-preview-thumbnail"
                >
                  <template #error>
                    <div class="image-slot">
                      <el-icon><Picture /></el-icon>
                    </div>
                  </template>
                </el-image>
              </template>
              <div class="sprite-preview-popover">
                <el-image
                  :src="row.cdn_url"
                  fit="contain"
                  style="max-width: 100%; max-height: 450px; background: white; padding: 4px; border-radius: 4px;"
                >
                  <template #error>
                    <div class="image-slot">
                      <el-icon><Picture /></el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
            </el-popover>
            <div v-else class="no-image">
              <el-icon><Picture /></el-icon>
              <span>暂无图片</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="尺寸" width="110">
        <template #default="{ row }">
          {{ row.width }} × {{ row.height }}
        </template>
      </el-table-column>
      <el-table-column prop="emoji_count" label="表情数量" width="110">
        <template #default="{ row }">
          <el-badge :value="row.emoji_count" class="item">
            <el-icon><Sunny /></el-icon>
          </el-badge>
        </template>
      </el-table-column>
      <el-table-column prop="file_size" label="文件大小" width="110">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '正常' : '异常' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="生成时间" min-width="150">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            @click="viewSpriteDetails(row)"
            :icon="View"
          >
            查看详情
          </el-button>
          <el-button
            v-if="row.cdn_url"
            type="success"
            size="small"
            @click="downloadSprite(row)"
            :icon="Download"
          >
            下载
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 雪碧图生成进度组件 -->
    <SpriteGenerationProgress
      v-model="showTaskDialog"
      :status="taskStatus"
      :message="taskMessage"
      :group-progress="groupProgress"
      :sprite-progress="spriteProgress"
      :show-group-progress="true"
      :show-sprite-progress="true"
    />

    <!-- 雪碧图详情对话框 -->
    <el-dialog v-model="showDetailsDialog" title="雪碧图详情" width="800px">
      <div v-if="selectedSprite" class="sprite-details">
        <div class="detail-section">
          <h3>基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="组号">{{ selectedSprite.sprite_group }}</el-descriptions-item>
            <el-descriptions-item label="文件名">{{ selectedSprite.filename }}</el-descriptions-item>
            <el-descriptions-item label="尺寸">{{ selectedSprite.width }} × {{ selectedSprite.height }}</el-descriptions-item>
            <el-descriptions-item label="表情数量">{{ selectedSprite.emoji_count }}</el-descriptions-item>
            <el-descriptions-item label="文件大小">{{ formatFileSize(selectedSprite.file_size) }}</el-descriptions-item>
            <el-descriptions-item label="生成时间">{{ formatTime(selectedSprite.created_at) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="detail-section" v-if="selectedSprite.cdn_url">
          <h3>雪碧图预览</h3>
          <div class="sprite-preview-large">
            <el-image
              :src="selectedSprite.cdn_url"
              :preview-src-list="[selectedSprite.cdn_url]"
              fit="contain"
            >
              <template #error>
                <div class="image-slot-large">
                  <el-icon><Picture /></el-icon>
                  <span>图片加载失败</span>
                </div>
              </template>
            </el-image>
          </div>
        </div>

        <div class="detail-section">
          <h3>包含的表情 ({{ spriteEmojis.length }})</h3>
          <el-table :data="spriteEmojis" max-height="300" border>
            <el-table-column prop="key" label="表情键" width="120" />
            <el-table-column prop="filename" label="文件名" />
            <el-table-column label="位置" width="150">
              <template #default="{ row }">
                ({{ row.sprite_position_x }}, {{ row.sprite_position_y }})
              </template>
            </el-table-column>
            <el-table-column prop="group_name" label="表情组" width="150" />
          </el-table>
        </div>
      </div>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, RefreshRight, View, Download, Picture, Sunny } from '@element-plus/icons-vue'
import SpriteGenerationProgress from '@/components/SpriteGenerationProgress.vue'
import { 
  getSpriteList, 
  getEmojiList, 
  regenerateSpritesStream,
  type EmojiSprite,
  type Emoji
} from '@/api/emoji'

// 响应式数据
const loading = ref(false)
const regenerating = ref(false)
const showDetailsDialog = ref(false)
const showTaskDialog = ref(false)
const spriteList = ref<EmojiSprite[]>([])
const selectedSprite = ref<EmojiSprite | null>(null)
const spriteEmojis = ref<Emoji[]>([])

// 任务相关
const taskStatus = ref<'running' | 'completed' | 'failed' | ''>('')
const taskMessage = ref('')
const taskId = ref<number | null>(null)
const totalGroups = ref(0)
const processedGroups = ref(0)
const currentGroupSprites = ref(0)
const currentGroupTotalSprites = ref(0)

// 计算属性
const totalEmojis = computed(() => {
  return spriteList.value.reduce((sum, sprite) => sum + sprite.emoji_count, 0)
})

const totalSize = computed(() => {
  return spriteList.value.reduce((sum, sprite) => sum + sprite.file_size, 0)
})

// 总组数进度百分比
const groupProgress = computed(() => {
  if (totalGroups.value === 0) return 0
  return Math.round((processedGroups.value / totalGroups.value) * 100)
})

// 当前组的雪碧图进度百分比
const spriteProgress = computed(() => {
  if (currentGroupTotalSprites.value === 0) return 0
  return Math.round((currentGroupSprites.value / currentGroupTotalSprites.value) * 100)
})

// 获取雪碧图列表
const fetchSpriteList = async () => {
  loading.value = true
  try {
    // 这里需要后端提供雪碧图列表接口
    // 暂时使用模拟数据
    const result = await getSpriteList()
    if (result.code === 0) {
      spriteList.value = result.data || []
    } else {
      ElMessage.error(result.msg || '获取雪碧图列表失败')
    }
  } catch (error) {
    console.error('获取雪碧图列表失败:', error)
    ElMessage.error('获取雪碧图列表失败')
  } finally {
    loading.value = false
  }
}

// 查看雪碧图详情
const viewSpriteDetails = async (sprite: EmojiSprite) => {
  selectedSprite.value = sprite
  
  // 获取该雪碧图包含的表情
  try {
    const response = await fetch(`/api/emoji/list?sprite_group=${sprite.sprite_group}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const result = await response.json()
      if (result.code === 0) {
        spriteEmojis.value = result.data.list || []
      }
    }
  } catch (error) {
    console.error('获取表情列表失败:', error)
    spriteEmojis.value = []
  }

  showDetailsDialog.value = true
}

// 下载雪碧图
const downloadSprite = (sprite: EmojiSprite) => {
  if (!sprite.cdn_url) {
    ElMessage.error('该雪碧图暂无下载链接')
    return
  }

  const link = document.createElement('a')
  link.href = sprite.cdn_url
  link.download = sprite.filename
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 重新生成雪碧图
const regenerateSprites = async () => {
  try {
    await ElMessageBox.confirm(
      '重新生成雪碧图将覆盖现有的雪碧图文件，确定要继续吗？',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    regenerating.value = true
    // 先显示对话框，然后再开始生成
    showTaskDialog.value = true
    taskStatus.value = 'running'
    taskMessage.value = '正在生成雪碧图...'
    totalGroups.value = 0
    processedGroups.value = 0
    currentGroupSprites.value = 0
    currentGroupTotalSprites.value = 0

    // 等待DOM更新后再开始生成
    await nextTick()

    // 使用SSE流式响应处理
    const handleSSEEvent = (event: string, data: any) => {
      console.log('SSE Event:', event, data)
      switch (event) {
        case 'start':
          ElMessage.info('开始生成雪碧图...')
          taskMessage.value = '开始生成雪碧图...'
          totalGroups.value = data.total_groups || 1
          processedGroups.value = 0
          break
        case 'group_start':
          taskMessage.value = `正在处理表情组: ${data.group_name} (${data.current}/${data.total})`
          processedGroups.value = data.current - 1
          currentGroupSprites.value = 0
          currentGroupTotalSprites.value = 1
          console.log(`正在处理表情组: ${data.group_name}`)
          break
        case 'sprite_progress':
          taskMessage.value = `${data.group_name}: 生成雪碧图 ${data.current}/${data.total}`
          currentGroupSprites.value = data.current
          currentGroupTotalSprites.value = data.total
          console.log(`${data.group_name}: 生成雪碧图 ${data.current}/${data.total}`)
          break
        case 'group_complete':
          ElMessage.success(`表情组 ${data.group_name} 生成完成，共 ${data.sprite_count} 个雪碧图`)
          processedGroups.value = data.current
          break
        case 'complete':
          ElMessage.success(`所有表情组生成完成！共 ${data.groups_count} 个组，${data.sprites_count} 个雪碧图`)
          processedGroups.value = totalGroups.value
          currentGroupSprites.value = currentGroupTotalSprites.value
          taskStatus.value = 'completed'
          taskMessage.value = '生成完成'
          fetchSpriteList()
          break
        case 'error':
          ElMessage.error(data.message)
          taskStatus.value = 'failed'
          taskMessage.value = `错误: ${data.message}`
          break
      }
    }

    await regenerateSpritesStream([], handleSSEEvent)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重新生成雪碧图失败:', error)
      ElMessage.error('操作失败')
      taskStatus.value = 'failed'
      taskMessage.value = '操作失败'
    }
  } finally {
    regenerating.value = false
  }
}

// 开始任务监控（保留用于兼容性，但不再使用）
const startTaskMonitoring = () => {
  if (!taskId.value) return

  showTaskDialog.value = true
  taskStatus.value = 'running'
  taskMessage.value = '正在生成雪碧图...'
  processedGroups.value = 0
  currentGroupSprites.value = 0

  const checkTask = async () => {
    if (!taskId.value) return

    try {
      const response = await fetch(`/api/emoji/task/${taskId.value}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const result = await response.json()
        if (result.code === 0) {
          const task = result.data
          currentGroupSprites.value = task.progress
          taskStatus.value = task.status
          taskMessage.value = task.message

          if (task.status === 'completed') {
            ElMessage.success('雪碧图生成完成')
            fetchSpriteList()
          } else if (task.status === 'failed') {
            ElMessage.error('雪碧图生成失败')
          } else if (task.status === 'running') {
            // 继续监控
            setTimeout(checkTask, 2000)
          }
        }
      }
    } catch (error) {
      console.error('获取任务状态失败:', error)
    }
  }

  checkTask()
}

// 工具函数
const formatFileSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(1)} MB`
}

const formatTime = (time: string | Date) => {
  const date = time instanceof Date ? time : new Date(time)
  return date.toLocaleString()
}

// 初始化
onMounted(() => {
  fetchSpriteList()
})
</script>

<style scoped lang="scss">
.emoji-sprites {
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

  .stats-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-bottom: 20px;

    .stat-card {
      .stat-item {
        text-align: center;
        padding: 10px;

        .stat-number {
          font-size: 2em;
          font-weight: bold;
          color: #409eff;
          margin-bottom: 5px;
        }

        .stat-label {
          color: #666;
          font-size: 14px;
        }
      }
    }
  }

  .preview-cell {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
  }

  .sprite-preview-thumbnail {
    width: 80px;
    height: 40px;
    border-radius: 4px;
    border: 1px solid #dcdfe6;
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 2px 12px rgba(64, 158, 255, 0.3);
    }

    .image-slot {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 100%;
      background: #f5f7fa;
      color: #909399;
    }
  }

  .sprite-preview-popover {
    text-align: center;
  }

  .sprite-preview {
    width: 80px;
    height: 40px;
    border-radius: 4px;

    .image-slot {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 100%;
      background: #f5f7fa;
      color: #909399;
    }
  }

  .no-image {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 40px;
    background: #f5f7fa;
    color: #909399;
    border-radius: 4px;
    font-size: 12px;
  }

  .item {
    margin-right: 10px;
  }
}

.sprite-details {
  .detail-section {
    margin-bottom: 30px;

    h3 {
      margin: 0 0 15px 0;
      color: #333;
      border-bottom: 2px solid #409eff;
      padding-bottom: 5px;
    }
  }

  .sprite-preview-large {
    text-align: center;

    .el-image {
      max-width: 100%;
      max-height: 400px;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
    }

    .image-slot-large {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      width: 300px;
      height: 200px;
      background: #f5f7fa;
      color: #909399;
    }
  }
}

</style>
