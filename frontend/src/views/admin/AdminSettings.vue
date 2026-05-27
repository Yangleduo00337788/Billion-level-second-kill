<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">系统设置</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- AI 服务配置 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6 lg:col-span-2">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
          AI 服务配置
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label class="text-sm text-gray-500 mb-1 block">AI 服务商</label>
            <n-select v-model:value="aiConfig.ai_provider" :options="providerOptions" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">API Key</label>
            <n-input v-model:value="aiConfig.ai_api_key" type="password" show-password-on="click" placeholder="sk-xxxx" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">模型名称</label>
            <n-input v-model:value="aiConfig.ai_model" placeholder="gpt-4o" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">API 基础地址</label>
            <n-input v-model:value="aiConfig.ai_base_url" :placeholder="aiConfig.ai_provider === 'deepseek' ? 'https://api.deepseek.com/v1' : 'https://api.openai.com/v1'" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">Max Tokens</label>
            <n-input-number v-model:value="aiConfig.ai_max_tokens" :min="1" :max="128000" class="w-full" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">Temperature</label>
            <n-input-number v-model:value="aiConfig.ai_temperature" :min="0" :max="2" :step="0.1" class="w-full" />
          </div>
          <div>
            <label class="text-sm text-gray-500 mb-1 block">每分钟请求限制</label>
            <n-input-number v-model:value="aiConfig.ai_rate_limit" :min="1" :max="1000" class="w-full" />
          </div>
        </div>
        <div class="flex items-center gap-3 mt-4">
          <n-button type="primary" :loading="savingAI" @click="saveAIConfig">
            保存 AI 配置
          </n-button>
          <n-button @click="testAI" :loading="testingAI">
            测试连接
          </n-button>
          <n-tag v-if="aiTestResult" :type="aiTestResult.connected ? 'success' : 'error'" size="small">
            {{ aiTestResult.message }}
          </n-tag>
        </div>
      </div>

      <!-- 基本信息 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          基本信息
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">系统名称</span>
            <span class="text-sm font-medium">推理引擎</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">版本</span>
            <span class="text-sm font-medium">v1.0.0</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Go 版本</span>
            <span class="text-sm font-medium">1.23+</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Vue 版本</span>
            <span class="text-sm font-medium">3.x</span>
          </div>
          <div class="flex justify-between py-2">
            <span class="text-sm text-gray-500">数据库</span>
            <span class="text-sm font-medium">MySQL 8</span>
          </div>
        </div>
      </div>

      <!-- 运行状态 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          运行状态
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between items-center py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">后端服务</span>
            <n-tag type="success" size="small">运行中</n-tag>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">数据库连接</span>
            <n-tag type="success" size="small">正常</n-tag>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Redis</span>
            <n-tag :type="redisStatus ? 'success' : 'warning'" size="small">{{ redisStatus ? '已连接' : '未配置' }}</n-tag>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">Elasticsearch</span>
            <n-tag type="warning" size="small">未连接</n-tag>
          </div>
          <div class="flex justify-between items-center py-2">
            <span class="text-sm text-gray-500">AI 服务</span>
            <n-tag :type="aiStatus ? 'success' : 'warning'" size="small">{{ aiStatus ? '可用' : '未配置' }}</n-tag>
          </div>
        </div>
      </div>

      <!-- 环境变量 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
          环境配置
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">服务端口</span>
            <span class="text-sm font-mono">8080</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">前端端口</span>
            <span class="text-sm font-mono">5173</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">JWT 过期时间</span>
            <span class="text-sm">7 天</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-50">
            <span class="text-sm text-gray-500">上传目录</span>
            <span class="text-sm font-mono">./uploads</span>
          </div>
          <div class="flex justify-between py-2">
            <span class="text-sm text-gray-500">日志级别</span>
            <span class="text-sm">INFO</span>
          </div>
        </div>
      </div>

      <!-- 快捷操作 -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-base font-semibold text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
          快捷操作
        </h3>
        <div class="space-y-3">
          <n-button block @click="$router.push('/admin/configs')">系统配置</n-button>
          <n-button block @click="$router.push('/admin/system-logs')">系统日志</n-button>
          <n-button block @click="$router.push('/admin/ai-stats')">AI 统计</n-button>
          <n-button block @click="$router.push('/admin/page-stats')">访问统计</n-button>
          <n-button block @click="checkHealth" :loading="checking">健康检查</n-button>
        </div>
      </div>
    </div>

    <!-- 健康检查结果 -->
    <n-modal v-model:show="showHealth" preset="card" title="健康检查结果" style="max-width: 400px">
      <div class="space-y-2">
        <div v-for="(value, key) in healthData" :key="key" class="flex justify-between py-2 border-b border-gray-50">
          <span class="text-sm text-gray-500">{{ key }}</span>
          <n-tag :type="value === 'ok' || value === true ? 'success' : 'error'" size="small">{{ value }}</n-tag>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { NTag, NButton, NInput, NSelect, NInputNumber, useMessage } from 'naive-ui'
import { get, put, post } from '@/api/request'

const message = useMessage()

const redisStatus = ref(false)
const aiStatus = ref(false)
const checking = ref(false)
const showHealth = ref(false)
const healthData = ref<any>({})
const savingAI = ref(false)
const testingAI = ref(false)
const aiTestResult = ref<any>(null)

const providerOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'DeepSeek', value: 'deepseek' },
]

const aiConfig = reactive({
  ai_provider: 'openai',
  ai_api_key: '',
  ai_base_url: '',
  ai_model: 'gpt-4o',
  ai_max_tokens: 2048,
  ai_temperature: 0.7,
  ai_rate_limit: 60,
})

async function checkHealth() {
  checking.value = true
  try {
    const res = await get<any>('/health')
    healthData.value = res.data || {}
    showHealth.value = true
  } catch {} finally { checking.value = false }
}

async function loadAIConfig() {
  try {
    const res = await get<any>('/admin/ai-config')
    if (res.data) {
      const d = res.data
      aiConfig.ai_provider = d.ai_provider || 'openai'
      aiConfig.ai_api_key = d.ai_api_key || ''
      aiConfig.ai_base_url = d.ai_base_url || ''
      aiConfig.ai_model = d.ai_model || 'gpt-4o'
      aiConfig.ai_max_tokens = parseInt(d.ai_max_tokens) || 2048
      aiConfig.ai_temperature = parseFloat(d.ai_temperature) || 0.7
      aiConfig.ai_rate_limit = parseInt(d.ai_rate_limit) || 60
      aiStatus.value = d.ai_api_key !== '' && d.ai_api_key !== '****'
    }
  } catch {}
}

async function saveAIConfig() {
  savingAI.value = true
  try {
    const payload: Record<string, string> = {}
    for (const [k, v] of Object.entries(aiConfig)) {
      payload[k] = String(v)
    }
    await put('/admin/ai-config', payload)
    message.success('AI 配置已保存')
    await loadAIConfig()
  } catch {} finally { savingAI.value = false }
}

async function testAI() {
  testingAI.value = true
  aiTestResult.value = null
  try {
    const res = await post<any>('/admin/ai-config/test')
    aiTestResult.value = res.data
  } catch {} finally { testingAI.value = false }
}

onMounted(() => {
  checkHealth()
  loadAIConfig()
})
</script>