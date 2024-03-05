package model

import "xorm.io/xorm"

type Product struct {
	Id int64 `json:"id"`
	// 制品处理器(1Maven、2NPM)
	Processor int64 `json:"processor"`
	// 制品组
	Group string `json:"group"`
	// 制品名
	Name string `json:"name"`
	// 添加时间
	AddTime string `json:"addTime"`
	// 更新时间
	UpdateTime string `json:"updateTime"`
}

type ProductModel struct {
	DB *xorm.Engine
}
