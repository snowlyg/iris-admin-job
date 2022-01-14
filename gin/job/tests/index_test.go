package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin-job/gin/job"
	rbac "github.com/snowlyg/iris-admin-rbac/gin"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
)

func TestJobList(t *testing.T) {
	TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
	TestClient.Login(rbac.LoginUrl, nil)
	if TestClient == nil {
		return
	}
	list := []httptest.Responses{
		// {
		// 	{Key: "id", Value: 1},
		// 	{Key: "updatedAt", Value: "", Type: "notempty"},
		// 	{Key: "createdAt", Value: "", Type: "notempty"},
		// 	{Key: "entryId", Value: 1},
		// 	{Key: "name", Value: ""},
		// 	{Key: "spec", Value: "@every 5m"},
		// 	{Key: "status", Value: "running"},
		// 	{Key: "desc", Value: ""},
		// },
	}

	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: httptest.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "list", Value: list},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	httptest.RequestParams = map[string]interface{}{"page": 1, "pageSize": 10, "orderBy": "id", "sort": "asc"}
	TestClient.GET(fmt.Sprintf("%s/list", Uri), pageKeys, httptest.RequestParams)
}

func TestModifyStatus(t *testing.T) {
	TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
	TestClient.Login(rbac.LoginUrl, nil)
	if TestClient == nil {
		return
	}

	req := &job.Request{BaseJob: job.BaseJob{EntryId: 789, Name: "updateStatus", Spec: job.DefaultCronJobSpec, Status: "running", Desc: "@every 5m"}}
	id, err := job.LogicCreate(req)
	if err != nil {
		t.Error(err)
		return
	}

	data := map[string]interface{}{
		"status": "stoped",
	}

	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	TestClient.POST(fmt.Sprintf("%s/modifyStatus/%d", Uri, id), pageKeys, data)
}

func TestModifyJobSpec(t *testing.T) {
	TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
	TestClient.Login(rbac.LoginUrl, nil)
	if TestClient == nil {
		return
	}
	req := &job.Request{BaseJob: job.BaseJob{EntryId: 678, Name: "updateJobSpec", Spec: job.DefaultCronJobSpec, Status: "running", Desc: "@every 5m"}}
	id, err := job.LogicCreate(req)
	if err != nil {
		t.Error(err)
		return
	}
	data := map[string]interface{}{
		"spec": "@every 5h",
	}
	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	TestClient.POST(fmt.Sprintf("%s/modifyJobSpec/%d", Uri, id), pageKeys, data)
}

func TestExecSerJob(t *testing.T) {
	TestClient = httptest.Instance(t, str.Join("http://", web.CONFIG.System.Addr), TestServer.GetEngine())
	TestClient.Login(rbac.LoginUrl, nil)
	if TestClient == nil {
		return
	}
	req := &job.Request{BaseJob: job.BaseJob{EntryId: 78934, Name: "exec_job", Spec: job.DefaultCronJobSpec, Status: "stoped", Desc: "@every 5m"}}
	id, err := job.LogicCreate(req)
	if err != nil {
		t.Error(err)
		return
	}

	pageKeys := httptest.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	TestClient.GET(fmt.Sprintf("%s/execJob/%d", Uri, id), pageKeys)
}
