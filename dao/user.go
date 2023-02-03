package dao

import (
	"bytes"
	"database/sql"
	"fmt"
	"rare_earth_mining_BE/model"
	"rare_earth_mining_BE/util"
	"strconv"
)

// 用于按照具有唯一性的字段来查找用户信息(不包括密码)
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUser(attribute string, value string) (u model.User, err error) {
	var row *sql.Row
	switch attribute {
	case "uID":
		//将uid的前面的"uid"字符去除并转化为int类型
		id, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(value), []byte("uID"))))
		if err != nil {
			return model.User{}, err
		}
		row = DB.QueryRow("select id, mail, nickname, position ,company, introduction, articleNum, collectNum, likeNum from users where id = ?", id)
	case "mail":
		row = DB.QueryRow("select id, mail, nickname, position ,company, introduction, articleNum, collectNum, likeNum from users where mail = ?", value)
	case "nickname":
		row = DB.QueryRow("select id, mail, nickname, position ,company, introduction, articleNum, collectNum, likeNum from users where nickname = ?", value)
	default:
		err = util.FieldsError
		return
	}

	if err = row.Err(); row.Err() != nil {
		return
	}

	err = row.Scan(&u.ID, &u.Mail, &u.Nickname, &u.Position, &u.Company, &u.Introduction, &u.ArticleNum, &u.CollectNum, &u.LikeNum)
	return
}

// 用于按照具有唯一性的字段来查找用户密码
// attributes: 用于查找user的字段,该字段必须有唯一性[uid,mail,nickname]
// value：字段的值
func SearchUserPassword(attribute string, value string) (uID string, password string, err error) {
	var row *sql.Row
	switch attribute {
	case "uID":
		//将uid的前面的"uID"字符去除并转化为int类型
		id, err := strconv.Atoi(string(bytes.TrimPrefix([]byte(value), []byte("uID"))))
		if err != nil {
			return "", "", err
		}
		row = DB.QueryRow("select ID, password from users where ID = ?", id)
	case "mail":
		row = DB.QueryRow("select ID, password from users where mail = ?", value)
	case "nickname":
		row = DB.QueryRow("select ID, password from users where nickname = ?", value)
	default:
		err = util.FieldsError
		return
	}

	if err = row.Err(); row.Err() != nil {
		return
	}

	err = row.Scan(&uID, &password)
	return
}

func CreateUser(u model.User) (err error) {
	fmt.Println("执行dao.CreatUser")
	_, err = DB.Exec("insert into users (mail, nickname, password) values (?,?,?)", u.Mail, u.Nickname, u.Password)
	return
}

func UserProfile(user model.User) (err error) {
	//oIDListStr := "\"" + strings.Trim(strings.Replace(fmt.Sprint(oIDList), " ", "\",\"", -1), "[]") + "\""
	fmt.Println(user)
	userStr := "\"" + user.Nickname + "\",\"" + user.Position + "\",\"" + user.Company + "\",\"" + user.Introduction + "\""
	fmt.Println("-----: ", userStr)
	_, err = DB.Exec("UPDATE users SET nickname = ?, position = ?, company = ?, introduction = ? WHERE ID = ?", user.Nickname, user.Position, user.Company, user.Introduction, user.ID)
	return
}
