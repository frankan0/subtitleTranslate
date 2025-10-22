package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/frank0/subtitleTranslate/api/handlers"
	"github.com/frank0/subtitleTranslate/api/middleware"
	"github.com/frank0/subtitleTranslate/internal/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置API路由
func SetupRouter() *gin.Engine {
	// 创建默认的gin路由器
	router := gin.Default()

	// 添加请求超时中间件 - 90秒超时
	router.Use(timeoutMiddleware(90 * time.Second))

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

	// 提供前端静态文件（SPA 支持）
	fsys := static.GetFileSystem()
	router.Use(func(c *gin.Context) {
		// 如果是 API 前缀，跳过静态文件处理
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.Next()
			return
		}
		// 尝试读取文件
		file, err := fsys.Open(c.Request.URL.Path)
		if err == nil {
			file.Close()
			// 文件存在，交给静态文件服务器
			http.FileServer(fsys).ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		// 文件不存在，回退到 index.html
		index, err := fsys.Open("index.html")
		if err != nil {
			c.String(404, "Not Found")
			return
		}
		index.Close()
		c.Request.URL.Path = "index.html"
		http.FileServer(fsys).ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})

	return router
}

// timeoutMiddleware 请求超时中间件
func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求的context
		c.Request = c.Request.WithContext(ctx)

		// 继续处理请求
		c.Next()

		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"error":   "请求处理超时，请稍后重试",
			})
		}
	}
}
