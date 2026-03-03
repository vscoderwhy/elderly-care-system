// pages/health/health.ts
const app = getApp<IAppOption>()

Page({
  data: {
    healthData: [],
    activeTab: 'week',
    filterType: '',
    records: [],
    page: 1,
    hasMore: true
  },

  onLoad() {
    this.loadHealthData()
    this.loadRecords()
  },

  // 加载最新健康数据
  loadHealthData() {
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
      url: `${app.globalData.baseUrl}/health/records/latest/${elderlyId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0 && res.data.data) {
          const healthMap: Record<string, any> = {
            'blood_pressure': { name: '血压', unit: 'mmHg', format: (v: string) => v },
            'blood_sugar': { name: '血糖', unit: 'mmol/L', format: (v: string) => v },
            'temperature': { name: '体温', unit: '℃', format: (v: string) => v },
            'weight': { name: '体重', unit: 'kg', format: (v: string) => v },
            'heart_rate': { name: '心率', unit: '次/分', format: (v: string) => v }
          }

          const healthData = res.data.data.map((item: any) => {
            const config = healthMap[item.record_type] || { name: item.record_type, unit: '', format: (v: string) => v }
            const value = config.format(item.value)

            // 判断状态
            let status = 'normal'
            let statusText = '正常'
            if (item.record_type === 'blood_pressure') {
              const systolic = parseInt(value.split('/')[0])
              if (systolic >= 140) { status = 'warning'; statusText = '偏高' }
              else if (systolic < 90) { status = 'warning'; statusText = '偏低' }
            } else if (item.record_type === 'blood_sugar') {
              const sugar = parseFloat(value)
              if (sugar > 7.0) { status = 'warning'; statusText = '偏高' }
              else if (sugar < 3.9) { status = 'warning'; statusText = '偏低' }
            } else if (item.record_type === 'temperature') {
              const temp = parseFloat(value)
              if (temp > 37.3) { status = 'warning'; statusText = '发热' }
              else if (temp < 36.0) { status = 'warning'; statusText = '偏低' }
            }

            return {
              ...item,
              name: config.name,
              value,
              unit: config.unit,
              time: item.recorded_at ? item.recorded_at.substring(5, 16) : '',
              status,
              statusText
            }
          })

          that.setData({ healthData })
        }
      }
    })
  },

  // 加载健康记录
  loadRecords() {
    const that = this
    const elderlyId = app.globalData.elderlyInfo?.id

    wx.request({
      url: `${app.globalData.baseUrl}/health/records`,
      method: 'GET',
      data: {
        elderly_id: elderlyId,
        record_type: that.data.filterType,
        page: that.data.page,
        page_size: 20
      },
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0) {
          const typeNames: Record<string, string> = {
            'blood_pressure': '血压',
            'blood_sugar': '血糖',
            'temperature': '体温',
            'weight': '体重',
            'heart_rate': '心率'
          }

          const newRecords = res.data.data.list.map((item: any) => ({
            ...item,
            typeName: typeNames[item.record_type] || item.record_type,
            time: item.recorded_at ? item.recorded_at.substring(0, 16) : '',
            recorder: item.recorder?.nickname || '系统'
          }))

          that.setData({
            records: that.data.page === 1 ? newRecords : [...that.data.records, ...newRecords],
            hasMore: newRecords.length >= 20
          })
        }
      }
    })
  },

  // 切换时间范围
  switchTab(e: any) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ activeTab: tab })
    // 重新加载图表数据
  },

  // 筛选记录类型
  filterRecords(e: any) {
    const type = e.currentTarget.dataset.type
    this.setData({ filterType: type, page: 1 })
    this.loadRecords()
  },

  // 加载更多
  loadMore() {
    this.setData({ page: this.data.page + 1 })
    this.loadRecords()
  }
})
