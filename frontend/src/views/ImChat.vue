<template>
  <div class="max-w-7xl mx-auto px-4 h-[calc(100vh-4rem)]">
    <div class="flex bg-white rounded-2xl border border-gray-100 h-full overflow-hidden">
      <div class="w-60 lg:w-80 border-r border-gray-100 flex flex-col flex-shrink-0 hidden md:flex">
        <div class="p-4 border-b border-gray-100">
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索联系人..."
              class="w-full pl-9 pr-3 py-2 text-sm bg-gray-50 border border-gray-100 rounded-lg outline-none text-gray-600 placeholder-gray-400 focus:border-primary/30 focus:bg-white transition-colors"
            />
          </div>
        </div>
        <div class="flex-1 overflow-y-auto">
          <div
            v-for="contact in filteredContacts"
            :key="contact.id"
            class="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors hover:bg-gray-50"
            :class="{ 'bg-primary/5': store.activeContactId === contact.id }"
            @click="selectContact(contact.id)"
          >
            <div class="relative flex-shrink-0">
              <n-avatar
                :src="contact.avatar"
                round
                size="medium"
              />
              <span
                v-if="isOnline(contact.id)"
                class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 border-2 border-white rounded-full"
              ></span>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-dark truncate">{{ contact.username }}</p>
                <span v-if="isOnline(contact.id)" class="text-xs text-green-500 ml-2">在线</span>
              </div>
              <p class="text-xs text-gray-400 truncate mt-0.5">
                {{ lastMessage(contact.id) }}
              </p>
            </div>
          </div>
          <div v-if="filteredContacts.length === 0" class="text-center py-12 text-gray-400 text-sm">
            {{ searchQuery ? '未找到联系人' : '暂无联系人' }}
          </div>
        </div>
      </div>

      <div class="flex-1 flex flex-col">
        <template v-if="activeContact">
          <div class="px-4 lg:px-6 py-4 border-b border-gray-100 flex items-center gap-3">
            <n-button size="tiny" quaternary circle class="md:hidden" @click="showSidebar = !showSidebar">
              <template #icon>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </template>
            </n-button>
            <n-avatar
              :src="activeContact.avatar"
              round
              size="small"
            />
            <div>
              <p class="text-sm font-medium text-dark">{{ activeContact.username }}</p>
              <p class="text-xs" :class="isOnline(activeContact.id) ? 'text-green-500' : 'text-gray-400'">
                {{ isOnline(activeContact.id) ? '在线' : '离线' }}
              </p>
            </div>
          </div>

          <div ref="messageListRef" class="flex-1 overflow-y-auto p-4 lg:p-6 space-y-3">
            <div v-for="msg in store.activeMessages" :key="msg.id" class="flex" :class="msg.sender_id === currentUserId ? 'justify-end' : 'justify-start'">
              <div
                class="max-w-[85%] lg:max-w-[70%] rounded-2xl px-4 py-2.5 text-sm leading-relaxed whitespace-pre-wrap break-words"
                :class="msg.sender_id === currentUserId ? 'bg-primary text-white rounded-br-md' : 'bg-gray-100 text-gray-800 rounded-bl-md'"
              >
                {{ msg.content }}
              </div>
            </div>
            <div v-if="store.activeMessages.length === 0" class="flex items-center justify-center h-full text-gray-400 text-sm">暂无消息，发送第一条消息吧</div>
          </div>

          <div class="px-4 lg:px-6 py-4 border-t border-gray-100">
            <div class="flex items-end gap-3">
              <div class="flex-1 relative">
                <textarea
                  v-model="inputText"
                  placeholder="输入消息..."
                  class="w-full px-4 py-2.5 text-sm bg-gray-50 border border-gray-100 rounded-xl outline-none resize-none text-gray-600 placeholder-gray-400 focus:border-primary/30 focus:bg-white transition-colors"
                  rows="2"
                  @keydown.enter.exact.prevent="send"
                ></textarea>
              </div>
              <n-button type="primary" size="medium" round :disabled="!inputText.trim()" @click="send">
                <template #icon>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19V5m0 0l-7 7m7-7l7 7" />
                  </svg>
                </template>
              </n-button>
            </div>
          </div>
        </template>

        <template v-else>
          <div class="flex-1 flex items-center justify-center text-gray-400">
            <div class="text-center">
              <div class="w-20 h-20 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-4">
                <svg class="w-10 h-10 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
              <p class="text-sm">选择联系人开始聊天</p>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { NButton, NAvatar, useMessage } from 'naive-ui'
import { useChatStore } from '@/stores/chat'
import { useUserStore } from '@/stores/user'
import type { User } from '@/types/api'

const store = useChatStore()
const userStore = useUserStore()
const message = useMessage()

const searchQuery = ref('')
const inputText = ref('')
const messageListRef = ref<HTMLElement | null>(null)
const showSidebar = ref(false)

const currentUserId = computed(() => userStore.user?.id ?? 0)

const filteredContacts = computed(() => {
  const list = store.contacts
  if (!searchQuery.value.trim()) return list
  const q = searchQuery.value.toLowerCase()
  return list.filter(c => c.username.toLowerCase().includes(q))
})

const activeContact = computed(() => {
  const id = store.activeContactId
  if (!id) return null
  return store.contacts.find(c => c.id === id) || null
})

function isOnline(userId: number) {
  return store.isOnline(userId)
}

function selectContact(userId: number) {
  store.setActiveContact(userId)
  showSidebar.value = false
}

function lastMessage(contactId: number) {
  const msgs = store.messages[contactId]
  if (!msgs || msgs.length === 0) return '暂无消息'
  const last = msgs[msgs.length - 1]
  return last.content.slice(0, 30)
}

function send() {
  const text = inputText.value.trim()
  if (!text || !store.activeContactId) return
  try {
    store.sendMessage(text, store.activeContactId)
    inputText.value = ''
  } catch (e: any) {
    message.error('发送失败')
  }
}

onMounted(() => {
  if (userStore.isAuthenticated && userStore.token) {
    store.connect(userStore.token)
    store.fetchContacts()
  }
})

onUnmounted(() => {
  store.disconnect()
})

watch(() => store.activeMessages.length, async () => {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
})

</script>
