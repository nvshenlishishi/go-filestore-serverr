## 3-1 MYSQL基础知识
### 为什么选Mysql
- Sql数据库与非SQL数据库
```
SQLServer Oracle MYSQL      MongoDB HBase Redis 
```
- Mysql的特点与优劣点
```
小
稳定
社区活跃
多平台
满足一般场景
缺高级功能 
```
- Mysql的适应场景
```
存储关系型数据 很多。。
```
### 服务架构变迁
- 用户-> 上传server -> 用户表|文件表|本地存储

### Mysql安装配置
- 安装模式
```
单点模式
主从模式
    Master DB->写日志->Bin Log
    Slave IO线程读取Bin log 
    Slave IO线程->写日志->Relay Log
    Slave RelayLog -> 读日志 -> SQL线程 -> 回放
多主模式
```
## 3-2 MYSQL主从数据同步演示
```
sudo docker ps 
sudo netstat -antup | grep docker
# 3306 3307

# Slave
mysql -uroot -h127.0.0.1 -P3307 -p

# Master
mysql -uroot -h127.0.0.1 -p
show master status;

# Slave
change master to master_host='192.168.2.244',master_user='reader','master_password='reader',master_log_file='binlog.000002',master_log_pos=0;
start slave;
show slave status\G;
# 查看Slave_IO_Running + Slave_SQL_Running = Yes

# Master
create database test1 default character set utf8;
show databases;

# Slave
show databases;

# Master 
create table tbl_test(`user` varchar(64) not null, `age` int(11) not null) default charset utf8;
show tables;

# Slave
show tables;

# Master 
insert into tbl_test(user, age) values('xiaoming',18);

# Slave 
select * from tbl_test;
```
## 3-3 文件表的设计及创建
```sql
create database fileserver default character set utf8;

use fileserver;
# 创建文件表
CREATE TABLE `tbl_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_hash` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime default NOW() COMMENT '创建日期',
  `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_hash`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

show create tbl_file;
```

### Mysql 分库分表
- 水平分表
```
假设分为256张文件表
按文件hash值后两位来切分
则以: tbl_${file_hash}[:-2]的规则到对应表进行存取
```
### Golang操作Mysql
- 访问Mysql
```
使用Go标准接口
增删改查操作 
```
## 3-4 编码实战: "云存储"系统之持久化元数据到文件表
- github.com/go-sql-driver/mysql
- OnFileUploadFinished
## 3-5 编码实战: "云存储"系统之从文件表中获取元数据
- GetFileMeta

## 3-6 Docker入门基础文档
### Linux下安装Docker社区版（以centos7为例）
- yum源使用阿里云的源
```
cd /etc/yum.repos.d/
# 下载阿里云yum源
wget http://mirrors.aliyun.com/repo/Centos-7.repo
mv CentOS-Base.repo CentOS-Base.repo.bak
mv Centos-7.repo CentOS-Base.repo
```
- 重置yum源
``` 
yum clean all
yum makecache
```
- 安装docker
``` 
# 查看阿里云上docker源信息
yum list docker-ce
# 安装docker最新社区版本
yum -y install docker-ce
# 查看docker版本
docker -v
# 启动docker
systemctl start docker
# 查看docker详细状态信息
docker info
```
### 常用操作
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

## 3-7 Ubuntu中通过Docker安装配置Mysql主从节点
- 以下docker相关的命令，需要在root用户环境下或通过sudo提升权限来进行操作。
### 1.拉取Mysql5.7镜像到本地
``` 
docker pull mysql:5.7

# 如果你只需要跑一个mysql实例，不做主从，那么执行以下命令即可，不用再做后面的参考步骤:
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

# 然后用Shell或客户端软件通过配置(用户名: root 密码: 123456 IP:你的本机IP 端口:3306)来登陆
```
### 2. 准备MYSQL配置文件
- mysql5.7安装后的默认配置文件在/etc/mysql/my.cnf
- 自定义的配置文件一般在/etc/mysql/conf.d路径下
- 创建/data/mysql/conf/master.conf 和 /data/mysql/conf/slave.conf 文件，用于配置主从
- /data/mysql/conf/master.conf
``` 
[client]
default-character-set=utf8
[mysql]
default-character-set=utf8
[mysqld]
log_bin = log # 开启二进制日志，用于从节点的历史复制回放
collation-server = utf8_unicode_ci
init-connect='SET NAMES utf8'
character-set-server = utf8
server_id = 1 # 需保证主库和从库的server_id不同，假设主库设为1
replicate-do-db=fileserver # 需要复制的数据库名，需复制多个数据库的话则重复设置这个选项
```
- /data/mysql/conf/slave.conf
``` 
[client]
default-character-set=utf8
[mysql]
default-character-set=utf8
[mysqld]
log_bin = log # 开启二进制日志，用于从节点的历史复制回放
collation-server = utf8_unicode_ci
init-connect='SET NAMES utf8'
character-set-server = utf8
server_id = 2 # 需保证主库和从库的server_id不同，假设从库设为2
replicate-do-db=fileserver # 需要复制的数据库名，需复制多个数据库的话则重复设置这个选项
```
### 3.Docker分别运行Mysql 主/从两个容器
- 将mysql主节点运行起来
``` 
mkdir -p /data/mysql/data_master
docker run -d --name mysql-master -p 13306:3306 -v /data/mysql/conf/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /data/mysql/data_master:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```
- 运行参数说明
``` 
docker run -d 
--name mysql-master                 // 容器的名称设为mysql-master
-p 13306:3306                       // 将host的13306端口映射到容器的3306端口
-v /data/mysql/conf/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf  // master.conf配置文件挂载
-v /data/mysql/data_master:/var/lib/mysql                           // mysql容器内数据挂载到host的/data/mysql/data_master,用于持久化
-e MYSQL_ROOT_PASSWORD=123456 mysql:5.7                             // mysql的root登陆密码为123456
```
- 将mysql从节点运行起来
``` 
mkdir -p /data/mysql/data_slave
docker run -d --name mysql-slave -p 13307:3306 -v /data/mysql/conf/slave.conf:/etc/mysql/mysq.conf.d/mysqld.conf -v /data/mysql/data_slave:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```
### 4.登陆MYSQL主节点配置同步信息
- 登陆mysql 
``` 
# 192.168.1.xx 是你本机的内网ip
mysql -u root -h 192.168.1.xx -P13306 -p123456
```
- 在mysql client中执行
``` 
mysql > grant replication slave on *.* to 'slave'@'%' identified by 'slave';
mysql > flush privileges;
mysql > create database fileserver default character set utf8mb4;
```
- 获取status，得到类似如下的输出:
```
mysql > show master status \G;
```
### 5.登陆MYSQL从节点配置同步信息
- 登陆mysql 
``` 
docker inspect --format='{{.NetworkSettings.IPAddress}}' mysql-master

# 192.168.1.xx 是你本机的内网ip
mysql -u root -h 192.168.1.xx -P13307 -p123456
```
- 在mysql client中执行
``` 
mysql > stop slave;
# 这个创建表可忽略，master_log_pos=0,  如果pos是最新的mysql master pos, 请自动创建表
# mysql > create database fileserver default character set utf8mb4;
mysql > change master to master_host='192.168.1.xx',master_port=13306,master_user='slave',master_password='slave',master_log_file='log.000000',master_log_pos=0;
mysql > start slave;
```
- 获取status，得到类似如下的输出:
```
mysql > show slave status \G;
```
- 可以尝试在主mysql的fileserver数据库里建表操作下，然后在从节点上检查数据是否已经同步过来。
## 3-8 本章小结
### 使用Mysql小结
- 使用sql.DB来管理数据库连接对象
- 通过sql.Open来创建协程安全的sql.DB对象
- 优先使用Prepared Statement

### 本章小结
- Mysql特点与应用场景
- 主从架构与文件表设计逻辑
- Golang与Mysql的亲密接触