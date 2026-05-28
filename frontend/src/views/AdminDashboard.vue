<template>
  <div class="dashboard-page">
    <div class="dashboard-header">
      <div>
        <h2>管理仪表盘</h2>
        <span class="header-sub">Billion SecKill 系统运行概览</span>
      </div>
      <button class="refresh-btn" @click="loadDashboard">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23 4 23 10 17 10"/>
          <polyline points="1 20 1 14 7 14"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
        刷新数据
      </button>
    </div>

    <div class="module-section" v-for="mod in modules" :key="mod.key">
      <div class="section-header">
        <svg class="section-icon" :style="{ color: mod.color }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path :d="mod.icon"/>
        </svg>
        <h3>{{ mod.title }}</h3>
        <span class="section-badge" :style="{ background: mod.badgeBg, color: mod.color }">{{ mod.badge }}</span>
      </div>

      <div class="stats-grid">
        <div v-for="stat in mod.stats" :key="stat.label" class="stat-card">
          <div class="stat-icon-box" :style="{ background: mod.iconBg }">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :style="{ color: mod.color }">
              <path :d="stat.icon"/>
            </svg>
          </div>
          <div class="stat-info">
            <span class="stat-label">{{ stat.label }}</span>
            <span class="stat-value" :style="{ color: mod.color }">{{ stat.value }}</span>
            <span class="stat-sub">{{ stat.sub }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAdminDashboard } from '../api/admin'

const dashboardData = ref({ user: {}, seckill: {}, im: {} })

const modules = ref([
  {
    key: 'user',
    title: '用户系统',
    badge: 'USER SYSTEM',
    color: '#4dabf7',
    badgeBg: 'rgba(77,171,247,0.12)',
    iconBg: 'rgba(77,171,247,0.08)',
    icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M9 7a4 4 0 1 0 0-8 4 4 0 0 0 0 8zM23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75',
    stats: []
  },
  {
    key: 'seckill',
    title: '秒杀引擎',
    badge: 'SECKILL ENGINE',
    color: '#f59f00',
    badgeBg: 'rgba(245,159,0,0.12)',
    iconBg: 'rgba(245,159,0,0.08)',
    icon: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z',
    stats: []
  },
  {
    key: 'im',
    title: '即时通讯',
    badge: 'IM SYSTEM',
    color: '#20c997',
    badgeBg: 'rgba(32,201,151,0.12)',
    iconBg: 'rgba(32,201,151,0.08)',
    icon: 'M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z',
    stats: []
  }
])

const statIcons = {
  userCount: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M9 7a4 4 0 1 0 0-8 4 4 0 0 0 0 8z',
  roleCount: 'M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5',
  blacklistCount: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z',
  productCount: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10',
  activityCount: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z',
  todayOrders: 'M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2M9 5a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2M9 5h6',
  qpsLimit: 'M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5',
  onlineUsers: 'M3 3h18v18H3zM8 12h8M12 8v8',
  activeSessions: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M23 21v-2a4 4 0 0 0-3-3.87',
  todayMessages: 'M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z',
  groupCount: 'M17 20h5v-2a3 3 0 0 0-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 0 1 5.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 0 1 9.288 0M15 7a3 3 0 1 1-6 0 3 3 0 0 1 6 0z'
}

function buildStats() {
  const d = dashboardData.value
  modules.value[0].stats = [
    { label: '注册用户', value: d.user?.userCount || 0, sub: '用户总数', icon: statIcons.userCount },
    { label: '角色数量', value: d.user?.roleCount || 0, sub: '系统角色', icon: statIcons.roleCount },
    { label: '黑名单', value: d.user?.blacklistCount || 0, sub: '已封禁目标', icon: statIcons.blacklistCount },
    { label: '活跃会话', value: '0', sub: '实时在线用户', icon: statIcons.activeSessions }
  ]
  modules.value[1].stats = [
    { label: '商品总数', value: d.seckill?.productCount || 0, sub: '可秒杀商品', icon: statIcons.productCount },
    { label: '秒杀活动', value: d.seckill?.activityCount || 0, sub: '进行中/即将开始', icon: statIcons.activityCount },
    { label: '今日订单', value: d.seckill?.todayOrders || 0, sub: '秒杀成交订单', icon: statIcons.todayOrders },
    { label: 'QPS限制', value: d.seckill?.qpsLimit || '10K', sub: '令牌桶上限', icon: statIcons.qpsLimit }
  ]
  modules.value[2].stats = [
    { label: '在线用户', value: d.im?.onlineUsers || 0, sub: 'Netty WebSocket', icon: statIcons.onlineUsers },
    { label: '活跃会话', value: d.im?.activeSessions || 0, sub: '单聊+群聊', icon: statIcons.activeSessions },
    { label: '今日消息', value: d.im?.todayMessages || 0, sub: 'Protobuf编码传输', icon: statIcons.todayMessages },
    { label: '群组数量', value: d.im?.groupCount || 0, sub: '活跃群组', icon: statIcons.groupCount }
  ]
}

async function loadDashboard() {
  try {
    const res = await getAdminDashboard()
    if (res.data) {
      dashboardData.value = res.data
    }
  } catch (e) {
    console.error('Dashboard load failed', e)
  }
  buildStats()
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.dashboard-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dashboard-header h2 {
  font-size: 26px;
  font-weight: 800;
  color: #e7e9ed;
  letter-spacing: -0.5px;
  margin: 0;
}

.header-sub {
  color: #595f70;
  font-size: 13px;
  margin-top: 4px;
  display: block;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.06);
  background: rgba(255,255,255,0.03);
  color: #8b93a5;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.refresh-btn svg {
  width: 16px;
  height: 16px;
}

.refresh-btn:hover {
  background: rgba(255,255,255,0.06);
  color: #e7e9ed;
}

.module-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-icon {
  width: 22px;
  height: 22px;
}

.section-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: #e7e9ed;
  margin: 0;
}

.section-badge {
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.stat-card {
  background: rgba(22,27,38,0.72);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  transition: all 0.2s;
}

.stat-card:hover {
  background: rgba(30,37,50,0.85);
  border-color: rgba(255,255,255,0.12);
  transform: translateY(-2px);
}

.stat-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon-box svg {
  width: 22px;
  height: 22px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 11px;
  color: #595f70;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  font-family: 'DM Mono', 'Consolas', monospace;
  margin: 2px 0;
}

.stat-sub {
  font-size: 11px;
  color: #8b93a5;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>