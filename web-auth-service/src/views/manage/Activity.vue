<template>
  <div class="activity-page">
    <div class="page-header">
      <h3>操作日志</h3>
      <el-button @click="refreshLogs" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card">
      <el-form :model="queryParams" inline>
        <el-form-item label="操作类型">
          <el-select v-model="queryParams.action" placeholder="全部" clearable style="width: 120px">
            <el-option label="登录" value="login" />
            <el-option label="QQ登录" value="qq_login" />
            <el-option label="静默登录" value="silent_login" />
            <el-option label="登出" value="logout" />
            <el-option label="SSO退出" value="sso_logout" />
            <el-option label="手动踢出" value="manual_kick" />
            <el-option label="自动踢出" value="auto_kick" />
            <el-option label="退出所有" value="logout_all" />
            <el-option label="设备过期" value="expired" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchLogs">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 日志列表 -->
    <el-card class="logs-card">
      <el-table 
        :data="logs" 
        v-loading="loading"
        empty-text="暂无日志记录"
      >
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="action" label="操作" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="device_id" label="设备ID" width="120">
          <template #default="{ row }">
            <span class="device-id">{{ row.device_id }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="message" label="描述" min-width="200" />
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { manageApi } from '../../api/manage'

const logs = ref([])
const loading = ref(false)
const total = ref(0)

const queryParams = ref({
  page: 1,
  page_size: 20,
  action: '',
  start_time: '',
  end_time: ''
})

// 获取日志列表
const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await manageApi.getLogs(queryParams.value)
    console.log('日志响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      logs.value = res.data.data.list || []
      total.value = res.data.data.total || 0
      console.log('日志数据:', logs.value) // 调试信息
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || '获取日志失败')
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error('获取日志失败')
  } finally {
    loading.value = false
  }
}

// 刷新日志
const refreshLogs = () => {
  fetchLogs()
}

// 搜索日志
const searchLogs = () => {
  queryParams.value.page = 1
  fetchLogs()
}

// 重置查询
const resetQuery = () => {
  queryParams.value = {
    page: 1,
    page_size: 20,
    action: '',
    start_time: '',
    end_time: ''
  }
  fetchLogs()
}

// 分页大小改变
const handleSizeChange = (val) => {
  queryParams.value.page_size = val
  queryParams.value.page = 1
  fetchLogs()
}

// 当前页改变
const handleCurrentChange = (val) => {
  queryParams.value.page = val
  fetchLogs()
}

// 获取操作类型样式
const getActionType = (action) => {
  const typeMap = {
    'login': 'success',
    'qq_login': 'success',
    'silent_login': 'info',
    'logout': 'warning',
    'sso_logout': 'warning',
    'manual_kick': 'danger',
    'auto_kick': 'danger',
    'logout_all': 'danger',
    'expired': 'danger'
  }
  return typeMap[action] || 'info'
}

// 获取操作标签
const getActionLabel = (action) => {
  const labelMap = {
    'login': '登录',
    'qq_login': 'QQ登录',
    'silent_login': '静默登录',
    'logout': '登出',
    'sso_logout': 'SSO退出',
    'manual_kick': '手动踢出',
    'auto_kick': '自动踢出',
    'logout_all': '退出所有',
    'expired': '设备过期'
  }
  return labelMap[action] || action
}

// 格式化日期时间
const formatDateTime = (timeStr) => {
  const time = new Date(timeStr)
  return time.toLocaleDateString() + ' ' + time.toLocaleTimeString()
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.activity-page {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
  font-weight: 500;
}

.filter-card {
  margin-bottom: 20px;
}

.logs-card {
  margin-bottom: 20px;
}

.device-id {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #606266;
}

.pagination-wrapper {
  margin-top: 20px;
  text-align: right;
}
</style>
