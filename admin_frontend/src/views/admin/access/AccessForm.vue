<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改权限' : '增加权限' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="模块名称" prop="module_name">
        <el-input v-model="form.module_name" />
      </el-form-item>
      <el-form-item label="操作名称" prop="action_name">
        <el-input v-model="form.action_name" />
      </el-form-item>
      <el-form-item label="类型" prop="type">
        <el-select v-model="form.type" placeholder="选择类型" style="width:100%">
          <el-option label="模块" :value="1" />
          <el-option label="菜单" :value="2" />
          <el-option label="操作" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item label="URL" prop="url">
        <el-input v-model="form.url" />
      </el-form-item>
      <el-form-item label="所属模块" prop="module_id">
        <el-select v-model="form.module_id" placeholder="选择所属模块" style="width:100%">
          <el-option label="顶级模块" :value="0" />
          <el-option
            v-for="m in topModules"
            :key="m.id"
            :label="m.module_name"
            :value="m.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input-number v-model="form.sort" :min="0" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input v-model="form.description" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          提交
        </el-button>
        <el-button @click="$router.back()">返回</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import * as api from '@/api/admin'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const topModules = ref<any[]>([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  module_name: '',
  action_name: '',
  type: 1 as number,
  url: '',
  module_id: 0 as number,
  sort: 10,
  status: 1,
  description: '',
})

const rules: FormRules = {
  module_name: [{ required: true, message: '请输入模块名称', trigger: 'blur' }],
}

async function loadTopModules() {
  const res: any = await api.getAccessTopModules()
  topModules.value = res.data || []
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getAccess(Number(id))
  const item = res.data?.access || res.data
  if (item) {
    form.module_name = item.module_name
    form.action_name = item.action_name || ''
    form.type = item.type
    form.url = item.url || ''
    form.module_id = item.module_id ?? 0
    form.sort = item.sort ?? 10
    form.status = item.status ?? 1
    form.description = item.description || ''
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = new URLSearchParams()
    data.append('module_name', form.module_name)
    data.append('action_name', form.action_name)
    data.append('type', String(form.type))
    data.append('url', form.url)
    data.append('module_id', String(form.module_id))
    data.append('sort', String(form.sort))
    data.append('status', String(form.status))
    data.append('description', form.description)

    if (isEdit.value) {
      await api.updateAccess(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createAccess(data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/access')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadTopModules()
  loadData()
})
</script>
