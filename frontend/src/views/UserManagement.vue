<template>
  <div class="admin-container">
    <header class="admin-header">
      <div class="header-inner">
        <div class="brand">
          <div class="brand-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
            </svg>
          </div>
          <div class="brand-text">
            <span class="brand-name">Flash Mall</span>
            <span class="brand-tagline">管理后台</span>
          </div>
        </div>
        
        <nav class="main-nav">
          <router-link to="/seckill" class="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
            </svg>
            秒杀专区
          </router-link>
          <router-link to="/users" class="nav-item active">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
            </svg>
            用户管理
          </router-link>
          <router-link to="/roles" class="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path d="M2 17l10 5 10-5"/>
              <path d="M2 12l10 5 10-5"/>
            </svg>
            角色管理
          </router-link>
        </nav>
        
        <div class="user-area">
          <button class="logout-btn" @click="handleLogout">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            退出
          </button>
        </div>
      </div>
    </header>
    
    <main class="admin-main">
      <div class="page-header">
        <div class="page-title">
          <h1>用户管理</h1>
          <p>管理系统用户账户、权限和状态</p>
        </div>
        <div class="page-stats">
          <div class="stat-card">
            <div class="stat-icon users">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
              </svg>
            </div>
            <div class="stat-info">
              <span class="stat-value">{{ total }}</span>
              <span class="stat-label">总用户数</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="content-card">
        <div class="card-header">
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input type="text" placeholder="搜索用户..." v-model="searchText" />
          </div>
          <button class="refresh-btn" @click="loadUsers">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/>
              <polyline points="1 20 1 14 7 14"/>
              <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
            </svg>
            刷新
          </button>
        </div>
        
        <div class="table-wrapper" v-loading="loading">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>用户名</th>
                <th>手机号</th>
                <th>邮箱</th>
                <th>状态</th>
                <th>最后登录</th>
                <th>注册时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.id">
                <td class="id-cell">{{ user.id }}</td>
                <td class="name-cell">
                  <div class="user-avatar">{{ user.username?.charAt(0).toUpperCase() }}</div>
                  <span>{{ user.username }}</span>
                </td>
                <td>{{ user.phone || '-' }}</td>
                <td>{{ user.email || '-' }}</td>
                <td>
                  <span class="status-badge" :class="getStatusClass(user.status)">
                    {{ getStatusText(user.status) }}
                  </span>
                </td>
                <td>{{ formatDateTime(user.lastLoginTime) }}</td>
                <td>{{ formatDateTime(user.createTime) }}</td>
                <td class="actions-cell">
                  <div class="action-buttons">
                    <button v-if="user.status !== 1" class="action-btn enable" @click="handleEnable(user)" title="启用">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="20 6 9 17 4 12"/>
                      </svg>
                    </button>
                    <button v-if="user.status === 1" class="action-btn disable" @click="handleDisable(user)" title="禁用">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"/>
                        <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                      </svg>
                    </button>
                    <button class="action-btn lock" @click="handleLock(user)" title="锁定">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                        <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                      </svg>
                    </button>
                    <button class="action-btn reset" @click="handleResetPwd(user)" title="重置密码">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L15 13"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="handleDelete(user)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          
          <div v-if="filteredUsers.length === 0 && !loading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="10"/>
              <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
              <line x1="9" y1="9" x2="9.01" y2="9"/>
              <line x1="15" y1="9" x2="15.01" y2="9"/>
            </svg>
            <p>暂无用户数据</p>
          </div>
        </div>
        
        <div class="pagination-bar" v-if="total > 0">
          <span class="page-info">共 {{ total }} 条记录</span>
          <div class="pagination-controls">
            <button class="page-btn" :disabled="pageNum === 1" @click="pageNum--; loadUsers()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 18 9 12 15 6"/>
              </svg>
            </button>
            <span class="page-num">{{ pageNum }}</span>
            <button class="page-btn" @click="pageNum++; loadUsers()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 18 15 12 9 6"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </main>
    
    <el-dialog v-model="resetPwdDialog" title="" width="400px" class="custom-dialog">
      <div class="dialog-content">
        <div class="dialog-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L15 13"/>
          </svg>
        </div>
        <h3>重置密码</h3>
        <p class="dialog-desc">用户: {{ currentUser?.username }}</p>
        <div class="form-group">
          <label>新密码</label>
          <input type="password" v-model="resetPwdForm.newPassword" placeholder="请输入新密码" />
        </div>
      </div>
      <template #footer>
        <button class="dialog-btn cancel" @click="resetPwdDialog = false">取消</button>
        <button class="dialog-btn confirm" @click="confirmResetPwd">确定</button>
      </template>
    </el-dialog>

    <el-dialog v-model="lockDialog" title="" width="400px" class="custom-dialog">
      <div class="dialog-content">
        <div class="dialog-icon warning">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
          </svg>
        </div>
        <h3>锁定用户</h3>
        <p class="dialog-desc">用户: {{ currentUser?.username }}</p>
        <div class="form-group">
          <label>锁定时长（分钟）</label>
          <input type="number" v-model="lockForm.lockMinutes" min="1" max="1440" />
        </div>
      </div>
      <template #footer>
        <button class="dialog-btn cancel" @click="lockDialog = false">取消</button>
        <button class="dialog-btn confirm" @click="confirmLock">确定</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, updateUserStatus, deleteUser, lockUser, resetPassword } from '../api/user'

const router = useRouter()
const loading = ref(false)
const users = ref([])
const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchText = ref('')

const resetPwdDialog = ref(false)
const lockDialog = ref(false)
const currentUser = ref(null)
const resetPwdForm = ref({ newPassword: '' })
const lockForm = ref({ lockMinutes: 30 })

const filteredUsers = computed(() => {
  if (!searchText.value) return users.value
  const search = searchText.value.toLowerCase()
  return users.value.filter(u => 
    u.username?.toLowerCase().includes(search) ||
    u.phone?.includes(search) ||
    u.email?.toLowerCase().includes(search)
  )
})

async function loadUsers() {
  loading.value = true
  try {
    const res = await getUserList({ pageNum: pageNum.value, pageSize: pageSize.value })
    users.value = res.data || []
    total.value = res.total || users.value.length
  } catch (error) {
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

function getStatusClass(status) {
  if (status === 1) return 'active'
  if (status === 0) return 'disabled'
  return 'locked'
}

function getStatusText(status) {
  if (status === 1) return '正常'
  if (status === 0) return '禁用'
  return '锁定'
}

function formatDateTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

async function handleEnable(user) {
  try {
    await updateUserStatus(user.id, 1)
    ElMessage.success('已启用')
    loadUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDisable(user) {
  try {
    await ElMessageBox.confirm('确定要禁用该用户吗？', '提示', { type: 'warning' })
    await updateUserStatus(user.id, 0)
    ElMessage.success('已禁用')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

function handleLock(user) {
  currentUser.value = user
  lockForm.value.lockMinutes = 30
  lockDialog.value = true
}

async function confirmLock() {
  try {
    await lockUser(currentUser.value.id, lockForm.value.lockMinutes)
    ElMessage.success('已锁定')
    lockDialog.value = false
    loadUsers()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

function handleResetPwd(user) {
  currentUser.value = user
  resetPwdForm.value.newPassword = ''
  resetPwdDialog.value = true
}

async function confirmResetPwd() {
  if (!resetPwdForm.value.newPassword) {
    ElMessage.warning('请输入新密码')
    return
  }
  try {
    await resetPassword(currentUser.value.id, resetPwdForm.value.newPassword)
    ElMessage.success('密码已重置')
    resetPwdDialog.value = false
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(user) {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？此操作不可恢复！', '警告', { type: 'warning' })
    await deleteUser(user.id)
    ElMessage.success('已删除')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

function handleLogout() {
  localStorage.removeItem('accessToken')
  router.push('/login')
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=Space+Grotesk:wght@500;700&display=swap');

.admin-container {
  min-height: 100vh;
  background: #f8f9fc;
  font-family: 'DM Sans', sans-serif;
}

.admin-header {
  background: #fff;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 24px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.brand-icon svg { width: 22px; height: 22px; }

.brand-text { display: flex; flex-direction: column; }

.brand-name {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #1a1a2e;
}

.brand-tagline {
  font-size: 11px;
  color: #6366f1;
  font-weight: 600;
}

.main-nav { display: flex; gap: 8px; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item svg { width: 18px; height: 18px; }
.nav-item:hover { background: #f1f5f9; color: #1a1a2e; }
.nav-item.active { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; }

.user-area { display: flex; align-items: center; gap: 12px; }

.logout-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #fef2f2;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #ef4444;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn svg { width: 18px; height: 18px; }
.logout-btn:hover { background: #fee2e2; }

.admin-main {
  max-width: 1280px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.page-title h1 {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.page-title p {
  font-size: 14px;
  color: #64748b;
  margin: 0;
}

.page-stats { display: flex; gap: 16px; }

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.users { background: #eef2ff; color: #6366f1; }
.stat-icon svg { width: 22px; height: 22px; }

.stat-info { display: flex; flex-direction: column; }

.stat-value {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-label { font-size: 12px; color: #64748b; }

.content-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.04);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: #f8fafc;
  border-radius: 10px;
  width: 300px;
}

.search-box svg { width: 18px; height: 18px; color: #94a3b8; }
.search-box input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #1a1a2e;
  outline: none;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #f8fafc;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn svg { width: 18px; height: 18px; }
.refresh-btn:hover { background: #f1f5f9; }

.table-wrapper { overflow-x: auto; }

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  padding: 14px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  background: #f8fafc;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table td {
  padding: 16px;
  font-size: 14px;
  color: #1a1a2e;
  border-bottom: 1px solid #f1f5f9;
}

.id-cell {
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 600;
  color: #64748b;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.active { background: #dcfce7; color: #22c55e; }
.status-badge.disabled { background: #fef2f2; color: #ef4444; }
.status-badge.locked { background: #fef3c7; color: #d97706; }

.actions-cell { width: 140px; }

.action-buttons { display: flex; gap: 6px; }

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn svg { width: 16px; height: 16px; }
.action-btn.enable { background: #dcfce7; color: #22c55e; }
.action-btn.enable:hover { background: #bbf7d0; }
.action-btn.disable { background: #fef2f2; color: #ef4444; }
.action-btn.disable:hover { background: #fee2e2; }
.action-btn.lock { background: #fef3c7; color: #d97706; }
.action-btn.lock:hover { background: #fde68a; }
.action-btn.reset { background: #eef2ff; color: #6366f1; }
.action-btn.reset:hover { background: #e0e7ff; }
.action-btn.delete { background: #f1f5f9; color: #64748b; }
.action-btn.delete:hover { background: #e2e8f0; color: #ef4444; }

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #94a3b8;
}

.empty-state svg { width: 48px; height: 48px; margin-bottom: 12px; }
.empty-state p { margin: 0; font-size: 14px; }

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-top: 1px solid #f1f5f9;
}

.page-info { font-size: 13px; color: #64748b; }

.pagination-controls { display: flex; align-items: center; gap: 8px; }

.page-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn svg { width: 16px; height: 16px; color: #64748b; }
.page-btn:hover:not(:disabled) { border-color: #6366f1; }
.page-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.page-num {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #1a1a2e;
  padding: 0 8px;
}

.dialog-content { text-align: center; padding: 10px 0; }

.dialog-icon {
  width: 60px;
  height: 60px;
  background: #eef2ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  color: #6366f1;
}

.dialog-icon.warning { background: #fef3c7; color: #d97706; }
.dialog-icon svg { width: 28px; height: 28px; }

.dialog-content h3 {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.dialog-desc { font-size: 13px; color: #64748b; margin: 0 0 20px 0; }

.form-group { text-align: left; margin-bottom: 8px; }

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #64748b;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  font-size: 14px;
  color: #1a1a2e;
  outline: none;
  transition: all 0.2s;
}

.form-group input:focus { border-color: #6366f1; }

.dialog-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-btn.cancel {
  background: #f1f5f9;
  color: #64748b;
}

.dialog-btn.cancel:hover { background: #e2e8f0; }

.dialog-btn.confirm {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
}

.dialog-btn.confirm:hover { filter: brightness(1.1); }
</style>
