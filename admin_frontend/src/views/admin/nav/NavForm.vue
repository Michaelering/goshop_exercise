<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改导航' : '增加导航' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="链接" prop="link">
        <el-input v-model="form.link" />
      </el-form-item>
      <el-form-item label="位置" prop="position">
        <el-input-number v-model="form.position" :min="1" />
      </el-form-item>
      <el-form-item label="新窗口" prop="is_opennew">
        <el-radio-group v-model="form.is_opennew">
          <el-radio :value="0">否</el-radio>
          <el-radio :value="1">是</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="关联" prop="relation">
        <el-input v-model="form.relation" />
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
import type { FormInstance, FormRules } from 'element-plus'
import * as api from '@/api/admin'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  link: '',
  position: 2,
  is_opennew: 0,
  relation: '',
  sort: 10,
  status: 1,
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getNav(Number(id))
  const item = res.data?.nav
  if (item) {
    form.title = item.title || ''
    form.link = item.link || ''
    form.position = item.position ?? 2
    form.is_opennew = item.is_opennew ?? 0
    form.relation = item.relation || ''
    form.sort = item.sort ?? 10
    form.status = item.status ?? 1
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = new URLSearchParams()
    data.append('title', form.title)
    data.append('link', form.link)
    data.append('position', String(form.position))
    data.append('is_opennew', String(form.is_opennew))
    data.append('relation', form.relation)
    data.append('sort', String(form.sort))
    data.append('status', String(form.status))

    if (isEdit.value) {
      await api.updateNav(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createNav(data)
      ElMessage.success('添加成功')
    }
    router.push({ path: '/admin/nav' })
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
