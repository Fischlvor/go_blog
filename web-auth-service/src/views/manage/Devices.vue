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
      <el-row :gutter="16" class="devices-row">
        <el-col :span="12" v-for="deviceGroup in groupedDevices" :key="deviceGroup.deviceId" class="device-col">
          <el-card class="device-card" :class="{ 'current-device': deviceGroup.isCurrent }">
            <div class="device-header">
              <div class="device-icon">
                <el-icon size="20">
                  <Monitor v-if="deviceGroup.deviceType === 'web'" />
                  <Cellphone v-else />
                </el-icon>
              </div>
              <div class="device-info">
                <div class="device-name-row">
                  <div class="device-name">
                    {{ deviceGroup.deviceName }}
                    <el-tag v-if="deviceGroup.isCurrent" type="primary" size="small">当前</el-tag>
                  </div>
                  <el-button 
                    v-if="!deviceGroup.isCurrent"
                    type="danger" 
                    size="small"
                    @click="confirmKickAllApps(deviceGroup)"
                    :loading="kickingDevices.includes(deviceGroup.deviceId)"
                  >
                    踢出所有应用
                  </el-button>
                </div>
                <div class="device-meta">
                  <span class="meta-item">{{ deviceGroup.ipAddress }}</span>
                  <span class="meta-item">{{ formatTime(deviceGroup.lastActiveAt) }}</span>
                </div>
              </div>
            </div>
            
            <!-- 应用列表 - 紧凑布局 -->
            <div class="device-apps">
              <div class="apps-header" v-if="deviceGroup.apps.length > 1">
                <span class="apps-count">{{ deviceGroup.apps.length }} 个应用</span>
              </div>
              <div class="apps-list">
                <div v-for="app in deviceGroup.apps" :key="app.id" class="app-item">
                  <div class="app-info">
                    <span class="app-name">{{ app.app_name || app.app_key || '未知应用' }}</span>
                    <span class="app-time">{{ formatTime(app.last_active_at) }}</span>
                  </div>
                  <el-button 
                    type="danger" 
                    size="small"
                    @click="confirmKickDevice(app)"
                    :loading="kickingDevices.includes(app.device_id)"
                    v-if="!app.is_current"
                  >
                    踢出
                  </el-button>
                </div>
              </div>
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

// 按设备分组
const groupedDevices = computed(() => {
  const groups = new Map()
  
  devices.value.forEach(device => {
    const deviceId = device.device_id
    
    if (!groups.has(deviceId)) {
      groups.set(deviceId, {
        deviceId: deviceId,
        deviceName: device.device_name,
        deviceType: device.device_type,
        ipAddress: device.ip_address,
        lastActiveAt: device.last_active_at,
        isCurrent: device.is_current,
        apps: []
      })
    }
    
    const group = groups.get(deviceId)
    group.apps.push(device)
    
    // 更新设备级别的信息（使用最新的）
    if (new Date(device.last_active_at) > new Date(group.lastActiveAt)) {
      group.lastActiveAt = device.last_active_at
      group.ipAddress = device.ip_address
    }
    
    // 如果任何应用是当前设备，则整个设备组标记为当前
    if (device.is_current) {
      group.isCurrent = true
    }
  })
  
  // 转换为数组并按最后活跃时间排序
  return Array.from(groups.values()).sort((a, b) => 
    new Date(b.lastActiveAt) - new Date(a.lastActiveAt)
  )
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

// 确认踢出单个应用
const confirmKickDevice = (device) => {
  ElMessageBox.confirm(
    `确定要踢出应用 "${device.app_name}" 吗？`,
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

// 确认踢出设备的所有应用
const confirmKickAllApps = (deviceGroup) => {
  ElMessageBox.confirm(
    `确定要踢出设备 "${deviceGroup.deviceName}" 的所有应用吗？`,
    '确认操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    await kickAllApps(deviceGroup)
  }).catch(() => {
    // 用户取消
  })
}

// 踢出设备的所有应用
const kickAllApps = async (deviceGroup) => {
  kickingDevices.value.push(deviceGroup.deviceId)
  try {
    // 逐个踢出该设备的所有应用
    for (const app of deviceGroup.apps) {
      if (!app.is_current) { // 跳过当前设备
        await manageApi.kickDevice(app.device_id)
      }
    }
    ElMessage.success(`设备 "${deviceGroup.deviceName}" 的所有应用已踢出`)
    await fetchDevices() // 刷新列表
  } catch (error) {
    console.error('踢出设备应用失败:', error)
    ElMessage.error('踢出设备应用失败')
  } finally {
    kickingDevices.value = kickingDevices.value.filter(id => id !== deviceGroup.deviceId)
  }
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

/* 等高布局 */
.devices-row {
  display: flex;
  flex-wrap: wrap;
}

.device-col {
  display: flex;
  margin-bottom: 16px;
}

.device-card {
  width: 100%;
  display: flex;
  flex-direction: column;
  transition: all 0.3s;
}

.device-card .el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.device-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.current-device {
  border: 2px solid #409eff;
}

.device-header {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}

.device-icon {
  color: #409eff;
  display: flex;
  align-items: center;
}

.device-info {
  flex: 1;
}

.device-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.device-name {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 6px;
}

.device-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.meta-item {
  display: inline-block;
}

/* 应用列表样式 - 紧凑布局 */
.device-apps {
  margin-top: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.apps-header {
  margin-bottom: 8px;
}

.apps-count {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.apps-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 180px;
  overflow-y: auto;
}

/* 美化滚动条 */
.apps-list::-webkit-scrollbar {
  width: 4px;
}

.apps-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 2px;
}

.apps-list::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 2px;
}

.apps-list::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

.app-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  background: #f8f9fa;
  border-radius: 4px;
  border-left: 3px solid #e4e7ed;
  gap: 12px;
}

.app-item:hover {
  background: #f0f2f5;
}

.app-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
}

.app-time {
  font-size: 11px;
  color: #c0c4cc;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}
</style>
