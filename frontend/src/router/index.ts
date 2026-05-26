import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Editor from '@/views/Editor.vue'
import ArticleDetail from '@/views/ArticleDetail.vue'
import PromptList from '@/views/PromptList.vue'
import PromptDetail from '@/views/PromptDetail.vue'
import Chat from '@/views/Chat.vue'
import ImChat from '@/views/ImChat.vue'
import UserProfile from '@/views/UserProfile.vue'
import SearchResults from '@/views/SearchResults.vue'
import AdminLayout from '@/views/admin/AdminLayout.vue'
import AdminDashboard from '@/views/admin/AdminDashboard.vue'
import AdminUsers from '@/views/admin/AdminUsers.vue'
import AdminArticles from '@/views/admin/AdminArticles.vue'
import AdminPrompts from '@/views/admin/AdminPrompts.vue'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    guest?: boolean
    requiresAdmin?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'Home', component: Home, meta: { title: '首页' } },
    { path: '/login', name: 'Login', component: Login, meta: { title: '登录', guest: true } },
    { path: '/register', name: 'Register', component: Register, meta: { title: '注册', guest: true } },
    { path: '/editor', name: 'Editor', component: Editor, meta: { title: '写文章', requiresAuth: true } },
    { path: '/editor/:id', name: 'EditorEdit', component: Editor, meta: { title: '编辑文章', requiresAuth: true } },
    { path: '/article/:id', name: 'ArticleDetail', component: ArticleDetail, meta: { title: '文章详情' } },
    { path: '/prompt', name: 'PromptList', component: PromptList, meta: { title: 'Prompt 社区' } },
    { path: '/prompt/:id', name: 'PromptDetail', component: PromptDetail, meta: { title: 'Prompt 详情' } },
    { path: '/chat', name: 'Chat', component: Chat, meta: { title: 'AI 助手', requiresAuth: true } },
    { path: '/user/:id', name: 'UserProfile', component: UserProfile, meta: { title: '用户主页' } },
    { path: '/im', name: 'ImChat', component: ImChat, meta: { title: '即时通讯', requiresAuth: true } },
    { path: '/search', name: 'SearchResults', component: SearchResults, meta: { title: '搜索' } },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { title: '管理后台', requiresAuth: true, requiresAdmin: true },
      children: [
        { path: '', name: 'AdminDashboard', component: AdminDashboard, meta: { title: '仪表盘' } },
        { path: 'users', name: 'AdminUsers', component: AdminUsers, meta: { title: '用户管理' } },
        { path: 'articles', name: 'AdminArticles', component: AdminArticles, meta: { title: '文章管理' } },
        { path: 'prompts', name: 'AdminPrompts', component: AdminPrompts, meta: { title: 'Prompt 管理' } }
      ]
    }
  ]
})

router.beforeEach((to, _from) => {
  document.title = `${to.meta.title || '推理引擎'} - 推理引擎`

  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) {
      return { name: 'Login', query: { redirect: to.fullPath } }
    }
  }

  if (to.meta.guest && to.name === 'Login') {
    const userStore = useUserStore()
    if (userStore.isAuthenticated) {
      return { name: 'Home' }
    }
  }

  if (to.meta.requiresAdmin) {
    const userStore = useUserStore()
    if (userStore.user?.role !== 'admin') {
      return { name: 'Home' }
    }
  }
})

export default router
