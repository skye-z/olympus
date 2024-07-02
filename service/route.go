package service

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/processor"
	"xorm.io/xorm"
)

func addPublicRoute(router *gin.Engine, engine *xorm.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/app"
		router.HandleContext(ctx)
	})
	// OAuth2登录路由
	as := NewAuthService(engine)
	if as != nil {
		router.GET("/oauth2/bind", as.Bind)
		router.GET("/oauth2/login", as.Login)
		router.GET("/oauth2/callback", as.Callback)
	}
	product := model.ProductModel{
		DB: engine,
	}
	version := model.VersionModel{
		DB: engine,
	}
	// 制品路由
	ps := NewProductService(product, version)
	router.GET("/api/product/stat", ps.Stat)
	router.GET("/api/product/search", ps.Search)
	router.GET("/api/product/number", ps.GetNumber)
	router.GET("/api/product/list", ps.GetList)
	router.GET("/api/product/info/:id", ps.GetInfo)
	router.GET("/api/product/version/:id", ps.GetVersionList)
	router.GET("/api/product/npm/:name", ps.GetNpmConfig)
	router.GET("/api/product/maven/:group/:name", ps.GetMavenConfig)
	router.GET("/api/product/maven/:group/:name/:version", ps.GetMavenConfig)
	router.GET("/api/product/go/:group/:name", ps.GetGoConfig)
	// Maven处理路由
	maven := processor.Maven{
		Product: product,
		Version: version,
	}
	router.HEAD("/maven/*param", maven.GetFile)
	router.GET("/maven/*param", maven.GetFile)
	// NPM处理路由
	npm := processor.Npm{
		Product: product,
		Version: version,
	}
	router.POST("/npm/*param", npm.GetFile)
	router.GET("/npm/*param", npm.GetFile)
	// Go处理路由
	goModule := processor.GoModule{
		Product: product,
		Version: version,
	}
	router.GET("/go/*param", goModule.GetFile)
	// Docker处理路由
	dockerModule := processor.DockerModule{
		Product: product,
		Version: version,
	}
	router.Any("/docker/*param", dockerModule.GetFile)
}

func addPrivateRoute(router *gin.Engine, engine *xorm.Engine) {
	// private := router.Group("").Use(AuthHandler())
	// {
	// private.GET("/api/device/info", ms.GetDeviceInfo)
	// private.GET("/api/system/use", ms.GetUse)
	// }
}
