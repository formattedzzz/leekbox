package api

import (
	"fmt"
	"leekbox/api/stream"
	"leekbox/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户模块的操作数据库接口
type RoomDB interface {
	GetUserById(uid int) (*model.User, error)
	GetRoomById(id int) (*model.RoomInfo, error)
	CreateNewRoom(room *model.Room) (*model.Room, error)
	UpdateRoomInfo(room *model.Room) (*model.Room, error)
	CreateNewComment(comment *model.Comment) (*model.Comment, error)
	GetRoomComments(room_id int, page int, limit int) ([]*model.CommentItem, error)
}

type RoomAPI struct {
	DB     RoomDB
	Stream *stream.StreamAPI
}

type RoomCreateForm struct {
	OwnerId int    `json:"owner_id" form:"owner_id" example:"35"`
	Title   string `json:"title" form:"title" binding:"required" example:"韭菜盒子直播间"`
	Desc    string `json:"desc" form:"desc" example:"来点介绍吧"`
	Avatar  string `json:"avatar" gorm:"avatar" example:"https://theshy.cc/img/avatar.png"`
	Status  int    `json:"status" form:"status" example:"0"`
	Dev     bool   `json:"dev" form:"dev" example:"false"`
}

// swagger:model RoomUpdateForm
type RoomUpdateForm struct {
	Id      int    `json:"id" form:"id" binding:"gte=1,required" example:"1"`
	OwnerId int    `json:"owner_id" form:"owner_id" binding:"required" example:"35"`
	Title   string `json:"title" form:"title" binding:"required" example:"新名称"`
	Desc    string `json:"desc" form:"desc" example:"新简介"`
	Avatar  string `json:"avatar" form:"avatar" example:"https://theshy.cc/img/avatar.png"`
	Status  int    `json:"status" form:"status" example:"0"`
	Deleted int    `json:"deleted" form:"deleted" example:"0"`
}

// @Summary 创建讨论组
// @Param body body RoomCreateForm true "结构体"
// @Security ApiKeyAuth
// @Router /api/room/create [post]
// @Success 200 {object} model.Resp
func (this *RoomAPI) CreateNewRoom(c *gin.Context) {
	body := RoomCreateForm{}
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
	if body.Dev && body.OwnerId > 0 {
		uid = body.OwnerId
	}
	room := new(model.Room)
	room.OwnerId = uid
	room.Title = body.Title
	room.Desc = body.Desc
	room.Avatar = body.Avatar
	room.Status = body.Status
	if newRoom, err := this.DB.CreateNewRoom(room); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, model.Return(20000, newRoom, model.API_SUCCESS))
	}
}

// @Summary 获取讨论组信息
// @Param Authorization header string false "token"
// @Param id path int true "讨论组ID"
// @Router /api/room/{id} [get]
// @Success 200 {object} model.Resp
func (this *RoomAPI) GetRoomInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, model.PARAMS_ERROR))
		return
	}
	fmt.Printf("id: %v\n", id)
	if room, err := this.DB.GetRoomById(id); err != nil {
		c.JSON(http.StatusOK, model.Return(40000, nil, err.Error()))
	} else {
		// 该接口token为可带可不带
		if authed, exist := c.Get("authed"); authed.(bool) == true && exist {
			userInfo := c.MustGet("userInfo").(model.User)
			if userInfo.Id == room.OwnerId {
				room.IsOwner = true
			} else {
				room.IsOwner = false
			}
		}
		c.JSON(http.StatusOK, model.Return(20000, room, model.API_SUCCESS))
	}
}

// @Summary 修改讨论组
// @Description 删除讨论组时deleted=1 此操作不可逆
// @Param body body RoomUpdateForm true "结构体"
// @Security ApiKeyAuth
// @Router /api/room/update [put]
// @Success 200 {object} model.Resp
func (this *RoomAPI) UpdateRoomInfo(c *gin.Context) {
	body := RoomUpdateForm{}
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
	if body.OwnerId != uid {
		c.JSON(http.StatusOK, model.Return(40000, nil, model.ROOM_MSG.ROOM_FORBIDDEN))
		return
	}
	room := new(model.Room)
	room.Id = body.Id
	room.OwnerId = uid
	room.Title = body.Title
	room.Desc = body.Desc
	room.Avatar = body.Avatar
	room.Status = body.Status
	room.Deleted = body.Deleted
	if room, err := this.DB.UpdateRoomInfo(room); err != nil {
		c.JSON(http.StatusInternalServerError, model.Return(50000, nil, err.Error()))
	} else {
		c.JSON(http.StatusOK, model.Return(20000, room, model.API_SUCCESS))
	}
}

type CommentCreateForm struct {
	RoomId  int    `json:"room_id" form:"room_id" binding:"required" example:"1"`
	Type    int    `json:"type" form:"type" example:"0"`
	Atsb    int    `json:"atsb" form:"atsb" example:"0"`
	Content string `json:"content" form:"content" binding:"required" example:"说点儿好听的吧～"`
	Attach  string `json:"attach" form:"attach" example:"{}"`
}

// @Summary 创建发言
// @Param body body CommentCreateForm true "结构体"
// @Security ApiKeyAuth
// @Router /api/comment/create [post]
// @Success 200 {object} model.Resp
func (this *RoomAPI) CreateNewComment(c *gin.Context) {
	body := CommentCreateForm{}
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
	comment := new(model.Comment)
	comment.Uid = uid
	comment.Content = body.Content
	comment.RoomId = body.RoomId
	comment.Attach = body.Attach
	comment.Type = body.Type
	comment.Atsb = body.Atsb
	if newComment, err := this.DB.CreateNewComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, model.Return(20000, newComment, model.API_SUCCESS))
	}
	var atsb *model.User = nil
	if comment.Atsb > 0 {
		if usr, err := this.DB.GetUserById(comment.Atsb); err == nil {
			atsb = usr
		}
	}
	this.Stream.PushRoomComment(comment, userInfo.(model.User), atsb)
}

type CommentQueryForm struct {
	Page   int `json:"page" form:"page" example:"1"`
	Limit  int `json:"limit" form:"limit" example:"30"`
	RoomId int `json:"room_id" form:"room_id" binding:"required"`
}

// @Summary 获取讨论组历史发言
// @Param Authorization header string false "token"
// @Param room_id query string true "讨论组ID"
// @Param limit query string false "limit" mininum(30) maxinum(100) default(30)
// @Param page query string false "page" mininum(1) default(1)
// @Router /api/room/comments [get]
// @Success 200 {object} model.Resp
func (this *RoomAPI) GetRoomComments(c *gin.Context) {
	query := new(CommentQueryForm)
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, model.Return(40000, nil, err.Error()))
		return
	}
	if query.Limit == 0 {
		query.Limit = 30
	}
	if query.Page == 0 {
		query.Page = 1
	}
	fmt.Printf("query: %v\n", query)
	if comments, err := this.DB.GetRoomComments(query.RoomId, query.Page, query.Limit); err != nil {
		c.JSON(http.StatusInternalServerError, model.Return(50000, query, err.Error()))
	} else {
		if authed, exist := c.Get("authed"); authed.(bool) == true && exist {
			userInfo := c.MustGet("userInfo").(model.User)
			// 对切片循环需要注意 如果元素不是指针类型的结构体 遍历中的item是个副本
			for _, comment := range comments {
				if comment.Uid == userInfo.Id {
					comment.IsOwner = true
				}
			}
		}
		c.JSON(http.StatusOK, model.Return(20000, comments, model.API_SUCCESS))
	}
}
