package mq

// 将要写到rabbitmq的数据的结构体
type TransferData struct {
	FileHash      string `json:"file_hash"`
	CurLocation   string `json:"cur_location"`
	DestLocation  string `json:"dest_location"`
	DestStoreType int    `json:"dest_store_type"`
}
