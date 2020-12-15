package rpc

import (
	"context"
	"go-filestore-server/config"
	upProto "go-filestore-server/service/upload/proto"
)

// Upload : upload结构体
type Upload struct{}

// UploadEntry : 获取上传入口
func (u *Upload) UploadEntry(
	ctx context.Context,
	req *upProto.ReqEntry,
	res *upProto.RespEntry) error {

	res.Entry = config.DefaultConfig.UploadEntry
	return nil
}
