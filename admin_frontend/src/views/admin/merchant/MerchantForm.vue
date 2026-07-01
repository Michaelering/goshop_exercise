<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改商户' : '增加商户' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="登录密码" prop="password">
        <el-input v-model="form.password" type="password" show-password :placeholder="isEdit ? '留空则不修改密码' : ''" />
      </el-form-item>
      <el-form-item label="店铺名称" prop="shop_name">
        <el-input v-model="form.shop_name" />
      </el-form-item>
      <el-form-item label="手机号" prop="mobile">
        <el-input v-model="form.mobile" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
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
  username: '',
  password: '',
  shop_name: '',
  mobile: '',
  email: '',
  status: 1,
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, message: '用户名至少2个字符', trigger: 'blur' },
  ],
  shop_name: [{ required: true, message: '请输入店铺名称', trigger: 'blur' }],
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getMerchant(Number(id))
  const item = res.data?.merchant || res.data
  if (item) {
    form.username = item.username
    form.shop_name = item.shop_name || ''
    form.mobile = item.mobile || ''
    form.email = item.email || ''
    form.status = item.status ?? 1
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = new URLSearchParams()
    data.append('username', form.username)
    data.append('password', form.password || '')
    data.append('shop_name', form.shop_name)
    data.append('mobile', form.mobile)
    data.append('email', form.email)
    data.append('status', String(form.status))

    if (isEdit.value) {
      await api.updateMerchant(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createMerchant(data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/merchant')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
