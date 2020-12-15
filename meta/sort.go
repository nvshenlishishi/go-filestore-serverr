package meta

import (
	"go-filestore-server/common"
	"time"
)

type ByUploadTime []FileMeta

func (b ByUploadTime) Len() int {
	return len(b)
}

func (b ByUploadTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(common.StandardTimeFormat, b[i].UploadAt)
	jTime, _ := time.Parse(common.StandardTimeFormat, b[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}
