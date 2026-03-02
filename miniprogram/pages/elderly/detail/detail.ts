// pages/elderly/detail/detail.ts
const app = getApp<IAppOption>()

Page({
  data: {
    elderlyId: 0,
    elderly: {} as any,
    careRecords: [] as any[],
    healthRecords: [] as any[],
    healthHistory: [] as any[],
    showHealthDialog: false,
    selectedHealthType: ''
  },

  onLoad(options: any) {
    this.setData({
      elderlyId: options.id
    })
    this.loadElderlyDetail()
    this.loadCareRecords()
    this.loadHealthRecords()
  },

  // 加载老人详情
  loadElderlyDetail() {
    app.request({
      url: `/elderly/${this.data.elderlyId}`,
      method: 'GET',
      success: (res: any) => {
        // 计算年龄
        if (res.birth_date) {
          const birth = new Date(res.birth_date)
          const today = new Date()
          let age = today.getFullYear() - birth.getFullYear()
          const monthDiff = today.getMonth() - birth.getMonth()
          if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
            age--
          }
          res.age = age
        }
        this.setData({
          elderly: res
        })
      }
    })
  },

  // 加载护理记录
  loadCareRecords() {
    app.request({
      url: `/care/records?elderly_id=${this.data.elderlyId}&page=1&page_size=10`,
      method: 'GET',
      success: (res: any) => {
        this.setData({
          careRecords: res.list || []
        })
      }
    })
  },

  // 加载健康记录
  loadHealthRecords() {
    app.request({
      url: `/health/records/latest/${this.data.elderlyId}`,
      method: 'GET',
      success: (res: any) => {
        this.setData({
          healthRecords: res || []
        })
      }
    })
  },

  // 获取健康值显示
  getHealthValue(record: any): string {
    if (record.record_type === 'blood_pressure' && record.value2) {
      return `${record.value}/${record.value2}`
    }
    return record.value
  },

  // 获取健康标签
  getHealthLabel(type: string): string {
    const labels: Record<string, string> = {
      blood_pressure: '血压',
      blood_sugar: '血糖',
      temperature: '体温',
      weight: '体重',
      heart_rate: '心率'
    }
    return labels[type] || type
  },

  // 获取健康单位
  getHealthUnit(type: string): string {
    const units: Record<string, string> = {
      blood_pressure: 'mmHg',
      blood_sugar: 'mmol/L',
      temperature: '℃',
      weight: 'kg',
      heart_rate: '次/分'
    }
    return units[type] || ''
  },

  // 显示健康历史
  showHealthHistory(e: any) {
    const type = e.currentTarget.dataset.type
    this.setData({
      selectedHealthType: type,
      showHealthDialog: true
    })
    this.loadHealthHistory(type)
  },

  // 加载健康历史记录
  loadHealthHistory(type: string) {
    app.request({
      url: `/health/records?elderly_id=${this.data.elderlyId}&record_type=${type}&page=1&page_size=10`,
      method: 'GET',
      success: (res: any) => {
        this.setData({
          healthHistory: res.list || []
        })
      }
    })
  },

  // 关闭健康历史弹窗
  closeHealthDialog() {
    this.setData({
      showHealthDialog: false
    })
  },

  // 查看全部记录
  viewAllRecords() {
    wx.navigateTo({
      url: '/pages/elderly/care-records/care-records?elderlyId=' + this.data.elderlyId
    })
  },

  // 预览图片
  previewImage(e: any) {
    const url = e.currentTarget.dataset.url
    const urls = e.currentTarget.dataset.urls
    wx.previewImage({
      current: url,
      urls: urls
    })
  },

  // 快速呼叫
  quickCall(e: any) {
    const type = e.currentTarget.dataset.type
    app.request({
      url: '/service/requests',
      method: 'POST',
      data: {
        elderly_id: this.data.elderlyId,
        type: type,
        notes: ''
      },
      success: () => {
        wx.showToast({
          title: '呼叫已发送',
          icon: 'success'
        })
      },
      fail: () => {
        wx.showToast({
          title: '发送失败',
          icon: 'error'
        })
      }
    })
  },

  // 跳转到服务呼叫页
  goToServiceCall() {
    wx.navigateTo({
      url: '/pages/service/call/call'
    })
  }
})
