package common

const (
	StandardTimeFormat = "2006-01-02 15:04:05" // 标准格式
)

// 存储类型（表示文件存到哪里)
const (
	_          int = iota // 0.
	StoreLocal            // 1.节点本地
	StoreCeph             // 2.Ceph集群
	StoreOSS              // 3.阿里OSS
	StoreMix              // 4.混合(Ceph+OSS)
	StoreAll              // 5.所有类型的存储都存储一份数据
)
