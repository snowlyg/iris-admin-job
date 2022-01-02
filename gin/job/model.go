package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	v3 "github.com/robfig/cron/v3"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"gorm.io/gorm"
)

// Job
type Job struct {
	gorm.Model
	BaseJob
}

type BaseJob struct {
	EntryId v3.EntryID `gorm:"uniqueIndex:entry_id;column:entry_id;type:char(15);default:0" json:"entryId" binding:"required"`  // 任务id
	Ip      string     `gorm:"uniqueIndex:service_name;column:ip;type:char(15);default:''" json:"ip"`                           // 设备内部ip
	Name    string     `gorm:"uniqueIndex:service_name;column:name;type:varchar(30);default:''" json:"name" binding:"required"` // 任务名称
	Spec    string     `gorm:"column:spec;type:varchar(10);default:''" json:"spec" binding:"required"`                          // 任务cron 配置
	Status  string     `gorm:"column:status;type:varchar(10);default:''" json:"status" binding:"required"`                      // 任务状态
	Desc    string     `gorm:"column:desc;type:varchar(50);default:''" json:"desc"`                                             // 任务描述
	Count   int        `gorm:"column:count;type:bigint;default:0" json:"count"`
	Message string     `gorm:"column:message;type:varchar(2000);default:''" json:"message"` // 执行信息
	LastAt  *time.Time `json:"lastAt"`
}

// Create 添加
func (item *Job) Create(db *gorm.DB) (uint, error) {
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return item.ID, fmt.Errorf("添加失败:%w", err)
	}
	return item.ID, nil
}

// Update 更新
func (item *Job) Update(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	data := map[string]interface{}{
		"status": item.Status,
		"desc":   item.Desc,
		"name":   item.Name,
		"spec":   item.Spec,
		"ip":     item.Ip,
	}
	err := db.Model(item).Scopes(scopes...).Updates(data).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return fmt.Errorf("更新失败:%w", err)
	}
	return nil
}

// Delete 删除
func (item *Job) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return fmt.Errorf("删除失败:%w", err)
	}
	return nil
}

// 定义内置任务
type BuiltinJob struct {
	sync.RWMutex
	jobs map[string]cron.Job
}
