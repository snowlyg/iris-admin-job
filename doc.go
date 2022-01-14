/**
<h1 align="center">IrisAdminJob</h1>

[IrisAdminJob](https://www.github.com/snowlyg/iris-admin-job) 项目为一个任务管模块插件,可以为 [IrisAdmin](https://www.github.com/snowlyg/iris-admin) 项目快速集成任务管理API.

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
	web.Start(wi)
}

```
*/

package doc
