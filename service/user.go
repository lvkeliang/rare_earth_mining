package service

import (
	"rare_earth_mining_BE/dao"
	"rare_earth_mining_BE/model"
)

// 用于查找用户信息(不包括密码)
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUser(attribute string, value string) (u model.User, err error) {
	u, err = dao.SearchUser(attribute, value)
	return
}

// 用于查找用户密码
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUserPassword(attribute string, value string) (uID string, password string, err error) {
	uID, password, err = dao.SearchUserPassword(attribute, value)
	return
}

func CreateUser(u model.User) error {
	//fmt.Println("执行service.CreateUser")
	err := dao.CreateUser(u)
	return err
}
