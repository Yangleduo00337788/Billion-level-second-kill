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

        <div class="text-center mt-6 relative z-10">
          <span class="text-sm text-gray-400">还没有账号？</span>
          <router-link to="/register" class="text-sm text-primary ml-1 hover:underline">立即注册</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import type { FormInst } from 'naive-ui'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref<FormInst | null>(null)
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

async function handleLogin() {
  // 先触发表单验证
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

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
