import request from '../utils/request'

export function getSeckillProducts() {
  return request({
    url: '/seckill/products',
    method: 'get'
  })
}

export function getSeckillProduct(seckillId) {
  return request({
    url: `/seckill/product/${seckillId}`,
    method: 'get'
  })
}

export function executeSeckill(seckillId) {
  return request({
    url: `/seckill/execute/${seckillId}`,
    method: 'post'
  })
}

export function getStock(seckillId) {
  return request({
    url: `/seckill/stock/${seckillId}`,
    method: 'get'
  })
}

export function getOrderResult(seckillId) {
  return request({
    url: `/seckill/order/result/${seckillId}`,
    method: 'get'
  })
}

export function getMyOrders() {
  return request({
    url: '/seckill/my-orders',
    method: 'get'
  })
}

export function payOrder(orderId) {
  return request({
    url: `/seckill/pay/${orderId}`,
    method: 'post'
  })
}

export function cancelOrder(orderId) {
  return request({
    url: `/seckill/cancel/${orderId}`,
    method: 'post'
  })
}