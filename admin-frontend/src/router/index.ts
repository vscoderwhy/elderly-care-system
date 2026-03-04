import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  // 指挥中心大屏 - 独立路由（不在Layout内）
  {
    path: '/dashboard/command',
    name: 'CommandCenter',
    component: () => import('@/views/Dashboard/CommandCenter.vue'),
    meta: {
      title: '指挥中心',
      requiresAuth: true,
      fullScreen: true
    }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard/overview',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        redirect: '/dashboard/overview'
      },
      {
        path: '/dashboard/overview',
        name: 'DashboardOverview',
        component: () => import('@/views/Dashboard/Overview.vue'),
        meta: {
          title: '概览仪表盘',
          requiresAuth: true
        }
      },
      {
        path: '/statistics/advanced',
        name: 'StatisticsAdvanced',
        component: () => import('@/views/Statistics/Advanced.vue'),
        meta: {
          title: '高级数据分析',
          requiresAuth: true
        }
      },
      {
        path: '/statistics/report',
        name: 'StatisticsReport',
        component: () => import('@/views/Statistics/Report.vue'),
        meta: {
          title: '报表中心',
          requiresAuth: true
        }
      },
      {
        path: '/monitor/realtime',
        name: 'MonitorRealtime',
        component: () => import('@/views/Monitor/Realtime.vue'),
        meta: {
          title: '实时监控',
          requiresAuth: true
        }
      },
      {
        path: '/elderly',
        name: 'Elderly',
        redirect: '/elderly/list',
        meta: {
          title: '老人管理',
          requiresAuth: true
        }
      },
      {
        path: '/elderly/list',
        name: 'ElderlyList',
        component: () => import('@/views/Elderly/List.vue'),
        meta: {
          title: '老人列表',
          requiresAuth: true
        }
      },
      {
        path: '/elderly/detail/:id',
        name: 'ElderlyDetail',
        component: () => import('@/views/Elderly/Detail.vue'),
        meta: {
          title: '老人详情',
          requiresAuth: true
        }
      },
      {
        path: '/elderly/health/:id',
        name: 'ElderlyHealth',
        component: () => import('@/views/Elderly/Health.vue'),
        meta: {
          title: '健康档案',
          requiresAuth: true
        }
      },
      {
        path: '/care',
        name: 'Care',
        redirect: '/care/tasks',
        meta: {
          title: '护理管理',
          requiresAuth: true
        }
      },
      {
        path: '/care/tasks',
        name: 'CareTasks',
        component: () => import('@/views/Care/Tasks.vue'),
        meta: {
          title: '护理任务',
          requiresAuth: true
        }
      },
      {
        path: '/care/records',
        name: 'CareRecords',
        component: () => import('@/views/Care/Records.vue'),
        meta: {
          title: '护理记录',
          requiresAuth: true
        }
      },
      {
        path: '/care/assessment',
        name: 'CareAssessment',
        component: () => import('@/views/Care/Assessment.vue'),
        meta: {
          title: '护理评估',
          requiresAuth: true
        }
      },
      {
        path: '/care/scheduler',
        name: 'CareScheduler',
        component: () => import('@/components/Care/TaskScheduler.vue'),
        meta: {
          title: '任务调度',
          requiresAuth: true
        }
      },
      {
        path: '/finance',
        name: 'Finance',
        redirect: '/finance/bills',
        meta: {
          title: '财务管理',
          requiresAuth: true
        }
      },
      {
        path: '/finance/bills',
        name: 'FinanceBills',
        component: () => import('@/views/Finance/Bills.vue'),
        meta: {
          title: '费用账单',
          requiresAuth: true
        }
      },
      {
        path: '/finance/payments',
        name: 'FinancePayments',
        component: () => import('@/views/Finance/Payments.vue'),
        meta: {
          title: '支付记录',
          requiresAuth: true
        }
      },
      {
        path: '/health',
        name: 'Health',
        redirect: '/health/medication',
        meta: {
          title: '健康管理',
          requiresAuth: true
        }
      },
      {
        path: '/health/medication',
        name: 'HealthMedication',
        component: () => import('@/views/Health/Medication.vue'),
        meta: {
          title: '药品管理',
          requiresAuth: true
        }
      },
      {
        path: '/health/alerts',
        name: 'HealthAlerts',
        component: () => import('@/components/Health/HealthAlertSystem.vue'),
        meta: {
          title: '健康预警',
          requiresAuth: true
        }
      },
      {
        path: '/visits',
        name: 'Visits',
        redirect: '/visits/appointments',
        meta: {
          title: '探视管理',
          requiresAuth: true
        }
      },
      {
        path: '/visits/appointments',
        name: 'VisitsAppointments',
        component: () => import('@/views/Visits/Appointments.vue'),
        meta: {
          title: '探视预约',
          requiresAuth: true
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: '404'
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 养老院管理系统`
  }

  // 检查是否需要登录
  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
