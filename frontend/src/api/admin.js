import request from '../utils/request'

export function getAdminDashboard() {
  return request({ url: '/admin/dashboard', method: 'get' })
}

export function getBlacklist() {
  return request({ url: '/admin/blacklist', method: 'get' })
}

export function addBlacklist(data) {
  return request({ url: '/admin/blacklist', method: 'post', data })
}

export function removeBlacklist(id) {
  return request({ url: `/admin/blacklist/${id}`, method: 'delete' })
}

export function getAllDevices() {
  return request({ url: '/admin/devices', method: 'get' })
}

export function getPermissions() {
  return request({ url: '/admin/permissions', method: 'get' })
}

export function savePermission(data) {
  return request({ url: '/admin/permissions', method: 'post', data })
}

export function updatePermission(id, data) {
  return request({ url: `/admin/permissions/${id}`, method: 'put', data })
}

export function deletePermission(id) {
  return request({ url: `/admin/permissions/${id}`, method: 'delete' })
}

export function getRolePermissions(roleId) {
  return request({ url: `/admin/roles/permissions/${roleId}`, method: 'get' })
}

export function assignPermissions(roleId, permissionIds) {
  return request({ url: `/admin/roles/${roleId}/permissions`, method: 'put', data: permissionIds })
}

export function getAdminProducts() {
  return request({ url: '/admin/seckill/products', method: 'get' })
}

export function saveProduct(data) {
  return request({ url: '/admin/seckill/products', method: 'post', data })
}

export function deleteProduct(id) {
  return request({ url: `/admin/seckill/products/${id}`, method: 'delete' })
}

export function getAdminActivities() {
  return request({ url: '/admin/seckill/activities', method: 'get' })
}

export function saveActivity(data) {
  return request({ url: '/admin/seckill/activities', method: 'post', data })
}

export function deleteActivity(id) {
  return request({ url: `/admin/seckill/activities/${id}`, method: 'delete' })
}

export function getAdminOrders() {
  return request({ url: '/admin/seckill/orders', method: 'get' })
}

export function getStockLogs() {
  return request({ url: '/admin/seckill/stock-logs', method: 'get' })
}

export function getAdminSessions() {
  return request({ url: '/admin/im/sessions', method: 'get' })
}

export function getAdminMessages() {
  return request({ url: '/admin/im/messages', method: 'get' })
}

export function getAdminGroups() {
  return request({ url: '/admin/im/groups', method: 'get' })
}

export function getPushLogs() {
  return request({ url: '/admin/im/push-logs', method: 'get' })
}