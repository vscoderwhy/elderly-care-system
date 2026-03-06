import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/views/Layout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/Statistics/Dashboard.vue')
      },
      {
        path: 'statistics/report',
        name: 'StatisticsReport',
        component: () => import('@/views/Statistics/Report.vue')
      },
      {
        path: 'elderly',
        name: 'Elderly',
        component: () => import('@/views/Elderly/Index.vue')
      },
      {
        path: 'elderly/create',
        name: 'ElderlyCreate',
        component: () => import('@/views/Elderly/Create.vue')
      },
      {
        path: 'elderly/:id',
        name: 'ElderlyDetail',
        component: () => import('@/views/Elderly/Detail.vue')
      },
      {
        path: 'care',
        name: 'Care',
        component: () => import('@/views/Care/Index.vue')
      },
      {
        path: 'service',
        name: 'ServiceRequests',
        component: () => import('@/views/Service/Requests.vue')
      },
      {
        path: 'bills',
        name: 'Bills',
        component: () => import('@/views/Bills/Index.vue')
      },
      {
        path: 'staff',
        name: 'Staff',
        component: () => import('@/views/Staff/Index.vue')
      },
      {
        path: 'rooms',
        name: 'Rooms',
        component: () => import('@/views/Rooms/Index.vue')
      },
      {
        path: 'medications',
        name: 'Medications',
        component: () => import('@/views/Medications/Index.vue')
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('@/views/Inventory/Index.vue')
      },
      {
        path: 'visits',
        name: 'Visits',
        component: () => import('@/views/Visits/Index.vue')
      },
      {
        path: 'alerts',
        name: 'Alerts',
        component: () => import('@/views/Alerts/Index.vue')
      },
      {
        path: 'export',
        name: 'Export',
        component: () => import('@/views/Export/Index.vue')
      },
      {
        path: 'system/users',
        name: 'SystemUsers',
        component: () => import('@/views/System/Users.vue')
      },
      {
        path: 'system/roles',
        name: 'SystemRoles',
        component: () => import('@/views/System/Roles.vue')
      },
      {
        path: 'system/menus',
        name: 'SystemMenus',
        component: () => import('@/views/System/Menus.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
