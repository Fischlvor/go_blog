<template>
  <div class="ai-assistant">
    <!-- AI助手图标 -->
    <div 
      class="ai-assistant-icon" 
      :class="{ 'active': isOpen }"
      @click="toggleChat"
      title="AI助手"
    >
      <el-icon size="24">
        <ChatDotRound v-if="!isOpen" />
        <Close v-else />
      </el-icon>
    </div>

    <!-- 聊天对话框 -->
    <el-dialog
      v-model="isOpen"
      title="AI助手"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      class="ai-chat-dialog"
      destroy-on-close
    >
      <div class="chat-container">
        <!-- 会话选择 -->
        <div class="session-selector" v-if="sessions.length > 0">
          <el-select 
            v-model="currentSessionId" 
            placeholder="选择会话"
            @change="loadMessages"
            style="width: 100%; margin-bottom: 10px;"
          >
            <el-option
              v-for="session in sessions"
              :key="session.id"
              :label="session.title"
              :value="session.id"
            />
          </el-select>
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

        <!-- 输入框 -->
        <div class="input-container">
          <el-input
            v-model="inputMessage"
            placeholder="输入您的问题..."
            type="textarea"
            :rows="3"
            @keydown.enter.prevent="sendMessage"
            :disabled="loading"
          />
          <el-button 
            type="primary" 
            @click="sendMessage"
            :loading="loading"
            :disabled="!inputMessage.trim()"
            style="margin-top: 10px; width: 100%;"
          >
            发送
          </el-button>
        </div>

        <!-- 新建会话按钮 -->
        <el-button 
          type="success" 
          @click="createNewSession"
          style="margin-top: 10px; width: 100%;"
        >
          新建会话
        </el-button>
      </div>
    </el-dialog>

    <!-- 新建会话对话框 -->
    <el-dialog
      v-model="showNewSessionDialog"
      title="新建会话"
      width="300px"
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
import { ChatDotRound, Close, Loading } from '@element-plus/icons-vue'
import { aiChatApi } from '@/api/ai-chat'
import { MdPreview } from 'md-editor-v3'

// 响应式数据
const isOpen = ref(false)
const loading = ref(false)
const streaming = ref(false)
const inputMessage = ref('')
const messages = ref<any[]>([])
const sessions = ref<any[]>([])
const currentSessionId = ref<number | null>(null)
const availableModels = ref<any[]>([])
const showNewSessionDialog = ref(false)
const newSessionForm = ref({
  title: '',
  model: ''
})
const messagesContainer = ref<HTMLElement>()

// 组件挂载时验证
onMounted(() => {
  console.log('AI助手组件已加载')
})

// 切换聊天窗口
const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    loadSessions()
    loadAvailableModels()
  }
}

// 加载可用模型
const loadAvailableModels = async () => {
  try {
    const response = await aiChatApi.getAvailableModels()
    availableModels.value = response.data
  } catch (error) {
    console.error('加载模型失败:', error)
  }
}

// 加载会话列表
const loadSessions = async () => {
  try {
    const response = await aiChatApi.getSessions()
    sessions.value = response.data.list || []
    if (sessions.value.length > 0 && !currentSessionId.value) {
      currentSessionId.value = sessions.value[0].id
      await loadMessages()
    } else if (sessions.value.length === 0) {
      // 如果没有会话，自动创建第一个会话
      await createFirstSession()
    }
  } catch (error) {
    console.error('加载会话失败:', error)
  }
}

// 加载消息
const loadMessages = async () => {
  if (!currentSessionId.value) return
  
  try {
    const response = await aiChatApi.getMessages(currentSessionId.value)
    messages.value = response.data.list || []
    await nextTick()
    scrollToBottom()
  } catch (error) {
    console.error('加载消息失败:', error)
  }
}

// 发送消息 - 支持流式响应
const sendMessage = async () => {
  if (!inputMessage.value.trim() || !currentSessionId.value) return
  
  const content = inputMessage.value
  inputMessage.value = ''
  loading.value = true
  streaming.value = false
  
  // 添加用户消息到界面
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
    streaming.value = true // 开始流式响应
    
    // 使用通用的流式API
    await aiChatApi.sendMessageStream({
      session_id: currentSessionId.value,
      content
    }, {
      onData: (data) => {
        // 更新AI消息内容
        if (data.content !== undefined) {
          messages.value[aiMessageIndex].content += data.content
          nextTick(() => {
            scrollToBottom()
          })
        }
      },
      onComplete: (data) => {
        // 流式响应完成，更新消息ID
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
    // 移除空的AI消息
    messages.value.splice(aiMessageIndex, 1)
  } finally {
    loading.value = false
    streaming.value = false // 结束流式响应
  }
}



// 创建新会话
const createNewSession = () => {
  showNewSessionDialog.value = true
  newSessionForm.value = {
    title: '',
    model: availableModels.value[0]?.name || ''
  }
}

// 创建第一个会话
const createFirstSession = async () => {
  if (availableModels.value.length === 0) {
    ElMessage.warning('暂无可用的AI模型')
    return
  }
  
  try {
    const firstModel = availableModels.value[0]
    const response = await aiChatApi.createSession({
      title: '新对话',
      model: firstModel.name
    })
    sessions.value.push(response.data)
    currentSessionId.value = response.data.id
    messages.value = []
    ElMessage.success('会话创建成功')
  } catch (error) {
    ElMessage.error('创建会话失败')
    console.error('创建会话失败:', error)
  }
}

// 确认创建会话
const confirmCreateSession = async () => {
  if (!newSessionForm.value.title || !newSessionForm.value.model) {
    ElMessage.warning('请填写完整信息')
    return
  }
  
  try {
    const response = await aiChatApi.createSession(newSessionForm.value)
    sessions.value.push(response.data)
    currentSessionId.value = response.data.id
    messages.value = []
    showNewSessionDialog.value = false
    ElMessage.success('会话创建成功')
  } catch (error) {
    ElMessage.error('创建会话失败')
    console.error('创建会话失败:', error)
  }
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
.ai-assistant {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1000;
}

.ai-assistant-icon {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  color: white;
}

.ai-assistant-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

.ai-assistant-icon.active {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.ai-chat-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.chat-container {
  height: 500px;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 12px;
  background: #fafafa;
}

.message {
  margin-bottom: 12px;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 80%;
  padding: 8px 12px;
  border-radius: 12px;
  word-wrap: break-word;
}

.message.user .message-content {
  background: #409eff;
  color: white;
}

.message.user .message-content .message-text {
  color: white;
  font-size: 14px;
  line-height: 1.5;
}

.message.assistant .message-content {
  background: white;
  border: 1px solid #e4e7ed;
}

/* Markdown预览样式 */
.message.assistant .message-content :deep(.md-preview) {
  background: transparent;
  border: none;
  padding: 0;
  margin: 0;
}

.message.assistant .message-content :deep(.md-preview .md-preview-html) {
  background: transparent;
  color: inherit;
  font-size: 14px;
  line-height: 1.5;
}

.message.assistant .message-content :deep(.md-preview pre) {
  background: #f5f5f5;
  border-radius: 4px;
  padding: 8px;
  margin: 4px 0;
  overflow-x: auto;
}

.message.assistant .message-content :deep(.md-preview code) {
  background: #f0f0f0;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.message.assistant .message-content :deep(.md-preview p) {
  margin: 4px 0;
}

.message.assistant .message-content :deep(.md-preview h1, .md-preview h2, .md-preview h3, .md-preview h4, .md-preview h5, .md-preview h6) {
  margin: 8px 0 4px 0;
  font-weight: bold;
}

.message.assistant .message-content :deep(.md-preview ul, .md-preview ol) {
  margin: 4px 0;
  padding-left: 20px;
}

.message.assistant .message-content :deep(.md-preview blockquote) {
  border-left: 3px solid #ddd;
  margin: 4px 0;
  padding-left: 10px;
  color: #666;
}

.message-time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.message.user .message-time {
  color: rgba(255, 255, 255, 0.8);
}

.input-container {
  margin-bottom: 10px;
}

.session-selector {
  margin-bottom: 10px;
}
</style> 