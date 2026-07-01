<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span><i class="el-icon-user"></i> 管理员列表</span>
        <el-button type="primary" @click="$router.push('/admin/manager/add')">
          增加管理员
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="username" label="管理员名称" />
      <el-table-column prop="mobile" label="电话" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="角色">
        <template #default="{ row }">
          <el-tag>{{ row.role?.title || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间">
        <template #default="{ row }">
          {{ formatTime(row.addTime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/manager/edit/' + row.id)">
            修改
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">
            删除
          </el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'

const list = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

function formatTime(ts: number) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  return d.toLocaleDateString('zh-CN') + ' ' + d.toLocaleTimeString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await request.get('/admin/manager', {
      params: { page: page.value, pageSize: pageSize.value },
    })
    list.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除此管理员吗？', '确认删除')
    await request.delete('/admin/manager/' + id)
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
.mt-3 {
  margin-top: 16px;
}
</style>
