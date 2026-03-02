// Mock 数据 - 用于开发和演示

// 模拟延迟
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

// Mock 用户数据
const mockUser = {
  id: 1,
  phone: '13800138000',
  name: '张先生',
  avatar: '',
  role: 'family'
}

// Mock 老人数据
const mockElderlyList = [
  {
    id: 1,
    name: '张大爷',
    gender: 'male',
    age: 78,
    room: 'A栋 301',
    bed: '1号床',
    careLevel: '二级护理',
    healthStatus: '良好',
    admissionDate: '2023-06-15',
    avatar: ''
  },
  {
    id: 2,
    name: '李奶奶',
    gender: 'female',
    age: 82,
    room: 'B栋 205',
    bed: '2号床',
    careLevel: '一级护理',
    healthStatus: '需关注',
    admissionDate: '2023-03-20',
    avatar: ''
  }
]

// Mock 护理记录
const mockCareRecords = [
  {
    id: 1,
    elderlyId: 1,
    elderlyName: '张大爷',
    itemType: '日常护理',
    content: '协助洗漱、更换衣物',
    caregiver: '护理员小王',
    time: '2024-03-02 08:30',
    images: []
  },
  {
    id: 2,
    elderlyId: 1,
    elderlyName: '张大爷',
    itemType: '用药提醒',
    content: '服用降压药、阿司匹林',
    caregiver: '护理员小李',
    time: '2024-03-02 09:00',
    images: []
  },
  {
    id: 3,
    elderlyId: 1,
    elderlyName: '张大爷',
    itemType: '康复训练',
    content: '进行肢体功能训练30分钟',
    caregiver: '康复师小张',
    time: '2024-03-02 10:30',
    images: []
  },
  {
    id: 4,
    elderlyId: 2,
    elderlyName: '李奶奶',
    itemType: '日常护理',
    content: '协助进食、测量血压',
    caregiver: '护理员小王',
    time: '2024-03-02 07:30',
    images: []
  },
  {
    id: 5,
    elderlyId: 2,
    elderlyName: '李奶奶',
    itemType: '健康检查',
    content: '血压120/80，心率72',
    caregiver: '护士小陈',
    time: '2024-03-02 14:00',
    images: []
  }
]

// Mock 账单数据
const mockBills = [
  {
    id: 1,
    elderlyId: 1,
    elderlyName: '张大爷',
    period: '2024年3月',
    totalAmount: 5800,
    paidAmount: 5800,
    status: 'paid',
    items: [
      { name: '床位费', amount: 2000 },
      { name: '护理费', amount: 2500 },
      { name: '餐费', amount: 1000 },
      { name: '其他费用', amount: 300 }
    ],
    createdAt: '2024-03-01'
  },
  {
    id: 2,
    elderlyId: 1,
    elderlyName: '张大爷',
    period: '2024年2月',
    totalAmount: 5800,
    paidAmount: 5800,
    status: 'paid',
    items: [
      { name: '床位费', amount: 2000 },
      { name: '护理费', amount: 2500 },
      { name: '餐费', amount: 1000 },
      { name: '其他费用', amount: 300 }
    ],
    createdAt: '2024-02-01'
  },
  {
    id: 3,
    elderlyId: 2,
    elderlyName: '李奶奶',
    period: '2024年3月',
    totalAmount: 6500,
    paidAmount: 0,
    status: 'unpaid',
    items: [
      { name: '床位费', amount: 2500 },
      { name: '护理费', amount: 3000 },
      { name: '餐费', amount: 800 },
      { name: '其他费用', amount: 200 }
    ],
    createdAt: '2024-03-01'
  }
]

// Mock API 实现
export const mockApi = {
  auth: {
    login: async (phone: string, password: string) => {
      await delay(500)
      if (phone === '13800138000' && password === '123456') {
        return {
          token: 'mock-jwt-token-' + Date.now(),
          user: mockUser
        }
      }
      throw new Error('手机号或密码错误')
    },
    wechatLogin: async () => {
      await delay(500)
      return {
        token: 'mock-jwt-token-' + Date.now(),
        user: mockUser
      }
    }
  },

  elderly: {
    list: async () => {
      await delay(300)
      return { list: mockElderlyList, total: mockElderlyList.length }
    },
    get: async (id: number) => {
      await delay(300)
      const elderly = mockElderlyList.find(e => e.id === id)
      if (!elderly) throw new Error('未找到老人信息')
      return elderly
    },
    familyList: async () => {
      await delay(300)
      return mockElderlyList
    }
  },

  care: {
    records: async (params: any) => {
      await delay(300)
      let records = mockCareRecords
      if (params.elderlyId) {
        records = records.filter(r => r.elderlyId === params.elderlyId)
      }
      return { list: records, total: records.length }
    },
    items: async () => {
      await delay(200)
      return [
        { id: 1, name: '日常护理' },
        { id: 2, name: '用药提醒' },
        { id: 3, name: '康复训练' },
        { id: 4, name: '健康检查' }
      ]
    }
  },

  bill: {
    list: async (params: any) => {
      await delay(300)
      let bills = mockBills
      if (params.elderlyId) {
        bills = bills.filter(b => b.elderlyId === params.elderlyId)
      }
      return { list: bills, total: bills.length }
    },
    get: async (id: number) => {
      await delay(300)
      const bill = mockBills.find(b => b.id === id)
      if (!bill) throw new Error('未找到账单')
      return bill
    }
  }
}
