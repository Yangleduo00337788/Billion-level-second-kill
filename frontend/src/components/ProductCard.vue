<template>
  <div class="product-card" @click="handleClick">
    <div class="card-image">
      <img v-if="product.mainImage" :src="product.mainImage" :alt="product.productName" />
      <div v-else class="image-placeholder">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
          <circle cx="8.5" cy="8.5" r="1.5"/>
          <polyline points="21 15 16 10 5 21"/>
        </svg>
        <span class="placeholder-text">{{ product.productName?.slice(0, 4) || '商品' }}</span>
      </div>
      
      <div class="card-tags">
        <span class="tag-seckill" v-if="isSeckill && product.status === 1">秒杀</span>
        <span class="tag-discount" v-if="discountPercent > 0">{{ discountPercent }}折</span>
        <span class="tag-hot" v-if="product.soldCount > 50">热卖</span>
      </div>
      
      <div class="card-overlay" v-if="showOverlay">
        <span v-if="product.seckillStock <= 0">已售罄</span>
        <span v-else-if="product.status === 2">已结束</span>
        <span v-else-if="product.status === 0">即将开始</span>
      </div>
      
      <div class="card-progress" v-if="isSeckill && product.status === 1 && product.seckillStock > 0">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: `${soldPercent}%` }"></div>
        </div>
        <span class="progress-text">已抢{{ soldPercent }}%</span>
      </div>
    </div>
    
    <div class="card-body">
      <h3 class="card-title">{{ product.productName || `秒杀商品 #${product.id}` }}</h3>
      
      <div class="card-price">
        <div class="price-current">
          <span class="price-symbol">¥</span>
          <span class="price-value">{{ displayPrice }}</span>
        </div>
        <div class="price-original" v-if="product.originalPrice && product.originalPrice > displayPrice">
          <span>¥{{ product.originalPrice }}</span>
        </div>
      </div>
      
      <div class="card-meta" v-if="isSeckill">
        <span class="meta-stock">库存 {{ product.seckillStock }}</span>
        <span class="meta-time">{{ formatTime(product.startTime) }}-{{ formatTime(product.endTime) }}</span>
      </div>
      
      <div class="card-meta" v-else>
        <span class="meta-sold">已售 {{ product.soldCount || 0 }}</span>
        <span class="meta-shop">{{ product.shopName || '官方旗舰店' }}</span>
      </div>
      
      <div class="card-actions">
        <button 
          class="btn-buy"
          :class="{ disabled: isDisabled }"
          :disabled="isDisabled"
          @click.stop="handleBuy"
        >
          <span v-if="product.seckillStock <= 0 && isSeckill">已售罄</span>
          <span v-else-if="product.status === 2 && isSeckill">已结束</span>
          <span v-else-if="product.status === 0 && isSeckill">即将开始</span>
          <span v-else>{{ isSeckill ? '立即抢购' : '加入购物车' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  product: {
    type: Object,
    required: true
  },
  isSeckill: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['buy', 'click'])

const router = useRouter()

const displayPrice = computed(() => {
  return props.isSeckill ? props.product.seckillPrice : props.product.price
})

const discountPercent = computed(() => {
  if (!props.product.originalPrice || !displayPrice.value) return 0
  const discount = Math.round((displayPrice.value / props.product.originalPrice) * 10)
  return discount < 10 ? discount : 0
})

const soldPercent = computed(() => {
  if (!props.product.totalStock || !props.product.seckillStock) return 50
  return Math.round(((props.product.totalStock - props.product.seckillStock) / props.product.totalStock) * 100)
})

const showOverlay = computed(() => {
  return props.isSeckill && (props.product.status !== 1 || props.product.seckillStock <= 0)
})

const isDisabled = computed(() => {
  return props.isSeckill && (props.product.status !== 1 || props.product.seckillStock <= 0)
})

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

function handleClick() {
  if (props.isSeckill) {
    router.push(`/product/${props.product.id}?type=seckill`)
  } else {
    router.push(`/product/${props.product.id}`)
  }
  emit('click', props.product)
}

function handleBuy() {
  emit('buy', props.product)
}
</script>

<style scoped>
.product-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.12);
}

.card-image {
  position: relative;
  width: 100%;
  height: 200px;
  background: #f8f8f8;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.product-card:hover .card-image img {
  transform: scale(1.05);
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
  width: 48px;
  height: 48px;
  color: #ccc;
}

.placeholder-text {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}

.card-tags {
  position: absolute;
  top: 8px;
  left: 8px;
  display: flex;
  gap: 6px;
}

.tag-seckill {
  background: linear-gradient(135deg, #FF5000, #FF3400);
  color: #fff;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.tag-discount {
  background: #FF5000;
  color: #fff;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
}

.tag-hot {
  background: #FFD700;
  color: #333;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-overlay span {
  background: rgba(255,255,255,0.9);
  color: #FF5000;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
}

.card-progress {
  position: absolute;
  bottom: 8px;
  left: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-bar {
  flex: 1;
  height: 4px;
  background: rgba(255,255,255,0.3);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #FF5000, #FF3400);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 11px;
  color: #fff;
  background: rgba(0,0,0,0.5);
  padding: 2px 6px;
  border-radius: 4px;
}

.card-body {
  padding: 16px;
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  line-height: 1.4;
  margin: 0 0 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-price {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 8px;
}

.price-current {
  display: flex;
  align-items: baseline;
}

.price-symbol {
  font-size: 14px;
  color: #FF5000;
  font-weight: 600;
}

.price-value {
  font-size: 20px;
  color: #FF5000;
  font-weight: 700;
}

.price-original {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
  margin-bottom: 12px;
}

.meta-stock, .meta-sold {
  color: #666;
}

.meta-time {
  color: #FF5000;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.btn-buy {
  width: 100%;
  padding: 10px 16px;
  background: linear-gradient(135deg, #FF5000, #FF3400);
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-buy:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}

.btn-buy.disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>