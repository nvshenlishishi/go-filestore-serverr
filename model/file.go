package model

import (
	"database/sql"
	"go-filestore-server/database/mysql"
	"go-filestore-server/logger"
)

// TableFile: 文件表结构体
type TableFile struct {
	FileHash string         `json:"file_hash"`
	FileName sql.NullString `json:"file_name"`
	FileSize sql.NullInt64  `json:"file_size"`
	FileAddr sql.NullString `json:"file_addr"`
}

// 插入:  文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file (`file_hash`,`file_name`,`file_size`,`file_addr`,`status`) value (?,?,?,?,1)")
	if err != nil {
		logger.Info("failed to prepare statement, err:\t", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			logger.Info("file with hash has been uploaded before:\t", filehash)
		}
		return true
	}
	return false
}

// 查询:  获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_hash, file_addr, file_name, file_size from tbl_file where file_hash = ? and status = 1 limit 1")
	if err != nil {
		logger.Info(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录，返回参数及错误均为nil
			return nil, nil
		} else {
			logger.Info(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

// 查询:  批量获取文件元信息
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_hash, file_addr, file_name, file_size from tbl_file where status=1 limit ?")
	if err != nil {
		logger.Info(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		logger.Info(err.Error())
		return nil, err
	}

	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tFiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tFile := TableFile{}
		err = rows.Scan(&tFile.FileHash, &tFile.FileAddr, &tFile.FileName, &tFile.FileSize)
		if err != nil {
			logger.Info(err.Error())
			break
		}
		tFiles = append(tFiles, tFile)
	}
	return tFiles, nil
}

// 更新: 更新文件的存储地址
func UpdateFileLocation(filehash string, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set `file_addr` = ? where `file_hash` = ? limit 1")
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		logger.Info(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			logger.Info("更新文件addr 失败, fileHash:\t", filehash)
		}
		return true
	}
	return false
}

// 文件是否已经上传过
func IsFileUploaded(filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("select 1 from tbl_file where filehash=? and status=1 limit 1")
	// TODO 测试中文输入，完成查询逻辑
	rows, err := stmt.Query(filehash)
	if err != nil {
		return false
	} else if rows == nil || !rows.Next() {
		return false
	}
	return true
}

// 文件删除
func OnFileRemoved(filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set status=2 where filehash=? and status=1 limit 1")
	if err != nil {
		logger.Info("failed to prepare statement, err:", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash)
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			logger.Info("file with hash not uploaded:\t", filehash)
		}
		return true
	}
	return false
}
