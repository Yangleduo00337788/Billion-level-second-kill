<template>
  <div class="tb-register-page">
    <div class="register-header-bar">
      <div class="header-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
          </svg>
        </div>
        <span class="logo-text">秒杀商城</span>
      </div>
      <div class="header-links">
        <router-link to="/seckill">首页</router-link>
        <span class="divider">|</span>
        <router-link to="/login">登录</router-link>
      </div>
    </div>

    <div class="register-content">
      <div class="register-banner">
        <div class="banner-image">
          <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=shopping%20registration%20welcome%20orange%20red%20gradient%20modern%20illustration%20happy%20people%20celebration&image_size=portrait_4_3" alt="欢迎注册" />
        </div>
        <div class="banner-overlay"></div>
        <div class="banner-text">
          <h2>加入我们</h2>
          <p>开启您的秒杀之旅</p>
          <div class="banner-features">
            <div class="feature-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              <span>专属优惠</span>
            </div>
            <div class="feature-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              <span>会员特权</span>
            </div>
          </div>
        </div>
      </div>

      <div class="register-box">
        <div class="register-title">
          <h3>新用户注册</h3>
          <p>注册即可享受专属秒杀优惠</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" class="register-form">
          <el-form-item prop="username">
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              <input 
                v-model="form.username" 
                type="text" 
                placeholder="请输入用户名"
                class="tb-input"
              />
            </div>
          </el-form-item>

          <el-form-item prop="phone">
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
                <line x1="12" y1="18" x2="12.01" y2="18"/>
              </svg>
              <input 
                v-model="form.phone" 
                type="text" 
                placeholder="请输入手机号"
                class="tb-input"
              />
            </div>
          </el-form-item>

          <el-form-item prop="password">
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <input 
                v-model="form.password" 
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码（6-20位）"
                class="tb-input"
              />
              <button class="toggle-password" type="button" @click="showPassword = !showPassword">
                <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <input 
                v-model="form.confirmPassword" 
                :type="showConfirmPassword ? 'text' : 'password'"
                placeholder="请再次输入密码"
                class="tb-input"
              />
              <button class="toggle-password" type="button" @click="showConfirmPassword = !showConfirmPassword">
                <svg v-if="!showConfirmPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>
          </el-form-item>

          <el-form-item prop="agree">
            <div class="agree-wrap">
              <input type="checkbox" v-model="form.agree" />
              <span>我已阅读并同意</span>
              <a class="link">《用户协议》</a>
              <a class="link">《隐私政策》</a>
            </div>
          </el-form-item>

          <button class="register-btn" type="button" :disabled="loading" @click="handleRegister">
            <span v-if="loading" class="loading-spinner"></span>
            <span v-else>注册</span>
          </button>

          <div class="other-register">
            <div class="divider-line">
              <span>其他注册方式</span>
            </div>
            <div class="register-icons">
              <button class="icon-btn wechat" type="button" @click="handleWechatRegister">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.5c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.049.186-.094.37-.094.37-.046.186.094.37.28.37.094 0 .186-.046.28-.14l1.86-1.38a.59.59 0 0 1 .28-.094c.094 0 .186.046.28.094.744.28 1.58.466 2.46.466.14 0 .28 0 .42-.012a6.787 6.787 0 0 1-.14-1.442c0-3.792 3.56-6.88 7.98-6.88.14 0 .28 0 .42.012-.744-3.2-4.02-5.5-7.98-5.5zm-2.46 4.02c.56 0 1.02.466 1.02 1.02s-.466 1.02-1.02 1.02-1.02-.466-1.02-1.02.466-1.02 1.02-1.02zm4.92 0c.56 0 1.02.466 1.02 1.02s-.466 1.02-1.02 1.02-1.02-.466-1.02-1.02.466-1.02 1.02-1.02zm5.856 3.2c-4.02 0-7.28 2.8-7.28 6.18 0 1.86 1.02 3.5 2.6 4.6a.468.468 0 0 1 .14.56l-.28 1.2c-.046.14-.094.28-.094.28-.046.14.094.28.234.28.094 0 .14-.046.234-.094l1.44-1.02a.468.468 0 0 1 .234-.046c.094 0 .14.046.234.046.56.186 1.2.28 1.86.28 4.02 0 7.28-2.8 7.28-6.18s-3.26-6.18-7.28-6.18zm-2.46 3.2c.466 0 .84.37.84.84s-.37.84-.84.84-.84-.37-.84-.84.37-.84.84-.84zm4.92 0c.466 0 .84.37.84.84s-.37.84-.84.84-.84-.37-.84-.84.37-.84.84-.84z"/>
                </svg>
                <span class="icon-label">微信</span>
              </button>
              <button class="icon-btn google" type="button" @click="handleGoogleRegister">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                <span class="icon-label">Google</span>
              </button>
            </div>
          </div>

          <div class="login-link">
            <span>已有账号?</span>
            <router-link to="/login">立即登录</router-link>
          </div>
        </el-form>
      </div>
    </div>

    <footer class="register-footer">
      <div class="footer-links">
        <a>关于我们</a>
        <span class="dot">·</span>
        <a>联系客服</a>
        <span class="dot">·</span>
        <a>隐私政策</a>
        <span class="dot">·</span>
        <a>用户协议</a>
      </div>
      <p class="copyright">© 2024 秒杀商城 版权所有</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

const form = reactive({
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
  agree: false
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const validateAgree = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请阅读并同意用户协议'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3-20之间', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在6-20之间', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  agree: [
    { validator: validateAgree, trigger: 'change' }
  ]
}

async function handleRegister() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  
  loading.value = true
  try {
    await userStore.register({
      username: form.username,
      phone: form.phone,
      password: form.password
    })
    ElMessage.success('注册成功')
    router.push('/login')
  } catch (error) {
    console.error('注册失败', error)
  } finally {
    loading.value = false
  }
}

function handleWechatRegister() {
  ElMessage.info('微信注册功能正在开发中...')
}

function handleGoogleRegister() {
  ElMessage.info('Google注册功能正在开发中...')
}
</script>

<style scoped>
.tb-register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #FFF0E6 0%, #FFF0E6 50%, #FFD8C4 100%);
  display: flex;
  flex-direction: column;
}

.register-header-bar {
  background: linear-gradient(135deg, #FF5000 0%, #FF6633 50%, #FF7A33 100%);
  padding: 12px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: rgba(255,255,255,0.95);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #FF5000;
}

.logo-icon svg {
  width: 20px;
  height: 20px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.header-links {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-links a {
  color: #fff;
  font-size: 14px;
  text-decoration: none;
}

.header-links .divider {
  color: rgba(255,255,255,0.5);
}

.register-content {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 20px;
  gap: 60px;
}

.register-banner {
  width: 420px;
  height: 480px;
  position: relative;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(255,80,0,0.2);
}

.banner-image {
  width: 100%;
  height: 100%;
}

.banner-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,80,0,0.3) 0%, rgba(255,102,51,0.2) 100%);
}

.banner-text {
  position: absolute;
  bottom: 40px;
  left: 40px;
  right: 40px;
  color: #fff;
}

.banner-text h2 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 8px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.2);
}

.banner-text p {
  font-size: 16px;
  opacity: 0.95;
  margin-bottom: 20px;
}

.banner-features {
  display: flex;
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255,255,255,0.25);
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  backdrop-filter: blur(4px);
}

.feature-item svg {
  width: 16px;
  height: 16px;
}

.register-box {
  width: 400px;
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);
}

.register-title {
  text-align: center;
  margin-bottom: 24px;
}

.register-title h3 {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 0 0 8px;
}

.register-title p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-wrap {
  display: flex;
  align-items: center;
  background: #f5f5f5;
  border-radius: 24px;
  padding: 14px 16px;
  border: 1px solid transparent;
  transition: all 0.2s;
  width: 100%;
  box-sizing: border-box;
}

.input-wrap:focus-within {
  background: #fff;
  border-color: #FF5000;
  box-shadow: 0 0 0 3px rgba(255,80,0,0.1);
}

.input-icon {
  width: 20px;
  height: 20px;
  color: #999;
  margin-right: 12px;
}

.tb-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #333;
  outline: none;
  width: 100%;
}

.tb-input::placeholder {
  color: #999;
}

.toggle-password {
  background: transparent;
  border: none;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  cursor: pointer;
}

.toggle-password svg {
  width: 18px;
  height: 18px;
}

.agree-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #666;
}

.agree-wrap input {
  accent-color: #FF5000;
}

.agree-wrap .link {
  color: #FF5000;
  cursor: pointer;
}

.register-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border: none;
  border-radius: 24px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 8px;
}

.register-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(255,80,0,0.3);
}

.register-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.other-register {
  margin-top: 24px;
}

.divider-line {
  text-align: center;
  position: relative;
  margin-bottom: 16px;
}

.divider-line::before,
.divider-line::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 30%;
  height: 1px;
  background: #e0e0e0;
}

.divider-line::before { left: 0; }
.divider-line::after { right: 0; }

.divider-line span {
  font-size: 12px;
  color: #999;
  background: #fff;
  padding: 0 12px;
}

.register-icons {
  display: flex;
  justify-content: center;
  gap: 24px;
}

.icon-btn {
  width: 100px;
  height: 48px;
  border-radius: 24px;
  border: 1px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fff;
}

.icon-btn svg {
  width: 24px;
  height: 24px;
}

.icon-label {
  font-size: 13px;
  color: #666;
}

.icon-btn.wechat {
  border-color: #07C160;
}

.icon-btn.wechat svg {
  color: #07C160;
}

.icon-btn.wechat:hover {
  background: #07C160;
  color: #fff;
}

.icon-btn.wechat:hover svg,
.icon-btn.wechat:hover .icon-label {
  color: #fff;
}

.icon-btn.google {
  border-color: #4285F4;
}

.icon-btn.google svg {
  color: #4285F4;
}

.icon-btn.google:hover {
  background: #4285F4;
}

.icon-btn.google:hover svg,
.icon-btn.google:hover .icon-label {
  color: #fff;
}

.icon-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.login-link {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #666;
}

.login-link a {
  color: #FF7490;
  font-weight: 500;
  text-decoration: none;
  margin-left: 4px;
}

.register-footer {
  padding: 24px;
  text-align: center;
  background: #fff;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 12px;
}

.footer-links a {
  font-size: 12px;
  color: #666;
  cursor: pointer;
}

.footer-links .dot {
  color: #ccc;
}

.copyright {
  font-size: 12px;
  color: #999;
}

@media (max-width: 900px) {
  .register-banner {
    display: none;
  }

  .register-content {
    padding: 20px;
  }

  .register-box {
    width: 100%;
    max-width: 400px;
  }
}
</style>