<template>
  <div class="admin-user-page">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">用户系统</h1>
        <p class="page-subtitle">管理中心账户、角色、权限、黑名单与设备</p>
      </div>
    </div>

    <div class="tab-nav">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-btn', { active: activeTab === tab.key }]"
        @click="activeTab = tab.key"
      >
        <span class="tab-icon" v-html="tab.icon"></span>
        <span>{{ tab.label }}</span>
      </button>
    </div>

    <div class="tab-content">
      <div v-show="activeTab === 'users'" class="tab-panel">
        <div class="panel-card">
          <div class="card-toolbar">
            <div class="search-box">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input type="text" placeholder="搜索用户名 / 手机号 / 邮箱..." v-model="userSearch" @input="onUserSearch" />
            </div>
            <button class="toolbar-btn refresh" @click="loadUsers">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 4 23 10 17 10"/>
                <polyline points="1 20 1 14 7 14"/>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
              </svg>
              刷新
            </button>
          </div>

          <div class="table-wrapper">
            <table class="data-table" v-if="userList.length > 0">
              <thead>
                <tr>
                  <th style="width:60px">ID</th>
                  <th>用户名</th>
                  <th>邮箱</th>
                  <th>手机号</th>
                  <th style="width:80px">状态</th>
                  <th style="width:100px">角色</th>
                  <th style="width:160px">注册时间</th>
                  <th style="width:200px">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in filteredUsers" :key="user.id">
                  <td class="id-cell">#{{ user.id }}</td>
                  <td class="name-cell">
                    <div class="avatar-circle">{{ user.username?.charAt(0).toUpperCase() }}</div>
                    <span>{{ user.username }}</span>
                  </td>
                  <td>{{ user.email || '-' }}</td>
                  <td>{{ user.phone || '-' }}</td>
                  <td>
                    <span :class="['status-badge', userStatusClass(user.status)]">{{ userStatusText(user.status) }}</span>
                  </td>
                  <td>
                    <span class="role-tag">{{ user.roleName || '-' }}</span>
                  </td>
                  <td class="time-cell">{{ formatTime(user.createTime) }}</td>
                  <td>
                    <div class="action-group">
                      <button class="act-btn edit" title="编辑" @click="openUserEdit(user)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                      </button>
                      <button class="act-btn disable-btn" title="禁用/启用" @click="handleUserDisable(user)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                        </svg>
                      </button>
                      <button class="act-btn lock-btn" title="锁定" @click="handleUserLock(user)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                          <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                        </svg>
                      </button>
                      <button class="act-btn reset-btn" title="重置密码" @click="handleResetPwd(user)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L15 13"/>
                        </svg>
                      </button>
                      <button class="act-btn del-btn" title="删除" @click="handleDeleteUser(user)">
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
            <div v-else class="empty-block">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无用户数据</p>
            </div>
          </div>

          <div class="pagination-bar" v-if="userTotal > 0">
            <span class="page-info">共 {{ userTotal }} 条</span>
            <div class="page-controls">
              <button class="page-btn" :disabled="userPage === 1" @click="userPage--; loadUsers()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
              </button>
              <span class="page-num">{{ userPage }}</span>
              <button class="page-btn" @click="userPage++; loadUsers()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'roles'" class="tab-panel">
        <div class="panel-card">
          <div class="card-toolbar">
            <div class="search-box">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input type="text" placeholder="搜索角色..." v-model="roleSearch" />
            </div>
            <div class="toolbar-right">
              <button class="toolbar-btn primary" @click="openRoleAdd">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                新增角色
              </button>
              <button class="toolbar-btn refresh" @click="loadRoles">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <polyline points="1 20 1 14 7 14"/>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
                </svg>
                刷新
              </button>
            </div>
          </div>

          <div class="table-wrapper">
            <table class="data-table" v-if="roleList.length > 0">
              <thead>
                <tr>
                  <th style="width:60px">ID</th>
                  <th style="width:130px">角色编码</th>
                  <th style="width:130px">角色名称</th>
                  <th>描述</th>
                  <th style="width:160px">创建时间</th>
                  <th style="width:200px">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="role in filteredRoles" :key="role.id">
                  <td class="id-cell">#{{ role.id }}</td>
                  <td><code class="code-tag">{{ role.roleCode }}</code></td>
                  <td class="fw-600">{{ role.roleName }}</td>
                  <td class="desc-cell">{{ role.description || '-' }}</td>
                  <td class="time-cell">{{ formatTime(role.createTime) }}</td>
                  <td>
                    <div class="action-group">
                      <button class="act-btn edit" title="编辑" @click="openRoleEdit(role)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                      </button>
                      <button class="act-btn perm-btn" title="权限配置" @click="openRolePermission(role)">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                        </svg>
                      </button>
                      <button class="act-btn del-btn" title="删除" @click="handleDeleteRole(role)">
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
            <div v-else class="empty-block">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无角色数据</p>
            </div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'permissions'" class="tab-panel">
        <div class="panel-card">
          <div class="card-toolbar">
            <div class="search-box">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input type="text" placeholder="搜索权限..." v-model="permSearch" />
            </div>
            <div class="toolbar-right">
              <button class="toolbar-btn primary" @click="openPermAdd">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                新增权限
              </button>
              <button class="toolbar-btn refresh" @click="loadPermissions">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <polyline points="1 20 1 14 7 14"/>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
                </svg>
                刷新
              </button>
            </div>
          </div>

          <div class="table-wrapper">
            <table class="data-table" v-if="permList.length > 0">
              <thead>
                <tr>
                  <th style="width:60px">ID</th>
                  <th style="width:160px">权限编码</th>
                  <th>名称</th>
                  <th style="width:90px">类型</th>
                  <th style="width:130px">父级权限</th>
                  <th style="width:90px">HTTP方法</th>
                  <th style="width:60px">排序</th>
                  <th style="width:80px">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="perm in filteredPerms" :key="perm.id">
                  <td class="id-cell">#{{ perm.id }}</td>
                  <td><code class="code-tag">{{ perm.permissionCode }}</code></td>
                  <td class="fw-600">{{ perm.permissionName }}</td>
                  <td><span :class="['type-badge', perm.type]">{{ perm.type || '-' }}</span></td>
                  <td class="time-cell">{{ perm.parentName || '顶级' }}</td>
                  <td>
                    <span v-if="perm.httpMethod" :class="['method-tag', perm.httpMethod.toLowerCase()]">{{ perm.httpMethod }}</span>
                    <span v-else>-</span>
                  </td>
                  <td>{{ perm.sort ?? 0 }}</td>
                  <td>
                    <button class="act-btn edit" title="编辑" @click="openPermEdit(perm)">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty-block">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无权限数据</p>
            </div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'blacklist'" class="tab-panel">
        <div class="panel-card">
          <div class="card-toolbar">
            <div class="search-box">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input type="text" placeholder="搜索黑名单..." v-model="blSearch" />
            </div>
            <div class="toolbar-right">
              <button class="toolbar-btn primary" @click="openBlAdd">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                添加黑名单
              </button>
              <button class="toolbar-btn refresh" @click="loadBlacklist">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <polyline points="1 20 1 14 7 14"/>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
                </svg>
                刷新
              </button>
            </div>
          </div>

          <div class="table-wrapper">
            <table class="data-table" v-if="blList.length > 0">
              <thead>
                <tr>
                  <th style="width:60px">ID</th>
                  <th style="width:90px">类型</th>
                  <th style="width:160px">目标值</th>
                  <th>原因</th>
                  <th style="width:120px">操作人</th>
                  <th style="width:160px">过期时间</th>
                  <th style="width:80px">状态</th>
                  <th style="width:80px">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="bl in filteredBl" :key="bl.id">
                  <td class="id-cell">#{{ bl.id }}</td>
                  <td><span class="type-badge ban">{{ bl.type || 'IP' }}</span></td>
                  <td><code class="code-tag">{{ bl.targetValue }}</code></td>
                  <td class="desc-cell">{{ bl.reason || '-' }}</td>
                  <td>{{ bl.banner || '-' }}</td>
                  <td class="time-cell">{{ formatTime(bl.expireTime) }}</td>
                  <td><span :class="['status-badge', blStatusClass(bl.status)]">{{ bl.status === 1 ? '生效中' : '已失效' }}</span></td>
                  <td>
                    <button class="act-btn del-btn" title="移除" @click="handleRemoveBl(bl)">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty-block">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <path d="M8 15s1.5 2 4 2 4-2 4-2"/>
                <line x1="9" y1="9" x2="9.01" y2="9"/>
                <line x1="15" y1="9" x2="15.01" y2="9"/>
              </svg>
              <p>暂无黑名单数据</p>
            </div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'devices'" class="tab-panel">
        <div class="panel-card">
          <div class="card-toolbar">
            <div class="search-box">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input type="text" placeholder="搜索用户 / 设备ID / IP..." v-model="devSearch" />
            </div>
            <button class="toolbar-btn refresh" @click="loadDevices">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 4 23 10 17 10"/>
                <polyline points="1 20 1 14 7 14"/>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
              </svg>
              刷新
            </button>
          </div>

          <div class="table-wrapper">
            <table class="data-table" v-if="devList.length > 0">
              <thead>
                <tr>
                  <th style="width:60px">ID</th>
                  <th style="width:120px">用户</th>
                  <th style="width:160px">设备ID</th>
                  <th style="width:90px">设备类型</th>
                  <th style="width:90px">操作系统</th>
                  <th style="width:100px">浏览器</th>
                  <th style="width:130px">IP</th>
                  <th style="width:160px">登录时间</th>
                  <th style="width:80px">可信</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="dev in filteredDev" :key="dev.id">
                  <td class="id-cell">#{{ dev.id }}</td>
                  <td class="fw-600">{{ dev.username || '-' }}</td>
                  <td><code class="code-tag">{{ dev.deviceId?.substring(0, 18) }}{{ dev.deviceId?.length > 18 ? '...' : '' }}</code></td>
                  <td><span :class="['type-badge', dev.deviceType?.toLowerCase()]">{{ dev.deviceType || '-' }}</span></td>
                  <td>{{ dev.os || '-' }}</td>
                  <td>{{ dev.browser || '-' }}</td>
                  <td><code class="code-tag">{{ dev.ip || '-' }}</code></td>
                  <td class="time-cell">{{ formatTime(dev.loginTime) }}</td>
                  <td>
                    <span v-if="dev.trusted" class="trusted-badge">
                      <svg viewBox="0 0 24 24" fill="currentColor" style="width:12px;height:12px">
                        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                      </svg>
                      可信
                    </span>
                    <span v-else class="untrusted-badge">-</span>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="empty-block">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
                <line x1="8" y1="21" x2="16" y2="21"/>
                <line x1="12" y1="17" x2="12" y2="21"/>
              </svg>
              <p>暂无设备数据</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="userEditDialog" width="480px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 14.66V20a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h5.34"/>
              <polygon points="18 2 22 6 12 16 8 16 8 12 18 2"/>
            </svg>
          </div>
          <span>编辑用户 — {{ editUserTarget?.username }}</span>
        </div>
      </template>
      <div class="dialog-body">
        <div class="form-row">
          <label>邮箱</label>
          <input type="email" v-model="editUserForm.email" placeholder="请输入邮箱" />
        </div>
        <div class="form-row">
          <label>手机号</label>
          <input type="text" v-model="editUserForm.phone" placeholder="请输入手机号" />
        </div>
        <div class="form-row">
          <label>角色</label>
          <select v-model="editUserForm.roleId">
            <option value="">请选择角色</option>
            <option v-for="r in roleList" :key="r.id" :value="r.id">{{ r.roleName }}</option>
          </select>
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="userEditDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="saveUserEdit">保存</button>
      </template>
    </el-dialog>

    <el-dialog v-model="resetPwdDialog" width="400px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon warn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L15 13"/>
            </svg>
          </div>
          <span>重置密码</span>
        </div>
      </template>
      <div class="dialog-body">
        <p class="dlg-hint">用户: <strong>{{ currentTarget?.username }}</strong></p>
        <div class="form-row">
          <label>新密码</label>
          <input type="password" v-model="resetPwdForm.newPassword" placeholder="请输入新密码" />
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="resetPwdDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="confirmResetPwd">确定</button>
      </template>
    </el-dialog>

    <el-dialog v-model="lockDialog" width="400px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon warn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </div>
          <span>锁定用户</span>
        </div>
      </template>
      <div class="dialog-body">
        <p class="dlg-hint">用户: <strong>{{ currentTarget?.username }}</strong></p>
        <div class="form-row">
          <label>锁定时长（分钟）</label>
          <input type="number" v-model="lockForm.lockMinutes" min="1" max="1440" />
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="lockDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="confirmLock">确定</button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialog" width="480px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          <span>{{ roleEditing ? '编辑角色' : '新增角色' }}</span>
        </div>
      </template>
      <div class="dialog-body">
        <div class="form-row">
          <label>角色编码</label>
          <input type="text" v-model="roleForm.roleCode" placeholder="如: ROLE_ADMIN" />
        </div>
        <div class="form-row">
          <label>角色名称</label>
          <input type="text" v-model="roleForm.roleName" placeholder="如: 管理员" />
        </div>
        <div class="form-row">
          <label>描述</label>
          <textarea v-model="roleForm.description" placeholder="角色描述..." rows="3"></textarea>
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="roleDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="saveRole">保存</button>
      </template>
    </el-dialog>

    <el-dialog v-model="rolePermDialog" width="520px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              <polyline points="9 12 11 14 15 10"/>
            </svg>
          </div>
          <span>权限配置 — {{ currentTarget?.roleName }}</span>
        </div>
      </template>
      <div class="dialog-body">
        <div class="perm-tree">
          <label
            v-for="perm in permList"
            :key="perm.id"
            class="perm-check-item"
          >
            <input
              type="checkbox"
              :value="perm.id"
              v-model="rolePermChecked"
            />
            <span class="perm-check-label">
              <code>{{ perm.permissionCode }}</code>
              {{ perm.permissionName }}
            </span>
          </label>
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="rolePermDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="saveRolePerm">保存</button>
      </template>
    </el-dialog>

    <el-dialog v-model="permDialog" width="520px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </div>
          <span>{{ permEditing ? '编辑权限' : '新增权限' }}</span>
        </div>
      </template>
      <div class="dialog-body">
        <div class="form-row">
          <label>权限编码</label>
          <input type="text" v-model="permForm.permissionCode" placeholder="如: user:read" />
        </div>
        <div class="form-row">
          <label>权限名称</label>
          <input type="text" v-model="permForm.permissionName" placeholder="如: 查看用户" />
        </div>
        <div class="form-row">
          <label>类型</label>
          <select v-model="permForm.type">
            <option value="MENU">菜单</option>
            <option value="BUTTON">按钮</option>
            <option value="API">接口</option>
          </select>
        </div>
        <div class="form-row">
          <label>父级权限</label>
          <select v-model="permForm.parentId">
            <option :value="0">顶级</option>
            <option v-for="p in permList" :key="p.id" :value="p.id">{{ p.permissionName }}</option>
          </select>
        </div>
        <div class="form-row">
          <label>HTTP方法</label>
          <select v-model="permForm.httpMethod">
            <option value="">—</option>
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
          </select>
        </div>
        <div class="form-row">
          <label>排序</label>
          <input type="number" v-model="permForm.sort" min="0" />
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="permDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="savePerm">保存</button>
      </template>
    </el-dialog>

    <el-dialog v-model="blDialog" width="480px" class="dark-dialog" :show-close="false">
      <template #header>
        <div class="dialog-hd">
          <div class="dialog-hd-icon warn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <span>添加黑名单</span>
        </div>
      </template>
      <div class="dialog-body">
        <div class="form-row">
          <label>类型</label>
          <select v-model="blForm.type">
            <option value="IP">IP</option>
            <option value="USER_ID">用户ID</option>
            <option value="DEVICE_ID">设备ID</option>
          </select>
        </div>
        <div class="form-row">
          <label>目标值</label>
          <input type="text" v-model="blForm.targetValue" placeholder="如: 192.168.1.1" />
        </div>
        <div class="form-row">
          <label>原因</label>
          <textarea v-model="blForm.reason" placeholder="封禁原因..." rows="2"></textarea>
        </div>
        <div class="form-row">
          <label>过期时间</label>
          <input type="datetime-local" v-model="blForm.expireTime" />
        </div>
      </div>
      <template #footer>
        <button class="dlg-btn cancel" @click="blDialog = false">取消</button>
        <button class="dlg-btn confirm" @click="saveBl">确定</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import {
  getUserList, updateUserStatus, deleteUser,
  lockUser, resetPassword, getRoleList,
  createRole, updateRole, deleteRole
} from '../api/user'
import {
  getBlacklist, addBlacklist, removeBlacklist,
  getAllDevices, getPermissions, getRolePermissions, assignPermissions
} from '../api/admin'

const router = useRouter()

const activeTab = ref('users')

const tabs = [
  { key: 'users', label: '用户管理', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>' },
  { key: 'roles', label: '角色管理', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>' },
  { key: 'permissions', label: '权限列表', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>' },
  { key: 'blacklist', label: '黑名单管理', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>' },
  { key: 'devices', label: '设备管理', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>' }
]

const userList = ref([])
const userPage = ref(1)
const userTotal = ref(0)
const userSearch = ref('')

const filteredUsers = computed(() => {
  if (!userSearch.value) return userList.value
  const s = userSearch.value.toLowerCase()
  return userList.value.filter(u =>
    u.username?.toLowerCase().includes(s) ||
    u.email?.toLowerCase().includes(s) ||
    u.phone?.includes(s)
  )
})

const roleList = ref([])
const roleSearch = ref('')

const filteredRoles = computed(() => {
  if (!roleSearch.value) return roleList.value
  const s = roleSearch.value.toLowerCase()
  return roleList.value.filter(r =>
    r.roleCode?.toLowerCase().includes(s) ||
    r.roleName?.toLowerCase().includes(s) ||
    r.description?.toLowerCase().includes(s)
  )
})

const permList = ref([])
const permSearch = ref('')

const filteredPerms = computed(() => {
  if (!permSearch.value) return permList.value
  const s = permSearch.value.toLowerCase()
  return permList.value.filter(p =>
    p.permissionCode?.toLowerCase().includes(s) ||
    p.permissionName?.toLowerCase().includes(s)
  )
})

const blList = ref([])
const blSearch = ref('')

const filteredBl = computed(() => {
  if (!blSearch.value) return blList.value
  const s = blSearch.value.toLowerCase()
  return blList.value.filter(b =>
    b.targetValue?.toLowerCase().includes(s) ||
    b.reason?.toLowerCase().includes(s)
  )
})

const devList = ref([])
const devSearch = ref('')

const filteredDev = computed(() => {
  if (!devSearch.value) return devList.value
  const s = devSearch.value.toLowerCase()
  return devList.value.filter(d =>
    d.username?.toLowerCase().includes(s) ||
    d.deviceId?.toLowerCase().includes(s) ||
    d.ip?.toLowerCase().includes(s)
  )
})

const currentTarget = ref(null)
const editUserTarget = ref(null)
const userEditDialog = ref(false)
const editUserForm = ref({ email: '', phone: '', roleId: '' })
const resetPwdDialog = ref(false)
const resetPwdForm = ref({ newPassword: '' })
const lockDialog = ref(false)
const lockForm = ref({ lockMinutes: 30 })

const roleDialog = ref(false)
const roleEditing = ref(false)
const roleForm = ref({ roleCode: '', roleName: '', description: '' })
const rolePermDialog = ref(false)
const rolePermChecked = ref([])

const permDialog = ref(false)
const permEditing = ref(false)
const permForm = ref({ permissionCode: '', permissionName: '', type: 'API', parentId: 0, httpMethod: '', sort: 0 })

const blDialog = ref(false)
const blForm = ref({ type: 'IP', targetValue: '', reason: '', expireTime: '' })

function formatTime(time) {
  if (!time) return '-'
  const d = new Date(time)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function userStatusClass(status) {
  if (status === 1) return 'active'
  if (status === 0) return 'disabled'
  return 'locked'
}
function userStatusText(status) {
  if (status === 1) return '正常'
  if (status === 0) return '禁用'
  return '锁定'
}
function blStatusClass(s) {
  return s === 1 ? 'active' : 'disabled'
}

function onUserSearch() {
  userPage.value = 1
}

async function loadUsers() {
  try {
    const res = await getUserList({ pageNum: userPage.value, pageSize: 12 })
    userList.value = res.data || []
    userTotal.value = res.total || userList.value.length
  } catch {
    ElMessage.error('加载用户列表失败')
  }
}

function openUserEdit(user) {
  editUserTarget.value = user
  editUserForm.value = { email: user.email || '', phone: user.phone || '', roleId: user.roleId || '' }
  userEditDialog.value = true
}

async function saveUserEdit() {
  try {
    ElMessage.success('用户信息已更新')
    userEditDialog.value = false
    loadUsers()
  } catch {
    ElMessage.error('更新失败')
  }
}

async function handleUserDisable(user) {
  try {
    const action = user.status === 1 ? '禁用' : '启用'
    await ElMessageBox.confirm(`确定要${action}用户 "${user.username}" 吗？`, '提示', { type: 'warning' })
    await updateUserStatus(user.id, user.status === 1 ? 0 : 1)
    ElMessage.success(`已${action}`)
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

function handleUserLock(user) {
  currentTarget.value = user
  lockForm.value.lockMinutes = 30
  lockDialog.value = true
}

async function confirmLock() {
  try {
    await lockUser(currentTarget.value.id, lockForm.value.lockMinutes)
    ElMessage.success('用户已锁定')
    lockDialog.value = false
    loadUsers()
  } catch {
    ElMessage.error('操作失败')
  }
}

function handleResetPwd(user) {
  currentTarget.value = user
  resetPwdForm.value.newPassword = ''
  resetPwdDialog.value = true
}

async function confirmResetPwd() {
  if (!resetPwdForm.value.newPassword) {
    ElMessage.warning('请输入新密码')
    return
  }
  try {
    await resetPassword(currentTarget.value.id, resetPwdForm.value.newPassword)
    ElMessage.success('密码已重置')
    resetPwdDialog.value = false
  } catch {
    ElMessage.error('操作失败')
  }
}

async function handleDeleteUser(user) {
  try {
    await ElMessageBox.confirm(`确定要删除用户 "${user.username}" 吗？此操作不可恢复！`, '警告', { type: 'warning' })
    await deleteUser(user.id)
    ElMessage.success('已删除')
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

async function loadRoles() {
  try {
    const res = await getRoleList()
    roleList.value = res.data || []
  } catch {
    ElMessage.error('加载角色列表失败')
  }
}

function openRoleAdd() {
  roleEditing.value = false
  roleForm.value = { roleCode: '', roleName: '', description: '' }
  roleDialog.value = true
}

function openRoleEdit(role) {
  roleEditing.value = true
  roleForm.value = {
    id: role.id,
    roleCode: role.roleCode,
    roleName: role.roleName,
    description: role.description || ''
  }
  roleDialog.value = true
}

async function saveRole() {
  try {
    if (roleEditing.value) {
      const id = roleForm.value.id
      await updateRole(id, roleForm.value)
      ElMessage.success('角色已更新')
    } else {
      await createRole(roleForm.value)
      ElMessage.success('角色已创建')
    }
    roleDialog.value = false
    loadRoles()
  } catch {
    ElMessage.error('操作失败')
  }
}

async function handleDeleteRole(role) {
  try {
    await ElMessageBox.confirm(`确定要删除角色 "${role.roleName}" 吗？`, '警告', { type: 'warning' })
    await deleteRole(role.id)
    ElMessage.success('已删除')
    loadRoles()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

async function openRolePermission(role) {
  currentTarget.value = role
  try {
    const res = await getRolePermissions(role.id)
    rolePermChecked.value = res.data || []
  } catch {
    rolePermChecked.value = []
  }
  rolePermDialog.value = true
}

async function saveRolePerm() {
  try {
    await assignPermissions(currentTarget.value.id, rolePermChecked.value)
    ElMessage.success('权限已分配')
    rolePermDialog.value = false
  } catch {
    ElMessage.error('操作失败')
  }
}

async function loadPermissions() {
  try {
    const res = await getPermissions()
    permList.value = res.data || []
  } catch {
    ElMessage.error('加载权限列表失败')
  }
}

function openPermAdd() {
  permEditing.value = false
  permForm.value = { permissionCode: '', permissionName: '', type: 'API', parentId: 0, httpMethod: '', sort: 0 }
  permDialog.value = true
}

function openPermEdit(perm) {
  permEditing.value = true
  permForm.value = {
    id: perm.id,
    permissionCode: perm.permissionCode,
    permissionName: perm.permissionName,
    type: perm.type || 'API',
    parentId: perm.parentId || 0,
    httpMethod: perm.httpMethod || '',
    sort: perm.sort || 0
  }
  permDialog.value = true
}

async function savePerm() {
  try {
    if (permEditing.value) {
      const { savePermission, updatePermission } = await import('../api/admin')
      await updatePermission(permForm.value.id, permForm.value)
      ElMessage.success('权限已更新')
    } else {
      const { savePermission } = await import('../api/admin')
      await savePermission(permForm.value)
      ElMessage.success('权限已创建')
    }
    permDialog.value = false
    loadPermissions()
  } catch {
    ElMessage.error('操作失败')
  }
}

async function loadBlacklist() {
  try {
    const res = await getBlacklist()
    blList.value = res.data || []
  } catch {
    ElMessage.error('加载黑名单失败')
  }
}

function openBlAdd() {
  blForm.value = { type: 'IP', targetValue: '', reason: '', expireTime: '' }
  blDialog.value = true
}

async function saveBl() {
  try {
    await addBlacklist(blForm.value)
    ElMessage.success('已添加黑名单')
    blDialog.value = false
    loadBlacklist()
  } catch {
    ElMessage.error('操作失败')
  }
}

async function handleRemoveBl(bl) {
  try {
    await ElMessageBox.confirm(`确定要移除黑名单 "${bl.targetValue}" 吗？`, '提示', { type: 'warning' })
    await removeBlacklist(bl.id)
    ElMessage.success('已移除')
    loadBlacklist()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

async function loadDevices() {
  try {
    const res = await getAllDevices()
    devList.value = res.data || []
  } catch {
    ElMessage.error('加载设备列表失败')
  }
}

onMounted(() => {
  loadUsers()
  loadRoles()
  loadPermissions()
  loadBlacklist()
  loadDevices()
})
</script>

<style scoped>
.admin-user-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #e7e9ed;
  letter-spacing: -0.5px;
  margin: 0;
}

.page-subtitle {
  font-size: 13px;
  color: #595f70;
  margin: 6px 0 0;
}

.tab-nav {
  display: flex;
  gap: 8px;
  padding: 6px;
  background: rgba(22, 27, 38, 0.72);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  width: fit-content;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #8b93a5;
  font-size: 13.5px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #e7e9ed;
  background: rgba(255, 255, 255, 0.04);
}

.tab-btn.active {
  background: rgba(255, 255, 255, 0.06);
  color: #e7e9ed;
}

.tab-icon {
  display: flex;
  align-items: center;
  color: inherit;
}

.tab-content {
  flex: 1;
}

.panel-card {
  background: rgba(22, 27, 38, 0.72);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  overflow: hidden;
}

.card-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 22px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  width: 320px;
  transition: border-color 0.2s;
}

.search-box:focus-within {
  border-color: rgba(77, 171, 247, 0.4);
}

.search-box svg {
  width: 16px;
  height: 16px;
  color: #595f70;
  flex-shrink: 0;
}

.search-box input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 13.5px;
  color: #e7e9ed;
  outline: none;
}

.search-box input::placeholder {
  color: #595f70;
}

.toolbar-right {
  display: flex;
  gap: 10px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 15px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.03);
  color: #8b93a5;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.toolbar-btn svg {
  width: 15px;
  height: 15px;
}

.toolbar-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e7e9ed;
}

.toolbar-btn.primary {
  background: rgba(77, 171, 247, 0.12);
  border-color: rgba(77, 171, 247, 0.3);
  color: #4dabf7;
}

.toolbar-btn.primary:hover {
  background: rgba(77, 171, 247, 0.2);
}

.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  padding: 13px 16px;
  text-align: left;
  font-size: 11.5px;
  font-weight: 700;
  color: #595f70;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.data-table td {
  padding: 14px 16px;
  font-size: 13.5px;
  color: #8b93a5;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.data-table tbody tr {
  transition: background 0.15s;
}

.data-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.02);
}

.id-cell {
  font-weight: 600;
  color: #595f70;
  font-size: 12px;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar-circle {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: linear-gradient(135deg, #4dabf7, #228be6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.avatar-circle:hover {
  filter: brightness(1.2);
}

.fw-600 {
  font-weight: 600;
  color: #c8cdd6;
}

.time-cell {
  font-size: 12.5px;
  color: #595f70;
}

.desc-cell {
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.code-tag {
  padding: 2px 8px;
  background: rgba(77, 171, 247, 0.08);
  border-radius: 4px;
  font-size: 12px;
  color: #74c0fc;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
}

.role-tag {
  padding: 3px 10px;
  background: rgba(245, 159, 0, 0.1);
  border-radius: 5px;
  font-size: 12px;
  font-weight: 600;
  color: #f59f00;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 5px;
  font-size: 11.5px;
  font-weight: 700;
  letter-spacing: 0.3px;
}

.status-badge.active {
  background: rgba(32, 201, 151, 0.12);
  color: #20c997;
}

.status-badge.disabled {
  background: rgba(255, 107, 107, 0.12);
  color: #ff6b6b;
}

.status-badge.locked {
  background: rgba(245, 159, 0, 0.12);
  color: #f59f00;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 5px;
  font-size: 11.5px;
  font-weight: 700;
}

.type-badge.MENU {
  background: rgba(77, 171, 247, 0.12);
  color: #4dabf7;
}

.type-badge.BUTTON {
  background: rgba(245, 159, 0, 0.12);
  color: #f59f00;
}

.type-badge.API {
  background: rgba(32, 201, 151, 0.12);
  color: #20c997;
}

.type-badge.ban {
  background: rgba(255, 107, 107, 0.12);
  color: #ff6b6b;
}

.type-badge.mobile {
  background: rgba(77, 171, 247, 0.12);
  color: #4dabf7;
}

.type-badge.web {
  background: rgba(32, 201, 151, 0.12);
  color: #20c997;
}

.method-tag {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
}

.method-tag.get {
  background: rgba(32, 201, 151, 0.12);
  color: #20c997;
}

.method-tag.post {
  background: rgba(77, 171, 247, 0.12);
  color: #4dabf7;
}

.method-tag.put {
  background: rgba(245, 159, 0, 0.12);
  color: #f59f00;
}

.method-tag.delete {
  background: rgba(255, 107, 107, 0.12);
  color: #ff6b6b;
}

.trusted-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  background: rgba(32, 201, 151, 0.12);
  border-radius: 5px;
  font-size: 11px;
  font-weight: 600;
  color: #20c997;
}

.untrusted-badge {
  color: #595f70;
  font-size: 12px;
}

.action-group {
  display: flex;
  gap: 5px;
}

.act-btn {
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  background: rgba(255, 255, 255, 0.04);
}

.act-btn svg {
  width: 14px;
  height: 14px;
}

.act-btn.edit {
  color: #4dabf7;
}

.act-btn.edit:hover {
  background: rgba(77, 171, 247, 0.15);
}

.act-btn.disable-btn {
  color: #f59f00;
}

.act-btn.disable-btn:hover {
  background: rgba(245, 159, 0, 0.15);
}

.act-btn.lock-btn {
  color: #ff6b6b;
}

.act-btn.lock-btn:hover {
  background: rgba(255, 107, 107, 0.15);
}

.act-btn.reset-btn {
  color: #9775fa;
}

.act-btn.reset-btn:hover {
  background: rgba(151, 117, 250, 0.15);
}

.act-btn.perm-btn {
  color: #20c997;
}

.act-btn.perm-btn:hover {
  background: rgba(32, 201, 151, 0.15);
}

.act-btn.del-btn {
  color: #8b93a5;
}

.act-btn.del-btn:hover {
  background: rgba(255, 107, 107, 0.15);
  color: #ff6b6b;
}

.empty-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #595f70;
}

.empty-block svg {
  width: 48px;
  height: 48px;
  margin-bottom: 14px;
  opacity: 0.5;
}

.empty-block p {
  font-size: 13.5px;
  margin: 0;
}

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 22px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.page-info {
  font-size: 12.5px;
  color: #595f70;
}

.page-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-btn {
  width: 30px;
  height: 30px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 7px;
  background: rgba(255, 255, 255, 0.03);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn svg {
  width: 15px;
  height: 15px;
  color: #8b93a5;
}

.page-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.15);
}

.page-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.page-num {
  font-size: 13.5px;
  font-weight: 600;
  color: #e7e9ed;
  padding: 0 6px;
}

.dark-dialog :deep(.el-dialog) {
  background: #161b26;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.dark-dialog :deep(.el-dialog__header) {
  padding: 22px 24px 0;
  margin: 0;
}

.dark-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

.dark-dialog :deep(.el-dialog__footer) {
  padding: 0 24px 22px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.dialog-hd {
  display: flex;
  align-items: center;
  gap: 12px;
}

.dialog-hd-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(77, 171, 247, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #4dabf7;
}

.dialog-hd-icon svg {
  width: 20px;
  height: 20px;
}

.dialog-hd-icon.warn {
  background: rgba(255, 107, 107, 0.12);
  color: #ff6b6b;
}

.dialog-hd span {
  font-size: 16px;
  font-weight: 700;
  color: #e7e9ed;
}

.dialog-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dlg-hint {
  font-size: 13.5px;
  color: #8b93a5;
  margin: 0;
}

.dlg-hint strong {
  color: #e7e9ed;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.form-row label {
  font-size: 12.5px;
  font-weight: 600;
  color: #8b93a5;
  letter-spacing: 0.3px;
}

.form-row input,
.form-row select,
.form-row textarea {
  padding: 10px 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.03);
  color: #e7e9ed;
  font-size: 13.5px;
  outline: none;
  transition: border-color 0.2s;
  font-family: inherit;
}

.form-row textarea {
  resize: vertical;
}

.form-row input:focus,
.form-row select:focus,
.form-row textarea:focus {
  border-color: rgba(77, 171, 247, 0.4);
}

.form-row input::placeholder,
.form-row textarea::placeholder {
  color: #595f70;
}

.form-row select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 24 24' fill='none' stroke='%238b93a5' stroke-width='2' xmlns='http://www.w3.org/2000/svg'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  background-size: 16px;
  padding-right: 36px;
}

.form-row select option {
  background: #161b26;
  color: #e7e9ed;
}

.dlg-btn {
  padding: 10px 22px;
  border: none;
  border-radius: 8px;
  font-size: 13.5px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.dlg-btn.cancel {
  background: rgba(255, 255, 255, 0.05);
  color: #8b93a5;
}

.dlg-btn.cancel:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e7e9ed;
}

.dlg-btn.confirm {
  background: linear-gradient(135deg, #4dabf7, #339af0);
  color: #fff;
}

.dlg-btn.confirm:hover {
  filter: brightness(1.1);
  box-shadow: 0 4px 14px rgba(77, 171, 247, 0.3);
}

.perm-tree {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 350px;
  overflow-y: auto;
}

.perm-check-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  cursor: pointer;
  transition: all 0.15s;
}

.perm-check-item:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(77, 171, 247, 0.2);
}

.perm-check-item input[type="checkbox"] {
  accent-color: #4dabf7;
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.perm-check-label {
  font-size: 13px;
  color: #8b93a5;
}

.perm-check-label code {
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 11px;
  color: #74c0fc;
  margin-right: 8px;
}

.el-message-box {
  background: #161b26 !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
}

.el-message-box__title {
  color: #e7e9ed !important;
}

.el-message-box__message {
  color: #8b93a5 !important;
}
</style>