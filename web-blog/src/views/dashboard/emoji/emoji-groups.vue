<template>
  <div class="emoji-groups">
    <div class="page-header">
      <h2>表情组管理</h2>
      <el-button type="primary" @click="showCreateDialog = true" :icon="Plus">
        新建表情组
      </el-button>
    </div>

    <!-- 表情组列表 -->
    <el-table :data="groupList" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="group_name" label="组名" width="200">
        <template #default="{ row }">
          <el-tag type="primary" size="large">{{ row.group_name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="300" />
      <el-table-column prop="emoji_count" label="表情数量" width="120">
        <template #default="{ row }">
          <el-badge :value="row.emoji_count" class="item">
            <el-icon><Picture /></el-icon>
          </el-badge>
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="editGroup(row)" :icon="Edit">
            编辑
          </el-button>
          <el-button
            type="warning"
            size="small"
            @click="regenerateGroupSprites(row)"
            :loading="regeneratingGroupKey === row.group_key"
            :icon="Refresh"
            :disabled="row.emoji_count === 0"
          >
            重生成
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click="deleteGroup(row)"
            :icon="Delete"
            :disabled="row.emoji_count > 0"
          >
            删除
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingGroup ? '编辑表情组' : '新建表情组'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="groupForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="组名" prop="group_name">
          <el-input
            v-model="groupForm.group_name"
            placeholder="请输入表情组名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="groupForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入表情组描述"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number
            v-model="groupForm.sort_order"
            :min="0"
            :max="999"
            placeholder="排序权重，数字越小越靠前"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelEdit">取消</el-button>
          <el-button type="primary" @click="saveGroup" :loading="saving">
            {{ editingGroup ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Picture, Refresh } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import SpriteGenerationProgress from '@/components/SpriteGenerationProgress.vue'
import { 
  getEmojiGroups, 
  createEmojiGroup, 
  updateEmojiGroup, 
  deleteEmojiGroup as apiDeleteEmojiGroup,
  regenerateSpritesStream,
  type EmojiGroup,
  type EmojiGroupCreateRequest,
  type EmojiGroupUpdateRequest
} from '@/api/emoji'

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const showCreateDialog = ref(false)
const editingGroup = ref<EmojiGroup | null>(null)
const groupList = ref<EmojiGroup[]>([])
const formRef = ref<FormInstance>()
const regeneratingGroupKey = ref<string>('')

// 任务进度相关
const showTaskDialog = ref(false)
const taskStatus = ref<'running' | 'completed' | 'failed' | ''>('')
const taskMessage = ref('')
const totalGroups = ref(0)
const processedGroups = ref(0)
const currentGroupSprites = ref(0)
const currentGroupTotalSprites = ref(0)

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

// 表单数据
const groupForm = reactive({
  group_name: '',
  description: '',
  sort_order: 0
})

// 表单验证规则
const formRules: FormRules = {
  group_name: [
    { required: true, message: '请输入表情组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述不能超过 200 个字符', trigger: 'blur' }
  ],
  sort_order: [
    { type: 'number', message: '排序必须是数字', trigger: 'blur' }
  ]
}

// 获取表情组列表
const fetchGroupList = async () => {
  loading.value = true
  try {
    const result = await getEmojiGroups()
    if (result.code === 0) {
      groupList.value = result.data || []
    } else {
      ElMessage.error(result.msg || '获取表情组列表失败')
    }
  } catch (error) {
    console.error('获取表情组列表失败:', error)
    ElMessage.error('获取表情组列表失败')
  } finally {
    loading.value = false
  }
}

// 编辑表情组
const editGroup = (group: EmojiGroup) => {
  editingGroup.value = group
  groupForm.group_name = group.group_name
  groupForm.description = group.description
  groupForm.sort_order = group.sort_order
  showCreateDialog.value = true
}

// 取消编辑
const cancelEdit = () => {
  showCreateDialog.value = false
  editingGroup.value = null
  resetForm()
}

// 重置表单
const resetForm = () => {
  groupForm.group_name = ''
  groupForm.description = ''
  groupForm.sort_order = 0
  formRef.value?.clearValidate()
}

// 保存表情组
const saveGroup = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    const data: EmojiGroupCreateRequest | EmojiGroupUpdateRequest = {
      group_name: groupForm.group_name,
      description: groupForm.description,
      sort_order: groupForm.sort_order
    }

    let result
    if (editingGroup.value) {
      result = await updateEmojiGroup(editingGroup.value.id, data)
    } else {
      result = await createEmojiGroup(data)
    }

    if (result.code === 0) {
      ElMessage.success(editingGroup.value ? '更新成功' : '创建成功')
      showCreateDialog.value = false
      editingGroup.value = null
      resetForm()
      fetchGroupList()
    } else {
      ElMessage.error(result.msg || '操作失败')
    }
  } catch (error) {
    console.error('保存表情组失败:', error)
    ElMessage.error('操作失败')
  } finally {
    saving.value = false
  }
}

// 删除表情组
const deleteGroup = async (group: EmojiGroup) => {
  if (group.emoji_count > 0) {
    ElMessage.warning('该表情组内还有表情，无法删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除表情组 "${group.group_name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const result = await apiDeleteEmojiGroup(group.id)
    if (result.code === 0) {
      ElMessage.success('删除成功')
      fetchGroupList()
    } else {
      ElMessage.error(result.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除表情组失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 重新生成表情组的雪碧图
const regenerateGroupSprites = async (group: EmojiGroup) => {
  try {
    await ElMessageBox.confirm(
      `确定要重新生成表情组 "${group.group_name}" 的雪碧图吗？这可能需要一些时间。`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    regeneratingGroupKey.value = group.group_key
    
    // 显示进度条对话框
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
    await regenerateSpritesStream([group.group_key], (event: string, data: any) => {
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
          fetchGroupList()
          break
        case 'error':
          ElMessage.error(data.message)
          taskStatus.value = 'failed'
          taskMessage.value = `错误: ${data.message}`
          break
      }
    })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重新生成雪碧图失败:', error)
      ElMessage.error('操作失败')
      taskStatus.value = 'failed'
      taskMessage.value = '操作失败'
    }
  } finally {
    regeneratingGroupKey.value = ''
  }
}

// 工具函数
const formatTime = (time: string) => {
  return new Date(time).toLocaleString()
}

// 初始化
onMounted(() => {
  fetchGroupList()
})
</script>

<style scoped lang="scss">
.emoji-groups {
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
  }

  .item {
    margin-right: 10px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
