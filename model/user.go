package model

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id        int        `json:"id"`
	Username  string     `json:"username"`
	RealName  string     `json:"real_name"`
	Phone     string     `json:"phone"`
	Password  string     `json:"password"`
	Role      int        `json:"role"`
	ApplyTime time.Time  `json:"apply_time"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u *User) TableName() string {
	return "user"
}

var layout = "2006-01-02"

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code"`
}

func Login(req *LoginReq) (*User, error) {
	var user User

	tx := db.WithContext(context.Background())
	err := tx.Model(&User{}).Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}
	if user.ApplyTime.Before(time.Now()) {
		return nil, errors.New("账号过期")
	}
	return &user, nil
}

type AddUserReq struct {
	Username  string `json:"username" binding:"required"`
	RealName  string `json:"real_name"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone"`
	ApplyTime string `json:"apply_time"`
}

func AddUser(req *AddUserReq) error {
	var applyTime time.Time

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if req.ApplyTime != "" {
		applyTime, err = time.ParseInLocation(layout, req.ApplyTime, time.Local)
		if err != nil {
			return err
		}
	} else {
		applyTime = time.Now()
	}
	err = db.Model(&User{}).Create(&User{
		Username:  req.Username,
		RealName:  req.RealName,
		Phone:     req.Phone,
		Password:  string(password),
		ApplyTime: applyTime,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

type UpdateUserReq struct {
	Id        int    `json:"id" binding:"required"`
	Username  string `json:"username"`
	RealName  string `json:"real_name"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	ApplyTime string `json:"apply_time"`
}

func UpdateUser(req *UpdateUserReq) error {
	var (
		password  []byte
		applyTime time.Time
		err       error
	)

	data := make(map[string]interface{})
	if req.ApplyTime != "" {
		applyTime, err = time.ParseInLocation(layout, req.ApplyTime, time.Local)
		if err != nil {
			return err
		}
		data["apply_time"] = applyTime
	} else {
		applyTime = time.Now()
	}

	if req.Password != "" {
		password, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data["password"] = password
	}
	if req.Username != "" {
		data["username"] = req.Username
	}
	if req.RealName != "" {
		data["real_name"] = req.RealName
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}

	err = db.Model(&User{}).Where("id = ?", req.Id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

type DeleteUserReq struct {
	Id int `json:"id"`
}

func DeleteUser(req *DeleteUserReq) error {
	err := db.Model(&User{}).Where("id = ?", req.Id).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

type ListUserReq struct {
	PageId   int    `json:"page_id"`
	PageSize int    `json:"page_size"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	StarTime string `json:"star_time"`
	EndTime  string `json:"end_time"`
}

type ListUserRsp struct {
	Count int64      `json:"count"`
	List  []UserItem `json:"list"`
}

type UserItem struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	RealName  string `json:"real_name"`
	Phone     string `json:"phone"`
	Role      int    `json:"role"`
	ApplyTime string `json:"apply_time"`
	CreatedAt string `json:"created_at"`
}

func ListUser(req *ListUserReq) (*ListUserRsp, error) {
	var (
		users []User
		count int64
		list  []UserItem
	)

	tx := db.WithContext(context.Background())
	tx = tx.Model(&User{})
	if req.Username != "" {
		tx = tx.Where("username like ?", "%"+req.Username+"%")
	}
	if req.RealName != "" {
		tx = tx.Where("real_name like ?", "%"+req.RealName+"%")
	}
	if req.Phone != "" {
		tx = tx.Where("phone = ?", req.Phone)
	}
	if req.StarTime != "" {
		tx = tx.Where("apply_time > ?", req.StarTime)
	}
	if req.EndTime != "" {
		tx = tx.Where("apply_time < ?", req.EndTime)
	}
	if req.PageId <= 0 {
		req.PageId = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	err := tx.Count(&count).
		Offset((req.PageId - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("id asc").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		list = append(list, UserItem{
			Id:        v.Id,
			Username:  v.Username,
			RealName:  v.RealName,
			Phone:     v.Phone,
			Role:      v.Role,
			ApplyTime: v.ApplyTime.Format(layout),
			CreatedAt: v.CreatedAt.Format(layout),
		})
	}
	return &ListUserRsp{List: list, Count: count}, nil
}
