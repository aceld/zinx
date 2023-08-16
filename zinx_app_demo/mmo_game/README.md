# Zinx Application - MMO Game Case Study

English | [简体中文](README-CN.md)

## 1.Introduction to the Application Case

 "In the previous chapters, we have completed some of the core functionalities of the Zinx framework. Now, we are going to build a server-side application based on Zinx. To illustrate this, we will use a game application server as an example of Zinx's practical application."

  The game scenario is an MMO (Massively Multiplayer Online) game with a Unity3D client. The server-side demo implements essential modules of an MMO game, including AOI (Area Of Interest) based broadcasting and global chat.

![mmo-game01](https://github.com/aceld/zinx/assets/7778936/b9576a95-d33f-4ffd-bf89-94a551af1977)

## 2.Server Application Basic Protocols

| MsgID | Client      | Server      | 描述                                                       |
| ----- | ----------- | ----------- | ---------------------------------------------------------- |
| 1     | -           | SyncPid     | Synchronize the player's login ID (used for player ID).|
| 2     | Talk        | -           | World chat.                                                   |
| 3     | MovePackege | -           | Movement.                                                       |
| 200   | -           | BroadCast   | Broadcast messages (Tp 1: world chat, 2: coordinates synchronization (spawn point sync), 3: action, 4: updated coordinates after movement). |
| 201   | -           | SyncPid     | Broadcast messages for disconnect/aoi disappearance.|
| 202   | -           | SyncPlayers | Synchronize the position information of nearby players (including oneself).                               |


## 3. Zinx Documentation

[ < Zinx Wiki : English > ](https://github.com/aceld/zinx/wiki)

[ < Zinx 文档 : 简体中文> ](https://www.yuque.com/aceld/tsgooa/sbvzgczh3hqz8q3l)

## 4.Tutorial 
### 4.1 Online Tutorial

| platform | Entry | 
| ---- | ---- | 
| <img src="https://user-images.githubusercontent.com/7778936/236784004-b6d99e26-b1ab-4bc3-988e-7a46108b85fe.png" width = "100" height = "100" alt="" align=center />| [Zinx Framework tutorial-Lightweight server based on Golang](https://dev.to/aceld/1building-basic-services-with-zinx-framework-296e)| 
|<img src="https://user-images.githubusercontent.com/7778936/236784168-6528a9b8-d37b-4b02-a37c-b9988d7508d8.jpeg" width = "100" height = "100" alt="" align=center />|[《Golang轻量级并发服务器框架zinx》](https://www.yuque.com/aceld)|


### 4.2 Online Tutorial Video

| platform | online video | 
| ---- | ---- | 
| <img src="https://s1.ax1x.com/2022/09/22/xFePUK.png" width = "100" height = "100" alt="" align=center />| [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.bilibili.com/video/av71067087)| 
| <img src="https://s1.ax1x.com/2022/09/22/xFeRVx.png" width = "100" height = "100" alt="" align=center />  | [![zinx-BiliBili](https://s2.ax1x.com/2019/10/13/uv340S.jpg)](https://www.douyin.com/video/6983301202939333891) |
| <img src="https://s1.ax1x.com/2022/09/23/xkQcng.png" width = "100" height = "100" alt="" align=center />| [![zinx-youtube](https://s2.ax1x.com/2019/10/14/KSurCR.jpg)](https://www.youtube.com/watch?v=U95iF-HMWsU&list=PL_GrAPKmuajzeNI8HBTi-k5NQO1g0rM-A)| 


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
