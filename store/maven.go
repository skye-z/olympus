package store

import (
	"log"
	"strings"
	"time"

	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/util"
)

const (
	mavenRepository = "./repository/maven/"
)

type MavenStore struct {
	RemoteURL string
	Product   model.ProductModel
	Version   model.VersionModel
}

// 获取文件
func (ms MavenStore) GetFile(path string) []byte {
	fileName := path[strings.LastIndex(path, "/")+1:]
	directory := path[:strings.LastIndex(path, "/")]
	if util.CheckExist(mavenRepository + path) {
		log.Println("[Store] maven from cache: " + npmRepository + directory + "/" + fileName)
		return util.ReadFile(mavenRepository + path)
	} else {
		log.Println("[Store] maven from online: " + path)
		content := util.GetData(ms.RemoteURL, path, true)
		if content != nil {
			util.SaveFile(mavenRepository+directory, fileName, content)
		}
		go ms.saveData(directory, fileName)
		return content
	}
}

// 保存数据
func (ms MavenStore) saveData(path, fileName string) {
	params := strings.Split(path, "/")
	if len(params) < 2 {
		return
	}
	// 从地址解析数据
	name := params[len(params)-2]
	number := params[len(params)-1]
	file := fileName[:strings.LastIndex(fileName, ".")]
	if name+"-"+number != file {
		return
	}
	group := strings.Join(params[:len(params)-2], ".")
	// 查询制品信息是否存在
	product := ms.Product.Query(1, group, name)
	if product == nil || product.Id == 0 {
		// 制品不存在 创建制品信息
		product = &model.Product{
			Processor: 1,
			Group:     group,
			Name:      name,
			AddTime:   time.Now().Unix(),
		}
		ms.Product.Add(product)
	}
	if product.Id > 0 {
		// 查询制品版本是否存在
		version := ms.Version.Query(product.Id, number)
		if version == nil || version.Id == 0 {
			// 版本不存在 创建版本信息
			version = &model.Version{
				PId:     product.Id,
				Number:  number,
				AddTime: time.Now().Unix(),
			}
			ms.Version.Add(version)
			// 更新制品
			product.UpdateTime = time.Now().Unix()
			ms.Product.Edit(product)
		}
	}
}
