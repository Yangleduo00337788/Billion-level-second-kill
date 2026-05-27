<template>
  <Layout>
    <div class="cart-page">
      <div class="page-header">
        <h1>我的购物车</h1>
        <span class="cart-count">共 {{ cartItems.length }} 件商品</span>
      </div>

      <div class="cart-content" v-if="cartItems.length > 0">
        <div class="cart-list">
          <div class="cart-header">
            <div class="col-check">
              <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
              <span>全选</span>
            </div>
            <div class="col-product">商品信息</div>
            <div class="col-price">单价</div>
            <div class="col-quantity">数量</div>
            <div class="col-total">小计</div>
            <div class="col-action">操作</div>
          </div>

          <div 
            v-for="item in cartItems"
            :key="item.id"
            class="cart-item"
            :class="{ selected: item.selected }"
          >
            <div class="col-check">
              <input type="checkbox" v-model="item.selected" @change="updateSelectAll" />
            </div>
            <div class="col-product">
              <div class="product-image">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="3" width="18" height="18" rx="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
              </div>
              <div class="product-info">
                <h3 class="product-name">{{ item.productName }}</h3>
                <p class="product-desc">
                  <span class="tag-seckill" v-if="item.isSeckill">秒杀</span>
                  {{ item.isSeckill ? '限时秒杀商品' : '普通商品' }}
                </p>
              </div>
            </div>
            <div class="col-price">
              <span class="price-current">¥{{ item.price }}</span>
              <span class="price-original" v-if="item.originalPrice">¥{{ item.originalPrice }}</span>
            </div>
            <div class="col-quantity">
              <div class="quantity-control">
                <button class="qty-btn" @click="decreaseQty(item)" :disabled="item.quantity <= 1">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                </button>
                <input type="number" v-model="item.quantity" class="qty-input" min="1" :max="item.maxQuantity" @change="fixQuantity(item)" />
                <button class="qty-btn" @click="increaseQty(item)" :disabled="item.quantity >= item.maxQuantity">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                  </svg>
                </button>
              </div>
            </div>
            <div class="col-total">
              <span class="total-price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
            </div>
            <div class="col-action">
              <button class="action-btn remove" @click="removeItem(item)">删除</button>
            </div>
          </div>
        </div>

        <div class="cart-footer">
          <div class="footer-left">
            <div class="select-all">
              <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
              <span>全选</span>
            </div>
            <button class="batch-remove" @click="removeSelected" :disabled="selectedCount === 0">
              删除选中 ({{ selectedCount }})
            </button>
          </div>
          <div class="footer-right">
            <div class="summary">
              <span class="selected-count">已选 {{ selectedCount }} 件</span>
              <div class="total-box">
                <span class="total-label">合计：</span>
                <span class="total-amount">¥{{ selectedTotal.toFixed(2) }}</span>
              </div>
            </div>
            <button class="checkout-btn" @click="handleCheckout" :disabled="selectedCount === 0">
              去结算
            </button>
          </div>
        </div>
      </div>

      <div class="cart-empty" v-else>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
          <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
        </svg>
        <p>购物车空空如也~</p>
        <router-link to="/seckill" class="go-shopping">去秒杀专区逛逛</router-link>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'

const router = useRouter()

const cartItems = ref([])
const selectAll = ref(false)

const selectedCount = computed(() => cartItems.value.filter(item => item.selected).length)

const selectedTotal = computed(() => {
  return cartItems.value
    .filter(item => item.selected)
    .reduce((sum, item) => sum + item.price * item.quantity, 0)
})

function loadCart() {
  cartItems.value = JSON.parse(localStorage.getItem('cart') || '[]')
  updateSelectAll()
}

function saveCart() {
  localStorage.setItem('cart', JSON.stringify(cartItems.value))
}

function toggleSelectAll() {
  cartItems.value.forEach(item => {
    item.selected = selectAll.value
  })
  saveCart()
}

function updateSelectAll() {
  selectAll.value = cartItems.value.length > 0 && cartItems.value.every(item => item.selected)
}

function decreaseQty(item) {
  if (item.quantity > 1) {
    item.quantity--
    saveCart()
  }
}

function increaseQty(item) {
  if (item.quantity < item.maxQuantity) {
    item.quantity++
    saveCart()
  }
}

function fixQuantity(item) {
  if (item.quantity < 1) item.quantity = 1
  if (item.quantity > item.maxQuantity) item.quantity = item.maxQuantity
  saveCart()
}

function removeItem(item) {
  ElMessageBox.confirm('确定要删除该商品吗？', '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    const index = cartItems.value.findIndex(i => i.id === item.id)
    if (index > -1) {
      cartItems.value.splice(index, 1)
      saveCart()
      ElMessage.success('已删除')
    }
  }).catch(() => {})
}

function removeSelected() {
  if (selectedCount.value === 0) return
  
  ElMessageBox.confirm(`确定要删除选中的 ${selectedCount.value} 件商品吗？`, '批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    cartItems.value = cartItems.value.filter(item => !item.selected)
    saveCart()
    ElMessage.success('已删除选中商品')
  }).catch(() => {})
}

function handleCheckout() {
  if (selectedCount.value === 0) {
    ElMessage.warning('请先选择要结算的商品')
    return
  }
  
  const selectedItems = cartItems.value.filter(item => item.selected)
  const seckillItems = selectedItems.filter(item => item.isSeckill)
  
  if (seckillItems.length > 0) {
    ElMessage.warning('秒杀商品需要单独下单，请前往秒杀专区抢购')
    return
  }
  
  ElMessage.success('正在跳转到结算页面...')
}

watch(cartItems, () => {
  saveCart()
}, { deep: true })

onMounted(() => {
  loadCart()
})
</script>

<style scoped>
.cart-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  padding: 20px;
  border-radius: 12px;
}

.page-header h1 {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.cart-count {
  font-size: 14px;
  color: #999;
}

.cart-content {
  background: #fff;
  border-radius: 12px;
}

.cart-list {
  padding: 0 20px;
}

.cart-header {
  display: flex;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #e5e5e5;
  font-size: 14px;
  color: #666;
}

.cart-item {
  display: flex;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid #e5e5e5;
  transition: background 0.2s;
}

.cart-item.selected {
  background: #FFF0E6;
}

.cart-item:last-child {
  border-bottom: none;
}

.col-check {
  width: 60px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.col-check input {
  width: 18px;
  height: 18px;
  accent-color: #FF5000;
}

.col-product {
  flex: 1;
  display: flex;
  gap: 16px;
}

.product-image {
  width: 80px;
  height: 80px;
  background: #f5f5f5;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-image svg {
  width: 32px;
  height: 32px;
  color: #ccc;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.product-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.product-desc {
  font-size: 12px;
  color: #999;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-seckill {
  background: #FF5000;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}

.col-price {
  width: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.price-current {
  font-size: 16px;
  font-weight: 600;
  color: #FF5000;
}

.price-original {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
}

.col-quantity {
  width: 120px;
  display: flex;
  justify-content: center;
}

.quantity-control {
  display: flex;
  align-items: center;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
}

.qty-btn {
  width: 32px;
  height: 32px;
  background: #f5f5f5;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #666;
}

.qty-btn:hover:not(:disabled) {
  background: #e5e5e5;
}

.qty-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.qty-btn svg {
  width: 14px;
  height: 14px;
}

.qty-input {
  width: 50px;
  height: 32px;
  border: none;
  border-left: 1px solid #e5e5e5;
  border-right: 1px solid #e5e5e5;
  text-align: center;
  font-size: 14px;
}

.col-total {
  width: 100px;
  display: flex;
  justify-content: center;
}

.total-price {
  font-size: 16px;
  font-weight: 600;
  color: #FF5000;
}

.col-action {
  width: 80px;
  display: flex;
  justify-content: center;
}

.action-btn {
  padding: 8px 16px;
  background: transparent;
  border: none;
  font-size: 14px;
  color: #666;
  cursor: pointer;
}

.action-btn:hover {
  color: #FF5000;
}

.cart-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: #f5f5f5;
  border-radius: 0 0 12px 12px;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.select-all {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.select-all input {
  width: 18px;
  height: 18px;
  accent-color: #FF5000;
}

.batch-remove {
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
}

.batch-remove:hover:not(:disabled) {
  border-color: #FF5000;
  color: #FF5000;
}

.batch-remove:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.summary {
  display: flex;
  align-items: center;
  gap: 16px;
}

.selected-count {
  font-size: 14px;
  color: #666;
}

.total-box {
  display: flex;
  align-items: baseline;
}

.total-label {
  font-size: 14px;
  color: #666;
}

.total-amount {
  font-size: 24px;
  font-weight: 700;
  color: #FF5000;
}

.checkout-btn {
  padding: 16px 40px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.checkout-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}

.checkout-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.cart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: #fff;
  border-radius: 12px;
}

.cart-empty svg {
  width: 80px;
  height: 80px;
  color: #ccc;
}

.cart-empty p {
  font-size: 16px;
  color: #999;
  margin: 16px 0;
}

.go-shopping {
  padding: 12px 32px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s;
}

.go-shopping:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}
</style>