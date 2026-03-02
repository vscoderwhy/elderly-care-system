// pages/login/login.ts
const app = getApp<IAppOption>()

Page({
  data: {
    phone: '',
    password: '',
    showPassword: false,
    loading: false
  },

  // 手机号输入
  onPhoneInput(e: any) {
    this.setData({
      phone: e.detail.value
    })
  },

  // 密码输入
  onPasswordInput(e: any) {
    this.setData({
      password: e.detail.value
    })
  },

  // 切换密码显示
  togglePassword() {
    this.setData({
      showPassword: !this.data.showPassword
    })
  },

  // 手机号登录
  handleLogin() {
    const { phone, password } = this.data

    if (!phone) {
      wx.showToast({
        title: '请输入手机号',
        icon: 'none'
      })
      return
    }

    if (!/^1\d{10}$/.test(phone)) {
      wx.showToast({
        title: '手机号格式不正确',
        icon: 'none'
      })
      return
    }

    if (!password) {
      wx.showToast({
        title: '请输入密码',
        icon: 'none'
      })
      return
    }

    this.setData({ loading: true })

    app.request({
      url: '/auth/login',
      method: 'POST',
      data: { phone, password },
      success: (res: any) => {
        const { token, user } = res
        app.setToken(token)
        app.globalData.userInfo = user
        wx.setStorageSync('userInfo', user)

        wx.showToast({
          title: '登录成功',
          icon: 'success'
        })

        setTimeout(() => {
          wx.switchTab({
            url: '/pages/index/index'
          })
        }, 1500)
      },
      fail: () => {
        this.setData({ loading: false })
      }
    })
  },

  // 微信登录
  handleWeChatLogin(e: any) {
    if (e.detail.errMsg !== 'getUserInfo:ok') {
      wx.showToast({
        title: '需要授权才能登录',
        icon: 'none'
      })
      return
    }

    wx.login({
      success: (res) => {
        if (res.code) {
          this.setData({ loading: true })

          app.request({
            url: '/auth/wechat-login',
            method: 'POST',
            data: {
              code: res.code,
              ...e.detail.userInfo
            },
            success: (result: any) => {
              const { token, user } = result
              app.setToken(token)
              app.globalData.userInfo = user
              wx.setStorageSync('userInfo', user)

              wx.showToast({
                title: '登录成功',
                icon: 'success'
              })

              setTimeout(() => {
                wx.switchTab({
                  url: '/pages/index/index'
                })
              }, 1500)
            },
            fail: () => {
              this.setData({ loading: false })
            }
          })
        }
      }
    })
  }
})
