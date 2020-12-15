package orm

import (
	"fmt"
	"go-filestore-server/database/mysql"
	"log"
)

func UserSignup(username string, passwd string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		log.Println("failed to insert, err:\t", err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		res.Suc = false
		return
	}
	res.Suc = false
	res.Msg = "无更新记录"
	return
}

func UserSignin(username string, encpwd string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		log.Println("username not found:\t", username)
		res.Suc = false
		res.Msg = "用户名未注册"
		return
	}

	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		res.Suc = true
		res.Data = true
		return
	}
	res.Suc = false
	res.Msg = "用户名/密码不匹配"
	return
}

func UpdateToken(username string, token string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values(?,?)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

func GetUserInfo(username string) (res ExecResult) {
	user := TableUser{}
	stmt, err := mysql.DBConn().Prepare(
		"select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = user
	fmt.Println("result:\t", res.Data)
	return
}

func UserExist(username string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"select 1 from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}
