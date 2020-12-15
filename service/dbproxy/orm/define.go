package orm

import "database/sql"

// 文件表结构
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 用户表
type TableUser struct {
	Username   string
	Email      string
	Phone      string
	SignupAt   string
	LastActive string
	Status     int
}

// 用户文件表
type TableUserFile struct {
	UserName   string
	FileHash   string
	FileName   string
	FileSize   int64
	UploadAt   string
	LastUpdate string
}

// sql函数执行的结果
type ExecResult struct {
	Suc  bool        `json:"suc"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
