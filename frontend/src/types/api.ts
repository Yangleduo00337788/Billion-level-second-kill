export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface User {
  id: number
  username: string
  email: string
  avatar: string
  bio: string
  role: string
  level: number
  follow_count: number
  fans_count: number
  article_count: number
  points: number
  created_at: string
}

export interface Article {
  id: number
  user_id: number
  title: string
  content: string
  summary: string
  cover: string
  status: string
  category_id: number
  tags: string[]
  view_count: number
  like_count: number
  comment_count: number
  favorite_count: number
  is_ai: boolean
  created_at: string
  updated_at: string
  user?: User
  category?: Category
  liked?: boolean
  favorited?: boolean
}

export interface Category {
  id: number
  name: string
  desc: string
  sort: number
}

export interface Comment {
  id: number
  article_id: number
  user_id: number
  parent_id: number
  content: string
  like_count: number
  created_at: string
  user: User
}

export interface Prompt {
  id: number
  user_id: number
  title: string
  content: string
  description: string
  category: string
  model: string
  tags: string[]
  usage_count: number
  like_count: number
  rating: number
  created_at: string
  user: User
}

export interface ChatMessage {
  id: number
  chat_id: number
  sender_id: number
  receiver_id: number
  content: string
  msg_type: string
  file_url: string
  status: string
  created_at: string
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  size: number
}
