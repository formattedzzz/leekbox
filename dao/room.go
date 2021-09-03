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

func (this *GormDB) GetRoomById(id int) (*model.RoomInfo, error) {
	room := model.RoomInfo{}
	if err := this.DB.Table("rooms").Model(&room).Preload("Owner").Find(&room, "id = ? and deleted = ?", id, 0).Error; err != nil {
		return nil, err
	}
	if room.Id > 0 {
		return &room, nil
	} else {
		return nil, fmt.Errorf(model.ROOM_MSG.ROOM_NOT_EXIST)
	}
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
