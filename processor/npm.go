package processor

import (
	"compress/gzip"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
	"github.com/skye-z/olympus/util"
)

const NpmRemoteURL = "https://registry.npmjs.org/"

type Npm struct {
	Product model.ProductModel
	Version model.VersionModel
}

func (n Npm) GetFile(ctx *gin.Context) {
	method := ctx.Request.Method
	param := ctx.Param("param")

	ms := store.NpmStore{
		RemoteURL: NpmRemoteURL,
		Product:   n.Product,
		Version:   n.Version,
	}

	if method == http.MethodPost {
		gzipReader, err := gzip.NewReader(ctx.Request.Body)
		if err != nil {
			util.ReturnMessage(ctx, false, "数据格式错误")
		}
		ctx.Data(200, "application/json", ms.GetSecurityFile(param[1:], ctx.ContentType(), gzipReader))
		ctx.Abort()
	} else {
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
}

func (n Npm) GetConfig(name string) []byte {
	ms := store.NpmStore{
		RemoteURL: NpmRemoteURL,
		Product:   n.Product,
		Version:   n.Version,
	}
	return ms.GetFile(name)
}
