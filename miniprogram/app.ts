// app.ts
App({
  globalData: {
    userInfo: null,
    token: null,
    baseUrl: 'http://localhost:8080/api'
  },

  onLaunch() {
    // 检查登录状态
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
    }

    // 检查更新
    this.checkUpdate()
  },

  // 检查小程序更新
  checkUpdate() {
    if (wx.canIUse('getUpdateManager')) {
      const updateManager = wx.getUpdateManager()

      updateManager.onCheckForUpdate((res) => {
        if (res.hasUpdate) {
          updateManager.onUpdateReady(() => {
            wx.showModal({
              title: '更新提示',
              content: '新版本已准备好，是否重启应用？',
              success: (res) => {
                if (res.confirm) {
                  updateManager.applyUpdate()
                }
              }
            })
          })

          updateManager.onUpdateFailed(() => {
            wx.showModal({
              title: '更新失败',
              content: '新版本下载失败，请检查网络',
              showCancel: false
            })
          })
        }
      })
    }
  },

  // API 请求封装
  request(options: any) {
    const { url, method = 'GET', data = {}, success, fail } = options

    wx.request({
      url: this.globalData.baseUrl + url,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': this.globalData.token ? `Bearer ${this.globalData.token}` : ''
      },
      success: (res: any) => {
        if (res.statusCode === 200) {
          if (res.data.code === 0) {
            success && success(res.data.data)
          } else {
            wx.showToast({
              title: res.data.message || '请求失败',
              icon: 'none'
            })
            fail && fail(res.data)
          }
        } else if (res.statusCode === 401) {
          // Token 失效，重新登录
          this.logout()
          wx.navigateTo({
            url: '/pages/login/login'
          })
        } else {
          wx.showToast({
            title: '网络错误',
            icon: 'none'
          })
          fail && fail(res)
        }
      },
      fail: (err: any) => {
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
        fail && fail(err)
      }
    })
  },

  // 保存 Token
  setToken(token: string) {
    this.globalData.token = token
    wx.setStorageSync('token', token)
  },

  // 退出登录
  logout() {
    this.globalData.token = null
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
  }
})
