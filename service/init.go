package service

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/processor"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func InitDB() *xorm.Engine {
	log.Println("[Data] load engine")
	engine, err := xorm.NewEngine("sqlite", "./local.store")
	if err != nil {
		panic(err)
	}
	return engine
}

func InitDBTable(engine *xorm.Engine) {
	log.Println("[Data] load data")
	err := engine.Sync2(new(model.User))
	if err != nil {
		panic(err)
	}
	err = engine.Sync2(new(model.Product))
	if err != nil {
		panic(err)
	}
	err = engine.Sync2(new(model.Version))
	if err != nil {
		panic(err)
	}
}

func InitRouter(page embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	log.Println("[Core] load page")
	pageFile, _ := fs.Sub(page, "page/dist")
	router.StaticFS("/app", http.FS(pageFile))
	return router
}

func InitRoute(router *gin.Engine, engine *xorm.Engine) {
	// 公共路由
	addPublicRoute(router, engine)

	// 私有路由
	addPrivateRoute(router, engine)

	router.NoRoute(func(c *gin.Context) {
		// 打印请求地址
		log.Printf("[NR] URL: %s\n", c.Request.URL.Path)
		c.Next()
	})
}

func RunRouter(router *gin.Engine, engine *xorm.Engine) {
	port := getPort()
	log.Println("[Core] service started, port is", port)
	// 启动服务
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	// 等待中断信号以优雅关闭服务器
	waitForInterrupt(engine)
}

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
}

func addPrivateRoute(router *gin.Engine, engine *xorm.Engine) {
	// private := router.Group("").Use(AuthHandler())
	// {
	// private.GET("/api/device/info", ms.GetDeviceInfo)
	// private.GET("/api/system/use", ms.GetUse)
	// }
}

// 等待关闭
func waitForInterrupt(engine *xorm.Engine) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Core] shutting down server")

	defer engine.Close()

	log.Println("[Core] server stopped")
}

// 获取端口号配置
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "27680"
	}
	return port
}
