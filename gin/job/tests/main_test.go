package tests

import (
	_ "embed"
	"fmt"
	"os"
	"testing"
	"time"

	gin_job "github.com/iris-admin-job/gin"
	"github.com/iris-admin-job/gin/job"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/server/web"
	web_tests "github.com/snowlyg/iris-admin/server/web/tests"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
)

var TestServer *web_gin.WebServer
var TestClient *httptest.Client

func TestMain(m *testing.M) {
	web.CONFIG.System.Level = "test"
	var uuid string
	uuid, TestServer = web_tests.BeforeTestMainGin(gin_job.PartyFunc, gin_job.SeedFunc)

	go job.StartJob() // 服务监控

	time.Sleep(5 * time.Second)
	code := m.Run()

	web_tests.AfterTestMain(uuid, true)
	job.StopJob()

	err := job.Remove()
	if err != nil {
		fmt.Println(err)
	}
	err = web.Remove()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}
