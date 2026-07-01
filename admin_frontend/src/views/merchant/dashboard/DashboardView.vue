<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#67c23a;">
              <el-icon :size="32"><Goods /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.goodsCount }}</div>
              <div class="stat-label">我的商品</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background:#409eff;">
              <el-icon :size="32"><List /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.orderCount }}</div>
              <div class="stat-label">关联订单</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="mt-4">
      <template #header>
        <span><i class="el-icon-info"></i> 店铺信息</span>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="店铺名称">{{ authStore.user?.shopName }}</el-descriptions-item>
        <el-descriptions-item label="商户账号">{{ authStore.user?.username }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import request from '@/api/request'
import { Goods, List } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const stats = ref({ goodsCount: 0, orderCount: 0 })

onMounted(async () => {
  try {
    const res: any = await request.get('/merchant/dashboard')
    stats.value = res.data
  } catch {}
})
</script>

<style scoped>
.stat-card { display: flex; align-items: center; gap: 20px; }
.stat-icon { width: 80px; height: 80px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: #fff; }
.stat-value { font-size: 32px; font-weight: 700; }
.stat-label { font-size: 14px; color: #909399; margin-top: 4px; }
.mt-4 { margin-top: 20px; }
</style>
