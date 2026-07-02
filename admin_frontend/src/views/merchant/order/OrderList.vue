<template>
  <el-card>
    <template #header>
      <span>包含我店铺商品的订单</span>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column label="订单号" width="180">
        <template #default="{ row }">
          <el-tag>{{ row.order_id }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="收货人" width="100">
        <template #default="{ row }">{{ row.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="联系电话" width="130">
        <template #default="{ row }">{{ row.phone || '-' }}</template>
      </el-table-column>
      <el-table-column label="收货地址" min-width="160">
        <template #default="{ row }">{{ row.address || '-' }}</template>
      </el-table-column>
      <el-table-column label="总金额" width="110">
        <template #default="{ row }">
          <span style="color:#f56c6c;font-weight:700;">¥{{ row.all_price?.toFixed(2) || '0.00' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusType(row.order_status)">
            {{ statusText(row.order_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="支付状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.pay_status === 1 ? 'success' : 'danger'">
            {{ row.pay_status === 1 ? '已支付' : '未支付' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="支付方式" width="100">
        <template #default="{ row }">
          <el-tag :type="row.pay_type === 1 ? 'success' : ''">{{ row.pay_type === 0 ? '支付宝' : row.pay_type === 1 ? '微信' : '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="下单时间" width="170">
        <template #default="{ row }">{{ formatTime(row.add_time) }}</template>
      </el-table-column>
      <el-table-column label="付款时间" width="170">
        <template #default="{ row }">{{ formatTime(row.pay_time) }}</template>
      </el-table-column>
      <el-table-column label="商品明细" min-width="200">
        <template #default="{ row }">
          <div v-for="item in row.order_item" :key="item.id" class="order-item-row">
            <el-tag size="small" type="warning">{{ item.product_title }}</el-tag>
            <span class="item-info">¥{{ item.product_price }} × {{ item.product_num }}</span>
          </div>
          <span v-if="!row.order_item || row.order_item.length === 0" style="color:#999">-</span>
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
    list.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.mt-3 { margin-top: 16px; }
.order-item-row { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.order-item-row:last-child { margin-bottom: 0; }
.item-info { font-size: 13px; color: #666; }
</style>
