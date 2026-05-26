<template>
  <Layout>
    <div class="product-detail-page" v-loading="loading">
      <div class="detail-container">
        <div class="product-gallery">
          <div class="gallery-main">
            <div class="main-image">
              <img v-if="product.mainImage" :src="product.mainImage" :alt="product.productName" />
              <div v-else class="image-placeholder">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
                <span>{{ product.productName?.slice(0, 6) || '商品图片' }}</span>
              </div>
            </div>
            <div class="gallery-thumbnails">
              <div class="thumbnail active">
                <div class="thumb-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                  </svg>
                </div>
              </div>
              <div class="thumbnail">
                <div class="thumb-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                  </svg>
                </div>
              </div>
              <div class="thumbnail">
                <div class="thumb-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="product-info">
          <div class="info-header">
            <div class="product-tags">
              <span class="tag-seckill" v-if="isSeckill">限时秒杀</span>
              <span class="tag-discount" v-if="discountPercent > 0">{{ discountPercent }}折</span>
              <span class="tag-hot" v-if="product.soldCount > 50">热卖</span>
            </div>
            <h1 class="product-title">{{ product.productName || '秒杀商品' }}</h1>
            <p class="product-subtitle">{{ product.title || '高品质商品 限时特惠' }}</p>
          </div>

          <div class="price-section" v-if="isSeckill">
            <div class="seckill-price-box">
              <div class="price-label">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
                </svg>
                <span>秒杀价</span>
              </div>
              <div class="price-value">
                <span class="currency">¥</span>
                <span class="amount">{{ product.seckillPrice }}</span>
              </div>
              <div class="price-original">
                <span>原价 ¥{{ product.originalPrice || product.price }}</span>
              </div>
            </div>
            
            <div class="seckill-countdown">
              <div class="countdown-header">
                <span class="countdown-label">{{ product.status === 1 ? '距结束' : '距开始' }}</span>
              </div>
              <div class="countdown-timer">
                <span class="cd-block">{{ countdownHours }}</span>
                <span class="cd-sep">:</span>
                <span class="cd-block">{{ countdownMinutes }}</span>
                <span class="cd-sep">:</span>
                <span class="cd-block">{{ countdownSeconds }}</span>
              </div>
            </div>
            
            <div class="seckill-progress">
              <div class="progress-info">
                <span class="progress-label">抢购进度</span>
                <span class="progress-percent">已抢 {{ soldPercent }}%</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: `${soldPercent}%` }"></div>
              </div>
              <div class="stock-info">
                <span>库存 {{ product.seckillStock }} 件</span>
                <span>限购 {{ product.personLimit || 1 }} 件/人</span>
              </div>
            </div>
          </div>

          <div class="price-section" v-else>
            <div class="normal-price-box">
              <div class="price-value">
                <span class="currency">¥</span>
                <span class="amount">{{ product.price }}</span>
              </div>
              <div class="price-original" v-if="product.originalPrice">
                <span>原价 ¥{{ product.originalPrice }}</span>
              </div>
            </div>
          </div>

          <div class="info-attrs">
            <div class="attr-row">
              <span class="attr-label">商品编号</span>
              <span class="attr-value">{{ product.id }}</span>
            </div>
            <div class="attr-row">
              <span class="attr-label">商品分类</span>
              <span class="attr-value">{{ categoryName }}</span>
            </div>
            <div class="attr-row">
              <span class="attr-label">发货时间</span>
              <span class="attr-value">付款后24小时内发货</span>
            </div>
            <div class="attr-row">
              <span class="attr-label">服务保障</span>
              <span class="attr-value">
                <span class="service-tag">正品保障</span>
                <span class="service-tag">七天退换</span>
                <span class="service-tag">极速发货</span>
              </span>
            </div>
          </div>

          <div class="quantity-section">
            <span class="quantity-label">购买数量</span>
            <div class="quantity-control">
              <button class="qty-btn minus" @click="quantity > 1 && quantity--" :disabled="quantity <= 1">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
              </button>
              <input type="number" v-model="quantity" class="qty-input" min="1" :max="maxQuantity" />
              <button class="qty-btn plus" @click="quantity < maxQuantity && quantity++" :disabled="quantity >= maxQuantity">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
              </button>
            </div>
            <span class="quantity-tip">剩余 {{ product.seckillStock || product.totalStock }} 件</span>
          </div>

          <div class="action-buttons">
            <button 
              class="btn-buy-now"
              :class="{ disabled: isDisabled }"
              :disabled="isDisabled"
              @click="handleBuyNow"
            >
              <span v-if="isSeckill">
                <span v-if="product.status === 0">即将开始</span>
                <span v-else-if="product.status === 2">已结束</span>
                <span v-else-if="product.seckillStock <= 0">已售罄</span>
                <span v-else>立即抢购</span>
              </span>
              <span v-else>立即购买</span>
            </button>
            <button class="btn-add-cart" v-if="!isSeckill" @click="handleAddCart">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
              </svg>
              加入购物车
            </button>
          </div>
        </div>
      </div>

      <div class="detail-tabs">
        <div class="tabs-header">
          <button 
            class="tab-btn"
            :class="{ active: activeTab === 'detail' }"
            @click="activeTab = 'detail'"
          >
            商品详情
          </button>
          <button 
            class="tab-btn"
            :class="{ active: activeTab === 'specs' }"
            @click="activeTab = 'specs'"
          >
            规格参数
          </button>
          <button 
            class="tab-btn"
            :class="{ active: activeTab === 'reviews' }"
            @click="activeTab = 'reviews'"
          >
            用户评价
          </button>
        </div>
        
        <div class="tabs-content">
          <div class="tab-panel" v-show="activeTab === 'detail'">
            <div class="detail-desc">
              <h3>商品介绍</h3>
              <p>{{ product.description || '这是一款高品质的商品，采用优质材料制作，具有出色的性能和耐用性。限时秒杀活动，价格优惠，数量有限，欢迎选购！' }}</p>
            </div>
            <div class="detail-images">
              <div class="detail-image-placeholder">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <rect x="3" y="3" width="18" height="18" rx="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
                <span>商品详情图 1</span>
              </div>
              <div class="detail-image-placeholder">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                  <rect x="3" y="3" width="18" height="18" rx="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
                <span>商品详情图 2</span>
              </div>
            </div>
          </div>
          
          <div class="tab-panel" v-show="activeTab === 'specs'">
            <div class="specs-table">
              <div class="spec-row">
                <span class="spec-name">品牌</span>
                <span class="spec-value">官方品牌</span>
              </div>
              <div class="spec-row">
                <span class="spec-name">型号</span>
                <span class="spec-value">{{ product.productName }}</span>
              </div>
              <div class="spec-row">
                <span class="spec-name">产地</span>
                <span class="spec-value">中国大陆</span>
              </div>
              <div class="spec-row">
                <span class="spec-name">颜色</span>
                <span class="spec-value">默认</span>
              </div>
              <div class="spec-row">
                <span class="spec-name">材质</span>
                <span class="spec-value">优质材料</span>
              </div>
              <div class="spec-row">
                <span class="spec-name">重量</span>
                <span class="spec-value">约500g</span>
              </div>
            </div>
          </div>
          
          <div class="tab-panel" v-show="activeTab === 'reviews'">
            <div class="reviews-summary">
              <div class="summary-score">
                <span class="score-value">4.8</span>
                <span class="score-label">好评率 98%</span>
              </div>
              <div class="summary-tags">
                <span class="review-tag">质量好(128)</span>
                <span class="review-tag">发货快(96)</span>
                <span class="review-tag">性价比高(85)</span>
                <span class="review-tag">包装精美(42)</span>
              </div>
            </div>
            <div class="reviews-list">
              <div class="review-item">
                <div class="review-header">
                  <span class="review-user">用户***123</span>
                  <span class="review-date">2024-01-15</span>
                </div>
                <div class="review-content">
                  <p>商品质量很好，发货速度也很快，非常满意！</p>
                </div>
              </div>
              <div class="review-item">
                <div class="review-header">
                  <span class="review-user">用户***456</span>
                  <span class="review-date">2024-01-12</span>
                </div>
                <div class="review-content">
                  <p>秒杀价格太优惠了，抢到了很开心，商品也很不错。</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <el-dialog 
        v-model="seckillDialogVisible" 
        title=""
        width="420px"
        :close-on-click-modal="false"
        class="seckill-result-dialog"
      >
        <div class="result-content">
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
          <div v-if="seckillResult.success" class="result-order">
            <span class="order-label">订单号</span>
            <span class="order-value">{{ seckillResult.orderNo }}</span>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <button class="dialog-btn secondary" @click="seckillDialogVisible = false">
              {{ seckillResult.success ? '稍后支付' : '继续抢购' }}
            </button>
            <button v-if="seckillResult.success" class="dialog-btn primary" @click="goPay">
              立即支付
            </button>
          </div>
        </template>
      </el-dialog>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'
import { getSeckillProduct, executeSeckill, getStock } from '../api/seckill'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const product = ref({})
const quantity = ref(1)
const activeTab = ref('detail')
const seckillDialogVisible = ref(false)
const seckillResult = ref({ success: false, message: '', orderNo: '' })

const countdownHours = ref('02')
const countdownMinutes = ref('00')
const countdownSeconds = ref('00')

let countdownTimer = null

const isSeckill = computed(() => route.query.type === 'seckill')

const discountPercent = computed(() => {
  if (!product.value.originalPrice) return 0
  const price = isSeckill.value ? product.value.seckillPrice : product.value.price
  const discount = Math.round((price / product.value.originalPrice) * 10)
  return discount < 10 ? discount : 0
})

const soldPercent = computed(() => {
  if (!product.value.totalStock || !product.value.seckillStock) return 50
  return Math.round(((product.value.totalStock - product.value.seckillStock) / product.value.totalStock) * 100)
})

const maxQuantity = computed(() => {
  if (isSeckill.value) {
    return Math.min(product.value.personLimit || 1, product.value.seckillStock || 1)
  }
  return product.value.totalStock || 99
})

const categoryName = computed(() => {
  const categories = ['手机数码', '电脑办公', '家用电器', '服装鞋包', '美妆护肤', '食品生鲜']
  return categories[product.value.categoryId - 1] || '其他'
})

const isDisabled = computed(() => {
  if (isSeckill.value) {
    return product.value.status !== 1 || product.value.seckillStock <= 0
  }
  return false
})

async function loadProduct() {
  loading.value = true
  const productId = route.params.id
  
  try {
    if (isSeckill.value) {
      const res = await getSeckillProduct(productId)
      product.value = res.data || {}
      
      try {
        const stockRes = await getStock(productId)
        product.value.seckillStock = stockRes.data
      } catch (e) {}
    } else {
      product.value = {
        id: productId,
        productName: '示例商品',
        price: 999,
        originalPrice: 1299,
        description: '这是一款高品质的商品',
        totalStock: 100,
        soldCount: 50
      }
    }
  } catch (error) {
    console.error('加载商品失败', error)
    ElMessage.error('商品不存在或已下架')
    router.push('/seckill')
  } finally {
    loading.value = false
  }
}

async function handleBuyNow() {
  if (isDisabled.value) return
  
  if (isSeckill.value) {
    try {
      await ElMessageBox.confirm(
        `确定要抢购「${product.value.productName}」吗？每位用户仅有一次秒杀资格！`,
        '抢购确认',
        {
          confirmButtonText: '立即抢购',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
    } catch { return }
    
    try {
      ElMessage.info('正在抢购中，请稍候...')
      const res = await executeSeckill(product.value.id)
      seckillResult.value = {
        success: res.data.success,
        message: res.data.message,
        orderNo: res.data.orderNo || ''
      }
      seckillDialogVisible.value = true
      if (res.data.success) {
        try {
          const stockRes = await getStock(product.value.id)
          product.value.seckillStock = stockRes.data
        } catch (e) {}
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
  } else {
    ElMessage.success('正在跳转到结算页面...')
    router.push('/checkout')
  }
}

function handleAddCart() {
  ElMessage.success(`已将 ${quantity.value} 件商品加入购物车`)
}

function goPay() {
  seckillDialogVisible.value = false
  router.push('/orders')
}

function updateCountdown() {
  if (!isSeckill.value || !product.value.endTime) return
  
  const now = new Date()
  const target = product.value.status === 1 
    ? new Date(product.value.endTime)
    : new Date(product.value.startTime)
  
  const diff = Math.max(0, Math.floor((target - now) / 1000))
  countdownHours.value = Math.floor(diff / 3600).toString().padStart(2, '0')
  countdownMinutes.value = Math.floor((diff % 3600) / 60).toString().padStart(2, '0')
  countdownSeconds.value = (diff % 60).toString().padStart(2, '0')
}

onMounted(() => {
  loadProduct()
  countdownTimer = setInterval(updateCountdown, 1000)
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})
</script>

<style scoped>
.product-detail-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-container {
  display: flex;
  gap: 40px;
  background: #fff;
  border-radius: 12px;
  padding: 32px;
}

.product-gallery {
  width: 400px;
}

.gallery-main {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.main-image {
  width: 400px;
  height: 400px;
  border-radius: 12px;
  overflow: hidden;
  background: #f8f8f8;
}

.main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f5f5, #e8e8e8);
}

.image-placeholder svg {
  width: 64px;
  height: 64px;
  color: #ccc;
}

.image-placeholder span {
  font-size: 14px;
  color: #999;
  margin-top: 12px;
}

.gallery-thumbnails {
  display: flex;
  gap: 8px;
}

.thumbnail {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease;
}

.thumbnail.active {
  border-color: #FF5000;
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}

.thumb-placeholder svg {
  width: 32px;
  height: 32px;
  color: #ccc;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-tags {
  display: flex;
  gap: 8px;
}

.tag-seckill {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.tag-discount {
  background: #FF5000;
  color: #fff;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
}

.tag-hot {
  background: #FFD700;
  color: #333;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.product-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
  line-height: 1.4;
}

.product-subtitle {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.price-section {
  background: linear-gradient(135deg, #FFF5F0, #FFE8E0);
  border-radius: 12px;
  padding: 24px;
}

.seckill-price-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.price-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #FF5000;
  font-size: 14px;
  font-weight: 600;
}

.price-label svg {
  width: 18px;
  height: 18px;
}

.price-value {
  display: flex;
  align-items: baseline;
}

.currency {
  font-size: 20px;
  color: #FF5000;
  font-weight: 600;
}

.amount {
  font-size: 36px;
  color: #FF5000;
  font-weight: 700;
}

.price-original {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
}

.seckill-countdown {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255,80,0,0.2);
}

.countdown-header {
  display: flex;
  flex-direction: column;
}

.countdown-label {
  font-size: 14px;
  color: #FF5000;
  font-weight: 500;
}

.countdown-timer {
  display: flex;
  align-items: center;
  gap: 4px;
}

.cd-block {
  font-size: 24px;
  font-weight: 700;
  color: #FF5000;
  background: #fff;
  padding: 8px 12px;
  border-radius: 6px;
}

.cd-sep {
  font-size: 20px;
  color: #FF5000;
  font-weight: 700;
}

.seckill-progress {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255,80,0,0.2);
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.progress-label {
  font-size: 14px;
  color: #666;
}

.progress-percent {
  font-size: 14px;
  color: #FF5000;
  font-weight: 500;
}

.progress-bar {
  height: 8px;
  background: rgba(255,80,0,0.2);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #FF5000, #FF3400);
  transition: width 0.3s ease;
}

.stock-info {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  font-size: 12px;
  color: #999;
}

.normal-price-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-attrs {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.attr-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.attr-label {
  width: 80px;
  font-size: 14px;
  color: #999;
}

.attr-value {
  font-size: 14px;
  color: #333;
  display: flex;
  gap: 8px;
}

.service-tag {
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
}

.quantity-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.quantity-label {
  font-size: 14px;
  color: #666;
}

.quantity-control {
  display: flex;
  align-items: center;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
}

.qty-btn {
  width: 40px;
  height: 40px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #666;
  transition: all 0.2s ease;
}

.qty-btn:hover:not(:disabled) {
  background: #f5f5f5;
}

.qty-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.qty-btn svg {
  width: 18px;
  height: 18px;
}

.qty-input {
  width: 60px;
  height: 40px;
  border: none;
  border-left: 1px solid #e5e5e5;
  border-right: 1px solid #e5e5e5;
  text-align: center;
  font-size: 14px;
}

.quantity-tip {
  font-size: 12px;
  color: #999;
}

.action-buttons {
  display: flex;
  gap: 16px;
}

.btn-buy-now {
  flex: 1;
  padding: 16px 32px;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border: none;
  border-radius: 12px;
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-buy-now:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(255,80,0,0.3);
}

.btn-buy-now.disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-add-cart {
  flex: 1;
  padding: 16px 32px;
  background: #fff;
  border: 2px solid #FF5000;
  border-radius: 12px;
  color: #FF5000;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.btn-add-cart svg {
  width: 20px;
  height: 20px;
}

.btn-add-cart:hover {
  background: #FFF5F0;
}

.detail-tabs {
  background: #fff;
  border-radius: 12px;
}

.tabs-header {
  display: flex;
  border-bottom: 1px solid #e5e5e5;
}

.tab-btn {
  padding: 16px 32px;
  background: transparent;
  border: none;
  font-size: 16px;
  color: #666;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
}

.tab-btn.active {
  color: #FF5000;
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: #FF5000;
}

.tabs-content {
  padding: 24px;
}

.tab-panel {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-desc h3 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px;
}

.detail-desc p {
  font-size: 14px;
  color: #666;
  line-height: 1.8;
}

.detail-images {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-image-placeholder {
  width: 100%;
  height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 8px;
}

.detail-image-placeholder svg {
  width: 48px;
  height: 48px;
  color: #ccc;
}

.detail-image-placeholder span {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}

.specs-table {
  display: flex;
  flex-direction: column;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
}

.spec-row {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #e5e5e5;
}

.spec-row:last-child {
  border-bottom: none;
}

.spec-name {
  width: 120px;
  font-size: 14px;
  color: #999;
  background: #f5f5f5;
  padding: 8px;
}

.spec-value {
  flex: 1;
  font-size: 14px;
  color: #333;
  padding: 8px;
}

.reviews-summary {
  display: flex;
  gap: 32px;
  padding: 24px;
  background: #f5f5f5;
  border-radius: 12px;
}

.summary-score {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.score-value {
  font-size: 36px;
  font-weight: 700;
  color: #FF5000;
}

.score-label {
  font-size: 12px;
  color: #999;
}

.summary-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.review-tag {
  padding: 8px 16px;
  background: #fff;
  border-radius: 8px;
  font-size: 13px;
  color: #666;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.review-item {
  padding: 16px;
  border-bottom: 1px solid #e5e5e5;
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.review-user {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.review-date {
  font-size: 12px;
  color: #999;
}

.review-content p {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.seckill-result-dialog .result-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
}

.result-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.result-icon.success {
  background: #E8F5E9;
  color: #4CAF50;
}

.result-icon.fail {
  background: #FFEBEE;
  color: #F44336;
}

.result-icon svg {
  width: 48px;
  height: 48px;
}

.result-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px;
}

.result-desc {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.result-order {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.order-label {
  font-size: 13px;
  color: #999;
}

.order-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.dialog-btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dialog-btn.secondary {
  background: #f5f5f5;
  border: none;
  color: #666;
}

.dialog-btn.primary {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border: none;
  color: #fff;
}

@media (max-width: 900px) {
  .detail-container {
    flex-direction: column;
  }
  
  .product-gallery {
    width: 100%;
  }
  
  .main-image {
    width: 100%;
    height: 300px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>