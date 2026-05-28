<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">系统配置</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 站点基本设置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          站点设置
        </h3>
        <div class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">站点名称</label>
            <n-input v-model:value="siteConfig.site_name" placeholder="推理引擎" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">站点描述</label>
            <n-input v-model:value="siteConfig.site_description" type="textarea" :rows="2" placeholder="AI 驱动的开发者社区" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">开放注册</label>
            <n-switch v-model:value="siteConfig.allow_register" />
          </div>
          <n-button type="primary" :loading="savingSite" @click="saveSiteConfig">保存站点设置</n-button>
        </div>
      </div>

      <!-- 数据库配置（只读展示） -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" /></svg>
          数据库配置
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">MySQL</span>
            <n-tag type="success" size="small">已连接</n-tag>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Redis</span>
            <n-tag :type="redisConnected ? 'success' : 'warning'" size="small">{{ redisConnected ? '已连接' : '未配置' }}</n-tag>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Elasticsearch</span>
            <n-tag type="warning" size="small">可选</n-tag>
          </div>
          <div class="flex justify-between py-2">
            <span class="text-sm text-gray-500">MinIO</span>
            <n-tag type="warning" size="small">可选</n-tag>
          </div>
        </div>
      </div>

      <!-- Redis 配置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" /></svg>
          Redis 配置
        </h3>
        <div class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">Redis 地址</label>
            <n-input v-model:value="redisConfig.host" placeholder="127.0.0.1" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">端口</label>
            <n-input-number v-model:value="redisConfig.port" :min="1" :max="65535" class="w-full" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">密码</label>
            <n-input v-model:value="redisConfig.password" type="password" show-password-on="click" placeholder="无密码可留空" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">数据库编号</label>
            <n-input-number v-model:value="redisConfig.db" :min="0" :max="15" class="w-full" />
          </div>
          <div class="flex gap-3">
            <n-button type="primary" :loading="savingRedis" @click="saveRedisConfig">保存 Redis 配置</n-button>
            <n-button @click="testRedis" :loading="testingRedis">测试连接</n-button>
          </div>
          <n-tag v-if="redisTestResult" :type="redisTestResult.ok ? 'success' : 'error'" size="small">{{ redisTestResult.message }}</n-tag>
        </div>
      </div>

      <!-- Elasticsearch 配置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          Elasticsearch 配置
        </h3>
        <div class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">ES 地址</label>
            <n-input v-model:value="esConfig.host" placeholder="127.0.0.1" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">端口</label>
            <n-input-number v-model:value="esConfig.port" :min="1" :max="65535" class="w-full" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">用户名</label>
            <n-input v-model:value="esConfig.username" placeholder="可选" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">密码</label>
            <n-input v-model:value="esConfig.password" type="password" show-password-on="click" placeholder="可选" />
          </div>
          <n-button type="primary" :loading="savingES" @click="saveESConfig">保存 ES 配置</n-button>
        </div>
      </div>

      <!-- 文件上传配置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" /></svg>
          文件上传
        </h3>
        <div class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">存储方式</label>
            <n-select v-model:value="uploadConfig.storage_type" :options="[
              { label: '本地存储', value: 'local' },
              { label: 'MinIO 对象存储', value: 'minio' }
            ]" />
          </div>
          <div v-if="uploadConfig.storage_type === 'minio'">
            <label class="text-sm text-gray-500 mb-1 block">MinIO 地址</label>
            <n-input v-model:value="uploadConfig.minio_endpoint" placeholder="http://localhost:9000" />
          </div>
          <div v-if="uploadConfig.storage_type === 'minio'">
            <label class="text-sm text-gray-500 mb-1 block">Bucket</label>
            <n-input v-model:value="uploadConfig.minio_bucket" placeholder="inference-engine" />
          </div>
          <n-button type="primary" :loading="savingUpload" @click="saveUploadConfig">保存上传配置</n-button>
        </div>
      </div>

      <!-- 安全配置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
          安全配置
        </h3>
        <div class="space-y-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">JWT 过期时间（天）</label>
            <n-input-number v-model:value="securityConfig.jwt_expire_days" :min="1" :max="365" class="w-full" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">请求限流（每分钟）</label>
            <n-input-number v-model:value="securityConfig.rate_limit" :min="1" :max="10000" class="w-full" />
          </div>
          <n-button type="primary" :loading="savingSecurity" @click="saveSecurityConfig">保存安全配置</n-button>
        </div>
      </div>
    </div>

    <!-- 高级：原始配置表 -->
    <div class="mt-6">
      <n-collapse>
        <n-collapse-item title="高级：查看全部配置项" name="raw">
          <div class="bg-white rounded-xl border border-gray-100 p-4">
            <n-data-table :columns="rawColumns" :data="allConfigs" :bordered="false" :single-line="false" size="small" />
          </div>
        </n-collapse-item>
      </n-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { NButton, NInput, NTag, useMessage, NCollapse, NCollapseItem } from 'naive-ui'
import { get, put } from '@/api/request'

const message = useMessage()
const redisConnected = ref(false)

// Saving states
const savingSite = ref(false)
const savingRedis = ref(false)
const testingRedis = ref(false)
const savingES = ref(false)
const savingUpload = ref(false)
const savingSecurity = ref(false)
const redisTestResult = ref<any>(null)

// Config groups
const siteConfig = reactive({ site_name: '', site_description: '', allow_register: true })
const redisConfig = reactive({ host: '127.0.0.1', port: 6379, password: '', db: 0 })
const esConfig = reactive({ host: '127.0.0.1', port: 9200, username: '', password: '' })
const uploadConfig = reactive({ storage_type: 'local', minio_endpoint: '', minio_bucket: '' })
const securityConfig = reactive({ jwt_expire_days: 7, rate_limit: 60 })

// All configs for raw view
const allConfigs = ref<any[]>([])

const rawColumns = [
  { title: 'Key', key: 'key', width: 200 },
  { title: 'Value', key: 'value', minWidth: 200, ellipsis: { tooltip: true } },
  { title: '说明', key: 'desc', minWidth: 200 }
]

async function fetchAllConfigs() {
  try {
    const res = await get<any>('/admin/configs')
    const items = Array.isArray(res.data) ? res.data : (res.data?.items || [])
    allConfigs.value = items

    // Populate config groups from flat config list
    for (const item of items) {
      switch (item.key) {
        case 'site_name': siteConfig.site_name = item.value; break
        case 'site_description': siteConfig.site_description = item.value; break
        case 'allow_register': siteConfig.allow_register = item.value === 'true'; break
        case 'redis_host': redisConfig.host = item.value; break
        case 'redis_port': redisConfig.port = parseInt(item.value) || 6379; break
        case 'redis_password': redisConfig.password = item.value; break
        case 'redis_db': redisConfig.db = parseInt(item.value) || 0; break
        case 'es_host': esConfig.host = item.value; break
        case 'es_port': esConfig.port = parseInt(item.value) || 9200; break
        case 'es_username': esConfig.username = item.value; break
        case 'es_password': esConfig.password = item.value; break
        case 'storage_type': uploadConfig.storage_type = item.value; break
        case 'minio_endpoint': uploadConfig.minio_endpoint = item.value; break
        case 'minio_bucket': uploadConfig.minio_bucket = item.value; break
        case 'jwt_expire_days': securityConfig.jwt_expire_days = parseInt(item.value) || 7; break
        case 'rate_limit': securityConfig.rate_limit = parseInt(item.value) || 60; break
      }
    }

    // Check health
    try {
      const health = await get<any>('/health')
      redisConnected.value = health.data?.redis === true
    } catch {}
  } catch {}
}

async function saveConfigKV(key: string, value: string, desc: string) {
  // Find existing config or create new
  const existing = allConfigs.value.find(c => c.key === key)
  if (existing) {
    await put(`/admin/configs/${existing.id}`, { value })
  } else {
    // Create via the backend - use update endpoint with a special create-or-update
    await put('/admin/ai-config', { [key]: value })
  }
}

async function saveSiteConfig() {
  savingSite.value = true
  try {
    await saveConfigKV('site_name', siteConfig.site_name, '站点名称')
    await saveConfigKV('site_description', siteConfig.site_description, '站点描述')
    await saveConfigKV('allow_register', String(siteConfig.allow_register), '是否开放注册')
    message.success('站点设置已保存')
    fetchAllConfigs()
  } catch {} finally { savingSite.value = false }
}

async function saveRedisConfig() {
  savingRedis.value = true
  try {
    await saveConfigKV('redis_host', redisConfig.host, 'Redis 地址')
    await saveConfigKV('redis_port', String(redisConfig.port), 'Redis 端口')
    await saveConfigKV('redis_password', redisConfig.password, 'Redis 密码')
    await saveConfigKV('redis_db', String(redisConfig.db), 'Redis 数据库编号')
    message.success('Redis 配置已保存，重启服务后生效')
    fetchAllConfigs()
  } catch {} finally { savingRedis.value = false }
}

async function testRedis() {
  testingRedis.value = true
  redisTestResult.value = null
  try {
    const res = await get<any>('/health')
    redisTestResult.value = { ok: res.data?.redis === true, message: res.data?.redis ? 'Redis 连接正常' : 'Redis 未连接，请检查配置' }
  } catch { redisTestResult.value = { ok: false, message: '连接失败' } } finally { testingRedis.value = false }
}

async function saveESConfig() {
  savingES.value = true
  try {
    await saveConfigKV('es_host', esConfig.host, 'ES 地址')
    await saveConfigKV('es_port', String(esConfig.port), 'ES 端口')
    await saveConfigKV('es_username', esConfig.username, 'ES 用户名')
    await saveConfigKV('es_password', esConfig.password, 'ES 密码')
    message.success('ES 配置已保存，重启服务后生效')
    fetchAllConfigs()
  } catch {} finally { savingES.value = false }
}

async function saveUploadConfig() {
  savingUpload.value = true
  try {
    await saveConfigKV('storage_type', uploadConfig.storage_type, '文件存储方式')
    if (uploadConfig.storage_type === 'minio') {
      await saveConfigKV('minio_endpoint', uploadConfig.minio_endpoint, 'MinIO 地址')
      await saveConfigKV('minio_bucket', uploadConfig.minio_bucket, 'MinIO Bucket')
    }
    message.success('上传配置已保存')
    fetchAllConfigs()
  } catch {} finally { savingUpload.value = false }
}

async function saveSecurityConfig() {
  savingSecurity.value = true
  try {
    await saveConfigKV('jwt_expire_days', String(securityConfig.jwt_expire_days), 'JWT 过期天数')
    await saveConfigKV('rate_limit', String(securityConfig.rate_limit), '请求限流')
    message.success('安全配置已保存')
    fetchAllConfigs()
  } catch {} finally { savingSecurity.value = false }
}

onMounted(fetchAllConfigs)
</script>