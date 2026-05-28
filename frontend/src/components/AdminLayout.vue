<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <router-link to="/admin" class="sidebar-brand">
        <div class="brand-logo">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
        </div>
        <div class="brand-text">
          <span class="brand-name">Billion SecKill</span>
          <span class="brand-sub">管理控制台</span>
        </div>
      </router-link>

      <nav class="sidebar-nav">
        <div class="nav-section-title">核心模块</div>

        <router-link to="/admin" class="nav-item" exact-active-class="active">
          <div class="nav-dot" style="background:#fff;"></div>
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
          </svg>
          <span>仪表盘</span>
        </router-link>

        <router-link to="/admin/user" class="nav-item user-item">
          <div class="nav-dot user-dot"></div>
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>用户系统</span>
        </router-link>

        <router-link to="/admin/seckill" class="nav-item seckill-item">
          <div class="nav-dot seckill-dot"></div>
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
          <span>秒杀引擎</span>
        </router-link>

        <router-link to="/admin/im" class="nav-item im-item">
          <div class="nav-dot im-dot"></div>
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <span>即时通讯</span>
        </router-link>

        <div class="nav-section-title">快捷入口</div>
        <a href="/" class="nav-item">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9 22 9 12 15 12 15 22"/>
          </svg>
          <span>返回商城</span>
        </a>
      </nav>

      <div class="sidebar-footer">
        <div class="footer-avatar">A</div>
        <div class="footer-info">
          <span class="footer-name">admin</span>
          <span class="footer-role">超级管理员</span>
        </div>
        <button class="logout-btn" @click="handleLogout" title="退出登录">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
        </button>
      </div>
    </aside>

    <main class="admin-main">
      <router-view />
    </main>

    <div class="grain-overlay"></div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: #0a0c10;
}

.admin-sidebar {
  width: 240px;
  flex-shrink: 0;
  background: #13161d;
  border-right: 1px solid rgba(255,255,255,0.06);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
}

.sidebar-brand {
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  text-decoration: none;
}

.brand-logo {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: linear-gradient(135deg, #4dabf7, #f59f00);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.brand-logo svg {
  width: 20px;
  height: 20px;
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-size: 15px;
  font-weight: 700;
  color: #e7e9ed;
  letter-spacing: -0.3px;
}

.brand-sub {
  font-size: 11px;
  color: #595f70;
  margin-top: 2px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-section-title {
  font-size: 10px;
  font-weight: 700;
  color: #595f70;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  padding: 16px 14px 6px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  color: #8b93a5;
  font-size: 13.5px;
  font-weight: 500;
  text-decoration: none;
  border: 1px solid transparent;
}

.nav-item:hover {
  background: rgba(255,255,255,0.03);
  color: #e7e9ed;
}

.nav-item.active,
.nav-item.router-link-exact-active {
  background: rgba(255,255,255,0.03);
  color: #e7e9ed;
  border-color: rgba(255,255,255,0.12);
}

.nav-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  opacity: 0.6;
}

.nav-item:hover .nav-dot,
.nav-item.active .nav-dot,
.nav-item.router-link-exact-active .nav-dot {
  opacity: 1;
}

.user-dot {
  background: #4dabf7;
  box-shadow: 0 0 8px rgba(77,171,247,0.25);
}

.seckill-dot {
  background: #f59f00;
  box-shadow: 0 0 8px rgba(245,159,0,0.25);
}

.im-dot {
  background: #20c997;
  box-shadow: 0 0 8px rgba(32,201,151,0.25);
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  opacity: 0.5;
}

.nav-item:hover .nav-icon,
.nav-item.active .nav-icon,
.nav-item.router-link-exact-active .nav-icon {
  opacity: 1;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  gap: 10px;
}

.footer-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(77,171,247,0.14);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  color: #4dabf7;
}

.footer-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.footer-name {
  font-size: 12px;
  color: #e7e9ed;
  font-weight: 600;
}

.footer-role {
  font-size: 11px;
  color: #595f70;
}

.logout-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: rgba(255,255,255,0.04);
  color: #8b93a5;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.logout-btn svg {
  width: 16px;
  height: 16px;
}

.logout-btn:hover {
  background: rgba(255,107,107,0.15);
  color: #ff6b6b;
}

.admin-main {
  margin-left: 240px;
  flex: 1;
  padding: 32px 36px;
}

.grain-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 999;
  opacity: 0.03;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
}
</style>