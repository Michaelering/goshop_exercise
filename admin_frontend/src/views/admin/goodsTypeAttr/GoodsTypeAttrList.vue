<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>商品类型属性列表</span>
        <el-button type="primary" @click="$router.push('/admin/goodsTypeAttr/add')">
          增加属性
        </el-button>
      </div>
    </template>

    <el-table :data="list" stripe v-loading="loading" style="width:100%">
      <el-table-column prop="title" label="属性标题" />
      <el-table-column prop="cate_id" label="所属类型ID" width="100" align="center" />
      <el-table-column label="属性类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.attr_type === 1">文本</el-tag>
          <el-tag v-else-if="row.attr_type === 2" type="success">文本域</el-tag>
          <el-tag v-else-if="row.attr_type === 3" type="warning">下拉</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="属性值">
        <template #default="{ row }">
          {{ row.attr_value ? row.attr_value.substring(0, 30) + (row.attr_value.length > 30 ? '...' : '') : '-' }}
        </template>
      </el-table-column>
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
          <el-button size="small" @click="$router.push('/admin/goodsTypeAttr/edit/' + row.id)">
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
    const res: any = await api.getGoodsTypeAttrList()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleStatusChange(id: number, val: boolean) {
  try {
    await api.changeStatus('goods_type_attr', 'status', id)
    ElMessage.success(val ? '已启用' : '已禁用')
    loadData()
  } catch {
    loadData()
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除吗？', '确认删除')
    await api.deleteGoodsTypeAttr(id)
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
