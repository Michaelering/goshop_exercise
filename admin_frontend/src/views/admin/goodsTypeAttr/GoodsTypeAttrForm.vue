<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '修改商品类型属性' : '增加商品类型属性' }}</span>
    </template>

    <el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="max-width:600px;">
      <el-form-item label="所属类型" prop="cate_id">
        <el-select v-model="form.cate_id" placeholder="选择所属类型" style="width:100%">
          <el-option
            v-for="t in goodsTypeList"
            :key="t.id"
            :label="t.title"
            :value="t.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="属性标题" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="属性类型" prop="attr_type">
        <el-select v-model="form.attr_type" placeholder="选择属性类型" style="width:100%">
          <el-option label="文本" :value="1" />
          <el-option label="文本域" :value="2" />
          <el-option label="下拉" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item label="属性值" prop="attr_value">
        <el-input
          v-model="form.attr_value"
          type="textarea"
          :rows="4"
          :placeholder="form.attr_type === 3 ? '每行一个可选值' : ''"
        />
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
const goodsTypeList = ref<any[]>([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  cate_id: null as number | null,
  title: '',
  attr_type: 1 as number,
  attr_value: '',
  sort: 10,
  status: 1,
})

const rules: FormRules = {
  cate_id: [{ required: true, message: '请选择所属类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入属性标题', trigger: 'blur' }],
  attr_type: [{ required: true, message: '请选择属性类型', trigger: 'change' }],
}

async function loadGoodsTypes() {
  const res: any = await api.getGoodsTypeList()
  goodsTypeList.value = res.data || []
}

async function loadData() {
  if (!isEdit.value) return
  const id = route.params.id
  const res: any = await api.getGoodsTypeAttr(Number(id))
  const item = res.data?.goodsTypeAttr || res.data
  if (item) {
    form.cate_id = item.cate_id
    form.title = item.title
    form.attr_type = item.attr_type
    form.attr_value = item.attr_value || ''
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
    data.append('cate_id', String(form.cate_id))
    data.append('title', form.title)
    data.append('attr_type', String(form.attr_type))
    data.append('attr_value', form.attr_value)
    data.append('sort', String(form.sort))
    data.append('status', String(form.status))

    if (isEdit.value) {
      await api.updateGoodsTypeAttr(Number(route.params.id), data)
      ElMessage.success('修改成功')
    } else {
      await api.createGoodsTypeAttr(data)
      ElMessage.success('添加成功')
    }
    router.push('/admin/goodsTypeAttr')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadGoodsTypes()
  loadData()
})
</script>
