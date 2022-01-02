package job

import (
	"fmt"

	"github.com/snowlyg/iris-admin/server/database/orm"
	"gorm.io/gorm"
)

type Response struct {
	orm.Model
	BaseJob
}

func (res *Response) First(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(&Job{}).Scopes(scopes...).First(res).Error
	if err != nil {
		return fmt.Errorf("获取失败:%w", err)
	}
	return nil
}

// Paginate 分页
type PageResponse struct {
	Item []*Response
}

func (res *PageResponse) Paginate(db *gorm.DB, pageScope func(db *gorm.DB) *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	db = db.Model(&Job{})
	var count int64
	err := db.Scopes(scopes...).Count(&count).Error
	if err != nil {
		return count, fmt.Errorf("获取总数失败:%w", err)
	}
	err = db.Scopes(pageScope).Find(&res.Item).Error
	if err != nil {
		return count, fmt.Errorf("获取分页数据失败:%w", err)
	}
	return count, nil
}

func (res *PageResponse) Find(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	db = db.Model(&Job{})
	err := db.Scopes(scopes...).Find(&res.Item).Error
	if err != nil {
		return fmt.Errorf("获取数据失败:%w", err)
	}
	return nil
}
