<template>
  <div class="devices-page">
    <div class="page-header">
      <h3>设备管理</h3>
      <el-button @click="refreshDevices" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <div class="devices-stats">
      <el-card>
        <div class="stats-content">
          <div class="stat-item">
            <div class="stat-number">{{ devices.length }}</div>
            <div class="stat-label">活跃设备</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ currentDeviceCount }}</div>
            <div class="stat-label">当前设备</div>
          </div>
        </div>
      </el-card>
    </div>

    <div class="devices-list">
      <el-row :gutter="16">
        <el-col :span="12" v-for="device in devices" :key="device.device_id">
          <el-card class="device-card" :class="{ 'current-device': device.is_current }">
            <div class="device-header">
              <div class="device-icon">
                <el-icon size="24">
                  <Monitor v-if="device.device_type === 'web'" />
                  <Cellphone v-else />
                </el-icon>
              </div>
              <div class="device-info">
                <div class="device-name">
                  {{ device.device_name }}
                  <el-tag v-if="device.is_current" type="primary" size="small">当前设备</el-tag>
                </div>
                <div class="device-details">
                  <div class="detail-item">
                    <span class="label">设备ID:</span>
                    <span class="value">{{ device.device_id }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="label">IP地址:</span>
                    <span class="value">{{ device.ip_address }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="label">最后活跃:</span>
                    <span class="value">{{ formatTime(device.last_active_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="device-actions" v-if="!device.is_current">
              <el-button 
                type="danger" 
                size="small"
                @click="confirmKickDevice(device)"
                :loading="kickingDevices.includes(device.device_id)"
              >
                踢出设备
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div v-if="devices.length === 0 && !loading" class="empty-state">
      <el-empty description="暂无活跃设备" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, Cellphone, Refresh } from '@element-plus/icons-vue'
import { manageApi } from '../../api/manage'

const devices = ref([])
const loading = ref(false)
const kickingDevices = ref([])

const currentDeviceCount = computed(() => {
  return devices.value.filter(device => device.is_current).length
})

// 获取设备列表
const fetchDevices = async () => {
  loading.value = true
  try {
    const res = await manageApi.getDevices()
    console.log('设备列表响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      devices.value = res.data.data || []
      console.log('设备列表数据:', devices.value) // 调试信息
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || '获取设备列表失败')
    }
  } catch (error) {
    console.error('获取设备列表失败:', error)
    ElMessage.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新设备列表
const refreshDevices = () => {
  fetchDevices()
}

// 确认踢出设备
const confirmKickDevice = (device) => {
  ElMessageBox.confirm(
    `确定要踢出设备 "${device.device_name}" 吗？`,
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    await kickDevice(device.device_id)
  }).catch(() => {
    // 用户取消
  })
}

// 踢出设备
const kickDevice = async (deviceId) => {
  kickingDevices.value.push(deviceId)
  try {
    const res = await manageApi.kickDevice(deviceId)
    console.log('踢出设备响应:', res) // 调试信息
    if (res.data && res.data.code === 0) {
      ElMessage.success('设备已踢出')
      await fetchDevices() // 刷新列表
    } else {
      console.error('响应格式错误:', res) // 调试信息
      ElMessage.error(res.data?.message || '踢出设备失败')
    }
  } catch (error) {
    console.error('踢出设备失败:', error)
    ElMessage.error('踢出设备失败')
  } finally {
    kickingDevices.value = kickingDevices.value.filter(id => id !== deviceId)
  }
}

// 格式化时间
const formatTime = (timeStr) => {
  const time = new Date(timeStr)
  const now = new Date()
  const diff = now - time
  
  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) { // 1天内
    return `${Math.floor(diff / 3600000)}小时前`
  } else {
    return time.toLocaleDateString() + ' ' + time.toLocaleTimeString()
  }
}

onMounted(() => {
  fetchDevices()
})
</script>

<style scoped>
.devices-page {
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

.devices-stats {
  margin-bottom: 20px;
}

.stats-content {
  display: flex;
  gap: 40px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.device-card {
  margin-bottom: 16px;
  transition: all 0.3s;
}

.device-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.current-device {
  border: 2px solid #409eff;
}

.device-header {
  display: flex;
  gap: 12px;
}

.device-icon {
  color: #409eff;
  display: flex;
  align-items: flex-start;
  padding-top: 4px;
}

.device-info {
  flex: 1;
}

.device-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.device-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item {
  font-size: 14px;
  display: flex;
  gap: 8px;
}

.label {
  color: #909399;
  min-width: 60px;
}

.value {
  color: #606266;
  word-break: break-all;
}

.device-actions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  text-align: right;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}
</style>
