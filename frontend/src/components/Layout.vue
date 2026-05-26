<template>
  <div class="tb-layout">
    <header class="tb-header">
      <div class="header-top">
        <div class="top-left">
          <span class="welcome">Hi，欢迎来秒杀商城</span>
          <router-link to="/login" class="top-link" v-if="!userStore.token">请登录</router-link>
          <router-link to="/register" class="top-link register" v-if="!userStore.token">免费注册</router-link>
          <span class="username" v-if="userStore.token">{{ userStore.username }}</span>
          <button class="top-link logout" v-if="userStore.token" @click="handleLogout">退出</button>
        </div>
        <div class="top-right">
          <router-link to="/orders" class="top-link">我的订单</router-link>
          <span class="divider">|</span>
          <router-link to="/cart" class="top-link cart-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-cart">
              <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
              <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
            </svg>
            购物车
            <span class="cart-count" v-if="cartCount > 0">{{ cartCount }}</span>
          </router-link>
          <span class="divider">|</span>
          <a class="top-link" @click="showCategoryPanel = true">商品分类</a>
          <span class="divider">|</span>
          <a class="top-link">联系客服</a>
        </div>
      </div>
      
      <div class="header-main">
        <div class="header-inner">
          <router-link to="/" class="logo-area">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
            </div>
            <div class="logo-text">
              <span class="logo-main">秒杀商城</span>
              <span class="logo-sub">Seckill Mall</span>
            </div>
          </router-link>
          
          <div class="search-area">
            <div class="search-box">
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
              <span>热门：</span>
              <a v-for="word in hotWords" :key="word" class="hot-word-link" @click.prevent="clickHotWord(word)">{{ word }}</a>
            </div>
          </div>
          
          <div class="header-actions">
            <router-link to="/seckill" class="action-btn seckill-btn">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
              <span>限时秒杀</span>
            </router-link>
          </div>
        </div>
      </div>
      
      <nav class="header-nav">
        <div class="nav-inner">
          <div class="nav-category" @click="showCategoryPanel = !showCategoryPanel">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
            <span>全部商品分类</span>
          </div>
          <div class="nav-links">
            <router-link to="/" class="nav-link" :class="{ active: currentRoute === '/' }">首页</router-link>
            <router-link to="/seckill" class="nav-link" :class="{ active: currentRoute === '/seckill' }">秒杀专区</router-link>
          </div>
        </div>
      </nav>
      
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
        <p class="copyright">© 2024 秒杀商城 版权所有 | 技术支持：Spring Boot 3 + Redis + Kafka</p>
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
const cartCount = ref(0)
const categories = ref([])
const products = ref([])

const hotWords = ['iPhone', 'MacBook', '耳机', '手表', '手机']

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
  background: #f5f5f5;
}

.tb-header {
  background: linear-gradient(135deg, #FF5000 0%, #FF6000 50%, #FF7000 100%);
  position: sticky;
  top: 0;
  z-index: 1000;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.header-top {
  background: rgba(0,0,0,0.1);
  border-bottom: 1px solid rgba(255,255,255,0.1);
  padding: 8px 0;
  display: flex;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 8px 20px;
}

.top-left, .top-right {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.welcome {
  color: rgba(255,255,255,0.9);
}

.top-link {
  color: rgba(255,255,255,0.9);
  cursor: pointer;
  text-decoration: none;
  transition: color 0.2s;
}

.top-link:hover {
  color: #fff;
}

.top-link.register {
  color: #fff;
  font-weight: 500;
}

.username {
  color: #fff;
  font-weight: 500;
}

.logout {
  background: transparent;
  border: none;
  padding: 0;
  font-size: 12px;
  color: rgba(255,255,255,0.9);
}

.divider {
  color: rgba(255,255,255,0.5);
  margin: 0 4px;
}

.cart-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: rgba(255,255,255,0.9);
}

.icon-cart {
  width: 14px;
  height: 14px;
}

.cart-count {
  background: #fff;
  color: #FF5000;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 10px;
}

.header-main {
  background: linear-gradient(135deg, #FF5000 0%, #FF7800 50%, #FF9000 100%);
  padding: 16px 0;
  border-bottom: none;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  flex-shrink: 0;
}

.logo-icon {
  width: 48px;
  height: 48px;
  background: rgba(255,255,255,0.95);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FF5000;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.logo-icon svg {
  width: 28px;
  height: 28px;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-main {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1px;
}

.logo-sub {
  font-size: 11px;
  color: rgba(255,255,255,0.8);
  letter-spacing: 2px;
}

.search-area {
  flex: 1;
  max-width: 560px;
  min-width: 400px;
}

.search-box {
  display: flex;
  border: 2px solid rgba(255,255,255,0.9);
  border-radius: 24px;
  overflow: hidden;
  background: rgba(255,255,255,0.95);
  box-shadow: 0 2px 12px rgba(0,0,0,0.15);
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: none;
  font-size: 14px;
  outline: none;
  background: transparent;
}

.search-input::placeholder {
  color: #999;
}

.search-btn {
  padding: 12px 28px;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border: none;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.search-btn:hover {
  background: linear-gradient(135deg, #FF3400, #FF2000);
}

.hot-words {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  font-size: 12px;
}

.hot-words span {
  color: rgba(255,255,255,0.8);
}

.hot-words .hot-word-link {
  color: rgba(255,255,255,0.9);
  cursor: pointer;
  transition: color 0.2s;
  text-decoration: none;
}

.hot-words .hot-word-link:hover {
  color: #fff;
  text-decoration: underline;
}

.header-actions {
  display: flex;
  gap: 16px;
  flex-shrink: 0;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 24px;
  border-radius: 24px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
  height: 48px;
  box-sizing: border-box;
}

.seckill-btn {
  background: rgba(255,255,255,0.95);
  color: #FF5000;
  border: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.seckill-btn svg {
  width: 18px;
  height: 18px;
}

.seckill-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.header-nav {
  background: linear-gradient(135deg, #FF5000 0%, #FF6800 50%, #FF8000 100%);
}

.nav-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
}

.nav-category {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(0,0,0,0.15);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  position: relative;
}

.nav-category svg {
  width: 18px;
  height: 18px;
}

.nav-links {
  display: flex;
  gap: 24px;
  padding: 12px 20px;
}

.nav-link {
  color: #fff;
  font-size: 14px;
  text-decoration: none;
  padding: 8px 0;
  position: relative;
  transition: color 0.2s;
}

.nav-link:hover {
  color: rgba(255,255,255,0.8);
}

.nav-link.active {
  color: #fff;
  font-weight: 600;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: #fff;
}

.category-panel {
  position: absolute;
  top: 100%;
  left: 20px;
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
  background: #FFF5F0;
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  width: 100%;
}

.tb-footer {
  background: #fff;
  margin-top: 40px;
  border-top: 1px solid #e5e5e5;
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
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