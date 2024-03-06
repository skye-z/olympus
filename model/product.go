package model

import "xorm.io/xorm"

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

func (model ProductModel) QueryProduct(processor int16, group, name string) *Product {
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
