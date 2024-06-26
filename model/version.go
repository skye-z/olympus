package model

import (
	"xorm.io/xorm"
)

type Version struct {
	Id int64 `json:"id"`
	// 制品编号
	PId int64 `json:"pid"`
	// 版本号
	Number string `json:"number"`
	// 添加时间
	AddTime int64 `json:"addTime"`
}

type VersionModel struct {
	DB *xorm.Engine
}

func (model VersionModel) Add(version *Version) bool {
	_, err := model.DB.Insert(version)
	return err == nil
}

func (model VersionModel) Edit(version *Version) bool {
	if version.Id == 0 {
		return false
	}
	_, err := model.DB.ID(version.Id).Update(version)
	return err == nil
}

func (model VersionModel) Del(version *Version) bool {
	if version.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(version)
	return err == nil
}

func (model VersionModel) GetList(pid int64) ([]Version, error) {
	var list []Version
	err := model.DB.Where("p_id = ?", pid).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (model VersionModel) Query(pid int64, number string) *Version {
	version := &Version{
		PId:    pid,
		Number: number,
	}
	has, _ := model.DB.Get(version)
	if !has {
		return version
	}
	return nil
}
