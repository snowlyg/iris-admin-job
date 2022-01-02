package job

import (
	"github.com/snowlyg/iris-admin/server/zap_server"
)

type YourJob struct {
	Name string
	Ip   string
}

func (j *YourJob) Run() {
	var message string
	err := yourJobRun()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		message = err.Error()
	}
	err = UpdateExecInfo(j.Name, j.Ip, message)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}
}

// yourJobRun
func yourJobRun() error {
	// do something here...
	return nil
}
