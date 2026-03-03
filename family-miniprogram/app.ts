// app.ts
App<IAppOption>({
  globalData: {
    userInfo: null,
    token: null,
    elderlyInfo: null,
    baseUrl: 'http://1.12.223.138:8080/api'
  },

  onLaunch() {
    // 检查登录状态
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
      // 获取用户信息
      this.getUserInfo()
    }
  },

  // 获取用户信息
  getUserInfo() {
    const that = this
    wx.request({
      url: `${this.globalData.baseUrl}/user/profile`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${this.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0) {
          that.globalData.userInfo = res.data.data
          // 获取关联的老人信息
          that.getElderlyInfo()
        }
      }
    })
  },

  // 获取关联的老人信息
  getElderlyInfo() {
    const that = this
    wx.request({
      url: `${this.globalData.baseUrl}/user/elderly-list`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${this.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0 && res.data.data.length > 0) {
          that.globalData.elderlyInfo = res.data.data[0]
        }
      }
    })
  }
})
