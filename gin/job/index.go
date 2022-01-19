package job

import (
	"github.com/gin-gonic/gin"
)

func Group(app *gin.RouterGroup) {
	portRouter := app.Group("job")
	{
		portRouter.GET("/list", All)
		portRouter.POST("/modifyStatus/:id", ModifyStatus)
		portRouter.POST("/modifyJobSpec/:id", ModifyJobSpec)
		portRouter.GET("/execJob/:id", ExecJob)
	}
}
