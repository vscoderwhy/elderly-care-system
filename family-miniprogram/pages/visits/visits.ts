// pages/visits/visits.ts
const app = getApp<IAppOption>()

Page({
  data: {
    upcomingVisits: [],
    visitList: []
  },

  onLoad() {
    this.loadVisits()
  },

  onShow() {
    this.loadVisits()
  },

  // 加载探视记录
  loadVisits() {
    const that = this
    const elderlyId = app.globalData.elderlyInfo?.id

    if (!elderlyId) {
      wx.showToast({
        title: '未关联老人信息',
        icon: 'none'
      })
      return
    }

    wx.request({
      url: `${app.globalData.baseUrl}/elderly/${elderlyId}/visits`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0 && res.data.data) {
          const visits = res.data.data
          const now = new Date()

          // 处理探视数据
          const processedVisits = visits.map((visit: any) => {
            const visitDate = new Date(visit.scheduled_date)
            const dateStr = visit.scheduled_date.substring(0, 10)
            const month = visitDate.getMonth() + 1
            const day = visitDate.getDate()

            let statusText = ''
            let statusClass = ''
            let canCancel = false

            switch (visit.status) {
              case 'pending':
                statusText = '待确认'
                statusClass = 'pending'
                canCancel = true
                break
              case 'confirmed':
                statusText = '已确认'
                statusClass = 'confirmed'
                canCancel = true
                break
              case 'completed':
                statusText = '已完成'
                statusClass = 'completed'
                break
              case 'cancelled':
                statusText = '已取消'
                statusClass = 'cancelled'
                break
            }

            return {
              ...visit,
              date: dateStr,
              month,
              day,
              timeRange: `${visit.start_time} - ${visit.end_time}`,
              statusText,
              statusClass,
              canCancel
            }
          })

          // 分离即将到来和已完成的探视
          const upcoming = processedVisits.filter((v: any) => {
            const visitDate = new Date(v.scheduled_date)
            return visitDate >= now && ['pending', 'confirmed'].includes(v.status)
          }).sort((a: any, b: any) => new Date(a.scheduled_date) - new Date(b.scheduled_date))

          const list = processedVisits.sort((a: any, b: any) => new Date(b.scheduled_date) - new Date(a.scheduled_date))

          that.setData({
            upcomingVisits: upcoming,
            visitList: list
          })
        }
      }
    })
  },

  // 去预约
  goToBook() {
    wx.navigateTo({
      url: '/pages/visit-book/visit-book'
    })
  },

  // 取消探视
  cancelVisit(e: any) {
    const id = e.currentTarget.dataset.id

    wx.showModal({
      title: '取消探视',
      content: '确认取消本次探视预约？',
      success(res: any) {
        if (res.confirm) {
          const that = this
          wx.request({
            url: `${app.globalData.baseUrl}/visits/${id}/cancel`,
            method: 'PUT',
            header: {
              'Authorization': `Bearer ${app.globalData.token}`
            },
            success(res: any) {
              if (res.data.code === 0) {
                wx.showToast({
                  title: '取消成功',
                  icon: 'success'
                })
                that.loadVisits()
              }
            }
          })
        }
      }
    })
  }
})
