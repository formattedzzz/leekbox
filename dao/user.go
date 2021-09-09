package dao

import (
	"leekbox/model"
)

func (this *GormDB) CreateNewUser(user *model.User) (*model.User, error) {
	if err := this.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (this *GormDB) GetUserByUid(user_id string) (*model.User, error) {
	user := new(model.User)
	if err := this.DB.First(user, "user_id = ?", user_id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (this *GormDB) GetUserById(uid int) (*model.User, error) {
	user := new(model.User)
	if err := this.DB.Find(user, uid).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (this *GormDB) GetUserSubRooms(uid int) ([]*model.RoomInfo, error) {
	rooms := []*model.RoomInfo{}
	sub := this.DB.Table("subscribes").Select("room_id").Where("uid = ?", uid)
	if err := this.DB.Table("rooms").Preload("Owner").Where("id IN (?) AND deleted = 0", sub).Find(&rooms).Error; err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if room.OwnerId == uid {
			room.IsOwner = true
		}
		room.Subscribed = true
	}
	return rooms, nil
}

func (this *GormDB) CheckUserExist(user_id string) bool {
	user := new(model.User)
	if err := this.DB.First(user, "user_id = ?", user_id).Error; err != nil {
		return false
	}
	if user.Id != 0 {
		return true
	}
	return false
}

func (this *GormDB) GetUserList(page, limit int) ([]model.User, error) {
	users := []model.User{}
	if err := this.DB.Limit(limit).Offset((page - 1) * limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (this *GormDB) UpdateUserInfo(user *model.User) (*model.User, error) {
	if err := this.DB.Model(user).Select("nick_name", "desc", "avatar", "phone", "email").Updates(*user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
