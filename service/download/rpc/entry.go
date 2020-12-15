package rpc

import (
	"context"
	"go-filestore-server/config"
	dlProto "go-filestore-server/service/download/proto"
)

// Dwonload :download结构体
type Download struct{}

// DownloadEntry : 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqEntry,
	res *dlProto.RespEntry) error {

	res.Entry = config.DefaultConfig.DownloadEntry
	return nil
}
