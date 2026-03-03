import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页', showTabbar: true }
      },
      {
        path: 'elderly',
        name: 'ElderlyList',
        component: () => import('@/views/Elderly/List.vue'),
        meta: { title: '我的家人', showTabbar: true }
      },
      {
        path: 'elderly/:id',
        name: 'ElderlyDetail',
        component: () => import('@/views/Elderly/Detail.vue'),
        meta: { title: '详情' }
      },
      {
        path: 'service/call/:elderlyId',
        name: 'ServiceCall',
        component: () => import('@/views/Service/Call.vue'),
        meta: { title: '服务呼叫' }
      },
      {
        path: 'care',
        name: 'CareRecords',
        component: () => import('@/views/Care/Records.vue'),
        meta: { title: '护理记录', showTabbar: true }
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/Tasks/Index.vue'),
        meta: { title: '今日任务', showTabbar: true }
      },
      {
        path: 'attendance',
        name: 'Attendance',
        component: () => import('@/views/Attendance/Index.vue'),
        meta: { title: '考勤打卡' }
      },
      {
        path: 'bills',
        name: 'Bills',
        component: () => import('@/views/Bills/List.vue'),
        meta: { title: '账单', showTabbar: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '我的', showTabbar: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
