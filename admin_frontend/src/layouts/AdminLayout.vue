<template>
  <el-container class="admin-layout">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '220px'" class="admin-sidebar">
      <div class="sidebar-logo">
        <i class="el-icon-s-shop"></i>
        <span v-show="!isCollapse">GoShop 管理</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <!-- 管理员管理 -->
        <el-sub-menu index="manager-sub">
          <template #title>
            <el-icon><User /></el-icon>
            <span>管理员管理</span>
          </template>
          <el-menu-item index="/admin/manager">管理员列表</el-menu-item>
          <el-menu-item index="/admin/manager/add">增加管理员</el-menu-item>
        </el-sub-menu>

        <!-- 角色管理 -->
        <el-sub-menu index="role-sub">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>角色管理</span>
          </template>
          <el-menu-item index="/admin/role">角色列表</el-menu-item>
          <el-menu-item index="/admin/role/add">增加角色</el-menu-item>
        </el-sub-menu>

        <!-- 权限管理 -->
        <el-sub-menu index="access-sub">
          <template #title>
            <el-icon><Lock /></el-icon>
            <span>权限管理</span>
          </template>
          <el-menu-item index="/admin/access">权限列表</el-menu-item>
          <el-menu-item index="/admin/access/add">增加权限</el-menu-item>
        </el-sub-menu>

        <!-- 商品管理 -->
        <el-sub-menu index="goods-sub">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="/admin/goods">商品列表</el-menu-item>
          <el-menu-item index="/admin/goods/add">增加商品</el-menu-item>
        </el-sub-menu>

        <!-- 商品分类 -->
        <el-sub-menu index="goodsCate-sub">
          <template #title>
            <el-icon><Collection /></el-icon>
            <span>商品分类</span>
          </template>
          <el-menu-item index="/admin/goodsCate">分类列表</el-menu-item>
          <el-menu-item index="/admin/goodsCate/add">增加分类</el-menu-item>
        </el-sub-menu>

        <!-- 商品类型 -->
        <el-sub-menu index="goodsType-sub">
          <template #title>
            <el-icon><Tickets /></el-icon>
            <span>商品类型</span>
          </template>
          <el-menu-item index="/admin/goodsType">类型列表</el-menu-item>
          <el-menu-item index="/admin/goodsType/add">增加类型</el-menu-item>
        </el-sub-menu>

        <!-- 商品类型属性 -->
        <el-sub-menu index="goodsTypeAttr-sub">
          <template #title>
            <el-icon><List /></el-icon>
            <span>类型属性</span>
          </template>
          <el-menu-item index="/admin/goodsTypeAttr">属性列表</el-menu-item>
          <el-menu-item index="/admin/goodsTypeAttr/add">增加属性</el-menu-item>
        </el-sub-menu>

        <!-- 导航管理 -->
        <el-sub-menu index="nav-sub">
          <template #title>
            <el-icon><Guide /></el-icon>
            <span>导航管理</span>
          </template>
          <el-menu-item index="/admin/nav">导航列表</el-menu-item>
          <el-menu-item index="/admin/nav/add">增加导航</el-menu-item>
        </el-sub-menu>

        <!-- 轮播图管理 -->
        <el-sub-menu index="focus-sub">
          <template #title>
            <el-icon><PictureFilled /></el-icon>
            <span>轮播图管理</span>
          </template>
          <el-menu-item index="/admin/focus">轮播图列表</el-menu-item>
          <el-menu-item index="/admin/focus/add">增加轮播图</el-menu-item>
        </el-sub-menu>

        <!-- 系统设置 -->
        <el-menu-item index="/admin/setting">
          <el-icon><Tools /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>

        <!-- 商户管理 -->
        <el-sub-menu index="merchant-sub">
          <template #title>
            <el-icon><Shop /></el-icon>
            <span>商户管理</span>
          </template>
          <el-menu-item index="/admin/merchant">商户列表</el-menu-item>
          <el-menu-item index="/admin/merchant/add">增加商户</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- 顶部导航 -->
      <el-header class="admin-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="toggleCollapse">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click">
            <span class="user-info">
              <el-icon><User /></el-icon>
              {{ authStore.user?.username || '管理员' }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Fold, Expand, Odometer, User, Setting, Goods, SwitchButton, Lock, Collection, Tickets, List, Guide, PictureFilled, Tools, Shop } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)

const activeMenu = computed(() => route.path)

function toggleCollapse() {
  isCollapse.value = !isCollapse.value
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}
.admin-sidebar {
  background-color: #304156;
  overflow-y: auto;
  overflow-x: hidden;
  transition: width 0.3s;
}
.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  gap: 8px;
}
.admin-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.collapse-btn {
  font-size: 20px;
  cursor: pointer;
}
.header-right .user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
}
.admin-main {
  background: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}
</style>
