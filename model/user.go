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

type LoginReq struct {
	// 用户名
	Username string `json:"username" binding:"required"`
	// 密码
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

