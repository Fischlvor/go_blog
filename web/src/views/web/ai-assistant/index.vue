<template>
  <div class="ai-assistant-page">
    <div class="main-container">
      <div class="chat-container">
        <!-- 会话选择器 -->
        <div class="session-header">
          <div class="session-selector">
            <el-select 
              v-model="currentSessionId" 
              placeholder="选择会话" 
              style="width: 200px;"
              @change="loadSessionMessages"
            >
              <el-option
                v-for="session in sessions"
                :key="session.id"
                :label="session.title"
                :value="session.id"
              />
            </el-select>
            <el-button 
              type="primary" 
              @click="createNewSession"
              style="margin-left: 10px;"
            >
              新建会话
            </el-button>
          </div>
        </div>

        <!-- 消息列表 -->
        <div class="messages-container" ref="messagesContainer">
          <div 
            v-for="message in messages" 
            :key="message.id"
            class="message"
            :class="message.role"
          >
            <div class="message-content">
              <!-- 用户消息使用普通文本显示 -->
              <div v-if="message.role === 'user'" class="message-text" v-html="formatMessage(message.content)"></div>
              <!-- AI助手消息使用Markdown渲染 -->
              <MdPreview 
                v-else
                :modelValue="message.content" 
                :preview-only="true"
                :show-code-row-number="false"
                :preview-theme="'github'"
                :code-theme="'github'"
              />
              <div class="message-time">{{ formatTime(message.created_at) }}</div>
            </div>
          </div>
          
          <!-- 加载中状态 -->
          <div v-if="loading" class="message assistant">
            <div class="message-content">
              <div class="message-text">
                <el-icon class="is-loading"><Loading /></el-icon>
                正在思考中...
              </div>
            </div>
          </div>
          
          <!-- 流式响应状态 -->
          <div v-if="streaming" class="message assistant">
            <div class="message-content">
              <div class="message-text">
                <el-icon class="is-loading"><Loading /></el-icon>
                AI正在回复中...
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-area">
          <el-input
            v-model="inputMessage"
            placeholder="输入您的问题..."
            type="textarea"
            :rows="4"
            @keydown.enter.prevent="sendMessage"
            :disabled="loading || !currentSessionId"
          />
          <div class="input-actions">
            <el-button 
              type="primary" 
              @click="sendMessage"
              :loading="loading"
              :disabled="!inputMessage.trim() || !currentSessionId"
              size="large"
            >
              发送
            </el-button>
            <el-button 
              @click="clearMessages"
              :disabled="!currentSessionId"
              size="large"
            >
              清空对话
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 新建会话对话框 -->
    <el-dialog
      v-model="showNewSessionDialog"
      title="新建会话"
      width="400px"
    >
      <el-form :model="newSessionForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="newSessionForm.title" placeholder="请输入会话标题" />
        </el-form-item>
        <el-form-item label="模型">
          <el-select v-model="newSessionForm.model" placeholder="选择AI模型" style="width: 100%;">
            <el-option
              v-for="model in availableModels"
              :key="model.name"
              :label="model.display_name"
              :value="model.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNewSessionDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateSession">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { aiChatApi } from '@/api/ai-chat'
import { MdPreview } from 'md-editor-v3'

// 响应式数据
const loading = ref(false)
const streaming = ref(false)
const inputMessage = ref('')
const messages = ref<any[]>([])
const sessions = ref<any[]>([])
const currentSessionId = ref<number | null>(null)
const showNewSessionDialog = ref(false)
const messagesContainer = ref<HTMLElement>()

// 新建会话表单
const newSessionForm = ref({
  title: '',
  model: 'gpt-3.5-turbo'
})

// 可用的AI模型
const availableModels = ref([
  { name: 'gpt-3.5-turbo', display_name: 'GPT-3.5 Turbo' },
  { name: 'gpt-4', display_name: 'GPT-4' },
  { name: 'claude-3-sonnet', display_name: 'Claude 3 Sonnet' }
])

// 页面加载时获取会话列表
onMounted(async () => {
  await loadSessions()
})

// 加载会话列表
const loadSessions = async () => {
  try {
    const response = await aiChatApi.getSessions()
    sessions.value = response.data.list || []
    if (sessions.value.length > 0 && !currentSessionId.value) {
      currentSessionId.value = sessions.value[0].id
      await loadSessionMessages()
    }
  } catch (error) {
    console.error('加载会话列表失败:', error)
    ElMessage.error('加载会话列表失败')
  }
}

// 加载会话消息
const loadSessionMessages = async () => {
  if (!currentSessionId.value) return
  
  try {
    const response = await aiChatApi.getMessages(currentSessionId.value)
    messages.value = response.data.list || []
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('加载消息失败:', error)
    ElMessage.error('加载消息失败')
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim() || !currentSessionId.value) return
  
  const content = inputMessage.value
  inputMessage.value = ''
  loading.value = true
  streaming.value = false
  
  // 添加用户消息到UI
  messages.value.push({
    id: Date.now(),
    role: 'user',
    content,
    created_at: new Date().toISOString()
  })
  
  await nextTick()
  scrollToBottom()
  
  // 创建AI回复消息占位符
  const aiMessageId = Date.now() + 1
  const aiMessageIndex = messages.value.length
  messages.value.push({
    id: aiMessageId,
    role: 'assistant',
    content: '',
    created_at: new Date().toISOString()
  })
  
  await nextTick()
  scrollToBottom()
  
  try {
    loading.value = false
    streaming.value = true
    
    // 使用流式API
    await aiChatApi.sendMessageStream({
      session_id: currentSessionId.value,
      content
    }, {
      onData: (data) => {
        if (data.content !== undefined) {
          messages.value[aiMessageIndex].content += data.content
          nextTick(() => {
            scrollToBottom()
          })
        }
      },
      onComplete: (data) => {
        if (data.message_id) {
          messages.value[aiMessageIndex].id = data.message_id
        }
      },
      onError: (error) => {
        console.error('流式请求失败:', error)
        ElMessage.error('发送消息失败')
      }
    })
  } catch (error) {
    ElMessage.error('发送消息失败')
    console.error('发送消息失败:', error)
    messages.value.splice(aiMessageIndex, 1)
  } finally {
    loading.value = false
    streaming.value = false
  }
}

// 新建会话
const createNewSession = () => {
  newSessionForm.value = {
    title: '',
    model: 'gpt-3.5-turbo'
  }
  showNewSessionDialog.value = true
}

// 确认新建会话
const confirmCreateSession = async () => {
  if (!newSessionForm.value.title.trim()) {
    ElMessage.warning('请输入会话标题')
    return
  }
  
  try {
    const response = await aiChatApi.createSession({
      title: newSessionForm.value.title,
      model: newSessionForm.value.model
    })
    
    ElMessage.success('会话创建成功')
    showNewSessionDialog.value = false
    
    // 重新加载会话列表
    await loadSessions()
    
    // 切换到新创建的会话
    currentSessionId.value = response.data.id
    messages.value = []
  } catch (error) {
    console.error('创建会话失败:', error)
    ElMessage.error('创建会话失败')
  }
}

// 清空对话
const clearMessages = () => {
  messages.value = []
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 格式化消息内容 - 用户消息的普通文本处理
const formatMessage = (content: string) => {
  if (!content) return ''
  // 处理换行符
  return content.replace(/\n/g, '<br>')
}

// 格式化时间
const formatTime = (time: string) => {
  return new Date(time).toLocaleTimeString()
}

// 监听消息变化，自动滚动
watch(messages, () => {
  nextTick(() => {
    scrollToBottom()
  })
})
</script>

<style scoped>
.ai-assistant-page {
  min-height: 100vh;
  background: transparent;
  padding-top: 20px;
}

.main-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: calc(100vh - 100px);
}

.chat-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
}

.session-header {
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
  background: #f8f9fa;
}

.session-selector {
  display: flex;
  align-items: center;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
}

.message {
  margin-bottom: 16px;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 16px;
  word-wrap: break-word;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.message.user .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.message.user .message-content .message-text {
  color: white;
  font-size: 14px;
  line-height: 1.6;
}

.message.assistant .message-content {
  background: white;
  border: 1px solid #e4e7ed;
}

/* Markdown预览样式 - 仅在AI助手页面内生效 */
.ai-assistant-page .message.assistant .message-content :deep(.md-preview) {
  background: white;
  border: none;
  padding: 0;
  margin: 0;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview .md-preview-html) {
  background: white;
  color: inherit;
  font-size: 14px;
  line-height: 1.6;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview pre) {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 12px;
  margin: 8px 0;
  overflow-x: auto;
  border: 1px solid #e4e7ed;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview code) {
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview p) {
  margin: 6px 0;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview h1, .md-preview h2, .md-preview h3, .md-preview h4, .md-preview h5, .md-preview h6) {
  margin: 12px 0 6px 0;
  font-weight: bold;
  color: #333;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview ul, .md-preview ol) {
  margin: 6px 0;
  padding-left: 20px;
}

.ai-assistant-page .message.assistant .message-content :deep(.md-preview blockquote) {
  border-left: 4px solid #667eea;
  margin: 8px 0;
  padding-left: 12px;
  color: #666;
  background: #f8f9fa;
  padding: 8px 12px;
  border-radius: 0 4px 4px 0;
}

.message-time {
  font-size: 12px;
  color: #999;
  margin-top: 6px;
}

.message.user .message-time {
  color: rgba(255, 255, 255, 0.8);
}

.input-area {
  padding: 20px;
  border-top: 1px solid #e4e7ed;
  background: white;
}

.input-actions {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}

.input-actions .el-button {
  flex: 1;
}

/* 滚动条样式 */
.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style> 