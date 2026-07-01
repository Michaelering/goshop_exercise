<template>
  <el-card>
    <template #header><span>{{ isEdit ? '修改商品' : '增加商品' }}</span></template>
    <el-form ref="formRef" :model="form" label-width="120px">
      <el-tabs>
        <el-tab-pane label="基本信息">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="商品名称" prop="title">
                <el-input v-model="form.title" />
              </el-form-item>
              <el-form-item label="副标题">
                <el-input v-model="form.subTitle" />
              </el-form-item>
              <el-form-item label="商品编号">
                <el-input v-model="form.goodsSn" />
              </el-form-item>
              <el-form-item label="商品分类" prop="cateId">
                <el-cascader
                  v-model="form.cateId"
                  :options="cateOptions"
                  :props="{ value: 'id', label: 'title', checkStrictly: true, emitPath: false }"
                  style="width:100%"
                />
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
            </el-col>
            <el-col :span="12">
              <el-form-item label="商品类型">
                <el-select v-model="form.goodsTypeId" placeholder="选择类型" style="width:100%">
                  <el-option v-for="t in goodsTypeList" :key="t.id" :label="t.title" :value="t.id" />
                </el-select>
              </el-form-item>
              <el-form-item label="商品颜色">
                <el-checkbox-group v-model="form.goodsColor">
                  <el-checkbox v-for="c in colorList" :key="c.id" :label="String(c.id)">
                    {{ c.colorName }}
                  </el-checkbox>
                </el-checkbox-group>
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
            </el-col>
          </el-row>
          <el-form-item label="商品标签">
            <el-checkbox v-model="form.isHot" :true-value="1" :false-value="0">热销</el-checkbox>
            <el-checkbox v-model="form.isBest" :true-value="1" :false-value="0">精品</el-checkbox>
            <el-checkbox v-model="form.isNew" :true-value="1" :false-value="0">新品</el-checkbox>
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="商品详情">
          <el-form-item label="关键词">
            <el-input v-model="form.goodsKeywords" />
          </el-form-item>
          <el-form-item label="商品描述">
            <el-input v-model="form.goodsDesc" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="商品内容">
            <el-input v-model="form.goodsContent" type="textarea" :rows="8" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="高级设置">
          <el-form-item label="关联商品">
            <el-input v-model="form.relationGoods" placeholder="商品ID逗号分隔" />
          </el-form-item>
          <el-form-item label="商品赠品">
            <el-input v-model="form.goodsGift" />
          </el-form-item>
          <el-form-item label="商品配件">
            <el-input v-model="form.goodsFitting" />
          </el-form-item>
          <el-form-item label="商品版本">
            <el-input v-model="form.goodsVersion" />
          </el-form-item>
          <el-form-item label="商品属性">
            <el-input v-model="form.goodsAttr" type="textarea" />
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
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
import * as api from '@/api/admin'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const formRef = ref()
const submitting = ref(false)

interface ColorItem {
  id: number
  colorName: string
}

const colorList: ColorItem[] = [
  { id: 1, colorName: '红色' },
  { id: 2, colorName: '蓝色' },
  { id: 3, colorName: '绿色' },
  { id: 4, colorName: '黑色' },
  { id: 5, colorName: '白色' },
  { id: 6, colorName: '灰色' },
  { id: 7, colorName: '金色' },
]

const goodsTypeList = ref<any[]>([])
const cateOptions = ref<any[]>([])

const form = reactive({
  title: '',
  subTitle: '',
  goodsSn: '',
  cateId: null as number | null,
  price: 0,
  marketPrice: 0,
  goodsNumber: 0,
  goodsTypeId: null as number | null,
  goodsColor: [] as string[],
  sort: 0,
  status: 1,
  isHot: 0,
  isBest: 0,
  isNew: 0,
  goodsKeywords: '',
  goodsDesc: '',
  goodsContent: '',
  relationGoods: '',
  goodsGift: '',
  goodsFitting: '',
  goodsVersion: '',
  goodsAttr: '',
})

function buildTree(list: any[]): any[] {
  const map: Record<number, any> = {}
  const tree: any[] = []
  list.forEach((item) => {
    map[item.id] = { ...item, children: [] }
  })
  list.forEach((item) => {
    if (item.pid && map[item.pid]) {
      map[item.pid].children.push(map[item.id])
    } else {
      tree.push(map[item.id])
    }
  })
  // 删除空的 children 数组，让 cascader 默认样式更好
  function cleanup(nodes: any[]) {
    nodes.forEach((node) => {
      if (node.children && node.children.length === 0) {
        delete node.children
      } else if (node.children) {
        cleanup(node.children)
      }
    })
  }
  cleanup(tree)
  return tree
}

async function loadGoodsCateList() {
  try {
    const res: any = await api.getGoodsCateList()
    const list = res.data?.list || res.data || []
    cateOptions.value = buildTree(list)
  } catch (e) {
    // ignore
  }
}

async function loadGoodsTypeList() {
  try {
    const res: any = await api.getGoodsTypeList()
    goodsTypeList.value = res.data?.list || res.data || []
  } catch (e) {
    // ignore
  }
}

async function loadGoods(id: number) {
  try {
    const res: any = await api.getGoods(id)
    const goods = res.data?.goods || res.data || {}
    form.title = goods.title || ''
    form.subTitle = goods.subTitle || ''
    form.goodsSn = goods.goodsSn || ''
    form.cateId = goods.cateId || null
    form.price = goods.price || 0
    form.marketPrice = goods.marketPrice || 0
    form.goodsNumber = goods.goodsNumber || 0
    form.goodsTypeId = goods.goodsTypeId || null
    form.sort = goods.sort || 0
    form.status = goods.status ?? 1
    form.isHot = goods.isHot ?? 0
    form.isBest = goods.isBest ?? 0
    form.isNew = goods.isNew ?? 0
    form.goodsKeywords = goods.goodsKeywords || ''
    form.goodsDesc = goods.goodsDesc || ''
    form.goodsContent = goods.goodsContent || ''
    form.relationGoods = goods.relationGoods || ''
    form.goodsGift = goods.goodsGift || ''
    form.goodsFitting = goods.goodsFitting || ''
    form.goodsVersion = goods.goodsVersion || ''
    form.goodsAttr = goods.goodsAttr || ''
    if (goods.goodsColor) {
      form.goodsColor = String(goods.goodsColor)
        .split(',')
        .filter((s: string) => s.trim())
    } else {
      form.goodsColor = []
    }
  } catch (e) {
    // ignore
  }
}

function buildParams(): URLSearchParams {
  const params = new URLSearchParams()
  params.append('title', form.title)
  params.append('subTitle', form.subTitle)
  params.append('goodsSn', form.goodsSn)
  if (form.cateId != null) params.append('cateId', String(form.cateId))
  params.append('price', String(form.price))
  params.append('marketPrice', String(form.marketPrice))
  params.append('goodsNumber', String(form.goodsNumber))
  if (form.goodsTypeId != null) params.append('goodsTypeId', String(form.goodsTypeId))
  params.append('goodsColor', form.goodsColor.join(','))
  params.append('sort', String(form.sort))
  params.append('status', String(form.status))
  params.append('isHot', String(form.isHot))
  params.append('isBest', String(form.isBest))
  params.append('isNew', String(form.isNew))
  params.append('goodsKeywords', form.goodsKeywords)
  params.append('goodsDesc', form.goodsDesc)
  params.append('goodsContent', form.goodsContent)
  params.append('relationGoods', form.relationGoods)
  params.append('goodsGift', form.goodsGift)
  params.append('goodsFitting', form.goodsFitting)
  params.append('goodsVersion', form.goodsVersion)
  params.append('goodsAttr', form.goodsAttr)
  params.append('merchant_id', '0')
  return params
}

async function handleSubmit() {
  submitting.value = true
  try {
    const params = buildParams()
    if (isEdit.value) {
      await api.updateGoods(Number(route.params.id), params)
      ElMessage.success('修改成功')
    } else {
      await api.createGoods(params)
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
  loadGoodsCateList()
  loadGoodsTypeList()
  if (isEdit.value) {
    loadGoods(Number(route.params.id))
  }
})
</script>
