package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
	"github.com/skye-z/olympus/processor"
	"github.com/skye-z/olympus/util"
)

type ProductService struct {
	Product model.ProductModel
	Version model.VersionModel
}

func NewProductService(product model.ProductModel, version model.VersionModel) *ProductService {
	ps := new(ProductService)
	ps.Product = product
	ps.Version = version
	return ps
}

// 获取制品统计
func (ps ProductService) Stat(ctx *gin.Context) {
	data := ps.Product.Stat()
	if data == nil {
		util.ReturnMessage(ctx, false, "获取制品统计信息失败")
	} else {
		util.ReturnData(ctx, true, data)
	}
}

// 搜索制品
func (ps ProductService) Search(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if len(keyword) == 0 {
		util.ReturnMessage(ctx, false, "搜索关键词不能为空")
		return
	} else {
		keyword = "%" + keyword + "%"
	}
	processor, err := strconv.Atoi(ctx.Query("processor"))
	if err != nil {
		processor = 0
	}
	list, err := ps.Product.GetList(processor, "", keyword, 1, 10)
	if err != nil {
		util.ReturnMessage(ctx, false, "搜索制品失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

// 获取制品数量
func (ps ProductService) GetNumber(ctx *gin.Context) {
	group := ctx.Query("group")
	keyword := ctx.Query("keyword")
	if len(keyword) == 0 {
		util.ReturnMessage(ctx, false, "搜索关键词不能为空")
		return
	} else {
		keyword = "%" + keyword + "%"
	}
	processor, err := strconv.Atoi(ctx.Query("processor"))
	if err != nil {
		processor = 0
	}
	number, err := ps.Product.GetNumber(processor, group, keyword)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取制品数量失败")
	} else {
		util.ReturnData(ctx, true, number)
	}
}

// 获取制品列表
func (ps ProductService) GetList(ctx *gin.Context) {
	group := ctx.Query("group")
	keyword := ctx.Query("keyword")
	if len(keyword) == 0 {
		util.ReturnMessage(ctx, false, "搜索关键词不能为空")
		return
	} else {
		keyword = "%" + keyword + "%"
	}
	processor, err := strconv.Atoi(ctx.Query("processor"))
	if err != nil {
		processor = 0
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	number, err := strconv.Atoi(ctx.Query("number"))
	if err != nil {
		number = 20
	}
	list, err := ps.Product.GetList(processor, group, keyword, page, number)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取制品列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}

type ProductInfo struct {
	Info    *model.Product  `json:"info"`
	Version []model.Version `json:"version"`
}

// 获取制品详情
func (ps ProductService) GetInfo(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	product := ps.Product.GetItem(id)
	if product == nil {
		util.ReturnMessage(ctx, false, "制品不存在")
		return
	}
	util.ReturnData(ctx, true, product)
}

// 获取制品详情
func (ps ProductService) GetNpmConfig(ctx *gin.Context) {
	name := ctx.Param("name")
	npmPro := &processor.Npm{
		Product: ps.Product,
		Version: ps.Version,
	}
	data := npmPro.GetConfig(name)
	if data == nil {
		util.ReturnMessage(ctx, false, "制品不存在")
		return
	}
	ctx.Data(200, "application/json; charset=utf-8", data)
	ctx.Abort()
}

// 获取版本列表
func (ps ProductService) GetVersionList(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	list, err := ps.Version.GetList(id)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取版本列表失败")
		return
	}
	util.ReturnData(ctx, true, list)
}
