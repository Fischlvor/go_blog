<template>
  <el-dialog 
    v-model="visible" 
    title="雪碧图生成进度" 
    width="500px" 
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="status !== 'running'"
    :destroy-on-close="false"
  >
    <div class="sprite-generation-progress">
      <!-- 总组数进度 -->
      <div class="progress-section" v-if="showGroupProgress">
        <div class="progress-label">
          <span class="label-text">总组数进度</span>
          <span class="label-percentage">{{ groupProgress }}%</span>
        </div>
        <el-progress
          :percentage="groupProgress"
          :status="status === 'completed' ? 'success' : status === 'failed' ? 'exception' : undefined"
        />
      </div>

      <!-- 当前组的雪碧图进度 -->
      <div class="progress-section" v-if="showSpriteProgress">
        <div class="progress-label">
          <span class="label-text">当前组雪碧图进度</span>
          <span class="label-percentage">{{ spriteProgress }}%</span>
        </div>
        <el-progress
          :percentage="spriteProgress"
          :status="status === 'completed' ? 'success' : status === 'failed' ? 'exception' : undefined"
        />
      </div>

      <!-- 消息显示 -->
      <div class="task-message">{{ message }}</div>

      <!-- 完成/失败提示 -->
      <div v-if="status === 'completed'" class="task-result">
        <el-alert title="任务完成" type="success" :closable="false" />
      </div>
      <div v-if="status === 'failed'" class="task-result">
        <el-alert title="任务失败" type="error" :closable="false" />
      </div>
    </div>

    <template #footer>
      <el-button @click="visible = false" :disabled="status === 'running'">
        {{ status === 'running' ? '生成中...' : '关闭' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  modelValue: boolean
  status: 'running' | 'completed' | 'failed' | ''
  message: string
  groupProgress?: number
  spriteProgress?: number
  showGroupProgress?: boolean
  showSpriteProgress?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = withDefaults(defineProps<Props>(), {
  groupProgress: 0,
  spriteProgress: 0,
  showGroupProgress: true,
  showSpriteProgress: true
})

const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

<style scoped lang="scss">
.sprite-generation-progress {
  .progress-section {
    margin-bottom: 24px;

    .progress-label {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      .label-text {
        color: #333;
        font-weight: 500;
      }

      .label-percentage {
        color: #409eff;
        font-weight: bold;
        font-size: 14px;
      }
    }
  }

  .task-message {
    margin: 20px 0;
    color: #666;
    text-align: center;
    min-height: 20px;
  }

  .task-result {
    margin-top: 20px;
  }
}
</style>
