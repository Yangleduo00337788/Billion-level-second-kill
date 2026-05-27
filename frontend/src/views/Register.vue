<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="glass-card p-8">
        <div class="text-center mb-8 relative z-10">
          <h1 class="text-2xl font-bold text-dark tracking-tight">推理引擎</h1>
          <p class="text-sm text-gray-400 mt-2">创建你的账号</p>
        </div>

        <n-form ref="formRef" :model="formData" :rules="rules" @submit.prevent="handleRegister" class="relative z-10">
          <n-form-item path="username" label="用户名">
            <n-input
              v-model:value="formData.username"
              placeholder="请输入用户名"
              size="large"
            />
          </n-form-item>

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

          <n-form-item path="confirmPassword" label="确认密码">
            <n-input
              v-model:value="formData.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              size="large"
              show-password-on="click"
            />
          </n-form-item>

          <button
            type="submit"
            class="glass-button-primary w-full py-3 text-sm font-medium mt-2"
            :class="{ 'opacity-60 pointer-events-none': loading }"
            @click="handleRegister"
          >
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </n-form>

        <div class="text-center mt-6 relative z-10">
          <span class="text-sm text-gray-400">已有账号？</span>
          <router-link to="/login" class="text-sm text-primary ml-1 hover:underline">
            立即登录
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import type { MessageApiInjection } from 'naive-ui/es/message/src/MessageProvider'

const router = useRouter()
const userStore = useUserStore()
const message = inject<MessageApiInjection>('message')!

const formRef = ref(null)
const loading = ref(false)

const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度2-20位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        return value === formData.password
      },
      message: '两次密码不一致',
      trigger: 'blur'
    }
  ]
}

async function handleRegister() {
  loading.value = true
  try {
    await userStore.register({
      username: formData.username,
      email: formData.email,
      password: formData.password
    })
    message.success('注册成功，请登录')
    router.push('/login')
  } catch {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>
