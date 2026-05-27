<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="glass-card p-8">
        <div class="text-center mb-8 relative z-10">
          <h1 class="text-2xl font-bold text-dark tracking-tight">推理引擎</h1>
          <p class="text-sm text-gray-400 mt-2">欢迎回来</p>
        </div>

        <n-form ref="formRef" :model="formData" :rules="rules" @submit.prevent="handleLogin" class="relative z-10">
          <n-form-item path="email" label="邮箱">
            <n-input
              v-model:value="formData.email"
              placeholder="请输入邮箱"
              size="large"
            />
          </n-form-item>

          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password-on="click"
            />
          </n-form-item>

          <button
            type="submit"
            class="glass-button-primary w-full py-3 text-sm font-medium mt-2"
            :class="{ 'opacity-60 pointer-events-none': loading }"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </n-form>

        <div class="relative flex items-center my-6 relative z-10">
          <div class="flex-1 border-t border-gray-200"></div>
          <span class="px-4 text-xs text-gray-400">或</span>
          <div class="flex-1 border-t border-gray-200"></div>
        </div>

        <div class="space-y-3 relative z-10">
          <button
            class="w-full flex items-center justify-center gap-3 py-2.5 px-4 rounded-xl border border-gray-200 bg-white/80 hover:bg-white text-sm font-medium text-gray-700 transition-all duration-200 hover:shadow-sm"
            @click="handleOAuthLogin('google')"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/><path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/><path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>
            使用 Google 登录
          </button>

          <button
            class="w-full flex items-center justify-center gap-3 py-2.5 px-4 rounded-xl border border-gray-200 bg-white/80 hover:bg-white text-sm font-medium text-gray-700 transition-all duration-200 hover:shadow-sm"
            @click="handleOAuthLogin('github')"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="#333"><path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/></svg>
            使用 GitHub 登录
          </button>

          <button
            class="w-full flex items-center justify-center gap-3 py-2.5 px-4 rounded-xl border border-gray-200 bg-white/80 hover:bg-white text-sm font-medium text-gray-700 transition-all duration-200 hover:shadow-sm"
            @click="handleOAuthLogin('wechat')"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="#07C160"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 01.213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.295.295a.326.326 0 00.167-.054l1.903-1.114a.864.864 0 01.717-.098 10.16 10.16 0 002.837.403c.276 0 .543-.027.811-.05a6.42 6.42 0 01-.253-1.78c0-3.54 3.28-6.41 7.326-6.41.18 0 .354.014.53.025-.838-3.2-4.153-5.47-8.069-5.47zM12.503 16.232c-3.54 0-6.41 2.493-6.41 5.573 0 3.08 2.87 5.573 6.41 5.573.725 0 1.424-.105 2.08-.3a.67.67 0 01.557.076l1.48.867a.253.253 0 00.13.042c.124 0 .228-.104.228-.228 0-.057-.022-.11-.038-.166l-.304-1.15a.457.457 0 01.165-.516C20.88 24.232 22 22.576 22 20.573c0-3.08-2.87-5.573-6.41-5.573h-.087zm-2.62 3.38a.91.91 0 110-1.82.91.91 0 010 1.82zm5.24 0a.91.91 0 110-1.82.91.91 0 010 1.82z"/></svg>
            使用微信登录
          </button>
        </div>

        <div class="text-center mt-6 relative z-10">
          <span class="text-sm text-gray-400">还没有账号？</span>
          <router-link to="/register" class="text-sm text-primary ml-1 hover:underline">立即注册</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)

const formData = reactive({
  email: '',
  password: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

function handleOAuthLogin(provider: string) {
  window.location.href = '/api/v1/auth/' + provider
}

// Handle OAuth redirect with token
onMounted(() => {
  const token = route.query.token as string
  if (token) {
    localStorage.setItem('token', token)
    userStore.token = token
    userStore.fetchProfile()
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  }
})

async function handleLogin() {
  loading.value = true
  try {
    await userStore.login({ email: formData.email, password: formData.password })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>
