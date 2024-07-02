package model

import (
	"log"

	"xorm.io/xorm"
)

type Product struct {
	Id int64 `json:"id"`
	// 制品处理器(1Maven、2NPM、3Go、4Docker)
	Processor int16 `json:"processor"`
	// 制品组
	Group string `json:"group"`
	// 制品名
	Name string `json:"name"`
	// 添加时间
	AddTime int64 `json:"addTime"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
}

type ProductWithVersion struct {
	Product       `xorm:"extends"`
	LatestVersion string `json:"lastVersion" xorm:"'number'"`
}

type ProductModel struct {
	DB *xorm.Engine
}

func (model ProductModel) Add(product *Product) bool {
	_, err := model.DB.Insert(product)
	return err == nil
}

func (model ProductModel) Edit(product *Product) bool {
	if product.Id == 0 {
		return false
	}
	_, err := model.DB.ID(product.Id).Update(product)
	return err == nil
}

func (model ProductModel) Del(product *Product) bool {
	if product.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(product)
	return err == nil
}

func (model ProductModel) GetItem(id int64) *Product {
	product := &Product{
		Id: id,
	}
	has, _ := model.DB.Get(product)
	if has {
		return product
	}
	return nil
}

func (model ProductModel) Query(processor int16, group, name string) *Product {
	product := &Product{
		Processor: processor,
		Group:     group,
		Name:      name,
	}
	has, _ := model.DB.Get(product)
	if has {
		return product
	}
	return nil
}

func (model ProductModel) GetList(processor int, group, name string, page int, num int) ([]ProductWithVersion, error) {
	var list []ProductWithVersion
	var cache *xorm.Session
	cache = model.DB.Table("product").
		Join("LEFT", "version", "product.id = version.p_id").
		GroupBy("product.id").
		Select("product.*, MAX(version.number) as `number`")
	if len(name) > 0 {
		cache = cache.Where("product.name LIKE ?", name)
	}
	if processor > 0 {
		cache = cache.Where("product.processor = ?", processor)
	}
	if len(group) > 0 {
		cache = cache.Where("product.group = ?", group)
	}
	err := cache.Desc("product.add_time").Limit(page*num, (page-1)*num).Find(&list)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return list, nil
}

func (model ProductModel) GetNumber(processor int, group, name string) (int64, error) {
	var product Product
	var cache *xorm.Session
	if len(name) > 0 {
		cache = model.DB.Where("name LIKE ?", name)
	} else {
		cache = model.DB.Where("")
	}
	if processor > 0 {
		cache = cache.Where("processor = ?", processor)
	}
	if len(group) > 0 {
		cache = cache.Where("group = ?", group)
	}
	number, err := cache.Count(product)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (model ProductModel) Stat() []map[string]interface{} {
	data, err := model.DB.QueryInterface("SELECT processor, COUNT( 1 ) AS `number` FROM product GROUP BY processor")
	if err == nil {
		return data
	}
	return nil
}
