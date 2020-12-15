package orm

import (
	"fmt"
	"go-filestore-server/database/mysql"
	"log"
	"time"
)

func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("" +
		"insert ignore into tbl_user_file (`user_name`,`file_hash`,`file_name`,`file_size`,`upload_at`) value(?,?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

func QueryUserFileMetas(username string, limit int64) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_hash, file_name, file_size, upload_at, last_update from tbl_user_file where user_name=? limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	var userFiles []TableUserFile
	for rows.Next() {
		ufile := TableUserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdate)
		if err != nil {
			log.Println(err.Error())
			break
		}
		fmt.Println(ufile)
		userFiles = append(userFiles, ufile)
	}
	res.Suc = true
	res.Data = userFiles
	return
}

func DeleteUserFile(username, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"update tbl_user_file set status=2 where user_name=? and file_hash=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

func RenameFileName(username, filehash, filename string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"update tbl_user_file set file_name=? where user_name=? and file_hash=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

func QueryUserFileMeta(username string, filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_hash, file_name, file_size, upload_at, last_update from tbl_user_file where user_name=? and file_hash=? limit 1")
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	ufile := TableUserFile{}
	if rows.Next() {
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdate)
		if err != nil {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = ufile
	return
}
