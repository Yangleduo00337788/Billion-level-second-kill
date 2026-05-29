<template>
  <div class="max-w-[1200px] mx-auto px-6 lg:px-10 py-6 sm:py-8">
    <button
      class="mb-4 inline-flex items-center gap-2 text-sm text-gray-500 hover:text-dark transition-colors glass-button px-3 py-1.5 rounded-xl"
      @click="goBack"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      返回
    </button>

    <div v-if="!article" class="text-center py-20 text-gray-400">
      <n-spin size="large" />
    </div>

    <template v-else>
      <div class="flex gap-8">
        <div class="flex-1 min-w-0 max-w-3xl mx-auto">
          <div v-if="article.cover" class="aspect-[21/9] rounded-2xl overflow-hidden mb-8">
            <img :src="article.cover" :alt="article.title" class="w-full h-full object-cover" />
          </div>

          <div class="mb-8">
            <div class="flex items-center gap-2 mb-4">
              <n-tag v-if="article.is_ai" size="tiny" :bordered="false" type="warning">AI 生成</n-tag>
              <n-tag v-if="article.category" size="tiny" :bordered="false">{{ article.category.name }}</n-tag>
            </div>
            <h1 class="text-2xl sm:text-3xl font-bold text-dark leading-tight mb-4">{{ article.title }}</h1>
            <div class="flex items-center gap-4">
              <div class="flex items-center gap-2 cursor-pointer" @click="goToUser(article.user_id)">
                <n-avatar :src="article.user?.avatar" round size="small" />
                <span class="text-sm font-medium text-gray-700">{{ article.user?.username }}</span>
              </div>
              <span class="text-xs text-gray-400">{{ formatDate(article.created_at) }}</span>
              <span class="text-xs text-gray-400">{{ formatCount(article.view_count) }} 阅读</span>
            </div>
          </div>

          <div class="markdown-body mb-12" v-html="renderedContent"></div>

          <div v-if="article.tags && article.tags.length" class="flex items-center gap-2 mb-8">
            <span class="text-sm text-gray-400">标签：</span>
            <n-tag v-for="tag in parseTags(article.tags)" :key="tag" size="small" :bordered="false">
              {{ tag }}
            </n-tag>
          </div>

          <div class="flex items-center justify-center gap-6 py-6 border-t border-b border-gray-100 mb-8">
            <n-button
              :type="article.liked ? 'primary' : 'default'"
              quaternary
              @click="handleLike"
            >
              {{ article.liked ? '❤️' : '🤍' }} {{ article.like_count }}
            </n-button>
            <n-button quaternary @click="handleFavorite">
              {{ article.favorited ? '⭐' : '☆' }} {{ article.favorite_count }}
            </n-button>
            <n-button quaternary @click="shareArticle">
              🔗 分享
            </n-button>
          </div>

          <div class="mb-8">
            <h3 class="text-lg font-semibold text-dark mb-4">评论 ({{ article.comment_count }})</h3>
            <div class="flex gap-3 mb-6">
              <n-input v-model:value="commentContent" type="textarea" placeholder="写下你的评论..." :rows="2" class="flex-1" />
              <n-button type="primary" size="small" :loading="submittingComment" @click="submitComment" class="self-end">发送</n-button>
            </div>
            <div v-if="comments.length === 0" class="text-center py-8 text-gray-400 text-sm">暂无评论</div>
            <div v-for="comment in comments" :key="comment.id" class="flex gap-3 py-4 border-b border-gray-50 last:border-0">
              <n-avatar :src="comment.user?.avatar" round size="small" />
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-sm font-medium text-dark">{{ comment.user?.username }}</span>
                  <span class="text-xs text-gray-400">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="text-sm text-gray-700">{{ comment.content }}</p>
              </div>
              <div v-if="userStore.user?.id === comment.user_id">
                <n-button size="tiny" quaternary @click="deleteComment(comment.id)">删除</n-button>
              </div>
            </div>
          </div>
        </div>

        <aside class="hidden xl:block w-64 flex-shrink-0">
          <div class="sticky top-24 space-y-6">
            <div class="bg-white rounded-xl p-5 border border-gray-100">
              <h4 class="text-sm font-semibold text-dark mb-3">目录</h4>
              <div v-for="heading in toc" :key="heading.id" class="text-sm mb-2 cursor-pointer hover:text-primary transition-colors" :style="{ paddingLeft: (heading.level - 1) * 12 + 'px' }" @click="scrollToHeading(heading.id)">
                {{ heading.text }}
              </div>
              <div v-if="toc.length === 0" class="text-xs text-gray-400">无目录</div>
            </div>
          </div>
        </aside>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { marked } from 'marked'
import { get, post, del } from '@/api/request'
import { useArticleStore } from '@/stores/article'
import { useUserStore } from '@/stores/user'
import type { Comment } from '@/types/api'
import { formatDate, formatCount, parseTags } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const articleStore = useArticleStore()
const userStore = useUserStore()

const article = ref<any>(null)
const comments = ref<Comment[]>([])
const commentContent = ref('')
const submittingComment = ref(false)
const toc = ref<Array<{ id: string; text: string; level: number }>>([])

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  const html = marked(article.value.content) as string
  const withIds = html.replace(/<h([1-3])(.*?)>(.*?)<\/h\1>/gi, (_m: string, level: string, attrs: string, text: string) => {
    const id = text.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/(^-|-$)/g, '')
    return `<h${level} id="${id}"${attrs}>${text}</h${level}>`
  })
  return withIds
})

function generateToc(html: string) {
  const div = document.createElement('div')
  div.innerHTML = html
  const headings = div.querySelectorAll('h1, h2, h3')
  toc.value = Array.from(headings).map(h => ({
    id: h.id || h.textContent?.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/(^-|-$)/g, '') || '',
    text: h.textContent || '',
    level: parseInt(h.tagName[1])
  }))
}

function scrollToHeading(id: string) {
  const el = document.getElementById(id)
  if (el) el.scrollIntoView({ behavior: 'smooth' })
}

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}


function goToUser(userId: number) {
  router.push(`/user/${userId}`)
}

function shareArticle() {
  if (navigator.share) {
    navigator.share({ title: article.value?.title, url: window.location.href })
  } else {
    navigator.clipboard.writeText(window.location.href)
    message.success('链接已复制')
  }
}

async function handleLike() {
  if (!article.value || !userStore.isAuthenticated) { message.warning('请先登录'); return }
  try {
    await post(`/articles/${article.value.id}/like`)
    article.value.liked = !article.value.liked
    article.value.like_count += article.value.liked ? 1 : -1
    // 刷新用户信息（积分可能变化）
    userStore.fetchProfile()
  } catch {}
}

async function handleFavorite() {
  if (!article.value || !userStore.isAuthenticated) { message.warning('请先登录'); return }
  try {
    await post(`/articles/${article.value.id}/favorite`)
    article.value.favorited = !article.value.favorited
    article.value.favorite_count += article.value.favorited ? 1 : -1
  } catch {}
}

async function fetchComments() {
  const id = Number(route.params.id)
  try {
    const res = await get<{items: Comment[]}>(`/articles/${id}/comments`, { page: 1, size: 20 })
    comments.value = res.data.items || res.data || []
  } catch { comments.value = [] }
}

async function submitComment() {
  if (!commentContent.value.trim()) return
  if (!userStore.isAuthenticated) { message.warning('请先登录'); return }
  submittingComment.value = true
  try {
    await post(`/articles/${route.params.id}/comments`, { content: commentContent.value })
    commentContent.value = ''
    message.success('评论成功')
    fetchComments()
    if (article.value) article.value.comment_count++
    // 刷新用户信息（积分可能变化）
    userStore.fetchProfile()
  } catch { message.error('评论失败') }
  finally { submittingComment.value = false }
}

async function deleteComment(commentId: number) {
  try {
    await del(`/comments/${commentId}`)
    comments.value = comments.value.filter(c => c.id !== commentId)
    if (article.value) article.value.comment_count--
    message.success('已删除')
  } catch { message.error('删除失败') }
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (isNaN(id)) {
    message.error('文章不存在')
    router.push('/')
    return
  }
  try {
    const data = await articleStore.fetchById(id)
    article.value = data
  } catch {
    message.error('文章不存在')
    router.push('/')
  }
  fetchComments()

  nextTick(() => {
    generateToc(renderedContent.value)
  })
})
</script>
