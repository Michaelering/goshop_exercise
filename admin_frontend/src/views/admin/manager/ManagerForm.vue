<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改管理员' : '增加管理员' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="管理员名称" prop="username">
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="登录密码" prop="password">
        <el-input v-model="form.password" type="password" show-password :placeholder="isEdit ? '留空则不修改密码' : ''" />
      </el-form-item>
      <el-form-item label="电话" prop="mobile">
        <el-input v-model="form.mobile" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="角色" prop="roleId">
        <el-select v-model="form.roleId" placeholder="选择角色" style="width:100%">
          <el-option
            v-for="r in roleList"
            :key="r.id"
            :label="r.title"
            :value="r.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ submitting ? '提交中...' : '提交' }}
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
import request from '@/api/request'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const roleList = ref<any[]>([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  username: '',
  password: '',
  mobile: '',
  email: '',
  roleId: null as number | null,
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入管理员名称', trigger: 'blur' }],
  roleId: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

async function loadRoles() {
  const res: any = await request.get('/admin/role')
  roleList.value = res.data || []
}

async function loadManager() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await request.get('/admin/manager/' + id)
  const mgr = res.data?.manager
  if (mgr) {
    form.username = mgr.username
    form.mobile = mgr.mobile
    form.email = mgr.email
    form.roleId = mgr.roleId
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
    data.append('mobile', form.mobile)
    data.append('email', form.email)
    data.append('role_id', String(form.roleId))

    if (isEdit.value) {
      await request.put('/admin/manager/' + route.params.id, data)
      ElMessage.success('修改成功')
    } else {
      await request.post('/admin/manager', data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/manager')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadRoles()
  loadManager()
})
</script>
