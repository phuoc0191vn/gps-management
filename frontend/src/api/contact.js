import request from '@/utils/request'

export function submitContact(data) {
  return request({
    url: '/contact',
    method: 'post',
    data
  })
}

export function submitDoneContact(id) {
  return request({
    url: '/contact/done/' + id,
    method: 'get',
  })
}

export function deleteContact(id) {
  return request({
    url: '/contact/' + id,
    method: 'delete',
  })
}

export function datatable(page, limit) {
  return request({
    url: '/contact?output=datatable&page=' + page + "&limit=" + limit,
    method: 'get',
  })
}
