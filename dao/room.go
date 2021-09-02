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
	if err := this.DB.Table("rooms").Model(&room).Preload("Owner").Find(&room, id).Error; err != nil {
		return nil, err
	}
	if room.Id > 0 {
		return &room, nil
	} else {
		return nil, fmt.Errorf("改讨论组不存在")
	}
}
