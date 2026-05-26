<template>
  <div class="im-container">
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-icon">💬</span>
          <span class="logo-text">即时通讯</span>
        </div>
        <div class="nav-links">
          <router-link to="/seckill">秒杀商城</router-link>
          <el-button type="danger" plain size="small" @click="handleLogout">退出</el-button>
        </div>
      </div>
    </header>
    
    <main class="main-content">
      <div class="chat-layout">
        <aside class="sidebar">
          <div class="sidebar-header">
            <h3>会话列表</h3>
          </div>
          <div class="session-list">
            <div 
              v-for="session in sessions" 
              :key="session.id"
              :class="['session-item', { active: currentSession?.id === session.id }]"
              @click="selectSession(session)"
            >
              <div class="session-avatar">{{ session.name?.charAt(0) || 'U' }}</div>
              <div class="session-info">
                <span class="session-name">{{ session.name }}</span>
                <span class="session-preview">{{ session.lastMessage || '暂无消息' }}</span>
              </div>
            </div>
            <el-empty v-if="sessions.length === 0" description="暂无会话" :image-size="80" />
          </div>
        </aside>
        
        <section class="chat-area">
          <div v-if="currentSession" class="chat-header">
            <span class="chat-title">{{ currentSession.name }}</span>
            <span class="chat-status" :class="{ online: wsConnected }">
              {{ wsConnected ? '在线' : '离线' }}
            </span>
          </div>
          
          <div ref="messageListRef" class="message-list">
            <div 
              v-for="msg in messages" 
              :key="msg.id"
              :class="['message-item', { self: msg.self }]"
            >
              <div class="message-avatar">{{ msg.sender?.charAt(0) || 'U' }}</div>
              <div class="message-content">
                <div class="message-sender">{{ msg.sender }}</div>
                <div class="message-text">{{ msg.content }}</div>
                <div class="message-time">{{ formatMsgTime(msg.time) }}</div>
              </div>
            </div>
            <el-empty v-if="messages.length === 0" description="开始聊天吧" :image-size="100" />
          </div>
          
          <div class="message-input">
            <el-input
              v-model="inputMessage"
              type="textarea"
              :rows="3"
              placeholder="输入消息..."
              @keyup.enter.ctrl="sendMessage"
            />
            <el-button type="primary" @click="sendMessage">
              发送
            </el-button>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const wsConnected = ref(false)
const ws = ref(null)
const sessions = ref([
  { id: 1, name: '官方客服', lastMessage: '欢迎来到秒杀商城' },
  { id: 2, name: '系统通知', lastMessage: '秒杀活动即将开始' }
])
const currentSession = ref(null)
const messages = ref([])
const inputMessage = ref('')
const messageListRef = ref(null)

function connectWebSocket() {
  const token = localStorage.getItem('accessToken')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  try {
    ws.value = new WebSocket(`ws://localhost:9999/ws?token=${token}`)
    
    ws.value.onopen = () => {
      wsConnected.value = true
      ElMessage.success('IM连接成功')
    }
    
    ws.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handleWsMessage(data)
      } catch (e) {
        console.error('解析消息失败', e)
      }
    }
    
    ws.value.onclose = () => {
      wsConnected.value = false
      ElMessage.warning('IM连接已断开')
    }
    
    ws.value.onerror = (error) => {
      wsConnected.value = false
      console.error('WebSocket错误', error)
    }
  } catch (error) {
    console.error('连接WebSocket失败', error)
  }
}

function handleWsMessage(data) {
  if (data.type === 'message') {
    messages.value.push({
      id: Date.now(),
      sender: data.sender || '对方',
      content: data.content,
      time: new Date(),
      self: false
    })
    scrollToBottom()
  }
}

function sendMessage() {
  if (!inputMessage.value.trim()) return
  if (!wsConnected.value) {
    ElMessage.warning('未连接到服务器')
    return
  }
  
  const msg = {
    type: 'message',
    sessionId: currentSession.value?.id,
    content: inputMessage.value,
    sender: userStore.username
  }
  
  try {
    ws.value.send(JSON.stringify(msg))
    messages.value.push({
      id: Date.now(),
      sender: userStore.username,
      content: inputMessage.value,
      time: new Date(),
      self: true
    })
    inputMessage.value = ''
    scrollToBottom()
  } catch (error) {
    ElMessage.error('发送失败')
  }
}

function selectSession(session) {
  currentSession.value = session
  messages.value = []
}

function scrollToBottom() {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

function formatMsgTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

async function handleLogout() {
  if (ws.value) {
    ws.value.close()
  }
  await userStore.logout()
  router.push('/login')
}

onMounted(() => {
  selectSession(sessions.value[0])
  connectWebSocket()
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})
</script>

<style scoped>
.im-container {
  min-height: 100vh;
  background: #1a1a2e;
}

.header {
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0 40px;
  height: 60px;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-links a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.3s;
}

.nav-links a:hover {
  color: #fff;
}

.main-content {
  height: calc(100vh - 60px);
  padding: 20px;
}

.chat-layout {
  height: 100%;
  display: flex;
  gap: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.sidebar {
  width: 280px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-header h3 {
  color: #fff;
  font-size: 16px;
}

.session-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s;
}

.session-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.session-item.active {
  background: rgba(102, 126, 234, 0.3);
}

.session-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 600;
}

.session-info {
  flex: 1;
  overflow: hidden;
}

.session-name {
  display: block;
  color: #fff;
  font-size: 14px;
  margin-bottom: 4px;
}

.session-preview {
  display: block;
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-area {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chat-title {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
}

.chat-status {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
}

.chat-status.online {
  background: rgba(103, 194, 106, 0.2);
  color: #67c23a;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.message-item {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.message-item.self {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  flex-shrink: 0;
}

.message-item.self .message-avatar {
  background: linear-gradient(135deg, #f5576c, #f093fb);
}

.message-content {
  max-width: 60%;
}

.message-sender {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  margin-bottom: 4px;
}

.message-text {
  background: rgba(255, 255, 255, 0.1);
  padding: 12px 16px;
  border-radius: 12px;
  color: #fff;
  font-size: 14px;
  word-break: break-word;
}

.message-item.self .message-text {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.message-time {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
  margin-top: 4px;
}

.message-input {
  padding: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  gap: 10px;
}

.message-input .el-textarea {
  flex: 1;
}

.message-input .el-button {
  align-self: flex-end;
}
</style>
