package server

import (
	"adp_backend/config"
	"adp_backend/model"
	"time"
)

func GetUserByUserName(e *config.Env, userName string) (*model.User, error) {
	var user model.User

	db := e.MysqlCli.Table(user.TableName()).Where("username = ?", userName).First(&user)
	if db.Error != nil {
		return nil, db.Error
	}

	return &user, nil
}

func FirstUser(e *config.Env) (*model.User, error) {
	var user model.User

	db := e.MysqlCli.Table(user.TableName()).First(&user)
	if db.Error != nil {
		return nil, db.Error
	}

	return &user, nil
}

func FindAllUser(e *config.Env, limit int32, offset int32, search, userName string) ([]model.User, int64, error) {
	var user []model.User
	var cnt int64

	db := e.MysqlCli.Where("name = ?", userName).Limit(int(limit)).Offset(int(offset - 1)).Find(&user).Count(&cnt)
	if db.Error != nil {
		return nil, 0, db.Error
	}

	return user, cnt, nil
}

func AddUser(e *config.Env, userName, passWord, mobile, email, remark string) (*model.User, error) {
	user := model.User{
		UserName:  userName,
		Password:  passWord,
		Mobile:    mobile,
		Email:     email,
		Remark:    remark,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := e.MysqlCli.Select(&user).Create(&user)
	if db.Error != nil {
		return nil, db.Error
	}

	return &user, nil
}

func UpdateUser(e *config.Env, userId int32, username, password, mobile, email, remark string) (*model.User, error) {
	var user model.User
	updateValues := map[string]interface{}{
		"UserName": username,
		"Password": password,
		"Mobile":   mobile,
		"Email":    email,
		"Remark":   remark,
	}

	db := e.MysqlCli.Model(&user).Where("ID = ?", userId).Updates(updateValues)
	if db.Error != nil {
		return nil, db.Error
	}

	return &user, nil
}

func DeleteUser(e *config.Env, userId int32) (*model.User, error) {
	var user model.User

	db := e.MysqlCli.Delete(&user, userId)
	if db.Error != nil {
		return nil, db.Error
	}

	return nil, nil
}

func CreateUser(e *config.Env, user model.User) error {
	db := e.MysqlCli.Table(user.TableName()).Create(&user)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
