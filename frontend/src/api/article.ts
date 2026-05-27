import { get, post, put, del } from './request'
import type { ApiResponse, Article, PaginatedData, Category } from '@/types/api'

export const articleApi = {
  list(params?: { page?: number; page_size?: number; category?: number; tag?: string; status?: string }) {
    return get<PaginatedData<Article>>('/articles', params)
  },

  getById(id: number) {
    return get<Article>(`/articles/${id}`)
  },

  create(data: { title: string; content: string; summary?: string; cover?: string; category_id?: number }) {
    return post<Article>('/articles', data)
  },

  update(id: number, data: Partial<Article>) {
    return put<Article>(`/articles/${id}`, data)
  },

  delete(id: number) {
    return del<void>(`/articles/${id}`)
  },

  like(id: number) {
    return post<void>(`/articles/${id}/like`)
  },

  favorite(id: number) {
    return post<void>(`/articles/${id}/favorite`)
  },

  getHot() {
    return get<Article[]>('/articles/hot')
  },

  getFeed(params?: { page?: number; page_size?: number }) {
    return get<PaginatedData<Article>>('/articles/feed', params)
  },
  getCategories() {
    return get<Category[]>('/articles/categories')
  }
}
