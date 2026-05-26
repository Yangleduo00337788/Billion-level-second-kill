import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types/api'
import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

const instance: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' }
})

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data
    if (res.code !== 0) {
      message.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return response
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      const msg = error.response.data?.message || '服务器错误'
      if (status === 401) {
        localStorage.removeItem('token')
        message.error('登录已过期，请重新登录')
        window.location.href = '/login'
      } else {
        message.error(msg)
      }
    } else {
      message.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export async function get<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  const res = await instance.get<ApiResponse<T>>(url, { params, ...config })
  return res.data
}

export async function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  const res = await instance.post<ApiResponse<T>>(url, data, config)
  return res.data
}

export async function put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  const res = await instance.put<ApiResponse<T>>(url, data, config)
  return res.data
}

export async function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  const res = await instance.delete<ApiResponse<T>>(url, config)
  return res.data
}

export default instance
