import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { get } from '@/api/request'
import type { User, ChatMessage } from '@/types/api'

export const useChatStore = defineStore('chat', () => {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const contacts = ref<User[]>([])
  const messages = ref<Record<number, ChatMessage[]>>({})
  const activeContactId = ref<number | null>(null)
  const onlineUsers = ref<number[]>([])

  const activeMessages = computed(() => {
    if (!activeContactId.value) return []
    return messages.value[activeContactId.value] || []
  })

  function connect(token: string) {
    if (ws.value) ws.value.close()
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const socket = new WebSocket(`${protocol}//${host}/ws?token=${token}`)

    socket.onopen = () => { connected.value = true }
    socket.onclose = () => { connected.value = false }
    socket.onerror = () => { connected.value = false }

    socket.onmessage = (event) => {
      try {
        const msg: ChatMessage = JSON.parse(event.data)
        const chatId = msg.sender_id === activeContactId.value ? msg.receiver_id : msg.sender_id
        if (!messages.value[chatId]) messages.value[chatId] = []
        messages.value[chatId]!.push(msg)
      } catch {}
    }

    ws.value = socket
  }

  function disconnect() {
    if (ws.value) ws.value.close()
    ws.value = null
    connected.value = false
  }

  function sendMessage(content: string, receiverId: number) {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) return
    ws.value.send(JSON.stringify({
      content,
      receiver_id: receiverId,
      msg_type: 'text'
    }))
  }

  async function fetchContacts() {
    try {
      const res = await get<{ online_users: number[] }>('/chat/list')
      onlineUsers.value = res.data.online_users || []
    } catch {}
  }

  function setActiveContact(userId: number) {
    activeContactId.value = userId
    if (!messages.value[userId]) messages.value[userId] = []
  }

  function isOnline(userId: number) {
    return onlineUsers.value.includes(userId)
  }

  return { connected, contacts, messages, onlineUsers, activeContactId, activeMessages, connect, disconnect, sendMessage, fetchContacts, setActiveContact, isOnline }
})
