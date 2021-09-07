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
	GetRoomById(id int) (*model.RoomInfo, error)
	CreateNewRoom(room *model.Room) (*model.Room, error)
	UpdateRoomInfo(room *model.Room) (*model.Room, error)
	CreateNewComment(comment *model.Comment) (*model.Comment, error)
}

type RoomAPI struct {
	DB     RoomDB
	Stream *stream.StreamAPI
}

type RoomCreateForm struct {
	OwnerId int    `json:"owner_id" form:"owner_id"`
	Title   string `json:"title" form:"title" binding:"required"`
	Desc    string `json:"desc" form:"desc"`
	Avatar  string `json:"avatar" gorm:"avatar"`
	Status  int    `json:"status" form:"status"`
	Dev     bool   `json:"dev" form:"dev"`
}

type RoomUpdateForm struct {
	Id      int    `json:"id" form:"id" binding:"gte=1,required"`
	OwnerId int    `json:"owner_id" form:"owner_id" binding:"required"`
	Title   string `json:"title" form:"title" binding:"required"`
	Desc    string `json:"desc" form:"desc"`
	Avatar  string `json:"avatar" form:"avatar"`
	Status  int    `json:"status" form:"status"`
	Deleted int    `json:"deleted" form:"deleted"`
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
	RoomId  int    `json:"room_id" form:"room_id" binding:"required"`
	Type    int    `json:"type" form:"type"`
	Content string `json:"content" form:"content" binding:"required"`
	Attach  string `json:"attach" form:"attach"`
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
	if newComment, err := this.DB.CreateNewComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, model.Return(20000, newComment, model.API_SUCCESS))
	}
	if client_list, ok := this.Stream.RoomClient[comment.RoomId]; ok {
		for _, client := range client_list {
			client.Sender <- comment
		}
	}
}
