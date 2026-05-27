<template>
  <Layout>
    <div class="admin-page">
      <div class="page-header">
        <div class="page-title">
          <h1>角色管理</h1>
          <p>管理系统角色和权限配置</p>
        </div>
        <button class="add-btn" @click="handleAdd">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          新增角色
        </button>
      </div>

      <div class="content-card">
        <div class="table-wrapper" v-loading="loading">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>角色名称</th>
                <th>角色编码</th>
                <th>描述</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="role in roles" :key="role.id">
                <td class="id-cell">{{ role.id }}</td>
                <td>
                  <div class="role-name">
                    <div class="role-icon">{{ role.roleName?.charAt(0) }}</div>
                    <span>{{ role.roleName }}</span>
                  </div>
                </td>
                <td><code class="code-tag">{{ role.roleCode }}</code></td>
                <td class="desc-cell">{{ role.description || '-' }}</td>
                <td>{{ formatDateTime(role.createTime) }}</td>
                <td class="actions-cell">
                  <div class="action-buttons">
                    <button class="action-btn edit" @click="handleEdit(role)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="handleDelete(role)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="roles.length === 0 && !loading" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path d="M2 17l10 5 10-5"/>
              <path d="M2 12l10 5 10-5"/>
            </svg>
            <p>暂无角色数据</p>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="" width="480px" class="custom-dialog">
      <div class="dialog-content">
        <div class="dialog-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
        </div>
        <h3>{{ isEdit ? '编辑角色' : '新增角色' }}</h3>
        <div class="form-group">
          <label>角色名称</label>
          <input type="text" v-model="roleForm.roleName" placeholder="请输入角色名称" />
        </div>
        <div class="form-group">
          <label>角色编码</label>
          <input type="text" v-model="roleForm.roleCode" placeholder="请输入角色编码" :disabled="isEdit" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="roleForm.description" placeholder="请输入描述" rows="3"></textarea>
        </div>
      </div>
      <template #footer>
        <button class="dialog-btn cancel" @click="dialogVisible = false">取消</button>
        <button class="dialog-btn confirm" @click="handleSubmit">确定</button>
      </template>
    </el-dialog>
  </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoleList, createRole, updateRole, deleteRole } from '../api/user'
import Layout from '../components/Layout.vue'

const router = useRouter()
const loading = ref(false)
const roles = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const roleForm = ref({ id: null, roleName: '', roleCode: '', description: '' })

async function loadRoles() {
  loading.value = true
  try {
    const res = await getRoleList()
    roles.value = res.data || []
  } catch (error) {
    ElMessage.error('加载角色列表失败')
  } finally {
    loading.value = false
  }
}

function formatDateTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function handleAdd() {
  isEdit.value = false
  roleForm.value = { id: null, roleName: '', roleCode: '', description: '' }
  dialogVisible.value = true
}

function handleEdit(role) {
  isEdit.value = true
  roleForm.value = { ...role }
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!roleForm.value.roleName || !roleForm.value.roleCode) {
    ElMessage.warning('请填写角色名称和编码')
    return
  }
  try {
    if (isEdit.value) {
      await updateRole(roleForm.value.id, roleForm.value)
      ElMessage.success('更新成功')
    } else {
      await createRole(roleForm.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadRoles()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(role) {
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '警告', { type: 'warning' })
    await deleteRole(role.id)
    ElMessage.success('已删除')
    loadRoles()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

onMounted(() => {
  loadRoles()
})
</script>

<style scoped>
.admin-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 8px 0;
}

.page-title p {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn svg {
  width: 18px;
  height: 18px;
}

.add-btn:hover {
  filter: brightness(1.1);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255,80,0,0.3);
}

.content-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  overflow: hidden;
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  padding: 14px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #999;
  background: #fafafa;
}

.data-table td {
  padding: 16px;
  font-size: 14px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
}

.id-cell {
  font-weight: 600;
  color: #999;
}

.role-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.role-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #FF5000, #FF3300);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.code-tag {
  padding: 4px 10px;
  background: #f5f5f5;
  border-radius: 6px;
  font-size: 13px;
  color: #FF5000;
  font-family: monospace;
}

.desc-cell {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.actions-cell {
  width: 100px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn.edit {
  background: #FFF0E6;
  color: #FF5000;
}

.action-btn.edit:hover {
  background: #FFD8C4;
}

.action-btn.delete {
  background: #f5f5f5;
  color: #999;
}

.action-btn.delete:hover {
  background: #e0e0e0;
  color: #F44336;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #ccc;
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.dialog-content {
  text-align: center;
  padding: 10px 0;
}

.dialog-icon {
  width: 60px;
  height: 60px;
  background: #FFF0E6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  color: #FF5000;
}

.dialog-icon svg {
  width: 28px;
  height: 28px;
}

.dialog-content h3 {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin: 0 0 20px 0;
}

.form-group {
  text-align: left;
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #666;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
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
}

.form-group input:disabled {
  background: #fafafa;
  cursor: not-allowed;
}

.dialog-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.dialog-btn.cancel {
  background: #f5f5f5;
  color: #666;
}

.dialog-btn.cancel:hover {
  background: #e0e0e0;
}

.dialog-btn.confirm {
  background: linear-gradient(135deg, #FF5000, #FF3300);
  color: #fff;
}

.dialog-btn.confirm:hover {
  filter: brightness(1.1);
}
</style>