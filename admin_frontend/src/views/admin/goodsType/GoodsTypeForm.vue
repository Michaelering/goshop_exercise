<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改商品类型' : '增加商品类型' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" />
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

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  description: '',
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getGoodsType(Number(id))
  const item = res.data?.goodsType || res.data
  if (item) {
    form.title = item.title
    form.description = item.description || ''
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
      await api.updateGoodsType(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createGoodsType(data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/goodsType')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
