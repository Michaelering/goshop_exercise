<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>商品列表</span>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="title" label="商品名称" min-width="200" show-overflow-tooltip />
      <el-table-column label="价格" width="120">
        <template #default="{ row }">
          <span style="color:#f56c6c;font-weight:700;">¥{{ row.price }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="clickCount" label="点击量" width="80" align="center" />
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="排序" width="80" align="center">
        <template #default="{ row }">
          <el-tag>{{ row.sort }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="goodsNumber" label="库存" width="80" align="center" />
      <el-table-column label="操作" width="120" align="center">
        <template #default="{ row }">
          <el-button size="small">修改</el-button>
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

async function loadData() {
  loading.value = true
  try {
    const res: any = await request.get('/admin/goods', {
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
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mt-3 { margin-top: 16px; }
</style>
