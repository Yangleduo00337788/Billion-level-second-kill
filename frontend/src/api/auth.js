import request from '../utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data: {
      username: data.username,
      password: data.password,
      deviceId: data.deviceId || 'web-browser',
      deviceType: 'WEB',
      loginType: 1
    }
  })
}

export function register(data) {
  return request({
    url: '/auth/register',
    method: 'post',
    data: {
      username: data.username,
      password: data.password,
      confirmPassword: data.confirmPassword,
      phone: data.phone,
      email: data.email
    }
  })
}

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

export function refreshToken(refreshToken) {
  return request({
    url: '/auth/refresh-token',
    method: 'post',
    data: {
      refreshToken: refreshToken,
      deviceId: 'web-browser'
    }
  })
}
