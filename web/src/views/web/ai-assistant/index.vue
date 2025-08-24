<template>
  <div class="ai-assistant-page">
    <div class="main-container">
      <!-- 左侧边栏 -->
      <div class="sidebar" :class="{ 'collapsed': sidebarCollapsed }">
        <div class="sidebar-header">
          <div class="avatar-section">
            <div class="avatar-name-container">
              <div class="avatar">AI</div>
              <span class="app-name">AI助手</span>
            </div>
            <el-button 
              type="text" 
              class="collapse-btn-header"
              @click="toggleSidebar"
            >
              <el-icon><ArrowLeft v-if="!sidebarCollapsed" /><ArrowRight v-else /></el-icon>
            </el-button>
          </div>
          <el-button 
            type="primary" 
            class="new-chat-btn"
            @click="createNewSession"
          >
            + 新对话
          </el-button>
        </div>
        
        <div class="conversation-section">
          <h3 class="conversation-title-fixed">历史对话</h3>
          <div class="conversation-list">
            <div class="conversation-items">
              <div 
                v-for="session in sessions" 
                :key="session.id"
                class="conversation-item"
                :class="{ 'active': currentSessionId === session.id }"
              >
                <div class="conversation-content" @click="selectSession(session.id)">
                  <el-icon><Document /></el-icon>
                  <span class="conversation-title">{{ session.title }}</span>
                </div>
                
                <!-- 会话操作按钮 -->
                <div class="conversation-actions">
                  <el-dropdown trigger="click" @command="handleSessionAction">
                    <el-button 
                      type="text" 
                      size="small"
                      class="action-btn"
                      @click.stop
                    >
                      <el-icon><MoreFilled /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :command="`rename:${session.id}`">
                          <el-icon><EditPen /></el-icon>
                          重命名
                        </el-dropdown-item>
                        <el-dropdown-item :command="`delete:${session.id}`" divided>
                          <el-icon><Delete /></el-icon>
                          <span style="color: #f56c6c;">删除</span>
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 主聊天区域 -->
      <div class="chat-container" :class="{ 'sidebar-collapsed': sidebarCollapsed }">

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
            </div>
          </div>
          
          <!-- 思考中状态 - 只显示加载状态，不显示对话框 -->
          <div v-if="thinking" class="thinking-status">
            <el-icon class="is-loading"><Loading /></el-icon>
            正在思考中...
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-area">
          <div class="chat-input-container">
            <!-- 上部分：输入区域 -->
            <div class="input-section">
              <el-input
                v-model="inputMessage"
                placeholder="请输入问题"
                type="textarea"
                :rows="1"
                @keydown.enter.prevent="sendMessage"
                :disabled="loading"
                class="chat-textarea"
              />
            </div>
            
            <!-- 下部分：工具栏 -->
            <div class="toolbar-section">
              <div class="toolbar-left">
                <el-select 
                  v-model="selectedModel" 
                  placeholder="选择模型" 
                  size="small"
                  class="model-selector"
                  :disabled="!isNewSession"
                >
                  <el-option
                    v-for="model in availableModels"
                    :key="model.name"
                    :label="model.display_name"
                    :value="model.name"
                  />
                </el-select>
              </div>
              <div class="toolbar-right">
                <el-button 
                  type="primary" 
                  @click="sendMessage"
                  :loading="loading"
                  :disabled="!inputMessage.trim()"
                  class="send-button"
                  size="large"
                >
                  <el-icon><ArrowUp /></el-icon>
                </el-button>
              </div>
            </div>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading, Document, ArrowLeft, ArrowRight, ArrowUp, Paperclip, User, Scissor, Microphone, MoreFilled, EditPen, Delete } from '@element-plus/icons-vue'
import { aiChatApi } from '@/api/ai-chat'
import { MdPreview } from 'md-editor-v3'

// 事件类型常量
const EVENT_MESSAGE = 1        // 正常消息
const EVENT_COMPLETE = 2       // 流式响应完成
const EVENT_TITLE_GENERATED = 3 // 标题生成完成

// 响应式数据
const loading = ref(false)
const streaming = ref(false)
const thinking = ref(false)
const inputMessage = ref('')
const messages = ref<any[]>([])
const sessions = ref<any[]>([])
const currentSessionId = ref<number | null>(null)
const showNewSessionDialog = ref(false)
const messagesContainer = ref<HTMLElement>()
const sidebarCollapsed = ref(false)
const selectedModel = ref('')
const isNewSession = ref(false)

// 新建会话表单
const newSessionForm = ref({
  title: '',
  model: ''
})

// 可用的AI模型
const availableModels = ref<any[]>([])

// 页面加载时获取会话列表和可用模型
onMounted(async () => {
  await loadAvailableModels()
  await loadSessions()
})

// 加载可用模型
const loadAvailableModels = async () => {
  try {
    const response = await aiChatApi.getAvailableModels()
    availableModels.value = response.data || []
  } catch (error) {
    console.error('加载可用模型失败:', error)
    ElMessage.error('加载可用模型失败')
  }
}

// 加载会话列表
const loadSessions = async () => {
  try {
    const response = await aiChatApi.getSessions()
    sessions.value = response.data.list || []
    
    // 如果有会话且当前没有选中的会话，选择第一个
    if (sessions.value.length > 0 && !currentSessionId.value) {
      currentSessionId.value = sessions.value[0].id
      // 使用会话的模型信息
      const currentSession = sessions.value[0]
      if (currentSession.model) {
        selectedModel.value = currentSession.model
      }
      await loadSessionMessages()
    }
    
    // 如果当前会话ID不在列表中，清空消息
    if (currentSessionId.value && !sessions.value.find(s => s.id === currentSessionId.value)) {
      messages.value = []
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
  thinking.value = true
  
  // 添加用户消息到UI
  messages.value.push({
    id: Date.now(),
    role: 'user',
    content,
    created_at: new Date().toISOString()
  })
  
  await nextTick()
  scrollToBottom()
  
  // 如果是新会话的第一次提问，先创建真正的会话
  if (isNewSession.value && currentSessionId.value) {
    try {
      const response = await aiChatApi.createSession({
        title: '新对话',
        model: selectedModel.value
      })
      
      // 更新会话ID和状态
      currentSessionId.value = response.data.id as number
      isNewSession.value = false
      
      // 重新加载会话列表
      await loadSessions()
      
    } catch (error) {
      console.error('创建会话失败:', error)
      ElMessage.error('创建会话失败')
      return
    }
  }
  
  // 不预先创建AI消息占位符，等收到第一个数据时再创建
  let aiMessageIndex = -1
  
  try {
    loading.value = false
    streaming.value = true
    
    // 使用流式API
    await aiChatApi.sendMessageStream({
      session_id: currentSessionId.value,
      content
    }, {
      onData: (data) => {
        // 收到数据时，根据event_id处理不同类型的消息
        if (data.content !== undefined || data.event_id) {
          thinking.value = false
          
          // 根据event_id处理不同类型的消息
          switch (data.event_id) {
            case EVENT_TITLE_GENERATED:
              // 收到标题生成完成标记，立即请求新标题
              if (currentSessionId.value) {
                updateSessionTitle(currentSessionId.value)
              }
              return
              
            case EVENT_COMPLETE:
              // 流式响应完成，不需要特殊处理
              return
              
            case EVENT_MESSAGE:
            default:
              // 正常消息内容，创建或更新AI消息
              if (data.content && data.content.trim()) {
                // 如果是第一个数据，创建AI消息
                if (aiMessageIndex === -1) {
                  aiMessageIndex = messages.value.length
                  messages.value.push({
                    id: data.message_id || Date.now(),
                    role: 'assistant',
                    content: data.content,
                    created_at: new Date().toISOString()
                  })
                } else {
                  // 后续数据，追加内容
                  messages.value[aiMessageIndex].content += data.content
                }
                
                nextTick(() => {
                  scrollToBottom()
                })
              }
              break
          }
        }
      },
      onComplete: async (data) => {
        // 流式响应完成，标题更新已通过event_id处理
        // 这里主要用于处理message_id等元数据
        if (data.message_id && aiMessageIndex !== -1) {
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
    // 如果已经创建了AI消息，则移除它
    if (aiMessageIndex !== -1) {
      messages.value.splice(aiMessageIndex, 1)
    }
  } finally {
    loading.value = false
    streaming.value = false
    thinking.value = false
  }
}

// 新建会话
const createNewSession = async () => {
  // 重新加载可用模型
  await loadAvailableModels()
  
  // 直接创建临时会话，不显示对话框
  isNewSession.value = true
  currentSessionId.value = Date.now() // 临时ID
  messages.value = []
  
  // 清空输入框
  inputMessage.value = ''
  
  ElMessage.success('新对话已创建，请选择模型并开始提问')
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

// 选择会话
const selectSession = (sessionId: number) => {
  currentSessionId.value = sessionId
  // 切换到已存在的会话，重置新会话状态
  isNewSession.value = false
  // 使用选中会话的模型信息
  const selectedSession = sessions.value.find(s => s.id === sessionId)
  if (selectedSession && selectedSession.model) {
    selectedModel.value = selectedSession.model
  }
  loadSessionMessages()
}

// 更新会话标题
const updateSessionTitle = async (sessionId: number) => {
  try {
    const response = await aiChatApi.getSessionDetail(sessionId)
    
    if (response.data.title && response.data.title !== '新对话') {
      // 更新本地会话列表中的标题
      const sessionIndex = sessions.value.findIndex(s => s.id === sessionId)
      if (sessionIndex !== -1) {
        sessions.value[sessionIndex].title = response.data.title
      }
      
      ElMessage.success('对话标题已自动生成')
    }
  } catch (error) {
    console.error('获取会话详情失败:', error)
    // 不显示错误提示，避免影响用户体验
  }
}

// 处理会话操作
const handleSessionAction = async (command: string) => {
  const [action, sessionId] = command.split(':')
  const session = sessions.value.find(s => s.id === parseInt(sessionId))

  if (!session) {
    ElMessage.error('会话不存在')
    return
  }

  if (action === 'rename') {
    try {
      const { value: newTitle } = await ElMessageBox.prompt('请输入新的会话标题', '重命名', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputValue: session.title,
        inputValidator: (value) => {
          if (!value || value.trim() === '') {
            return '标题不能为空'
          }
          return true
        }
      })
      
      if (newTitle && newTitle.trim()) {
        await aiChatApi.updateSession({
          session_id: session.id,
          title: newTitle.trim()
        })
        ElMessage.success('会话标题已更新')
        session.title = newTitle.trim()
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('更新会话标题失败:', error)
        ElMessage.error('更新会话标题失败')
      }
    }
  } else if (action === 'delete') {
    try {
      await ElMessageBox.confirm(`确定要删除会话 "${session.title}" 吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
      
      await aiChatApi.deleteSession({
        session_id: session.id
      })
      ElMessage.success('会话已删除')
      sessions.value = sessions.value.filter(s => s.id !== session.id)
      if (currentSessionId.value === session.id) {
        currentSessionId.value = null
        messages.value = []
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除会话失败:', error)
        ElMessage.error('删除会话失败')
      }
    }
  }
}

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
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
  max-width: 60vw;
  margin: 0 auto;
  background: transparent;
  padding: 0px;
}

.main-container {
  padding: 0px;
  width: 100%;
  height: 100%;
  display: flex;
}

.chat-container {
  background: white;
  border-radius: 0;
  overflow: hidden;
  height: 73vh;
  display: flex;
  flex-direction: column;
  flex: 1;
  transition: margin-left 0.3s ease;
}

.chat-container.sidebar-collapsed {
  margin-left: 0;
}



.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: white;
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
  border-radius: 12px;
  word-wrap: break-word;
}

.message.user .message-content {
  background: #e8e8e8;
  color: #333;
}

.message.user .message-content .message-text {
  color: #333;
  font-size: 14px;
  line-height: 1.6;
}

.message.assistant .message-content {
  background: white;
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

/* 针对github主题的样式调整 */
.ai-assistant-page .message.assistant .message-content :deep(.github-theme *) {
  margin-top: 0px;
  margin-bottom: 0px;
}

/* 思考中状态样式 */
.thinking-status {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #666;
  font-size: 14px;
  gap: 8px;
}

.thinking-status .el-icon {
  font-size: 16px;
  color: #409eff;
}

.input-area {
  padding: 10px;
  background: white;
}

/* 对话框样式 */
.chat-input-container {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  padding: 8px;
}

.input-section {
  padding: 0px;
  background: white;
}

.chat-textarea {
  border: none;
  resize: none;
  background: white;
}

.chat-textarea :deep(.el-textarea__inner) {
  border: none;
  box-shadow: none;
  padding: 0;
  font-size: 14px;
  line-height: 1.6;
  min-height: 80px;
  background: white;
  color: #333;
}

.chat-textarea :deep(.el-textarea__inner::-webkit-resizer) {
  display: none;
}

.chat-textarea :deep(.el-textarea__inner:focus) {
  box-shadow: none;
  background: white;
}

.chat-textarea :deep(.el-textarea__inner::placeholder) {
  color: #999;
}

.toolbar-section {
  display: flex;
  justify-content: space-between; /* Changed to space-between for left and right */
  align-items: center;
  padding-top: 8px;
  background: white;
  min-height: 32px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  background: transparent;
  border: none;
  color: #666;
  padding: 6px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 13px;
}

.toolbar-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.toolbar-btn .el-icon {
  margin-right: 4px;
  font-size: 14px;
}

.send-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  border: none;
  box-shadow: none;
  transition: all 0.2s ease;
}

.send-button:hover {
  background: #e0e0e0;
  transform: none;
}

.send-button:disabled {
  background: #f0f0f0;
  opacity: 0.5;
}

.send-button:not(:disabled) {
  background: #409eff;
}

.send-button:not(:disabled):hover {
  background: #337ecc;
}

.send-button:not(:disabled) .el-icon {
  color: white;
}

.send-button .el-icon {
  font-size: 16px;
  color: #666;
}

.model-selector {
  width: 120px;
  margin-right: 10px;
}

.model-selector :deep(.el-input__wrapper) {
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
}

.model-selector :deep(.el-input__wrapper:hover) {
  border-color: #c0c4cc;
}

.model-selector :deep(.el-input__wrapper.is-focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}


/* 侧边栏样式 */
.sidebar {
  width: 280px;
  height: 73vh;
  background: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 60px;
}

.sidebar-header {
  padding: 10px 20px 10px 20px;
  position: relative;
  border-bottom: 1px solid #e4e7ed;
}

.avatar-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.avatar-name-container {
  display: flex;
  align-items: center;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #667eea;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 12px;
}

.app-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.new-chat-btn {
  width: 100%;
  height: 40px;
  margin-bottom: 0px;
}

.collapse-btn-header {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  padding: 0;
  background: #f0f0f0;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.collapse-btn-header:hover {
  background: #e0e0e0;
}

.sidebar.collapsed .collapse-btn-header {
  position: static;
  top: auto;
  right: auto;
  margin: 20px auto 0;
  display: block;
}

.conversation-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  min-height: 0;
}

.conversation-title-fixed {
  margin: 0;
  padding: 12px 0px 12px 20px;
  font-size: 14px;
  color: #666;
  font-weight: 500;
  background: white;
  flex-shrink: 0;
}

.conversation-list {
  flex: 1;
  padding: 0 20px 20px 20px;
  overflow-y: auto;
  background: white;
  /* 确保滚动条不占用内容空间 */
  scrollbar-width: none;
  -ms-overflow-style: none;
  min-height: 0;
}

.conversation-items {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 0px 12px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  gap: 12px;
  min-height: 30px;
}

.conversation-item:hover {
  background: #e9ecef;
}

.conversation-item.active {
  background: #e3f2fd;
  color: #1976d2;
}

.conversation-item .el-icon {
  font-size: 16px;
  color: #666;
  flex-shrink: 0;
}

.conversation-content {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.conversation-title {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.conversation-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.conversation-item:hover .conversation-actions {
  opacity: 1;
}

.action-btn {
  background: transparent;
  border: none;
  color: #666;
  padding: 6px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 13px;
}

.action-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.action-btn .el-icon {
  margin-right: 4px;
  font-size: 14px;
}


/* 侧边栏收缩时的样式 */
.sidebar.collapsed .new-chat-btn,
.sidebar.collapsed .conversation-section {
  display: none;
}

.sidebar.collapsed .avatar-section {
  justify-content: center;
}

.sidebar.collapsed .avatar-name-container {
  display: none;
}

.sidebar.collapsed .collapse-btn-header {
  margin: 0;
  width: 32px;
  height: 32px;
  flex-shrink: 0;
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

/* 会话列表滚动条样式 - 完全隐藏滚动条 */
.conversation-list::-webkit-scrollbar {
  display: none;
}
</style> 