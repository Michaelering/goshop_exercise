import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { getToken } from '@/utils/token'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'manager',
        name: 'AdminManager',
        component: () => import('@/views/admin/manager/ManagerList.vue'),
        meta: { title: '管理员列表' },
      },
      {
        path: 'manager/add',
        name: 'AdminManagerAdd',
        component: () => import('@/views/admin/manager/ManagerForm.vue'),
        meta: { title: '增加管理员' },
      },
      {
        path: 'manager/edit/:id',
        name: 'AdminManagerEdit',
        component: () => import('@/views/admin/manager/ManagerForm.vue'),
        meta: { title: '修改管理员' },
      },
      // Role routes
      {
        path: 'role',
        name: 'AdminRole',
        component: () => import('@/views/admin/role/RoleList.vue'),
        meta: { title: '角色列表' },
      },
      {
        path: 'role/add',
        name: 'AdminRoleAdd',
        component: () => import('@/views/admin/role/RoleForm.vue'),
        meta: { title: '增加角色' },
      },
      {
        path: 'role/edit/:id',
        name: 'AdminRoleEdit',
        component: () => import('@/views/admin/role/RoleForm.vue'),
        meta: { title: '修改角色' },
      },
      // Goods routes (placeholder)
      {
        path: 'goods',
        name: 'AdminGoods',
        component: () => import('@/views/admin/goods/GoodsList.vue'),
        meta: { title: '商品列表' },
      },
    ],
  },
  {
    path: '/merchant',
    name: 'MerchantLayout',
    component: () => import('@/layouts/MerchantLayout.vue'),
    meta: { requiresAuth: true, role: 'merchant' },
    redirect: '/merchant/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'MerchantDashboard',
        component: () => import('@/views/merchant/dashboard/DashboardView.vue'),
        meta: { title: '店铺仪表盘' },
      },
      {
        path: 'goods',
        name: 'MerchantGoods',
        component: () => import('@/views/merchant/goods/GoodsList.vue'),
        meta: { title: '商品管理' },
      },
      {
        path: 'order',
        name: 'MerchantOrder',
        component: () => import('@/views/merchant/order/OrderList.vue'),
        meta: { title: '订单管理' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login',
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const hasToken = !!getToken()

  if (to.meta.requiresAuth) {
    if (!hasToken) {
      next({ path: '/login' })
    } else {
      next()
    }
  } else {
    if (hasToken && to.path === '/login') {
      // 已登录跳转到首页，由各 layout 自行重定向
      next({ path: '/' })
    } else {
      next()
    }
  }
})

export default router
