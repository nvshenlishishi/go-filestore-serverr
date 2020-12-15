package orm

import (
	"database/sql"
	"go-filestore-server/database/mysql"
	"log"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file (`file_hash`,`file_name`,`file_size`,`file_addr`,`status`) value (?,?,?,?,1)")
	if err != nil {
		log.Println("failed to prepare statement, err:\t", err.Error())
		res.Suc = false
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		return
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Println("file with hash has been uploaded before:\t", filehash)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

func GetFileMeta(filehash string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_hash, file_addr, file_name,file_size from tbl_file where file_hash=? and status=1 limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Suc = true
			res.Data = nil
			return
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = tfile
	return
}

func GetFileMetaList(limit int64) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare("select file_hash, file_addr, file_name, file_size from tbl_file where status =1 limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
		if err != nil {
			log.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	res.Suc = true
	res.Data = tfiles
	return
}

func UpdateFileLocation(filehash string, fileaddr string) (res ExecResult) {
	stmt, err := mysql.DBConn().Prepare(
		"udpate tbl_file set `file_addr`=? where `file_sha1`=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			log.Println("更新记录:\t", filehash)
			res.Suc = false
			res.Msg = "无更新记录"
			return
		}
		res.Suc = true
		return
	} else {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
}
