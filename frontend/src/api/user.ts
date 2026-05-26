import { get, post, put, del } from './request'
import type { ApiResponse, User, Article, PaginatedData } from '@/types/api'

export const userApi = {
  register(data: { username: string; email: string; password: string }) {
    return post<{ user: User; token: string }>('/auth/register', data)
  },

  login(data: { email: string; password: string }) {
    return post<{ user: User; token: string }>('/auth/login', data)
  },

  getProfile() {
    return get<User>('/user/profile')
  },

  updateProfile(data: Partial<User>) {
    return put<User>('/user/profile', data)
  },

  getUserById(id: number) {
    return get<User>(`/user/${id}`)
  },

  follow(id: number) {
    return post<void>(`/user/${id}/follow`)
  },

  getFollowers(id: number, params?: { page?: number; size?: number }) {
    return get<PaginatedData<User>>(`/user/${id}/followers`, params)
  },

  getFollowing(id: number, params?: { page?: number; size?: number }) {
    return get<PaginatedData<User>>(`/user/${id}/following`, params)
  },

  unfollow(id: number) {
    return del<void>(`/user/${id}/follow`)
  },

  getUserArticles(id: number, params?: { page?: number; size?: number; status?: string }) {
    return get<PaginatedData<Article>>(`/user/${id}/articles`, params)
  }
}
