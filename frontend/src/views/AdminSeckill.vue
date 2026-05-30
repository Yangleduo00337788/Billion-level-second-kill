<template>
  <div class="admin-seckill">
    <div class="page-header">
      <div class="header-left">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
        </div>
        <div class="header-text">
          <h1 class="header-title">秒杀引擎管理</h1>
          <p class="header-sub">管理秒杀商品、活动、订单与库存流水</p>
        </div>
      </div>
      <div class="header-right">
        <div class="realtime-indicator">
          <span class="indicator-dot"></span>
          <span>实时监控中</span>
        </div>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon-wrapper" style="--card-color: #4dabf7;">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">商品总数</span>
          <span class="stat-value">{{ stats.totalProducts }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon-wrapper" style="--card-color: #f59f00;">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">秒杀活动</span>
          <span class="stat-value">{{ stats.totalActivities }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon-wrapper" style="--card-color: #20c997;">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">今日订单</span>
          <span class="stat-value">{{ stats.todayOrders }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon-wrapper" style="--card-color: #ff6b6b;">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">QPS 限制</span>
          <span class="stat-value">{{ stats.qpsLimit }}</span>
        </div>
      </div>
    </div>

    <div class="tab-bar">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="tab-icon">
          <path :d="tab.icon"/>
        </svg>
        <span>{{ tab.label }}</span>
        <div class="tab-indicator" v-if="activeTab === tab.key"></div>
      </button>
    </div>

    <!-- Tab: 商品管理 -->
    <div class="tab-content" v-show="activeTab === 'products'">
      <div class="content-toolbar">
        <div class="toolbar-left">
          <span class="toolbar-title">商品列表</span>
          <span class="toolbar-count">{{ products.length }} 条记录</span>
        </div>
        <button class="btn-primary" @click="openProductDialog()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加商品
        </button>
      </div>

      <div class="table-wrapper">
        <table class="data-table" v-if="products.length > 0">
          <thead>
            <tr>
              <th>ID</th>
              <th>商品图片</th>
              <th>商品名称</th>
              <th>原价</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in products" :key="item.id">
              <td class="td-id">#{{ item.id }}</td>
              <td>
                <img
                  v-if="item.image"
                  :src="item.image"
                  class="product-thumb"
                  @error="onImgError"
                />
                <div v-else class="product-thumb-placeholder">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <polyline points="21 15 16 10 5 21"/>
                  </svg>
                </div>
              </td>
              <td class="td-name">{{ item.name || '--' }}</td>
              <td class="td-price">&yen;{{ (item.originalPrice || 0).toFixed(2) }}</td>
              <td>
                <span class="badge" :class="'badge-' + getProductStatusType(item.status)">
                  {{ getProductStatusLabel(item.status) }}
                </span>
              </td>
              <td class="td-time">{{ formatTime(item.createTime) }}</td>
              <td class="td-actions">
                <button class="btn-table btn-edit" @click="openProductDialog(item)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                  编辑
                </button>
                <button class="btn-table btn-delete" @click="handleDeleteProduct(item)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-else class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
              <line x1="12" y1="11" x2="12" y2="17"/>
              <line x1="9" y1="14" x2="15" y2="14"/>
            </svg>
          </div>
          <h3 class="empty-title">暂无商品数据</h3>
          <p class="empty-desc">点击"添加商品"按钮创建第一个秒杀商品</p>
          <button class="btn-primary empty-btn" @click="openProductDialog()">
            添加商品
          </button>
        </div>
      </div>
    </div>

    <!-- Tab: 秒杀活动 -->
    <div class="tab-content" v-show="activeTab === 'activities'">
      <div class="content-toolbar">
        <div class="toolbar-left">
          <span class="toolbar-title">活动列表</span>
          <span class="toolbar-count">{{ activities.length }} 条记录</span>
        </div>
        <button class="btn-primary" @click="openActivityDialog()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加活动
        </button>
      </div>

      <div class="table-wrapper">
        <table class="data-table" v-if="activities.length > 0">
          <thead>
            <tr>
              <th>ID</th>
              <th>活动名称</th>
              <th>关联商品</th>
              <th>秒杀价</th>
              <th>库存</th>
              <th>开始时间</th>
              <th>结束时间</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in activities" :key="item.id">
              <td class="td-id">#{{ item.id }}</td>
              <td class="td-name">{{ item.name || '--' }}</td>
              <td>{{ getProductName(item.productId) }}</td>
              <td class="td-price sec-price">&yen;{{ (item.seckillPrice || 0).toFixed(2) }}</td>
              <td class="td-stock">
                <span class="stock-num">{{ item.stock || 0 }}</span>
                <span class="stock-text">件</span>
              </td>
              <td class="td-time">{{ formatTime(item.startTime) }}</td>
              <td class="td-time">{{ formatTime(item.endTime) }}</td>
              <td>
                <span class="badge" :class="'badge-' + getActivityStatusType(item)">
                  {{ getActivityStatusLabel(item) }}
                </span>
              </td>
              <td class="td-actions">
                <button class="btn-table btn-edit" @click="openActivityDialog(item)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                  编辑
                </button>
                <button class="btn-table btn-delete" @click="handleDeleteActivity(item)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-else class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
              <line x1="10" y1="14" x2="14" y2="18"/>
              <line x1="14" y1="14" x2="10" y2="18"/>
            </svg>
          </div>
          <h3 class="empty-title">暂无秒杀活动</h3>
          <p class="empty-desc">点击"添加活动"创建秒杀活动，关联商品后即可上线</p>
          <button class="btn-primary empty-btn" @click="openActivityDialog()">
            添加活动
          </button>
        </div>
      </div>
    </div>

    <!-- Tab: 订单管理 -->
    <div class="tab-content" v-show="activeTab === 'orders'">
      <div class="content-toolbar">
        <div class="toolbar-left">
          <span class="toolbar-title">订单列表</span>
          <span class="toolbar-count">{{ orders.length }} 条记录</span>
        </div>
        <div class="toolbar-right">
          <button class="btn-ghost" @click="loadOrders">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
            刷新
          </button>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="data-table" v-if="orders.length > 0">
          <thead>
            <tr>
              <th>订单ID</th>
              <th>用户ID</th>
              <th>秒杀ID</th>
              <th>商品名称</th>
              <th>价格</th>
              <th>状态</th>
              <th>创建时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in orders" :key="item.id">
              <td class="td-id">#{{ item.id || item.orderId }}</td>
              <td>{{ item.userId || '--' }}</td>
              <td>{{ item.seckillId || '--' }}</td>
              <td class="td-name">{{ item.productName || getProductName(item.productId) || '--' }}</td>
              <td class="td-price">&yen;{{ (item.price || 0).toFixed(2) }}</td>
              <td>
                <span class="badge" :class="'badge-' + getOrderStatusType(item.status)">
                  {{ getOrderStatusLabel(item.status) }}
                </span>
              </td>
              <td class="td-time">{{ formatTime(item.createTime) }}</td>
            </tr>
          </tbody>
        </table>

        <div v-else class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="16" y1="13" x2="8" y2="13"/>
              <line x1="16" y1="17" x2="8" y2="17"/>
              <line x1="10" y1="9" x2="8" y2="11"/>
            </svg>
          </div>
          <h3 class="empty-title">暂无订单数据</h3>
          <p class="empty-desc">秒杀活动开始后，产生的订单将在此处显示</p>
        </div>
      </div>
    </div>

    <!-- Tab: 库存流水 -->
    <div class="tab-content" v-show="activeTab === 'stockLogs'">
      <div class="content-toolbar">
        <div class="toolbar-left">
          <span class="toolbar-title">库存流水</span>
          <span class="toolbar-count">{{ stockLogs.length }} 条记录</span>
        </div>
        <div class="toolbar-right">
          <button class="btn-ghost" @click="loadStockLogs">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
            刷新
          </button>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="data-table" v-if="stockLogs.length > 0">
          <thead>
            <tr>
              <th>ID</th>
              <th>秒杀ID</th>
              <th>关联商品</th>
              <th>操作类型</th>
              <th>变更量</th>
              <th>变更后库存</th>
              <th>Redis状态</th>
              <th>时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in stockLogs" :key="item.id">
              <td class="td-id">#{{ item.id }}</td>
              <td>{{ item.seckillId || '--' }}</td>
              <td class="td-name">{{ item.productName || getProductName(item.productId) || '--' }}</td>
              <td>
                <span class="op-tag" :class="'op-tag-' + getOpTypeClass(item.operationType)">
                  {{ getOpTypeLabel(item.operationType) }}
                </span>
              </td>
              <td>
                <span class="change-num" :class="getChangeClass(item.changeAmount)">
                  {{ formatChangeAmount(item.changeAmount) }}
                </span>
              </td>
              <td class="td-stock">
                <span class="stock-num">{{ item.stockAfter ?? '--' }}</span>
              </td>
              <td>
                <span class="badge" :class="'badge-' + getRedisStatusType(item.redisStatus)">
                  {{ getRedisStatusLabel(item.redisStatus) }}
                </span>
              </td>
              <td class="td-time">{{ formatTime(item.createTime || item.time) }}</td>
            </tr>
          </tbody>
        </table>

        <div v-else class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 12a9 9 0 1 1-9-9"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
              <line x1="8" y1="8" x2="12" y2="12"/>
              <line x1="16" y1="8" x2="12" y2="12"/>
            </svg>
          </div>
          <h3 class="empty-title">暂无库存流水</h3>
          <p class="empty-desc">秒杀活动中的库存变更记录将在此处显示</p>
        </div>
      </div>
    </div>

    <!-- 商品对话框 -->
    <el-dialog
      v-model="productDialogVisible"
      :title="editingProduct.id ? '编辑商品' : '添加商品'"
      width="560px"
      destroy-on-close
      class="admin-dialog"
    >
      <el-form :model="productForm" label-position="top" class="admin-form">
        <el-form-item label="商品名称">
          <el-input v-model="productForm.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="原价">
          <el-input-number v-model="productForm.originalPrice" :min="0" :precision="2" :step="10" style="width: 100%" placeholder="请输入原价" />
        </el-form-item>
        <el-form-item label="商品图片URL">
          <el-input v-model="productForm.image" placeholder="请输入商品图片URL" />
        </el-form-item>
        <el-form-item label="商品描述">
          <el-input v-model="productForm.description" type="textarea" :rows="3" placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="productForm.status" style="width: 100%">
            <el-option label="上架" :value="1" />
            <el-option label="下架" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <button class="btn-ghost" @click="productDialogVisible = false">取消</button>
          <button class="btn-primary" @click="handleSaveProduct" :disabled="saving">
            {{ saving ? '保存中...' : '确认保存' }}
          </button>
        </div>
      </template>
    </el-dialog>

    <!-- 活动对话框 -->
    <el-dialog
      v-model="activityDialogVisible"
      :title="editingActivity.id ? '编辑活动' : '添加活动'"
      width="560px"
      destroy-on-close
      class="admin-dialog"
    >
      <el-form :model="activityForm" label-position="top" class="admin-form">
        <el-form-item label="活动名称">
          <el-input v-model="activityForm.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="关联商品">
          <el-select v-model="activityForm.productId" style="width: 100%" placeholder="请选择关联商品">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="秒杀价格">
          <el-input-number v-model="activityForm.seckillPrice" :min="0" :precision="2" :step="10" style="width: 100%" placeholder="请输入秒杀价格" />
        </el-form-item>
        <el-form-item label="秒杀库存">
          <el-input-number v-model="activityForm.stock" :min="0" :step="1" style="width: 100%" placeholder="请输入库存数量" />
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="activityForm.startTime"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="activityForm.endTime"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="activityForm.status" style="width: 100%">
            <el-option label="待开始" :value="0" />
            <el-option label="进行中" :value="1" />
            <el-option label="已结束" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <button class="btn-ghost" @click="activityDialogVisible = false">取消</button>
          <button class="btn-primary" @click="handleSaveActivity" :disabled="saving">
            {{ saving ? '保存中...' : '确认保存' }}
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdminProducts,
  saveProduct,
  deleteProduct,
  getAdminActivities,
  saveActivity,
  deleteActivity,
  getAdminOrders,
  getStockLogs
} from '../api/admin'

const activeTab = ref('products')
const saving = ref(false)

const products = ref([])
const activities = ref([])
const orders = ref([])
const stockLogs = ref([])

const productDialogVisible = ref(false)
const productForm = ref({ name: '', originalPrice: 0, image: '', description: '', status: 1 })
const editingProduct = ref({})

const activityDialogVisible = ref(false)
const activityForm = ref({ name: '', productId: null, seckillPrice: 0, stock: 0, startTime: '', endTime: '', status: 0 })
const editingActivity = ref({})

const imgErrorSrc = new Set()

const stats = computed(() => {
  const today = new Date().toDateString()
  const todayOrders = orders.value.filter(o => {
    if (!o.createTime) return false
    return new Date(o.createTime).toDateString() === today
  }).length
  return {
    totalProducts: products.value.length,
    totalActivities: activities.value.length,
    todayOrders,
    qpsLimit: 10000
  }
})

const tabs = [
  {
    key: 'products',
    label: '商品管理',
    icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
  },
  {
    key: 'activities',
    label: '秒杀活动',
    icon: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z'
  },
  {
    key: 'orders',
    label: '订单管理',
    icon: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z M14 2v6h6 M16 13H8 M16 17H8 M10 9H8'
  },
  {
    key: 'stockLogs',
    label: '库存流水',
    icon: 'M21 12a9 9 0 1 1-9-9 M22 4l-10 10.01-3-3.01'
  }
]

onMounted(() => {
  loadProducts()
  loadActivities()
  loadOrders()
  loadStockLogs()
})

async function loadProducts() {
  try {
    const res = await getAdminProducts()
    products.value = res.data || []
  } catch (e) {
    ElMessage.error('加载商品列表失败')
  }
}

async function loadActivities() {
  try {
    const res = await getAdminActivities()
    activities.value = res.data || []
  } catch (e) {
    ElMessage.error('加载活动列表失败')
  }
}

async function loadOrders() {
  try {
    const res = await getAdminOrders()
    orders.value = res.data || []
  } catch (e) {
    ElMessage.error('加载订单列表失败')
  }
}

async function loadStockLogs() {
  try {
    const res = await getStockLogs()
    stockLogs.value = res.data || []
  } catch (e) {
    ElMessage.error('加载库存流水失败')
  }
}

function openProductDialog(item) {
  if (item) {
    editingProduct.value = item
    productForm.value = {
      name: item.name || '',
      originalPrice: item.originalPrice || 0,
      image: item.image || '',
      description: item.description || '',
      status: item.status !== undefined ? item.status : 1
    }
  } else {
    editingProduct.value = {}
    productForm.value = { name: '', originalPrice: 0, image: '', description: '', status: 1 }
  }
  productDialogVisible.value = true
}

async function handleSaveProduct() {
  if (!productForm.value.name) {
    ElMessage.warning('请输入商品名称')
    return
  }
  saving.value = true
  try {
    const data = { ...productForm.value }
    if (editingProduct.value.id) {
      data.id = editingProduct.value.id
    }
    await saveProduct(data)
    ElMessage.success(editingProduct.value.id ? '商品更新成功' : '商品添加成功')
    productDialogVisible.value = false
    await loadProducts()
  } catch (e) {
    ElMessage.error('保存商品失败')
  } finally {
    saving.value = false
  }
}

async function handleDeleteProduct(item) {
  try {
    await ElMessageBox.confirm(
      `确定要删除商品「${item.name}」吗？此操作不可恢复。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
    await deleteProduct(item.id)
    ElMessage.success('商品已删除')
    await loadProducts()
  } catch (e) {
    if (e !== 'cancel' && e?.toString() !== 'cancel') {
      ElMessage.error('删除商品失败')
    }
  }
}

function openActivityDialog(item) {
  if (item) {
    editingActivity.value = item
    activityForm.value = {
      name: item.name || '',
      productId: item.productId || null,
      seckillPrice: item.seckillPrice || 0,
      stock: item.stock || 0,
      startTime: item.startTime || '',
      endTime: item.endTime || '',
      status: item.status !== undefined ? item.status : 0
    }
  } else {
    editingActivity.value = {}
    activityForm.value = {
      name: '',
      productId: null,
      seckillPrice: 0,
      stock: 0,
      startTime: '',
      endTime: '',
      status: 0
    }
  }
  activityDialogVisible.value = true
}

async function handleSaveActivity() {
  if (!activityForm.value.name) {
    ElMessage.warning('请输入活动名称')
    return
  }
  if (!activityForm.value.productId) {
    ElMessage.warning('请选择关联商品')
    return
  }
  saving.value = true
  try {
    const data = { ...activityForm.value }
    if (editingActivity.value.id) {
      data.id = editingActivity.value.id
    }
    await saveActivity(data)
    ElMessage.success(editingActivity.value.id ? '活动更新成功' : '活动添加成功')
    activityDialogVisible.value = false
    await loadActivities()
  } catch (e) {
    ElMessage.error('保存活动失败')
  } finally {
    saving.value = false
  }
}

async function handleDeleteActivity(item) {
  try {
    await ElMessageBox.confirm(
      `确定要删除活动「${item.name}」吗？此操作不可恢复。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
    await deleteActivity(item.id)
    ElMessage.success('活动已删除')
    await loadActivities()
  } catch (e) {
    if (e !== 'cancel' && e?.toString() !== 'cancel') {
      ElMessage.error('删除活动失败')
    }
  }
}

function getProductName(productId) {
  if (!productId) return '--'
  const p = products.value.find(x => x.id === productId)
  return p ? p.name : `商品#${productId}`
}

function formatTime(time) {
  if (!time) return '--'
  try {
    const d = new Date(time)
    if (isNaN(d.getTime())) return '--'
    const pad = n => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
  } catch {
    return '--'
  }
}

function onImgError(e) {
  imgErrorSrc.add(e.target.src)
  e.target.style.display = 'none'
  if (e.target.parentElement) {
    const placeholder = e.target.parentElement.querySelector('.product-thumb-placeholder')
    if (!placeholder) {
      const div = document.createElement('div')
      div.className = 'product-thumb-placeholder'
      div.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>`
      e.target.parentElement.appendChild(div)
    }
  }
}

function getProductStatusType(status) {
  return status === 1 ? 'success' : 'default'
}

function getProductStatusLabel(status) {
  return status === 1 ? '上架' : '下架'
}

function getActivityStatusType(item) {
  if (item.status === 1) return 'warning'
  if (item.status === 2) return 'default'
  const now = Date.now()
  if (item.startTime && now < new Date(item.startTime).getTime()) return 'info'
  if (item.endTime && now > new Date(item.endTime).getTime()) return 'default'
  if (item.startTime && item.endTime &&
    now >= new Date(item.startTime).getTime() &&
    now <= new Date(item.endTime).getTime()) return 'success'
  return 'info'
}

function getActivityStatusLabel(item) {
  if (item.status === 1) return '进行中'
  if (item.status === 2) return '已结束'
  const now = Date.now()
  if (item.startTime && now < new Date(item.startTime).getTime()) return '待开始'
  if (item.endTime && now > new Date(item.endTime).getTime()) return '已过期'
  if (item.startTime && item.endTime &&
    now >= new Date(item.startTime).getTime() &&
    now <= new Date(item.endTime).getTime()) return '进行中'
  return '待开始'
}

function getOrderStatusType(status) {
  const map = { 0: 'info', 1: 'success', 2: 'warning', 3: 'default', 4: 'danger' }
  return map[status] || 'info'
}

function getOrderStatusLabel(status) {
  const map = { 0: '待支付', 1: '已支付', 2: '已发货', 3: '已完成', 4: '已取消' }
  return map[status] || '未知'
}

function getOpTypeLabel(type) {
  const map = { 0: '初始化', 1: '扣减', 2: '回滚', 3: '增加' }
  return map[type] !== undefined ? map[type] : (type || '--')
}

function getOpTypeClass(type) {
  const map = { 0: 'info', 1: 'danger', 2: 'warning', 3: 'success' }
  return map[type] || 'info'
}

function getChangeClass(amount) {
  if (amount === undefined || amount === null) return ''
  if (amount < 0) return 'negative'
  if (amount > 0) return 'positive'
  return 'zero'
}

function formatChangeAmount(amount) {
  if (amount === undefined || amount === null) return '--'
  if (amount > 0) return `+${amount}`
  return String(amount)
}

function getRedisStatusType(status) {
  const map = { 0: 'default', 1: 'success', 2: 'danger' }
  return map[status] || 'default'
}

function getRedisStatusLabel(status) {
  const map = { 0: '未同步', 1: '已同步', 2: '同步失败' }
  return map[status] || '未知'
}
</script>

<style scoped>
.admin-seckill {
  min-height: 100vh;
  padding-bottom: 60px;
}

/* ---- Page Header ---- */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(245,159,0,0.18), rgba(245,159,0,0.04));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #f59f00;
  flex-shrink: 0;
}

.header-icon svg {
  width: 22px;
  height: 22px;
}

.header-title {
  font-size: 22px;
  font-weight: 700;
  color: #e7e9ed;
  letter-spacing: -0.3px;
  margin: 0;
}

.header-sub {
  font-size: 13px;
  color: #8b93a5;
  margin: 4px 0 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.realtime-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(32,201,151,0.08);
  border: 1px solid rgba(32,201,151,0.15);
  color: #20c997;
  font-size: 13px;
  font-weight: 500;
}

.indicator-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #20c997;
  animation: pulse-dot 2s infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* ---- Stats Grid ---- */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 28px;
}

.stat-card {
  background: rgba(22,27,38,0.72);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 22px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s;
}

.stat-card:hover {
  border-color: rgba(255,255,255,0.12);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}

.stat-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: color-mix(in srgb, var(--card-color) 12%, transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--card-color);
  flex-shrink: 0;
}

.stat-icon-wrapper svg {
  width: 24px;
  height: 24px;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12.5px;
  color: #8b93a5;
  font-weight: 500;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  color: #e7e9ed;
  letter-spacing: -1px;
  line-height: 1;
}

/* ---- Tab Bar ---- */
.tab-bar {
  display: flex;
  gap: 4px;
  margin-bottom: 24px;
  background: rgba(22,27,38,0.48);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  padding: 4px;
}

.tab-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 11px 20px;
  border-radius: 9px;
  border: none;
  background: transparent;
  color: #8b93a5;
  font-size: 13.5px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.25s;
  position: relative;
}

.tab-item:hover {
  color: #e7e9ed;
  background: rgba(255,255,255,0.03);
}

.tab-item.active {
  color: #f59f00;
  background: rgba(245,159,0,0.08);
}

.tab-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.tab-indicator {
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 3px;
  border-radius: 2px;
  background: #f59f00;
}

/* ---- Toolbar ---- */
.content-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.toolbar-title {
  font-size: 15px;
  font-weight: 600;
  color: #e7e9ed;
}

.toolbar-count {
  font-size: 12.5px;
  color: #8b93a5;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ---- Buttons ---- */
.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 18px;
  border-radius: 8px;
  border: none;
  background: linear-gradient(135deg, #f59f00, #e68a00);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  filter: brightness(1.1);
  box-shadow: 0 4px 16px rgba(245,159,0,0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  filter: none;
  box-shadow: none;
}

.btn-primary svg {
  width: 16px;
  height: 16px;
}

.btn-ghost {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 18px;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.04);
  color: #8b93a5;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-ghost:hover {
  color: #e7e9ed;
  border-color: rgba(255,255,255,0.2);
  background: rgba(255,255,255,0.07);
}

.btn-ghost svg {
  width: 16px;
  height: 16px;
}

/* ---- Table ---- */
.table-wrapper {
  background: rgba(22,27,38,0.72);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  border-bottom: 1px solid rgba(255,255,255,0.08);
}

.data-table th {
  padding: 14px 18px;
  text-align: left;
  font-size: 11.5px;
  font-weight: 700;
  color: #595f70;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  background: rgba(255,255,255,0.02);
  white-space: nowrap;
}

.data-table td {
  padding: 14px 18px;
  font-size: 13.5px;
  color: #8b93a5;
  border-bottom: 1px solid rgba(255,255,255,0.04);
  vertical-align: middle;
}

.data-table tbody tr {
  transition: background 0.15s;
}

.data-table tbody tr:hover {
  background: rgba(255,255,255,0.02);
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}

.td-id {
  font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  font-size: 12.5px;
  color: #595f70;
}

.td-name {
  color: #e7e9ed;
  font-weight: 500;
}

.td-price {
  font-weight: 600;
  color: #ff6b6b;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.sec-price {
  color: #f59f00;
}

.td-stock {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.stock-num {
  font-weight: 700;
  color: #e7e9ed;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.stock-text {
  font-size: 11px;
  color: #595f70;
}

.td-time {
  font-size: 12.5px;
  color: #8b93a5;
  white-space: nowrap;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.td-actions {
  display: flex;
  gap: 8px;
}

.btn-table {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid rgba(255,255,255,0.08);
  background: rgba(255,255,255,0.03);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-table svg {
  width: 14px;
  height: 14px;
}

.btn-edit {
  color: #4dabf7;
}

.btn-edit:hover {
  background: rgba(77,171,247,0.12);
  border-color: rgba(77,171,247,0.25);
}

.btn-delete {
  color: #ff6b6b;
}

.btn-delete:hover {
  background: rgba(255,107,107,0.12);
  border-color: rgba(255,107,107,0.25);
}

/* ---- Product Thumb ---- */
.product-thumb {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid rgba(255,255,255,0.06);
  background: rgba(255,255,255,0.03);
}

.product-thumb-placeholder {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #595f70;
}

.product-thumb-placeholder svg {
  width: 20px;
  height: 20px;
}

/* ---- Badges ---- */
.badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 11.5px;
  font-weight: 600;
  letter-spacing: 0.2px;
}

.badge-success {
  background: rgba(32,201,151,0.12);
  color: #20c997;
}

.badge-warning {
  background: rgba(245,159,0,0.12);
  color: #f59f00;
}

.badge-danger {
  background: rgba(255,107,107,0.12);
  color: #ff6b6b;
}

.badge-info {
  background: rgba(77,171,247,0.12);
  color: #4dabf7;
}

.badge-default {
  background: rgba(255,255,255,0.06);
  color: #595f70;
}

/* ---- Operation Tags ---- */
.op-tag {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11.5px;
  font-weight: 600;
}

.op-tag-danger {
  background: rgba(255,107,107,0.12);
  color: #ff6b6b;
}

.op-tag-warning {
  background: rgba(245,159,0,0.12);
  color: #f59f00;
}

.op-tag-success {
  background: rgba(32,201,151,0.12);
  color: #20c997;
}

.op-tag-info {
  background: rgba(77,171,247,0.12);
  color: #4dabf7;
}

.change-num {
  font-weight: 600;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.change-num.positive {
  color: #20c997;
}

.change-num.negative {
  color: #ff6b6b;
}

.change-num.zero {
  color: #595f70;
}

/* ---- Empty State ---- */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 40px;
}

.empty-icon {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  color: #595f70;
}

.empty-icon svg {
  width: 36px;
  height: 36px;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: #8b93a5;
  margin: 0 0 8px;
}

.empty-desc {
  font-size: 13px;
  color: #595f70;
  margin: 0 0 24px;
}

.empty-btn {
  padding: 10px 24px;
  font-size: 13.5px;
}

/* ---- Dialog ---- */
:deep(.admin-dialog) {
  --el-dialog-bg-color: #13161d;
  --el-dialog-border-color: rgba(255,255,255,0.08);
  --el-text-color-primary: #e7e9ed;
  --el-text-color-regular: #8b93a5;
  --el-border-color: rgba(255,255,255,0.1);
  --el-border-color-light: rgba(255,255,255,0.08);
  --el-bg-color: #1a1e28;
  --el-fill-color-blank: #1a1e28;
  --el-color-primary: #f59f00;
  --el-color-primary-light-3: #e68a00;
  --el-input-bg-color: #1a1e28;
  --el-input-border-color: rgba(255,255,255,0.1);
  --el-input-hover-border-color: rgba(255,255,255,0.2);
  --el-input-focus-border-color: #f59f00;
  --el-input-text-color: #e7e9ed;
  --el-input-placeholder-color: #595f70;
  border-radius: 14px;
}

:deep(.admin-dialog .el-dialog__header) {
  padding: 24px 28px 0;
  border-bottom: none;
}

:deep(.admin-dialog .el-dialog__title) {
  font-size: 17px;
  font-weight: 700;
  color: #e7e9ed;
}

:deep(.admin-dialog .el-dialog__body) {
  padding: 20px 28px;
}

:deep(.admin-dialog .el-dialog__footer) {
  padding: 0 28px 24px;
}

:deep(.admin-form .el-form-item__label) {
  color: #8b93a5;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 6px;
}

:deep(.admin-form .el-input__wrapper) {
  background: #1a1e28;
  border-color: rgba(255,255,255,0.1);
  border-radius: 8px;
  box-shadow: none;
  padding: 2px 12px;
}

:deep(.admin-form .el-input__wrapper:hover) {
  border-color: rgba(255,255,255,0.2);
}

:deep(.admin-form .el-input__wrapper.is-focus) {
  border-color: #f59f00;
  box-shadow: 0 0 0 1px rgba(245,159,0,0.2);
}

:deep(.admin-form .el-input__inner) {
  color: #e7e9ed;
  font-size: 13.5px;
  height: 38px;
  line-height: 38px;
}

:deep(.admin-form .el-input__inner::placeholder) {
  color: #595f70;
}

:deep(.admin-form .el-textarea__inner) {
  background: #1a1e28;
  border-color: rgba(255,255,255,0.1);
  border-radius: 8px;
  color: #e7e9ed;
  font-size: 13.5px;
  padding: 10px 12px;
  resize: none;
}

:deep(.admin-form .el-textarea__inner:hover) {
  border-color: rgba(255,255,255,0.2);
}

:deep(.admin-form .el-textarea__inner:focus) {
  border-color: #f59f00;
  box-shadow: 0 0 0 1px rgba(245,159,0,0.2);
}

:deep(.admin-form .el-textarea__inner::placeholder) {
  color: #595f70;
}

:deep(.admin-form .el-select .el-input__wrapper) {
  background: #1a1e28;
}

:deep(.admin-form .el-input-number) {
  width: 100%;
}

:deep(.admin-form .el-input-number .el-input__wrapper) {
  background: #1a1e28;
}

:deep(.admin-form .el-input-number__decrease),
:deep(.admin-form .el-input-number__increase) {
  background: rgba(255,255,255,0.04);
  color: #8b93a5;
  border-color: rgba(255,255,255,0.08);
}

:deep(.admin-form .el-input-number__decrease:hover),
:deep(.admin-form .el-input-number__increase:hover) {
  color: #f59f00;
}

:deep(.admin-dialog .el-date-editor .el-input__wrapper) {
  background: #1a1e28;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* ---- MessageBox Overrides ---- */
:deep(.el-message-box) {
  --el-bg-color: #13161d;
  --el-border-color-lighter: rgba(255,255,255,0.08);
  --el-text-color-primary: #e7e9ed;
  --el-text-color-regular: #8b93a5;
}

:deep(.el-message-box__title) {
  color: #e7e9ed;
}

:deep(.el-message-box__message) {
  color: #8b93a5;
}

/* ---- Select Dropdown Overrides ---- */
:deep(.el-select-dropdown) {
  background: #1a1e28;
  border: 1px solid rgba(255,255,255,0.08);
}

:deep(.el-select-dropdown__item) {
  color: #8b93a5;
}

:deep(.el-select-dropdown__item.hover) {
  background: rgba(255,255,255,0.04);
}

:deep(.el-select-dropdown__item.selected) {
  color: #f59f00;
  font-weight: 600;
}

/* ---- Popper Overrides ---- */
:deep(.el-popper.is-light) {
  background: #1a1e28;
  border: 1px solid rgba(255,255,255,0.08);
  color: #8b93a5;
}

/* ---- Responsive ---- */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .tab-bar {
    flex-wrap: wrap;
  }

  .tab-item {
    flex: 0 0 calc(50% - 2px);
  }

  .data-table {
    font-size: 12px;
  }

  .td-actions {
    flex-direction: column;
    gap: 4px;
  }
}
</style>