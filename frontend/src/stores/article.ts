import { defineStore } from 'pinia'
import { ref } from 'vue'
import { articleApi } from '@/api/article'
import type { Article, PaginatedData } from '@/types/api'

export const useArticleStore = defineStore('article', () => {
  const articles = ref<Article[]>([])
  const currentArticle = ref<Article | null>(null)
  const feed = ref<Article[]>([])
  const hot = ref<Article[]>([])
  const total = ref(0)
  const page = ref(1)
  const loading = ref(false)

  async function fetchList(params?: { page?: number; size?: number; category?: number; tag?: string; status?: string }) {
    loading.value = true
    try {
      const res = await articleApi.list(params)
      if (params?.page && params.page > 1) {
        articles.value = [...articles.value, ...res.data.items]
      } else {
        articles.value = res.data.items
      }
      total.value = res.data.total
      page.value = params?.page || 1
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: number) {
    const res = await articleApi.getById(id)
    currentArticle.value = res.data
    return res.data
  }

  async function create(data: { title: string; content: string; summary?: string; cover?: string; category_id?: number }) {
    const res = await articleApi.create(data)
    return res.data
  }

  async function update(id: number, data: Partial<Article>) {
    const res = await articleApi.update(id, data)
    currentArticle.value = res.data
    return res.data
  }

  async function remove(id: number) {
    await articleApi.delete(id)
    articles.value = articles.value.filter(a => a.id !== id)
  }

  async function like(id: number) {
    await articleApi.like(id)
    const article = articles.value.find(a => a.id === id)
    if (article) {
      article.liked = !article.liked
      article.like_count += article.liked ? 1 : -1
    }
  }

  async function favorite(id: number) {
    await articleApi.favorite(id)
    const article = articles.value.find(a => a.id === id)
    if (article) {
      article.favorited = !article.favorited
      article.favorite_count += article.favorited ? 1 : -1
    }
  }

  async function fetchFeed(params?: { page?: number; size?: number }) {
    const res = await articleApi.getFeed(params)
    feed.value = res.data.items
  }

  async function fetchHot() {
    const res = await articleApi.getHot()
    hot.value = res.data
  }

  return { articles, currentArticle, feed, hot, total, page, loading, fetchList, fetchById, create, update, remove, like, favorite, fetchFeed, fetchHot }
})
