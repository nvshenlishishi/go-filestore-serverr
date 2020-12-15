## ubuntu下离线安装harbor1.6
- 修改本地dns
- 假设本地ip为192.168.200.212
``` 
# vim /etc/hosts
192.168.200.212 hub.fileserver.com
```
- 下载harbor1.6
``` 
cd /data/apps
wget https://storage.googleapis.com/harbor-releases/harbor-offline-installer-v1.6.1.tgz
tar xvf harbor-offline-installer-v1.6.1.tgz
```
- 修改harbor.cfg
``` 
# vim harbor.cfg, 测试环境下可以只修改以下两个参数
hostname = hub.fileserver.com
harbor_admin_password = <自定义admin的登录密码>
```
- 一键安装harbor
``` 
# 执行./install.sh脚本
./install.sh
```
- ui登录测试
``` 
# 打开浏览器输入地址
http://hub.fileserver.com
初始密码：admin/Harbor12345
```
- 测试docker login
``` 
默认情况下,docker是要去访问https://hub.fileserver.com的．但目前没有配置https;
所以我们可以先改下本地docker的配置:

# vim /etc/docker/daemon.json
{
 "insecure-registries" : [
    "hub.fileserver.com"
  ]
}

然后重启docker: sudo systemctl restart docker
再重试login: sudo docker login hub.fileserver.com

推送镜像测试: 
docker tag fileserver/apigw hub.fileserver.com/filestore/apigw
docker push hub.fileserver.com/filestore/apigw
```
- 后续启动或停止harbor
``` 
cd /data/apps/harbor
# 启动
sudo docker-compose start
# 停止
sudo docker-compose stop
```