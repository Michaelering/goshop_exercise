<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>订单列表</span>
        <el-input
          v-model="keyword"
          placeholder="搜索订单号/收货人/电话"
          style="width:260px"
          clearable
          @keyup.enter="search"
          @clear="search"
        >
          <template #append>
            <el-button @click="search">搜索</el-button>
          </template>
        </el-input>
      </div>
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
      <el-table-column label="电话" width="130">
        <template #default="{ row }">{{ row.phone || '-' }}</template>
      </el-table-column>
      <el-table-column label="地址" min-width="160" show-overflow-tooltip>
        <template #default="{ row }">{{ row.address || '-' }}</template>
      </el-table-column>
      <el-table-column label="总金额" width="110">
        <template #default="{ row }">
          <span style="color:#f56c6c;font-weight:700;">¥{{ row.all_price?.toFixed(2) || '0.00' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="110">
        <template #default="{ row }">
          <el-select
            :model-value="row.order_status"
            size="small"
            style="width:100px"
            @change="(val: number) => handleStatusChange(row.id, 'order_status', val)"
          >
            <el-option v-for="(label, val) in orderStatusMap" :key="val" :label="label" :value="Number(val)" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="支付" width="90">
        <template #default="{ row }">
          <el-switch
            :model-value="row.pay_status === 1"
            size="small"
            @change="(val: boolean) => handleStatusChange(row.id, 'pay_status', val ? 1 : 0)"
          />
        </template>
      </el-table-column>
      <el-table-column label="下单时间" width="170">
        <template #default="{ row }">{{ formatTime(row.add_time) }}</template>
      </el-table-column>
      <el-table-column label="商品" min-width="180">
        <template #default="{ row }">
          <div v-for="item in row.order_item" :key="item.id" class="order-item-row">
            <el-tag size="small" type="warning">{{ item.product_title }}</el-tag>
            <span class="item-info">¥{{ item.product_price }} × {{ item.product_num }}</span>
          </div>
          <span v-if="!row.order_item?.length" style="color:#999">-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openDetail(row)">详情</el-button>
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

    <!-- 订单详情弹窗 -->
    <el-dialog v-model="dialogVisible" title="订单详情" width="700px">
      <template v-if="currentOrder">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.order_id }}</el-descriptions-item>
          <el-descriptions-item label="总金额">
            <span style="color:#f56c6c;font-weight:700;">¥{{ currentOrder.all_price?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="收货人">{{ currentOrder.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ currentOrder.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="收货地址" :span="2">{{ currentOrder.address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-select
              :model-value="currentOrder.order_status"
              size="small"
              style="width:110px"
              @change="(val: number) => handleStatusChange(currentOrder.id, 'order_status', val)"
            >
              <el-option v-for="(label, val) in orderStatusMap" :key="val" :label="label" :value="Number(val)" />
            </el-select>
          </el-descriptions-item>
          <el-descriptions-item label="支付状态">
            <el-switch
              :model-value="currentOrder.pay_status === 1"
              @change="(val: boolean) => handleStatusChange(currentOrder.id, 'pay_status', val ? 1 : 0)"
            />
          </el-descriptions-item>
          <el-descriptions-item label="下单时间">{{ formatTime(currentOrder.add_time) }}</el-descriptions-item>
          <el-descriptions-item label="付款时间">{{ formatTime(currentOrder.pay_time) }}</el-descriptions-item>
          <el-descriptions-item label="物流公司">{{ currentOrder.logistics_company || '-' }}</el-descriptions-item>
          <el-descriptions-item label="运单号">{{ currentOrder.waybill_no || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>商品明细</el-divider>
        <el-table :data="currentOrder.order_item" size="small">
          <el-table-column prop="product_title" label="商品名称" />
          <el-table-column prop="product_price" label="单价" width="100">
            <template #default="{ row }">¥{{ row.product_price }}</template>
          </el-table-column>
          <el-table-column prop="product_num" label="数量" width="80" />
          <el-table-column prop="goods_version" label="规格" width="100">
            <template #default="{ row }">{{ row.goods_version || '-' }}</template>
          </el-table-column>
          <el-table-column prop="goods_color" label="颜色" width="100">
            <template #default="{ row }">{{ row.goods_color || '-' }}</template>
          </el-table-column>
        </el-table>
      </template>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

const list = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref('')

const dialogVisible = ref(false)
const currentOrder = ref<any>(null)

const orderStatusMap: Record<number, string> = {
  0: '已下单', 1: '已付款', 2: '已配货', 3: '已发货', 4: '交易成功', 5: '退货', 6: '已取消',
}

function formatTime(ts: number) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await request.get('/admin/order', {
      params: { page: page.value, pageSize: pageSize.value, keyword: keyword.value },
    })
    list.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadData()
}

async function handleStatusChange(id: number, field: string, value: number) {
  try {
    const data = new URLSearchParams()
    data.append(field, String(value))
    await request.put('/admin/order/' + id, data)
    ElMessage.success('状态已更新')
    // 更新本地列表中的值
    const item = list.value.find((o: any) => o.id === id)
    if (item) (item as any)[field] = value
    if (currentOrder.value?.id === id) (currentOrder.value as any)[field] = value
  } catch {
    loadData()
  }
}

async function openDetail(row: any) {
  try {
    const res: any = await request.get('/admin/order/' + row.id)
    currentOrder.value = res.data
    dialogVisible.value = true
  } catch {
    ElMessage.error('获取订单详情失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.mt-3 { margin-top: 16px; }
.order-item-row { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.order-item-row:last-child { margin-bottom: 0; }
.item-info { font-size: 13px; color: #666; }
</style>
