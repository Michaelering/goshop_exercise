<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>导航列表</span>
        <el-button type="primary" @click="$router.push('/admin/nav/add')">
          增加导航
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="link" label="链接" />
      <el-table-column prop="position" label="位置" />
      <el-table-column label="新窗口">
        <template #default="{ row }">
          <el-tag :type="row.is_opennew === 1 ? 'success' : 'info'">
            {{ row.is_opennew === 1 ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" />
      <el-table-column label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="添加时间">
        <template #default="{ row }">
          {{ formatTime(row.add_time) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/nav/edit/' + row.id)">
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
import * as api from '@/api/admin'

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
    const res: any = await api.getNavList({ page: page.value, pageSize: pageSize.value })
    list.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除吗？', '确认删除')
    await api.deleteNav(id)
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
