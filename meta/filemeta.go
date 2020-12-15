package meta

import (
	"go-filestore-server/model"
	"sort"
)

type FileMeta struct {
	FileHash string `json:"file_hash"` // 文件Hash
	FileName string `json:"file_name"` // 文件名称
	FileSize int64  `json:"file_size"` // 文件大小
	Location string `json:"location"`  // 文件位置
	UploadAt string `json:"upload_at"` // 上次时间
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileHash] = fmeta
}

// 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return model.OnFileUploadFinished(fmeta.FileHash, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := model.GetFileMeta(fileSha1)
	if tfile == nil || err != nil {
		return nil, err
	}

	fmeta := FileMeta{
		FileHash: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil
}

// 获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// 批量从mysql获取文件元信息
func GetLastFileMetasDB(limit int) ([]FileMeta, error) {
	tfiles, err := model.GetFileMetaList(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}

	tfilesMeta := make([]FileMeta, len(tfiles))
	for i := 0; i < len(tfilesMeta); i++ {
		tfilesMeta[i] = FileMeta{
			FileHash: tfiles[i].FileHash,
			FileName: tfiles[i].FileName.String,
			FileSize: tfiles[i].FileSize.Int64,
			Location: tfiles[i].FileAddr.String,
		}
	}
	return tfilesMeta, nil
}

// 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
