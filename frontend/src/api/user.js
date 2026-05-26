import request from '../utils/request'

export function getUserList(params = {}) {
  return request({
    url: '/users/list',
    method: 'get',
    params
  })
}

export function updateUserStatus(userId, status) {
  return request({
    url: `/users/${userId}/status`,
    method: 'put',
    params: { status }
  })
}

export function deleteUser(userId) {
  return request({
    url: `/users/${userId}`,
    method: 'delete'
  })
}

export function lockUser(userId, lockMinutes = 30) {
  return request({
    url: `/users/${userId}/lock`,
    method: 'post',
    params: { lockMinutes }
  })
}

export function resetPassword(userId, newPassword) {
  return request({
    url: `/users/${userId}/password`,
    method: 'put',
    params: { newPassword }
  })
}

export function getRoleList() {
  return request({
    url: '/roles/list',
    method: 'get'
  })
}

export function createRole(data) {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

export function updateRole(roleId, data) {
  return request({
    url: `/roles/${roleId}`,
    method: 'put',
    data
  })
}

export function deleteRole(roleId) {
  return request({
    url: `/roles/${roleId}`,
    method: 'delete'
  })
}
