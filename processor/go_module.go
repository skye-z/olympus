package processor

import (
	"mime"
	"path/filepath"

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

	ctx.Data(200, mimeType, gs.GetFile(param[1:], mimeType))
	ctx.Abort()
}

func (g GoModule) GetConfig(group, name string) []byte {
	gs := store.GoStore{
		RemoteURL: GoRemoteURL,
		Product:   g.Product,
		Version:   g.Version,
	}
	return gs.GetFile(group+"/"+name+"/@v/list", "application/json")
}
