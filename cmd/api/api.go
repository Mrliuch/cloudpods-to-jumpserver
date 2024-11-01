package api

import (
	"cloudpods-webhook/pkg/cloudpods"
	"cloudpods-webhook/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunApi() *gin.Engine {
	// 启动api服务
	Router := gin.Default()
	// 定义一个处理POST请求的路由
	Router.POST("/cloud", func(c *gin.Context) {
		var notice common.Notice
		// 绑定请求体中的JSON数据到User结构体
		if err := c.ShouldBindJSON(&notice); err != nil {
			// 如果绑定失败，返回400 Bad Request
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		n := &cloudpods.Notices{Notices: &notice}
		if err := n.DealEvent(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
				"code":    -1,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "succeed",
			"code":    0,
		})
	})
	return Router

}
