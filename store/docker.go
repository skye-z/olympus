package store

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/util"
)

const (
	dockerRepository = "./repository/docker/"
)

type DockerStore struct {
	RemoteURL string
	Product   model.ProductModel
	Version   model.VersionModel
}

type RespStore struct {
	Code   int
	Header http.Header
	Data   []byte
}

// 获取文件
func (ds DockerStore) GetFile(path string, header http.Header) *RespStore {
	var directory string
	var fileName string
	var group string
	var name string
	var version string
	save := false

	if path == "v2" || path == "v2/" {
		directory = dockerRepository
		fileName = "v2.json"
	} else {
		if strings.Contains(path, "/sha256:") {
			var paths []string
			if strings.Contains(path, "/blobs/sha256:") {
				paths = strings.Split(path, "/blobs/sha256:")
			} else {
				paths = strings.Split(path, "/sha256:")
			}
			params := strings.Split(paths[0], "/")
			group = params[1]
			name = params[2]
			version = paths[1]
			directory = dockerRepository + group + "/" + name
			fileName = paths[1]
			save = true
		} else {
			params := strings.Split(path, "/")
			group = params[1]
			name = params[2]
			version = params[len(params)-1]
			directory = dockerRepository + group + "/" + name
			fileName = version + ".json"
			save = true
		}
	}

	if util.CheckExist(directory + "/" + fileName) {
		log.Println("[Store] docker from cache: " + directory + "/" + fileName)
		rs := &RespStore{}

		if util.CheckExist(directory + "/" + fileName + ".resp") {
			cache := util.ReadFile(directory + "/" + fileName + ".resp")
			json.Unmarshal(cache, &rs)
		}
		rs.Data = util.ReadFile(directory + "/" + fileName)

		return rs
	} else {
		log.Println("[Store] docker from online: " + path)
		resp := util.GetResp(ds.RemoteURL, path, header, false)
		rs := &RespStore{}
		if resp != nil {
			headerCache := make(map[string][]string)
			for k, v := range resp.Header {
				headerCache[k] = v
			}

			rs.Code = resp.StatusCode
			rs.Header = headerCache

			if resp.StatusCode != 200 {
				return rs
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("[Req] error:", err)
				return nil
			}
			util.SaveFile(directory, fileName, body)

			rsc, err := json.Marshal(rs)
			if err == nil {
				util.SaveFile(directory, fileName+".resp", rsc)
			}

			rs.Data = body

			if save && version != "" {
				go ds.saveData(group, name, version)
			}
		}
		return rs
	}
}

// 保存数据
func (ds DockerStore) saveData(group, name, number string) {
	// 查询制品信息是否存在
	product := ds.Product.Query(4, group, name)
	if product == nil || product.Id == 0 {
		// 制品不存在 创建制品信息
		product = &model.Product{
			Processor: 4,
			Group:     group,
			Name:      name,
			AddTime:   time.Now().Unix(),
		}
		ds.Product.Add(product)
	}
	if product.Id > 0 {
		// 查询制品版本是否存在
		version := ds.Version.Query(product.Id, number)
		if version == nil || version.Id == 0 {
			// 版本不存在 创建版本信息
			version = &model.Version{
				PId:     product.Id,
				Number:  number,
				AddTime: time.Now().Unix(),
			}
			ds.Version.Add(version)
			// 更新制品
			product.UpdateTime = time.Now().Unix()
			ds.Product.Edit(product)
		}
	}
}
