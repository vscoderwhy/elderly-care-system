// pages/index/index.ts
const app = getApp<IAppOption>()

Page({
  data: {
    userInfo: null,
    elderlyList: [],
    recentCare: []
  },

  onLoad() {
    this.checkLogin()
  },

  onShow() {
    if (app.globalData.token) {
      this.loadData()
    }
  },

  // 检查登录状态
  checkLogin() {
    if (!app.globalData.token) {
      wx.navigateTo({
        url: '/pages/login/login'
      })
      return
    }
    this.setData({
      userInfo: app.globalData.userInfo
    })
    this.loadData()
  },

  // 加载数据
  loadData() {
    this.loadElderlyList()
    this.loadRecentCare()
  },

  // 加载老人列表
  loadElderlyList() {
    app.request({
      url: '/user/elderly-list',
      method: 'GET',
      success: (res: any) => {
        this.setData({
          elderlyList: res
        })
      }
    })
  },

  // 加载最新护理记录
  loadRecentCare() {
    if (this.data.elderlyList.length === 0) return

    const elderlyId = this.data.elderlyList[0].id
    app.request({
      url: `/care/records?elderly_id=${elderlyId}&page=1&page_size=5`,
      method: 'GET',
      success: (res: any) => {
        this.setData({
          recentCare: res.list || []
        })
      }
    })
  },

  // 跳转到老人详情
  goToElderly(e: any) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/elderly/detail/detail?id=${id}`
    })
  },

  // 跳转到护理记录
  goToCareRecords() {
    if (this.data.elderlyList.length === 0) {
      wx.showToast({
        title: '暂无关联老人',
        icon: 'none'
      })
      return
    }
    wx.navigateTo({
      url: '/pages/elderly/care-records/care-records'
    })
  },

  // 跳转到账单
  goToBills() {
    wx.switchTab({
      url: '/pages/bills/list/list'
    })
  },

  // 跳转到服务呼叫
  goToServiceCall() {
    if (this.data.elderlyList.length === 0) {
      wx.showToast({
        title: '暂无关联老人',
        icon: 'none'
      })
      return
    }
    wx.navigateTo({
      url: '/pages/service/call/call'
    })
  },

  // 联系客服
  contact() {
    wx.makePhoneCall({
      phoneNumber: '400-123-4567'
    })
  }
})
