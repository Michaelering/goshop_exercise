<template>
  <el-card>
    <template #header><span>{{ isEdit ? '修改商品' : '增加商品' }}</span></template>
    <el-form ref="formRef" :model="form" label-width="120px">
      <el-form-item label="商品名称" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="商品编号">
        <el-input v-model="form.goodsSn" />
      </el-form-item>
      <el-form-item label="商品分类" prop="cateId">
        <el-input-number v-model="form.cateId" :min="0" />
      </el-form-item>
      <el-form-item label="商品价格" prop="price">
        <el-input-number v-model="form.price" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="市场价格">
        <el-input-number v-model="form.marketPrice" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="商品库存">
        <el-input-number v-model="form.goodsNumber" :min="0" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="form.sort" :min="0" />
      </el-form-item>
      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="关键词">
        <el-input v-model="form.goodsKeywords" />
      </el-form-item>
      <el-form-item label="商品描述">
        <el-input v-model="form.goodsDesc" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="商品内容">
        <el-input v-model="form.goodsContent" type="textarea" :rows="8" />
      </el-form-item>
      <el-form-item style="margin-top:20px">
        <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
        <el-button @click="$router.back()">返回</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const formRef = ref()
const submitting = ref(false)

const form = reactive({
  title: '',
  goodsSn: '',
  cateId: 0,
  price: 0,
  marketPrice: 0,
  goodsNumber: 0,
  sort: 0,
  status: 1,
  goodsKeywords: '',
  goodsDesc: '',
  goodsContent: '',
})

async function loadGoods(id: number) {
  try {
    const res: any = await request.get('/merchant/goods/' + id)
    const goods = res.data?.goods || res.data || {}
    form.title = goods.title || ''
    form.goodsSn = goods.goodsSn || ''
    form.cateId = goods.cateId || 0
    form.price = goods.price || 0
    form.marketPrice = goods.marketPrice || 0
    form.goodsNumber = goods.goodsNumber || 0
    form.sort = goods.sort || 0
    form.status = goods.status ?? 1
    form.goodsKeywords = goods.goodsKeywords || ''
    form.goodsDesc = goods.goodsDesc || ''
    form.goodsContent = goods.goodsContent || ''
  } catch (e) {
    // ignore
  }
}

function buildParams(): URLSearchParams {
  const params = new URLSearchParams()
  params.append('title', form.title)
  params.append('goodsSn', form.goodsSn)
  params.append('cateId', String(form.cateId))
  params.append('price', String(form.price))
  params.append('marketPrice', String(form.marketPrice))
  params.append('goodsNumber', String(form.goodsNumber))
  params.append('sort', String(form.sort))
  params.append('status', String(form.status))
  params.append('goodsKeywords', form.goodsKeywords)
  params.append('goodsDesc', form.goodsDesc)
  params.append('goodsContent', form.goodsContent)
  return params
}

async function handleSubmit() {
  submitting.value = true
  try {
    const params = buildParams()
    if (isEdit.value) {
      await request.put('/merchant/goods/' + route.params.id, params)
      ElMessage.success('修改成功')
    } else {
      await request.post('/merchant/goods', params)
      ElMessage.success('添加成功')
    }
    router.back()
  } catch (e) {
    // error handled by request interceptor
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (isEdit.value) {
    loadGoods(Number(route.params.id))
  }
})
</script>
