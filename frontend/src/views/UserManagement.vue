<template>
  <Layout>
    <div class="admin-page">
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
    </div>

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
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, updateUserStatus, deleteUser, lockUser, resetPassword } from '../api/user'
import Layout from '../components/Layout.vue'

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

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.admin-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.page-title p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.page-stats {
  display: flex;
  gap: 16px;
}

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

.stat-icon.users {
  background: #FFF0E6;
  color: #FF5000;
}

.stat-icon svg {
  width: 22px;
  height: 22px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.content-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: #f8f8f8;
  border-radius: 8px;
  width: 300px;
}

.search-box svg {
  width: 18px;
  height: 18px;
  color: #999;
}

.search-box input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #333;
  outline: none;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #f8f8f8;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn svg {
  width: 18px;
  height: 18px;
}

.refresh-btn:hover {
  background: #f0f0f0;
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  padding: 14px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #999;
  background: #fafafa;
}

.data-table td {
  padding: 16px;
  font-size: 14px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
}

.id-cell {
  font-weight: 600;
  color: #999;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
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

.status-badge.active {
  background: #E8F5E9;
  color: #4CAF50;
}

.status-badge.disabled {
  background: #FFF0E6;
  color: #FF5000;
}

.status-badge.locked {
  background: #FFF8E1;
  color: #F9A825;
}

.actions-cell {
  width: 180px;
}

.action-buttons {
  display: flex;
  gap: 6px;
}

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

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn.enable {
  background: #E8F5E9;
  color: #4CAF50;
}

.action-btn.enable:hover {
  background: #C8E6C9;
}

.action-btn.disable {
  background: #FFF0E6;
  color: #FF5000;
}

.action-btn.disable:hover {
  background: #FFD8C4;
}

.action-btn.lock {
  background: #FFF8E1;
  color: #F9A825;
}

.action-btn.lock:hover {
  background: #FFECB3;
}

.action-btn.reset {
  background: #F3E5F5;
  color: #9C27B0;
}

.action-btn.reset:hover {
  background: #E1BEE7;
}

.action-btn.delete {
  background: #f5f5f5;
  color: #999;
}

.action-btn.delete:hover {
  background: #e0e0e0;
  color: #F44336;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #ccc;
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

.page-info {
  font-size: 13px;
  color: #999;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn svg {
  width: 16px;
  height: 16px;
  color: #666;
}

.page-btn:hover:not(:disabled) {
  border-color: #FF5000;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-num {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  padding: 0 8px;
}

.dialog-content {
  text-align: center;
  padding: 10px 0;
}

.dialog-icon {
  width: 60px;
  height: 60px;
  background: #FFF0E6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  color: #FF5000;
}

.dialog-icon.warning {
  background: #FFF8E1;
  color: #F9A825;
}

.dialog-icon svg {
  width: 28px;
  height: 28px;
}

.dialog-content h3 {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin: 0 0 8px 0;
}

.dialog-desc {
  font-size: 13px;
  color: #999;
  margin: 0 0 20px 0;
}

.form-group {
  text-align: left;
  margin-bottom: 8px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #666;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  font-size: 14px;
  color: #333;
  outline: none;
  transition: all 0.2s;
}

.form-group input:focus {
  border-color: #FF5000;
}

.dialog-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-btn.cancel {
  background: #f5f5f5;
  color: #666;
}

.dialog-btn.cancel:hover {
  background: #e0e0e0;
}

.dialog-btn.confirm {
  background: linear-gradient(135deg, #FF5000, #FF3300);
  color: #fff;
}

.dialog-btn.confirm:hover {
  filter: brightness(1.1);
}
</style>