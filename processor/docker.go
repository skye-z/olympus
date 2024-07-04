package processor

import (
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
	header := ctx.Request.Header

	ds := store.DockerStore{
		RemoteURL: DockerRemoteURL,
		Product:   d.Product,
		Version:   d.Version,
	}

	resp := ds.GetFile(param[1:], header)

	ctx.Status(resp.Code)
	for k, v := range resp.Header {
		ctx.Header(k, v[0])
	}
	ctx.Writer.Write(resp.Data)
	ctx.Abort()
}

func (d DockerModule) GetConfig(ctx *gin.Context, group, name, version string) []byte {
	ds := store.DockerStore{
		RemoteURL: DockerRemoteURL,
		Product:   d.Product,
		Version:   d.Version,
	}
	header := ctx.Request.Header
	var data *store.RespStore
	if version == "" {
		data = ds.GetFile("v2/"+group+"/"+name+"/manifests/latest", header)
	} else {
		data = ds.GetFile("v2/"+group+"/"+name+"/manifests/"+version, header)
	}
	if data.Data == nil {
		return nil
	}
	return data.Data
}
