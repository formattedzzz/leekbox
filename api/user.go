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
	UpdateUserInfo(user *model.User) (*model.User, error)
}

type UserAPI struct {
	DB        UserDB
	UserEvent *UserEvent
}

type UserLoginForm struct {
	UserId string `json:"user_id" form:"user_id" binding:"required"`
	Pass   string `json:"pass" form:"pass" binding:"required"`
}
type UserSignForm struct {
	UserId   string `json:"user_id" form:"user_id" binding:"required"`
	Pass     string `json:"pass" form:"pass" binding:"required"`
	Repass   string `json:"repass" form:"repass" binding:"required,eqfield=Pass"`
	NickName string `json:"nick_name" form:"nick_name" binding:"required"`
}
type UserCheckForm struct {
	UserId string `json:"user_id" xml:"user_id" form:"user_id" binding:"required"`
}

type UserUpdateForm struct {
	Id       int    `json:"id" form:"id" binding:"gte=1,required"`
	NickName string `json:"nick_name" form:"nick_name"`
	Desc     string `json:"desc" form:"desc"`
	Avatar   string `json:"avatar" form:"avatar"`
	Email    string `json:"email" form:"email" binding:"email"`
	Phone    string `json:"phone" form:"phone"`
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
// @Param body body UserSignForm true "结构体"
// @Router /api/user/signup [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserSignup(c *gin.Context) {
	body := UserSignForm{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	// 检查用户ID是否合法
	if user_id, err := preHandleUserId(body.UserId); err != nil {
		c.JSON(http.StatusOK, model.Return(40000, user_id, err.Error()))
		return
	} else {
		body.UserId = user_id
	}
	fmt.Printf("\n%+v\n", body)
	if this.DB.CheckUserExist(body.UserId) {
		c.JSON(http.StatusOK, model.Return(40000, nil, model.USER_MSG.USER_EXISTED))
		return
	}
	newUser := model.User{
		UserId:   body.UserId,
		Pass:     utils.MD5(body.Pass),
		NickName: body.NickName,
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

// @Summary 修改用户信息
// @Security ApiKeyAuth
// @Param body body UserUpdateForm true "结构体"
// @Router /api/user/update [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UpdateUserInfo(c *gin.Context) {
	body := UserUpdateForm{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	userInfo := c.MustGet("userInfo")
	if userInfo == nil {
		c.JSON(http.StatusOK, model.Return(50000, nil, model.UNHANDLED_ERROR))
		return
	}
	uid := userInfo.(model.User).Id
	if uid != body.Id {
		c.JSON(http.StatusForbidden, model.Return(40300, nil, model.USER_MSG.USER_FORBIDDEN))
		return
	}
	fmt.Printf("\n%+v\n", body)
	updateUser := new(model.User)
	updateUser.Id = uid
	updateUser.Avatar = body.Avatar
	updateUser.NickName = body.NickName
	updateUser.Desc = body.Desc
	updateUser.Email = body.Email
	updateUser.Phone = body.Phone
	if user, err := this.DB.UpdateUserInfo(updateUser); err != nil {
		c.JSON(http.StatusInternalServerError, model.Return(50000, nil, err.Error()))
	} else {
		c.JSON(http.StatusOK, model.Return(50000, user, model.API_SUCCESS))
	}
}

// @Summary 用户登录
// @Param body body UserLoginForm true "结构体"
// @Router /api/user/login [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) UserLogin(c *gin.Context) {
	loginData := UserLoginForm{}
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
// @Param body body UserCheckForm true "结构体"
// @Router /api/user/check [post]
// @Success 200 {object} model.Resp
func (this *UserAPI) CheckUserId(c *gin.Context) {
	checkform := UserCheckForm{}
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
