package startup

import (
	"gitee.com/cristiane/micro-mall-api/router"
	"gitee.com/cristiane/micro-mall-api/vars"
	"github.com/gin-gonic/gin"
)

func RegisterHttpRoute() *gin.Engine {
	return router.InitRouter()
}

func RegisterTasks() []vars.CronTask {
	return []vars.CronTask{}
}
