<template>
  <div class="tb-pc-page">
    <!-- 顶部导航栏 -->
    <header class="pc-topbar">
      <div class="topbar-inner">
        <div class="topbar-left">
          <span class="welcome">Hi，欢迎来秒杀商城</span>
          <router-link to="/login" class="topbar-link" v-if="!userStore.token">请登录</router-link>
          <router-link to="/register" class="topbar-link" v-if="!userStore.token">免费注册</router-link>
          <span class="username" v-if="userStore.token">{{ userStore.username }}</span>
          <button class="logout-link" v-if="userStore.token" @click="handleLogout">退出</button>
        </div>
        <div class="topbar-right">
          <a class="topbar-link" @click="showOrders = true">我的订单</a>
          <a class="topbar-link">购物车</a>
          <a class="topbar-link">收藏夹</a>
          <a class="topbar-link">商品分类</a>
          <a class="topbar-link">卖家中心</a>
          <a class="topbar-link">联系客服</a>
        </div>
      </div>
    </header>

    <!-- 搜索区域 -->
    <div class="pc-search-area">
      <div class="search-inner">
        <div class="logo-area">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
            </svg>
          </div>
          <div class="logo-text">
            <span class="logo-main">秒杀商城</span>
            <span class="logo-sub">Seckill Mall</span>
          </div>
        </div>
        <div class="search-box">
          <div class="search-input-wrap">
            <input 
              type="text" 
              v-model="searchKeyword" 
              placeholder="搜索秒杀商品..."
              class="search-input"
              @keyup.enter="handleSearch"
            />
            <button class="search-camera">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/>
                <circle cx="12" cy="13" r="4"/>
              </svg>
            </button>
          </div>
          <button class="search-submit" @click="handleSearch">搜索</button>
          <div class="hot-words">
            <span>热门搜索：</span>
            <a v-for="word in hotWords" :key="word" @click="searchKeyword = word">{{ word }}</a>
          </div>
        </div>
        <div class="qr-code">
          <div class="qr-box">
            <span>手机扫码</span>
            <div class="qr-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7"/>
                <rect x="14" y="3" width="7" height="7"/>
                <rect x="14" y="14" width="7" height="7"/>
                <rect x="3" y="14" width="7" height="7"/>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 主导航 -->
    <nav class="pc-main-nav">
      <div class="nav-inner">
        <div class="nav-category">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="12" x2="21" y2="12"/>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
          <span>主题市场</span>
        </div>
        <div class="nav-links">
          <a class="nav-link active">秒杀首页</a>
          <a class="nav-link">天猫</a>
          <a class="nav-link">聚划算</a>
          <a class="nav-link">超市</a>
          <a class="nav-link">拍卖</a>
          <a class="nav-link">飞猪旅行</a>
          <a class="nav-link">天天特卖</a>
        </div>
      </div>
    </nav>

    <!-- 主体内容 -->
    <div class="pc-main-content">
      <div class="content-inner">
        <!-- 左侧分类 -->
        <aside class="category-sidebar">
          <div class="category-title">商品分类</div>
          <div class="category-list">
            <div 
              v-for="cat in categories" 
              :key="cat.id"
              class="category-item"
              :class="{ active: selectedCategory === cat.id }"
              @click="selectedCategory = cat.id"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path :d="cat.icon"/>
              </svg>
              <span>{{ cat.name }}</span>
            </div>
          </div>
        </aside>

        <!-- 中间内容区 -->
        <div class="center-content">
          <!-- Banner区 -->
          <div class="banner-section">
            <div class="main-banner">
              <div class="banner-slide active">
                <div class="banner-content-pc">
                  <div class="banner-text-pc">
                    <h2>限时秒杀</h2>
                    <p>全场低至1折起</p>
                    <div class="banner-countdown">
                      <span>距结束</span>
                      <div class="countdown-box">
                        <span class="cd-num">{{ countdownHours }}</span>
                        <span class="cd-sep">:</span>
                        <span class="cd-num">{{ countdownMinutes }}</span>
                        <span class="cd-sep">:</span>
                        <span class="cd-num">{{ countdownSeconds }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="side-banners">
              <div class="side-banner">
                <span>我的订单</span>
                <p>查看秒杀订单</p>
              </div>
              <div class="side-banner">
                <span>秒杀提醒</span>
                <p>设置开抢提醒</p>
              </div>
            </div>
          </div>

          <!-- 筛选栏 -->
          <div class="filter-bar">
            <div class="filter-left">
              <button 
                class="filter-btn" 
                :class="{ active: filter === 'all' }"
                @click="filter = 'all'"
              >
                全部商品
              </button>
              <button 
                class="filter-btn" 
                :class="{ active: filter === 'active' }"
                @click="filter = 'active'"
              >
                正在秒杀
              </button>
              <button 
                class="filter-btn" 
                :class="{ active: filter === 'pending' }"
                @click="filter = 'pending'"
              >
                即将开始
              </button>
              <button 
                class="filter-btn" 
                :class="{ active: filter === 'ended' }"
                @click="filter = 'ended'"
              >
                已结束
              </button>
            </div>
            <div class="filter-right">
              <span class="result-count">共 {{ filteredProducts.length }} 件商品</span>
              <div class="sort-options">
                <button 
                  class="sort-btn" 
                  :class="{ active: sortBy === 'default' }"
                  @click="sortBy = 'default'"
                >
                  综合
                </button>
                <button 
                  class="sort-btn" 
                  :class="{ active: sortBy === 'price' }"
                  @click="sortBy = 'price'"
                >
                  价格
                  <svg class="sort-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="6 9 12 15 18 9"/>
                  </svg>
                </button>
                <button 
                  class="sort-btn" 
                  :class="{ active: sortBy === 'stock' }"
                  @click="sortBy = 'stock'"
                >
                  库存
                </button>
              </div>
            </div>
          </div>

          <!-- 商品网格 -->
          <div class="products-grid" v-loading="loading">
            <div 
              v-for="(product, index) in sortedProducts" 
              :key="product.id"
              class="grid-item"
              :style="{ animationDelay: `${index * 0.03}s` }"
            >
              <div class="item-image">
                <div class="image-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
                <div class="item-tags">
                  <span class="tag-seckill" v-if="product.status === 1">秒杀</span>
                  <span class="tag-discount" v-if="product.originalPrice">
                    {{ Math.round((1 - product.seckillPrice / product.originalPrice) * 100) }}折
                  </span>
                </div>
                <div class="item-overlay" v-if="product.status !== 1 || product.seckillStock <= 0">
                  <span v-if="product.seckillStock <= 0">已售罄</span>
                  <span v-else-if="product.status === 2">已结束</span>
                  <span v-else-if="product.status === 0">即将开始</span>
                </div>
              </div>
              <div class="item-info">
                <h3 class="item-title">{{ product.productName || `秒杀商品 #${product.id}` }}</h3>
                <div class="item-price-row">
                  <div class="item-price">
                    <span class="price-symbol">¥</span>
                    <span class="price-value">{{ product.seckillPrice }}</span>
                  </div>
                  <div class="item-original" v-if="product.originalPrice">
                    <span>¥{{ product.originalPrice }}</span>
                  </div>
                </div>
                <div class="item-sales">
                  <span class="sales-count">已抢{{ Math.round(100 - getStockPercent(product)) }}%</span>
                  <div class="sales-bar">
                    <div class="sales-fill" :style="{ width: `${100 - getStockPercent(product)}%` }"></div>
                  </div>
                </div>
                <div class="item-meta">
                  <span class="meta-stock">库存 {{ product.seckillStock }}</span>
                  <span class="meta-time">{{ formatTime(product.startTime) }}-{{ formatTime(product.endTime) }}</span>
                </div>
                <button 
                  class="item-buy-btn"
                  :class="{ disabled: product.status !== 1 || product.seckillStock <= 0 }"
                  :disabled="product.status !== 1 || product.seckillStock <= 0"
                  @click="handleSeckill(product)"
                >
                  <span v-if="product.seckillStock <= 0">已售罄</span>
                  <span v-else-if="product.status === 2">已结束</span>
                  <span v-else-if="product.status === 0">即将开始</span>
                  <span v-else>立即抢购</span>
                </button>
              </div>
            </div>

            <!-- 空状态 -->
            <div v-if="!loading && sortedProducts.length === 0" class="grid-empty">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无秒杀商品</p>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination-bar" v-if="sortedProducts.length > 0">
            <button class="page-btn" :disabled="currentPage === 1" @click="currentPage--">上一页</button>
            <div class="page-numbers">
              <button 
                v-for="page in totalPages" 
                :key="page"
                class="page-num"
                :class="{ active: currentPage === page }"
                @click="currentPage = page"
              >
                {{ page }}
              </button>
            </div>
            <button class="page-btn" :disabled="currentPage === totalPages" @click="currentPage++">下一页</button>
            <span class="page-info">共 {{ totalPages }} 页</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部 -->
    <footer class="pc-footer">
      <div class="footer-inner">
        <div class="footer-links">
          <a>关于我们</a>
          <a>联系客服</a>
          <a>合作伙伴</a>
          <a>营销中心</a>
          <a>隐私政策</a>
          <a>用户协议</a>
        </div>
        <p class="copyright">© 2024 秒杀商城 版权所有</p>
      </div>
    </footer>

    <!-- 秒杀结果弹窗 -->
    <el-dialog 
      v-model="seckillDialogVisible" 
      title=""
      width="400px"
      :close-on-click-modal="false"
      class="pc-dialog"
    >
      <div class="result-panel">
        <div class="result-icon" :class="seckillResult.success ? 'success' : 'fail'">
          <svg v-if="seckillResult.success" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"/>
          </svg>
        </div>
        <h3 class="result-title">{{ seckillResult.success ? '恭喜！秒杀成功' : '秒杀失败' }}</h3>
        <p class="result-desc">{{ seckillResult.message }}</p>
        <div v-if="seckillResult.success" class="result-order-info">
          <span class="order-label">订单号</span>
          <span class="order-value">{{ seckillResult.orderNo }}</span>
        </div>
      </div>
      <template #footer>
        <button class="pc-dialog-btn primary" @click="seckillDialogVisible = false">
          {{ seckillResult.success ? '立即支付' : '继续抢购' }}
        </button>
      </template>
    </el-dialog>

    <!-- 订单抽屉 -->
    <el-drawer
      v-model="showOrders"
      title="我的订单"
      size="600px"
      class="orders-drawer"
    >
      <div class="orders-content">
        <div class="orders-tabs">
          <button class="orders-tab active">全部订单</button>
          <button class="orders-tab">待付款</button>
          <button class="orders-tab">已完成</button>
          <button class="orders-tab">已取消</button>
        </div>
        <div v-loading="ordersLoading" class="orders-list">
          <div v-for="order in orders" :key="order.id" class="order-card">
            <div class="order-header">
              <span class="order-date">{{ formatDateTime(order.createTime) }}</span>
              <span class="order-no">订单号：{{ order.orderNo }}</span>
              <span class="order-status-tag" :class="getStatusClassByOrder(order)">
                {{ getStatusTextByOrder(order) }}
              </span>
            </div>
            <div class="order-body">
              <div class="order-product">
                <div class="product-image">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  </svg>
                </div>
                <div class="product-detail">
                  <h4>{{ order.productName }}</h4>
                  <p>秒杀商品</p>
                </div>
                <div class="product-price-qty">
                  <span class="price">¥{{ order.totalAmount }}</span>
                  <span class="qty">x1</span>
                </div>
              </div>
            </div>
            <div class="order-footer-pc">
              <div class="order-total-pc">
                共1件商品 实付款：<span class="total-amount">¥{{ order.totalAmount }}</span>
              </div>
              <div class="order-actions-pc" v-if="order.status === 0">
                <div class="countdown-pc" :class="{ urgent: getCountdownSeconds(order) < 300 }">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                  </svg>
                  <span>剩余 {{ formatCountdown(order) }}</span>
                </div>
                <button class="action-btn-pc cancel" @click="handleCancel(order.orderNo)">取消订单</button>
                <button class="action-btn-pc pay" @click="handlePay(order.orderNo)">立即付款</button>
              </div>
            </div>
          </div>
          <div v-if="orders.length === 0" class="orders-empty-pc">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
            </svg>
            <p>暂无订单</p>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { getSeckillProducts, executeSeckill, getStock, getMyOrders, payOrder, cancelOrder } from '../api/seckill'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const products = ref([])
const filter = ref('all')
const sortBy = ref('default')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = 20
const selectedCategory = ref(0)
const seckillDialogVisible = ref(false)
const seckillResult = ref({ success: false, message: '', orderNo: '' })
const showOrders = ref(false)
const ordersLoading = ref(false)
const orders = ref([])

const countdownHours = ref('02')
const countdownMinutes = ref('00')
const countdownSeconds = ref('00')

let stockTimer = null
let countdownTimer = null
let bannerCountdownTimer = null

const hotWords = ['手机', '电脑', '耳机', '手表', '包包']

const categories = [
  { id: 0, name: '全部商品', icon: 'M4 6h16M4 12h16M4 18h16' },
  { id: 1, name: '手机数码', icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z' },
  { id: 2, name: '电脑办公', icon: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' },
  { id: 3, name: '家用电器', icon: 'M13 10V3L4 14h7v7l9-11h-7z' },
  { id: 4, name: '服装鞋包', icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z' },
  { id: 5, name: '美妆护肤', icon: 'M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z' },
  { id: 6, name: '食品生鲜', icon: 'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z' },
  { id: 7, name: '运动户外', icon: 'M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064' },
]

const activeCount = computed(() => products.value.filter(p => p.status === 1).length)

const filteredProducts = computed(() => {
  let list = products.value
  if (filter.value === 'active') list = list.filter(p => p.status === 1)
  else if (filter.value === 'pending') list = list.filter(p => p.status === 0)
  else if (filter.value === 'ended') list = list.filter(p => p.status === 2)
  if (searchKeyword.value) {
    list = list.filter(p => p.productName?.includes(searchKeyword.value))
  }
  return list
})

const sortedProducts = computed(() => {
  let list = [...filteredProducts.value]
  if (sortBy.value === 'price') {
    list.sort((a, b) => a.seckillPrice - b.seckillPrice)
  } else if (sortBy.value === 'stock') {
    list.sort((a, b) => b.seckillStock - a.seckillStock)
  }
  const start = (currentPage.value - 1) * pageSize
  return list.slice(start, start + pageSize)
})

const totalPages = computed(() => Math.ceil(filteredProducts.value.length / pageSize) || 1)

async function loadProducts() {
  loading.value = true
  try {
    const res = await getSeckillProducts()
    products.value = res.data || []
  } catch (error) {
    console.error('加载商品失败', error)
  } finally {
    loading.value = false
  }
}

async function refreshStock() {
  for (const product of products.value) {
    try {
      const res = await getStock(product.id)
      product.seckillStock = res.data
    } catch (e) {}
  }
}

function getStockPercent(product) {
  if (!product.totalStock) return 50
  return (product.seckillStock / product.totalStock) * 100
}

function getStatusClassByOrder(order) {
  if (order.status === 0) return 'pending'
  if (order.status === 1) return 'paid'
  return 'cancelled'
}

function getStatusTextByOrder(order) {
  if (order.status === 0) return '待付款'
  if (order.status === 1) return '已完成'
  return '已取消'
}

async function handleSeckill(product) {
  if (product.status === 0) {
    ElMessage.info('已设置秒杀提醒')
    return
  }
  
  try {
    await ElMessageBox.confirm('确定要抢购该商品吗？每位用户仅有一次秒杀资格！', '抢购确认', {
      confirmButtonText: '立即抢购',
      cancelButtonText: '取消',
      type: 'warning'
    })
  } catch { return }
  
  try {
    ElMessage.info('正在抢购中...')
    const res = await executeSeckill(product.id)
    seckillResult.value = {
      success: res.data.success,
      message: res.data.message,
      orderNo: res.data.orderNo || ''
    }
    seckillDialogVisible.value = true
    if (res.data.success) {
      await refreshStock()
      setTimeout(() => loadOrders(), 1500)
    }
  } catch (error) {
    const errorMsg = error.message || '抢购失败'
    if (errorMsg.includes('已参与过') || errorMsg.includes('已经参与')) {
      ElMessage.warning('您已参与过该商品的秒杀，无法重复下单！')
    } else {
      seckillResult.value = { success: false, message: errorMsg, orderNo: '' }
      seckillDialogVisible.value = true
    }
  }
}

async function loadOrders() {
  ordersLoading.value = true
  try {
    const res = await getMyOrders()
    orders.value = res.data || []
  } catch (error) {
    console.error('加载订单失败', error)
  } finally {
    ordersLoading.value = false
  }
}

async function handlePay(orderNo) {
  try {
    await ElMessageBox.confirm('确认支付该订单吗？', '支付确认', {
      confirmButtonText: '确认支付',
      cancelButtonText: '取消',
      type: 'info'
    })
    await payOrder(orderNo)
    ElMessage.success('支付成功')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '支付失败')
      await loadOrders()
    }
  }
}

async function handleCancel(orderNo) {
  try {
    await ElMessageBox.confirm('确认取消该订单吗？取消后将无法再次秒杀该商品！', '取消订单', {
      confirmButtonText: '确认取消',
      cancelButtonText: '返回',
      type: 'warning'
    })
    await cancelOrder(orderNo)
    ElMessage.success('订单已取消')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || '取消失败')
  }
}

function handleSearch() {
  currentPage.value = 1
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

function formatDateTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getFullYear()}-${(date.getMonth()+1).toString().padStart(2,'0')}-${date.getDate().toString().padStart(2,'0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

function getCountdownSeconds(order) {
  if (!order.createTime || order.status !== 0) return 0
  const createTime = new Date(order.createTime).getTime()
  const expireTime = createTime + 15 * 60 * 1000
  return Math.max(0, Math.floor((expireTime - Date.now()) / 1000))
}

function formatCountdown(order) {
  const seconds = getCountdownSeconds(order)
  if (seconds <= 0) return '已超时'
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${minutes}分${secs.toString().padStart(2, '0')}秒`
}

function updateBannerCountdown() {
  const now = new Date()
  const end = new Date()
  end.setHours(12, 0, 0, 0)
  if (now.getHours() >= 12) {
    end.setHours(24, 0, 0, 0)
  }
  const diff = Math.max(0, Math.floor((end - now) / 1000))
  const hours = Math.floor(diff / 3600)
  const minutes = Math.floor((diff % 3600) / 60)
  const secs = diff % 60
  countdownHours.value = hours.toString().padStart(2, '0')
  countdownMinutes.value = minutes.toString().padStart(2, '0')
  countdownSeconds.value = secs.toString().padStart(2, '0')
}

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadProducts()
  loadOrders()
  stockTimer = setInterval(refreshStock, 5000)
  countdownTimer = setInterval(() => { orders.value = [...orders.value] }, 1000)
  bannerCountdownTimer = setInterval(updateBannerCountdown, 1000)
  updateBannerCountdown()
})

onUnmounted(() => {
  if (stockTimer) clearInterval(stockTimer)
  if (countdownTimer) clearInterval(countdownTimer)
  if (bannerCountdownTimer) clearInterval(bannerCountdownTimer)
})
</script>

<style scoped>
/* ===== 顶部导航栏 ===== */
.pc-topbar {
  background: #f5f5f5;
  border-bottom: 1px solid #e5e5e5;
}

.topbar-inner {
  max-width: 1200