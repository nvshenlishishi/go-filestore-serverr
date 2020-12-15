# centos安装docker
## 1.yum配置阿里云的源
```
cd /etc/yum.repos.d/

#下载阿里云yum源
wget http://mirrors.aliyun.com/repo/Centos-7.repo
mv CentOS-Base.repo CentOS-Base.repo.bak
mv Centos-7.repo CentOS-Base.repo
```
## 2.重置yum源
``` 
yum clean all
yum makecache
```
## 3.开始安装docker
``` 
yum list docker-ce
yum -y install docker-ce
docker -v
systemctl start docker
docker info
```

### 其他常用操作
- 拉取镜像
``` 
// 下载官方nginx镜像源，等同于docker pull docker.io/nginx
// 也等同于 docker pull docker.io/nginx:lastest
docker pull nginx

// 下载国内镜像源的ubuntu镜像，并指定版本为18.04
docker pull registry.docker-cn.com/library/ubuntu:18.04
```
- 推送镜像
``` 
// 推送镜像到 docker hub, 需要先注册账户
docker push <你的用户名>/<你打包时定义的镜像名>:<标签，版本号>

// 推送到你的一个私有镜像仓库，需要提前搭建好仓库服务(比如用harbor来搭建)
docker push <私有镜像库域名, 如a.b.com>/<项目名称>/镜像名:<标签>
```
- 打包镜像
``` 
// 提取准备好一个Dockerfile,在Dockerfile相同路径下执行
docker build -t <指定一个完整的镜像名, 比如testsvr:v1.0> .
// 即可打包出一个本地镜像， 然后再通过docker push就可以推送到远端镜像仓库
```
- 启动容器
``` 
docker run -d               // -d 表示通过daemon方式来启动
-p 13306:3306               // 端口映射，将host主机的13306端口和docker容器的3306端口映射起来
-v /etc/mysql:/var/mysql    // 目录挂载，将容器内的/var/mysql目录挂载到host主机到/etc/mysql目录，可以实现容器内这个目录下的数据持久化
mysql                       // 镜像名，指定加载哪个镜像   
```
- 重启或停止或删除容器应用
``` 
docker ps                // 列出目前正在运行的容器列表
docker ps -a             // 列出所有的容器列表
docker start <容器id>     // 通过容器id来重启某个容器， 批量操作的话，直接在参数后面再跟对应容器id即可
docker stop <容器id>      // 通过容器id来关闭某个容器， 批量操作的话，直接在参数后面再跟对应容器id即可
docker rm <容器id>        // 通过容器id来删掉某个已经停止掉的容器
docker rm -f <容器id>     // 通过容器id来删掉某个正在运行的容器
```
- 删除本地镜像
``` 
docker rmi  <镜像id>
docker rmi -f <镜像id> // 强制删除
```
- 查看容器日志
``` 
docker logs -f <容器id>
docker inspect <容器id> // 从返回结果中找到LogPath,运行的历史日志会在这个文件里找到
```
- 进入容器内
``` 
docker exec -it <容器id> /bin/bash        // 进入容器内并进入它的shell终端
docker exec -it <容器id> <shell命令>       // 在容器内执行shell命令，比如
docker exce -it <容器id> ls -l            // 查看容器内系统根目录下所有文件或文件夹
## 进入容器后，可以直接通过exit命令退出容器
```