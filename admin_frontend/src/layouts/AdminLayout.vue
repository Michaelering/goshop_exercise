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
        <!-- 动态菜单：从后端 /admin/menu 获取 -->
        <template v-for="item in menuList" :key="item.id">
          <!-- 多子菜单 → el-sub-menu -->
          <el-sub-menu
            v-if="item.children && item.children.length > 1"
            :index="item.url_prefix + '-sub'"
          >
            <template #title>
              <el-icon><component :is="getIcon(item.url_prefix)" /></el-icon>
              <span>{{ item.module_name }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.id"
              :index="getChildRoute(item.url_prefix, child.action_name)"
            >
              {{ child.action_name }}
            </el-menu-item>
          </el-sub-menu>

          <!-- 0-1 个子菜单 → 直接渲染为菜单项（点击进入列表页） -->
          <el-menu-item
            v-else
            :index="getChildRoute(item.url_prefix, item.children?.[0]?.action_name || '')"
          >
            <el-icon><component :is="getIcon(item.url_prefix)" /></el-icon>
            <template #title>{{ item.module_name }}</template>
          </el-menu-item>
        </template>
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
              <el-tag v-if="authStore.user?.roleTitle" size="small" type="warning" style="margin-left:6px">
                {{ authStore.user.roleTitle }}
              </el-tag>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getToken } from '@/utils/token'
import request from '@/api/request'
import { adminLogout } from '@/api/auth'
import {
  Fold, Expand, Odometer, User, Setting, Goods, SwitchButton,
  Lock, Collection, Tickets, List, Guide, PictureFilled, Tools, Shop, Menu
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)
const menuList = ref<any[]>([])

const activeMenu = computed(() => route.path)

// 图标映射：url_prefix → Element Plus icon component
const iconMap: Record<string, any> = {
  dashboard: Odometer,
  manager: User,
  role: Setting,
  access: Lock,
  goods: Goods,
  goodsCate: Collection,
  goodsType: Tickets,
  goodsTypeAttr: List,
  nav: Guide,
  focus: PictureFilled,
  setting: Tools,
  merchant: Shop,
}

function getIcon(prefix: string) {
  return iconMap[prefix] || Menu
}

// 根据子菜单决定路由目标
function getChildRoute(parentPrefix: string, actionName: string): string {
  // 0 或 1 个子菜单 → 直接进入列表页
  if (!actionName) {
    return '/admin/' + parentPrefix
  }
  // 如"增加角色" → /admin/role/add
  if (actionName.includes('增加')) {
    return '/admin/' + parentPrefix + '/add'
  }
  // 默认列表页
  return '/admin/' + parentPrefix
}

async function loadMenu() {
  try {
    const res: any = await request.get('/admin/menu')
    if (res && res.data) {
      menuList.value = res.data
    }
  } catch (e) {
    console.error('加载菜单失败:', e)
  }
}

// 当 token 发生变化时重新加载菜单（处理登录后首次挂载）
watch(() => authStore.isLoggedIn, (val) => {
  if (val) loadMenu()
})

function toggleCollapse() {
  isCollapse.value = !isCollapse.value
}

function handleLogout() {
  adminLogout().finally(() => {})
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  // 仅在已登录状态下加载菜单
  if (getToken()) {
    loadMenu()
  }
})
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
