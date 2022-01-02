package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iris-admin-job/gin/job"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
)

// Party v1 模块
func Party(app *gin.RouterGroup) {
	v1 := app.Group("/api/v1")
	{
		job.Group(v1) // 任务
	}
}

var LoginResponse = httptest.Responses{
	{Key: "status", Value: http.StatusOK},
	{Key: "message", Value: "请求成功"},
	{Key: "data",
		Value: httptest.Responses{
			{Key: "accessToken", Value: "", Type: "notempty"},
		},
	},
}
var LogoutResponse = httptest.Responses{
	{Key: "status", Value: http.StatusOK},
	{Key: "message", Value: "请求成功"},
}

// 加载模块
var PartyFunc = func(wi *web_gin.WebServer) {
	// 初始化驱动
	err := multi.InitDriver(&multi.Config{DriverType: "jwt", HmacSecret: nil})
	if err != nil {
		zap_server.ZAPLOG.Panic("err")
	}
	Party(wi.GetRouterGroup("/"))
}

//  填充数据
var SeedFunc = func(wi *web_gin.WebServer, mc *migration.MigrationCmd) {
	mc.AddMigration(job.GetMigration())
}
