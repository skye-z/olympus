package model

import "xorm.io/xorm"

type Version struct {
	Id int64 `json:"id"`
	// 制品编号
	PId int64 `json:"pid"`
	// 版本号
	Number string `json:"name"`
	// 添加时间
	AddTime int64 `json:"addTime"`
}

type VersionModel struct {
	DB *xorm.Engine
}

func (model VersionModel) AddVersion(version *Version) bool {
	_, err := model.DB.Insert(version)
	return err == nil
}

func (model VersionModel) EditVersion(version *Version) bool {
	if version.Id == 0 {
		return false
	}
	_, err := model.DB.ID(version.Id).Update(version)
	return err == nil
}

func (model VersionModel) DelVersion(version *Version) bool {
	if version.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(version)
	return err == nil
}

func (model VersionModel) QueryVersion(pid int64, number string) *Version {
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
