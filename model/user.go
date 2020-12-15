package model

import (
	"fmt"
	"go-filestore-server/database/mysql"
	"go-filestore-server/logger"
)

// User: 用户表结构体
type User struct {
	UserName     string `json:"user_name"`      // 用户名
	Email        string `json:"email"`          // 邮箱
	Phone        string `json:"phone"`          // 手机号
	SignupAt     string `json:"signup_at"`      // 登陆
	LastActiveAt string `json:"last_active_at"` // 最后一次活跃时间
	Status       int    `json:"status"`         // 状态
}

// 增加: 通过用户名+密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		logger.Info("failed to insert err:\t", err.Error())
		fmt.Println("failed to insert err:\t", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("failed to insert, err:\t", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// 查询: 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name =  ? limit 1")
	if err != nil {
		logger.Info("err:\t", err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		logger.Info(err.Error())
		return false
	} else if rows == nil {
		logger.Info("username not found:\t", username)
		return false
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 更新: 刷新用户登陆的token
func UpdateToken(username string, token string) bool {
	// replace into 首先尝试插入数据到表中，1.如果发现表中已经有此行数据（根据主键或者唯一索引判断）则先删除此行数据，然后插入新的数据。2.否则，直接插入新数据。
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token (`user_name`,`user_token`) value (?,?)")
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	return true
}

// 查询: 查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mysql.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		logger.Info(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.UserName, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// 查询用户是否存在
func UserExist(username string) (bool, error) {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_user where user_name=? limit 1")
	if err != nil {
		logger.Info(err.Error())
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}
