import request from '@/utils/request'


export function getEnableDevice() {
  return request({
    url: '/device/status?status=1',
    method: 'get',
  })
}

export function addDevice(data) {
  return request({
    url: '/device',
    method: 'post',
    data
  })
}

export function getDetail(id) {
  return request({
    url: '/device/detail/' + id,
    method: 'get',
  })
}

export function toggleStatus(id, status) {
  return request({
    url: '/device/toggle-status/' + id + '?status=' + status,
    method: 'get',
  })
}

export function updateDevice(id, data) {
  return request({
    url: '/device/' + id,
    method: 'patch',
    data
  })
}

export function deleteDevice(id) {
  return request({
    url: '/device/' + id,
    method: 'delete',
  })
}

export function datatable(page, limit) {
  return request({
    url: '/device?output=datatable&page=' + page + "&limit=" + limit,
    method: 'get',
  })
}
