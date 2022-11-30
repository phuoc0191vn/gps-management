import request from '@/utils/request'

export function datatable(page, limit) {
  return request({
    url: '/report?output=datatable&page=' + page + "&limit=" + limit,
    method: 'get',
  })
}

export function generateReport(id, startTime, endTime) {
  return request({
    url: '/report/generate/' + id + '?startTime=' + startTime + '&endTime=' + endTime,
    method: 'get',
  })
}

export function downloadReport(id) {
  return request({
    url: '/report/download/' + id,
    method: 'get',
  })
}

export function deleteReport(id) {
  return request({
    url: '/report/delete/' + id,
    method: 'get',
  })
}

