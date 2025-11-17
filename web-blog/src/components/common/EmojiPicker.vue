<template>
  <div class="emoji-picker">
    <!-- Emoji触发按钮 -->
    <button 
      @click="togglePicker" 
      class="emoji-trigger-btn"
      :class="{ active: isVisible }"
    >
      😀
    </button>

    <!-- Emoji面板 -->
    <div 
      v-if="isVisible" 
      class="emoji-panel"
      @click.stop
    >
      <div class="emoji-panel-header">
        <span>选择表情</span>
        <button @click="closePicker" class="close-btn">×</button>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>加载表情中...</span>
      </div>

      <!-- Emoji网格 -->
      <div v-else class="emoji-grid" ref="emojiGrid">
        <div 
          v-for="emoji in visibleEmojis" 
          :key="emoji.key"
          class="emoji-item"
          :class="[
            'emoji',
            `emoji-sprite-${emoji.spriteGroup}`,
            `emoji-${emoji.key}`
          ]"
          :title="emoji.key"
          @click="selectEmoji(emoji)"
        ></div>
        
        <!-- 加载更多指示器 -->
        <div v-if="hasMore" ref="loadMoreTrigger" class="load-more-trigger">
          <div class="loading-spinner"></div>
        </div>
      </div>
    </div>

    <!-- 遮罩层 -->
    <div 
      v-if="isVisible" 
      class="emoji-overlay"
      @click="closePicker"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { getAllEmojis, type EmojiInfo } from '@/utils/emojiParser'

interface Props {
  modelValue?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'select', emoji: EmojiInfo): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态管理
const isVisible = ref(false)
const loading = ref(false)
const allEmojis = ref<EmojiInfo[]>([])
const visibleEmojis = ref<EmojiInfo[]>([])
const currentPage = ref(0)
const pageSize = 48 // 每页显示48个emoji (6x8)
const emojiGrid = ref<HTMLElement>()
const loadMoreTrigger = ref<HTMLElement>()
let intersectionObserver: IntersectionObserver | null = null

// 计算属性
const hasMore = computed(() => {
  return visibleEmojis.value.length < allEmojis.value.length
})

// 切换面板显示
const togglePicker = () => {
  if (isVisible.value) {
    closePicker()
  } else {
    openPicker()
  }
}

// 打开面板
const openPicker = async () => {
  isVisible.value = true
  
  if (allEmojis.value.length === 0) {
    await loadEmojis()
  }
  
  // 延迟设置观察器，确保DOM已渲染
  await nextTick()
  setupIntersectionObserver()
}

// 关闭面板
const closePicker = () => {
  isVisible.value = false
  cleanupIntersectionObserver()
}

// 懒加载emoji数据
const loadEmojis = async () => {
  loading.value = true
  
  try {
    // 模拟异步加载（实际项目中可能从API获取）
    await new Promise(resolve => setTimeout(resolve, 300))
    
    allEmojis.value = await getAllEmojis()
    loadMoreEmojis()
  } catch (error) {
    console.error('加载emoji失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载更多emoji
const loadMoreEmojis = () => {
  const startIndex = currentPage.value * pageSize
  const endIndex = Math.min(startIndex + pageSize, allEmojis.value.length)
  
  const newEmojis = allEmojis.value.slice(startIndex, endIndex)
  visibleEmojis.value.push(...newEmojis)
  
  currentPage.value++
}

// 设置交叉观察器（用于无限滚动）
const setupIntersectionObserver = () => {
  if (!loadMoreTrigger.value) return
  
  intersectionObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting && hasMore.value && !loading.value) {
          loadMoreEmojis()
        }
      })
    },
    {
      rootMargin: '50px'
    }
  )
  
  intersectionObserver.observe(loadMoreTrigger.value)
}

// 清理观察器
const cleanupIntersectionObserver = () => {
  if (intersectionObserver) {
    intersectionObserver.disconnect()
    intersectionObserver = null
  }
}

// 选择emoji
const selectEmoji = (emoji: EmojiInfo) => {
  const emojiText = `:emoji:${emoji.key}:`
  emit('update:modelValue', (props.modelValue || '') + emojiText)
  emit('select', emoji)
  closePicker()
}

// 监听面板外点击
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (isVisible.value && !target.closest('.emoji-picker')) {
    closePicker()
  }
}

// 生命周期
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  cleanupIntersectionObserver()
})

// 监听可见性变化
watch(isVisible, (newVal) => {
  if (!newVal) {
    // 重置状态
    currentPage.value = 0
    visibleEmojis.value = []
  }
})
</script>

<style scoped>
.emoji-picker {
  position: relative;
  display: inline-block;
}

.emoji-trigger-btn {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.2s;
}

.emoji-trigger-btn:hover,
.emoji-trigger-btn.active {
  border-color: #409eff;
  background: #f0f9ff;
}

.emoji-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.1);
}

.emoji-panel {
  position: absolute;
  top: 100%;
  left: 0;
  width: 320px;
  max-height: 400px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  overflow: hidden;
}

.emoji-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.emoji-panel-header span {
  font-weight: 500;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #666;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #409eff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.emoji-grid {
  padding: 16px;
  max-height: 320px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.emoji-item {
  width: 32px;
  height: 32px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-item:hover {
  background: #f0f9ff;
  transform: scale(1.1);
}

.load-more-trigger {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  padding: 16px;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .emoji-panel {
    width: 280px;
  }
  
  .emoji-grid {
    grid-template-columns: repeat(6, 1fr);
  }
}
</style>
