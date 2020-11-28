package controller

import (
	"admin/middleware"
	"admin/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// @Summary 用户登录
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer"
// @Param object body model.LoginReq false "查询参数"
// @Success 0 {object} util.Response
// @Router /user/login [post]
func Login(c *gin.Context) {
	var err error

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
	Response(c, map[string]string{"token": token, "role": fmt.Sprintf("%d", user.Role)})
}