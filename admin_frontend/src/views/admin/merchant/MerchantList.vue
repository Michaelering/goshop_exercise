<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>商户列表</span>
        <el-button type="primary" @click="$router.push('/admin/merchant/add')">
          增加商户
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="shop_name" label="店铺名称" />
      <el-table-column prop="mobile" label="手机号" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.status === 1"
            @change="(val: boolean) => handleStatusChange(row.id, val)"
          />
        </template>
      </el-table-column>
      <el-table-column label="创建时间">
        <template #default="{ row }">
          {{ formatTime(row.add_time) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/merchant/edit/' + row.id)">
            修改
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as api from '@/api/admin'

const list = ref<any[]>([])
const loading = ref(false)

function formatTime(ts: number) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await api.getMerchantList()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleStatusChange(id: number, val: boolean) {
  try {
    await api.changeStatus('merchant', 'status', id)
    ElMessage.success(val ? '已启用' : '已禁用')
    loadData()
  } catch {
    loadData()
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除吗？', '确认删除')
    await api.deleteMerchant(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // cancelled
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
</style>
