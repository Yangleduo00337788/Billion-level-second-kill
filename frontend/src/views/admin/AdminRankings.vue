<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">排行榜</h1>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Hot Articles -->
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-4 h-4 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" /></svg>
          热门文章
        </h3>
        <div class="space-y-3">
          <div v-for="(item, i) in hotArticles" :key="item.id" class="flex items-start gap-3 p-2 rounded-lg hover:bg-gray-50 transition">
            <span class="flex-shrink-0 w-5 h-5 rounded-full text-xs font-bold flex items-center justify-center" :class="i < 3 ? 'bg-orange-500 text-white' : 'bg-gray-100 text-gray-500'">{{ i + 1 }}</span>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-dark truncate">{{ item.title }}</p>
              <p class="text-xs text-gray-400 mt-0.5">
                <span class="text-blue-500">{{ item.view_count || 0 }}</span> 浏览 ·
                <span class="text-red-400">{{ item.like_count || 0 }}</span> 点赞
              </p>
            </div>
          </div>
          <p v-if="!hotArticles.length" class="text-sm text-gray-400 text-center py-6">暂无数据</p>
        </div>
      </div>

      <!-- Hot Prompts -->
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
          热门 Prompt
        </h3>
        <div class="space-y-3">
          <div v-for="(item, i) in hotPrompts" :key="item.id" class="flex items-start gap-3 p-2 rounded-lg hover:bg-gray-50 transition">
            <span class="flex-shrink-0 w-5 h-5 rounded-full text-xs font-bold flex items-center justify-center" :class="i < 3 ? 'bg-purple-500 text-white' : 'bg-gray-100 text-gray-500'">{{ i + 1 }}</span>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-dark truncate">{{ item.title }}</p>
              <p class="text-xs text-gray-400 mt-0.5">
                <span class="text-purple-500">{{ item.usage_count || 0 }}</span> 使用 ·
                <span class="text-red-400">{{ item.like_count || 0 }}</span> 点赞
              </p>
            </div>
          </div>
          <p v-if="!hotPrompts.length" class="text-sm text-gray-400 text-center py-6">暂无数据</p>
        </div>
      </div>

      <!-- Active Users -->
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" /></svg>
          活跃用户
        </h3>
        <div class="space-y-3">
          <div v-for="(item, i) in activeUsers" :key="item.id" class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition">
            <span class="flex-shrink-0 w-5 h-5 rounded-full text-xs font-bold flex items-center justify-center" :class="i < 3 ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-500'">{{ i + 1 }}</span>
            <n-avatar :src="item.avatar" round size="small" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-dark truncate">{{ item.username }}</p>
              <p class="text-xs text-gray-400">{{ item.article_count || 0 }} 文章 · {{ item.fans_count || 0 }} 粉丝</p>
            </div>
          </div>
          <p v-if="!activeUsers.length" class="text-sm text-gray-400 text-center py-6">暂无数据</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NAvatar } from 'naive-ui'
import { get } from '@/api/request'

const hotArticles = ref<any[]>([])
const hotPrompts = ref<any[]>([])
const activeUsers = ref<any[]>([])

async function fetchData() {
  try {
    const [articlesRes, promptsRes, usersRes] = await Promise.all([
      get<any>('/admin/hot-articles'),
      get<any>('/admin/hot-prompts'),
      get<any>('/admin/active-users')
    ])
    hotArticles.value = Array.isArray(articlesRes.data) ? articlesRes.data : []
    hotPrompts.value = Array.isArray(promptsRes.data) ? promptsRes.data : []
    activeUsers.value = Array.isArray(usersRes.data) ? usersRes.data : []
  } catch {}
}

onMounted(fetchData)
</script>