package processor

import (
	"mime"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
)

const DockerRemoteURL = "https://registry.hub.docker.com/"

type DockerModule struct {
	Product model.ProductModel
	Version model.VersionModel
}

func (d DockerModule) GetFile(ctx *gin.Context) {
	param := ctx.Param("param")
	body := ctx.Request.Body
	header := ctx.Request.Header
	method := ctx.Request.Method

	ds := store.DockerStore{
		RemoteURL: DockerRemoteURL,
		Product:   d.Product,
		Version:   d.Version,
	}

	ext := filepath.Ext(param)
	mimeType := mime.TypeByExtension(ext)

	resp := ds.GetFile(param[1:], mimeType, method, body, header)

	ctx.Status(resp.Code)
	for k, v := range resp.Header {
		ctx.Header(k, v[0])
	}
	ctx.Writer.Write(resp.Data)
	ctx.Abort()
}
