package main

import (
	"embed"

	"github.com/skye-z/olympus/service"
	"github.com/skye-z/olympus/util"
)

//go:embed page/dist/*
var page embed.FS

func main() {
	util.InitConfig()
	engine := service.InitDB()
	go service.InitDBTable(engine)
	router := service.InitRouter(page)
	service.InitRoute(router, engine)
	service.RunRouter(router, engine)
}
