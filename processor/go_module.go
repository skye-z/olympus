package processor

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
)

const GoRemoteURL = "https://goproxy.cn/"

type GoModule struct {
	Product model.ProductModel
	Version model.VersionModel
}

func (g GoModule) GetFile(ctx *gin.Context) {
	param := ctx.Param("param")

	gs := store.GoStore{
		RemoteURL: GoRemoteURL,
		Product:   g.Product,
		Version:   g.Version,
	}

	ext := filepath.Ext(param)
	mimeType := mime.TypeByExtension(ext)

	ctx.Data(200, mimeType, gs.GetFile(param[1:]))
	ctx.Abort()
}

func (g GoModule) GetConfig(group, name, version string) []byte {
	gs := store.GoStore{
		RemoteURL: GoRemoteURL,
		Product:   g.Product,
		Version:   g.Version,
	}
	var data []byte
	url := strings.ReplaceAll(group, ".", "/")
	if version == "" {
		url = url + "/" + name
		data = gs.GetFile(url + "/maven-metadata.xml")
		if data == nil {
			url := strings.ReplaceAll(url, "-", "_")
			data = gs.GetFile(url + "/maven-metadata.xml")
		}
	} else {
		url = url + "/" + name + "/" + version + "/" + name
		data = gs.GetFile(url + "-" + version + ".pom")
		if data == nil {
			url := strings.ReplaceAll(url, "-", "_")
			data = gs.GetFile(url + "-" + version + ".pom")
		}
	}
	return data
}
