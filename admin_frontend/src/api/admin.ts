import request from './request'

// ====== 权限管理 (Access) ======
export function getAccessList() {
  return request.get('/admin/access')
}
export function getAccessTopModules() {
  return request.get('/admin/access/topModules')
}
export function getAccess(id: number) {
  return request.get('/admin/access/' + id)
}
export function createAccess(data: URLSearchParams) {
  return request.post('/admin/access', data)
}
export function updateAccess(id: number, data: URLSearchParams) {
  return request.put('/admin/access/' + id, data)
}
export function deleteAccess(id: number) {
  return request.delete('/admin/access/' + id)
}

// ====== 商品分类 (GoodsCate) ======
export function getGoodsCateList() {
  return request.get('/admin/goodsCate')
}
export function getGoodsCateTree() {
  return request.get('/admin/goodsCate/tree')
}
export function getGoodsCate(id: number) {
  return request.get('/admin/goodsCate/' + id)
}
export function createGoodsCate(data: URLSearchParams) {
  return request.post('/admin/goodsCate', data)
}
export function updateGoodsCate(id: number, data: URLSearchParams) {
  return request.put('/admin/goodsCate/' + id, data)
}
export function deleteGoodsCate(id: number) {
  return request.delete('/admin/goodsCate/' + id)
}

// ====== 商品类型 (GoodsType) ======
export function getGoodsTypeList() {
  return request.get('/admin/goodsType')
}
export function getGoodsType(id: number) {
  return request.get('/admin/goodsType/' + id)
}
export function createGoodsType(data: URLSearchParams) {
  return request.post('/admin/goodsType', data)
}
export function updateGoodsType(id: number, data: URLSearchParams) {
  return request.put('/admin/goodsType/' + id, data)
}
export function deleteGoodsType(id: number) {
  return request.delete('/admin/goodsType/' + id)
}

// ====== 商品类型属性 (GoodsTypeAttr) ======
export function getGoodsTypeAttrList() {
  return request.get('/admin/goodsTypeAttr')
}
export function getGoodsTypeAttr(id: number) {
  return request.get('/admin/goodsTypeAttr/' + id)
}
export function getGoodsTypeAttrByType(cateId: number) {
  return request.get('/admin/goodsTypeAttr/byType/' + cateId)
}
export function createGoodsTypeAttr(data: URLSearchParams) {
  return request.post('/admin/goodsTypeAttr', data)
}
export function updateGoodsTypeAttr(id: number, data: URLSearchParams) {
  return request.put('/admin/goodsTypeAttr/' + id, data)
}
export function deleteGoodsTypeAttr(id: number) {
  return request.delete('/admin/goodsTypeAttr/' + id)
}

// ====== 导航管理 (Nav) ======
export function getNavList(params?: { page?: number; pageSize?: number }) {
  return request.get('/admin/nav', { params })
}
export function getNav(id: number) {
  return request.get('/admin/nav/' + id)
}
export function createNav(data: URLSearchParams) {
  return request.post('/admin/nav', data)
}
export function updateNav(id: number, data: URLSearchParams) {
  return request.put('/admin/nav/' + id, data)
}
export function deleteNav(id: number) {
  return request.delete('/admin/nav/' + id)
}

// ====== 轮播图管理 (Focus) ======
export function getFocusList() {
  return request.get('/admin/focus')
}
export function getFocus(id: number) {
  return request.get('/admin/focus/' + id)
}
export function createFocus(data: URLSearchParams) {
  return request.post('/admin/focus', data)
}
export function updateFocus(id: number, data: URLSearchParams) {
  return request.put('/admin/focus/' + id, data)
}
export function deleteFocus(id: number) {
  return request.delete('/admin/focus/' + id)
}

// ====== 系统设置 (Setting) ======
export function getSetting() {
  return request.get('/admin/setting')
}
export function updateSetting(data: URLSearchParams) {
  return request.put('/admin/setting', data)
}

// ====== 商户管理 (Merchant) ======
export function getMerchantList() {
  return request.get('/admin/merchant')
}
export function getMerchant(id: number) {
  return request.get('/admin/merchant/' + id)
}
export function createMerchant(data: URLSearchParams) {
  return request.post('/admin/merchant', data)
}
export function updateMerchant(id: number, data: URLSearchParams) {
  return request.put('/admin/merchant/' + id, data)
}
export function deleteMerchant(id: number) {
  return request.delete('/admin/merchant/' + id)
}

// ====== 角色管理 (Role) ======
export function getRoleList() {
  return request.get('/admin/role')
}
export function getRole(id: number) {
  return request.get('/admin/role/' + id)
}
export function createRole(data: URLSearchParams) {
  return request.post('/admin/role', data)
}
export function updateRole(id: number, data: URLSearchParams) {
  return request.put('/admin/role/' + id, data)
}
export function deleteRole(id: number) {
  return request.delete('/admin/role/' + id)
}
export function getRoleAuth(id: number) {
  return request.get('/admin/role/' + id + '/auth')
}
export function saveRoleAuth(id: number, data: URLSearchParams) {
  return request.post('/admin/role/' + id + '/auth', data)
}

// ====== 商品管理 (Goods) ======
export function getGoodsList(params?: { page?: number; pageSize?: number; keyword?: string }) {
  return request.get('/admin/goods', { params })
}
export function getGoods(id: number) {
  return request.get('/admin/goods/' + id)
}
export function createGoods(data: URLSearchParams) {
  return request.post('/admin/goods', data)
}
export function updateGoods(id: number, data: URLSearchParams) {
  return request.put('/admin/goods/' + id, data)
}
export function deleteGoods(id: number) {
  return request.delete('/admin/goods/' + id)
}

// ====== 通用操作 ======
export function changeStatus(table: string, field: string, id: number) {
  const data = new URLSearchParams()
  data.append('table', table)
  data.append('field', field)
  data.append('id', String(id))
  return request.put('/admin/changeStatus', data)
}

export function changeNum(table: string, field: string, id: number, num: number) {
  const data = new URLSearchParams()
  data.append('table', table)
  data.append('field', field)
  data.append('id', String(id))
  data.append('num', String(num))
  return request.put('/admin/changeNum', data)
}
