<template>
  <div class="tb-layout">
    <header class="tb-header">
      <!-- 推广横幅 -->
      <div class="promo-bar">
        <div class="promo-inner">
          <div class="promo-left">
            <div class="promo-icon">
              <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=taobao%20mascot%20cute%20cartoon%20character%20orange%20small%20icon&image_size=square" alt="" />
            </div>
            <span class="promo-text">首页专属福利限时领，好货低价别错过</span>
          </div>
          <button class="promo-btn" @click="$router.push('/seckill')">去逛逛</button>
        </div>
      </div>

      <!-- 快捷入口 -->
      <div class="quick-bar">
        <div class="quick-inner">
          <div class="quick-item" v-for="(item, index) in quickEntries" :key="index" @click="handleEntryClick(item)">
            <div class="quick-icon" :style="{ background: item.bg }">
              <svg viewBox="0 0 24 24" fill="currentColor"><path :d="item.icon"/></svg>
            </div>
            <span class="quick-name">{{ item.name }}</span>
          </div>
        </div>
      </div>

      <!-- 搜索区 -->
      <div class="header-main">
        <div class="header-inner">
          <router-link to="/" class="logo-area">
            <div class="logo-text">
              <span class="logo-main">秒杀商城</span>
              <span class="logo-sub">Taobao.com</span>
            </div>
            <div class="logo-slogan">
              <span>热卖</span>
              <span>商品</span>
            </div>
          </router-link>

          <div class="search-area">
            <div class="search-box">
              <div class="search-type">
                <span>宝贝</span>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </div>
              <input
                type="text"
                v-model="searchKeyword"
                placeholder="搜索秒杀商品..."
                class="search-input"
                @keyup.enter="handleSearch"
              />
              <button class="search-btn" type="button" @click.prevent="handleSearch">搜索</button>
            </div>
            <div class="hot-words">
              <a v-for="word in hotWords" :key="word" class="hot-word-link" @click.prevent="clickHotWord(word)">{{ word }}</a>
            </div>
          </div>

          <div class="header-tools">
            <template v-if="userStore.isLoggedIn()">
              <router-link to="/im" class="tool-link im-link" :class="{ active: currentRoute === '/im' }">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-chat">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                </svg>
                消息
              </router-link>
              <router-link to="/orders" class="tool-link">我的订单</router-link>
              <router-link to="/cart" class="tool-link cart-link">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-cart">
                  <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
                  <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
                </svg>
                购物车
                <span class="cart-count" v-if="cartCount > 0">{{ cartCount }}</span>
              </router-link>
              <div class="tool-user">
                <button class="tool-link tool-user-btn" @click="showUserMenu = !showUserMenu">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-user">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                  {{ userStore.username || '个人中心' }}
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="admin-arrow" :class="{ open: showUserMenu }"><polyline points="6 9 12 15 18 9"/></svg>
                </button>
                <div class="admin-dropdown" v-show="showUserMenu" @mouseleave="showUserMenu = false">
                  <router-link to="/profile" class="admin-drop-item" :class="{ active: currentRoute === '/profile' }">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                      <circle cx="12" cy="7" r="4"/>
                    </svg>
                    个人资料
                  </router-link>
                  <a class="admin-drop-item logout-item" @click="handleLogout">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                      <polyline points="16 17 21 12 16 7"/>
                      <line x1="21" y1="12" x2="9" y2="12"/>
                    </svg>
                    退出登录
                  </a>
                </div>
              </div>
            </template>
            <template v-else>
              <router-link to="/login" class="tool-link login-link">登录</router-link>
              <router-link to="/register" class="tool-link register-link-btn">免费注册</router-link>
            </template>
          </div>
        </div>
      </div>

      <!-- 导航栏 -->
      <nav class="header-nav" ref="navRef">
        <div class="nav-inner">
          <div class="nav-category" :class="{ active: showCategoryPanel }" @click="showCategoryPanel = !showCategoryPanel">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="nav-cat-icon">
              <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
            <span>主题市场</span>
          </div>
          <div class="nav-links">
            <router-link to="/" class="nav-link" :class="{ active: currentRoute === '/' }">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>
              首页
            </router-link>
            <router-link to="/seckill" class="nav-link" :class="{ active: currentRoute === '/seckill' }">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
              秒杀专区
            </router-link>
            
          </div>
        </div>

        <div class="category-panel" v-show="showCategoryPanel" @mouseleave="showCategoryPanel = false">
          <div class="category-list">
            <div
              v-for="cat in categories"
              :key="cat.id"
              class="category-item"
              @click="goCategory(cat.id)"
            >
              <div class="cat-icon" :style="{ background: cat.color }">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path :d="cat.icon"/>
                </svg>
              </div>
              <span class="cat-name">{{ cat.name }}</span>
            </div>
          </div>
        </div>
      </nav>
    </header>

    <main class="tb-main">
      <slot></slot>
    </main>

    <footer class="tb-footer">
      <div class="footer-inner">
        <div class="footer-services">
          <div class="service-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <span>正品保障</span>
          </div>
          <div class="service-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/>
              <circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/>
            </svg>
            <span>极速发货</span>
          </div>
          <div class="service-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
              <polyline points="22,6 12,13 2,6"/>
            </svg>
            <span>七天退换</span>
          </div>
          <div class="service-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72"/>
            </svg>
            <span>24h客服</span>
          </div>
        </div>
        <div class="footer-links">
          <a>关于我们</a>
          <span class="dot">·</span>
          <a>联系客服</a>
          <span class="dot">·</span>
          <a>隐私政策</a>
          <span class="dot">·</span>
          <a>用户协议</a>
        </div>
        <p class="copyright"> 2024 秒杀商城 版权所有 | 技术支持：Spring Boot 3 + Redis + Kafka</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { getSeckillProducts } from '../api/seckill'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const searchKeyword = ref('')
const showCategoryPanel = ref(false)
const showUserMenu = ref(false)
const cartCount = ref(0)
const categories = ref([])
const products = ref([])
const navRef = ref(null)

const hotWords = ['新款手机', '电脑', '耳机', '手表', '运动鞋']

const quickEntries = [
  { name: '国家补贴', bg: '#52C41A', icon: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z', route: '/seckill' },
  { name: '淘宝秒杀', bg: '#FF4D4F', icon: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z', route: '/seckill' },
  { name: 'U先试用', bg: '#FAAD14', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', route: '/seckill' },
  { name: '百亿补贴', bg: '#FF5000', icon: 'M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z', route: '/seckill' },
  { name: '领券中心', bg: '#EB2F96', icon: 'M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z', route: '/seckill' },
  { name: '聚划算', bg: '#FF5000', icon: 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2', route: '/seckill' },
]

const currentRoute = computed(() => route.path)

async function loadCategories() {
  try {
    const res = await getSeckillProducts()
    products.value = res.data || []
    const categoryMap = new Map()
    products.value.forEach(p => {
      if (p.categoryId && p.categoryName) {
        categoryMap.set(p.categoryId, p.categoryName)
      }
    })
    const colors = ['#FF5000', '#1677FF', '#52C41A', '#FAAD14', '#EB2F96', '#13C2C2', '#722ED1', '#595959']
    const icons = [
      'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z',
      'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z',
      'M13 10V3L4 14h7v7l9-11h-7z',
      'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z',
      'M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z',
      'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z'
    ]
    let idx = 0
    categoryMap.forEach((name, id) => {
      categories.value.push({
        id,
        name,
        color: colors[idx % colors.length],
        icon: icons[idx % icons.length]
      })
      idx++
    })
    if (categories.value.length === 0) {
      categories.value = [
        { id: 1, name: '手机数码', color: '#FF5000', icon: icons[0] },
        { id: 2, name: '电脑办公', color: '#1677FF', icon: icons[1] },
        { id: 3, name: '家用电器', color: '#52C41A', icon: icons[2] },
        { id: 4, name: '服装鞋包', color: '#FAAD14', icon: icons[3] },
      ]
    }
  } catch (e) {
    categories.value = [
      { id: 1, name: '手机数码', color: '#FF5000', icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z' },
      { id: 2, name: '电脑办公', color: '#1677FF', icon: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' },
    ]
  }
}

function updateCartCount() {
  const cart = JSON.parse(localStorage.getItem('cart') || '[]')
  cartCount.value = cart.length
}

function handleSearch() {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/seckill', query: { keyword: searchKeyword.value.trim() } })
  }
}

function clickHotWord(word) {
  searchKeyword.value = word
  handleSearch()
}

function goCategory(catId) {
  showCategoryPanel.value = false
  router.push({ path: '/seckill', query: { category: catId } })
}

function handleEntryClick(item) {
  if (item.route) {
    router.push(item.route)
  }
}

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadCategories()
  updateCartCount()
})
</script>

<style scoped>
.tb-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f4f4f4;
}

/* ===== 推广横幅 ===== */
.promo-bar {
  background: linear-gradient(90deg, #FF5000, #FF7A33);
  padding: 8px 0;
}

.promo-inner {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.promo-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.promo-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.promo-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.promo-text {
  color: #fff;
  font-size: 15px;
  font-weight: 500;
}

.promo-btn {
  padding: 6px 18px;
  background: #fff;
  border: none;
  border-radius: 16px;
  color: #FF5000;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

/* ===== 快捷入口 ===== */
.quick-bar {
  background: #fff;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.quick-inner {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 60px;
  display: flex;
  justify-content: center;
  gap: 48px;
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: transform 0.2s;
}

.quick-item:hover {
  transform: translateY(-2px);
}

.quick-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.quick-icon svg {
  width: 24px;
  height: 24px;
}

.quick-name {
  font-size: 12px;
  color: #666;
}

/* ===== 搜索区 ===== */
.header-main {
  background: #fff;
  padding: 20px 0 0;
}

.header-inner {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 60px;
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  flex-shrink: 0;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-main {
  font-size: 32px;
  font-weight: 700;
  color: #FF5000;
  letter-spacing: 2px;
  line-height: 1;
}

.logo-sub {
  font-size: 11px;
  color: #FF5000;
  letter-spacing: 1px;
  margin-top: 2px;
}

.logo-slogan {
  display: flex;
  flex-direction: column;
  border-left: 2px solid #FF5000;
  padding-left: 8px;
  margin-left: 4px;
}

.logo-slogan span {
  font-size: 14px;
  color: #FF5000;
  font-weight: 600;
  line-height: 1.3;
}

.search-area {
  flex: 1;
  max-width: none;
  min-width: 0;
}

.search-box {
  display: flex;
  border: 2px solid #FF5000;
  border-radius: 20px;
  overflow: hidden;
  background: #fff;
}

.search-type {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 12px;
  background: #f5f5f5;
  border-right: 1px solid #eee;
  font-size: 12px;
  color: #666;
  cursor: pointer;
  white-space: nowrap;
}

.search-type svg {
  width: 12px;
  height: 12px;
}

.search-input {
  flex: 1;
  padding: 10px 16px;
  border: none;
  font-size: 14px;
  outline: none;
  background: transparent;
}

.search-input::placeholder {
  color: #ccc;
}

.search-btn {
  padding: 10px 32px;
  background: linear-gradient(90deg, #FF7A33, #FF5000);
  border: none;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 0 18px 18px 0;
}

.search-btn:hover {
  background: linear-gradient(90deg, #FF6611, #FF3300);
}

.hot-words {
  display: flex;
  gap: 16px;
  margin-top: 6px;
  font-size: 12px;
  padding-left: 4px;
}

.hot-words .hot-word-link {
  color: #999;
  cursor: pointer;
  transition: color 0.2s;
  text-decoration: none;
}

.hot-words .hot-word-link:hover {
  color: #FF5000;
}

/* 右侧工具 */
.header-tools {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.tool-link {
  color: #666;
  font-size: 13px;
  text-decoration: none;
  transition: color 0.2s;
}

.tool-link:hover {
  color: #FF5000;
}

.cart-link {
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-cart {
  width: 16px;
  height: 16px;
}

.im-link {
  display: flex;
  align-items: center;
  gap: 4px;
}

.im-link.active {
  color: #FF5000;
}

.icon-chat {
  width: 16px;
  height: 16px;
}

.tool-user-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: none;
  cursor: pointer;
  font-family: inherit;
}

.tool-user-btn.active {
  color: #FF5000;
}

.icon-user {
  width: 16px;
  height: 16px;
}

.tool-user {
  position: relative;
}

.tool-user .admin-dropdown {
  right: 0;
  left: auto;
}

.dropdown-divider {
  height: 1px;
  background: #eee;
  margin: 6px 0;
}

.logout-item {
  color: #999;
}

.logout-item:hover {
  background: #FFF0E6;
  color: #FF5000;
}

.login-link {
  font-weight: 600;
  color: #FF5000;
}

.login-link:hover {
  color: #FF3300;
}

.register-link-btn {
  padding: 6px 16px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 16px;
  color: #fff !important;
  font-weight: 600;
}

.register-link-btn:hover {
  filter: brightness(1.1);
  color: #fff !important;
}

.cart-count {
  background: #FF5000;
  color: #fff;
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 8px;
}

/* ===== 导航栏 ===== */
.header-nav {
  background: linear-gradient(90deg, #FF7A33, #FF5000);
  margin-top: 12px;
  position: relative;
}

.nav-inner {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 60px;
  display: flex;
  align-items: center;
}

.nav-category {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 24px;
  background: rgba(0,0,0,0.15);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.nav-category:hover,
.nav-category.active {
  background: rgba(0,0,0,0.25);
}

.nav-cat-icon {
  width: 16px;
  height: 16px;
}

.nav-links {
  display: flex;
  gap: 0;
  padding: 0;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255,255,255,0.85);
  font-size: 14px;
  text-decoration: none;
  padding: 10px 24px;
  transition: all 0.2s;
  position: relative;
}

.nav-link svg {
  width: 16px;
  height: 16px;
}

.nav-link:hover {
  color: #fff;
  background: rgba(255,255,255,0.1);
}

.nav-link.active {
  color: #fff;
  font-weight: 600;
  background: rgba(255,255,255,0.2);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 3px;
  background: #fff;
  border-radius: 2px;
}



.admin-arrow {
  width: 14px;
  height: 14px;
  margin-left: 2px;
  transition: transform 0.2s;
}

.admin-arrow.open {
  transform: rotate(180deg);
}

.admin-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  min-width: 140px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 8px 28px rgba(0,0,0,0.12);
  padding: 8px;
  z-index: 200;
  animation: fadeInDown 0.15s ease;
}

@keyframes fadeInDown {
  from { opacity: 0; transform: translateY(-6px); }
  to { opacity: 1; transform: translateY(0); }
}

.admin-drop-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  color: #333;
  text-decoration: none;
  transition: all 0.15s;
}

.admin-drop-item svg {
  width: 18px;
  height: 18px;
  color: #999;
}

.admin-drop-item:hover {
  background: #FFF0E6;
  color: #FF5000;
}

.admin-drop-item:hover svg {
  color: #FF5000;
}

.admin-drop-item.active {
  background: #FFF0E6;
  color: #FF5000;
  font-weight: 500;
}

.admin-drop-item.active svg {
  color: #FF5000;
}

.category-panel {
  position: absolute;
  top: 100%;
  left: 60px;
  width: 200px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  z-index: 1001;
  padding: 12px;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  background: #FFF0E6;
}

.cat-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.cat-icon svg {
  width: 18px;
  height: 18px;
}

.cat-name {
  font-size: 14px;
  color: #333;
}

.tb-main {
  flex: 1;
  max-width: 100%;
  margin: 0 auto;
  padding: 0 60px;
  width: 100%;
}

.tb-footer {
  background: #fff;
  margin-top: 40px;
  border-top: 1px solid #e5e5e5;
}

.footer-inner {
  max-width: 100%;
  margin: 0 auto;
  padding: 40px 60px;
}

.footer-services {
  display: flex;
  justify-content: center;
  gap: 60px;
  margin-bottom: 30px;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #666;
}

.service-item svg {
  width: 32px;
  height: 32px;
  color: #FF5000;
}

.service-item span {
  font-size: 14px;
  font-weight: 500;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 16px;
}

.footer-links a {
  font-size: 12px;
  color: #666;
  cursor: pointer;
}

.footer-links .dot {
  color: #ccc;
}

.copyright {
  text-align: center;
  font-size: 12px;
  color: #999;
}
</style>
