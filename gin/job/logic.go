package job

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/snowlyg/iris-admin/server/cron_server"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"gorm.io/gorm"
)

const DefaultCronJobSpec = "@every 5m"

var ErrCronExpression = errors.New("错误的CRON表达式")

// StartJob
func StartJob() {
	cron_server.CronInstance().Start()
}

// StopJob
func StopJob() {
	cron_server.CronInstance().Stop()
	clearJob()
}

// OnceJob
func OnceJob(cmd cron.Job, d time.Duration) {
	once := time.Now().Add(d)
	onceSpec := fmt.Sprintf("%d %d %d %d %d %d", once.Second(), once.Minute(), once.Hour(), once.Day(), once.Month(), once.Weekday())
	cron_server.CronInstance().AddJob(onceSpec, cmd)
	cron_server.CronInstance().Run()
}

// clearJob
func clearJob() error {
	err := database.Instance().Unscoped().Where("1=1").Delete(&Job{}).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}
	return nil
}

// LogicExecJob 执行任务
func LogicExecJob(id uint) error {
	response := &Response{}
	err := response.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}

	if response.Status == "running" {
		return cron_server.ErrStartServer
	}

	once := time.Now().Add(2 * time.Second)
	onceSpec := fmt.Sprintf("%d %d %d %d %d %d", once.Second(), once.Minute(), once.Hour(), once.Day(), once.Month(), once.Weekday())
	_, err = cron_server.CronInstance().AddJob(onceSpec, BuiltinJobs.GetBuiltinJob(response.Name))
	if err != nil {
		return err
	}
	cron_server.CronInstance().Run()

	return nil
}

// LogicCreate 添加
func LogicCreate(req *Request) (uint, error) {
	j := &Job{BaseJob: req.BaseJob}
	return j.Create(database.Instance())
}

// LogicUpdate 更新
func LogicUpdate(id uint, req *Request) error {
	j := &Job{BaseJob: req.BaseJob}
	return j.Update(database.Instance(), scope.IdScope(id))
}

// LogicModifyJobSpec 更新任务条件
func LogicModifyJobSpec(id uint, spec string) error {
	r := regexp.MustCompile(`@(annually|yearly|monthly|weekly|daily|hourly|reboot)|(@every (\d+(m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	if !r.MatchString(spec) {
		return ErrCronExpression
	}

	response := &Response{}
	err := response.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}

	// 停止
	if response.Status == "running" {
		cron_server.CronInstance().Remove(response.EntryId)
		cron_server.CronInstance().Run()
	}

	data := map[string]interface{}{"spec": spec}
	err = database.Instance().Model(&Job{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}

	return nil
}

// LogicModifyStatus 更新状态
func LogicModifyStatus(id uint, status string) error {
	data := map[string]interface{}{"status": status}
	response := &Response{}
	err := response.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}

	if response.Status == status {
		return nil
	}

	// 停止
	if status == "stoped" {
		cron_server.CronInstance().Remove(response.EntryId)
		cron_server.CronInstance().Run()
	} else if status == "running" {
		entryId, err := cron_server.CronInstance().AddJob(response.Spec, BuiltinJobs.GetBuiltinJob(response.Name))
		if err != nil {
			return err
		}
		data["entry_id"] = entryId
		cron_server.CronInstance().Run()
	}

	err = database.Instance().Model(&Job{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}

	return nil
}

// syncJob
func syncJob(jobName, spec, desc string, job cron.Job) error {

	if spec == "" {
		spec = DefaultCronJobSpec
	}

	response := &Response{}
	err := response.First(database.Instance(), NameScope(jobName))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		entryId, err := cron_server.CronInstance().AddJob(spec, job)
		if err != nil {
			zap_server.ZAPLOG.Error(err.Error())
			return err
		}
		req := &Request{BaseJob: BaseJob{EntryId: entryId, Name: jobName, Spec: spec, Status: "running", Desc: desc}}
		_, err = LogicCreate(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateExecInfo
func UpdateExecInfo(jobName, message string) error {
	response := &Response{}
	err := response.First(database.Instance(), NameScope(jobName))
	if err != nil {
		return err
	}
	count := response.Count
	count++
	data := map[string]interface{}{"count": count, "message": message, "last_at": time.Now()}
	err = database.Instance().Model(&Job{}).Where("id = ?", response.Id).Updates(data).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}

	return nil
}
