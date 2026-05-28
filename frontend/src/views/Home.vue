<template>
  <div class="max-w-[1440px] mx-auto px-6 lg:px-10 py-8">
    <div class="flex gap-8">
      <!-- Left Sidebar: Categories -->
      <aside class="hidden md:block w-52 flex-shrink-0">
        <div class="sticky top-24">
          <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
            <button
              class="w-full flex items-center justify-between px-4 py-3 text-sm font-semibold text-dark border-b border-gray-100 hover:bg-gray-50 transition-colors"
              @click="sidebarCollapsed = !sidebarCollapsed"
            >
              <span>分类导航</span>
              <svg class="w-4 h-4 text-gray-400 transition-transform duration-300" :class="{ '-rotate-90': sidebarCollapsed }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            <div class="overflow-hidden transition-all duration-300" :style="{ maxHeight: sidebarCollapsed ? '0' : '500px' }">
              <div class="py-2">
                <button
                  v-for="tag in hotTags"
                  :key="tag"
                  class="w-full text-left px-4 py-2.5 text-sm transition-all duration-200 flex items-center gap-3"
                  :class="selectedTag === tag
                    ? 'text-primary font-medium bg-primary/5 border-r-2 border-primary'
                    : 'text-gray-600 hover:text-dark hover:bg-gray-50'"
                  @click="selectedTag = tag"
                >
                  <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="selectedTag === tag ? 'bg-primary' : 'bg-gray-300'"></span>
                  {{ tag }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- Main Content -->
      <div class="flex-1 min-w-0">
        <!-- Featured Recommendations -->
        <div v-if="featuredArticles.length > 0" class="mb-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-bold text-dark">精选推荐</h2>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="item in featuredArticles"
              :key="item.recommend_id"
              class="bg-gradient-to-br from-primary/5 to-primary/10 rounded-2xl border border-primary/20 overflow-hidden cursor-pointer hover:shadow-lg hover:-translate-y-1 transition-all duration-300"
              @click="item.type === 'article' ? goToArticle(item.id) : goToPrompt(item.id)"
            >
              <div v-if="item.cover" class="aspect-[16/9] overflow-hidden">
                <img :src="item.cover" :alt="item.title" class="w-full h-full object-cover" />
              </div>
              <div class="p-4">
                <div class="flex items-center gap-2 mb-2">
                  <span class="px-2 py-0.5 text-xs font-medium rounded-full" :class="item.type === 'article' ? 'bg-blue-100 text-blue-600' : 'bg-purple-100 text-purple-600'">
                    {{ item.type === 'article' ? '文章' : 'Prompt' }}
                  </span>
                  <span class="px-2 py-0.5 text-xs bg-yellow-100 text-yellow-600 rounded-full">推荐</span>
                </div>
                <h3 class="font-semibold text-dark text-base mb-2 line-clamp-2">{{ item.title }}</h3>
                <p v-if="item.summary" class="text-sm text-gray-500 mb-3 line-clamp-2">{{ item.summary }}</p>
                <div v-if="item.author" class="flex items-center gap-2">
                  <n-avatar :src="item.avatar" round size="tiny" />
                  <span class="text-xs text-gray-500">{{ item.author }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Mobile: horizontal tag scroll -->
        <div class="flex gap-3 mb-6 overflow-x-auto pb-2 scrollbar-none md:hidden">
          <button
            v-for="tag in hotTags"
            :key="tag"
            class="flex-shrink-0 text-xs font-medium px-4 py-1.5 transition-all duration-300"
            :class="selectedTag === tag
              ? 'glass-button-primary text-white'
              : 'glass-button text-gray-600'"
            @click="selectedTag = tag"
          >
            {{ tag }}
          </button>
        </div>

        <div v-if="loading && articles.length === 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="i in 6" :key="i" class="bg-white rounded-2xl p-4 animate-pulse border border-gray-100">
            <div class="bg-gray-100 h-40 rounded-xl mb-3"></div>
            <div class="bg-gray-100 h-4 rounded-lg w-3/4 mb-2"></div>
            <div class="bg-gray-100 h-3 rounded-lg w-full mb-1"></div>
            <div class="bg-gray-100 h-3 rounded-lg w-2/3"></div>
          </div>
        </div>

        <div v-else-if="articles.length === 0" class="text-center py-20 text-gray-400">
          <p>暂无文章</p>
        </div>

        <template v-else>
          <div class="columns-1 md:columns-2 gap-4">
            <div
              v-for="article in articles"
              :key="article.id"
              class="break-inside-avoid mb-4"
            >
              <div
                class="bg-white rounded-2xl border border-gray-100 overflow-hidden cursor-pointer hover:shadow-md hover:-translate-y-0.5 transition-all duration-300"
                @click="goToArticle(article.id)"
              >
                <div v-if="article.cover" class="aspect-[16/10] overflow-hidden">
                  <img :src="article.cover" :alt="article.title" class="w-full h-full object-cover" />
                </div>
                <div class="p-4">
                  <div class="flex items-center gap-2 mb-2">
                    <n-tag v-if="article.is_ai" size="tiny" :bordered="false" type="warning">AI</n-tag>
                    <n-tag v-if="article.category" size="tiny" :bordered="false">{{ article.category.name }}</n-tag>
                  </div>
                  <h3 class="font-semibold text-dark text-base mb-2 line-clamp-2">{{ article.title }}</h3>
                  <p class="text-sm text-gray-500 mb-3 line-clamp-2">{{ article.summary }}</p>
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <n-avatar :src="article.user?.avatar" round size="small" />
                      <span class="text-xs text-gray-500">{{ article.user?.username }}</span>
                    </div>
                    <div class="flex items-center gap-3 text-xs text-gray-400">
                      <span>{{ formatCount(article.view_count) }} 阅读</span>
                      <span>{{ formatCount(article.like_count) }} 赞</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="loading" class="text-center py-8 text-gray-400">加载中...</div>
          <div v-if="hasMore" ref="sentinel" class="h-4"></div>
        </template>
      </div>

      <!-- Right Sidebar -->
      <aside class="hidden lg:block w-72 flex-shrink-0">
        <div class="sticky top-24 space-y-6">
          <div class="bg-white rounded-2xl border border-gray-100 p-5">
            <h3 class="text-sm font-semibold text-dark mb-4">热门文章</h3>
            <div>
              <div v-for="(item, i) in hotArticles" :key="item.id" class="flex items-start gap-3 mb-4 last:mb-0 cursor-pointer group" @click="goToArticle(item.id)">
                <span class="text-lg font-bold text-gray-200 w-6 flex-shrink-0 group-hover:text-primary transition-colors">{{ i + 1 }}</span>
                <div class="min-w-0">
                  <p class="text-sm text-gray-700 line-clamp-2 group-hover:text-primary transition-colors">{{ item.title }}</p>
                  <span class="text-xs text-gray-400">{{ formatCount(item.view_count) }} 阅读</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Sidebar Recommendations -->
          <div v-if="sidebarRecommendations.length > 0" class="bg-white rounded-2xl border border-gray-100 p-5">
            <h3 class="text-sm font-semibold text-dark mb-4">推荐阅读</h3>
            <div>
              <div v-for="item in sidebarRecommendations" :key="item.recommend_id" class="mb-4 last:mb-0 cursor-pointer group" @click="item.type === 'article' ? goToArticle(item.id) : goToPrompt(item.id)">
                <div class="flex items-start gap-3">
                  <div v-if="item.cover" class="w-16 h-12 rounded-lg overflow-hidden flex-shrink-0">
                    <img :src="item.cover" :alt="item.title" class="w-full h-full object-cover" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm text-gray-700 line-clamp-2 group-hover:text-primary transition-colors">{{ item.title }}</p>
                    <div class="flex items-center gap-1 mt-1">
                      <span class="text-xs px-1.5 py-0.5 rounded" :class="item.type === 'article' ? 'bg-blue-50 text-blue-500' : 'bg-purple-50 text-purple-500'">{{ item.type === 'article' ? '文章' : 'Prompt' }}</span>
                      <span v-if="item.author" class="text-xs text-gray-400">{{ item.author }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-2xl border border-gray-100 p-5">
            <h3 class="text-sm font-semibold text-dark mb-4">推荐作者</h3>
            <div>
              <div v-for="creator in recommendedCreators" :key="creator.id" class="flex items-center gap-3 mb-4 last:mb-0 cursor-pointer" @click="goToUser(creator.id)">
                <n-avatar :src="creator.avatar" round size="small" />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-dark truncate">{{ creator.username }}</p>
                  <p class="text-xs text-gray-400 truncate">{{ creator.bio || '这个人很懒，什么都没写' }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useArticleStore } from '@/stores/article'
import { formatCount } from '@/utils/helpers'
import { get } from '@/api/request'
import { articleApi } from '@/api/article'

const router = useRouter()
const articleStore = useArticleStore()

const selectedTag = ref('推荐')
const sentinel = ref<HTMLElement | null>(null)
const sidebarCollapsed = ref(false)

const categoryMap = ref<Record<string, number>>({})
const hotTags = ref(['推荐', '最新', '热门'])

async function fetchCategories() {
  try {
    const res = await articleApi.getCategories()
    const cats = res.data || []
    cats.forEach((c: any) => {
      categoryMap.value[c.name] = c.id
      if (!hotTags.value.includes(c.name)) hotTags.value.push(c.name)
    })
  } catch {}
}

const articles = computed(() => articleStore.articles)
const loading = computed(() => articleStore.loading)
const hotArticles = computed(() => articleStore.hot)
const hasMore = computed(() => articles.value.length < articleStore.total)

const recommendedCreators = ref<Array<{ id: number; username: string; avatar: string; bio: string }>>([])
const featuredArticles = ref<Array<any>>([])
const sidebarRecommendations = ref<Array<any>>([])

function goToUser(id: number) { router.push('/user/' + id) }
function goToArticle(id: number) { router.push('/article/' + id) }
function goToPrompt(id: number) { router.push('/prompt/' + id) }

async function fetchFeatured() {
  try {
    const res = await get<any>('/recommendations?position=homepage_top')
    featuredArticles.value = Array.isArray(res.data) ? res.data : []
  } catch {}
}

async function fetchSidebarRecommendations() {
  try {
    const res = await get<any>('/recommendations?position=homepage_sidebar')
    sidebarRecommendations.value = Array.isArray(res.data) ? res.data : []
  } catch {}
}

async function loadMore() {
  if (articleStore.loading) return
  await articleStore.fetchList({ page: Math.floor(articles.value.length / 10) + 1, page_size: 10 })
}

async function fetchRecommendedCreators() {
  try {
    const res = await get('/user/list')
    recommendedCreators.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch {}
}

watch(selectedTag, () => {
  const tag = selectedTag.value
  const params: Record<string, any> = { page: 1, page_size: 10, status: 'published' }
  if (categoryMap.value[tag]) {
    params.category = categoryMap.value[tag]
  }
  articleStore.fetchList(params)
})

onMounted(async () => {
  fetchCategories()
  fetchFeatured()
  fetchSidebarRecommendations()
  await articleStore.fetchList({ page: 1, page_size: 10 })
  articleStore.fetchHot()
  fetchRecommendedCreators()

  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value) loadMore()
  }, { threshold: 0.1 })

  if (sentinel.value) observer.observe(sentinel.value)
  onUnmounted(() => observer.disconnect())
})
</script>