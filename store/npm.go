package store

import (
	"bytes"
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
	npmRepository = "./repository/npm/"
)

type NpmStore struct {
	RemoteURL string
	Product   model.ProductModel
	Version   model.VersionModel
}

// 获取文件
func (ns NpmStore) GetFile(path string) []byte {
	fileName := ""
	directory := ""
	index := strings.LastIndex(path, "/")
	if index == -1 {
		fileName = path + ".json"
		directory = path
	} else {
		fileName = path[strings.LastIndex(path, "/")+1:]
		directory = strings.ReplaceAll(path[:strings.LastIndex(path, "/")], "/-", "")
	}
	// 无后缀文件统一标为JSON
	if !strings.Contains(fileName, ".") {
		fileName = fileName + ".json"
	}

	if util.CheckExist(npmRepository + directory + "/" + fileName) {
		log.Println("[Store] npm from cache: " + npmRepository + directory + "/" + fileName)
		return util.ReadFile(npmRepository + directory + "/" + fileName)
	} else {
		log.Println("[Store] npm from online: " + path)
		content := ns.getRemoteData(path)
		if content == nil {
			return nil
		}
		if index == -1 {
			result := bytes.Replace(content, []byte("https://registry.npmjs.org/"), []byte("http://localhost:27680/npm/"), -1)
			util.SaveFile(npmRepository+directory, fileName, result)
			return result
		} else {
			util.SaveFile(npmRepository+directory, fileName, content)
			go ns.saveData(path)
			return content
		}
	}
}

// 获取远程数据
func (ns NpmStore) GetSecurityFile(path, contentType string, body io.ReadCloser) []byte {
	resp, err := http.Post(ns.RemoteURL+path, contentType, body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return res
}

// 获取远程数据
func (ns NpmStore) getRemoteData(path string) []byte {
	resp, err := http.Get(ns.RemoteURL + path)
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
func (ns NpmStore) saveData(path string) {
	params := strings.Split(path, "/")
	name := ""
	number := ""
	length := len(params)
	if length == 2 {
		return
	} else if length > 2 {
		name = params[0]
		cache := params[length-1]
		cache = strings.Replace(cache, name+"-", "", -1)
		number = cache[0:strings.LastIndex(cache, ".")]
	}
	// 查询制品信息是否存在
	product := ns.Product.Query(2, "", name)
	if product == nil || product.Id == 0 {
		// 制品不存在 创建制品信息
		product = &model.Product{
			Processor: 2,
			Group:     "",
			Name:      name,
			AddTime:   time.Now().Unix(),
		}
		ns.Product.Add(product)
	}
	if product.Id > 0 {
		// 查询制品版本是否存在
		version := ns.Version.Query(product.Id, number)
		if version == nil || version.Id == 0 {
			// 版本不存在 创建版本信息
			version = &model.Version{
				PId:     product.Id,
				Number:  number,
				AddTime: time.Now().Unix(),
			}
			ns.Version.Add(version)
			// 更新制品
			product.UpdateTime = time.Now().Unix()
			ns.Product.Edit(product)
		}
	}
}
