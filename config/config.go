package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Env                  string `json:"env"`                      // 环境
	UploadLBHost         string `json:"upload_lb_host"`           // 上传服务LB地址
	DownloadLBHost       string `json:"download_lb_host"`         // 下载服务LB地址
	TraceAgentHost       string `json:"trace_agent_host"`         // tracing agent地址
	PwdSalt              string `json:"pwd_salt"`                 // 盐值
	CephAccessKey        string `json:"ceph_access_key"`          // Ceph访问Key-------------Ceph
	CephSecretKey        string `json:"ceph_secret_key"`          // Ceph访问密钥
	CephGWEndpoint       string `json:"ceph_gw_endpoint"`         // Ceph Gateway地址
	MysqlUser            string `json:"mysql_user"`               // mysql用户名-------------Mysql
	MysqlPwd             string `json:"mysql_pwd"`                // mysql密码
	MysqlHost            string `json:"mysql_host"`               // mysql ip
	MysqlPort            string `json:"mysql_port"`               // mysql port
	MysqlDb              string `json:"mysql_db"`                 // mysql db
	MysqlCharset         string `json:"mysql_charset"`            // mysql charset
	MysqlMaxConn         int    `json:"mysql_max_conn"`           // mysql max connect
	RedisHost            string `json:"redis_host"`               // redis地址---------------Redis
	RedisPass            string `json:"redis_pass"`               // redis密码
	OSSBucket            string `json:"oss_bucket"`               // oss bucket名------------OSS
	OSSEndpoint          string `json:"oss_endpoint"`             // oss endpoint
	OSSAccessKey         string `json:"oss_access_key"`           // oss 访问key
	OSSAccessSecret      string `json:"oss_access_secret"`        // oss 访问secret
	AsyncTransferEnable  bool   `json:"async_transfer_enable"`    // 是否开启文件异步转移(默认同步)----Rabbit
	RabbitURL            string `json:"rabbit_url"`               // rabbitmq服务的入口url
	TransExchangeName    string `json:"trans_exchange_name"`      // 用户文件transfer的交换机
	TransOSSQueueName    string `json:"trans_oss_queue_name"`     // oss转移队列名
	TransOSSErrQueueName string `json:"trans_oss_err_queue_name"` // oss转移失败后写入另一个队列的队列名
	TransOSSRoutingKey   string `json:"trans_oss_routing_key"`    // routingkey
	TempLocalRootDir     string `json:"temp_local_root_dir"`      // 本地临时存储地址路径
	TempPartRootDir      string `json:"temp_part_root_dir"`       // 分块文件在本地临时存储地址的路径
	CurrentStoreType     int    `json:"current_store_type"`       // 设置当前文件的存储类型
	OSSRootDir           string `json:"oss_root_dir"`             // OSS的存储路径前缀
	CephRootDir          string `json:"ceph_root_dir"`            // Ceph的存储路径前缀
	DownloadEntry        string `json:"download_entry"`           // UploadEntry : 配置上传入口地址
	DownloadServiceHost  string `json:"download_service_host"`    // DownloadServiceHost : 上传服务监听的地址
	UploadEntry          string `json:"upload_entry"`             // UploadEntry : 配置上传入口地址
	UploadServiceHost    string `json:"upload_service_host"`      // 上传服务监听的地址--------Service
	UploadMicroHost      string `json:"upload_micro_host"`        //
	ConsulAddr           string `json:"consul_addr"`              //
}

var DefaultConfig *Configuration

func InitConfig(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println(err.Error())
		panic("无配置文件")
	}
	file, _ := os.Open(filename)
	defer file.Close()

	decoder := json.NewDecoder(file)
	DefaultConfig = &Configuration{}

	err := decoder.Decode(DefaultConfig)
	if err != nil {
		panic(err)
	}
}
