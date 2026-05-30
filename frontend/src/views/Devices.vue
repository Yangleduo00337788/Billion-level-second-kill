<template>
  <Layout>
    <div class="devices-page">
      <div class="page-header">
        <div class="header-left">
          <h1>设备管理</h1>
          <p>管理您的登录设备和可信设备</p>
        </div>
        <button class="remove-all-btn" @click="handleRemoveAllOthers" :disabled="devices.length <= 1">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
          </svg>
          踢掉其他设备
        </button>
      </div>

      <div class="devices-list" v-loading="loading">
        <div v-for="device in devices" :key="device.deviceId" :class="['device-card', { current: device.current }]">
          <div class="device-icon">
            <svg v-if="device.deviceType === 'MOBILE'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
              <line x1="12" y1="18" x2="12.01" y2="18"/>
            </svg>
            <svg v-else-if="device.deviceType === 'WEB'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
              <line x1="8" y1="21" x2="16" y2="21"/>
              <line x1="12" y1="17" x2="12" y2="21"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="4" y="2" width="16" height="20" rx="2" ry="2"/>
              <line x1="12" y1="6" x2="12.01" y2="6"/>
              <line x1="12" y1="18" x2="12.01" y2="18"/>
            </svg>
          </div>

          <div class="device-info">
            <div class="device-name-row">
              <span class="device-name">{{ device.deviceName || getDeviceLabel(device.deviceType) }}</span>
              <span v-if="device.current" class="current-badge">
                <svg viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="12" r="4"/></svg>
                当前设备
              </span>
              <span v-if="device.trusted" class="trusted-badge">可信</span>
            </div>
            <div class="device-meta">
              <span>{{ device.ip || '未知IP' }}</span>
              <span class="meta-dot">·</span>
              <span>{{ formatDateTime(device.lastLoginTime) }}</span>
              <span class="meta-dot">·</span>
              <span>{{ getOsLabel(device.os) }}</span>
            </div>
          </div>

          <div class="device-actions">
            <button
              v-if="!device.current"
              :class="['trust-btn', { trusted: device.trusted }]"
              @click="handleToggleTrust(device)"
            >
              <svg v-if="device.trusted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                <polyline points="9 12 11 14 15 10"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
              {{ device.trusted ? '取消可信' : '设为可信' }}
            </button>
            <button
              v-if="!device.current"
              class="remove-btn"
              @click="handleRemove(device)"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              移除
            </button>
          </div>
        </div>

        <div v-if="devices.length === 0 && !loading" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
            <line x1="8" y1="21" x2="16" y2="21"/>
            <line x1="12" y1="17" x2="12" y2="21"/>
          </svg>
          <p>暂无设备记录</p>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyDevices, removeDevice, removeAllOtherDevices, setDeviceTrusted } from '../api/user'
import Layout from '../components/Layout.vue'

const loading = ref(false)
const devices = ref([])

async function loadDevices() {
  loading.value = true
  try {
    const res = await getMyDevices()
    devices.value = res.data || []
  } catch (error) {
    ElMessage.error('加载设备列表失败')
  } finally {
    loading.value = false
  }
}

function formatDateTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function getDeviceLabel(type) {
  if (type === 'MOBILE') return '手机设备'
  if (type === 'WEB') return '浏览器'
  return '其他设备'
}

function getOsLabel(os) {
  if (!os) return '未知系统'
  if (os.includes('Windows')) return 'Windows'
  if (os.includes('Mac')) return 'macOS'
  if (os.includes('Linux')) return 'Linux'
  if (os.includes('Android')) return 'Android'
  if (os.includes('iPhone') || os.includes('iPad')) return 'iOS'
  return os
}

async function handleRemove(device) {
  try {
    await ElMessageBox.confirm('确定要移除该设备吗？', '提示', { type: 'warning' })
    await removeDevice(device.deviceId)
    ElMessage.success('设备已移除')
    loadDevices()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

async function handleRemoveAllOthers() {
  try {
    await ElMessageBox.confirm('将移除除当前设备外的所有其他设备，确定继续？', '警告', { type: 'warning' })
    await removeAllOtherDevices()
    ElMessage.success('已踢掉其他设备')
    loadDevices()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

async function handleToggleTrust(device) {
  try {
    await setDeviceTrusted(device.deviceId, !device.trusted)
    ElMessage.success(device.trusted ? '已取消可信' : '已设为可信')
    loadDevices()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped>
.devices-page {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-left h1 {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 6px 0;
}

.header-left p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.remove-all-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: #FFF0E6;
  border: 1px solid #FF5000;
  border-radius: 10px;
  color: #FF5000;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.remove-all-btn svg {
  width: 18px;
  height: 18px;
}

.remove-all-btn:hover:not(:disabled) {
  background: #FF5000;
  color: #fff;
}

.remove-all-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.devices-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.device-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  transition: all 0.2s;
}

.device-card.current {
  border: 2px solid #FF5000;
  box-shadow: 0 2px 12px rgba(255,80,0,0.12);
}

.device-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #FFF0E6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FF5000;
  flex-shrink: 0;
}

.device-icon svg {
  width: 24px;
  height: 24px;
}

.device-info {
  flex: 1;
  min-width: 0;
}

.device-name-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.device-name {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.current-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  background: #E8F5E9;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
  color: #4CAF50;
}

.current-badge svg {
  width: 10px;
  height: 10px;
}

.trusted-badge {
  padding: 3px 10px;
  background: #FFF0E6;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 600;
  color: #FF5000;
}

.device-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #999;
}

.meta-dot {
  color: #ddd;
}

.device-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.trust-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #f5f5f5;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.trust-btn svg {
  width: 16px;
  height: 16px;
}

.trust-btn.trusted {
  background: #FFF0E6;
  color: #FF5000;
}

.trust-btn:hover {
  background: #f0f0f0;
}

.trust-btn.trusted:hover {
  background: #FFD8C4;
}

.remove-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #f5f5f5;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  color: #999;
  cursor: pointer;
  transition: all 0.2s;
}

.remove-btn svg {
  width: 16px;
  height: 16px;
}

.remove-btn:hover {
  background: #e0e0e0;
  color: #F44336;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  background: #fff;
  border-radius: 12px;
}

.empty-state svg {
  width: 64px;
  height: 64px;
  color: #ccc;
}

.empty-state p {
  font-size: 14px;
  color: #999;
  margin-top: 16px;
}
</style>