<template>
  <Layout>
    <div class="orders-page">
      <div class="page-header">
        <h1>我的订单</h1>
        <p>查看和管理您的所有订单</p>
      </div>

      <div class="orders-tabs">
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'all' }"
          @click="activeTab = 'all'"
        >
          全部订单
          <span class="tab-count" v-if="orders.length">{{ orders.length }}</span>
        </button>
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'pending' }"
          @click="activeTab = 'pending'"
        >
          待付款
          <span class="tab-count" v-if="pendingOrders.length">{{ pendingOrders.length }}</span>
        </button>
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'paid' }"
          @click="activeTab = 'paid'"
        >
          已完成
        </button>
        <button 
          class="tab-btn"
          :class="{ active: activeTab === 'cancelled' }"
          @click="activeTab = 'cancelled'"
        >
          已取消
        </button>
      </div>

      <div class="orders-content" v-loading="loading">
        <div class="orders-list" v-if="filteredOrders.length > 0">
          <div 
            v-for="order in filteredOrders"
            :key="order.id"
            class="order-card"
          >
            <div class="order-header">
              <div class="header-left">
                <span class="order-date">{{ formatDateTime(order.createTime) }}</span>
                <span class="order-no">订单号：{{ order.orderNo }}</span>
              </div>
              <div class="header-right">
                <span class="order-status" :class="getStatusClass(order.status)">
                  {{ getStatusText(order.status) }}
                </span>
              </div>
            </div>

            <div class="order-body">
              <div class="order-product">
                <div class="product-image">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="3" width="18" height="18" rx="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
                <div class="product-info">
                  <h3 class="product-name">{{ order.productName }}</h3>
                  <p class="product-desc">秒杀商品 | 限时特惠</p>
                </div>
                <div class="product-price">
                  <span class="price-label">单价</span>
                  <span class="price-value">¥{{ order.productPrice }}</span>
                  <span class="quantity">x{{ order.quantity }}</span>
                </div>
              </div>
            </div>

            <div class="order-footer">
              <div class="footer-left">
                <span class="total-label">订单总额</span>
                <span class="total-value">¥{{ order.totalAmount }}</span>
              </div>
              <div class="footer-right">
                <div class="payment-countdown" v-if="order.status === 0 && countdownSeconds(order) > 0">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                  </svg>
                  <span :class="{ urgent: countdownSeconds(order) < 300 }">
                    剩余 {{ formatCountdown(order) }}
                  </span>
                </div>
                <div class="order-actions">
                  <button 
                    class="action-btn cancel"
                    v-if="order.status === 0"
                    @click="handleCancel(order)"
                  >
                    取消订单
                  </button>
                  <button 
                    class="action-btn pay"
                    v-if="order.status === 0"
                    @click="handlePay(order)"
                  >
                    立即付款
                  </button>
                  <button 
                    class="action-btn detail"
                    v-if="order.status !== 0"
                    @click="viewDetail(order)"
                  >
                    查看详情
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="empty-state" v-else>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <path d="M16 10a4 4 0 0 1-8 0"/>
          </svg>
          <p>暂无订单</p>
          <router-link to="/seckill" class="empty-btn">去抢购</router-link>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '../components/Layout.vue'
import { getMyOrders, payOrder, cancelOrder } from '../api/seckill'

const router = useRouter()

const loading = ref(false)
const orders = ref([])
const activeTab = ref('all')

let countdownTimer = null

const pendingOrders = computed(() => orders.value.filter(o => o.status === 0))

const filteredOrders = computed(() => {
  if (activeTab.value === 'all') return orders.value
  if (activeTab.value === 'pending') return orders.value.filter(o => o.status === 0)
  if (activeTab.value === 'paid') return orders.value.filter(o => o.status === 1)
  if (activeTab.value === 'cancelled') return orders.value.filter(o => o.status === 2)
  return orders.value
})

async function loadOrders() {
  loading.value = true
  try {
    const res = await getMyOrders()
    orders.value = res.data || []
  } catch (error) {
    console.error('加载订单失败', error)
  } finally {
    loading.value = false
  }
}

async function handlePay(order) {
  try {
    await ElMessageBox.confirm(
      `确认支付订单「${order.orderNo}」吗？支付金额 ¥${order.totalAmount}`,
      '支付确认',
      {
        confirmButtonText: '确认支付',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    await payOrder(order.orderNo)
    ElMessage.success('支付成功！')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '支付失败')
      await loadOrders()
    }
  }
}

async function handleCancel(order) {
  try {
    await ElMessageBox.confirm(
      '确认取消该订单吗？取消后将无法再次秒杀该商品！',
      '取消订单',
      {
        confirmButtonText: '确认取消',
        cancelButtonText: '返回',
        type: 'warning'
      }
    )
    await cancelOrder(order.orderNo)
    ElMessage.success('订单已取消')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '取消失败')
    }
  }
}

function viewDetail(order) {
  router.push(`/product/${order.seckillId}?type=seckill`)
}

function getStatusClass(status) {
  if (status === 0) return 'pending'
  if (status === 1) return 'paid'
  return 'cancelled'
}

function getStatusText(status) {
  if (status === 0) return '待付款'
  if (status === 1) return '已完成'
  return '已取消'
}

function formatDateTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getFullYear()}-${(date.getMonth()+1).toString().padStart(2,'0')}-${date.getDate().toString().padStart(2,'0')} ${date.getHours().toString().padStart(2,'0')}:${date.getMinutes().toString().padStart(2,'0')}`
}

function countdownSeconds(order) {
  if (!order.createTime || order.status !== 0) return 0
  const createTime = new Date(order.createTime).getTime()
  const expireTime = createTime + 15 * 60 * 1000
  return Math.max(0, Math.floor((expireTime - Date.now()) / 1000))
}

function formatCountdown(order) {
  const seconds = countdownSeconds(order)
  if (seconds <= 0) return '已超时'
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${minutes}分${secs.toString().padStart(2,'0')}秒`
}

onMounted(() => {
  loadOrders()
  countdownTimer = setInterval(() => {
    orders.value = [...orders.value]
  }, 1000)
})

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})
</script>

<style scoped>
.orders-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.page-header p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.orders-tabs {
  display: flex;
  gap: 8px;
  background: #fff;
  padding: 16px;
  border-radius: 12px;
}

.tab-btn {
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-btn.active {
  background: #FF5000;
  color: #fff;
}

.tab-count {
  background: rgba(255,80,0,0.1);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.tab-btn.active .tab-count {
  background: rgba(255,255,255,0.2);
}

.orders-content {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  min-height: 400px;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  overflow: hidden;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f5f5f5;
}

.header-left {
  display: flex;
  gap: 16px;
}

.order-date {
  font-size: 14px;
  color: #666;
}

.order-no {
  font-size: 14px;
  color: #999;
}

.order-status {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.order-status.pending {
  background: #FFF0E6;
  color: #FF5000;
}

.order-status.paid {
  background: #E8F5E9;
  color: #4CAF50;
}

.order-status.cancelled {
  background: #f5f5f5;
  color: #999;
}

.order-body {
  padding: 16px;
}

.order-product {
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
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.product-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.product-desc {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.product-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.price-label {
  font-size: 12px;
  color: #999;
}

.price-value {
  font-size: 16px;
  font-weight: 600;
  color: #FF5000;
}

.quantity {
  font-size: 12px;
  color: #999;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9f9f9;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.total-label {
  font-size: 14px;
  color: #666;
}

.total-value {
  font-size: 18px;
  font-weight: 700;
  color: #FF7490;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.payment-countdown {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #FF5000;
}

.payment-countdown svg {
  width: 16px;
  height: 16px;
}

.payment-countdown .urgent {
  color: #F44336;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn.cancel {
  background: transparent;
  border: 1px solid #e5e5e5;
  color: #666;
}

.action-btn.pay {
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border: none;
  color: #fff;
}

.action-btn.detail {
  background: transparent;
  border: 1px solid #FF5000;
  color: #FF5000;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
}

.empty-state svg {
  width: 64px;
  height: 64px;
  color: #ccc;
}

.empty-state p {
  font-size: 16px;
  color: #999;
  margin: 16px 0;
}

.empty-btn {
  padding: 12px 24px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s ease;
}

.empty-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}
</style>