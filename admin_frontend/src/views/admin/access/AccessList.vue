<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>权限列表</span>
        <el-button type="primary" @click="$router.push('/admin/access/add')">
          增加权限
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="module_name" label="模块名称" />
      <el-table-column prop="action_name" label="操作名称" />
      <el-table-column label="类型">
        <template #default="{ row }">
          <el-tag v-if="row.type === 1">模块</el-tag>
          <el-tag v-else-if="row.type === 2" type="success">菜单</el-tag>
          <el-tag v-else-if="row.type === 3" type="warning">操作</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="url" label="URL" />
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.status === 1"
            @change="(val: boolean) => handleStatusChange(row.id, val)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/access/edit/' + row.id)">
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

async function loadData() {
  loading.value = true
  try {
    const res: any = await api.getAccessList()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleStatusChange(id: number, val: boolean) {
  try {
    await api.changeStatus('access', 'status', id)
    ElMessage.success(val ? '已启用' : '已禁用')
    loadData()
  } catch {
    loadData()
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除吗？', '确认删除')
    await api.deleteAccess(id)
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
