<template>
  <el-card>
    <template #header>
      <span>角色授权 - {{ roleTitle }}</span>
    </template>

    <el-tree
      ref="treeRef"
      :data="accessList"
      show-checkbox
      node-key="id"
      default-expand-all
      :default-checked-keys="checkedKeys"
      :props="{ children: 'children', label: 'module_name' }"
    />

    <el-button type="primary" style="margin-top:20px;" :loading="submitting" @click="handleSubmit">
      保存授权
    </el-button>
    <el-button @click="$router.back()">返回</el-button>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as api from '@/api/admin'

const route = useRoute()
const treeRef = ref()
const submitting = ref(false)
const accessList = ref<any[]>([])
const checkedKeys = ref<number[]>([])
const roleTitle = ref('')

function collectCheckedKeys(nodes: any[], keys: number[]) {
  for (const node of nodes) {
    if (node.checked) {
      keys.push(node.id)
    }
    if (node.children && node.children.length > 0) {
      collectCheckedKeys(node.children, keys)
    }
  }
}

async function loadData() {
  const roleId = Number(route.params.id)
  const res: any = await api.getRoleAuth(roleId)
  if (res.data) {
    roleTitle.value = res.data.roleTitle || ''
    accessList.value = res.data.accessList || []
    const keys: number[] = []
    collectCheckedKeys(accessList.value, keys)
    checkedKeys.value = keys
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    const checkedNodes = treeRef.value?.getCheckedKeys() || []
    const halfCheckedNodes = treeRef.value?.getHalfCheckedKeys() || []
    const allCheckedIds = [...checkedNodes, ...halfCheckedNodes]

    const data = new URLSearchParams()
    allCheckedIds.forEach((id: number) => {
      data.append('access_node[]', String(id))
    })

    await api.saveRoleAuth(Number(route.params.id), data)
    ElMessage.success('授权成功')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
