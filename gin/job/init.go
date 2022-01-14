package job

import "github.com/snowlyg/iris-admin/server/database"

// 每次启动清理任务
func Init() {
	database.Instance().Unscoped().Delete(&Job{})
}
