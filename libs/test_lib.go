package libs

import "github.com/gin-gonic/gin"

func InitGinForTesting() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}
