<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改商品分类' : '增加商品分类' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="分类图片" prop="cate_img">
        <el-input v-model="form.cate_img" placeholder="请输入图片URL" />
      </el-form-item>
      <el-form-item label="链接" prop="link">
        <el-input v-model="form.link" />
      </el-form-item>
      <el-form-item label="模板" prop="template">
        <el-input v-model="form.template" />
      </el-form-item>
      <el-form-item label="上级分类" prop="pid">
        <el-select v-model="form.pid" placeholder="选择上级分类" style="width:100%">
          <el-option label="顶级分类" :value="0" />
          <el-option
            v-for="item in cateList"
            :key="item.id"
            :label="item.title"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="副标题" prop="sub_title">
        <el-input v-model="form.sub_title" />
      </el-form-item>
      <el-form-item label="关键字" prop="keywords">
        <el-input v-model="form.keywords" />
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input v-model="form.description" type="textarea" :rows="2" />
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
const cateList = ref<any[]>([])
const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  cate_img: '',
  link: '',
  template: '',
  pid: 0,
  sub_title: '',
  keywords: '',
  description: '',
  sort: 10,
  status: 1,
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
}

async function loadCateList() {
  const res: any = await api.getGoodsCateList()
  cateList.value = res.data || []
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getGoodsCate(Number(id))
  const item = res.data
  if (item) {
    form.title = item.title || ''
    form.cate_img = item.cate_img || ''
    form.link = item.link || ''
    form.template = item.template || ''
    form.pid = item.pid ?? 0
    form.sub_title = item.sub_title || ''
    form.keywords = item.keywords || ''
    form.description = item.description || ''
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
    data.append('cate_img', form.cate_img)
    data.append('link', form.link)
    data.append('template', form.template)
    data.append('pid', String(form.pid))
    data.append('sub_title', form.sub_title)
    data.append('keywords', form.keywords)
    data.append('description', form.description)
    data.append('sort', String(form.sort))
    data.append('status', String(form.status))

    if (isEdit.value) {
      await api.updateGoodsCate(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createGoodsCate(data)
      ElMessage.success('添加成功')
    }
    router.push({ path: '/admin/goodsCate' })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadCateList()
  loadData()
})
</script>
