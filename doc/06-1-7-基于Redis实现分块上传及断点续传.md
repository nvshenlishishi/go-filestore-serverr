## 6-1 分块上传与断点续传原理
### 分块上传与断点续传
- 两个概念
``` 
- 分块上传: 文件切成多块，独立传输，上传完成后合并
- 断点续传: 传输暂停或异常中断后，可基于原来进度重传

```
- 几点说明
```
- 小文件不建议分块上传
- 可以并行上传分块，并且可以无序传输
- 分块上传能极大提高传输效率
- 减少传输失败后重试的流量及时间
```
- 具体流程
``` 
1. 初始化上传        Initiate Multipart Upload
2. 上传分块（并行)    Upload Part -> Upload Abort
3. 通知上传完成       Complete Multipart Upload  Upload Query
```
### 服务架构变迁
- 用户->分块上传
- 分块上传<->本地存储
- 分块上传->Redis缓存|Hash计算
- 分块上传->用户文件表
- 分块上传->唯一文件表

## 6-2 编码实战: Go实现Redis连接池（存储分块信息）
### 分块上传通用接口
``` 
- InitiateMultipartUploadHandler    初始化分块信息
- UploadPartHandler                 上传分块
- CompleteUploadPartHandler         通知分块上传完成
- CancelUploadPartHandler           取消上传分块
- MultipartUploadStatusHandler      查看分块上传的整体状态
```

### 接口:上传初始化
- 判断是否已经上传过
- 生成唯一上传ID
- 缓存分块初始化信息

### redis操作
``` 
redis-cli
auth testupload

keys *
quit

```
### redis连接池
- newRedisPool

## 6-3 编码实战: 实现初始化分块上传接口
### InitialMultipartUploadHandler
- 1.解析用户请求参数
- 2.获得redis的一个连接
- 3.生成分块上传的初始化信息
- 4.将初始化信息写入到redis缓存
- 5.将响应初始化数据返回到客户端


## 6-4 编码实战: 实现分块上传接口
### UploadPartHandler
- 1.解析用户请求参数
- 2.获得redis的一个连接
- 3.获得文件句柄，用于存储分块内容
- 4.更新redis缓存状态
- 5.返回处理结果到客户端

## 6-5 编码实战: 实现分块合并接口
### CompleteUploadHandler
- 1.解析用户请求参数
- 2.获得redis的一个连接
- 3.通过uploadid查询redis并判断是否所有分块上传完成
- 4.TODO: 合并分块
- 5.更新唯一文件表及用户文件表
- 6.响应处理结果

## 6-6 分块上传场景测试+小结
- test_mpupload.go
```go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) // 每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		// 遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}

	return nil
}

func main() {
	username := "admin"
	token := "54eefa7dbd5bcf852c52fecd816f2a315c61832c"
	filehash := "dfa39cac093a7a9c94d25130671ec474d51a2995"

	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		"http://localhost:8080/file/mpupload/init",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"132489256"},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	filename := "/data/pkg/go1.10.3.linux-amd64.tar.gz"
	tURL := "http://localhost:8080/file/mpupload/uppart?" +
		"username=admin&token=" + token + "&uploadid=" + uploadID
	multipartUpload(filename, tURL, chunkSize)

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		"http://localhost:8080/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"132489256"},
			"filename": {"go1.10.3.linux-amd64.tar.gz"},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}
```
- go run main.go
- go run test_mpupload.go

- 手动合并验证
```
cat `ls | sort -n` > /tmp/a
shalsum /tmp/a
```
### 上传取消
- 删除已存在的分块文件
- 删除redis缓存状态
- 更新mysql文件status

### 上传状态查询
- 检查分块上传状态是否有效
- 获取分块初始化信息
- 获取已上传的分块信息

### 本章小结
- 1.分块上传与断点续传的概念
- 2.分块上传流程的讲解
- 3.几个重要接口的逻辑实现与演示

## 6-7 文件断点下载原理
### 断点续传下载的效果
- 客户端能够实现分段下载，中断传输后只要记住上次下载的位置，就能够续传而不需要重传
- 客户端可以实现进度条展示，能够手动暂停传输和继续传输

### 断点续传下载相关的几个HTTP头
- Accept-Ranges: bytes
``` 
服务端响应的header,用于告诉客户端我支持断点续传，你可以指定Range来下载文件的某一部分
```
- Range: 100- 或 100-1000, 自定义
```
客户端请求的header, 用于告诉服务端我想下载文件的哪一部分内容
100-表示下载100字节之后的文件内容
100-1000 表示下载offset为100-1000以内的这一段文件内容
如果不指定Range,默认是希望下载整个文件内容
```
- Content-Range: bytes=0-500/1000
``` 
服务端响应的header, 用于告诉客户端我返回的文件内容区间是多少
0-500是指文件的前面500个字节，而整个文件大小为1000
```
- Content-Length: 500
``` 
服务端响应的header, 用于告诉客户端我返回的内容长度是多少
500表示当前总共返回来500个字节的内容
```
- Accept: 比如image/gif, image/jpeg 或 */*
``` 
客户端请求的header,用于告诉服务端我可以接受的响应内容（文件）类型，比如image/gif,image/jpeg,*/*表示我什么类型都接受
```
- Last-Modified: 
``` 
服务端响应的header,非必须，用于告诉客户端这个文件资源最后一次的修改时间
如果客户端在下载文件的过程中，资源被修改来
可以通过Last-Modified来感知，从而客户端可能要考虑重新下载
```