package handler

import (
	"androidServer/api"
	"androidServer/app"
	"androidServer/app/log"
	"github.com/gin-gonic/gin"
)

func HttpRun(addr string) error {
	r := gin.New()

	r.GET("/HeartBeat", api.HeartBeatHandler().Process)

	user := r.Group(app.Conf.PathPrefix + "/users")
	{
		////枚举用户
		//user.GET("", func(c *gin.Context) {
		//	api.ListUsersHandler().Process(c)
		//})
		//导出枚举用户
		user.POST("/addUser", func(c *gin.Context) {
			api.AddUserHandler().Process(c)
		})
		//注册
		user.POST("/register", func(c *gin.Context) {
			api.RegisterUserHandler().Process(c)
		})
		//登陆
		user.POST("/login", func(c *gin.Context) {
			api.LoginUserHandler().Process(c)
		})

	}
	word := r.Group(app.Conf.PathPrefix + "/words")
	{
		word.POST("/addWord", func(c *gin.Context) {
			api.AddWordHandler().Process(c)
		})
	}

	files := r.Group(app.Conf.PathPrefix + "/files")
	{
		// 上传文件
		files.POST("", func(c *gin.Context) {
			api.UploadFileHandler().Process(c)
		})
	}
	log.Infof("HTTP server is running on %s", addr)
	_ = r.Run(addr)
	log.Errorf("HTTP server will be down %s", addr)
	return nil
}
