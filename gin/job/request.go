package job

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

type Request struct {
	BaseJob
}

type ModifyStatusRequest struct {
	Status string `json:"status" form:"status"`
}
type ModifyJobSpecRequest struct {
	Spec string `json:"spec" form:"spec"`
}
type HoldServiceJobRequest struct {
	Name string `json:"name" form:"name"`
}

func (req *Request) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return orm.ErrParamValidate
	}
	return nil
}

type ReqPaginate struct {
	orm.Paginate
	Status string `json:"status" form:"status"`
	Name   string `json:"name" form:"name"`
}
