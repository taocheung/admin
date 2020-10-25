package model

import (
	"gorm.io/gorm"
	"time"
)

type Resource struct {
	Id int `json:"id"`
	Phone string `json:"phone"`
	Account string `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Resource) TableName() string {
	return "resource"
}

func ResourceImport(data []Resource) (int64, error) {
	var (
		i int64
	)

	tx := db.Session(&gorm.Session{PrepareStmt: true}).Begin()

	for _, v := range data {
		result := tx.Model(&Resource{}).Where("account = ?", v.Account).FirstOrCreate(&v)
		if result.Error != nil {
			tx.Rollback()
			return 0, result.Error
		}
		i++
	}
	tx.Commit()
	return i, nil
}

func ResourceExport(ids []int) ([]Resource, error) {
	var list []Resource
	err := db.Model(&Resource{}).Where("id in ?", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

type ResourceListRsp struct {
	Id int `json:"id"`
	Phone string `json:"phone"`
	Account string `json:"account"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func ResourceList(account []string) ([]Resource, error) {
	var list []Resource

	if len(account) == 0 {
		err := db.Model(&Resource{}).Order("id asc").Limit(10).Find(&list).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := db.Model(&Resource{}).Where("account in ?", account).Find(&list).Error
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}