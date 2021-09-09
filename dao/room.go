package dao

import (
	"fmt"
	"leekbox/model"
)

func (this *GormDB) CreateNewRoom(room *model.Room) (*model.Room, error) {
	if err := this.DB.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (this *GormDB) CreateNewComment(comment *model.Comment) (*model.Comment, error) {
	if err := this.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (this *GormDB) CreateSubscribe(sub *model.Subscribe) (*model.Subscribe, error) {
	current := new(model.Subscribe)
	this.DB.Where(sub).Find(current)
	if current.Id > 0 {
		return current, nil
	}
	if err := this.DB.Create(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

func (this *GormDB) CancelSubscribe(sub *model.Subscribe) error {
	return this.DB.Where(sub).Delete(sub).Error
}

func (this *GormDB) GetRoomById(id int) (*model.RoomInfo, error) {
	room := new(model.RoomInfo)
	if err := this.DB.Table("rooms").Model(room).Preload("Owner").Find(room, "id = ? and deleted = ?", id, 0).Error; err != nil {
		return nil, err
	}
	if room.Id > 0 {
		return room, nil
	} else {
		return nil, fmt.Errorf(model.ROOM_MSG.ROOM_NOT_EXIST)
	}
}

func (this *GormDB) CheckRoomSubsscribed(room_id, uid int) bool {
	sub := new(model.Subscribe)
	this.DB.Find(sub, "room_id = ? and uid = ?", room_id, uid)
	if sub.Id > 0 {
		return true
	}
	return false
}

func (this *GormDB) GetRoomComments(room_id int, page int, limit int) ([]*model.CommentItem, error) {
	comments := []*model.CommentItem{}
	offset := limit * (page - 1)
	if err := this.DB.Table("comments").Model(comments).Preload("Owner").Preload("Refer").Limit(limit).Offset(offset).Order("id desc").Find(&comments, "room_id = ?", room_id).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (this *GormDB) UpdateRoomInfo(room *model.Room) (*model.Room, error) {
	currentRoom := new(model.Room)
	this.DB.Find(currentRoom, room.Id)
	if currentRoom.Id > 0 {
		if currentRoom.Deleted != 0 {
			return nil, fmt.Errorf(model.ROOM_MSG.ROOM_DELETED)
		}
		if currentRoom.OwnerId != room.OwnerId {
			return nil, fmt.Errorf(model.ROOM_MSG.ROOM_NOT_MATCH)
		}
	}
	if err := this.DB.Model(room).Select("title", "desc", "avatar", "status", "deleted").Updates(*room).Error; err != nil {
		return nil, err
	}
	this.DB.Find(currentRoom, room.Id)
	return currentRoom, nil
}
