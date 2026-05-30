<template>
  <Layout>
    <div class="im-page">
      <div class="chat-layout">
        <aside class="sidebar">
          <div class="sidebar-header">
            <h3>会话列表</h3>
            <span class="ws-status" :class="{ online: wsConnected }">
              {{ wsConnected ? '在线' : '离线' }}
            </span>
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
            <div class="chat-avatar">{{ currentSession.name?.charAt(0) || 'U' }}</div>
            <span class="chat-title">{{ currentSession.name }}</span>
          </div>

          <div v-if="!currentSession" class="chat-placeholder">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <p>选择一个会话开始聊天</p>
          </div>

          <template v-else>
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
              <div v-if="messages.length === 0" class="empty-chat">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                </svg>
                <p>开始聊天吧</p>
              </div>
            </div>

            <div class="message-input">
              <el-input
                v-model="inputMessage"
                type="textarea"
                :rows="3"
                placeholder="输入消息..."
                @keyup.enter.ctrl="sendMessage"
              />
              <button class="send-btn" @click="sendMessage" :disabled="!wsConnected || !inputMessage.trim()">
                <svg viewBox="0 0 24 24" fill="currentColor"><path d="M2 21l21-9L2 3v7l15 2-15 2v7z"/></svg>
                发送
              </button>
            </div>
          </template>
        </section>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import Layout from '../components/Layout.vue'

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

onMounted(() => {
  if (sessions.value.length > 0) {
    selectSession(sessions.value[0])
  }
  connectWebSocket()
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})
</script>

<style scoped>
.im-page {
  height: calc(100vh - 160px);
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.chat-layout {
  height: 100%;
  display: flex;
}

.sidebar {
  width: 280px;
  background: #fafafa;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.ws-status {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 10px;
  background: #f0f0f0;
  color: #999;
}

.ws-status.online {
  background: #E8F5E9;
  color: #67C23A;
}

.session-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.session-item:hover {
  background: #f0f0f0;
}

.session-item.active {
  background: #FFF5F7;
}

.session-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #FF7490, #FF5A78);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.session-info {
  flex: 1;
  overflow: hidden;
}

.session-name {
  display: block;
  font-size: 14px;
  color: #333;
  font-weight: 500;
  margin-bottom: 4px;
}

.session-preview {
  display: block;
  color: #999;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  gap: 12px;
}

.chat-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #FF7490, #FF5A78);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 600;
  font-size: 14px;
}

.chat-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.chat-placeholder {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #ccc;
  gap: 12px;
}

.chat-placeholder svg {
  width: 64px;
  height: 64px;
}

.chat-placeholder p {
  font-size: 14px;
  color: #999;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
}

.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 12px;
}

.empty-chat svg {
  width: 48px;
  height: 48px;
  color: #ccc;
}

.empty-chat p {
  font-size: 14px;
  color: #999;
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
  background: linear-gradient(135deg, #E8E8E8, #D0D0D0);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  font-size: 14px;
  flex-shrink: 0;
}

.message-item.self .message-avatar {
  background: linear-gradient(135deg, #FF7490, #FF5A78);
  color: #fff;
}

.message-content {
  max-width: 60%;
}

.message-sender {
  color: #999;
  font-size: 12px;
  margin-bottom: 4px;
}

.message-item.self .message-sender {
  text-align: right;
}

.message-text {
  background: #fff;
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid #eee;
  color: #333;
  font-size: 14px;
  word-break: break-word;
  line-height: 1.6;
}

.message-item.self .message-text {
  background: linear-gradient(135deg, #FF7490, #FF5A78);
  color: #fff;
  border: none;
}

.message-time {
  color: #ccc;
  font-size: 11px;
  margin-top: 4px;
}

.message-item.self .message-time {
  text-align: right;
}

.message-input {
  padding: 16px 20px;
  border-top: 1px solid #eee;
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.message-input .el-textarea {
  flex: 1;
}

.send-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #FF7490, #FF5A78);
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255,116,144,0.3);
}

.send-btn:disabled {
  background: #e0e0e0;
  cursor: not-allowed;
}

.send-btn svg {
  width: 16px;
  height: 16px;
}
</style>