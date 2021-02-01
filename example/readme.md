### 启动connector服务

#### 启动zookeeper(docker-compose.yml)
```yaml
# (默认账号：admin 密码：admin)
version: '3.1'
services:
  zoo1:
    image: zookeeper
    hostname: zoo1
    ports:
      - 2181:2181
    environment:
      ZOO_MY_ID: 1
      ZOO_SERVERS: server.1=0.0.0.0:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo2:
    image: zookeeper
    hostname: zoo2
    ports:
      - 2182:2181
    environment:
      ZOO_MY_ID: 2
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=0.0.0.0:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo3:
    image: zookeeper
    hostname: zoo3
    ports:
      - 2183:2181
    environment:
      ZOO_MY_ID: 3
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=0.0.0.0:2888:3888;2181

  node-zk-browser:
    image: fify/node-zk-browser
    hostname: node-zk-browser
    ports:
      - "3000:3000"
    environment:
      ZK_HOST: zoo1:2181
```

#### 更新包
```
go get -u 
```


#### 启动
```bash
# 启动 connector
cd connector; go run connector/main.go
# 启动 game1
cd game1; go run game1/main.go
```
