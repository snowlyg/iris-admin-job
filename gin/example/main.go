package main

import (
	job_gin "github.com/snowlyg/iris-admin-job/gin"
	"github.com/snowlyg/iris-admin-job/gin/job"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func main() {
	wi := web_gin.Init()
	v1 := wi.GetRouterGroup("/api/v1")
	{
		job_gin.Party(v1)
	}
	go func() {
		job.BuiltinJobs.AddBuiltinJob("yourJobRun", "@every 1m", "yourJobRun", &YourJob{})
		job.StartJob()
	}()
	web.Start(wi)
}

type YourJob struct {
	Name string
	// ....
}

func (j *YourJob) Run() {
	var message string
	err := yourJobRun()
	if err != nil {
		message = err.Error()
	}
	err = job.UpdateExecInfo(j.Name, message)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}
}

// yourJobRun
func yourJobRun() error {
	// do something here...
	return nil
}
