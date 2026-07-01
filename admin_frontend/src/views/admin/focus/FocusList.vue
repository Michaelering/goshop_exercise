<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>轮播图列表</span>
        <el-button type="primary" @click="$router.push('/admin/focus/add')">
          增加轮播图
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="title" label="标题" />
      <el-table-column label="类型">
        <template #default="{ row }">
          <el-tag :type="row.focus_type === 1 ? 'primary' : 'success'">
            {{ row.focus_type === 1 ? 'PC' : '移动' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="图片">
        <template #default="{ row }">
          <el-image
            v-if="row.focus_img"
            :src="row.focus_img"
            style="width:80px; height:40px;"
            fit="cover"
            :preview-src-list="[row.focus_img]"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="link" label="链接" />
      <el-table-column prop="sort" label="排序" />
      <el-table-column label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="$router.push('/admin/focus/edit/' + row.id)">
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
    const res: any = await api.getFocusList()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除吗？', '确认删除')
    await api.deleteFocus(id)
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
