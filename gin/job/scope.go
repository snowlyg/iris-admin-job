package job

import "gorm.io/gorm"

// NameScope 根据 name 查询
// - name 任务名称
func NameScope(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

// StatusScope 根据 status 查询
// - status 数据status
func StatusScope(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}
