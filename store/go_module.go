package store

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/util"
)

const (
	goRepository = "./repository/go/"
)

type GoStore struct {
	RemoteURL string
	Product   model.ProductModel
	Version   model.VersionModel
}

// 获取文件
func (gs GoStore) GetFile(path, mimeType string) []byte {
	params := strings.Split(path, "/")
	group := params[0]
	if len(params) > 4 {
		group = group + "/" + params[1]
	}

	if strings.Contains(path, "/@v/") {
		name := params[len(params)-3]
		version := params[len(params)-1]
		if version[0:1] == "v" {
			version = version[1:]
		}
		version = strings.ReplaceAll(version, ".zip", "")

		extend := ""
		if mimeType == "application/zip" {
			extend = ".zip"
		}

		filePath := goRepository + group + "/" + name + "/" + version + extend
		if util.CheckExist(filePath) {
			log.Println("[Store] go from cache: " + filePath)
			return util.ReadFile(filePath)
		} else {
			log.Println("[Store] go from online: " + path)
			content := gs.getRemoteData(path)
			if content != nil {
				util.SaveFile(goRepository+group+"/"+name, version+extend, content)
			}
			go gs.saveData(group, name, version)
			return content
		}
	} else {
		name := params[len(params)-2]
		filePath := goRepository + group + "/" + name + "/lastest.json"
		if util.CheckExist(filePath) {
			log.Println("[Store] go from cache: " + filePath)
			return util.ReadFile(filePath)
		} else {
			log.Println("[Store] go from online: " + path)
			content := gs.getRemoteData(path)
			if content != nil {
				util.SaveFile(goRepository+group+"/"+name, "lastest.json", content)
			}
			return content
		}
	}
}

// 获取远程数据
func (gs GoStore) getRemoteData(path string) []byte {
	resp, err := http.Get(gs.RemoteURL + path)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return body
}

// 保存数据
func (gs GoStore) saveData(group, name, number string) {
	// 查询制品信息是否存在
	product := gs.Product.Query(3, group, name)
	if product == nil || product.Id == 0 {
		// 制品不存在 创建制品信息
		product = &model.Product{
			Processor: 3,
			Group:     group,
			Name:      name,
			AddTime:   time.Now().Unix(),
		}
		gs.Product.Add(product)
	}
	if product.Id > 0 {
		// 查询制品版本是否存在
		version := gs.Version.Query(product.Id, number)
		if version == nil || version.Id == 0 {
			// 版本不存在 创建版本信息
			version = &model.Version{
				PId:     product.Id,
				Number:  number,
				AddTime: time.Now().Unix(),
			}
			gs.Version.Add(version)
			// 更新制品
			product.UpdateTime = time.Now().Unix()
			gs.Product.Edit(product)
		}
	}
}
