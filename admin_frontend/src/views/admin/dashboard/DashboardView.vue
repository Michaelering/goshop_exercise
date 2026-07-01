<template>
  <div class="dashboard">
    <el-row :gutter="20" class="mb-4">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#409eff;">
              <el-icon :size="32"><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.managerCount }}</div>
              <div class="stat-label">管理员总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#67c23a;">
              <el-icon :size="32"><Goods /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.goodsCount }}</div>
              <div class="stat-label">商品总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#e6a23c;">
              <el-icon :size="32"><List /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.orderCount }}</div>
              <div class="stat-label">订单总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <span>欢迎回来</span>
      </template>
      <p>当前登录：<strong>{{ authStore.user?.username }}</strong></p>
      <p>欢迎使用 GoShop 后台管理系统</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import request from '@/api/request'
import { UserFilled, Goods, List } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const stats = ref({
  managerCount: 0,
  goodsCount: 0,
  orderCount: 0,
})

onMounted(async () => {
  try {
    const res: any = await request.get('/admin/dashboard')
    stats.value = res.data
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.mb-4 {
  margin-bottom: 20px;
}
.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
}
.stat-icon {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
}
.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}
</style>
