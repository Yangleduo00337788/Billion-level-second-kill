import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { post } from '@/api/request'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Editor from '@/views/Editor.vue'
import ArticleDetail from '@/views/ArticleDetail.vue'
import PromptList from '@/views/PromptList.vue'
import PromptDetail from '@/views/PromptDetail.vue'
import PromptCreate from '@/views/PromptCreate.vue'
import Chat from '@/views/Chat.vue'
import ImChat from '@/views/ImChat.vue'
import UserProfile from '@/views/UserProfile.vue'
import SearchResults from '@/views/SearchResults.vue'
import Notifications from '@/views/Notifications.vue'
import AdminLayout from '@/views/admin/AdminLayout.vue'
import AdminDashboard from '@/views/admin/AdminDashboard.vue'
import AdminUsers from '@/views/admin/AdminUsers.vue'
import AdminArticles from '@/views/admin/AdminArticles.vue'
import AdminPrompts from '@/views/admin/AdminPrompts.vue'
import AdminComments from '@/views/admin/AdminComments.vue'
import AdminCategories from '@/views/admin/AdminCategories.vue'
import AdminAnnouncements from '@/views/admin/AdminAnnouncements.vue'
import AdminReports from '@/views/admin/AdminReports.vue'
import AdminRankings from '@/views/admin/AdminRankings.vue'
import AdminAuditLogs from '@/views/admin/AdminAuditLogs.vue'
import AdminAIStats from '@/views/admin/AdminAIStats.vue'
import AdminSystemConfig from '@/views/admin/AdminSystemConfig.vue'
import AdminNotifications from '@/views/admin/AdminNotifications.vue'
import AdminSettings from '@/views/admin/AdminSettings.vue'
import AdminUserTags from '@/views/admin/AdminUserTags.vue'
import AdminContentReviews from '@/views/admin/AdminContentReviews.vue'
import AdminRecommend from '@/views/admin/AdminRecommend.vue'
import AdminSensitiveWords from '@/views/admin/AdminSensitiveWords.vue'
import AdminIPBlacklist from '@/views/admin/AdminIPBlacklist.vue'
import AdminPointsRules from '@/views/admin/AdminPointsRules.vue'
import AdminInviteCodes from '@/views/admin/AdminInviteCodes.vue'
import AdminLoginLogs from '@/views/admin/AdminLoginLogs.vue'
import AdminSystemLogs from '@/views/admin/AdminSystemLogs.vue'
import AdminPageStats from '@/views/admin/AdminPageStats.vue'

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
    { path: '/prompt/create', name: 'PromptCreate', component: PromptCreate, meta: { title: '发布 Prompt', requiresAuth: true } },
    { path: '/prompt/:id', name: 'PromptDetail', component: PromptDetail, meta: { title: 'Prompt 详情' } },
    { path: '/chat', name: 'Chat', component: Chat, meta: { title: 'AI 助手', requiresAuth: true } },
    { path: '/user/:id', name: 'UserProfile', component: UserProfile, meta: { title: '用户主页' } },
    { path: '/im', name: 'ImChat', component: ImChat, meta: { title: '即时通讯', requiresAuth: true } },
    { path: '/search', name: 'SearchResults', component: SearchResults, meta: { title: '搜索' } },
    { path: '/notifications', name: 'Notifications', component: Notifications, meta: { title: '消息通知', requiresAuth: true } },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { title: '管理后台', requiresAuth: true, requiresAdmin: true },
      children: [
        { path: '', name: 'AdminDashboard', component: AdminDashboard, meta: { title: '仪表盘' } },
        { path: 'users', name: 'AdminUsers', component: AdminUsers, meta: { title: '用户管理' } },
        { path: 'user-tags', name: 'AdminUserTags', component: AdminUserTags, meta: { title: '用户标签' } },
        { path: 'articles', name: 'AdminArticles', component: AdminArticles, meta: { title: '文章管理' } },
        { path: 'prompts', name: 'AdminPrompts', component: AdminPrompts, meta: { title: 'Prompt 管理' } },
        { path: 'comments', name: 'AdminComments', component: AdminComments, meta: { title: '评论管理' } },
        { path: 'categories', name: 'AdminCategories', component: AdminCategories, meta: { title: '分类管理' } },
        { path: 'content-reviews', name: 'AdminContentReviews', component: AdminContentReviews, meta: { title: '内容审核' } },
        { path: 'recommend', name: 'AdminRecommend', component: AdminRecommend, meta: { title: '推荐位管理' } },
        { path: 'announcements', name: 'AdminAnnouncements', component: AdminAnnouncements, meta: { title: '公告管理' } },
        { path: 'reports', name: 'AdminReports', component: AdminReports, meta: { title: '举报管理' } },
        { path: 'rankings', name: 'AdminRankings', component: AdminRankings, meta: { title: '排行榜' } },
        { path: 'audit-logs', name: 'AdminAuditLogs', component: AdminAuditLogs, meta: { title: '操作日志' } },
        { path: 'login-logs', name: 'AdminLoginLogs', component: AdminLoginLogs, meta: { title: '登录日志' } },
        { path: 'sensitive-words', name: 'AdminSensitiveWords', component: AdminSensitiveWords, meta: { title: '敏感词管理' } },
        { path: 'ip-blacklist', name: 'AdminIPBlacklist', component: AdminIPBlacklist, meta: { title: 'IP 黑名单' } },
        { path: 'points-rules', name: 'AdminPointsRules', component: AdminPointsRules, meta: { title: '积分规则' } },
        { path: 'invite-codes', name: 'AdminInviteCodes', component: AdminInviteCodes, meta: { title: '邀请码管理' } },
        { path: 'ai-stats', name: 'AdminAIStats', component: AdminAIStats, meta: { title: 'AI 统计' } },
        { path: 'page-stats', name: 'AdminPageStats', component: AdminPageStats, meta: { title: '访问统计' } },
        { path: 'system-logs', name: 'AdminSystemLogs', component: AdminSystemLogs, meta: { title: '系统日志' } },
        { path: 'configs', name: 'AdminSystemConfig', component: AdminSystemConfig, meta: { title: '系统配置' } },
        { path: 'notifications', name: 'AdminNotifications', component: AdminNotifications, meta: { title: '通知管理' } },
        { path: 'settings', name: 'AdminSettings', component: AdminSettings, meta: { title: '系统设置' } }
      ]
    }
  ]
})

// Record page view on route change
router.afterEach((to) => {
  // Skip admin pages and API calls
  if (!to.path.startsWith('/admin') && !to.path.startsWith('/api/')) {
    post('/page-view', { path: to.fullPath }).catch(() => {})
  }
})

router.beforeEach((to, _from) => {
  document.title = `${to.meta.title || '推理引擎'} - 推理引擎`
  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) return { name: 'Login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guest) {
    const userStore = useUserStore()
    if (userStore.isAuthenticated) return { name: 'Home' }
  }
  if (to.meta.requiresAdmin) {
    const userStore = useUserStore()
    if (!userStore.isAuthenticated) return { name: 'Login', query: { redirect: to.fullPath } }
    if (!userStore.user || userStore.user.role !== 'admin') return { name: 'Home' }
  }
})

export default router