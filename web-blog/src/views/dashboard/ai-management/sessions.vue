<template>
  <div class="ai-session-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI会话管理</span>
        </div>
      </template>

      <!-- 搜索区域 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户UUID">
          <el-input v-model="searchForm.userUuid" placeholder="请输入用户UUID" clearable />
        </el-form-item>
        <el-form-item label="模型">
          <el-input v-model="searchForm.model" placeholder="请输入模型名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" width="80">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 5px;">
              <el-popover width="280" v-if="row.user">
                <template #reference>
                  <el-avatar :src="row.user.avatar" alt=""/>
                </template>
                <template #default>
                  <user-card :uuid="''"
                             :user-card-info="{uuid:row.user.uuid,username:row.user.username,avatar:row.user.avatar,address:row.user.address,signature:row.user.signature}"/>
                </template>
              </el-popover>
              <el-avatar v-else alt=""/>
              <span v-if="!row.user" style="font-size: 12px; color: #999;">无用户信息</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="user_uuid" label="用户UUID" width="100" />
        <el-table-column prop="title" label="会话标题" show-overflow-tooltip />
        <el-table-column prop="model" label="AI模型" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
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

    <!-- 会话详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="会话详情"
      width="800px"
    >
      <div v-if="sessionDetail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="会话ID">{{ sessionDetail.id }}</el-descriptions-item>
          <el-descriptions-item label="用户">
            <div v-if="sessionDetail.user" style="display: flex; align-items: center; gap: 10px;">
              <el-avatar :src="sessionDetail.user.avatar" alt=""/>
              <span>{{ sessionDetail.user.username }}</span>
            </div>
            <span v-else>{{ sessionDetail.user_uuid }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="会话标题">{{ sessionDetail.title }}</el-descriptions-item>
          <el-descriptions-item label="AI模型">{{ sessionDetail.model }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(sessionDetail.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(sessionDetail.updated_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="messages-section">
          <h3>消息列表</h3>
          <el-table :data="sessionDetail.messages" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="role" label="角色" width="100">
              <template #default="{ row }">
                <el-tag :type="row.role === 'user' ? 'primary' : 'success'">
                  {{ row.role === 'user' ? '用户' : 'AI助手' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="content" label="内容" show-overflow-tooltip />
            <el-table-column prop="tokens" label="Token数" width="100" />
            <el-table-column prop="created_at" label="时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { aiSessionApi } from '@/api/ai-management'
import UserCard from '@/components/widgets/UserCard.vue'

// 响应式数据
const loading = ref(false)
const detailDialogVisible = ref(false)
const sessionDetail = ref<any>(null)

// 搜索表单
const searchForm = reactive({
  userUuid: '',
  model: ''
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
      user_uuid: searchForm.userUuid,
      model: searchForm.model
    }
    const res = await aiSessionApi.list(params)
    console.log('会话列表数据:', res.data.list)
    tableData.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('获取AI会话列表失败:', error)
    ElMessage.error('获取AI会话列表失败')
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
    userUuid: '',
    model: ''
  })
  pagination.page = 1
  getList()
}

// 查看详情
const handleView = async (row: any) => {
  try {
    const res = await aiSessionApi.get(row.id)
    sessionDetail.value = res.data
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取会话详情失败:', error)
    ElMessage.error('获取会话详情失败')
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这个会话吗？删除后将无法恢复！', '提示', {
      type: 'warning'
    })
    
    await aiSessionApi.delete(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除AI会话失败:', error)
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
.ai-session-management {
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

.messages-section {
  margin-top: 20px;
}

.messages-section h3 {
  margin-bottom: 15px;
  color: #333;
}
</style> 