package processor

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
)

type Maven struct {
	Product model.ProductModel
	Version model.VersionModel
}

func (m Maven) GetFile(ctx *gin.Context) {
	param := ctx.Param("param")

	ms := store.MavenStore{
		RemoteURL: "https://repo1.maven.org/maven2/",
		Product:   m.Product,
		Version:   m.Version,
	}

	ext := filepath.Ext(param)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" && strings.Contains(param, "/-/") {
		mimeType = "application/octet-stream"
	} else if mimeType == "" {
		mimeType = "application/xml"
	}

	ctx.Data(200, mimeType, ms.GetFile(param[1:]))
	ctx.Abort()
}
