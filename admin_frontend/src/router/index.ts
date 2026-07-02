import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { getToken, getUserType } from '@/utils/token'

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
      {
        path: 'role/auth/:id',
        name: 'AdminRoleAuth',
        component: () => import('@/views/admin/role/RoleAuth.vue'),
        meta: { title: '角色授权' },
      },
      {
        path: 'access',
        name: 'AdminAccess',
        component: () => import('@/views/admin/access/AccessList.vue'),
        meta: { title: '权限列表' },
      },
      {
        path: 'access/add',
        name: 'AdminAccessAdd',
        component: () => import('@/views/admin/access/AccessForm.vue'),
        meta: { title: '增加权限' },
      },
      {
        path: 'access/edit/:id',
        name: 'AdminAccessEdit',
        component: () => import('@/views/admin/access/AccessForm.vue'),
        meta: { title: '修改权限' },
      },
      {
        path: 'goods',
        name: 'AdminGoods',
        component: () => import('@/views/admin/goods/GoodsList.vue'),
        meta: { title: '商品列表' },
      },
      {
        path: 'goods/add',
        name: 'AdminGoodsAdd',
        component: () => import('@/views/admin/goods/GoodsForm.vue'),
        meta: { title: '增加商品' },
      },
      {
        path: 'goods/edit/:id',
        name: 'AdminGoodsEdit',
        component: () => import('@/views/admin/goods/GoodsForm.vue'),
        meta: { title: '修改商品' },
      },
      {
        path: 'goodsCate',
        name: 'AdminGoodsCate',
        component: () => import('@/views/admin/goodsCate/GoodsCateList.vue'),
        meta: { title: '商品分类' },
      },
      {
        path: 'goodsCate/add',
        name: 'AdminGoodsCateAdd',
        component: () => import('@/views/admin/goodsCate/GoodsCateForm.vue'),
        meta: { title: '增加分类' },
      },
      {
        path: 'goodsCate/edit/:id',
        name: 'AdminGoodsCateEdit',
        component: () => import('@/views/admin/goodsCate/GoodsCateForm.vue'),
        meta: { title: '修改分类' },
      },
      {
        path: 'goodsType',
        name: 'AdminGoodsType',
        component: () => import('@/views/admin/goodsType/GoodsTypeList.vue'),
        meta: { title: '商品类型' },
      },
      {
        path: 'goodsType/add',
        name: 'AdminGoodsTypeAdd',
        component: () => import('@/views/admin/goodsType/GoodsTypeForm.vue'),
        meta: { title: '增加类型' },
      },
      {
        path: 'goodsType/edit/:id',
        name: 'AdminGoodsTypeEdit',
        component: () => import('@/views/admin/goodsType/GoodsTypeForm.vue'),
        meta: { title: '修改类型' },
      },
      {
        path: 'goodsTypeAttr',
        name: 'AdminGoodsTypeAttr',
        component: () => import('@/views/admin/goodsTypeAttr/GoodsTypeAttrList.vue'),
        meta: { title: '类型属性' },
      },
      {
        path: 'goodsTypeAttr/add',
        name: 'AdminGoodsTypeAttrAdd',
        component: () => import('@/views/admin/goodsTypeAttr/GoodsTypeAttrForm.vue'),
        meta: { title: '增加属性' },
      },
      {
        path: 'goodsTypeAttr/edit/:id',
        name: 'AdminGoodsTypeAttrEdit',
        component: () => import('@/views/admin/goodsTypeAttr/GoodsTypeAttrForm.vue'),
        meta: { title: '修改属性' },
      },
      {
        path: 'nav',
        name: 'AdminNav',
        component: () => import('@/views/admin/nav/NavList.vue'),
        meta: { title: '导航列表' },
      },
      {
        path: 'nav/add',
        name: 'AdminNavAdd',
        component: () => import('@/views/admin/nav/NavForm.vue'),
        meta: { title: '增加导航' },
      },
      {
        path: 'nav/edit/:id',
        name: 'AdminNavEdit',
        component: () => import('@/views/admin/nav/NavForm.vue'),
        meta: { title: '修改导航' },
      },
      {
        path: 'focus',
        name: 'AdminFocus',
        component: () => import('@/views/admin/focus/FocusList.vue'),
        meta: { title: '轮播图列表' },
      },
      {
        path: 'focus/add',
        name: 'AdminFocusAdd',
        component: () => import('@/views/admin/focus/FocusForm.vue'),
        meta: { title: '增加轮播图' },
      },
      {
        path: 'focus/edit/:id',
        name: 'AdminFocusEdit',
        component: () => import('@/views/admin/focus/FocusForm.vue'),
        meta: { title: '修改轮播图' },
      },
      {
        path: 'setting',
        name: 'AdminSetting',
        component: () => import('@/views/admin/setting/SettingForm.vue'),
        meta: { title: '系统设置' },
      },
      {
        path: 'merchant',
        name: 'AdminMerchant',
        component: () => import('@/views/admin/merchant/MerchantList.vue'),
        meta: { title: '商户列表' },
      },
      {
        path: 'merchant/add',
        name: 'AdminMerchantAdd',
        component: () => import('@/views/admin/merchant/MerchantForm.vue'),
        meta: { title: '增加商户' },
      },
      {
        path: 'merchant/edit/:id',
        name: 'AdminMerchantEdit',
        component: () => import('@/views/admin/merchant/MerchantForm.vue'),
        meta: { title: '修改商户' },
      },
      // 订单管理（管理员）
      {
        path: 'order',
        name: 'AdminOrder',
        component: () => import('@/views/admin/order/OrderList.vue'),
        meta: { title: '订单列表' },
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
        path: 'goods/add',
        name: 'MerchantGoodsAdd',
        component: () => import('@/views/merchant/goods/GoodsForm.vue'),
        meta: { title: '增加商品' },
      },
      {
        path: 'goods/edit/:id',
        name: 'MerchantGoodsEdit',
        component: () => import('@/views/merchant/goods/GoodsForm.vue'),
        meta: { title: '修改商品' },
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

router.beforeEach((to, _from, next) => {
  const hasToken = !!getToken()
  const userType = getUserType()

  if (to.meta.requiresAuth) {
    // 无 token → 登录页
    if (!hasToken) {
      next({ path: '/login' })
      return
    }
    // 有 token 但角色不匹配 → 重定向到对应首页
    if (to.meta.role && userType && to.meta.role !== userType) {
      const targetPath = userType === 'admin' ? '/admin/dashboard' : '/merchant/dashboard'
      next({ path: targetPath })
      return
    }
    next()
  } else {
    // 已登录用户访问登录页 → 跳转对应首页
    if (hasToken && to.path === '/login') {
      const targetPath = userType === 'admin' ? '/admin/dashboard' : '/merchant/dashboard'
      next({ path: targetPath })
    } else {
      next()
    }
  }
})

export default router
