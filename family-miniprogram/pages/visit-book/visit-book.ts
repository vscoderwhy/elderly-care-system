// pages/visit-book/visit-book.ts
const app = getApp<IAppOption>()

Page({
  data: {
    today: '',
    relations: ['配偶', '子女', '孙子女', '亲属', '朋友', '其他'],
    form: {
      date: '',
      startTime: '09:00',
      endTime: '10:00',
      visitors: 1,
      visitorName: '',
      phone: '',
      relationIndex: -1,
      notes: ''
    },
    submitting: false
  },

  onLoad() {
    // 设置今天日期
    const today = new Date()
    this.setData({
      today: `${today.getFullYear()}-${(today.getMonth() + 1).toString().padStart(2, '0')}-${today.getDate().toString().padStart(2, '0')}`
    })
  },

  // 日期选择
  onDateChange(e: any) {
    this.setData({
      'form.date': e.detail.value
    })
  },

  // 开始时间
  onStartTimeChange(e: any) {
    this.setData({
      'form.startTime': e.detail.value
    })
  },

  // 结束时间
  onEndTimeChange(e: any) {
    this.setData({
      'form.endTime': e.detail.value
    })
  },

  // 人数输入
  onVisitorsInput(e: any) {
    this.setData({
      'form.visitors': e.detail.value
    })
  },

  // 姓名输入
  onNameInput(e: any) {
    this.setData({
      'form.visitorName': e.detail.value
    })
  },

  // 电话输入
  onPhoneInput(e: any) {
    this.setData({
      'form.phone': e.detail.value
    })
  },

  // 关系选择
  onRelationChange(e: any) {
    this.setData({
      'form.relationIndex': e.detail.value
    })
  },

  // 备注输入
  onNotesInput(e: any) {
    this.setData({
      'form.notes': e.detail.value
    })
  },

  // 提交预约
  submitBook(e: any) {
    const that = this
    const { form } = this.data
    const elderlyId = app.globalData.elderlyInfo?.id

    // 表单验证
    if (!form.date) {
      wx.showToast({
        title: '请选择探视日期',
        icon: 'none'
      })
      return
    }

    if (!form.startTime || !form.endTime) {
      wx.showToast({
        title: '请选择探视时间',
        icon: 'none'
      })
      return
    }

    if (form.startTime >= form.endTime) {
      wx.showToast({
        title: '结束时间需大于开始时间',
        icon: 'none'
      })
      return
    }

    if (!form.visitorName) {
      wx.showToast({
        title: '请输入探视人姓名',
        icon: 'none'
      })
      return
    }

    if (!form.phone || !/^1[3-9]\d{9}$/.test(form.phone)) {
      wx.showToast({
        title: '请输入正确的手机号',
        icon: 'none'
      })
      return
    }

    if (form.relationIndex < 0) {
      wx.showToast({
        title: '请选择与老人关系',
        icon: 'none'
      })
      return
    }

    this.setData({ submitting: true })

    // 提交数据
    wx.request({
      url: `${app.globalData.baseUrl}/visits`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      data: {
        elderly_id: elderlyId,
        scheduled_date: form.date,
        start_time: form.startTime,
        end_time: form.endTime,
        visitor_name: form.visitorName,
        visitor_phone: form.phone,
        visitor_count: form.visitors,
        relation: that.data.relations[form.relationIndex],
        notes: form.notes
      },
      success(res: any) {
        if (res.data.code === 0) {
          wx.showToast({
            title: '预约成功',
            icon: 'success'
          })
          setTimeout(() => {
            wx.navigateBack()
          }, 1500)
        } else {
          wx.showToast({
            title: res.data.message || '预约失败',
            icon: 'none'
          })
        }
      },
      fail() {
        wx.showToast({
          title: '网络请求失败',
          icon: 'none'
        })
      },
      complete() {
        that.setData({ submitting: false })
      }
    })
  }
})
