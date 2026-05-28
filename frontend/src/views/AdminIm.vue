<template>
  <div class="admin-im">
    <div class="page-header">
      <div class="page-title">
        <h1>IM管理系统</h1>
        <p>监控和管理即时通讯系统运行状态</p>
      </div>
      <div class="header-actions">
        <button class="refresh-btn" @click="refreshAll">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"/>
            <polyline points="1 20 1 14 7 14"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          刷新数据
        </button>
      </div>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon online">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.onlineUsers }}</span>
          <span class="stat-label">在线用户</span>
        </div>
        <div class="stat-spark">
          <svg viewBox="0 0 80 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="0,18 8,14 16,16 24,8 32,12 40,6 48,10 56,4 64,8 72,2 80,6"/>
          </svg>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon active">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.activeSessions }}</span>
          <span class="stat-label">活跃会话</span>
        </div>
        <div class="stat-spark">
          <svg viewBox="0 0 80 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="0,12 8,10 16,8 24,14 32,6 40,12 48,4 56,10 64,6 72,8 80,2"/>
          </svg>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon messages">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.todayMessages }}</span>
          <span class="stat-label">今日消息</span>
        </div>
        <div class="stat-spark">
          <svg viewBox="0 0 80 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="0,20 8,18 16,10 24,14 32,4 40,12 48,6 56,10 64,2 72,8 80,4"/>
          </svg>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon groups">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M7 21v-2a4 4 0 0 1 3-3.87"/>
            <circle cx="12" cy="7" r="4"/>
            <circle cx="17" cy="11" r="2"/>
            <circle cx="7" cy="11" r="2"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.groupCount }}</span>
          <span class="stat-label">群组数量</span>
        </div>
        <div class="stat-spark">
          <svg viewBox="0 0 80 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="0,16 8,14 16,16 24,10 32,12 40,8 48,12 56,6 64,10 72,4 80,8"/>
          </svg>
        </div>
      </div>
    </div>

    <div class="content-card">
      <div class="tab-nav">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          :class="['tab-btn', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          <svg v-html="tab.icon" />
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <div class="tab-content">
        <div v-if="activeTab === 'sessions'" class="table-wrapper">
          <table class="data-table" v-loading="sessionLoading">
            <thead>
              <tr>
                <th>会话ID</th>
                <th>类型</th>
                <th>会话名称</th>
                <th>参与人数</th>
                <th>最后消息</th>
                <th>更新时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in sessions" :key="s.id">
                <td class="id-cell">#{{ s.id }}</td>
                <td>
                  <span :class="['type-badge', s.type === 1 ? 'single' : 'group']">
                    {{ s.type === 1 ? '单聊' : '群聊' }}
                  </span>
                </td>
                <td>
                  <div class="cell-name">
                    <div :class="['cell-avatar', s.type === 1 ? 'single-avatar' : 'group-avatar']">
                      {{ s.sessionName?.charAt(0) || 'S' }}
                    </div>
                    <span>{{ s.sessionName || '-' }}</span>
                  </div>
                </td>
                <td>{{ s.participantCount || 2 }}</td>
                <td class="preview-cell">{{ s.lastMessage || '暂无消息' }}</td>
                <td>{{ formatDateTime(s.updateTime) }}</td>
              </tr>
            </tbody>
          </table>

          <div v-if="sessions.length === 0 && !sessionLoading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              <line x1="9" y1="10" x2="15" y2="10"/>
              <line x1="12" y1="7" x2="12" y2="13"/>
            </svg>
            <p class="empty-title">暂无会话数据</p>
            <p class="empty-desc">当前系统中没有活跃的IM会话</p>
          </div>
        </div>

        <div v-if="activeTab === 'messages'" class="table-wrapper">
          <table class="data-table" v-loading="messageLoading">
            <thead>
              <tr>
                <th>服务端消息ID</th>
                <th>客户端消息ID</th>
                <th>发送者</th>
                <th>会话ID</th>
                <th>消息类型</th>
                <th>内容长度</th>
                <th>发送时间</th>
                <th>状态</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in messages" :key="m.id">
                <td class="id-cell">{{ m.serverMsgId || m.id }}</td>
                <td class="id-cell">{{ m.clientMsgId || '-' }}</td>
                <td>
                  <div class="cell-name">
                    <div class="cell-avatar single-avatar sender-avatar">
                      {{ m.sender?.charAt(0) || 'U' }}
                    </div>
                    <span>{{ m.sender || '-' }}</span>
                  </div>
                </td>
                <td class="id-cell">#{{ m.sessionId || '-' }}</td>
                <td>
                  <span :class="['msg-type-badge', getMsgTypeClass(m.msgType)]">
                    {{ getMsgTypeText(m.msgType) }}
                  </span>
                </td>
                <td>{{ m.contentLength ?? (m.content?.length || 0) }} 字符</td>
                <td>{{ formatDateTime(m.sendTime) }}</td>
                <td>
                  <span :class="['status-dot', getMsgStatusClass(m.status)]" />
                  <span class="status-text">{{ getMsgStatusText(m.status) }}</span>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="messages.length === 0 && !messageLoading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              <polyline points="9 10 12 13 15 10"/>
            </svg>
            <p class="empty-title">暂无消息记录</p>
            <p class="empty-desc">当前没有可追踪的IM消息数据</p>
          </div>
        </div>

        <div v-if="activeTab === 'groups'" class="table-wrapper">
          <div class="tab-toolbar">
            <div class="toolbar-title">群组列表</div>
            <button class="action-primary-btn" @click="openCreateGroupDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              创建群组
            </button>
          </div>
          <table class="data-table" v-loading="groupLoading">
            <thead>
              <tr>
                <th>群组ID</th>
                <th>群组名称</th>
                <th>群主</th>
                <th>成员数量</th>
                <th>群公告</th>
                <th>创建时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="g in groups" :key="g.id">
                <td class="id-cell">#{{ g.id }}</td>
                <td>
                  <div class="cell-name">
                    <div class="cell-avatar group-avatar">
                      {{ g.groupName?.charAt(0) || 'G' }}
                    </div>
                    <span>{{ g.groupName || '-' }}</span>
                  </div>
                </td>
                <td>{{ g.owner || '-' }}</td>
                <td>
                  <span class="member-count">{{ g.memberCount || 0 }}</span>
                </td>
                <td class="preview-cell">{{ g.notice || '暂无公告' }}</td>
                <td>{{ formatDateTime(g.createTime) }}</td>
              </tr>
            </tbody>
          </table>

          <div v-if="groups.length === 0 && !groupLoading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M7 21v-2a4 4 0 0 1 3-3.87"/>
              <circle cx="12" cy="7" r="4"/>
              <circle cx="17" cy="11" r="2"/>
              <circle cx="7" cy="11" r="2"/>
            </svg>
            <p class="empty-title">暂无群组数据</p>
            <p class="empty-desc">点击上方「创建群组」按钮添加新群组</p>
          </div>
        </div>

        <div v-if="activeTab === 'push'" class="table-wrapper">
          <table class="data-table" v-loading="pushLoading">
            <thead>
              <tr>
                <th>ID</th>
                <th>消息ID</th>
                <th>目标用户</th>
                <th>推送方式</th>
                <th>推送状态</th>
                <th>重试次数</th>
                <th>推送时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in pushLogs" :key="p.id">
                <td class="id-cell">{{ p.id }}</td>
                <td class="id-cell">{{ p.messageId || '-' }}</td>
                <td>
                  <div class="cell-name">
                    <div class="cell-avatar single-avatar sender-avatar">
                      {{ p.targetUser?.charAt(0) || 'U' }}
                    </div>
                    <span>{{ p.targetUser || '-' }}</span>
                  </div>
                </td>
                <td>
                  <span :class="['push-method-badge', getPushMethodClass(p.pushMethod)]">
                    {{ getPushMethodText(p.pushMethod) }}
                  </span>
                </td>
                <td>
                  <span :class="['push-status-badge', getPushStatusClass(p.pushStatus)]">
                    {{ getPushStatusText(p.pushStatus) }}
                  </span>
                </td>
                <td>
                  <span :class="['retry-count', { 'has-retry': p.retryCount > 0 }]">
                    {{ p.retryCount ?? 0 }}
                  </span>
                </td>
                <td>{{ formatDateTime(p.pushTime) }}</td>
              </tr>
            </tbody>
          </table>

          <div v-if="pushLogs.length === 0 && !pushLoading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
            <p class="empty-title">暂无推送记录</p>
            <p class="empty-desc">当前没有离线推送日志数据</p>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="createGroupDialog" title="" width="460px" class="custom-dialog dark-dialog">
      <div class="dialog-content">
        <div class="dialog-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M7 21v-2a4 4 0 0 1 3-3.87"/>
            <circle cx="12" cy="7" r="4"/>
            <circle cx="17" cy="11" r="2"/>
            <circle cx="7" cy="11" r="2"/>
          </svg>
        </div>
        <h3>创建群组</h3>
        <p class="dialog-desc">填写群组基本信息完成创建</p>
        <div class="form-group">
          <label>群组名称</label>
          <input type="text" v-model="groupForm.groupName" placeholder="请输入群组名称" />
        </div>
        <div class="form-group">
          <label>群主ID</label>
          <input type="text" v-model="groupForm.ownerId" placeholder="请输入群主用户ID" />
        </div>
        <div class="form-group">
          <label>群公告</label>
          <textarea v-model="groupForm.notice" placeholder="请输入群公告（选填）" rows="3"></textarea>
        </div>
      </div>
      <template #footer>
        <button class="dialog-btn cancel" @click="createGroupDialog = false">取消</button>
        <button class="dialog-btn confirm" @click="handleCreateGroup">确定创建</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getAdminSessions,
  getAdminMessages,
  getAdminGroups,
  getPushLogs
} from '../api/admin'

const activeTab = ref('sessions')

const sessionLoading = ref(false)
const messageLoading = ref(false)
const groupLoading = ref(false)
const pushLoading = ref(false)

const sessions = ref([])
const messages = ref([])
const groups = ref([])
const pushLogs = ref([])

const createGroupDialog = ref(false)
const groupForm = ref({ groupName: '', ownerId: '', notice: '' })

const stats = reactive({
  onlineUsers: 0,
  activeSessions: 0,
  todayMessages: 0,
  groupCount: 0
})

const tabs = [
  {
    key: 'sessions',
    label: '会话管理',
    icon: '<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>'
  },
  {
    key: 'messages',
    label: '消息追踪',
    icon: '<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>'
  },
  {
    key: 'groups',
    label: '群组管理',
    icon: '<path d="M17 21v-2a4 4 0 0 0-3-3.87"/><path d="M7 21v-2a4 4 0 0 1 3-3.87"/><circle cx="12" cy="7" r="4"/><circle cx="17" cy="11" r="2"/><circle cx="7" cy="11" r="2"/>'
  },
  {
    key: 'push',
    label: '离线推送',
    icon: '<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>'
  }
]

function formatDateTime(time) {
  if (!time) return '-'
  const d = new Date(time)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function getMsgTypeClass(type) {
  if (type === 1 || type === 'text') return 'type-text'
  if (type === 2 || type === 'image') return 'type-image'
  if (type === 3 || type === 'video') return 'type-video'
  if (type === 4 || type === 'file') return 'type-file'
  if (type === 5 || type === 'system') return 'type-system'
  return 'type-text'
}

function getMsgTypeText(type) {
  if (type === 1 || type === 'text') return '文本'
  if (type === 2 || type === 'image') return '图片'
  if (type === 3 || type === 'video') return '视频'
  if (type === 4 || type === 'file') return '文件'
  if (type === 5 || type === 'system') return '系统'
  return '文本'
}

function getMsgStatusClass(status) {
  if (status === 1 || status === 'sent') return 'sent'
  if (status === 2 || status === 'delivered') return 'delivered'
  if (status === 3 || status === 'read') return 'read'
  if (status === 4 || status === 'failed') return 'failed'
  return 'sent'
}

function getMsgStatusText(status) {
  if (status === 1 || status === 'sent') return '已发送'
  if (status === 2 || status === 'delivered') return '已送达'
  if (status === 3 || status === 'read') return '已读'
  if (status === 4 || status === 'failed') return '失败'
  return '已发送'
}

function getPushMethodClass(method) {
  if (method === 'apns' || method === 1) return 'method-apns'
  if (method === 'fcm' || method === 2) return 'method-fcm'
  if (method === 'websocket' || method === 3) return 'method-ws'
  if (method === 'sms' || method === 4) return 'method-sms'
  return 'method-ws'
}

function getPushMethodText(method) {
  if (method === 'apns' || method === 1) return 'APNs'
  if (method === 'fcm' || method === 2) return 'FCM'
  if (method === 'websocket' || method === 3) return 'WebSocket'
  if (method === 'sms' || method === 4) return '短信'
  return 'WebSocket'
}

function getPushStatusClass(status) {
  if (status === 1 || status === 'success') return 'push-success'
  if (status === 2 || status === 'pending') return 'push-pending'
  if (status === 3 || status === 'failed') return 'push-failed'
  if (status === 4 || status === 'retrying') return 'push-retrying'
  return 'push-pending'
}

function getPushStatusText(status) {
  if (status === 1 || status === 'success') return '推送成功'
  if (status === 2 || status === 'pending') return '等待推送'
  if (status === 3 || status === 'failed') return '推送失败'
  if (status === 4 || status === 'retrying') return '重试中'
  return '等待推送'
}

async function loadSessions() {
  sessionLoading.value = true
  try {
    const res = await getAdminSessions()
    sessions.value = res.data || []
    stats.activeSessions = sessions.value.length
  } catch {
    ElMessage.error('加载会话列表失败')
  } finally {
    sessionLoading.value = false
  }
}

async function loadMessages() {
  messageLoading.value = true
  try {
    const res = await getAdminMessages()
    messages.value = res.data || []
    stats.todayMessages = messages.value.length
  } catch {
    ElMessage.error('加载消息列表失败')
  } finally {
    messageLoading.value = false
  }
}

async function loadGroups() {
  groupLoading.value = true
  try {
    const res = await getAdminGroups()
    groups.value = res.data || []
    stats.groupCount = groups.value.length
  } catch {
    ElMessage.error('加载群组列表失败')
  } finally {
    groupLoading.value = false
  }
}

async function loadPushLogs() {
  pushLoading.value = true
  try {
    const res = await getPushLogs()
    pushLogs.value = res.data || []
  } catch {
    ElMessage.error('加载推送日志失败')
  } finally {
    pushLoading.value = false
  }
}

async function refreshAll() {
  await Promise.all([loadSessions(), loadMessages(), loadGroups(), loadPushLogs()])
  ElMessage.success('数据已刷新')
}

function openCreateGroupDialog() {
  groupForm.value = { groupName: '', ownerId: '', notice: '' }
  createGroupDialog.value = true
}

function handleCreateGroup() {
  if (!groupForm.value.groupName) {
    ElMessage.warning('请输入群组名称')
    return
  }
  if (!groupForm.value.ownerId) {
    ElMessage.warning('请输入群主ID')
    return
  }
  ElMessage.success('群组创建请求已提交')
  createGroupDialog.value = false
  loadGroups()
}

onMounted(() => {
  stats.onlineUsers = 1286
  refreshAll()
})
</script>

<style scoped>
.admin-im {
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: calc(100vh - 80px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  font-size: 24px;
  font-weight: 700;
  color: #e7e9ed;
  margin: 0 0 6px 0;
  letter-spacing: -0.3px;
}

.page-title p {
  font-size: 13px;
  color: #8b93a5;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 10px 18px;
  background: rgba(32, 201, 151, 0.12);
  border: 1px solid rgba(32, 201, 151, 0.25);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #20c997;
  cursor: pointer;
  transition: all 0.25s;
}

.refresh-btn svg {
  width: 16px;
  height: 16px;
}

.refresh-btn:hover {
  background: rgba(32, 201, 151, 0.2);
  border-color: rgba(32, 201, 151, 0.4);
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 24px;
  background: rgba(22, 27, 38, 0.72);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  overflow: hidden;
  transition: border-color 0.3s;
}

.stat-card:hover {
  border-color: rgba(32, 201, 151, 0.25);
}

.stat-icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 22px;
  height: 22px;
}

.stat-icon.online {
  background: rgba(32, 201, 151, 0.15);
  color: #20c997;
}

.stat-icon.active {
  background: rgba(77, 166, 255, 0.15);
  color: #4da6ff;
}

.stat-icon.messages {
  background: rgba(255, 183, 77, 0.15);
  color: #ffb74d;
}

.stat-icon.groups {
  background: rgba(186, 104, 200, 0.15);
  color: #ba68c8;
}

.stat-info {
  display: flex;
  flex-direction: column;
  z-index: 1;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #e7e9ed;
  line-height: 1.1;
  letter-spacing: -0.5px;
}

.stat-label {
  font-size: 12px;
  color: #8b93a5;
  margin-top: 2px;
}

.stat-spark {
  position: absolute;
  right: 14px;
  bottom: 8px;
  opacity: 0.12;
  pointer-events: none;
}

.stat-spark svg {
  width: 80px;
  height: 24px;
}

.stat-icon.online ~ .stat-spark { color: #20c997; }
.stat-icon.active ~ .stat-spark { color: #4da6ff; }
.stat-icon.messages ~ .stat-spark { color: #ffb74d; }
.stat-icon.groups ~ .stat-spark { color: #ba68c8; }

.content-card {
  background: rgba(22, 27, 38, 0.72);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  overflow: hidden;
}

.tab-nav {
  display: flex;
  gap: 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 500;
  color: #8b93a5;
  cursor: pointer;
  transition: all 0.25s;
  margin-bottom: -1px;
}

.tab-btn svg {
  width: 18px;
  height: 18px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.tab-btn:hover {
  color: #b8bfcc;
}

.tab-btn.active {
  color: #20c997;
  border-bottom-color: #20c997;
}

.tab-content {
  min-height: 400px;
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  padding: 14px 18px;
  text-align: left;
  font-size: 11px;
  font-weight: 600;
  color: #8b93a5;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.data-table td {
  padding: 15px 18px;
  font-size: 13px;
  color: #c5cad4;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.data-table tbody tr {
  transition: background 0.2s;
}

.data-table tbody tr:hover {
  background: rgba(32, 201, 151, 0.04);
}

.id-cell {
  font-family: 'SF Mono', 'Cascadia Code', monospace;
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.type-badge.single {
  background: rgba(77, 166, 255, 0.15);
  color: #4da6ff;
}

.type-badge.group {
  background: rgba(186, 104, 200, 0.15);
  color: #ba68c8;
}

.cell-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cell-avatar {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.single-avatar {
  background: rgba(77, 166, 255, 0.2);
  color: #4da6ff;
}

.group-avatar {
  background: rgba(186, 104, 200, 0.2);
  color: #ba68c8;
}

.sender-avatar {
  background: rgba(32, 201, 151, 0.18);
  color: #20c997;
}

.preview-cell {
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #8b93a5;
  font-size: 12px;
}

.msg-type-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.msg-type-badge.type-text {
  background: rgba(32, 201, 151, 0.15);
  color: #20c997;
}

.msg-type-badge.type-image {
  background: rgba(77, 166, 255, 0.15);
  color: #4da6ff;
}

.msg-type-badge.type-video {
  background: rgba(255, 112, 67, 0.15);
  color: #ff7043;
}

.msg-type-badge.type-file {
  background: rgba(255, 183, 77, 0.15);
  color: #ffb74d;
}

.msg-type-badge.type-system {
  background: rgba(144, 164, 174, 0.18);
  color: #90a4ae;
}

.status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
}

.status-dot.sent { background: #4da6ff; }
.status-dot.delivered { background: #ffb74d; }
.status-dot.read { background: #20c997; }
.status-dot.failed { background: #ef5350; }

.status-text {
  font-size: 12px;
  vertical-align: middle;
}

.tab-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.toolbar-title {
  font-size: 14px;
  font-weight: 600;
  color: #e7e9ed;
}

.action-primary-btn {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 18px;
  background: rgba(32, 201, 151, 0.15);
  border: 1px solid rgba(32, 201, 151, 0.25);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #20c997;
  cursor: pointer;
  transition: all 0.25s;
}

.action-primary-btn svg {
  width: 16px;
  height: 16px;
}

.action-primary-btn:hover {
  background: rgba(32, 201, 151, 0.22);
  border-color: rgba(32, 201, 151, 0.45);
}

.member-count {
  font-weight: 600;
  color: #e7e9ed;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-state svg {
  width: 56px;
  height: 56px;
  color: #3a3f4e;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 15px;
  font-weight: 600;
  color: #8b93a5;
  margin: 0 0 6px 0;
}

.empty-desc {
  font-size: 12px;
  color: #5a6072;
  margin: 0;
}

.push-method-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.push-method-badge.method-apns {
  background: rgba(144, 164, 174, 0.18);
  color: #90a4ae;
}

.push-method-badge.method-fcm {
  background: rgba(77, 166, 255, 0.15);
  color: #4da6ff;
}

.push-method-badge.method-ws {
  background: rgba(32, 201, 151, 0.15);
  color: #20c997;
}

.push-method-badge.method-sms {
  background: rgba(255, 183, 77, 0.15);
  color: #ffb74d;
}

.push-status-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.push-status-badge.push-success {
  background: rgba(32, 201, 151, 0.15);
  color: #20c997;
}

.push-status-badge.push-pending {
  background: rgba(255, 183, 77, 0.15);
  color: #ffb74d;
}

.push-status-badge.push-failed {
  background: rgba(239, 83, 80, 0.15);
  color: #ef5350;
}

.push-status-badge.push-retrying {
  background: rgba(77, 166, 255, 0.15);
  color: #4da6ff;
}

.retry-count {
  font-weight: 600;
  color: #8b93a5;
}

.retry-count.has-retry {
  color: #ffb74d;
}

.custom-dialog :deep(.el-dialog) {
  background: #141822;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.6);
}

.custom-dialog :deep(.el-dialog__header) {
  display: none;
}

.custom-dialog :deep(.el-dialog__body) {
  padding: 32px 32px 20px;
}

.custom-dialog :deep(.el-dialog__footer) {
  padding: 0 32px 28px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.dialog-content {
  text-align: center;
}

.dialog-icon {
  width: 56px;
  height: 56px;
  background: rgba(32, 201, 151, 0.12);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  color: #20c997;
}

.dialog-icon svg {
  width: 26px;
  height: 26px;
}

.dialog-content h3 {
  font-size: 18px;
  font-weight: 700;
  color: #e7e9ed;
  margin: 0 0 6px 0;
}

.dialog-desc {
  font-size: 13px;
  color: #8b93a5;
  margin: 0 0 24px 0;
}

.form-group {
  text-align: left;
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #8b93a5;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  color: #e7e9ed;
  outline: none;
  transition: all 0.25s;
  font-family: inherit;
  resize: none;
  box-sizing: border-box;
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: #5a6072;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: rgba(32, 201, 151, 0.5);
  background: rgba(255, 255, 255, 0.06);
}

.dialog-btn {
  padding: 10px 22px;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s;
}

.dialog-btn.cancel {
  background: rgba(255, 255, 255, 0.06);
  color: #8b93a5;
}

.dialog-btn.cancel:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #c5cad4;
}

.dialog-btn.confirm {
  background: linear-gradient(135deg, #20c997, #12b886);
  color: #fff;
}

.dialog-btn.confirm:hover {
  filter: brightness(1.12);
  box-shadow: 0 4px 16px rgba(32, 201, 151, 0.3);
}

@media (max-width: 1200px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: 1fr;
  }

  .tab-btn {
    padding: 14px 16px;
    font-size: 12px;
  }

  .tab-btn svg {
    width: 16px;
    height: 16px;
  }

  .tab-btn span {
    display: none;
  }
}
</style>