// pages/service/call/call.ts
const app = getApp<IAppOption>()

Page({
  data: {
    elderlyList: [] as any[],
    selectedElderly: null as any,
    selectedType: '',
    notes: '',
    submitting: false,
    canSubmit: false,
    historyList: [] as any[],
    serviceTypes: [
      { name: '护理帮助', value: '护理', icon: '/images/service-care.png' },
      { name: '送餐服务', value: '送餐', icon: '/images/service-meal.png' },
      { name: '打扫卫生', value: '打扫', icon: '/images/service-clean.png' },
      { name: '紧急呼叫', value: '紧急', icon: '/images/service-urgent.png' },
      { name: '医疗协助', value: '医疗', icon: '/images/service-medical.png' },
      { name: '其他帮助', value: '其他', icon: '/images/service-other.png' }
    ],
    quickNotes: {
      '护理': ['需要翻身', '需要换药', '需要喂饭', '需要洗澡协助'],
      '送餐': ['早餐', '午餐', '晚餐', '加餐'],
      '打扫': ['房间打扫', '更换床单', '清理卫生间', '倒垃圾'],
      '紧急': ['身体不适', '摔倒', '呼吸困难', '其他紧急情况'],
      '医疗': ['测量血压', '测量血糖', '服药提醒', '就医陪同'],
      '其他': ['外出陪同', '购物帮助', '设备故障', '其他需求']
    }
  },

  onLoad() {
    this.loadElderlyList()
    this.loadHistory()
  },

  onShow() {
    this.loadHistory()
  },

  // 更新提交按钮状态
  updateCanSubmit() {
    const canSubmit = !!(this.data.selectedElderly && this.data.selectedType)
    this.setData({ canSubmit })
  },

  // 加载老人列表
  loadElderlyList() {
    app.request({
      url: '/user/elderly-list',
      method: 'GET',
      success: (res: any) => {
        this.setData({
          elderlyList: res || []
        })
        if (res && res.length > 0) {
          this.setData({
            selectedElderly: res[0]
          }, () => {
            this.updateCanSubmit()
          })
        }
      }
    })
  },

  // 加载历史记录
  loadHistory() {
    app.request({
      url: '/service/requests?page=1&page_size=10',
      method: 'GET',
      success: (res: any) => {
        this.setData({
          historyList: res.list || []
        })
      }
    })
  },

  // 选择老人
  onElderlyChange(e: any) {
    const index = e.detail.value
    this.setData({
      selectedElderly: this.data.elderlyList[index]
    }, () => {
      this.updateCanSubmit()
    })
  },

  // 选择服务类型
  selectType(e: any) {
    const value = e.currentTarget.dataset.value
    this.setData({
      selectedType: value
    }, () => {
      this.updateCanSubmit()
    })
  },

  // 输入备注
  onNotesInput(e: any) {
    this.setData({
      notes: e.detail.value
    })
  },

  // 追加快捷备注
  appendNote(e: any) {
    const note = e.currentTarget.dataset.note
    const currentNotes = this.data.notes
    const newNotes = currentNotes ? `${currentNotes}，${note}` : note
    this.setData({
      notes: newNotes
    })
  },

  // 提交呼叫
  submitCall() {
    if (!this.data.canSubmit) return

    this.setData({ submitting: true })

    app.request({
      url: '/service/requests',
      method: 'POST',
      data: {
        elderly_id: this.data.selectedElderly.id,
        type: this.data.selectedType,
        notes: this.data.notes
      },
      success: () => {
        wx.showToast({
          title: '呼叫已发送',
          icon: 'success'
        })
        this.setData({
          selectedType: '',
          notes: ''
        }, () => {
          this.updateCanSubmit()
        })
        this.loadHistory()
      },
      fail: () => {
        wx.showToast({
          title: '发送失败',
          icon: 'error'
        })
      },
      complete: () => {
        this.setData({ submitting: false })
      }
    })
  },

  // 获取状态文本
  getStatusText(status: string): string {
    const map: Record<string, string> = {
      pending: '待处理',
      processing: '处理中',
      completed: '已完成'
    }
    return map[status] || status
  },

  // 格式化时间
  formatTime(time: string): string {
    if (!time) return ''
    const date = new Date(time)
    const now = new Date()
    const diff = now.getTime() - date.getTime()

    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`

    return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  }
})
