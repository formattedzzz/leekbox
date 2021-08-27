package handler

import (
	"fmt"
	"leekbox/dao"
	"leekbox/model"
	"net/http"

	"leekbox/utils"

	"github.com/gin-gonic/gin"
)

type SignForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
	Pass   string `json:"pass" xml:"pass" form:"pass" binding:"required"`
}
type CheckForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
}

// @Tags 用户相关
// @Summary 用户注册
// @Accept json
// @Param body body SignForm true "结构体"
// @Router /api/user/signup [post]
func UserSignup(c *gin.Context) {
	signdata := SignForm{}
	if err := c.ShouldBind(&signdata); err != nil {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	fmt.Printf("%+v\n\n", signdata)
	currentUser := model.User{}
	dao.DB.First(&currentUser, "user_id = ?", signdata.UserId)
	if currentUser.Id != 0 {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: model.USER_EXISTED,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	newUser := model.User{
		UserId: signdata.UserId,
		Pass:   utils.MD5(signdata.Pass),
	}
	resp := model.Resp{}
	if user, err := dao.CreateNewUser(&newUser); err != nil {
		resp.Code = 50000
		resp.Message = err.Error()
		resp.Data = nil
		c.JSON(http.StatusOK, resp)
	} else {
		resp.Code = 20000
		resp.Message = model.USER_SIGNUP_SUCCESS
		resp.Data = user
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 用户登录
// @Accept json
// @Param body body SignForm true "结构体"
// @Router /api/user/login [post]
func UserLogin(c *gin.Context) {
	signdata := SignForm{}
	if err := c.ShouldBind(&signdata); err != nil {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	fmt.Printf("%+v\n\n", signdata)
	currentUser := model.User{}
	dao.DB.First(&currentUser, "user_id = ?", signdata.UserId)
	if currentUser.Id == 0 {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: model.USER_NOT_EXISTED,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp := model.Resp{}
	if currentUser.Pass != utils.MD5(signdata.Pass) {
		resp.Code = 40000
		resp.Data = nil
		resp.Message = model.USER_PASS_INVALID
		c.JSON(http.StatusOK, resp)
		return
	}
	if token, err := GenToken(currentUser); err == nil {
		resp.Code = 20000
		resp.Data = map[string]string{
			"token": token,
		}
		resp.Message = model.USER_LOGIN_SUCCESS
		c.JSON(http.StatusOK, resp)
	} else {
		resp.Code = 50000
		resp.Data = nil
		resp.Message = err.Error()
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 获取用户信息
// @Accept json
// @Param Authorization header string true "token"
// @Router /api/user/info [get]
func GetUserInfo(c *gin.Context) {
	// id, err := strconv.Atoi(c.Param("id"))
	userInfo := c.MustGet("userInfo")
	if userInfo == nil {
		resp := model.Resp{
			Code:    50000,
			Data:    nil,
			Message: model.UNHANDLED_ERROR,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	resp := model.Resp{
		Code:    20000,
		Data:    dao.GetUserInfoById(userInfo.(model.User).Id),
		Message: model.USER_INFO_SUCCEED,
	}
	c.JSON(http.StatusOK, resp)

}

// @Summary 检查用户名是否可用
// @Accept json
// @Param body body CheckForm true "结构体"
// @Router /api/user/check [post]
func CheckUserId(c *gin.Context) {
	checkform := CheckForm{}
	if err := c.ShouldBind(&checkform); err != nil {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	currentUser := model.User{}
	dao.DB.First(&currentUser, "user_id = ?", checkform.UserId)
	resp := model.Resp{
		Code: 20000,
		Data: map[string]interface{}{
			"existed": utils.If(currentUser.Id != 0, true, false),
			"user_id": checkform.UserId,
		},
		Message: utils.If(currentUser.Id != 0, model.USER_EXISTED, model.USER_ACCESSIABLE).(string),
	}
	c.JSON(http.StatusOK, resp)
}
