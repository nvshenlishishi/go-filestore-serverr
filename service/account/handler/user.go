package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/model"
	proto "go-filestore-server/service/account/proto"
	dbcli "go-filestore-server/service/dbproxy/client"
	"go-filestore-server/util"
	"time"
)

// User: 用于实现UserServiceHandler接口的对象
type User struct{}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + config.DefaultConfig.PwdSalt))
	return tokenPrefix + ts[:8]
}

func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	username := req.Username
	passwd := req.Password

	if len(username) < 3 || len(passwd) < 5 {
		resp.Code = common.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}

	encPasswd := util.Sha1([]byte(passwd + config.DefaultConfig.PwdSalt))
	suc := model.UserSignup(username, encPasswd)
	if suc {
		resp.Code = common.StatusOK
		resp.Message = "注册成功"
	} else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "注册失败"
	}
	return nil
}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {
	username := req.Username
	password := req.Password

	encPasswd := util.Sha1([]byte(password + config.DefaultConfig.PwdSalt))

	dbResp, err := dbcli.UserSignin(username, encPasswd)
	if err != nil || !dbResp.Suc {
		resp.Code = common.StatusLoginFailed
		return nil
	}

	token := GenToken(username)
	upRes, err := dbcli.UpdateToken(username, token)
	if err != nil || !upRes.Suc {
		resp.Code = common.StatusServerError
		return nil
	}

	resp.Code = common.StatusOK
	resp.Token = token
	return nil
}

func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	dbResp, err := dbcli.GetUserInfo(req.Username)
	if err != nil {
		resp.Code = common.StatusServerError
		resp.Message = "服务错误"
		return nil
	}
	// 查不到对应的用户信息
	if !dbResp.Suc {
		resp.Code = common.StatusUserNotExists
		resp.Message = "用户不存在"
		return nil
	}
	fmt.Println("resp:\t", dbResp.Data)
	user := dbcli.ToTableUser(dbResp.Data)

	resp.Code = common.StatusOK
	resp.Username = user.Username
	resp.SignupAt = user.SignupAt
	resp.LastActiveAt = user.LastActive
	resp.Status = int32(user.Status)
	resp.Email = user.Email
	resp.Phone = user.Phone
	return nil
}

// UserFiles : 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}
	fmt.Println(dbResp.Data)
	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	fmt.Println(userFiles)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

// UserFiles : 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}
