/**
<h1 align="center">IrisAdminJob</h1>

[IrisAdminJob](https://www.github.com/snowlyg/iris-admin-job) 项目为一个任务管模块插件,可以为 [IrisAdmin](https://www.github.com/snowlyg/iris-admin) 项目快速集成任务管理API.

##### 下载

```sh
  go get -u github.com/snowlyg/iris-admin-job@latest
```

#####
- 简单使用(only for gin)
```go
package main

import (
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
  job "github.com/snowlyg/iris-admin-job/gin"
)

func main() {
  wi := web_gin.Init()
  v1 := wi.GetRouterGroup("/api/v1")
	{
		job.Party(v1)
	}
  go func(){
    	job.BuiltinJobs.AddBuiltinJob("yourJobRun", "@every 1m", "yourJobRun", &YourJob{})
	    job.StartJob()
  }()
	web.Start(wi)
}

```

#### 添加任务
```go
package job

import (
	"github.com/snowlyg/iris-admin/server/zap_server"
  	"github.com/snowlyg/iris-admin-job/gin/job"
)

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

```

##### 启动任务

```go
	job.BuiltinJobs.AddBuiltinJob("yourJobRun", "@every 1m", "yourJobRun", &YourJob{})
	job.StartJob()
```
*/

package doc
