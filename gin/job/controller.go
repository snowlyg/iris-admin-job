package job

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/web/web_gin/request"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	"gorm.io/gorm"
)

// All 列表
func All(ctx *gin.Context) {
	req := &ReqPaginate{}
	if errs := ctx.ShouldBind(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	items := &PageResponse{}
	var scopes []func(db *gorm.DB) *gorm.DB
	if req.Name != "" {
		scopes = append(scopes, NameScope(req.Name))
	}
	if req.Status != "" {
		scopes = append(scopes, StatusScope(req.Status))
	}
	total, err := orm.Pagination(database.Instance(), items, req.PaginateScope(), scopes...)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.OkWithData(response.PageResult{
		List:     items.Item,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, ctx)
}

// ModifyStatus 更新状态
func ModifyStatus(ctx *gin.Context) {
	var req request.GetById
	if errs := ctx.ShouldBindUri(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	mStatus := &ModifyStatusRequest{}
	if errs := ctx.ShouldBindJSON(mStatus); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := LogicModifyStatus(req.Id, mStatus.Status)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}

// ModifyJobSpec 更新任务条件
func ModifyJobSpec(ctx *gin.Context) {
	var req request.GetById
	if errs := ctx.ShouldBindUri(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	jobSpec := &ModifyJobSpecRequest{}
	if errs := ctx.ShouldBindJSON(jobSpec); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := LogicModifyJobSpec(req.Id, jobSpec.Spec)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}

// HoldService 服务保持
func HoldService(ctx *gin.Context) {
	holdSerJob := &HoldServiceJob{}
	if errs := ctx.ShouldBindJSON(holdSerJob); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := LogicHoldService(holdSerJob)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	} else {
		response.Ok(ctx)
	}
}

// ExecJob 执行任务
func ExecJob(ctx *gin.Context) {
	var req request.GetById
	if errs := ctx.ShouldBindUri(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := LogicExecJob(req.Id)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	} else {
		response.Ok(ctx)
	}
}
