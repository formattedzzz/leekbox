package api

import (
	"fmt"
	"leekbox/model"
	"net/http"
	"regexp"

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

type LoginForm struct {
	UserId string `json:"user_id" form:"user_id" binding:"required"`
	Pass   string `json:"pass" form:"pass" binding:"required"`
}
type SignForm struct {
	LoginForm
	Repass string `json:"repass" form:"repass" binding:"required,eqfield=Pass"`
}
type CheckForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
}

func preHandleUserId(user_id string) (string, error) {
	if user_id == "" {
		return "", fmt.Errorf("用户ID为空")
	}
	user_id = regexp.MustCompile(`^\s+`).ReplaceAllString(user_id, "")
	user_id = regexp.MustCompile(`\s+$`).ReplaceAllString(user_id, "")
	if match := regexp.MustCompile(`\s`).MatchString(user_id); match == true {
		return "", fmt.Errorf("用户ID不能含有空格")
	}
	if match := regexp.MustCompile(`^\w+$`).MatchString(user_id); match == false {
		return "", fmt.Errorf("用户ID只能包含字母数字下划线")
	}
	return user_id, nil
}

// @Summary 用户注册
// @Param body body SignForm true "结构体"
// @Router /api/user/signup [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserSignup(c *gin.Context) {
	signdata := SignForm{}
	if err := c.ShouldBind(&signdata); err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	// 检查用户ID是否合法
	if user_id, err := preHandleUserId(signdata.UserId); err != nil {
		c.JSON(http.StatusOK, model.Return(40000, user_id, err.Error()))
		return
	} else {
		signdata.UserId = user_id
	}
	fmt.Printf("\n%+v\n", signdata)
	if this.DB.CheckUserExist(signdata.UserId) {
		c.JSON(http.StatusOK, model.Return(40000, nil, model.USER_MSG.USER_EXISTED))
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
		c.JSON(http.StatusInternalServerError, resp)
	} else {
		resp.Code = 20000
		resp.Message = model.API_SUCCESS
		resp.Data = user
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary 用户登录
// @Param body body LoginForm true "结构体"
// @Router /api/user/login [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserLogin(c *gin.Context) {
	loginData := LoginForm{}
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	fmt.Printf("\n%+v\n", loginData)
	user, err := this.DB.GetUserByUid(loginData.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, model.USER_MSG.USER_NOT_EXISTED))
		return
	}
	if user.Pass != utils.MD5(loginData.Pass) {
		c.JSON(http.StatusOK, model.Return(40000, nil, model.USER_MSG.USER_PASS_INVALID))
		return
	}
	resp := model.Resp{}
	if token, err := GenToken(*user); err == nil {
		resp.Code = 20000
		resp.Data = map[string]string{
			"token": token,
		}
		resp.Message = model.API_SUCCESS
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
	userInfo := c.MustGet("userInfo")
	if userInfo == nil {
		c.JSON(http.StatusOK, model.Return(50000, nil, model.UNHANDLED_ERROR))
		return
	}
	uid := userInfo.(model.User).Id
	resp := model.Resp{
		Code:    20000,
		Data:    nil,
		Message: model.API_SUCCESS,
	}
	if user, err := this.DB.GetUserById(uid); err != nil {
		resp.Code = 40000
		resp.Message = model.USER_MSG.USER_NOT_EXISTED
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
		Message: utils.If(existed, model.USER_MSG.USER_EXISTED, model.API_SUCCESS).(string),
	}
	c.JSON(http.StatusOK, resp)
}
