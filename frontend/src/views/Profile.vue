<template>
  <Layout>
    <div class="profile-page">
      <div class="page-header">
        <h1>个人资料</h1>
        <p>管理您的账户信息和基本资料</p>
      </div>

      <div class="profile-card">
        <div class="profile-avatar-section">
          <div class="avatar-upload">
            <div class="avatar-preview" :style="{ background: avatarColor }">
              {{ profile.nickname?.charAt(0) || profile.realName?.charAt(0) || 'U' }}
            </div>
            <button class="upload-btn">更换头像</button>
          </div>
        </div>

        <div class="profile-form" v-loading="loading">
          <div class="form-row">
            <div class="form-group">
              <label>昵称</label>
              <input type="text" v-model="form.nickname" placeholder="请输入昵称" />
            </div>
            <div class="form-group">
              <label>真实姓名</label>
              <input type="text" v-model="form.realName" placeholder="请输入真实姓名" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>性别</label>
              <div class="gender-select">
                <button
                  v-for="g in genders"
                  :key="g.value"
                  :class="['gender-btn', { active: form.gender === g.value }]"
                  @click="form.gender = g.value"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path :d="g.icon"/></svg>
                  {{ g.label }}
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>生日</label>
              <input type="date" v-model="form.birthday" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>邮箱</label>
              <input type="email" v-model="form.email" placeholder="请输入邮箱" />
            </div>
            <div class="form-group">
              <label>手机号</label>
              <input type="tel" v-model="form.phone" placeholder="请输入手机号" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group full">
              <label>个人简介</label>
              <textarea v-model="form.bio" placeholder="介绍一下自己..." rows="4"></textarea>
            </div>
          </div>

          <div class="form-actions">
            <button class="save-btn" @click="handleSave" :disabled="saving">
              <svg v-if="!saving" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                <polyline points="17 21 17 13 7 13 7 21"/>
                <polyline points="7 3 7 8 15 8"/>
              </svg>
              <span v-if="saving" class="loading-spinner"></span>
              {{ saving ? '保存中...' : '保存资料' }}
            </button>
            <button class="reset-btn" @click="loadProfile">重置</button>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyProfile, updateMyProfile } from '../api/user'
import Layout from '../components/Layout.vue'

const loading = ref(false)
const saving = ref(false)
const profile = ref({})
const colors = ['#FF5000', '#1677FF', '#52C41A', '#FAAD14', '#EB2F96', '#13C2C2', '#722ED1']
const avatarColor = ref(colors[0])

const form = reactive({
  nickname: '',
  realName: '',
  gender: 0,
  birthday: '',
  email: '',
  phone: '',
  bio: ''
})

const genders = [
  { value: 0, label: '保密', icon: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20z' },
  { value: 1, label: '男', icon: 'M15 3h6v6M10 14L20 4M21 14v7h-7M3 10v7h7' },
  { value: 2, label: '女', icon: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zM12 6a3 3 0 1 0 0 6 3 3 0 0 0 0-6z' }
]

async function loadProfile() {
  loading.value = true
  try {
    const res = await getMyProfile()
    const data = res.data || {}
    profile.value = data
    form.nickname = data.nickname || ''
    form.realName = data.realName || ''
    form.gender = data.gender ?? 0
    form.birthday = data.birthday || ''
    form.email = data.email || ''
    form.phone = data.phone || ''
    form.bio = data.bio || ''
    avatarColor.value = colors[form.gender % colors.length]
  } catch (error) {
    ElMessage.error('加载资料失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await updateMyProfile({
      nickname: form.nickname,
      realName: form.realName,
      gender: form.gender,
      birthday: form.birthday,
      email: form.email,
      phone: form.phone,
      bio: form.bio
    })
    ElMessage.success('资料已保存')
    loadProfile()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 6px 0;
}

.page-header p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.profile-card {
  background: #fff;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
}

.profile-avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 36px;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

.avatar-preview {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
}

.upload-btn {
  padding: 8px 20px;
  background: #f5f5f5;
  border: 1px solid #e5e5e5;
  border-radius: 20px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-btn:hover {
  background: #FFF0E6;
  border-color: #FF5000;
  color: #FF5000;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.full {
  grid-column: 1 / -1;
}

.form-group label {
  font-size: 13px;
  font-weight: 600;
  color: #666;
}

.form-group input,
.form-group textarea {
  padding: 12px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 10px;
  font-size: 14px;
  color: #333;
  outline: none;
  transition: all 0.2s;
  font-family: inherit;
  resize: none;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #FF5000;
  box-shadow: 0 0 0 3px rgba(255,80,0,0.1);
}

.gender-select {
  display: flex;
  gap: 10px;
}

.gender-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
}

.gender-btn svg {
  width: 16px;
  height: 16px;
}

.gender-btn.active {
  border-color: #FF5000;
  background: #FFF0E6;
  color: #FF5000;
}

.gender-btn:hover:not(.active) {
  border-color: #ddd;
  background: #fafafa;
}

.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
}

.save-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border: none;
  border-radius: 10px;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.save-btn svg {
  width: 18px;
  height: 18px;
}

.save-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}

.save-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.reset-btn {
  padding: 14px 24px;
  background: #f5f5f5;
  border: 1px solid #e5e5e5;
  border-radius: 10px;
  font-size: 15px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.reset-btn:hover {
  background: #e0e0e0;
}

.loading-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 600px) {
  .profile-card {
    padding: 24px;
  }
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>