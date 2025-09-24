package routes

import (
	"net/http"

	"github.com/frank0/subtitleTranslate/api/handlers"
	"github.com/frank0/subtitleTranslate/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置API路由
func SetupRouter() *gin.Engine {
	// 创建默认的gin路由器
	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 添加自定义中间件
	router.Use(middleware.Logger())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API路由组
	api := router.Group("/api")
	{
		// 字幕翻译路由
		subtitle := api.Group("/subtitle")
		{
			// 翻译字幕文件
			subtitle.POST("/translate", handlers.TranslateSubtitle)
		}
	}

	return router
}
