package processor

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/store"
)

const MavenRemoteURL = "https://repo1.maven.org/maven2/"

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

func (m Maven) GetConfig(group, name, version string) []byte {
	ms := store.MavenStore{
		RemoteURL: "https://repo1.maven.org/maven2/",
		Product:   m.Product,
		Version:   m.Version,
	}
	var data []byte
	url := strings.ReplaceAll(group, ".", "/")
	if version == "" {
		url = url + "/" + name
		data = ms.GetFile(url + "/maven-metadata.xml")
		if data == nil {
			url := strings.ReplaceAll(url, "-", "_")
			data = ms.GetFile(url + "/maven-metadata.xml")
		}
	} else {
		url = url + "/" + name + "/" + version + "/" + name
		data = ms.GetFile(url + "-" + version + ".pom")
		if data == nil {
			url := strings.ReplaceAll(url, "-", "_")
			data = ms.GetFile(url + "-" + version + ".pom")
		}
	}
	return data
}
