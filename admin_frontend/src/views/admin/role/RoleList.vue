<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>角色列表</span>
        <el-button type="primary" @click="$router.push('/admin/role/add')">增加角色</el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="title" label="角色名称">
        <template #default="{ row }">
          {{ row.title }}
          <el-tag v-if="row.isBuiltin === 1" size="small" type="warning" style="margin-left:8px">内置</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" />
      <el-table-column label="创建时间">
        <template #default="{ row }">
          {{ formatTime(row.addTime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/role/edit/' + row.id)">修改</el-button>
          <el-button size="small" type="warning" @click="$router.push('/admin/role/auth/' + row.id)">授权</el-button>
          <el-button
            size="small"
            type="danger"
            :disabled="row.isBuiltin === 1"
            @click="handleDelete(row.id)"
          >
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
import request from '@/api/request'

const list = ref<any[]>([])
const loading = ref(false)

function formatTime(ts: number) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  return d.toLocaleDateString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await request.get('/admin/role')
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除此角色吗？', '确认删除')
    await request.delete('/admin/role/' + id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
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
