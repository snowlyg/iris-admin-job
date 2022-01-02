package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin-job/gin/job"
)

// Party job模块
func Party(app *gin.RouterGroup) {
	job.Group(app)
}
