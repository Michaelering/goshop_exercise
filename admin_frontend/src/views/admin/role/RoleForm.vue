<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改角色' : '增加角色' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width:500px;">
      <el-form-item label="角色名称" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
        <el-button @click="$router.back()">返回</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/api/request'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  description: '',
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
}

async function loadRole() {
  if (!isEdit.value) return
  const res: any = await request.get('/admin/role/' + route.params.id)
  if (res.data) {
    form.title = res.data.title
    form.description = res.data.description
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = new URLSearchParams()
    data.append('title', form.title)
    data.append('description', form.description)

    if (isEdit.value) {
      await request.put('/admin/role/' + route.params.id, data)
      ElMessage.success('修改成功')
    } else {
      await request.post('/admin/role', data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/role')
  } finally {
    submitting.value = false
  }
}

onMounted(loadRole)
</script>
