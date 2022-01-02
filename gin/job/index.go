package job

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin-rbac/gin/middleware"
)

func Group(app *gin.RouterGroup) {
	portRouter := app.Group("job", middleware.Auth(), middleware.CasbinHandler(), middleware.OperationRecord())
	{
		portRouter.GET("/list", All)
		portRouter.POST("/modifyStatus/:id", ModifyStatus)
		portRouter.POST("/modifyJobSpec/:id", ModifyJobSpec)
		portRouter.GET("/execJob/:id", ExecJob)
	}
}
