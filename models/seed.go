package models

import "fmt"

// RunMigrationAndSeed 播种内置角色和权限数据
// 注意：数据库迁移功能已注释，如需重新启用请取消注释 runMigration()
func RunMigrationAndSeed() {
	// runMigration()
	SeedBuiltinData()
}

// runMigration 执行数据库 DDL 迁移，适配新模型字段
// 迁移完成后会执行数据库修复和清理操作
// 当前已禁用：如需使用请取消函数首尾注释
/*
func runMigration() {
	// 使用 GORM AutoMigrate 创建新字段（url_prefix, http_methods, parent_id, is_builtin）
	if err := DB.AutoMigrate(&Role{}, &Access{}); err != nil {
		fmt.Println("AutoMigrate 警告:", err)
	}

	migrator := DB.Migrator()

	// ---- access 表：为 parent_id 创建索引（修复慢查询）----
	if migrator.HasTable(&Access{}) && migrator.HasColumn(&Access{}, "parent_id") {
		ensureIndex("access", "idx_access_parent_id", "parent_id")
	}

	// ---- role 表 ----
	if migrator.HasTable(&Role{}) && migrator.HasColumn(&Role{}, "is_builtin") {
		DB.Exec("UPDATE role SET is_builtin = 0 WHERE is_builtin IS NULL")
	}

	// ---- access 表：检测是否为旧数据结构，若是则清空后重播 ----
	if migrator.HasColumn(&Access{}, "url") && migrator.HasColumn(&Access{}, "url_prefix") {
		// 检查 url_prefix 是否全为空（说明旧数据未正确迁移）
		var emptyCount int64
		DB.Model(&Access{}).Where("type = 1 AND (url_prefix = '' OR url_prefix IS NULL)").Count(&emptyCount)
		var totalType1 int64
		DB.Model(&Access{}).Where("type = 1").Count(&totalType1)

		if totalType1 > 0 && emptyCount == totalType1 {
			// 旧数据结构，url_prefix 全为空 → 清空并让种子重新创建
			fmt.Println("检测到旧 access 数据，正在清空...")
			DB.Exec("DELETE FROM role_access")
			DB.Exec("DELETE FROM access")
			fmt.Println("旧 access / role_access 数据已清空，将由种子重建")
		} else if migrator.HasColumn(&Access{}, "url") {
			// url_prefix 已有数据，执行增量迁移
			DB.Exec("UPDATE access SET url_prefix = url WHERE (url_prefix = '' OR url_prefix IS NULL) AND url != ''")
			fmt.Println("access.url → url_prefix 增量迁移完成")
		}
	}

	// ---- access 表：迁移 module_id → parent_id ----
	if migrator.HasColumn(&Access{}, "module_id") {
		if migrator.HasColumn(&Access{}, "parent_id") {
			DB.Exec("UPDATE access SET parent_id = module_id WHERE parent_id = 0 AND module_id > 0")
			fmt.Println("access.module_id → parent_id 数据迁移完成")
		}
	}

	// ---- access 表：设置 http_methods 默认值 ----
	if migrator.HasColumn(&Access{}, "http_methods") {
		DB.Exec("UPDATE access SET http_methods = 'GET,POST,PUT,DELETE' WHERE type = 1 AND (http_methods = '' OR http_methods IS NULL)")
	}

	// ---- manager 表：处理 is_super（迁移到角色） ----
	if migrator.HasTable(&Manager{}) && migrator.HasColumn(&Manager{}, "is_super") {
		// 所有 is_super=1 的管理员都迁移到超管角色
		DB.Exec("UPDATE manager SET role_id = 1 WHERE is_super = 1")
		fmt.Println("manager.is_super → role_id=1 数据迁移完成")
	}

	fmt.Println("数据库迁移完成")

	// 修复旧自定义角色：给有管理员引用但无权限的自定义角色分配管理员权限
	fixOrphanRoles()

	// 清理不应该有子项的模块（仪表盘、系统设置只有一个子菜单，冗余）
	cleanupOrphanChildren()
}
*/

// SeedBuiltinData 播种内置角色和权限数据
func SeedBuiltinData() {
	seedBuiltinRoles()
	seedBuiltinAccess()
	seedSuperAdminPermissions()
	seedAdminPermissions()
}

// seedBuiltinRoles 创建/更新内置角色
func seedBuiltinRoles() {
	var count int64
	DB.Model(&Role{}).Where("is_builtin = 1").Count(&count)
	if count >= 2 {
		return // 已播种
	}

	builtinRoles := []Role{
		{Id: 1, Title: "超级管理员", Description: "系统内置超级管理员，拥有所有权限", Status: 1, IsBuiltin: 1, AddTime: int(GetUnix())},
		{Id: 2, Title: "管理员", Description: "平台管理员，管理商品、分类、订单、商户等", Status: 1, IsBuiltin: 1, AddTime: int(GetUnix())},
	}
	for _, role := range builtinRoles {
		// 使用 Save 而非 Create，处理 id=1/2 可能已被占用的情况
		var existing Role
		DB.Where("id = ?", role.Id).Find(&existing)
		if existing.Id > 0 {
			DB.Model(&existing).Updates(map[string]interface{}{
				"title":       role.Title,
				"description": role.Description,
				"is_builtin":  1,
				"status":      1,
			})
		} else {
			DB.Create(&role)
		}
	}
	fmt.Println("内置角色播种完成")
}

// access 种子数据定义
type accessSeed struct {
	ModuleName  string
	ActionName  string
	Type        int    // 1=模块组 2=菜单项
	UrlPrefix   string // 模块组才有
	HttpMethods string // 模块组才有
	ParentId    int    // 种子内的临时ID，会被替换为真实ID
	Sort        int
	TempId      int // 临时ID，用于建立父子关系
}

// seedBuiltinAccess 创建内置权限节点
func seedBuiltinAccess() {
	var count int64
	DB.Model(&Access{}).Count(&count)
	if count > 0 {
		return // 已播种
	}

	seeds := []accessSeed{
		// 仪表盘
		{TempId: 1, ModuleName: "仪表盘", ActionName: "", Type: 1, UrlPrefix: "dashboard", HttpMethods: "GET", Sort: 90},
		// 仪表盘是单页模块，无需子菜单

		// 管理员管理
		{TempId: 10, ModuleName: "管理员管理", ActionName: "", Type: 1, UrlPrefix: "manager", HttpMethods: "GET,POST,PUT,DELETE", Sort: 80},
		{TempId: 11, ModuleName: "管理员管理", ActionName: "管理员列表", Type: 2, ParentId: 10, Sort: 80},
		{TempId: 12, ModuleName: "管理员管理", ActionName: "增加管理员", Type: 2, ParentId: 10, Sort: 79},

		// 角色管理
		{TempId: 20, ModuleName: "角色管理", ActionName: "", Type: 1, UrlPrefix: "role", HttpMethods: "GET,POST,PUT,DELETE", Sort: 75},
		{TempId: 21, ModuleName: "角色管理", ActionName: "角色列表", Type: 2, ParentId: 20, Sort: 75},
		{TempId: 22, ModuleName: "角色管理", ActionName: "增加角色", Type: 2, ParentId: 20, Sort: 74},

		// 权限管理
		{TempId: 30, ModuleName: "权限管理", ActionName: "", Type: 1, UrlPrefix: "access", HttpMethods: "GET,POST,PUT,DELETE", Sort: 70},
		{TempId: 31, ModuleName: "权限管理", ActionName: "权限列表", Type: 2, ParentId: 30, Sort: 70},
		{TempId: 32, ModuleName: "权限管理", ActionName: "增加权限", Type: 2, ParentId: 30, Sort: 69},

		// 商品管理
		{TempId: 40, ModuleName: "商品管理", ActionName: "", Type: 1, UrlPrefix: "goods", HttpMethods: "GET,POST,PUT,DELETE", Sort: 65},
		{TempId: 41, ModuleName: "商品管理", ActionName: "商品列表", Type: 2, ParentId: 40, Sort: 65},
		{TempId: 42, ModuleName: "商品管理", ActionName: "增加商品", Type: 2, ParentId: 40, Sort: 64},

		// 商品分类
		{TempId: 50, ModuleName: "商品分类", ActionName: "", Type: 1, UrlPrefix: "goodsCate", HttpMethods: "GET,POST,PUT,DELETE", Sort: 60},
		{TempId: 51, ModuleName: "商品分类", ActionName: "分类列表", Type: 2, ParentId: 50, Sort: 60},
		{TempId: 52, ModuleName: "商品分类", ActionName: "增加分类", Type: 2, ParentId: 50, Sort: 59},

		// 商品类型
		{TempId: 60, ModuleName: "商品类型", ActionName: "", Type: 1, UrlPrefix: "goodsType", HttpMethods: "GET,POST,PUT,DELETE", Sort: 55},
		{TempId: 61, ModuleName: "商品类型", ActionName: "类型列表", Type: 2, ParentId: 60, Sort: 55},
		{TempId: 62, ModuleName: "商品类型", ActionName: "增加类型", Type: 2, ParentId: 60, Sort: 54},

		// 类型属性
		{TempId: 70, ModuleName: "类型属性", ActionName: "", Type: 1, UrlPrefix: "goodsTypeAttr", HttpMethods: "GET,POST,PUT,DELETE", Sort: 50},
		{TempId: 71, ModuleName: "类型属性", ActionName: "属性列表", Type: 2, ParentId: 70, Sort: 50},
		{TempId: 72, ModuleName: "类型属性", ActionName: "增加属性", Type: 2, ParentId: 70, Sort: 49},

		// 导航管理
		{TempId: 80, ModuleName: "导航管理", ActionName: "", Type: 1, UrlPrefix: "nav", HttpMethods: "GET,POST,PUT,DELETE", Sort: 45},
		{TempId: 81, ModuleName: "导航管理", ActionName: "导航列表", Type: 2, ParentId: 80, Sort: 45},
		{TempId: 82, ModuleName: "导航管理", ActionName: "增加导航", Type: 2, ParentId: 80, Sort: 44},

		// 轮播图管理
		{TempId: 90, ModuleName: "轮播图管理", ActionName: "", Type: 1, UrlPrefix: "focus", HttpMethods: "GET,POST,PUT,DELETE", Sort: 40},
		{TempId: 91, ModuleName: "轮播图管理", ActionName: "轮播图列表", Type: 2, ParentId: 90, Sort: 40},
		{TempId: 92, ModuleName: "轮播图管理", ActionName: "增加轮播图", Type: 2, ParentId: 90, Sort: 39},

		// 系统设置
		{TempId: 100, ModuleName: "系统设置", ActionName: "", Type: 1, UrlPrefix: "setting", HttpMethods: "GET,PUT", Sort: 35},
		// 系统设置是单页模块，无需子菜单

		// 商户管理
		{TempId: 110, ModuleName: "商户管理", ActionName: "", Type: 1, UrlPrefix: "merchant", HttpMethods: "GET,POST,PUT,DELETE", Sort: 30},
		{TempId: 111, ModuleName: "商户管理", ActionName: "商户列表", Type: 2, ParentId: 110, Sort: 30},
		{TempId: 112, ModuleName: "商户管理", ActionName: "增加商户", Type: 2, ParentId: 110, Sort: 29},
	}

	// 先创建所有 Type=1 的模块组
	tempToReal := make(map[int]int)
	for _, s := range seeds {
		if s.Type == 1 {
			access := Access{
				ModuleName:  s.ModuleName,
				ActionName:  s.ActionName,
				Type:        s.Type,
				UrlPrefix:   s.UrlPrefix,
				HttpMethods: s.HttpMethods,
				ParentId:    0,
				Sort:        s.Sort,
				Status:      1,
				AddTime:     int(GetUnix()),
			}
			DB.Create(&access)
			tempToReal[s.TempId] = access.Id
		}
	}

	// 再创建 Type=2 的菜单项
	for _, s := range seeds {
		if s.Type == 2 {
			realParentId := tempToReal[s.ParentId]
			access := Access{
				ModuleName:  s.ModuleName,
				ActionName:  s.ActionName,
				Type:        s.Type,
				UrlPrefix:   "",
				HttpMethods: "",
				ParentId:    realParentId,
				Sort:        s.Sort,
				Status:      1,
				AddTime:     int(GetUnix()),
			}
			DB.Create(&access)
			tempToReal[s.TempId] = access.Id
		}
	}

	fmt.Println("内置权限节点播种完成")
}

// seedSuperAdminPermissions 给超级管理员(id=1)分配所有权限
func seedSuperAdminPermissions() {
	var count int64
	DB.Model(&RoleAccess{}).Where("role_id = 1").Count(&count)
	if count > 0 {
		return
	}

	var accessList []Access
	DB.Where("type = 1").Find(&accessList)
	for _, a := range accessList {
		DB.Create(&RoleAccess{RoleId: 1, AccessId: a.Id})
	}
	fmt.Println("超级管理员权限播种完成")
}

// seedAdminPermissions 给管理员(id=2)分配运营权限
func seedAdminPermissions() {
	var count int64
	DB.Model(&RoleAccess{}).Where("role_id = 2").Count(&count)
	if count > 0 {
		return
	}

	// 管理员可访问的模块（排除 manager, role, access）
	adminPrefixes := []string{"dashboard", "goods", "goodsCate", "goodsType", "goodsTypeAttr", "nav", "focus", "setting", "merchant"}
	for _, prefix := range adminPrefixes {
		var access Access
		DB.Where("url_prefix = ? AND type = 1", prefix).Find(&access)
		if access.Id > 0 {
			DB.Create(&RoleAccess{RoleId: 2, AccessId: access.Id})
		}
	}
	fmt.Println("管理员权限播种完成")
}

// cleanupOrphanChildren 删除单页模块下不应存在的子菜单（旧种子数据残留）
// 当前已禁用
/*
func cleanupOrphanChildren() {
	DB.Exec("DELETE FROM access WHERE type = 2 AND parent_id IN (SELECT id FROM (SELECT id FROM access WHERE url_prefix IN ('dashboard','setting') AND type = 1) AS tmp)")
}
*/

// fixOrphanRoles 给旧自定义角色（有管理员在使用但无任何权限）自动分配管理员权限
// 当前已禁用
/*
func fixOrphanRoles() {
	// 找出所有被 manager 引用的角色
	rows, err := DB.Raw("SELECT DISTINCT role_id FROM manager WHERE role_id > 0 AND role_id NOT IN (1, 2)").Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	adminPrefixes := []string{"dashboard", "goods", "goodsCate", "goodsType", "goodsTypeAttr", "nav", "focus", "setting", "merchant"}
	for rows.Next() {
		var roleId int
		rows.Scan(&roleId)
		// 检查此角色是否有任何权限
		var accessCount int64
		DB.Model(&RoleAccess{}).Where("role_id = ?", roleId).Count(&accessCount)
		if accessCount == 0 {
			fmt.Printf("修复孤儿角色 id=%d: 自动分配管理员权限\n", roleId)
			for _, prefix := range adminPrefixes {
				var access Access
				DB.Where("url_prefix = ? AND type = 1", prefix).Find(&access)
				if access.Id > 0 {
					DB.Create(&RoleAccess{RoleId: roleId, AccessId: access.Id})
				}
			}
		}
	}
}
*/

// ensureIndex 安全创建索引（如果存在则跳过）
// 当前已禁用
/*
func ensureIndex(table, indexName, column string) {
	// 检查索引是否已存在
	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", table, indexName).Scan(&count)
	if count == 0 {
		sql := "CREATE INDEX " + indexName + " ON " + table + " (" + column + ")"
		if err := DB.Exec(sql).Error; err != nil {
			fmt.Printf("创建索引 %s 失败（可能已存在）: %v\n", indexName, err)
		} else {
			fmt.Printf("索引 %s 创建成功\n", indexName)
		}
	}
}
*/
