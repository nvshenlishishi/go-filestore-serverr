# docker 安装mysql主从
- 前提是安装好docker（MAC环境）
## 0.准备工作
```
# 0.提取创建好mysql/data目录
cd 
mkdir -p docker/mysql/data_master/
mkdir -p docker/mysql/data_slave/
mkdir -p docker/mysql/data/conf/
cd docker/mysql/data/conf/
# Mac 路径 /Users/xx/docker/mysql/data/conf

创建master.conf+slave.conf
```
- master.conf
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
- slave.conf
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
## 1.拉取Mysql5.7镜像到本地
```
docker pull mysql:5.7
```
## 2.运行master节点和slave节点
``` 
# 自己改路径
# master
docker run -d --name mysql-master -p 13306:3306 -v /Users/xx/docker/mysql/data/conf/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /Users/xx/docker/mysql/data_master:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
# slave
docker run -d --name mysql-slave  -p 13307:3306 -v /Users/xx/docker/mysql/data/conf/slave.conf:/etc/mysql/mysql.conf.d/mysqld.cnf  -v /Users/xx/docker/mysql/data_slave:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
```
## 3.登陆主从节点修改配置
``` 
# ifconfig 查询本机ip （192.168.xx.xx 是你本机的内网ip）

# 登陆主节点
mysql -u root -h 192.168.xx.xx -P13306 -p123456
mysql -u root -h 127.0.0.1 -P13306 -p123456
mysql > grant replication slave on *.* to 'slave'@'%' identified by 'slave';
mysql > flush privileges;
mysql > show master status \G;
# 注意对应的master_log_file名称


# 登陆从节点
mysql -u root -h 192.168.xx.xx -P13307 -p123456
mysql > stop slave;
mysql > change master to master_host='192.168.xx.xx',master_port=13306,master_user='slave',master_password='slave',master_log_file='log.00000x',master_log_pos=0;
mysql > start slave;

# 验证
## 主节点创建表验证
mysql > create database fileserver default character set utf8mb4;
```

## 4.运行Mysql和操作
```
mysql -u root -h 192.168.xx.xx -P13306 -p123456
mysql -u root -h 192.168.xx.xx -P13307 -p123456

# 停止重启Docker
docker stop mysql-slave
docker stop mysql-master

docker start mysql-master
docker start mysql-slave
```