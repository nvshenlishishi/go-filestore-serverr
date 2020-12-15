## 11-1 通过kubeadm单机安装k8s

## 11-2 安装k8s(v1.14.1)可视化管理


## 11-3 Docker与Docker-Compose基本概念
### 已有的部署方式
``` 
upload
download
transfer
apigw
account
dbproxy
```
- 多机多实例部署
### docker的几个优点
- 限制容器的cpu及内存等资源消耗
- 应用依赖环境的隔离
- 快速扩容，动态起停容器实例

### docker-compose
- compose是docker容器进行编排的工具
- 默认的模版文件是docker-compose.yml
- 非常适合组合使用多个容器进行开发的场景
- 通过dockerfile定义容器环境，打包成镜像
- 通过docker-compose.yml定义各应用服务
- 通过docker-compose up命令来启动所有容器


## 11-4 基于容器的微服务反向代理利器Traefik
### traefik架构
### traefik两大概念
- frontend 用于控制访问的路由规则，支持单个规则及正则匹配
- backend 用于匹配一组服务实例，通过轮询方式来选择转发的目标

### traefik特点
- 安装简单，无需安装依赖
- 监控后台，自动更新路由配置
- 支持两种均衡模式: 加权轮询，动态轮询
- 前后台均支持https


## 11-5 基于Docker-compose与Traefik的容器化部署演示
``` 
cd deploy/service_dc/

sudo docker-compose up --scale upload=2 --scale download=2 -d


cd deploy/traefik_dc
sudo docker-compose up -d
```
``` 
vi /etc/hosts

192.168.2.244 upload.fileserver.com
192.168.2.244 download.fileserver.com
192.168.2.244 apigw.fileserver.com

```

## 11-6 kubernetes基础原理
### Kubernetes是什么
- k8s是一个分布式系统的支撑平台
- 底层可以基于docker来包装应用
- 以集群的方式来运行/管理跨机器的容器应用
- 解决了docker跨机器场景的容器通讯问题
- 拥有自动修复能力
- 提供部署运行/资源调度/服务发现/动态伸缩等一系列功能

### 为什么选择kubernetes
- k8s天生适用于微服务架构应用
- k8s有强大的横向扩展能力
- 提供来完善的管理工具集
- 使得我们拥有随时迁移整体业务系统的能力

## 11-7 基于kubernetes的容器化部署演示
- service_k8s/batch_deploy.sh
``` 

```