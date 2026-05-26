<template>
  <Layout>
    <div class="home-page">
      <div class="main-layout">
        <aside class="left-sidebar">
          <div class="sidebar-card announcement-card">
            <div class="card-header">
              <div class="header-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
                  <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
                </svg>
              </div>
              <h3>活动公告</h3>
            </div>
            <div class="announcement-list">
              <div class="announcement-item" v-for="(item, index) in announcements" :key="index">
                <div class="announcement-tag" :class="item.type">{{ item.tag }}</div>
                <div class="announcement-content">
                  <h4>{{ item.title }}</h4>
                  <p>{{ item.desc }}</p>
                  <span class="announcement-time">{{ item.time }}</span>
                </div>
              </div>
            </div>
          </div>
        </aside>

        <div class="main-content">
          <section class="banner-section">
            <div class="banner-carousel" @mouseenter="pauseAutoPlay" @mouseleave="resumeAutoPlay">
              <div class="carousel-wrapper">
                <div 
                  class="carousel-slide"
                  :class="{ active: currentSlide === 0 }"
                >
                  <div class="slide-content slide-1">
                    <div class="slide-text">
                      <h2>限时秒杀</h2>
                      <p>全场低至1折起</p>
                      <router-link to="/seckill" class="slide-btn">立即抢购</router-link>
                    </div>
                  </div>
                </div>
                <div 
                  class="carousel-slide"
                  :class="{ active: currentSlide === 1 }"
                >
                  <div class="slide-content slide-2">
                    <div class="slide-text">
                      <h2>新品首发</h2>
                      <p>最新数码产品抢先体验</p>
                      <router-link to="/seckill" class="slide-btn">查看详情</router-link>
                    </div>
                  </div>
                </div>
                <div 
                  class="carousel-slide"
                  :class="{ active: currentSlide === 2 }"
                >
                  <div class="slide-content slide-3">
                    <div class="slide-text">
                      <h2>品质保障</h2>
                      <p>正品保证 七天无忧退换</p>
                      <router-link to="/seckill" class="slide-btn">了解更多</router-link>
                    </div>
                  </div>
                </div>
              </div>
              <div class="carousel-indicators">
                <button 
                  v-for="(_, index) in 3"
                  :key="index"
                  class="indicator"
                  :class="{ active: currentSlide === index }"
                  @click="goToSlide(index)"
                ></button>
              </div>
              <button class="carousel-arrow prev" @click="prevSlide">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15 18 9 12 15 6"/>
                </svg>
              </button>
              <button class="carousel-arrow next" @click="nextSlide">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9 18 15 12 9 6"/>
                </svg>
              </button>
            </div>
          </section>

          <section class="seckill-section">
            <div class="section-header">
              <div class="header-left">
                <div class="section-icon">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
                  </svg>
                </div>
                <h2 class="section-title">限时秒杀</h2>
                <span class="section-badge">HOT</span>
              </div>
              <div class="header-right">
                <div class="seckill-countdown" v-if="currentSeckillProduct">
                  <span class="countdown-label">{{ countdownLabel }}</span>
                  <div class="countdown-numbers">
                    <span class="cd-num">{{ countdownHours }}</span>
                    <span class="cd-sep">:</span>
                    <span class="cd-num">{{ countdownMinutes }}</span>
                    <span class="cd-sep">:</span>
                    <span class="cd-num">{{ countdownSeconds }}</span>
                  </div>
                </div>
                <router-link to="/seckill" class="header-link">
                  查看全部
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="9 18 15 12 9 6"/>
                  </svg>
                </router-link>
              </div>
            </div>
            
            <div class="seckill-grid" v-loading="loading">
              <ProductCard
                v-for="product in seckillProducts.slice(0, 4)"
                :key="product.id"
                :product="product"
                :isSeckill="true"
                @buy="handleSeckill"
                @click="goProductDetail"
              />
              
              <div v-if="!loading && seckillProducts.length === 0" class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                  <line x1="9" y1="9" x2="9.01" y2="9"/>
                  <line x1="15" y1="9" x2="15.01" y2="9"/>
                </svg>
                <p>暂无秒杀商品</p>
                <router-link to="/seckill" class="empty-btn">查看更多</router-link>
              </div>
            </div>
          </section>

          <section class="category-section">
            <div class="section-header">
              <h2 class="section-title">商品分类</h2>
            </div>
            <div class="category-grid">
              <div 
                v-for="cat in categories" 
                :key="cat.id"
                class="category-item"
                @click="goCategory(cat.id)"
              >
                <div class="category-icon" :style="{ background: cat.color }">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path :d="cat.icon"/>
                  </svg>
                </div>
                <span class="category-name">{{ cat.name }}</span>
                <span class="category-count">{{ getCategoryCount(cat.id) }}件</span>
              </div>
            </div>
          </section>

          <section class="features-section">
            <div class="features-grid">
              <div class="feature-item">
                <div class="feature-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  </svg>
                </div>
                <h3>正品保障</h3>
                <p>100%正品承诺</p>
              </div>
              <div class="feature-item">
                <div class="feature-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="1" y="3" width="15" height="13"/>
                    <polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/>
                    <circle cx="5.5" cy="18.5" r="2.5"/>
                    <circle cx="18.5" cy="18.5" r="2.5"/>
                  </svg>
                </div>
                <h3>极速发货</h3>
                <p>24小时内发货</p>
              </div>
              <div class="feature-item">
                <div class="feature-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                    <polyline points="22,6 12,13 2,6"/>
                  </svg>
                </div>
                <h3>七天退换</h3>
                <p>不满意可退换</p>
              </div>
              <div class="feature-item">
                <div class="feature-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
                    <line x1="12" y1="17" x2="12.01" y2="17"/>
                  </svg>
                </div>
                <h3>专业客服</h3>
                <p>7x24在线服务</p>
              </div>
            </div>
          </section>
        </div>

        <aside class="right-sidebar">
          <div class="sidebar-card hot-rank-card">
            <div class="card-header">
              <div class="header-icon hot">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <h3>热销榜单</h3>
              <span class="rank-badge">TOP 10</span>
            </div>
            <div class="rank-list">
              <div 
                class="rank-item" 
                v-for="(product, index) in hotProducts" 
                :key="product.id"
                @click="goProductDetail(product)"
              >
                <div class="rank-number" :class="{ top3: index < 3 }">{{ index + 1 }}</div>
                <div class="rank-product-image">
                  <img :src="product.productImage || product.image" alt="" />
                </div>
                <div class="rank-product-info">
                  <h4 class="rank-product-name">{{ product.productName || product.name }}</h4>
                  <div class="rank-product-price">
                    <span class="current-price">¥{{ product.seckillPrice }}</span>
                    <span class="original-price">¥{{ product.originalPrice }}</span>
                  </div>
                  <div class="rank-sales">已售 {{ product.soldCount || Math.floor(Math.random() * 1000) }}件</div>
                </div>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'
import ProductCard from '../components/ProductCard.vue'
import { getSeckillProducts, executeSeckill, getStock } from '../api/seckill'

const router = useRouter()

const loading = ref(false)
const seckillProducts = ref([])
const hotProducts = ref([])
const categories = ref([])
const currentSlide = ref(0)
const currentSeckillProduct = ref(null)

const countdownHours = ref('00')
const countdownMinutes = ref('00')
const countdownSeconds = ref('00')
const countdownLabel = ref('距结束')

let slideTimer = null
let countdownTimer = null

const announcements = ref([
  { tag: '限时', type: 'hot', title: '618年中大促即将开启', desc: '全场商品低至5折，更有神秘惊喜等你来抢', time: '2024-06-18' },
  { tag: '新品', type: 'new', title: '新款手机首发预约', desc: '预约享专属优惠，首发限量抢购', time: '2024-05-20' },
  { tag: '公告', type: 'normal', title: '平台服务升级通知', desc: '物流配送全面升级，极速送达更便捷', time: '2024-05-15' },
  { tag: '活动', type: 'event', title: '新人专享礼包', desc: '注册即送优惠券，首单立减50元', time: '长期有效' },
])

const defaultCategories = [
  { id: 1, name: '手机数码', icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z', color: '#FF5000' },
  { id: 2, name: '电脑办公', icon: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z', color: '#1677FF' },
  { id: 3, name: '家用电器', icon: 'M13 10V3L4 14h7v7l9-11h-7z', color: '#52C41A' },
  { id: 4, name: '服装鞋包', icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z', color: '#FAAD14' },
]

function getCategoryCount(catId) {
  return seckillProducts.value.filter(p => p.categoryId === catId).length || 0
}

async function loadProducts() {
  loading.value = true
  try {
    const res = await getSeckillProducts()
    seckillProducts.value = res.data || []
    
    hotProducts.value = [...seckillProducts.value]
      .sort((a, b) => (b.soldCount || 0) - (a.soldCount || 0))
      .slice(0, 10)
    
    const categoryMap = new Map()
    seckillProducts.value.forEach(p => {
      if (p.categoryId && p.categoryName) {
        categoryMap.set(p.categoryId, p.categoryName)
      }
    })
    
    const colors = ['#FF5000', '#1677FF', '#52C41A', '#FAAD14', '#EB2F96', '#13C2C2']
    const icons = defaultCategories.map(c => c.icon)
    
    if (categoryMap.size > 0) {
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
    } else {
      categories.value = defaultCategories
    }
    
    const activeProducts = seckillProducts.value.filter(p => p.status === 1)
    if (activeProducts.length > 0) {
      currentSeckillProduct.value = activeProducts[0]
    } else {
      const pendingProducts = seckillProducts.value.filter(p => p.status === 0)
      if (pendingProducts.length > 0) {
        currentSeckillProduct.value = pendingProducts[0]
        countdownLabel.value = '距开始'
      }
    }
    
    updateCountdown()
  } catch (error) {
    console.error('加载商品失败', error)
    categories.value = defaultCategories
  } finally {
    loading.value = false
  }
}

function updateCountdown() {
  if (!currentSeckillProduct.value) return
  
  const now = new Date()
  const product = currentSeckillProduct.value
  let targetTime
  
  if (product.status === 1) {
    countdownLabel.value = '距结束'
    targetTime = new Date(product.endTime)
  } else if (product.status === 0) {
    countdownLabel.value = '距开始'
    targetTime = new Date(product.startTime)
  } else {
    countdownHours.value = '00'
    countdownMinutes.value = '00'
    countdownSeconds.value = '00'
    return
  }
  
  const diff = Math.max(0, Math.floor((targetTime - now) / 1000))
  countdownHours.value = Math.floor(diff / 3600).toString().padStart(2, '0')
  countdownMinutes.value = Math.floor((diff % 3600) / 60).toString().padStart(2, '0')
  countdownSeconds.value = (diff % 60).toString().padStart(2, '0')
}

function nextSlide() {
  currentSlide.value = (currentSlide.value + 1) % 3
}

function prevSlide() {
  currentSlide.value = (currentSlide.value - 1 + 3) % 3
}

function goToSlide(index) {
  currentSlide.value = index
}

function pauseAutoPlay() {
  if (slideTimer) clearInterval(slideTimer)
}

function resumeAutoPlay() {
  slideTimer = setInterval(nextSlide, 4000)
}

async function handleSeckill(product) {
  if (product.status === 0) {
    ElMessage.info('已设置秒杀提醒')
    return
  }
  
  if (product.seckillStock <= 0) {
    ElMessage.warning('商品已售罄')
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
    if (res.data.success) {
      ElMessage.success(`秒杀成功！订单号：${res.data.orderNo}`)
      try {
        const stockRes = await getStock(product.id)
        product.seckillStock = stockRes.data
      } catch (e) {}
    } else {
      ElMessage.error(res.data.message || '秒杀失败')
    }
  } catch (error) {
    const errorMsg = error.message || '抢购失败'
    if (errorMsg.includes('已参与过') || errorMsg.includes('已经参与')) {
      ElMessage.warning('您已参与过该商品的秒杀，无法重复下单！')
    } else {
      ElMessage.error(errorMsg)
    }
  }
}

function goProductDetail(product) {
  router.push(`/product/${product.id}?type=seckill`)
}

function goCategory(catId) {
  router.push({ path: '/seckill', query: { category: catId } })
}

onMounted(() => {
  loadProducts()
  slideTimer = setInterval(nextSlide, 4000)
  countdownTimer = setInterval(updateCountdown, 1000)
})

onUnmounted(() => {
  if (slideTimer) clearInterval(slideTimer)
  if (countdownTimer) clearInterval(countdownTimer)
})
</script>

<style scoped>
.home-page {
  padding: 20px;
}

.main-layout {
  display: flex;
  gap: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.left-sidebar {
  width: 200px;
  flex-shrink: 0;
}

.right-sidebar {
  width: 280px;
  flex-shrink: 0;
}

.main-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.sidebar-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.header-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
}

.header-icon svg {
  width: 16px;
  height: 16px;
}

.header-icon.hot {
  background: linear-gradient(135deg, #FAAD14, #D48806);
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.rank-badge {
  background: linear-gradient(135deg, #FAAD14, #D48806);
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-item {
  display: flex;
  gap: 10px;
  padding: 10px;
  border-radius: 8px;
  background: #f9f9f9;
  transition: all 0.2s;
  cursor: pointer;
}

.announcement-item:hover {
  background: #f5f5f5;
}

.announcement-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.announcement-tag.hot {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
}

.announcement-tag.new {
  background: linear-gradient(135deg, #1677FF, #0052CC);
  color: #fff;
}

.announcement-tag.event {
  background: linear-gradient(135deg, #52C41A, #237804);
  color: #fff;
}

.announcement-tag.normal {
  background: #f0f0f0;
  color: #666;
}

.announcement-content {
  flex: 1;
}

.announcement-content h4 {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin: 0 0 4px;
  line-height: 1.3;
}

.announcement-content p {
  font-size: 12px;
  color: #666;
  margin: 0 0 4px;
  line-height: 1.4;
}

.announcement-time {
  font-size: 11px;
  color: #999;
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rank-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 8px;
  background: #f9f9f9;
  transition: all 0.2s;
  cursor: pointer;
}

.rank-item:hover {
  background: #f5f5f5;
  transform: translateX(4px);
}

.rank-number {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  color: #999;
  background: #f0f0f0;
}

.rank-number.top3 {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
}

.rank-product-image {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.rank-product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.rank-product-info {
  flex: 1;
  min-width: 0;
}

.rank-product-name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  margin: 0 0 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rank-product-price {
  display: flex;
  align-items: center;
  gap: 6px;
}

.current-price {
  font-size: 14px;
  font-weight: 700;
  color: #FF5000;
}

.original-price {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.rank-sales {
  font-size: 11px;
  color: #999;
}

.banner-section {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
}

.banner-carousel {
  position: relative;
  height: 300px;
}

.carousel-wrapper {
  position: relative;
  height: 100%;
}

.carousel-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  transition: opacity 0.5s ease;
}

.carousel-slide.active {
  opacity: 1;
}

.slide-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.slide-1 {
  background: linear-gradient(135deg, #FF5000, #FF3400);
}

.slide-2 {
  background: linear-gradient(135deg, #1677FF, #0052CC);
}

.slide-3 {
  background: linear-gradient(135deg, #52C41A, #237804);
}

.slide-text {
  text-align: center;
  color: #fff;
}

.slide-text h2 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 8px;
}

.slide-text p {
  font-size: 18px;
  margin: 0 0 24px;
  opacity: 0.9;
}

.slide-btn {
  display: inline-block;
  padding: 12px 32px;
  background: #fff;
  color: #FF5000;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.slide-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.carousel-indicators {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
}

.indicator {
  width: 24px;
  height: 4px;
  background: rgba(255,255,255,0.5);
  border: none;
  border-radius: 2px;
  cursor: pointer;
  transition: all 0.2s;
}

.indicator.active {
  background: #fff;
  width: 32px;
}

.carousel-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  background: rgba(255,255,255,0.3);
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.carousel-arrow:hover {
  background: rgba(255,255,255,0.5);
}

.carousel-arrow.prev {
  left: 20px;
}

.carousel-arrow.next {
  right: 20px;
}

.carousel-arrow svg {
  width: 24px;
  height: 24px;
}

.seckill-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.section-icon svg {
  width: 18px;
  height: 18px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.section-badge {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.seckill-countdown {
  display: flex;
  align-items: center;
  gap: 12px;
}

.countdown-label {
  font-size: 14px;
  color: #FF5000;
}

.countdown-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.cd-num {
  font-size: 20px;
  font-weight: 700;
  color: #FF5000;
  background: #FFF5F0;
  padding: 6px 10px;
  border-radius: 6px;
}

.cd-sep {
  font-size: 16px;
  color: #FF5000;
  font-weight: 700;
}

.header-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #FF5000;
  font-size: 14px;
  text-decoration: none;
}

.header-link svg {
  width: 16px;
  height: 16px;
}

.seckill-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  min-height: 200px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  grid-column: 1 / -1;
}

.empty-state svg {
  width: 48px;
  height: 48px;
  color: #ccc;
}

.empty-state p {
  font-size: 14px;
  color: #999;
  margin: 12px 0;
}

.empty-btn {
  padding: 8px 24px;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  text-decoration: none;
}

.category-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  background: #f5f5f5;
  transform: translateY(-2px);
}

.category-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.category-icon svg {
  width: 28px;
  height: 28px;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.category-count {
  font-size: 12px;
  color: #999;
}

.features-section {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 32px;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
}

.feature-icon {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #FFF5F0, #FFE8E0);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FF5000;
}

.feature-icon svg {
  width: 32px;
  height: 32px;
}

.feature-item h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.feature-item p {
  font-size: 12px;
  color: #999;
  margin: 0;
}

@media (max-width: 1400px) {
  .left-sidebar {
    width: 180px;
  }
  
  .right-sidebar {
    width: 240px;
  }
}

@media (max-width: 1200px) {
  .left-sidebar,
  .right-sidebar {
    display: none;
  }
  
  .seckill-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .category-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .banner-carousel {
    height: 200px;
  }
  
  .slide-text h2 {
    font-size: 24px;
  }
  
  .seckill-grid {
    grid-template-columns: 1fr;
  }
  
  .category-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>