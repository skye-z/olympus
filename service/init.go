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
	"github.com/skye-z/olympus/util"
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
		log.Printf("[NoR] URL: %s\n", c.Request.Method)
		log.Printf("[NoR] URL: %s\n", c.Request.URL.Path)
		c.Next()
	})
}

func RunRouter(router *gin.Engine, engine *xorm.Engine) {
	port := getPort()
	log.Println("[Core] service started, port is", port)
	// 启动服务
	go func() {
		cert := util.GetString("ssl.cert")
		if cert == "" {
			if err := router.Run(":" + port); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		} else {
			key := util.GetString("ssl.key")
			if err := router.RunTLS(":"+port, cert, key); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		}
	}()
	// 等待中断信号以优雅关闭服务器
	waitForInterrupt(engine)
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
