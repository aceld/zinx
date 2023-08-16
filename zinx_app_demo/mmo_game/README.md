# Zinx应用-MMO游戏案例

## 一、应用案例介绍

​	“ 好了，以上Zinx的框架的一些核心功能我们已经完成了，那么接下来我们就要基于Zinx完成一个服务端的应用程序了，整理用一个游戏应用服务器作为Zinx的一个应用案例。”

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


## 三、Zinx详细教程

[ < Zinx Wiki : English > ](https://github.com/aceld/zinx/wiki)

[ < Zinx 文档 : 简体中文> ](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)

---
### 关于作者：

作者：`Aceld(刘丹冰)`
简书号：`IT无崖子`

`mail`:
[danbing.at@gmail.com](mailto:danbing.at@gmail.com)
`github`:
[https://github.com/aceld](https://github.com/aceld)
`原创书籍gitbook`:
[http://legacy.gitbook.com/@aceld](http://legacy.gitbook.com/@aceld)

### Zinx技术讨论社区

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236775008-6bd488e3-249a-4d43-8885-7e3889e11e2d.png" width = "100" height = "100" alt="" align=center />| https://discord.gg/xQ8Xxfyfcz| 
| <img src="https://user-images.githubusercontent.com/7778936/236775137-5381f8a6-f534-49c4-8628-e52bf245c3bc.jpeg" width = "100" height = "100" alt="" align=center />  | 加微信: `ace_ld`  或扫二维码，备注`zinx`即可。</br><img src="https://user-images.githubusercontent.com/7778936/236781258-2f0371bd-5797-49e8-a74c-680e9f15843d.png" width = "150" height = "150" alt="" align=center /> |
|<img src="https://user-images.githubusercontent.com/7778936/236778547-9cdadfb6-0f62-48ac-851a-b940389038d0.jpeg" width = "100" height = "100" alt="" align=center />|<img src="https://s1.ax1x.com/2020/07/07/UFyUdx.th.jpg" height = "150"  alt="" align=center /> **WeChat Public Account** |
|<img src="https://user-images.githubusercontent.com/7778936/236779000-70f16c8f-0eec-4b5f-9faa-e1d5229a43e0.png" width = "100" height = "100" alt="" align=center />|<img src="https://s1.ax1x.com/2020/07/07/UF6Y9S.th.png" width = "150" height = "150" alt="" align=center /> **QQ Group** |
