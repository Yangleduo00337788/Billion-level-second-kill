<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="bg-white rounded-2xl p-8 border border-gray-100 shadow-sm">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-dark tracking-tight">推理引擎</h1>
          <p class="text-sm text-gray-400 mt-2">欢迎回来</p>
        </div>

        <n-form ref="formRef" :model="formData" :rules="rules" @submit.prevent="handleLogin">
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

          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            attr-type="submit"
            class="mt-2"
          >
            登录
          </n-button>
        </n-form>

        <div class="text-center mt-6">
          <span class="text-sm text-gray-400">还没有账号？</span>
          <router-link to="/register" class="text-sm text-primary ml-1 hover:underline">
            立即注册
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
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
