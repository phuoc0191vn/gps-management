import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

export function addUser(data) {
  return request({
    url: '/account',
    method: 'post',
    data
  })
}

export function updateUser(id, data) {
  return request({
    url: '/account/' + id,
    method: 'patch',
    data
  })
}

export function resetUser(id) {
  return request({
    url: '/account/reset/' + id,
    method: 'get',
  })
}

export function deleteUser(id) {
  return request({
    url: '/account/' + id,
    method: 'delete',
  })
}

export function userInfo() {
  return request({
    url: '/account/info',
    method: 'get',
  })
}

export function userDetail(id) {
  return request({
    url: '/account/detail/' + id,
    method: 'get',
  })
}

export function datatable(page, limit) {
  return request({
    url: '/account?output=datatable&page=' + page + "&limit=" + limit,
    method: 'get',
  })
}

export function getInfo(token) {
  return request({
    url: '/vue-admin-template/user/info',
    method: 'get',
    params: {token}
  })
}

export function logout() {
  return request({
    url: '/vue-admin-template/user/logout',
    method: 'post'
  })
}
