<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div v-if="loading" class="text-center py-20"><n-spin size="large" /></div>
    <template v-else-if="profile">
      <div class="bg-white rounded-2xl border border-gray-100 p-6 lg:p-8 mb-6">
        <div class="flex items-start gap-4 lg:gap-6">
          <n-avatar :src="profile.avatar" round :size="64" class="lg:w-20 lg:h-20" />
          <div class="flex-1 min-w-0">
            <h1 class="text-xl lg:text-2xl font-bold text-dark mb-1">{{ profile.username }}</h1>
            <p class="text-sm text-gray-500 mb-3">{{ profile.bio || '这个人很懒，什么都没写' }}</p>
            <div class="flex items-center gap-4 lg:gap-6 text-sm">
              <span class="text-gray-500"><strong class="text-dark">{{ formatCount(profile.article_count) }}</strong> 文章</span>
              <span class="text-gray-500"><strong class="text-dark">{{ formatCount(profile.follow_count) }}</strong> 关注</span>
              <span class="text-gray-500"><strong class="text-dark">{{ formatCount(profile.fans_count) }}</strong> 粉丝</span>
            </div>
          </div>
          <div v-if="isOwner" class="flex-shrink-0">
            <n-button size="small" quaternary @click="showEditModal = true">编辑资料</n-button>
          </div>
          <div v-else class="flex-shrink-0">
            <n-button size="small" :type="isFollowing ? 'default' : 'primary'" :loading="followLoading" @click="handleFollow">
              {{ isFollowing ? '已关注' : '关注' }}
            </n-button>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-1 mb-6 border-b border-gray-100 overflow-x-auto pb-px">
        <n-button v-for="tab in tabs" :key="tab.key" :type="activeTab === tab.key ? 'primary' : 'default'" :bordered="false" size="small" @click="activeTab = tab.key">{{ tab.label }}</n-button>
      </div>

      <div v-if="activeTab === 'articles'">
        <div v-if="userArticles.length === 0" class="text-center py-12 text-gray-400 text-sm">暂无文章</div>
        <div v-for="article in userArticles" :key="article.id" class="bg-white rounded-xl border border-gray-100 p-5 mb-4 cursor-pointer hover:shadow-sm transition-shadow" @click="router.push(`/article/${article.id}`)">
          <h3 class="font-semibold text-dark mb-2">{{ article.title }}</h3>
          <p class="text-sm text-gray-500 line-clamp-2 mb-3">{{ article.summary }}</p>
          <div class="flex items-center gap-4 text-xs text-gray-400">
            <span>{{ formatDate(article.created_at) }}</span>
            <span>{{ formatCount(article.view_count) }} 阅读</span>
            <span>{{ formatCount(article.like_count) }} 赞</span>
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'prompts'">
        <div v-if="userPrompts.length === 0" class="text-center py-12 text-gray-400 text-sm">暂无 Prompt</div>
        <div v-for="p in userPrompts" :key="p.id" class="bg-white rounded-xl border border-gray-100 p-5 mb-4 cursor-pointer hover:shadow-sm transition-shadow" @click="router.push(`/prompt/${p.id}`)">
          <h3 class="font-semibold text-dark mb-2">{{ p.title }}</h3>
          <p class="text-sm text-gray-500 line-clamp-2 mb-3">{{ p.description }}</p>
          <div class="flex items-center gap-3 text-xs text-gray-400">
            <n-tag size="tiny" :bordered="false">{{ p.category }}</n-tag>
            <span>使用 {{ formatCount(p.usage_count) }} 次</span>
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'favorites'">
        <div v-if="userFavorites.length === 0" class="text-center py-12 text-gray-400 text-sm">暂无收藏</div>
        <div v-for="article in userFavorites" :key="article.id" class="bg-white rounded-xl border border-gray-100 p-5 mb-4 cursor-pointer hover:shadow-sm transition-shadow" @click="router.push(`/article/${article.id}`)">
          <h3 class="font-semibold text-dark mb-2">{{ article.title }}</h3>
          <p class="text-sm text-gray-500 line-clamp-2">{{ article.summary }}</p>
        </div>
      </div>
    </template>
  </div>

  <n-modal v-model:show="showEditModal" title="编辑资料" preset="card" style="width: 420px;" closable>
    <n-form>
      <n-form-item label="用户名" required>
        <n-input v-model:value="editForm.username" />
      </n-form-item>
      <n-form-item label="简介">
        <n-input v-model:value="editForm.bio" type="textarea" :rows="3" />
      </n-form-item>
      <n-form-item label="头像链接">
        <n-input v-model:value="editForm.avatar" placeholder="输入头像 URL" />
      </n-form-item>
    </n-form>
    <template #footer>
      <div class="flex justify-end gap-2">
        <n-button @click="showEditModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="saveProfile">保存</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { get, put, post, del } from '@/api/request'
import type { User, Article, Prompt } from '@/types/api'
import { formatDate, formatCount } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

const profile = ref<User | null>(null)
const loading = ref(true)
const followLoading = ref(false)
const isFollowing = ref(false)
const activeTab = ref('articles')
const userArticles = ref<Article[]>([])
const userPrompts = ref<Prompt[]>([])
const userFavorites = ref<Article[]>([])
const showEditModal = ref(false)
const saving = ref(false)

const tabs = [
  { key: 'articles', label: '文章' },
  { key: 'prompts', label: 'Prompt' },
  { key: 'favorites', label: '收藏' }
]

const editForm = ref({ username: '', bio: '', avatar: '' })

const isOwner = computed(() => userStore.user?.id === profile.value?.id)

async function fetchProfile() {
  loading.value = true
  const id = Number(route.params.id)
  if (isNaN(id)) { router.push('/'); return }
  try {
    const res = await get<User>(`/user/${id}`)
    profile.value = res.data
    editForm.value = { username: res.data.username, bio: res.data.bio || '', avatar: res.data.avatar || '' }
    if (userStore.isAuthenticated) {
      try {
        const followRes = await get(`/user/${id}/followers`, { page: 1, size: 1 })
        isFollowing.value = followRes.data.items?.some((f: any) => f.id === userStore.user?.id) || false
      } catch {}
    }
  } catch { message.error('用户不存在'); router.push('/') }
  finally { loading.value = false }
}

async function fetchArticles() {
  const id = Number(route.params.id)
  try {
    const res = await get(`/articles`, { user_id: id, page: 1, size: 20 })
    userArticles.value = res.data.items || []
  } catch { userArticles.value = [] }
}

async function fetchPrompts() {
  const id = Number(route.params.id)
  try {
    const res = await get(`/prompts`, { user_id: id, page: 1, size: 20 })
    userPrompts.value = res.data.items || []
  } catch { userPrompts.value = [] }
}

async function fetchFavorites() {
  const id = Number(route.params.id)
  try {
    const res = await get(`/user/${id}/favorites`, { page: 1, size: 20 })
    userFavorites.value = res.data?.items || res.data || []
  } catch { userFavorites.value = [] }
}

watch(activeTab, (tab) => {
  if (tab === 'articles') fetchArticles()
  else if (tab === 'prompts') fetchPrompts()
  else if (tab === 'favorites') fetchFavorites()
})

async function handleFollow() {
  if (!userStore.isAuthenticated) { message.warning('请先登录'); return }
  followLoading.value = true
  try {
    if (isFollowing.value) {
      await del(`/user/${profile.value!.id}/follow`)
      isFollowing.value = false
      if (profile.value) profile.value.fans_count--
    } else {
      await post(`/user/${profile.value!.id}/follow`)
      isFollowing.value = true
      if (profile.value) profile.value.fans_count++
    }
  } catch { message.error('操作失败') }
  finally { followLoading.value = false }
}

async function saveProfile() {
  if (!editForm.value.username.trim()) {
    message.warning('用户名不能为空')
    return
  }
  saving.value = true
  try {
    await put('/user/profile', editForm.value)
    showEditModal.value = false
    message.success('保存成功')
    fetchProfile()
  } catch { message.error('保存失败') }
  finally { saving.value = false }
}

onMounted(() => { fetchProfile(); fetchArticles() })
</script>
