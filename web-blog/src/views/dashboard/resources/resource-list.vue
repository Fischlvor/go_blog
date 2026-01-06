<template>
  <div class="resource-list">
    <div class="title">
      <el-row>资源列表</el-row>
      <el-button-group>
        <el-button type="primary" icon="Upload" @click="uploadDialogVisible = true">
          上传资源
        </el-button>
        <el-button type="danger" icon="Delete" @click="handleBulkDeleteClick">
          批量删除
        </el-button>
      </el-button-group>
    </div>

    <div class="resource-list-request">
      <el-form :inline="true" :model="listRequest">
        <el-form-item label="文件名称">
          <el-input v-model="listRequest.file_name" placeholder="请输入文件名称" clearable />
        </el-form-item>
        <el-form-item label="文件类型">
          <el-select v-model="listRequest.mime_type" placeholder="全部" style="width: 200px" clearable>
            <el-option
              v-for="item in mimeTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getResourceTableData">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table ref="multipleTableRef" :data="tableData" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="60" />
      <el-table-column label="预览" width="100">
        <template #default="scope">
          <!-- 视频预览（使用后端返回的缩略图URL） -->
          <div v-if="getFileCategory(scope.row.mime_type) === 'video'" class="preview-cell">
            <el-image
              v-if="scope.row.thumbnail_url"
              :src="scope.row.thumbnail_url"
              fit="cover"
              style="width: 60px; height: 60px; border-radius: 4px;"
            >
              <template #error>
                <div class="preview-fallback">
                  <el-icon :size="24" color="#409EFF"><VideoPlay /></el-icon>
                </div>
              </template>
            </el-image>
            <div v-else class="preview-fallback">
              <el-icon :size="24" color="#409EFF"><VideoPlay /></el-icon>
            </div>
            <div class="video-play-overlay">
              <el-icon :size="16" color="#fff"><VideoPlay /></el-icon>
            </div>
          </div>
          <!-- 图片预览 -->
          <el-image
            v-else-if="getFileCategory(scope.row.mime_type) === 'image'"
            :src="scope.row.file_url"
            :preview-src-list="[scope.row.file_url]"
            fit="cover"
            style="width: 60px; height: 60px; border-radius: 4px;"
          />
          <!-- 音频图标 -->
          <div v-else-if="getFileCategory(scope.row.mime_type) === 'audio'" class="preview-icon">
            <el-icon :size="32" color="#E6A23C"><Headset /></el-icon>
          </div>
          <!-- 文档图标 -->
          <div v-else-if="getFileCategory(scope.row.mime_type) === 'document'" class="preview-icon">
            <el-icon :size="32" color="#F56C6C"><Document /></el-icon>
          </div>
          <!-- 其他文件图标 -->
          <div v-else class="preview-icon">
            <el-icon :size="32" color="#909399"><Folder /></el-icon>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="file_name" label="文件名" min-width="200" show-overflow-tooltip />
      <el-table-column prop="file_size" label="大小" width="120">
        <template #default="scope">
          {{ formatFileSize(scope.row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="mime_type" label="类型" width="150" />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <!-- 非视频文件不显示状态 -->
          <span v-if="getFileCategory(scope.row.mime_type) !== 'video'">-</span>
          <!-- 视频转码状态 -->
          <el-tag v-else-if="scope.row.transcode_status === 1" type="warning" size="small">转码中</el-tag>
          <el-tag v-else-if="scope.row.transcode_status === 2" type="success" size="small">已转码</el-tag>
          <el-tag v-else-if="scope.row.transcode_status === 3" type="danger" size="small">转码失败</el-tag>
          <el-tag v-else type="info" size="small">待转码</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <!-- 视频：只有转码成功才能复制，复制转码后的URL -->
          <el-button
            v-if="getFileCategory(scope.row.mime_type) === 'video'"
            type="primary"
            link
            :disabled="scope.row.transcode_status !== 2"
            @click="copyUrl(scope.row.transcode_url)"
          >
            复制链接
          </el-button>
          <!-- 非视频：直接复制原始URL -->
          <el-button
            v-else
            type="primary"
            link
            @click="copyUrl(scope.row.file_url)"
          >
            复制链接
          </el-button>
          <el-button type="danger" link @click="handleDeleteClick(scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />

    <!-- 删除确认对话框 -->
    <el-dialog v-model="deleteDialogVisible" width="500" align-center destroy-on-close>
      <template #header>删除资源</template>
      您已选中 [{{ idsToDelete.length }}] 项资源，删除后将无法恢复，是否确认删除？
      <template #footer>
        <el-button type="primary" @click="handleDelete">确定</el-button>
        <el-button @click="deleteDialogVisible = false">取消</el-button>
      </template>
    </el-dialog>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传资源"
      width="600"
      align-center
      destroy-on-close
      :close-on-click-modal="false"
    >
      <ResourceUpload @success="handleUploadSuccess" @cancel="uploadDialogVisible = false" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, VideoPlay, Headset, Folder } from '@element-plus/icons-vue'
import {
  getResourceList,
  deleteResource,
  formatFileSize,
  getFileCategory,
  type ResourceItem,
  type ResourceListRequest
} from '@/api/resource'
import ResourceUpload from './resource-upload.vue'

const multipleTableRef = ref()
const tableData = ref<ResourceItem[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedRows = ref<ResourceItem[]>([])
const deleteDialogVisible = ref(false)
const uploadDialogVisible = ref(false)
const idsToDelete = ref<number[]>([])

const route = useRoute()
const router = useRouter()

const listRequest = reactive<ResourceListRequest>({
  page: 1,
  page_size: 10,
  file_name: '',
  mime_type: ''
})

const mimeTypeOptions = [
  { value: '', label: '全部' },
  { value: 'image/jpeg', label: '图片 (JPEG)' },
  { value: 'image/png', label: '图片 (PNG)' },
  { value: 'image/gif', label: '图片 (GIF)' },
  { value: 'image/webp', label: '图片 (WebP)' },
  { value: 'video/mp4', label: '视频 (MP4)' },
  { value: 'audio/mpeg', label: '音频 (MP3)' },
  { value: 'application/pdf', label: 'PDF' },
  { value: 'application/zip', label: 'ZIP' }
]


const getResourceTableData = async () => {
  listRequest.page = page.value
  listRequest.page_size = pageSize.value

  const res = await getResourceList(listRequest)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        file_name: listRequest.file_name || undefined,
        mime_type: listRequest.mime_type || undefined,
        page: listRequest.page,
        page_size: listRequest.page_size
      }
    })
  }
}

const handleSelectionChange = (rows: ResourceItem[]) => {
  selectedRows.value = rows
}

const handleBulkDeleteClick = () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要删除的资源')
    return
  }
  idsToDelete.value = selectedRows.value.map((row) => row.id)
  deleteDialogVisible.value = true
}

const handleDeleteClick = (row: ResourceItem) => {
  idsToDelete.value = [row.id]
  deleteDialogVisible.value = true
}

const handleDelete = async () => {
  const res = await deleteResource({ ids: idsToDelete.value })
  if (res.code === 0) {
    ElMessage.success(res.msg || '删除成功')
    deleteDialogVisible.value = false
    getResourceTableData()
  }
}

const copyUrl = async (url: string) => {
  try {
    await navigator.clipboard.writeText(url)
    ElMessage.success('链接已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

const handleUploadSuccess = () => {
  uploadDialogVisible.value = false
  getResourceTableData()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  getResourceTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getResourceTableData()
}

onMounted(() => {
  listRequest.file_name = (route.query.file_name as string) || ''
  listRequest.mime_type = (route.query.mime_type as string) || ''
  page.value = Number(route.query.page) || 1
  pageSize.value = Number(route.query.page_size) || 10
})

watch(
  () => route.query,
  (newQuery) => {
    listRequest.file_name = (newQuery.file_name as string) || ''
    listRequest.mime_type = (newQuery.mime_type as string) || ''
    listRequest.page = Number(newQuery.page) || 1
    listRequest.page_size = Number(newQuery.page_size) || 10
  },
  { immediate: true }
)

nextTick(() => {
  getResourceTableData()
})
</script>

<style scoped lang="scss">
.resource-list {
  background-color: rgba(255, 255, 255, 0.7);
  padding: 10px;
  min-height: calc(100vh - 178px);

  .title {
    display: flex;

    .el-row {
      font-size: 24px;
    }

    .el-button-group {
      margin-left: auto;
      margin-top: auto;
      margin-bottom: auto;

      .el-button {
        margin-left: 12px;
      }
    }
  }

  .resource-list-request {
    padding-top: 10px;
    display: flex;

    .el-form {
      margin-left: auto;
    }
  }

  .el-table {
    border: 1px solid #dcdfe6;

    .preview-cell {
      position: relative;
      display: inline-block;
      width: 60px;
      height: 60px;

      .video-play-overlay {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 28px;
        height: 28px;
        background-color: rgba(0, 0, 0, 0.5);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .preview-fallback {
      width: 60px;
      height: 60px;
      background-color: #f5f7fa;
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .preview-icon {
      width: 60px;
      height: 60px;
      background-color: #f5f7fa;
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .el-pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
