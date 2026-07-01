<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改轮播图' : '增加轮播图' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="类型" prop="focus_type">
        <el-radio-group v-model="form.focus_type">
          <el-radio :value="1">PC</el-radio>
          <el-radio :value="2">移动</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="图片" prop="focus_img">
        <el-input v-model="form.focus_img" placeholder="请输入图片URL" />
      </el-form-item>
      <el-form-item label="链接" prop="link">
        <el-input v-model="form.link" />
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
  focus_type: 1,
  focus_img: '',
  link: '',
  sort: 10,
  status: 1,
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getFocus(Number(id))
  const item = res.data?.focus || res.data
  if (item) {
    form.title = item.title || ''
    form.focus_type = item.focus_type ?? 1
    form.focus_img = item.focus_img || ''
    form.link = item.link || ''
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
    data.append('focus_type', String(form.focus_type))
    data.append('focus_img', form.focus_img)
    data.append('link', form.link)
    data.append('sort', String(form.sort))
    data.append('status', String(form.status))

    if (isEdit.value) {
      await api.updateFocus(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createFocus(data)
      ElMessage.success('添加成功')
    }
    router.push({ path: '/admin/focus' })
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
