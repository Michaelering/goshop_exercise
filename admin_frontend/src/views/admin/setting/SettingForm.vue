<template>
  <el-card>
    <template #header>
      <span>系统设置</span>
    </template>

    <el-form ref="formRef" :model="form" label-width="120px" style="max-width:800px;">
      <el-divider content-position="left">网站信息</el-divider>

      <el-form-item label="网站标题" prop="site_title">
        <el-input v-model="form.site_title" />
      </el-form-item>
      <el-form-item label="网站LOGO" prop="site_logo">
        <el-input v-model="form.site_logo" placeholder="请输入LOGO图片URL" />
      </el-form-item>
      <el-form-item label="网站关键字" prop="site_keywords">
        <el-input v-model="form.site_keywords" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item label="网站描述" prop="site_description">
        <el-input v-model="form.site_description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="默认图片" prop="no_picture">
        <el-input v-model="form.no_picture" placeholder="请输入默认图片URL" />
      </el-form-item>
      <el-form-item label="备案号" prop="site_icp">
        <el-input v-model="form.site_icp" />
      </el-form-item>
      <el-form-item label="联系电话" prop="site_tel">
        <el-input v-model="form.site_tel" />
      </el-form-item>
      <el-form-item label="搜索关键词" prop="search_keywords">
        <el-input v-model="form.search_keywords" />
      </el-form-item>
      <el-form-item label="统计代码" prop="tongji_code">
        <el-input v-model="form.tongji_code" type="textarea" :rows="3" />
      </el-form-item>

      <el-divider content-position="left">OSS配置</el-divider>

      <el-form-item label="AppId" prop="appid">
        <el-input v-model="form.appid" />
      </el-form-item>
      <el-form-item label="AppSecret" prop="app_secret">
        <el-input v-model="form.app_secret" />
      </el-form-item>
      <el-form-item label="EndPoint" prop="end_point">
        <el-input v-model="form.end_point" />
      </el-form-item>
      <el-form-item label="Bucket" prop="bucket_name">
        <el-input v-model="form.bucket_name" />
      </el-form-item>
      <el-form-item label="OSS状态" prop="oss_status">
        <el-radio-group v-model="form.oss_status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="OSS域名" prop="oss_domain">
        <el-input v-model="form.oss_domain" />
      </el-form-item>
      <el-form-item label="缩略图尺寸" prop="thumbnail_size">
        <el-input v-model="form.thumbnail_size" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">保存设置</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import * as api from '@/api/admin'

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  site_title: '',
  site_logo: '',
  site_keywords: '',
  site_description: '',
  no_picture: '',
  site_icp: '',
  site_tel: '',
  search_keywords: '',
  tongji_code: '',
  appid: '',
  app_secret: '',
  end_point: '',
  bucket_name: '',
  oss_status: 0,
  oss_domain: '',
  thumbnail_size: '',
})

async function loadData() {
  const res: any = await api.getSetting()
  const item = res.data
  if (item) {
    form.site_title = item.site_title || ''
    form.site_logo = item.site_logo || ''
    form.site_keywords = item.site_keywords || ''
    form.site_description = item.site_description || ''
    form.no_picture = item.no_picture || ''
    form.site_icp = item.site_icp || ''
    form.site_tel = item.site_tel || ''
    form.search_keywords = item.search_keywords || ''
    form.tongji_code = item.tongji_code || ''
    form.appid = item.appid || ''
    form.app_secret = item.app_secret || ''
    form.end_point = item.end_point || ''
    form.bucket_name = item.bucket_name || ''
    form.oss_status = item.oss_status ?? 0
    form.oss_domain = item.oss_domain || ''
    form.thumbnail_size = item.thumbnail_size || ''
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    const data = new URLSearchParams()
    data.append('site_title', form.site_title)
    data.append('site_logo', form.site_logo)
    data.append('site_keywords', form.site_keywords)
    data.append('site_description', form.site_description)
    data.append('no_picture', form.no_picture)
    data.append('site_icp', form.site_icp)
    data.append('site_tel', form.site_tel)
    data.append('search_keywords', form.search_keywords)
    data.append('tongji_code', form.tongji_code)
    data.append('appid', form.appid)
    data.append('app_secret', form.app_secret)
    data.append('end_point', form.end_point)
    data.append('bucket_name', form.bucket_name)
    data.append('oss_status', String(form.oss_status))
    data.append('oss_domain', form.oss_domain)
    data.append('thumbnail_size', form.thumbnail_size)

    await api.updateSetting(data)
    ElMessage.success('保存成功')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>
