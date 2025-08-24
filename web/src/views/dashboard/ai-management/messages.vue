<template>
  <div class="ai-message-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI消息管理</span>
        </div>
      </template>

      <!-- 搜索区域 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="会话ID">
          <el-input v-model="searchForm.sessionId" placeholder="请输入会话ID" clearable />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
            <el-option label="用户" value="user" />
            <el-option label="AI助手" value="assistant" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="session_id" label="会话ID" width="100" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'user' ? 'primary' : 'success'">
              {{ row.role === 'user' ? '用户' : 'AI助手' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="tokens" label="Token数" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleView(row)">查看详情</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 消息详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="消息详情"
      width="800px"
    >
      <div v-if="messageDetail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="消息ID">{{ messageDetail.id }}</el-descriptions-item>
          <el-descriptions-item label="会话ID">{{ messageDetail.session_id }}</el-descriptions-item>
          <el-descriptions-item label="角色">
            <el-tag :type="messageDetail.role === 'user' ? 'primary' : 'success'">
              {{ messageDetail.role === 'user' ? '用户' : 'AI助手' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Token数">{{ messageDetail.tokens }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(messageDetail.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(messageDetail.updated_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="content-section">
          <h3>消息内容</h3>
          <div class="message-content">
            <pre>{{ messageDetail.content }}</pre>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { aiMessageApi } from '@/api/ai-management'

// 响应式数据
const loading = ref(false)
const detailDialogVisible = ref(false)
const messageDetail = ref<any>(null)

// 搜索表单
const searchForm = reactive({
  sessionId: '',
  role: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

// 表格数据
const tableData = ref([])

// 获取列表数据
const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      session_id: searchForm.sessionId,
      role: searchForm.role
    }
    const res = await aiMessageApi.list(params)
    console.log('消息列表数据:', res.data.list)
    tableData.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('获取AI消息列表失败:', error)
    ElMessage.error('获取AI消息列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  getList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    sessionId: '',
    role: ''
  })
  pagination.page = 1
  getList()
}

// 查看详情
const handleView = async (row: any) => {
  try {
    const res = await aiMessageApi.get(row.id)
    messageDetail.value = res.data
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取消息详情失败:', error)
    ElMessage.error('获取消息详情失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条消息吗？删除后将无法恢复！', '提示', {
      type: 'warning'
    })
    
    await aiMessageApi.delete(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除AI消息失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 分页大小改变
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  getList()
}

// 当前页改变
const handleCurrentChange = (page: number) => {
  pagination.page = page
  getList()
}

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style scoped>
.ai-message-management {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.content-section {
  margin-top: 20px;
}

.message-content {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  margin-top: 10px;
}

.message-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;
}
</style> 