<template>
  <Layout>
    <div class="seckill-page">
      <!-- 秒杀 Banner 头 -->
      <div class="seckill-banner">
        <div class="banner-content">
          <div class="banner-left">
            <div class="banner-icon">
              <svg viewBox="0 0 24 24" fill="currentColor"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
            </div>
            <div class="banner-text">
              <h1>限时秒杀</h1>
              <p>超值优惠 抢购不停</p>
            </div>
          </div>
          <div class="banner-right">
            <div class="session-tabs">
              <button class="session-tab" :class="{ active: activeSession === 'current' }" @click="switchSession('current')">
                <span class="tab-time">{{ currentSessionTime }}</span>
                <span class="tab-label">正在秒杀</span>
                <span class="tab-count">{{ currentSessionCount }}件</span>
              </button>
              <button class="session-tab" :class="{ active: activeSession === 'next' }" @click="switchSession('next')">
                <span class="tab-time">{{ nextSessionTime }}</span>
                <span class="tab-label">即将开始</span>
                <span class="tab-count">{{ nextSessionCount }}件</span>
              </button>
            </div>
            <div class="session-countdown">
              <span class="countdown-label">{{ activeSession === 'current' ? '距结束' : '距开始' }}</span>
              <div class="countdown-numbers">
                <span class="cd-num">{{ countdownHours }}</span><span class="cd-sep">:</span>
                <span class="cd-num">{{ countdownMinutes }}</span><span class="cd-sep">:</span>
                <span class="cd-num">{{ countdownSeconds }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="seckill-content">
        <aside class="sidebar">
          <div class="sidebar-section">
            <h3 class="sidebar-title">商品分类</h3>
            <div class="category-list">
              <button class="category-btn" :class="{ active: selectedCategory === 0 }" @click="selectedCategory = 0">
                全部商品
              </button>
              <button v-for="cat in categories" :key="cat.id" class="category-btn" :class="{ active: selectedCategory === cat.id }" @click="selectedCategory = cat.id">
                {{ cat.name }}
                <span class="cat-count">{{ getCategoryCount(cat.id) }}</span>
              </button>
            </div>
          </div>
          <div class="sidebar-section">
            <h3 class="sidebar-title">筛选条件</h3>
            <div class="filter-group">
              <label class="filter-label">价格区间</label>
              <div class="price-range">
                <input type="number" v-model.number="priceMin" placeholder="最低价" class="price-input" />
                <span class="price-sep">-</span>
                <input type="number" v-model.number="priceMax" placeholder="最高价" class="price-input" />
              </div>
              <button class="apply-filter-btn" @click="applyPriceFilter">应用</button>
            </div>
            <div class="filter-group">
              <label class="filter-label">折扣力度</label>
              <div class="discount-options">
                <button class="discount-btn" :class="{ active: discountFilter === 'all' }" @click="discountFilter = 'all'">全部</button>
                <button class="discount-btn" :class="{ active: discountFilter === '3' }" @click="discountFilter = '3'">3折以下</button>
                <button class="discount-btn" :class="{ active: discountFilter === '5' }" @click="discountFilter = '5'">5折以下</button>
                <button class="discount-btn" :class="{ active: discountFilter === '7' }" @click="discountFilter = '7'">7折以下</button>
              </div>
            </div>
            <div class="filter-group">
              <label class="filter-label">库存状态</label>
              <div class="stock-options">
                <button class="stock-btn" :class="{ active: stockFilter === 'all' }" @click="stockFilter = 'all'">全部</button>
                <button class="stock-btn" :class="{ active: stockFilter === 'available' }" @click="stockFilter = 'available'">有库存</button>
                <button class="stock-btn" :class="{ active: stockFilter === 'low' }" @click="stockFilter = 'low'">库存紧张</button>
              </div>
            </div>
          </div>
        </aside>

        <div class="main-content">
          <div class="filter-bar">
            <div class="filter-tabs">
              <button class="filter-tab" :class="{ active: statusFilter === 'all' }" @click="statusFilter = 'all'">全部</button>
              <button class="filter-tab" :class="{ active: statusFilter === 'active' }" @click="statusFilter = 'active'">正在秒杀</button>
              <button class="filter-tab" :class="{ active: statusFilter === 'pending' }" @click="statusFilter = 'pending'">即将开始</button>
              <button class="filter-tab" :class="{ active: statusFilter === 'ended' }" @click="statusFilter = 'ended'">已结束</button>
            </div>
            <div class="filter-right">
              <span class="result-count">共 {{ filteredProducts.length }} 件商品</span>
              <div class="sort-options">
                <button class="sort-btn" :class="{ active: sortBy === 'default' }" @click="sortBy = 'default'">综合</button>
                <button class="sort-btn" :class="{ active: sortBy === 'price-asc' }" @click="sortBy = 'price-asc'">价格↑</button>
                <button class="sort-btn" :class="{ active: sortBy === 'price-desc' }" @click="sortBy = 'price-desc'">价格↓</button>
                <button class="sort-btn" :class="{ active: sortBy === 'stock' }" @click="sortBy = 'stock'">库存</button>
                <button class="sort-btn" :class="{ active: sortBy === 'discount' }" @click="sortBy = 'discount'">折扣</button>
              </div>
            </div>
          </div>

          <div class="products-grid" v-loading="loading">
            <div v-for="(product, index) in paginatedProducts" :key="product.id" class="product-card" @click="goProductDetail(product)">
              <div class="card-image">
                <div class="image-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
                <div class="card-tags">
                  <span class="tag-seckill" v-if="product.status === 1">秒杀</span>
                  <span class="tag-discount" v-if="getDiscount(product) > 0">{{ getDiscount(product) }}折</span>
                </div>
                <div class="card-overlay" v-if="product.status !== 1 || product.seckillStock <= 0">
                  <span v-if="product.seckillStock <= 0">已售罄</span>
                  <span v-else-if="product.status === 2">已结束</span>
                  <span v-else-if="product.status === 0">即将开始</span>
                </div>
                <div class="card-progress" v-if="product.status === 1 && product.seckillStock > 0">
                  <div class="progress-bar"><div class="progress-fill" :style="{ width: `${getSoldPercent(product)}%` }"></div></div>
                  <span class="progress-text">已抢{{ getSoldPercent(product) }}%</span>
                </div>
              </div>
              <div class="card-body">
                <h3 class="card-title">{{ product.productName || `秒杀商品 #${product.id}` }}</h3>
                <div class="card-price">
                  <div class="price-current"><span class="price-symbol">¥</span><span class="price-value">{{ product.seckillPrice }}</span></div>
                  <div class="price-original" v-if="product.originalPrice && product.originalPrice > product.seckillPrice"><span>¥{{ product.originalPrice }}</span></div>
                </div>
                <div class="card-meta">
                  <span class="meta-stock">库存 {{ product.seckillStock }}</span>
                  <span class="meta-time">{{ formatTime(product.startTime) }} - {{ formatTime(product.endTime) }}</span>
                </div>
                <div class="card-countdown" v-if="product.status === 0">
                  <span class="countdown-label-small">距开始</span>
                  <span class="countdown-value-small">{{ getProductCountdown(product) }}</span>
                </div>
                <button class="btn-buy" :class="{ disabled: product.status !== 1 || product.seckillStock <= 0 }" :disabled="product.status !== 1 || product.seckillStock <= 0" @click.stop="handleSeckill(product)">
                  <span v-if="product.seckillStock <= 0">已售罄</span>
                  <span v-else-if="product.status === 2">已结束</span>
                  <span v-else-if="product.status === 0">即将开始</span>
                  <span v-else>立即抢购</span>
                </button>
              </div>
            </div>
            <div v-if="!loading && paginatedProducts.length === 0" class="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/><path d="M8 15s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无符合条件的秒杀商品</p>
              <button class="reset-btn" @click="resetFilters">重置筛选</button>
            </div>
          </div>

          <div class="pagination-bar" v-if="filteredProducts.length > pageSize">
            <button class="page-btn prev" :disabled="currentPage === 1" @click="currentPage--">上一页</button>
            <div class="page-numbers">
              <button v-for="page in visiblePages" :key="page" class="page-num" :class="{ active: currentPage === page }" @click="currentPage = page">{{ page }}</button>
            </div>
            <button class="page-btn next" :disabled="currentPage === totalPages" @click="currentPage++">下一页</button>
            <span class="page-info">共 {{ totalPages }} 页</span>
          </div>
        </div>
      </div>

      <el-dialog v-model="seckillDialogVisible" title="" width="420px" :close-on-click-modal="false" class="seckill-result-dialog">
        <div class="result-content">
          <div class="result-icon" :class="seckillResult.success ? 'success' : 'fail'">
            <svg v-if="seckillResult.success" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>
            <svg v-else viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"/></svg>
          </div>
          <h3 class="result-title">{{ seckillResult.success ? '恭喜！秒杀成功' : '秒杀失败' }}</h3>
          <p class="result-desc">{{ seckillResult.message }}</p>
          <div v-if="seckillResult.success" class="result-order"><span class="order-label">订单号</span><span class="order-value">{{ seckillResult.orderNo }}</span></div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <button class="dialog-btn secondary" @click="seckillDialogVisible = false">{{ seckillResult.success ? '稍后支付' : '继续抢购' }}</button>
            <button v-if="seckillResult.success" class="dialog-btn primary" @click="goPay">立即支付</button>
          </div>
        </template>
      </el-dialog>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'
import { getSeckillProducts, executeSeckill, getStock } from '../api/seckill'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const products = ref([])
const categories = ref([])
const selectedCategory = ref(0)
const statusFilter = ref('all')
const sortBy = ref('default')
const priceMin = ref(null)
const priceMax = ref(null)
const discountFilter = ref('all')
const stockFilter = ref('all')
const currentPage = ref(1)
const pageSize = 12
const activeSession = ref('current')

const countdownHours = ref('00')
const countdownMinutes = ref('00')
const countdownSeconds = ref('00')

const seckillDialogVisible = ref(false)
const seckillResult = ref({ success: false, message: '', orderNo: '' })

let countdownTimer = null
let stockTimer = null

const defaultCategories = [
  { id: 1, name: '手机数码' }, { id: 2, name: '电脑办公' }, { id: 3, name: '家用电器' }, { id: 4, name: '服装鞋包' },
]

const currentSessionTime = computed(() => {
  const hour = new Date().getHours()
  if (hour < 10) return '10:00'; if (hour < 12) return '12:00'; if (hour < 14) return '14:00'
  if (hour < 16) return '16:00'; if (hour < 18) return '18:00'; if (hour < 20) return '20:00'
  if (hour < 22) return '22:00'; return '00:00'
})

const nextSessionTime = computed(() => {
  const hour = new Date().getHours()
  if (hour < 10) return '12:00'; if (hour < 12) return '14:00'; if (hour < 14) return '16:00'
  if (hour < 16) return '18:00'; if (hour < 18) return '20:00'; if (hour < 20) return '22:00'
  if (hour < 22) return '00:00'; return '10:00'
})

const currentSessionCount = computed(() => products.value.filter(p => p.status === 1).length)
const nextSessionCount = computed(() => products.value.filter(p => p.status === 0).length)

const filteredProducts = computed(() => {
  let list = [...products.value]
  if (activeSession.value === 'current') list = list.filter(p => p.status === 1)
  else if (activeSession.value === 'next') list = list.filter(p => p.status === 0)
  if (statusFilter.value === 'active') list = list.filter(p => p.status === 1)
  else if (statusFilter.value === 'pending') list = list.filter(p => p.status === 0)
  else if (statusFilter.value === 'ended') list = list.filter(p => p.status === 2)
  if (selectedCategory.value > 0) list = list.filter(p => p.categoryId === selectedCategory.value)
  if (priceMin.value !== null && priceMin.value >= 0) list = list.filter(p => p.seckillPrice >= priceMin.value)
  if (priceMax.value !== null && priceMax.value > 0) list = list.filter(p => p.seckillPrice <= priceMax.value)
  if (discountFilter.value !== 'all') {
    const max = Number(discountFilter.value)
    list = list.filter(p => p.originalPrice && (p.seckillPrice / p.originalPrice) * 10 <= max)
  }
  if (stockFilter.value === 'available') list = list.filter(p => p.seckillStock > 0)
  else if (stockFilter.value === 'low') list = list.filter(p => p.seckillStock > 0 && p.seckillStock <= 10)
  const kw = route.query.keyword
  if (kw) list = list.filter(p => p.productName?.includes(kw))
  const cq = route.query.category
  if (cq) list = list.filter(p => p.categoryId === Number(cq))
  return list
})

const sortedProducts = computed(() => {
  let list = [...filteredProducts.value]
  if (sortBy.value === 'price-asc') list.sort((a, b) => a.seckillPrice - b.seckillPrice)
  else if (sortBy.value === 'price-desc') list.sort((a, b) => b.seckillPrice - a.seckillPrice)
  else if (sortBy.value === 'stock') list.sort((a, b) => b.seckillStock - a.seckillStock)
  else if (sortBy.value === 'discount') list.sort((a, b) => {
    const da = a.originalPrice ? (a.seckillPrice / a.originalPrice) : 1
    const db = b.originalPrice ? (b.seckillPrice / b.originalPrice) : 1
    return da - db
  })
  return list
})

const totalPages = computed(() => Math.ceil(sortedProducts.value.length / pageSize) || 1)
const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return sortedProducts.value.slice(start, start + pageSize)
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value
  if (total <= 5) { for (let i = 1; i <= total; i++) pages.push(i) }
  else if (current <= 3) pages.push(1, 2, 3, 4, 5)
  else if (current >= total - 2) pages.push(total - 4, total - 3, total - 2, total - 1, total)
  else pages.push(current - 2, current - 1, current, current + 1, current + 2)
  return pages
})

function getCategoryCount(catId) { return products.value.filter(p => p.categoryId === catId).length }
function getDiscount(product) { if (!product.originalPrice) return 0; const d = Math.round((product.seckillPrice / product.originalPrice) * 10); return d < 10 ? d : 0 }
function getSoldPercent(product) { if (!product.totalStock || !product.seckillStock) return 50; return Math.round(((product.totalStock - product.seckillStock) / product.totalStock) * 100) }
function formatTime(time) { if (!time) return ''; const d = new Date(time); return `${d.getHours().toString().padStart(2,'0')}:${d.getMinutes().toString().padStart(2,'0')}` }
function getProductCountdown(product) { if (!product.startTime) return ''; const diff = Math.max(0, Math.floor((new Date(product.startTime) - new Date()) / 1000)); if (diff <= 0) return '即将开始'; const h = Math.floor(diff / 3600); const m = Math.floor((diff % 3600) / 60); return h > 0 ? `${h}小时${m}分钟` : `${m}分钟` }

async function loadProducts() {
  loading.value = true
  try {
    const res = await getSeckillProducts()
    products.value = res.data || []
    const cm = new Map()
    products.value.forEach(p => { if (p.categoryId && p.categoryName) cm.set(p.categoryId, p.categoryName) })
    if (cm.size > 0) { cm.forEach((name, id) => categories.value.push({ id, name })) }
    else categories.value = defaultCategories
    updateSessionCountdown()
  } catch (e) { console.error('加载商品失败', e); categories.value = defaultCategories } finally { loading.value = false }
}

async function refreshStock() {
  for (const p of products.value) { if (p.status === 1 && p.seckillStock > 0) { try { const r = await getStock(p.id); p.seckillStock = r.data } catch (e) {} } }
}

function switchSession(s) { activeSession.value = s; currentPage.value = 1; updateSessionCountdown() }

function updateSessionCountdown() {
  const now = new Date()
  let target = null
  if (activeSession.value === 'current') {
    const ap = products.value.filter(p => p.status === 1)
    target = ap.length > 0 ? new Date(ap[0].endTime) : new Date(new Date().setHours(new Date().getHours() + 2, 0, 0, 0))
  } else {
    const pp = products.value.filter(p => p.status === 0)
    target = pp.length > 0 ? new Date(pp[0].startTime) : new Date(new Date().setHours(new Date().getHours() + 4, 0, 0, 0))
  }
  if (target) { const diff = Math.max(0, Math.floor((target - now) / 1000)); countdownHours.value = Math.floor(diff / 3600).toString().padStart(2, '0'); countdownMinutes.value = Math.floor((diff % 3600) / 60).toString().padStart(2, '0'); countdownSeconds.value = (diff % 60).toString().padStart(2, '0') }
}

function applyPriceFilter() { currentPage.value = 1 }
function resetFilters() { selectedCategory.value = 0; statusFilter.value = 'all'; priceMin.value = null; priceMax.value = null; discountFilter.value = 'all'; stockFilter.value = 'all'; sortBy.value = 'default'; currentPage.value = 1 }

async function handleSeckill(product) {
  if (product.status === 0) { ElMessage.info('已设置秒杀提醒，开抢前会通知您'); return }
  if (product.seckillStock <= 0) { ElMessage.warning('商品已售罄'); return }
  try { await ElMessageBox.confirm(`确定要抢购「${product.productName}」吗？`, '抢购确认', { confirmButtonText: '立即抢购', cancelButtonText: '取消', type: 'warning' }) } catch { return }
  try {
    ElMessage.info('正在抢购中，请稍候...')
    const res = await executeSeckill(product.id)
    seckillResult.value = { success: res.data.success, message: res.data.message, orderNo: res.data.orderNo || '' }
    seckillDialogVisible.value = true
    if (res.data.success) await refreshStock()
  } catch (error) {
    const em = error.message || '抢购失败'
    if (em.includes('已参与过') || em.includes('已经参与')) ElMessage.warning('您已参与过该商品的秒杀，无法重复下单！')
    else { seckillResult.value = { success: false, message: em, orderNo: '' }; seckillDialogVisible.value = true }
  }
}

function goProductDetail(p) { router.push(`/product/${p.id}?type=seckill`) }
function goPay() { seckillDialogVisible.value = false; router.push('/orders') }

watch([statusFilter, selectedCategory, discountFilter, stockFilter, sortBy, activeSession], () => { currentPage.value = 1 })

onMounted(() => { loadProducts(); countdownTimer = setInterval(updateSessionCountdown, 1000); stockTimer = setInterval(refreshStock, 5000) })
onUnmounted(() => { if (countdownTimer) clearInterval(countdownTimer); if (stockTimer) clearInterval(stockTimer) })
</script>

<style scoped>
.seckill-page { display: flex; flex-direction: column; gap: 16px; padding-top: 12px; }

.seckill-banner { background: linear-gradient(135deg, #FF5000, #FF3300); border-radius: 12px; overflow: hidden; }
.banner-content { padding: 24px; display: flex; justify-content: space-between; align-items: center; }
.banner-left { display: flex; align-items: center; gap: 16px; }
.banner-icon { width: 64px; height: 64px; background: rgba(255,255,255,0.2); border-radius: 16px; display: flex; align-items: center; justify-content: center; color: #fff; }
.banner-icon svg { width: 36px; height: 36px; }
.banner-text h1 { font-size: 28px; font-weight: 700; color: #fff; margin: 0; }
.banner-text p { font-size: 14px; color: rgba(255,255,255,0.9); margin: 4px 0 0; }
.banner-right { display: flex; align-items: center; gap: 24px; }

.session-tabs { display: flex; gap: 12px; }
.session-tab { display: flex; flex-direction: column; align-items: center; padding: 12px 20px; background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); border-radius: 8px; cursor: pointer; transition: all 0.2s; }
.session-tab.active { background: #fff; border-color: #fff; }
.tab-time { font-size: 20px; font-weight: 700; color: #fff; }
.session-tab.active .tab-time { color: #FF5000; }
.tab-label { font-size: 12px; color: rgba(255,255,255,0.8); margin-top: 4px; }
.session-tab.active .tab-label { color: #FF5000; }
.tab-count { font-size: 11px; color: rgba(255,255,255,0.6); margin-top: 2px; }
.session-tab.active .tab-count { color: #FF5000; }

.session-countdown { display: flex; align-items: center; gap: 12px; }
.countdown-label { font-size: 14px; color: #fff; }
.countdown-numbers { display: flex; align-items: center; gap: 4px; }
.cd-num { font-size: 24px; font-weight: 700; color: #fff; background: rgba(255,255,255,0.2); padding: 6px 10px; border-radius: 6px; }
.cd-sep { font-size: 20px; color: #fff; font-weight: 700; }

.seckill-content { display: flex; gap: 16px; }
.sidebar { width: 240px; display: flex; flex-direction: column; gap: 12px; flex-shrink: 0; }
.sidebar-section { background: #fff; border-radius: 12px; padding: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.sidebar-title { font-size: 14px; font-weight: 600; color: #333; margin: 0 0 12px; }
.category-list { display: flex; flex-direction: column; gap: 6px; }
.category-btn { display: flex; justify-content: space-between; align-items: center; padding: 10px 12px; background: transparent; border: none; border-radius: 8px; font-size: 13px; color: #666; cursor: pointer; transition: all 0.2s; text-align: left; }
.category-btn:hover { background: #f5f5f5; }
.category-btn.active { background: #FFF0E6;
  color: #FF5000; font-weight: 500; }
.cat-count { font-size: 12px; color: #999; }

.filter-group { margin-bottom: 16px; }
.filter-label { font-size: 12px; color: #999; margin-bottom: 6px; display: block; }
.price-range { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.price-input { width: 80px; padding: 8px 12px; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 13px; }
.price-input:focus { border-color: #FF5000; outline: none; }
.price-sep { color: #999; }
.apply-filter-btn { width: 100%; padding: 8px; background: #FF5000; border: none; border-radius: 6px; color: #fff; font-size: 12px; cursor: pointer; }
.discount-options, .stock-options { display: flex; flex-wrap: wrap; gap: 8px; }
.discount-btn, .stock-btn { padding: 6px 12px; background: transparent; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 12px; color: #666; cursor: pointer; transition: all 0.2s; }
.discount-btn.active, .stock-btn.active { background: #FF5000; border-color: #FF5000; color: #fff; }

.main-content { flex: 1; display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.filter-bar { background: #fff; border-radius: 12px; padding: 16px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.filter-tabs { display: flex; gap: 8px; }
.filter-tab { padding: 8px 16px; background: transparent; border: none; border-radius: 8px; font-size: 14px; color: #666; cursor: pointer; transition: all 0.2s; }
.filter-tab.active { background: #FF5000; color: #fff; }
.filter-right { display: flex; align-items: center; gap: 16px; }
.result-count { font-size: 13px; color: #999; }
.sort-options { display: flex; gap: 8px; }
.sort-btn { padding: 6px 12px; background: transparent; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 13px; color: #666; cursor: pointer; transition: all 0.2s; }
.sort-btn.active { border-color: #FF5000; color: #FF5000; }

.products-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px; min-height: 400px; }
.product-card { background: #fff; border-radius: 12px; overflow: hidden; cursor: pointer; transition: all 0.3s; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.product-card:hover { transform: translateY(-4px); box-shadow: 0 12px 24px rgba(0,0,0,0.12); }
.card-image { position: relative; height: 160px; background: #f8f8f8; }
.image-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; }
.image-placeholder svg { width: 48px; height: 48px; color: #ccc; }
.card-tags { position: absolute; top: 8px; left: 8px; display: flex; gap: 6px; }
.tag-seckill { background: linear-gradient(135deg, #FF5000, #FF3300); color: #fff; padding: 4px 10px; border-radius: 4px; font-size: 12px; font-weight: 600; }
.tag-discount { background: #FF5000; color: #fff; padding: 4px 10px; border-radius: 4px; font-size: 12px; }
.card-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.card-overlay span { background: rgba(255,255,255,0.9); color: #FF5000; padding: 12px 24px; border-radius: 8px; font-size: 16px; font-weight: 600; }
.card-progress { position: absolute; bottom: 8px; left: 8px; right: 8px; display: flex; align-items: center; gap: 8px; }
.progress-bar { flex: 1; height: 4px; background: rgba(255,255,255,0.3); border-radius: 2px; overflow: hidden; }
.progress-fill { height: 100%; background: linear-gradient(90deg, #FF5000, #FF3300); }
.progress-text { font-size: 11px; color: #fff; background: rgba(0,0,0,0.5); padding: 2px 6px; border-radius: 4px; }
.card-body { padding: 12px; }
.card-title { font-size: 14px; font-weight: 500; color: #333; margin: 0 0 8px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-price { display: flex; align-items: baseline; gap: 8px; margin-bottom: 8px; }
.price-symbol { font-size: 14px; color: #FF5000; font-weight: 600; }
.price-value { font-size: 18px; color: #FF5000; font-weight: 700; }
.price-original { font-size: 12px; color: #999; text-decoration: line-through; }
.card-meta { display: flex; justify-content: space-between; font-size: 12px; color: #999; margin-bottom: 8px; }
.meta-time { color: #FF5000; }
.card-countdown { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; padding: 6px 8px; background: #FFF0E6; border-radius: 6px; }
.countdown-label-small { font-size: 12px; color: #FF5000; }
.countdown-value-small { font-size: 12px; color: #FF5000; font-weight: 600; }
.btn-buy { width: 100%; padding: 10px; background: linear-gradient(135deg, #FF5000, #FF3300); border: none; border-radius: 8px; color: #fff; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.2s; }
.btn-buy:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 4px 12px rgba(255,80,0,0.3); }
.btn-buy.disabled { background: #ccc; cursor: not-allowed; }

.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 60px; grid-column: 1 / -1; }
.empty-state svg { width: 64px; height: 64px; color: #ccc; }
.empty-state p { font-size: 14px; color: #999; margin: 16px 0; }
.reset-btn { padding: 12px 24px; background: #FF5000; border: none; border-radius: 8px; color: #fff; font-size: 14px; cursor: pointer; }
.empty-btn { padding: 12px 24px; background: linear-gradient(135deg, #FF5000, #FF3300); border: none; border-radius: 8px; color: #fff; font-size: 14px; cursor: pointer; }

.pagination-bar { background: #fff; border-radius: 12px; padding: 16px; display: flex; justify-content: center; align-items: center; gap: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.page-btn { padding: 8px 16px; background: transparent; border: 1px solid #e5e5e5; border-radius: 8px; font-size: 14px; color: #666; cursor: pointer; transition: all 0.2s; }
.page-btn:hover:not(:disabled) { border-color: #FF5000; color: #FF5000; }
.page-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.page-numbers { display: flex; gap: 8px; }
.page-num { padding: 8px 12px; background: transparent; border: 1px solid #e5e5e5; border-radius: 6px; font-size: 13px; color: #666; cursor: pointer; transition: all 0.2s; }
.page-num.active { background: #FF5000; border-color: #FF5000; color: #fff; }
.page-info { font-size: 13px; color: #999; }

.seckill-result-dialog .result-content { display: flex; flex-direction: column; align-items: center; padding: 24px; }
.result-icon { width: 80px; height: 80px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-bottom: 16px; }
.result-icon.success { background: #E8F5E9; color: #4CAF50; }
.result-icon.fail { background: #FFEBEE; color: #F44336; }
.result-icon svg { width: 48px; height: 48px; }
.result-title { font-size: 20px; font-weight: 600; color: #333; margin: 0 0 8px; }
.result-desc { font-size: 14px; color: #666; margin: 0; }
.result-order { display: flex; align-items: center; gap: 8px; margin-top: 16px; padding: 12px 16px; background: #f5f5f5; border-radius: 8px; }
.order-label { font-size: 13px; color: #999; }
.order-value { font-size: 14px; color: #333; font-weight: 500; }
.dialog-footer { display: flex; justify-content: center; gap: 12px; }
.dialog-btn { padding: 12px 24px; border-radius: 8px; font-size: 14px; cursor: pointer; }
.dialog-btn.secondary { background: #f5f5f5; border: none; color: #666; }
.dialog-btn.primary { background: linear-gradient(135deg, #FF5000, #FF3300); border: none; color: #fff; }

@media (max-width: 1200px) { .sidebar { width: 200px; } }
@media (max-width: 900px) { .seckill-content { flex-direction: column; } .sidebar { width: 100%; } .products-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 600px) { .products-grid { grid-template-columns: 1fr; } .banner-content { flex-direction: column; gap: 16px; } .session-tabs { flex-direction: column; } }
</style>