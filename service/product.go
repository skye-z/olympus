package service

import (
	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
)

type ProductService struct {
	ProductModel model.ProductModel
	VersionModel model.VersionModel
}

func NewProductService(product model.ProductModel, version model.VersionModel) *ProductService {
	ps := new(ProductService)
	ps.ProductModel = product
	ps.VersionModel = version
	return ps
}

// 获取制品统计
func (ps ProductService) Stat(ctx *gin.Context) {

}

// 搜索制品
func (ps ProductService) Search(ctx *gin.Context) {

}

// 获取制品数量
func (ps ProductService) GetNumber(ctx *gin.Context) {

}

// 获取制品列表
func (ps ProductService) GetList(ctx *gin.Context) {

}

// 获取制品详情
func (ps ProductService) GetInfo(ctx *gin.Context) {

}
