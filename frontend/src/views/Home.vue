<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="flex gap-8">
      <div class="flex-1 min-w-0">
        <div class="flex gap-3 mb-6 overflow-x-auto pb-2 scrollbar-none">
          <n-button
            v-for="tag in hotTags"
            :key="tag"
            :type="selectedTag === tag ? 'primary' : 'default'"
            size="tiny"
            :bordered="false"
            class="flex-shrink-0"
            @click="selectedTag = tag"
          >
            {{ tag }}
          </n-button>
        </div>

        <div v-if="loading && articles.length === 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="i in 6" :key="i" class="bg-white rounded-xl p-4 animate-pulse">
            <div class="bg-gray-100 h-40 rounded-lg mb-3"></div>
            <div class="bg-gray-100 h-4 rounded w-3/4 mb-2"></div>
            <div class="bg-gray-100 h-3 rounded w-full mb-1"></div>
            <div class="bg-gray-100 h-3 rounded w-2/3"></div>
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
                class="bg-white rounded-xl overflow-hidden border border-gray-100 transition-all duration-200 hover:shadow-lg hover:-translate-y-1 cursor-pointer"
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
                      <n-avatar
                        :src="article.user?.avatar"
                        round
                        size="small"
                      />
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

          <div v-if="loading" class="text-center py-8 text-gray-400">
            加载中...
          </div>

          <div v-if="hasMore" ref="sentinel" class="h-4"></div>
        </template>
      </div>

      <aside class="hidden lg:block w-80 flex-shrink-0">
        <div class="sticky top-24 space-y-6">
          <div class="bg-white rounded-xl p-5 border border-gray-100">
            <h3 class="text-sm font-semibold text-dark mb-4">热门文章</h3>
            <div v-for="(item, i) in hotArticles" :key="item.id" class="flex items-start gap-3 mb-4 last:mb-0 cursor-pointer group" @click="goToArticle(item.id)">
              <span class="text-lg font-bold text-gray-200 w-6 flex-shrink-0 group-hover:text-primary transition-colors">{{ i + 1 }}</span>
              <div class="min-w-0">
                <p class="text-sm text-gray-700 line-clamp-2 group-hover:text-primary transition-colors">{{ item.title }}</p>
                <span class="text-xs text-gray-400">{{ formatCount(item.view_count) }} 阅读</span>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl p-5 border border-gray-100">
            <h3 class="text-sm font-semibold text-dark mb-4">推荐作者</h3>
            <div v-for="creator in recommendedCreators" :key="creator.id" class="flex items-center gap-3 mb-4 last:mb-0 cursor-pointer" @click="goToArticle(creator.id)">
              <n-avatar :src="creator.avatar" round size="small" />
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-dark truncate">{{ creator.username }}</p>
                <p class="text-xs text-gray-400 truncate">{{ creator.bio || '这个人很懒，什么都没写' }}</p>
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

const router = useRouter()
const articleStore = useArticleStore()

const selectedTag = ref('推荐')
const sentinel = ref<HTMLElement | null>(null)

const hotTags = ['推荐', '最新', '热门', 'AI', '编程', '设计', '产品', '科技']

const articles = computed(() => articleStore.articles)
const loading = computed(() => articleStore.loading)
const hotArticles = computed(() => articleStore.hot)
const hasMore = computed(() => articles.value.length < articleStore.total)

const recommendedCreators = ref<Array<{ id: number; username: string; avatar: string; bio: string }>>([])

function goToArticle(id: number) {
  router.push(`/article/${id}`)
}

async function loadMore() {
  if (articleStore.loading) return
  await articleStore.fetchList({ page: Math.floor(articles.value.length / 10) + 1, size: 10 })
}

const creatorList = ref<any[]>([])

async function fetchRecommendedCreators() {
  try {
    const res = await get('/user/list', { page: 1, size: 5, sort: 'article_count' })
    recommendedCreators.value = res.data?.items || []
  } catch {}
}

watch(selectedTag, () => {
  articleStore.fetchList({ page: 1, size: 10 })
})

onMounted(async () => {
  await articleStore.fetchList({ page: 1, size: 10 })
  articleStore.fetchHot()
  fetchRecommendedCreators()

  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value) {
      loadMore()
    }
  }, { threshold: 0.1 })

  if (sentinel.value) {
    observer.observe(sentinel.value)
  }

  onUnmounted(() => observer.disconnect())
})
</script>
