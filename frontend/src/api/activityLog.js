import request from '@/utils/request'

export function getActivityLog(accountID, date) {
  return request({
    url: '/activity-log/info/' + accountID + '?date=' + date,
    method: 'get',
  })
}

export function getCurrentLocation(accountID) {
  return request({
    url: '/activity-log/current-location/' + accountID,
    method: 'get',
  })
}
