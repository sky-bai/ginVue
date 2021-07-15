package routers

import (
	v1 "ginVue/api/v1"
	"ginVue/middleware"
	"ginVue/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Log())
	r.Use(middleware.Cors())
	//r.LoadHTMLGlob("static/bai/index.html")
	//r.Static("bai/static","static/bai/static")
	//r.StaticFile("bai/favicon.ico","static/bai/favicon.ico")
	//r.GET("bai", func(c *gin.Context) {
	//	c.HTML(200,"index.html",nil)
	//})
	gin.SetMode(utils.AppMode)
	auth := r.Group("api/v1")
	{
		// 用户模块的路由接口
		auth.GET("users", v1.GetUsers)
		auth.POST("user/add", v1.AddUser)
		auth.POST("login", v1.Login)
		auth.GET("user/:id", v1.GetUserInfo)

		// 分类模块的路由
		auth.GET("category", v1.GetCategory)
		// 获取单个文章
		auth.GET("article/info/:id", v1.GetArticle)
		// 获取文章列表
		auth.GET("articles", v1.GetArticles)
		// 查询分类下的所有文章
		auth.GET("article/catelist/:id", v1.GetCategoryArticles)
		// 查询单个分类
		auth.GET("category/:id", v1.GetCateInfo)

	}
	router := r.Group("api/v1")
	router.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		router.DELETE("user/:id", v1.DeleteUser)
		router.PUT("user/:id", v1.EditUser)

		// 分类模块的路由
		router.POST("category/add", v1.AddCategory)
		router.DELETE("category/:id", v1.DeleteCategory)
		router.PUT("category/:id", v1.EditCategory)

		// 文章模块的路由
		// 添加文章
		router.POST("article/add", v1.AddArticle)
		// 删除文章
		router.DELETE("article/:id", v1.DeleteArticle)
		// 编辑文章
		router.PUT("article/:id", v1.EditArticle)

		// 上传文件
		router.POST("upload", v1.UpLoad)
	}

	r.Run(":8081")
}
