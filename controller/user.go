package controller

import (
	"admin/middleware"
	"admin/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			Error(c, err)
		}
	}()

	var req model.LoginReq
	if err = c.Bind(&req); err != nil {
		Error(c, err)
		return
	}
	user, err := model.Login(&req)
	if err != nil {
		Error(c, err)
		return
	}

	token, err := middleware.NewJWT().CreateToken(middleware.CustomClaims{
		ID:             user.Id,
		Username:       user.Username,
		Role:           user.Role,
		StandardClaims: jwt.StandardClaims{},
	})
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, map[string]string{"token": token})
}

func AddUser(c *gin.Context)  {
	var req model.AddUserReq
	if err := c.Bind(&req); err != nil {
		Error(c, err)
		return
	}
	err := model.AddUser(&req)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, nil)
}


func UpdateUser(c *gin.Context)  {
	var req model.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		Error(c, err)
		return
	}
	err := model.UpdateUser(&req)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, nil)
}


func DeleteUser(c *gin.Context)  {
	var req model.DeleteUserReq
	if err := c.Bind(&req); err != nil {
		Error(c, err)
		return
	}
	err := model.DeleteUser(&req)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, nil)
}


func ListUser(c *gin.Context)  {
	var req model.ListUserReq
	if err := c.Bind(&req); err != nil {
		Error(c, err)
		return
	}
	list, err := model.ListUser(&req)
	if err != nil {
		Error(c, err)
		return
	}
	Response(c, list)
}