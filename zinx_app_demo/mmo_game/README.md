# Zinx应用-MMO游戏案例

## 一、应用案例介绍

​	好了，以上Zinx的框架的一些核心功能我们已经完成了，那么接下来我们就要基于Zinx完成一个服务端的应用程序了，整理用一个游戏应用服务器作为Zinx的一个应用案例。

​	游戏场景是一款MMO大型多人在线游戏，带unity3d 客户端的服务器端demo，该demo实现了mmo游戏的基础模块aoi(基于兴趣范围的广播), 世界聊天等。

![13-Zinx游戏-示例图.png](https://upload-images.jianshu.io/upload_images/11093205-593bb6246327e900.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

## 二、服务器应用基础协议

| MsgID | Client      | Server      | 描述                                                         |
| ----- | ----------- | ----------- | ------------------------------------------------------------ |
| 1     | -           | SyncPid     | 同步玩家本次登录的ID(用来标识玩家)                           |
| 2     | Talk        | -           | 世界聊天                                                     |
| 3     | MovePackege | -           | 移动                                                         |
| 200   | -           | BroadCast   | 广播消息(Tp 1 世界聊天 2 坐标(出生点同步) 3 动作 4 移动之后坐标信息更新) |
| 201   | -           | SyncPid     | 广播消息 掉线/aoi消失在视野                                  |
| 202   | -           | SyncPlayers | 同步周围的人位置信息(包括自己)                               |


## 三、Zinx详细教程(代码教程同步更新)
[《Zinx框架教程-基于Golang的轻量级并发服务器》](https://www.jianshu.com/p/23d07c0a28e5)

---
###关于作者：

作者：`Aceld(刘丹冰)`
简书号：`IT无崖子`

`mail`:
[danbing.at@gmail.com](mailto:danbing.at@gmail.com)
`github`:
[https://github.com/aceld](https://github.com/aceld)
`原创书籍gitbook`:
[http://legacy.gitbook.com/@aceld](http://legacy.gitbook.com/@aceld)

###Zinx技术讨论社区

QQ技术讨论群:

![gopool5.jpeg](https://upload-images.jianshu.io/upload_images/11093205-6cdfd381e8ffa127.jpeg?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

欢迎大家加入，获取更多相关学习资料
