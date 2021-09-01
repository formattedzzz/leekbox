package api

import (
	"fmt"
	"leekbox/model"
	"net/http"

	"leekbox/utils"

	"github.com/gin-gonic/gin"
)

// 用户相关事件回调中心
type UserEvent struct {
	AddedCallback []func(id int) error
	DeledCallback []func(id int) error
}

func (this *UserEvent) RegisterAddEvent(cb func(uid int) error) {
	this.AddedCallback = append(this.AddedCallback, cb)
}
func (this *UserEvent) DispatchAddEvent(uid int) error {
	for _, cb := range this.AddedCallback {
		if err := cb(uid); err != nil {
			return err
		}
	}
	return nil
}
func (this *UserEvent) RegisterDelEvent(cb func(uid int) error) {
	this.DeledCallback = append(this.DeledCallback, cb)
}
func (this *UserEvent) DispatchDelEvent(uid int) error {
	for _, cb := range this.DeledCallback {
		if err := cb(uid); err != nil {
			return err
		}
	}
	return nil
}

// 用户模块的操作数据库接口
type UserDB interface {
	CreateNewUser(user *model.User) (*model.User, error)
	GetUserByUid(id string) (*model.User, error)
	GetUserById(id int) (*model.User, error)
	CheckUserExist(user_id string) bool
}

type UserAPI struct {
	DB        UserDB
	UserEvent *UserEvent
}

type SignForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
	Pass   string `json:"pass" xml:"pass" form:"pass" binding:"required"`
}
type CheckForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
}

// @Summary 用户注册
// @Param body body SignForm true "结构体"
// @Router /api/user/signup [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserSignup(c *gin.Context) {
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
	fmt.Printf("\n%+v\n", signdata)
	if this.DB.CheckUserExist(signdata.UserId) {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: model.USER_EXISTED,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	newUser := model.User{
		UserId: signdata.UserId,
		Pass:   utils.MD5(signdata.Pass),
	}
	resp := model.Resp{}
	if user, err := this.DB.CreateNewUser(&newUser); err != nil {
		resp.Code = 50000
		resp.Message = err.Error()
		resp.Data = nil
		c.JSON(http.StatusBadGateway, resp)
	} else {
		resp.Code = 20000
		resp.Message = model.USER_SIGNUP_SUCCESS
		resp.Data = user
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 用户登录
// @Param body body SignForm true "结构体"
// @Router /api/user/login [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserLogin(c *gin.Context) {
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
	fmt.Printf("\n%+v\n", signdata)
	user, err := this.DB.GetUserByUid(signdata.UserId)
	if err != nil {
		resp := model.Resp{
			Code:    40000,
			Data:    nil,
			Message: model.USER_NOT_EXISTED,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp := model.Resp{}
	if user.Pass != utils.MD5(signdata.Pass) {
		resp.Code = 40000
		resp.Data = nil
		resp.Message = model.USER_PASS_INVALID
		c.JSON(http.StatusOK, resp)
		return
	}
	if token, err := GenToken(*user); err == nil {
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
// @Security ApiKeyAuth
// @Router /api/user/info [get]
// @Success 200 {object} model.Resp
func (this *UserAPI) GetUserInfo(c *gin.Context) {
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
	uid := userInfo.(model.User).Id
	resp := model.Resp{
		Code:    20000,
		Data:    nil,
		Message: model.USER_INFO_SUCCEED,
	}
	if user, err := this.DB.GetUserById(uid); err != nil {
		resp.Code = 40000
		resp.Message = model.USER_NOT_EXISTED
	} else {
		resp.Data = user
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 检查用户名是否可用
// @Param body body CheckForm true "结构体"
// @Router /api/user/check [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) CheckUserId(c *gin.Context) {
	checkform := CheckForm{}
	if err := c.ShouldBind(&checkform); err != nil {
		return
	}
	existed := this.DB.CheckUserExist(checkform.UserId)
	resp := model.Resp{
		Code: 20000,
		Data: map[string]interface{}{
			"existed": existed,
			"user_id": checkform.UserId,
		},
		Message: utils.If(existed, model.USER_EXISTED, model.USER_ACCESSIABLE).(string),
	}
	c.JSON(http.StatusOK, resp)
}
