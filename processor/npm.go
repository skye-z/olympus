package processor

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
)

type Npm struct {
	Product model.ProductModel
	Version model.VersionModel
}

func (n Npm) GetFile(ctx *gin.Context) {
	param := ctx.Param("param")

	ms := store.NpmStore{
		RemoteURL: "https://registry.npmjs.org/",
		Product:   n.Product,
		Version:   n.Version,
	}

	ext := filepath.Ext(param)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" && strings.Contains(param, "/-/") {
		mimeType = "application/octet-stream"
	} else if mimeType == "" {
		mimeType = "application/json"
	}

	ctx.Data(200, mimeType, ms.GetFile(param[1:]))
	ctx.Abort()
}
