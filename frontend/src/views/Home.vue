<template>
  <Layout>
    <div class="home-page">
      <div class="main-layout">
        <aside class="left-sidebar">
          <div class="sidebar-section">
            <h3 class="sidebar-title">商品分类</h3>
            <div class="menu-list">
              <div class="menu-item" v-for="(cat, index) in menuCategories" :key="index" @click="handleMenuClick(cat)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path :d="cat.icon"/>
                </svg>
                <span>{{ cat.name }}</span>
              </div>
            </div>
          </div>
        </aside>

        <div class="main-content">
          <section class="banner-section">
            <div class="banner-carousel" @mouseenter="pauseAutoPlay" @mouseleave="resumeAutoPlay">
              <div class="carousel-wrapper">
                <div class="carousel-slide" :class="{ active: currentSlide === 0 }">
                  <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=e-commerce%20fashion%20sale%20banner%20with%20summer%20clothing%20models%20colorful%20modern%20design&image_size=landscape_16_9" alt="618抢先购" class="slide-img" />
                  <div class="slide-overlay">
                    <div class="slide-tag">618抢先购</div>
                    <div class="slide-text">
                      <h2>服饰时尚</h2>
                      <p>叠券低至7.3折起</p>
                    </div>
                  </div>
                </div>
                <div class="carousel-slide" :class="{ active: currentSlide === 1 }">
                  <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=technology%20digital%20products%20sale%20banner%20smartphone%20laptop%20gadgets%20blue%20purple%20gradient&image_size=landscape_16_9" alt="数码潮品" class="slide-img" />
                  <div class="slide-overlay">
                    <div class="slide-tag">新品首发</div>
                    <div class="slide-text">
                      <h2>数码潮品</h2>
                      <p>最新科技抢先体验</p>
                    </div>
                  </div>
                </div>
                <div class="carousel-slide" :class="{ active: currentSlide === 2 }">
                  <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=quality%20guaranteed%20products%20banner%20premium%20goods%20shopping%20green%20fresh&image_size=landscape_16_9" alt="正品好货" class="slide-img" />
                  <div class="slide-overlay">
                    <div class="slide-tag">品质保障</div>
                    <div class="slide-text">
                      <h2>正品好货</h2>
                      <p>全场正品保证 七天无忧退换</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="carousel-indicators">
                <button v-for="(_, index) in 3" :key="index" class="indicator" :class="{ active: currentSlide === index }" @click="goToSlide(index)"></button>
              </div>
            </div>
          </section>

          <section class="seckill-section">
            <div class="section-header">
              <div class="header-left">
                <div class="section-icon">
                  <svg viewBox="0 0 24 24" fill="currentColor"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
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
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
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
        </div>

        <aside class="right-sidebar">
          <div class="activity-cards">
            <div class="activity-card card-red" @click="$router.push('/seckill')">
              <div class="card-text">
                <h4>618淘宝下单</h4>
                <p>赢华为平板</p>
              </div>
              <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=huawei%20tablet%20computer%20product%20icon%20minimal&image_size=square_hd" alt="" class="card-img" />
            </div>
            <div class="activity-card card-orange" @click="$router.push('/seckill')">
              <div class="card-text">
                <h4>iPhone17pro</h4>
                <p>直降千元 领300元券</p>
              </div>
              <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=smartphone%20iPhone%20product%20icon%20minimal%20white&image_size=square_hd" alt="" class="card-img" />
            </div>
            <div class="activity-card card-pink" @click="$router.push('/seckill')">
              <div class="card-text">
                <h4>酷暑来袭</h4>
                <p>升温爆款5折起</p>
              </div>
              <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=summer%20cooling%20products%20ice%20cream%20fan%20icon%20minimal&image_size=square_hd" alt="" class="card-img" />
            </div>
            <div class="activity-card card-yellow" @click="$router.push('/seckill')">
              <div class="card-text">
                <h4>运动户外</h4>
                <p>尖货爆款直降</p>
              </div>
              <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=sports%20outdoor%20running%20shoes%20product%20icon%20minimal&image_size=square_hd" alt="" class="card-img" />
            </div>
          </div>
        </aside>
      </div>

      <section class="recommend-section">
        <div class="recommend-header">
          <div class="header-icon">
            <svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>
          </div>
          <div class="header-text">
            <h3>猜你喜欢</h3>
            <p>精选好物推荐</p>
          </div>
        </div>
        <div class="recommend-grid">
          <div class="recommend-item" v-for="(product, index) in recommendProducts" :key="index" @click="goProductDetail(product)">
            <div class="rec-image">
              <img v-if="product.mainImage" :src="product.mainImage" :alt="product.productName" />
              <div v-else class="rec-img-placeholder">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
              </div>
            </div>
            <div class="rec-info">
              <h4 class="rec-title">{{ product.productName || `精选商品 #${product.id}` }}</h4>
              <div class="rec-price">
                <span class="rec-current">¥{{ product.seckillPrice || product.price }}</span>
                <span class="rec-original" v-if="product.originalPrice">¥{{ product.originalPrice }}</span>
              </div>
              <div class="rec-tags">
                <span class="rec-tag" v-if="product.soldCount > 100">热卖</span>
                <span class="rec-tag red" v-if="product.status === 1">秒杀中</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </Layout>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'
import ProductCard from '../components/ProductCard.vue'
import { getSeckillProducts, executeSeckill, getStock } from '../api/seckill'

const router = useRouter()

const loading = ref(false)
const seckillProducts = ref([])
const recommendProducts = ref([])
const currentSlide = ref(0)
const currentSeckillProduct = ref(null)

const countdownHours = ref('00')
const countdownMinutes = ref('00')
const countdownSeconds = ref('00')
const countdownLabel = ref('距结束')

let slideTimer = null
let countdownTimer = null

const menuCategories = [
  { name: '618现货抢先购 / 大牌5折起', icon: 'M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z', keyword: '618' },
  { name: '电脑 / 配件 / 办公 / 文具', icon: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z', keyword: '电脑' },
  { name: '工业品 / 商业 / 农业 / 定制', icon: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z', keyword: '工业品' },
  { name: '家电 / 手机 / 通信 / 数码', icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z', keyword: '手机' },
  { name: '家具 / 家装 / 家居 / 厨具', icon: 'M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z', keyword: '家具' },
  { name: '女装 / 男装 / 内衣 / 配饰', icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z', keyword: '女装' },
  { name: '女鞋 / 男鞋 / 运动 / 户外', icon: 'M22 12h-4l-3 9L9 3l-3 9H2', keyword: '运动鞋' },
  { name: '汽车 / 珠宝 / 文玩 / 箱包', icon: 'M5 18v3M10 18v3M15 18v3M20 18v3M3 12h18M5 12l2-7h10l2 7', keyword: '汽车' },
  { name: '食品 / 鲜花 / 酒水 / 健康', icon: 'M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 16a1 1 0 1 1 1-1 1 1 0 0 1-1 1zm1-5h-2V7h2z', keyword: '食品' },
  { name: '母婴 / 童装 / 玩具 / 宠物', icon: 'M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 16a1 1 0 1 1 1-1 1 1 0 0 1-1 1zm1-5h-2V7h2z', keyword: '母婴' },
  { name: '美妆 / 个护 / 娱乐 / 图书', icon: 'M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z', keyword: '美妆' },
]

async function loadProducts() {
  loading.value = true
  try {
    const res = await getSeckillProducts()
    seckillProducts.value = res.data || []
    recommendProducts.value = [...seckillProducts.value].sort(() => Math.random() - 0.5).slice(0, 10)
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

function handleMenuClick(cat) {
  if (cat.keyword) {
    router.push({ path: '/seckill', query: { keyword: cat.keyword } })
  }
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
  padding: 0;
}

.main-layout {
  display: flex;
  gap: 12px;
  max-width: 100%;
  margin: 12px auto 0;
  padding: 0 60px;
}

.left-sidebar {
  width: 240px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sidebar-section {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.sidebar-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid #FF5000;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: #666;
  transition: all 0.2s;
}

.menu-item:hover {
  background: #FFF0E6;
  color: #FF5000;
}

.menu-item svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.main-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.banner-section {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.banner-carousel {
  position: relative;
  height: 340px;
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

.slide-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.slide-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 40px 32px 24px;
  background: linear-gradient(transparent, rgba(0,0,0,0.6));
}

.slide-tag {
  display: inline-block;
  background: #FF5000;
  color: #fff;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
}

.slide-text {
  color: #fff;
}

.slide-text h2 {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 4px;
}

.slide-text p {
  font-size: 16px;
  margin: 0;
  opacity: 0.9;
}

.carousel-indicators {
  position: absolute;
  bottom: 12px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
  z-index: 2;
}

.indicator {
  width: 8px;
  height: 8px;
  background: rgba(255,255,255,0.5);
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
}

.indicator.active {
  background: #fff;
  width: 20px;
  border-radius: 4px;
}

.seckill-section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-icon {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.section-icon svg {
  width: 16px;
  height: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.section-badge {
  background: linear-gradient(135deg, #FF5000, #FF3300);
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.seckill-countdown {
  display: flex;
  align-items: center;
  gap: 8px;
}

.countdown-label {
  font-size: 13px;
  color: #FF5000;
}

.countdown-numbers {
  display: flex;
  align-items: center;
  gap: 3px;
}

.cd-num {
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  background: #FF5000;
  padding: 3px 6px;
  border-radius: 4px;
}

.cd-sep {
  font-size: 14px;
  color: #FF5000;
  font-weight: 700;
}

.header-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #999;
  font-size: 13px;
  text-decoration: none;
}

.header-link:hover {
  color: #FF5000;
}

.header-link svg {
  width: 14px;
  height: 14px;
}

.seckill-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
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
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  text-decoration: none;
}

.right-sidebar {
  width: 260px;
  flex-shrink: 0;
}

.activity-cards {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.activity-card {
  border-radius: 12px;
  padding: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s;
  min-height: 72px;
}

.activity-card:hover {
  transform: translateY(-2px);
}

.card-red {
  background: linear-gradient(135deg, #FF4D4F, #FF7875);
}

.card-orange {
  background: linear-gradient(135deg, #FF7A45, #FF9C6E);
}

.card-pink {
  background: linear-gradient(135deg, #FF7A33, #FF9A55);
}

.card-yellow {
  background: linear-gradient(135deg, #FAAD14, #FFC53D);
}

.card-text h4 {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px;
}

.card-text p {
  font-size: 12px;
  color: rgba(255,255,255,0.9);
  margin: 0;
}

.card-img {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  object-fit: cover;
  background: rgba(255,255,255,0.2);
}

.recommend-section {
  max-width: 100%;
  margin: 20px auto 0;
  padding: 0 60px;
}

.recommend-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.header-icon {
  width: 28px;
  height: 28px;
  color: #FF5000;
}

.header-icon svg {
  width: 100%;
  height: 100%;
}

.header-text h3 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.header-text p {
  font-size: 12px;
  color: #999;
  margin: 2px 0 0;
}

.recommend-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.recommend-item {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.recommend-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.rec-image {
  width: 100%;
  height: 180px;
  background: #f8f8f8;
  overflow: hidden;
}

.rec-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.rec-img-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ccc;
}

.rec-img-placeholder svg {
  width: 48px;
  height: 48px;
}

.rec-info {
  padding: 12px;
}

.rec-title {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  margin: 0 0 8px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.rec-price {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 6px;
}

.rec-current {
  font-size: 18px;
  font-weight: 700;
  color: #FF5000;
}

.rec-original {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.rec-tags {
  display: flex;
  gap: 6px;
}

.rec-tag {
  font-size: 11px;
  color: #FF5000;
  background: #FFF0E6;
  padding: 2px 6px;
  border-radius: 4px;
}

.rec-tag.red {
  background: #FF5000;
  color: #fff;
}

@media (max-width: 1200px) {
  .left-sidebar,
  .right-sidebar {
    display: none;
  }

  .seckill-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .recommend-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .banner-carousel {
    height: 200px;
  }

  .slide-text h2 {
    font-size: 24px;
  }

  .seckill-grid {
    grid-template-columns: 1fr;
  }

  .recommend-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
