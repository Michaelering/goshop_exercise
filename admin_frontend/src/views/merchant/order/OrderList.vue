<template>
  <el-card>
    <template #header>
      <span>包含我店铺商品的订单</span>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column label="订单号" width="160">
        <template #default="{ row }">
          <el-tag>{{ row.orderId }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="收货人" width="100" />
      <el-table-column prop="phone" label="联系电话" width="130" />
      <el-table-column label="总金额" width="100">
        <template #default="{ row }">
          <span style="color:#f56c6c;font-weight:700;">¥{{ row.allPrice }}</span>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.orderStatus)">
            {{ statusText(row.orderStatus) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="支付状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.payStatus === 1 ? 'success' : 'danger'">
            {{ row.payStatus === 1 ? '已支付' : '未支付' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="下单时间" width="160">
        <template #default="{ row }">{{ formatTime(row.addTime) }}</template>
      </el-table-column>
      <el-table-column label="商品" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="item in row.orderItem" :key="item.id" size="small" class="mr-1">
            {{ item.productTitle }} x{{ item.productNum }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="total > 0"
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      class="mt-3"
      @current-change="loadData"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const list = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

function formatTime(ts: number) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString('zh-CN')
}

function statusType(s: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'info', 2: 'primary', 3: 'primary', 4: 'success', 5: 'danger', 6: 'info' }
  return map[s] || 'info'
}

function statusText(s: number) {
  const map: Record<number, string> = { 0: '已下单', 1: '已付款', 2: '已配货', 3: '已发货', 4: '交易成功', 5: '退货', 6: '已取消' }
  return map[s] || '未知'
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await request.get('/merchant/order', {
      params: { page: page.value, pageSize: pageSize.value },
    })
    list.value = res.data?.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.mt-3 { margin-top: 16px; }
.mr-1 { margin-right: 4px; margin-bottom: 4px; }
</style>
