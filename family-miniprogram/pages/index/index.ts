// pages/index/index.ts
const app = getApp<IAppOption>()

Page({
  data: {
    elderlyInfo: null,
    todayCare: [],
    notices: []
  },

  onLoad() {
    this.loadData()
  },

  onShow() {
    // 每次显示页面刷新数据
    this.loadTodayCare()
  },

  // 加载数据
  loadData() {
    const elderlyInfo = app.globalData.elderlyInfo
    if (elderlyInfo) {
      this.setData({ elderlyInfo })
    }
    this.loadTodayCare()
    this.loadNotices()
  },

  // 加载今日护理
  loadTodayCare() {
    const that = this
    wx.request({
      url: `${app.globalData.baseUrl}/care/my-tasks`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0 && res.data.data) {
          const careList = res.data.data.map((item: any) => ({
            id: item.id,
            time: item.recorded_at ? item.recorded_at.substring(11, 16) : '',
            name: item.care_item?.name || '护理',
            status: item.completed_at ? 'done' : 'pending',
            statusText: item.completed_at ? '已完成' : '待执行'
          }))
          that.setData({ todayCare: careList })
        }
      }
    })
  },

  // 加载通知
  loadNotices() {
    // 模拟通知数据
    this.setData({
      notices: [
        { id: 1, title: '定期体检通知', time: '2小时前', type: 'info' },
        { id: 2, title: '探视预约确认', time: '昨天', type: 'success' }
      ]
    })
  },

  // 跳转到健康页面
  goToHealth() {
    wx.switchTab({
      url: '/pages/health/health'
    })
  },

  // 跳转到探视页面
  goToVisits() {
    wx.switchTab({
      url: '/pages/visits/visits'
    })
  },

  // 跳转到账单页面
  goToBills() {
    wx.switchTab({
      url: '/pages/bills/bills'
    })
  },

  // 呼叫服务
  makeCall() {
    wx.showModal({
      title: '呼叫服务',
      content: '确认呼叫护理服务？',
      success(res: any) {
        if (res.confirm) {
          // 创建服务请求
          wx.request({
            url: `${app.globalData.baseUrl}/service/requests`,
            method: 'POST',
            header: {
              'Authorization': `Bearer ${app.globalData.token}`
            },
            data: {
              elderly_id: app.globalData.elderlyInfo?.id,
              request_type: 'family_call',
              description: '家属呼叫服务'
            },
            success(res: any) {
              if (res.data.code === 0) {
                wx.showToast({
                  title: '呼叫成功',
                  icon: 'success'
                })
              }
            }
          })
        }
      }
    })
  },

  // 查看通知详情
  viewNotice(e: any) {
    const id = e.currentTarget.dataset.id
    wx.showToast({
      title: '通知详情功能开发中',
      icon: 'none'
    })
  }
})
