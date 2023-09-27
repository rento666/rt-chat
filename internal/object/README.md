本目录为`object`,是为存放一些结构体的定义、所属函数

目前有：
```text
client：
    介绍：客户端结构体，存放每个websocket连接的一些信息
    数据库：不保存
message：
    介绍：消息结构体，存放用户发送的websocket消息
    数据库：保存，表名为messages
user：
    介绍：用户结构体，存放关于用户的一些信息
    数据库：保存，有user_basic基础信息表和user_auth账户授权表
```