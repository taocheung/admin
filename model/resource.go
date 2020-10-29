package model

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Resource struct {
	Id        int       `json:"id"`
	Phone     string    `json:"phone"`
	Account   string    `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Resource) TableName() string {
	return "resource"
}

func ResourceImport(data []Resource) (int64, error) {
	tx := db.Session(&gorm.Session{PrepareStmt: true})

	sql := "INSERT IGNORE INTO resource (account, phone) VALUES "
	for _, v := range data {
		sql += fmt.Sprintf("('%s', '%s'),", v.Account, v.Phone)
	}
	result := tx.Exec(strings.Trim(sql, ","))
	if result.Error != nil {

	}
	return result.RowsAffected, nil
}

type ResourceExportReq struct {
	ID []int `json:"id"`
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
	Id        int    `json:"id"`
	Phone     string `json:"phone"`
	Account   string `json:"account"`
	Status    string `json:"status"`
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
