package tests

import (
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/snowlyg/httptest"
	job_gin "github.com/snowlyg/iris-admin-job/gin"
	"github.com/snowlyg/iris-admin-job/gin/job"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/common"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
)

var TestServer *web_gin.WebServer
var TestClient *httptest.Client

func TestMain(m *testing.M) {
	web.CONFIG.System.Level = "test"
	var uuid string
	uuid, TestServer = common.BeforeTestMainGin(job_gin.PartyFunc, job_gin.SeedFunc)

	go job.StartJob() // 服务监控

	time.Sleep(5 * time.Second)
	code := m.Run()

	common.AfterTestMain(uuid, true)
	job.StopJob()

	web.Remove()
	os.Exit(code)
}
