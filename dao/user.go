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
