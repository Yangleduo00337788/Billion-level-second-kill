<template>
  <div class="max-w-7xl mx-auto px-4 h-[calc(100vh-4rem)]">
    <div class="flex bg-white rounded-2xl border border-gray-100 h-full overflow-hidden">
      <div class="w-60 lg:w-80 border-r border-gray-100 flex flex-col flex-shrink-0 hidden md:flex">
        <div class="p-4 border-b border-gray-100 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-dark">AI 助手</h2>
          <n-button size="tiny" quaternary circle @click="newSession">
            <template #icon>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
            </template>
          </n-button>
        </div>
        <div class="flex-1 overflow-y-auto">
          <div
            v-for="(session, i) in sessions"
            :key="session.id"
            class="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors hover:bg-gray-50 group"
            :class="{ 'bg-gray-50': activeSession === i }"
            @click="switchSession(i)"
          >
            <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
              <span class="text-xs font-medium text-primary">AI</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-dark truncate">{{ session.title }}</p>
              <p class="text-xs text-gray-400 truncate mt-0.5">{{ session.messages[session.messages.length - 1]?.content?.slice(0, 30) || '新对话' }}</p>
            </div>
            <n-button size="tiny" quaternary circle class="opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" @click.stop="deleteSession(i)">
              <template #icon>
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </template>
            </n-button>
          </div>
          <div v-if="sessions.length === 0" class="text-center py-12 text-gray-400 text-sm">暂无对话</div>
        </div>
      </div>

      <div class="flex-1 flex flex-col">
        <div class="md:hidden px-4 py-3 border-b border-gray-100 flex items-center gap-3">
          <n-button size="tiny" quaternary circle @click="showSidebar = !showSidebar">
            <template #icon>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </template>
          </n-button>
          <h2 class="text-sm font-semibold text-dark">AI 助手</h2>
        </div>

        <div v-if="currentSession" class="flex-1 overflow-y-auto p-4 lg:p-6 space-y-4">
          <div v-for="(msg, i) in currentSession.messages" :key="i" class="flex gap-3" :class="msg.role === 'user' ? 'flex-row-reverse' : ''">
            <div v-if="msg.role === 'assistant'" class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
              <span class="text-xs font-medium text-primary">AI</span>
            </div>
            <div
              :class="msg.role === 'user' ? 'bg-primary text-white' : 'bg-gray-100 text-gray-800'"
              class="max-w-[85%] lg:max-w-[70%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed whitespace-pre-wrap"
            >
              <template v-if="msg.role === 'assistant' && i === currentSession.messages.length - 1 && isStreaming">
                {{ msg.content }}<span class="inline-block w-1.5 h-4 bg-primary animate-pulse ml-0.5 align-text-bottom"></span>
              </template>
              <template v-else>{{ msg.content }}</template>
            </div>
          </div>
          <div v-if="currentSession.messages.length === 0" class="flex items-center justify-center h-full text-gray-400">
            <div class="text-center">
              <div class="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
              </div>
              <p class="text-sm">开始和 AI 对话</p>
              <p class="text-xs mt-1">在下方输入你的问题</p>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-100 p-4">
          <div class="flex items-center gap-3">
            <n-input
              v-model:value="inputMessage"
              type="textarea"
              placeholder="输入消息..."
              :rows="1"
              :autosize="{ minRows: 1, maxRows: 4 }"
              class="flex-1"
              :disabled="isStreaming"
              @keydown.enter.prevent="sendMessage"
            />
            <n-button type="primary" circle @click="sendMessage" :disabled="!inputMessage.trim() || isStreaming">
              <template #icon>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19V5m0 0l-7 7m7-7l7 7" />
                </svg>
              </template>
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { useMessage } from 'naive-ui'

const userStore = useUserStore()
const message = useMessage()

interface ChatMsg {
  role: 'user' | 'assistant'
  content: string
}

interface Session {
  id: number
  title: string
  messages: ChatMsg[]
}

const sessions = ref<Session[]>([])
const activeSession = ref(0)
const inputMessage = ref('')
const isStreaming = ref(false)
const showSidebar = ref(false)

const currentSession = computed(() => sessions.value[activeSession.value])

function newSession() {
  sessions.value.push({ id: Date.now(), title: '新对话', messages: [] })
  activeSession.value = sessions.value.length - 1
}

function switchSession(i: number) {
  activeSession.value = i
  showSidebar.value = false
}

function deleteSession(i: number) {
  sessions.value.splice(i, 1)
  if (sessions.value.length === 0) {
    newSession()
  } else if (activeSession.value >= sessions.value.length) {
    activeSession.value = sessions.value.length - 1
  }
}

async function sendMessage() {
  if (!inputMessage.value.trim() || isStreaming.value) return
  if (!userStore.isAuthenticated) {
    message.warning('请先登录使用 AI 助手')
    return
  }

  const content = inputMessage.value.trim()
  inputMessage.value = ''

  if (!currentSession.value) newSession()
  const session = currentSession.value!
  if (session.messages.length === 0) {
    session.title = content.slice(0, 30) + (content.length > 30 ? '...' : '')
  }

  session.messages.push({ role: 'user', content })
  session.messages.push({ role: 'assistant', content: '' })
  isStreaming.value = true

  try {
    const messages = session.messages.slice(0, -1).map(m => ({ role: m.role, content: m.content }))
    const token = localStorage.getItem('token')

    const response = await fetch('/api/v1/ai/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ messages })
    })

    if (!response.ok) throw new Error('AI 服务请求失败')

    const reader = response.body?.getReader()
    if (!reader) throw new Error('No reader')

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6)
          if (data === '[DONE]') break
          try {
            const parsed = JSON.parse(data)
            if (parsed.content) {
              const lastMsg = session.messages[session.messages.length - 1]
              if (lastMsg?.role === 'assistant') {
                lastMsg.content += parsed.content
              }
            }
          } catch {}
        }
      }
    }
  } catch (e: any) {
    message.error(e.message || 'AI 响应失败')
    const lastMsg = session.messages[session.messages.length - 1]
    if (lastMsg?.role === 'assistant' && !lastMsg.content) {
      lastMsg.content = '抱歉，我遇到了问题，请稍后重试。'
    }
  } finally {
    isStreaming.value = false
  }
}

newSession()
</script>
