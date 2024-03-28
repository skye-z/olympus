package model

import (
	"xorm.io/xorm"
)

type Product struct {
	Id int64 `json:"id"`
	// 制品处理器(1Maven、2NPM)
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

type ProductModel struct {
	DB *xorm.Engine
}

func (model ProductModel) AddProduct(product *Product) bool {
	_, err := model.DB.Insert(product)
	return err == nil
}

func (model ProductModel) EditProduct(product *Product) bool {
	if product.Id == 0 {
		return false
	}
	_, err := model.DB.ID(product.Id).Update(product)
	return err == nil
}

func (model ProductModel) DelProduct(product *Product) bool {
	if product.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(product)
	return err == nil
}

func (model ProductModel) GetProduct(processor int16, group, name string) *Product {
	product := &Product{
		Processor: processor,
		Group:     group,
		Name:      name,
	}
	has, _ := model.DB.Get(product)
	if !has {
		return product
	}
	return nil
}

func (model ProductModel) GetList(processor int, group, name string, page int, num int) ([]Product, error) {
	var list []Product
	var cache *xorm.Session
	if processor > 0 {
		cache = model.DB.Where("processor = ?", processor)
	}
	if len(group) > 0 {
		cache = model.DB.Where("group = ?", group)
	}
	if len(name) > 0 {
		cache = model.DB.Where("name LIKE ?", name)
	}
	err := cache.Limit(page*num, (page-1)*num).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (model ProductModel) GetNumber(processor int, group, name string) (int64, error) {
	var product Product
	var cache *xorm.Session
	if processor > 0 {
		cache = model.DB.Where("processor = ?", processor)
	}
	if len(group) > 0 {
		cache = model.DB.Where("group = ?", group)
	}
	if len(name) > 0 {
		cache = model.DB.Where("name LIKE ?", name)
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
