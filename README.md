# floppyisadog
一，安装docker
1，下载安装
	docker --version
	Docker version 20.10.21, build baeda1f
2，配置腾讯云加速镜像：https://mirror.ccs.tencentyun.com

二，创建docker虚拟网络
目标：创建一个名为tars的桥接(bridge)虚拟网络，网关172.25.0.1，网段为172.25.0.0/16
docker network create -d bridge --subnet=172.25.0.0/16 --gateway=172.25.0.1 tars

docker network create -d bridge --subnet=172.25.0.0/16 --gateway=172.25.0.1 tars
7ee05446463adc398eb6627facd0cac1e6a5b0761f8b76087d2eb7772291a756




三，启动mysql服务
目标：版本5.6; 名称：tars-mysql; 用户：root； 密码：123456； ip：172.25.0.2；port：3306:3306
docker run -d -p 3306:3306 \
    --net=tars \
    -e MYSQL_ROOT_PASSWORD="123456" \
    --ip="172.25.0.2" \
    -v ~/framework-mysql:/var/lib/mysql \
    --name=tars-mysql \
    --platform linux/amd64 \
    mysql:5.6

问题：
1，platform问题
2，mount目录问题


四，部署框架服务
1，docker pull tarscloud/framework:latest

2，启动最新框架
docker run -d \
    --name=tars-framework \
    --platform linux/amd64 \
    --net=tars \
    -e MYSQL_HOST="172.25.0.2" \
    -e MYSQL_ROOT_PASSWORD="123456" \
    -e MYSQL_USER=root \
    -e MYSQL_PORT=3306 \
    -e REBUILD=false \
    -e INET=eth0 \
    -e SLAVE=false \
    --ip="172.25.0.3" \
    -v ~/framework:/data/tars \
    -p 3000:3000 \
    -p 3001:3001 \
    tarscloud/framework


Token：

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhZG1pbiIsImlhdCI6MTY2OTUyNjk2NCwiZXhwIjoxNzY0MzA3NzY0fQ.hMIOD7yWf42GaPyqWHs8TxC2NJJLGHoP5gvAFEN8Frg


GW扩展：
9000 proxy
9001 flow control
9002 web



五，部署应用阶段
1，拉取镜像
docker pull tarscloud/tars-node:latest

2,启动镜像
docker run -d \
    --name=tars-node \
    --platform linux/amd64 \
    --net=tars \
    -e INET=eth0 \
    -e WEB_HOST="http://172.25.0.3:3000" \
    --ip="172.25.0.5" \
    -v ~/tars:/data/tars \
    -p 9000-9030:9000-9030 \
    tarscloud/tars-node:latest



 make tar GOOS=linux
 一，创建schema
1，./migratetool migrate


docker run -d \
    --name=myaccount-spa \
    --platform linux/amd64 \
    --net=tars \
    --ip="172.25.0.8" \
    myaccount-spa:1.0

docker build -t myaccount-spa:1.0 .

http://mirrors.cloud.tencent.com/docker-ce/
"https://mirror.ccs.tencentyun.com"


docker run -d \
    --name=app-spa \
    --platform linux/amd64 \
    --net=tars \
    --ip="172.25.0.9" \
    app-spa:1.0

