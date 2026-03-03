// pages/bills/bills.ts
const app = getApp<IAppOption>()

Page({
  data: {
    stats: {
      pendingAmount: '0.00',
      paidAmount: '0.00'
    },
    bills: []
  },

  onLoad() {
    this.loadBills()
  },

  onShow() {
    this.loadBills()
  },

  // 加载账单列表
  loadBills() {
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
      url: `${app.globalData.baseUrl}/bills`,
      method: 'GET',
      data: { elderly_id: elderlyId },
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success(res: any) {
        if (res.data.code === 0) {
          const bills = res.data.data.list || []
          let pendingAmount = 0
          let paidAmount = 0

          const processedBills = bills.map((bill: any) => {
            const createdDate = new Date(bill.created_at)
            const month = createdDate.getMonth() + 1
            const year = createdDate.getFullYear()

            let statusText = ''
            let statusClass = ''
            let canPay = false

            switch (bill.status) {
              case 'pending':
                statusText = '待支付'
                statusClass = 'pending'
                canPay = true
                pendingAmount += bill.total_amount
                break
              case 'paid':
                statusText = '已支付'
                statusClass = 'paid'
                paidAmount += bill.total_amount
                break
              case 'overdue':
                statusText = '已逾期'
                statusClass = 'overdue'
                canPay = true
                pendingAmount += bill.total_amount
                break
            }

            // 计算到期日期
            const dueDate = createdDate
            dueDate.setDate(dueDate.getDate() + 7)

            return {
              ...bill,
              title: `${year}年${month}月费用`,
              amount: bill.total_amount.toFixed(2),
              period: `${year}.${month.toString().padStart(2, '0')}`,
              statusText,
              statusClass,
              canPay,
              dueDate: `到期日: ${dueDate.getMonth() + 1}月${dueDate.getDate()}日`
            }
          })

          that.setData({
            bills: processedBills,
            stats: {
              pendingAmount: pendingAmount.toFixed(2),
              paidAmount: paidAmount.toFixed(2)
            }
          })
        }
      }
    })
  },

  // 支付账单
  payBill(e: any) {
    const id = e.currentTarget.dataset.id

    wx.showLoading({ title: '处理中...' })

    // 调用微信支付
    const that = this
    wx.request({
      url: `${app.globalData.baseUrl}/bills/${id}/pay`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      data: {
        payment_method: 'wechat'
      },
      success(res: any) {
        wx.hideLoading()
        if (res.data.code === 0 && res.data.data.payment_params) {
          // 调起微信支付
          wx.requestPayment({
            ...res.data.data.payment_params,
            success() {
              wx.showToast({
                title: '支付成功',
                icon: 'success'
              })
              that.loadBills()
            },
            fail(err: any) {
              if (err.errMsg !== 'requestPayment:fail cancel') {
                wx.showToast({
                  title: '支付失败',
                  icon: 'none'
                })
              }
            }
          })
        }
      },
      fail() {
        wx.hideLoading()
        wx.showToast({
          title: '请求失败',
          icon: 'none'
        })
      }
    })
  }
})
