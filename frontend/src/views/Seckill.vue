<template>
  <div class="seckill-container">
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-icon">⚡</span>
          <span class="logo-text">秒杀商城</span>
        </div>
        <div class="user-info">
          <el-button type="primary" plain size="small" @click="showOrders = true">
            我的订单
          </el-button>
          <span class="username">{{ userStore.username }}</span>
          <el-button type="danger" plain size="small" @click="handleLogout">退出</el-button>
        </div>
      </div>
    </header>
    
    <main class="main-content">
      <div class="page-title">
        <h2>限时秒杀</h2>
        <p>抢购正在进行中，手慢无！</p>
      </div>
      
      <div v-loading="loading" class="product-grid">
        <div 
          v-for="product in products" 
          :key="product.id" 
          class="product-card"
        >
          <div class="product-image">
            <div class="product-badge" v-if="product.status === 1">秒杀中</div>
            <div class="product-badge ended" v-else-if="product.status === 2">已结束</div>
            <div class="product-badge pending" v-else>即将开始</div>
            <img v-if="product.mainImage" :src="product.mainImage" class="product-img" alt="商品图片" />
            <div v-else class="product-placeholder">
              <span>商品图片</span>
            </div>
          </div>
          
          <div class="product-info">
            <h3 class="product-name">{{ product.productName || '秒杀商品 #' + product.id }}</h3>
            <div class="price-row">
              <span class="seckill-price">¥{{ product.seckillPrice }}</span>
              <span class="original-price">¥{{ product.originalPrice || (product.seckillPrice * 1.5).toFixed(2) }}</span>
            </div>
            <div class="stock-info">
              <span>商品库存: {{ product.totalStock || '-' }} 件</span>
              <span>秒杀库存: {{ product.seckillStock }} 件</span>
              <span>限购: {{ product.personLimit }} 件/人</span>
            </div>
            <div class="time-info">
              <el-icon><Clock /></el-icon>
              <span>{{ formatTime(product.startTime) }} - {{ formatTime(product.endTime) }}</span>
            </div>
          </div>
          
          <div class="product-action">
            <el-button 
              type="danger" 
              size="large"
              class="seckill-btn"
              :disabled="product.status !== 1 || product.seckillStock <= 0"
              @click="handleSeckill(product)"
            >
              <el-icon><Lightning /></el-icon>
              <span v-if="product.seckillStock <= 0">已售罄</span>
              <span v-else-if="product.status === 2">已结束</span>
              <span v-else-if="product.status === 0">即将开始</span>
              <span v-else>立即秒杀</span>
            </el-button>
          </div>
        </div>
        
        <el-empty v-if="!loading && products.length === 0" description="暂无秒杀商品" />
      </div>
    </main>
    
    <el-dialog 
      v-model="seckillDialogVisible" 
      title="秒杀结果" 
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="seckill-result">
        <el-icon v-if="seckillResult.success" class="success-icon"><CircleCheckFilled /></el-icon>
        <el-icon v-else class="fail-icon"><CircleCloseFilled /></el-icon>
        <p class="result-message">{{ seckillResult.message }}</p>
        <p v-if="seckillResult.success" class="order-no">订单号: {{ seckillResult.orderNo }}</p>
      </div>
      <template #footer>
        <el-button type="primary" @click="seckillDialogVisible = false">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog 
      v-model="showOrders" 
      title="我的订单" 
      width="700px"
    >
      <div v-loading="ordersLoading">
        <div style="margin-bottom: 15px; text-align: right;">
          <el-button type="primary" size="small" @click="loadOrders">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
        <div v-for="order in orders" :key="order.id" class="order-item">
          <div class="order-info">
            <div class="order-no">订单号: {{ order.orderNo }}</div>
            <div class="order-product">{{ order.productName }}</div>
            <div class="order-price">¥{{ order.totalAmount }}</div>
            <div class="order-time">{{ formatDateTime(order.createTime) }}</div>
          </div>
          <div class="order-status">
            <el-tag v-if="order.status === 0" type="warning">待支付</el-tag>
            <el-tag v-else-if="order.status === 1" type="success">已支付</el-tag>
            <el-tag v-else type="info">已取消</el-tag>
          </div>
          <div class="order-actions">
            <el-button 
              v-if="order.status === 0" 
              type="primary" 
              size="small"
              @click="handlePay(order.orderNo)"
            >
              立即支付
            </el-button>
            <el-button 
              v-if="order.status === 0" 
              type="danger" 
              size="small"
              plain
              @click="handleCancel(order.orderNo)"
            >
              取消订单
            </el-button>
          </div>
        </div>
        <el-empty v-if="orders.length === 0" description="暂无订单" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { getSeckillProducts, executeSeckill, getStock, getMyOrders, payOrder, cancelOrder } from '../api/seckill'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const products = ref([])
const seckillDialogVisible = ref(false)
const seckillResult = ref({
  success: false,
  message: '',
  orderNo: ''
})

const showOrders = ref(false)
const ordersLoading = ref(false)
const orders = ref([])

let stockTimer = null

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
    } catch (e) {
      console.error('刷新库存失败', e)
    }
  }
}

async function handleSeckill(product) {
  try {
    await ElMessageBox.confirm('确定要秒杀该商品吗？每位用户仅有一次秒杀资格，请谨慎下单！', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
  } catch {
    return
  }
  
  try {
    ElMessage.info('正在秒杀中...')
    const res = await executeSeckill(product.id)
    
    seckillResult.value = {
      success: res.data.success,
      message: res.data.message,
      orderNo: res.data.orderNo || ''
    }
    seckillDialogVisible.value = true
    
    if (res.data.success) {
      await refreshStock()
      setTimeout(() => {
        loadOrders()
      }, 1500)
    }
  } catch (error) {
    const errorMsg = error.message || '秒杀失败'
    if (errorMsg.includes('已参与过') || errorMsg.includes('已经参与')) {
      ElMessage.warning('您已参与过该商品的秒杀，无法重复下单！')
    } else {
      seckillResult.value = {
        success: false,
        message: errorMsg,
        orderNo: ''
      }
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
    ElMessage.success('订单已取消，您将无法再次秒杀该商品')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '取消失败')
    }
  }
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

function formatDateTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return `${date.getMonth()+1}/${date.getDate()} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadProducts()
  loadOrders()
  stockTimer = setInterval(refreshStock, 5000)
})

onUnmounted(() => {
  if (stockTimer) {
    clearInterval(stockTimer)
  }
})
</script>

<style scoped>
.seckill-container {
  min-height: 100vh;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
}

.header {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0 40px;
  height: 70px;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #f5576c, #f093fb);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: #fff;
  font-size: 14px;
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px;
}

.page-title {
  text-align: center;
  margin-bottom: 40px;
}

.page-title h2 {
  font-size: 36px;
  color: #fff;
  margin-bottom: 10px;
}

.page-title p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 16px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 30px;
}

.product-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.product-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(245, 87, 108, 0.3);
  border-color: rgba(245, 87, 108, 0.5);
}

.product-image {
  position: relative;
  height: 200px;
  background: linear-gradient(135deg, #2d3436 0%, #000000 100%);
}

.product-badge {
  position: absolute;
  top: 15px;
  left: 15px;
  background: linear-gradient(135deg, #f5576c, #f093fb);
  color: #fff;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.product-badge.ended {
  background: #666;
}

.product-badge.pending {
  background: linear-gradient(135deg, #11998e, #38ef7d);
}

.product-placeholder {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.3);
  font-size: 16px;
}

.product-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-info {
  padding: 20px;
}

.product-name {
  color: #fff;
  font-size: 18px;
  margin-bottom: 15px;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 15px;
}

.seckill-price {
  font-size: 28px;
  font-weight: 700;
  color: #f5576c;
}

.original-price {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  text-decoration: line-through;
}

.stock-info {
  display: flex;
  gap: 20px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
  margin-bottom: 10px;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
}

.product-action {
  padding: 0 20px 20px;
}

.seckill-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #f5576c, #f093fb);
  border: none;
}

.seckill-btn:hover:not(:disabled) {
  transform: scale(1.02);
}

.seckill-btn:disabled {
  background: #444;
}

.seckill-result {
  text-align: center;
  padding: 20px 0;
}

.success-icon {
  font-size: 60px;
  color: #67c23a;
  margin-bottom: 15px;
}

.fail-icon {
  font-size: 60px;
  color: #f56c6c;
  margin-bottom: 15px;
}

.result-message {
  font-size: 18px;
  color: #333;
  margin-bottom: 10px;
}

.order-no {
  color: #666;
  font-size: 14px;
}

.order-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #eee;
  gap: 15px;
}

.order-info {
  flex: 1;
}

.order-no {
  font-size: 12px;
  color: #999;
}

.order-product {
  font-size: 16px;
  font-weight: 600;
  margin: 5px 0;
}

.order-price {
  font-size: 18px;
  color: #f5576c;
  font-weight: 600;
}

.order-time {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.order-status {
  width: 80px;
  text-align: center;
}

.order-actions {
  display: flex;
  gap: 10px;
}
</style>
