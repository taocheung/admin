package model

import "time"

type Resource struct {
	Id int `json:"id"`
	Phone string `json:"phone"`
	Account string `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Resource) TableName() string {
	return "resource"
}

func ResourceImport(data []Resource) error {
	err := db.Model(&Resource{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func ResourceExport(ids []int) ([]Resource, error) {
	var list []Resource
	err := db.Model(&Resource{}).Where("id in ?", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func ResourceList(account []string) ([]Resource, error) {
	var list []Resource

	err := db.Model(&Resource{}).Where("account in ?", account).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}