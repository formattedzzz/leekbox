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
		return nil, fmt.Errorf("该讨论组不存在")
	}
}

func (this *GormDB) UpdateRoomInfo(room *model.Room) error {
	if err := this.DB.Model(room).Updates(*room).Error; err != nil {
		return err
	}
	return nil
}
