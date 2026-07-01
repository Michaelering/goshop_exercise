<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h2>{{ isMerchant ? '商户管理中心' : '后台管理系统' }}</h2>
        <p>{{ isMerchant ? '请使用商户账号登录' : '请使用管理员账号登录' }}</p>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            :prefix-icon="User"
            placeholder="请输入账号"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :prefix-icon="Lock"
            placeholder="请输入密码"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item prop="verifyValue">
          <el-input
            v-model="form.verifyValue"
            :prefix-icon="Key"
            placeholder="验证码"
            size="large"
            class="captcha-input"
          >
            <template #append>
              <img
                :src="captchaImg"
                class="captcha-img"
                @click="loadCaptcha"
                alt="验证码"
                title="点击刷新"
              />
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            class="login-btn"
            size="large"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
        <div class="login-footer">
          <a href="javascript:;" @click="toggleMode">
            {{ isMerchant ? '管理员登录 →' : '商户登录 →' }}
          </a>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, Lock, Key } from '@element-plus/icons-vue'
import { getCaptcha } from '@/api/auth'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const captchaId = ref('')
const captchaImg = ref('')
const isMerchant = ref(false)

const form = reactive({
  username: '',
  password: '',
  verifyValue: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  verifyValue: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
}

async function loadCaptcha() {
  try {
    const res: any = await getCaptcha()
    captchaId.value = res.captchaId
    captchaImg.value = res.captchaImage
  } catch {
    // ignore
  }
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const type = isMerchant.value ? 'merchant' : 'admin'
    const success = await authStore.login({
      ...form,
      captchaId: captchaId.value,
    }, type)

    if (success) {
      ElMessage.success('登录成功')
      router.push(type === 'admin' ? '/admin/dashboard' : '/merchant/dashboard')
    }
  } finally {
    loading.value = false
    loadCaptcha() // 刷新验证码
  }
}

function toggleMode() {
  isMerchant.value = !isMerchant.value
  form.verifyValue = ''
  loadCaptcha()
}

onMounted(() => {
  loadCaptcha()
})
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0d47a1 0%, #1976d2 50%, #42a5f5 100%);
}
.login-card {
  width: 420px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}
.login-header {
  background: linear-gradient(135deg, #0d47a1, #1976d2);
  color: #fff;
  padding: 30px;
  text-align: center;
}
.login-header h2 {
  margin: 0 0 8px;
  font-weight: 700;
  font-size: 22px;
}
.login-header p {
  margin: 0;
  opacity: 0.85;
  font-size: 14px;
}
.login-form {
  padding: 30px;
}
.captcha-input :deep(.el-input-group__append) {
  padding: 0;
  overflow: hidden;
}
.captcha-img {
  height: 38px;
  cursor: pointer;
  display: block;
}
.login-btn {
  width: 100%;
  font-size: 16px;
  padding: 22px 0;
}
.login-footer {
  text-align: center;
  margin-top: 8px;
}
.login-footer a {
  color: #409eff;
  font-size: 14px;
}
</style>

<script lang="ts">
import { ElMessage } from 'element-plus'
export default {}
</script>
