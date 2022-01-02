package job

import (
	"errors"
	"regexp"
	"time"

	"github.com/snowlyg/iris-admin/server/cron_server"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

var ErrCronExpression = errors.New("错误的CRON表达式")
var DefaultCronJobSpec = "@every 5m"

// StartJob
func StartJob() {
	cron_server.CronInstance().Start()
}

func StopJob() {
	cron_server.CronInstance().Stop()
}

// LogicExecJob 执行任务
func LogicExecJob(id uint) error {
	// response := &Response{}
	// err := response.First(database.Instance(), scope.IdScope(id))
	// if err != nil {
	// 	return err
	// }

	// if response.Status == "running" {
	// 	return cron_server.ErrStartServer
	// }

	// BuiltinJob.RLock()
	// defer BuiltinJob.RUnlock()

	// jobServer := BuiltinJob.jobs[response.Name]

	// once := time.Now().Add(2 * time.Second)
	// onceSpec := fmt.Sprintf("%d %d %d %d %d %d", once.Second(), once.Minute(), once.Hour(), once.Day(), once.Month(), once.Weekday())
	// _, err = cron_server.CronInstance().AddJob(onceSpec, jobServer)
	// if err != nil {
	// 	return err
	// }
	// cron_server.CronInstance().Run()

	return nil
}

// LogicCreate 添加
func LogicCreate(req *Request) (uint, error) {
	return orm.Create(database.Instance(), &Job{BaseJob: req.BaseJob})
}

// LogicUpdate 更新
func LogicUpdate(id uint, req *Request) error {
	return orm.Update(database.Instance(), id, &Job{BaseJob: req.BaseJob})
}

// LogicModifyJobSpec 更新任务条件
func LogicModifyJobSpec(id uint, spec string) error {
	r := regexp.MustCompile(`@(annually|yearly|monthly|weekly|daily|hourly|reboot)|(@every (\d+(m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	if !r.MatchString(spec) {
		return ErrCronExpression
	}

	response := &Response{}
	err := orm.First(database.Instance(), response, scope.IdScope(id))
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
	err := orm.First(database.Instance(), response, scope.IdScope(id))
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
		// BuiltinJob.RLock()
		// defer BuiltinJob.RUnlock()
		// entryId, err := cron_server.CronInstance().AddJob(response.Spec, BuiltinJob.jobs[response.Name])
		// if err != nil {
		// 	return err
		// }
		// data["entry_id"] = entryId
		cron_server.CronInstance().Run()
	}

	err = database.Instance().Model(&Job{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return err
	}

	return nil
}

// SyncJob
func SyncJob(jobName, desc string) error {
	// BuiltinJob.RLock()
	// defer BuiltinJob.RUnlock()
	// response := &Response{}
	// err := response.First(database.Instance(), NameScope(jobName), IpScope(global.LocalIP()))
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	entryId, err := cron_server.CronInstance().AddJob(DefaultCronJobSpec, BuiltinJob.jobs[jobName])
	// 	if err != nil {
	// 		zap_server.ZAPLOG.Error(err.Error())
	// 		return err
	// 	}
	// 	zap_server.ZAPLOG.Debug("任务持久化....")
	// 	req := &Request{BaseJob: BaseJob{EntryId: entryId, Name: jobName, Ip: global.LocalIP(), Spec: DefaultCronJobSpec, Status: "running", Desc: desc}}
	// 	_, err = LogicCreate(req)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if err == nil {
	// 	entryId, err := cron_server.CronInstance().AddJob(response.Spec, BuiltinJob.jobs[jobName])
	// 	if err != nil {
	// 		zap_server.ZAPLOG.Error(err.Error())
	// 		return err
	// 	}
	// 	zap_server.ZAPLOG.Debug("任务更新....")
	// 	data := map[string]interface{}{"status": "running", "entry_id": entryId}
	// 	err = database.Instance().Model(&Job{}).Where("id = ?", response.Id).Updates(data).Error
	// 	if err != nil {
	// 		zap_server.ZAPLOG.Error(err.Error())
	// 		return err
	// 	}
	// }

	return nil
}

// UpdateExecInfo
func UpdateExecInfo(jobName, ip, message string) error {
	response := &Response{}
	err := response.First(database.Instance(), NameScope(jobName), IpScope(ip))
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
