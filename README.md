# go-filestore-server
## 单机运行
### 准备工作
- 1.准备mysql数据库，运行sql
- 2.运行
``` 
go run main.go
# 注册
http://localhost:8080/static/view/signup.html
# 登陆
http://localhost:8080/static/view/signin.html
# 主页
http://localhost:8080/static/view/home.html
```
## 微服务运行
### 准备工作
- micro生成
``` 
protoc --proto_path=service/account/proto --go_out=service/account/proto --micro_out=service/account/proto service/account/proto/user.proto
protoc --proto_path=service/dbproxy/proto --go_out=service/dbproxy/proto --micro_out=service/dbproxy/proto service/dbproxy/proto/proxy.proto
protoc --proto_path=service/download/proto --go_out=service/download/proto --micro_out=service/download/proto service/download/proto/download.proto
protoc --proto_path=service/upload/proto --go_out=service/upload/proto --micro_out=service/upload/proto service/upload/proto/upload.proto
```
- 运行
``` 
cd service
./start-all.sh

# 注册
http://localhost:8080/static/view/signup.html
# 登陆
http://localhost:8080/static/view/signin.html
# 主页
http://localhost:8080/static/view/home.html

./stop-all.sh
```