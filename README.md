# go-websocket
Golang分布式消息推送实现方案
1、每台接口服务器独立管理自己的websocket连接池.
2、管理机通过Redis发布订阅模式将消息发布到订阅频道的接口机.
3、接口机接收消息，发送给websocket连接的设备.
