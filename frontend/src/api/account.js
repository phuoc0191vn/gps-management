import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

export function addAccount(data) {
  return request({
    url: '/account',
    method: 'post',
    data
  })
}

export function updateAccount(id, data) {
  return request({
    url: '/account/' + id,
    method: 'patch',
    data
  })
}

export function resetAccount(id) {
  return request({
    url: '/account/reset/' + id,
    method: 'get',
  })
}

export function deleteAccount(id) {
  return request({
    url: '/account/' + id,
    method: 'delete',
  })
}

export function accountInfo() {
  return request({
    url: '/account/info',
    method: 'get',
  })
}

export function accountDetail(id) {
  return request({
    url: '/account/detail/' + id,
    method: 'get',
  })
}

export function getChildAccounts() {
  return request({
    url: '/account/child-accounts',
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
