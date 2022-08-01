<a name="ibxlh"></a>
## **一、背景及目的**
- 让团队更好地了解新人对技能的掌握情况
- 熟悉简单的Web API后台架构
- 熟悉使用Go实现HTTP API（JSON、文件）
- 熟悉使用Go实现基于TCP的RPC框架（设计和实现通信协议）
- 熟悉基于Auth Token的鉴权机制和流程
- 熟悉使用Go对MySQL、Redis进行基本操作
- 对任务进度和时间有所意识
- 对代码规范、测试、文档、性能调优需要有所意识

<a name="EpNKW"></a>
## **二、逻辑架构设计**
<a name="vvESq"></a>
### 系统架构图
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1658470040023-de40819d-5d6a-4ec6-ab0e-bdb1c83f1099.jpeg)
<a name="ORSG3"></a>
### 页面转换
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1658893344577-8f986f11-a1a9-431c-ba27-fbfa2d877669.jpeg)
<a name="BCBXT"></a>
### 目录结构
```shell
.
├── README.md
├── benchmark
│   ├── getUserInfo.lua
│   ├── signIn.lua
│   └── signUp.lua
├── cmd
│   ├── httpServer
│   │   └── httpServer.go
│   └── tcpServer
│       └── tcpServer.go
├── config
│   └── config.go
├── constant
│   └── errCodeEnum.go
├── database
│   ├── MySqlPool.go
│   └── RedisPool.go
├── go.mod
├── go.sum
├── img
├── internal
│   ├── controller
│   │   └── userController.go
│   ├── entity
│   │   ├── httpRequest.go
│   │   ├── httpResponse.go
│   │   ├── userDO.go
│   │   ├── userDTO.go
│   │   └── userVO.go
│   ├── manager
│   │   ├── sessionManager.go
│   │   ├── userManager.go
│   │   └── userManager_test.go
│   ├── mapper
│   │   ├── userMapper.go
│   │   └── userMapper_test.go
│   └── service
│       ├── userService.go
│       └── userService_test.go
├── logger
│   ├── log.txt
│   └── logger.go
├── rpc
│   ├── client
│   │   ├── rpcClient.go
│   │   └── rpcClient_test.go
│   ├── codec
│   │   └── rpcCodec.go
│   ├── network
│   │   └── transport.go
│   ├── rpcEntity
│   │   ├── rpcRequest.go
│   │   └── rpcResponse.go
│   ├── server
│   │   └── rpcServer.go
│   └── service
│       └── rpcService.go
└── web
    ├── signIn.html
    ├── signUp.html
    └── userInfo.html

```
<a name="r57di"></a>
## **三、核心逻辑详细设计**
<a name="iJrpz"></a>
### 注册流程
![](https://cdn.nlark.com/yuque/__puml/5b4d215a31e517a5df3d9445a36d0b04.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmn6XnnIvnlKjmiLfmmK_lkKblrZjlnKhcbmFjdGl2YXRlIG15c3FsXG5cbm15c3FsIC0-IHNlcnZpY2U6IOi_lOWbnueUqOaIt-S_oeaBr1xucm5vdGUgb3ZlciBzZXJ2aWNlXG7lr4bnoIHliqDlr4ZcbmVuZHJub3RlXG5cbnNlcnZpY2UgLT4gbXlzcWw6IOWtmOWCqOeUqOaIt-S_oeaBr1xubXlzcWwgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIG15c3FsXG5cbnNlcnZpY2UgLT4gY29udHJvbGxlcjogcnBjIHJlc3BvbnNlXG5kZWFjdGl2YXRlIHNlcnZpY2Vcblxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lpITnkIbnu5PmnpxcbmVuZHJub3RlXG5cbmNvbnRyb2xsZXIgLT4gd2ViOiBodHRwIHJlc3BvbnNlXG5cblxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC81YjRkMjE1YTMxZTUxN2E1ZGYzZDk0NDVhMzZkMGIwNC5zdmciLCJpZCI6IktLOUtSIiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="HQIEL"></a>
### 登录流程
![](https://cdn.nlark.com/yuque/__puml/550869d9c453748f1b46e5a0533f3e18.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmn6XnnIvnlKjmiLfmmK_lkKblrZjlnKhcbmFjdGl2YXRlIG15c3FsXG5cbm15c3FsIC0-IHNlcnZpY2U6IOi_lOWbnueUqOaIt-S_oeaBr1xuZGVhY3RpdmF0ZSBteXNxbFxuXG5ybm90ZSBvdmVyIHNlcnZpY2VcbuWvhueggeagoemqjFxuZW5kcm5vdGVcbnJub3RlIG92ZXIgc2VydmljZVxu55Sf5oiQc2Vzc2lvbklEXG5lbmRybm90ZVxuc2VydmljZSAtPiByZWRpczog57yT5a2Yc2Vzc2lvbklE5ZKM55So5oi35L-h5oGvXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHJlZGlzXG5cbnNlcnZpY2UgLT4gY29udHJvbGxlcjogcnBjIHJlc3BvbnNlXG5kZWFjdGl2YXRlIHNlcnZpY2Vcblxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lpITnkIbnu5PmnpxcbmVuZHJub3RlXG5cbmNvbnRyb2xsZXIgLT4gd2ViOiBodHRwIHJlc3BvbnNlXG5cblxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC81NTA4NjlkOWM0NTM3NDhmMWI0NmU1YTA1MzNmM2UxOC5zdmciLCJpZCI6ImRIQ2k4IiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="boHMv"></a>
### 登出
![](https://cdn.nlark.com/yuque/__puml/8cd1b7b882daf91db4cd723c3f6ac130.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IHJlZGlzOiDliKDpmaRzZXNzaW9uSURcbmFjdGl2YXRlIHJlZGlzXG5cbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvOGNkMWI3Yjg4MmRhZjkxZGI0Y2Q3MjNjM2Y2YWMxMzAuc3ZnIiwiaWQiOiJZWWVwVyIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="XvW7B"></a>
### 查看用户信息流程
![](https://cdn.nlark.com/yuque/__puml/a32aa0e8be8fe64f28b7cfee013d9c44.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IHJlZGlzOiDpqozor4FzZXNzaW9uSURcbmFjdGl2YXRlIHJlZGlzXG5cbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5ybm90ZSBvdmVyIHNlcnZpY2VcbnJlZGlz5a2Y5Zyo5YiZ6L-U5Zue5pWw5o2uXG5lbmRybm90ZVxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmn6Xor6LnlKjmiLfkv6Hmga9cbmFjdGl2YXRlIG15c3FsXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiByZWRpczog57yT5a2Y55So5oi35L-h5oGvXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHJlZGlzXG5cbnNlcnZpY2UgLT4gY29udHJvbGxlcjogcnBjIHJlc3BvbnNlXG5kZWFjdGl2YXRlIHNlcnZpY2Vcblxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lpITnkIbnu5PmnpxcbmVuZHJub3RlXG5cbmNvbnRyb2xsZXIgLT4gd2ViOiBodHRwIHJlc3BvbnNlXG5cblxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC9hMzJhYTBlOGJlOGZlNjRmMjhiN2NmZWUwMTNkOWM0NC5zdmciLCJpZCI6IkRSTzNCIiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="cchy9"></a>
### 修改头像流程
![](https://cdn.nlark.com/yuque/__puml/4fcbf41852b62f02d3b24942fde703ed.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbnJub3RlIG92ZXIgY29udHJvbGxlclxu5L-d5a2Y5Zu-54mHXG5lbmRybm90ZVxuYWN0aXZhdGUgY29udHJvbGxlclxuXG5jb250cm9sbGVyIC0-IHNlcnZpY2U6IHJwYyByZXF1ZXN0XG5hY3RpdmF0ZSBzZXJ2aWNlXG5cbnNlcnZpY2UgLT4gcmVkaXM6IOmqjOivgXNlc3Npb25JRFxuYWN0aXZhdGUgcmVkaXNcblxucmVkaXMgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHJlZGlzXG5cblxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmm7TmlrDnlKjmiLfkv6Hmga9cbmFjdGl2YXRlIG15c3FsXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiByZWRpczog5Yig6Zmk55So5oi35L-h5oGv57yT5a2YXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHJlZGlzXG5cbnNlcnZpY2UgLT4gY29udHJvbGxlcjogcnBjIHJlc3BvbnNlXG5kZWFjdGl2YXRlIHNlcnZpY2Vcblxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lpITnkIbnu5PmnpxcbmVuZHJub3RlXG5cbmNvbnRyb2xsZXIgLT4gd2ViOiBodHRwIHJlc3BvbnNlXG5cblxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC80ZmNiZjQxODUyYjYyZjAyZDNiMjQ5NDJmZGU3MDNlZC5zdmciLCJpZCI6IlJCSUhDIiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="ebBA2"></a>
### 修改昵称流程
![](https://cdn.nlark.com/yuque/__puml/8c4ef190d3bb1b26ce3817b0a7e0d07a.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJjb250cm9sbGVyXCIgYXMgY29udHJvbGxlclxucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxucGFydGljaXBhbnQgXCJyZWRpc1wiIGFzIHJlZGlzXG5wYXJ0aWNpcGFudCBcIm15c3FsXCIgYXMgbXlzcWxcblxuYWN0aXZhdGUgd2ViXG53ZWIgLT4gY29udHJvbGxlcjogaHR0cCByZXF1ZXN0XG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWPguaVsOagoemqjFxuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IHJlZGlzOiDpqozor4FzZXNzaW9uSURcbmFjdGl2YXRlIHJlZGlzXG5cbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5cblxuc2VydmljZSAtPiBteXNxbDog5pu05paw55So5oi35L-h5oGvXG5hY3RpdmF0ZSBteXNxbFxubXlzcWwgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIG15c3FsXG5cbnNlcnZpY2UgLT4gcmVkaXM6IOWIoOmZpOeUqOaIt-S_oeaBr-e8k-WtmFxuYWN0aXZhdGUgcmVkaXNcbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvOGM0ZWYxOTBkM2JiMWIyNmNlMzgxN2IwYTdlMGQwN2Euc3ZnIiwiaWQiOiJyNElNbyIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="iFQbi"></a>
### 鉴权
- 用户登录后校验用户名密码是否正确，正确则生成sessionID(全局唯一)，并将sessionID - username缓存到redis中(设置过期时间)
- 将生成的sessionID返回给客户端(set-cookie)，后续每次请求在cookie中携带sessionID
- 验证sessionID是否过期，如果过期则需要重新登录
<a name="kjyBR"></a>
### 密码加密
| **加密算法** | **加密方式** | **问题** |
| --- | --- | --- |
| SHA-256, SHA-1, md5 | 单向hash | 易破解 |
| salt | 密码先进行一次 MD5（或其它哈希算法）加密；将得到的 MD5 值前后加上随机串，再进行一次 MD5 加密 |  |
| scrypt |  | 大量内存计算，耗时久 |

<a name="VUHOM"></a>
### RPC
<a name="PawGC"></a>
#### 模块划分
**client**：RPC客户端，用于发起远程服务调用<br />**server**：RPC服务端，用于根据服务名称选择服务，并返回调用结果<br />**service**：RPC具体服务，实现反射处理远程调用具体逻辑，得到本地方法运行结果
<a name="SJToI"></a>
#### 流程图
![](https://cdn.nlark.com/yuque/__puml/6c077f780a37e331527fa22071e8d874.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcImNvbnRyb2xsZXJcIiBhcyBjb250cm9sbGVyXG5wYXJ0aWNpcGFudCBcIlJQQyBDbGllbnRcIiBhcyBjbGllbnRcbnBhcnRpY2lwYW50IFwiUlBDIFNlcnZlclwiIGFzIHNlcnZlclxucGFydGljaXBhbnQgXCJSUEMgU2VydmljZVwiIGFzIHN2Y1xucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxuXG5jbGllbnQgLT4gY2xpZW50OiBNYWtlQ2xpZW50XG5ybm90ZSBvdmVyIGNsaWVudFxuTWFrZUNsaWVudFxu5Yib5bu6dGNw6L-e5o6l5rGgXG5lbmRybm90ZVxuXG5cbnNlcnZlciAtPiAgc2VydmVyOiBNYWtlU2VydmVyXG5cblxucm5vdGUgb3ZlciBzZXJ2ZXJcbk1ha2VTZXJ2ZXJcbuWIm-W7unNlcnZpY2VOYW1l5ZKMc2VydmljZeaYoOWwhOihqFxuZW5kcm5vdGVcblxuc3ZjIC0-IHN2YzogTWFrZVNlcnZpY2VcbnJub3RlIG92ZXIgc3ZjXG5NYWtlU2VydmljZVxu5Yib5bu6bWV0aG9kTmFtZeWSjG1ldGhvZOaYoOWwhOihqFxuZW5kcm5vdGVcblxuc2VydmVyIC0-IHN2YzogcmVnaXN0ZXJcbnJub3RlIG92ZXIgc3ZjXG7lsIZzZXJ2aWNl5re75Yqg5Yiwc2VydmVy56uv55qEbWFw5LitXG5lbmRybm90ZVxuXG5cbnNlcnZlciAtPiBzZXJ2ZXI6IEFjY2VwdFxcbuebkeWQrOerr-WPo1xuYWN0aXZhdGUgc2VydmVyXG5cblxuXG5jb250cm9sbGVyIC0-IGNsaWVudDog6LCD55SoY2xpZW50LmNhbGwg5Y-R6LW3cnBj6LCD55SoXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5hY3RpdmF0ZSBjbGllbnRcblxucm5vdGUgb3ZlciBjbGllbnRcbmVuY29kZVxuZW5kcm5vdGVcblxuY2xpZW50IC0-IHNlcnZlcjogdGNwIGNvbm5lY3QgLi4uLi5cblxucm5vdGUgb3ZlciBzZXJ2ZXJcbmRlY29kZVxuZW5kcm5vdGVcblxucm5vdGUgb3ZlciBzZXJ2ZXJcbuWIpOaWrW1hcOS4reaYr-WQpuWMheWQq1xu6K-35rGC55qEc2VydmljZU5hbWVcbmVuZHJub3RlXG5cbnNlcnZlciAtPiBzdmM6IOmAieaLqXNlcnZpY2VcbmFjdGl2YXRlIHN2Y1xucm5vdGUgb3ZlciBzdmNcbuWIpOaWrW1hcOS4reaYr-WQpuWMheWQq1xu6K-35rGC55qEbWV0aG9kTmFtZVxuZW5kcm5vdGVcblxuc3ZjIC0-IHNlcnZpY2U6IOWPjeWwhOiwg-eUqOacrOWcsOaWueazlVxuYWN0aXZhdGUgc2VydmljZVxuc2VydmljZSAtPiBzdmM6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5zdmMgLT4gc2VydmVyOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgc3ZjXG5cbnJub3RlIG92ZXIgc2VydmVyXG5lbmNvZGVcbmVuZHJub3RlXG5cbnNlcnZlciAtPiBjbGllbnQ6IHRjcCBjb25uZWN0IC4uLi4uXG5cblxucm5vdGUgb3ZlciBjbGllbnRcbmRlY29kZVxuZW5kcm5vdGVcbmNsaWVudC0-Y29udHJvbGxlcjog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIGNsaWVudFxuZGVhY3RpdmF0ZSBjb250cm9sbGVyXG5cbkBlbmR1bWwiLCJ1cmwiOiJodHRwczovL2Nkbi5ubGFyay5jb20veXVxdWUvX19wdW1sLzZjMDc3Zjc4MGEzN2UzMzE1MjdmYTIyMDcxZThkODc0LnN2ZyIsImlkIjoiVkV3QjkiLCJtYXJnaW4iOnsidG9wIjp0cnVlLCJib3R0b20iOnRydWV9LCJjYXJkIjoiZGlhZ3JhbSJ9)![]
<a name="W5L88"></a>
## 四、接口设计
<a name="lws9T"></a>
### 错误码
| **err_code** | **err_msg** | **注释** |
| --- | --- | --- |
| 0 | ServerError | 服务端错误 |
| 1 | DataBaseError | 数据库错误 |
| 2 | InvalidSessionError | 过期Session |
| 3 | UserExistedError | 用户已存在 |
| 4 | UserNotExistError | 用户不存在 |
| 5 | PasswordError | 密码错误 |
| 6 | InvalidParamsError | 参数错误 |
| 7 | success | 成功 |

<a name="pfuyO"></a>
### 注册
**Post  api/entrytask/user/signup**<br />**入参：**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| username | string | 用户名<br />长度限制：[4,13] |
| password | string | 密码<br />长度限制：[4,13] |

```json
{
    "errCode":"7",
    "errMsg":"success",
    "data": ""
}
{
    "errCode":"6",
    "errMsg":"InvalidParamsError",
    "data":""
}
```
<a name="JyYo5"></a>
### 登录
**Post  api/entrytask/user/signin**<br />**入参**

| **字段名称** | **类型** | **注释** |
| --- | --- | --- |
| username | string | 用户名<br />长度限制：[4,13] |
| password | string | 密码<br />长度限制：[4,13] |

**返回值**

| **字段名称** | **类型** | **注释** |
| --- | --- | --- |
| sessionID | string | 设置在set-cookie中返回 |

```json
{
    "errCode":"7",
    "errMsg":"success",
    "data":""
}
{
    "errCode":"5",
    "errMsg":"PasswordError",
    "data":""
}
```
<a name="FJntJ"></a>
### 登出
**GET  api/entrytask/user/signout**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| sessionID | string | 从cookie中获取 |

```json
{
  "errCode":"7",
  "errMsg":"success",
  "data":""
}
{
    "errCode":"2",
    "errMsg":"InvalidSessionError",
    "data":""
}
```
<a name="GbzKu"></a>
### 查看用户信息
**GET  api/entrytask/user/get_user_info**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| sessionID | string | 从cookie中获取 |

**返回值**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| username | string | 用户名 |
| nickname | string | 昵称 |
| profilePath | string | 图片路径 |

```json
{
  "errCode":"7",
  "errMsg":"success",
  "data":{
    "username":"xq",
    "nickname":"nick",
    "profilePath":"test-2022-08-01.jpg",
  }
}
{
    "errCode":"2",
    "errMsg":"InvalidSessionError",
    "data":""
}
```
<a name="EvlbF"></a>
### 更新头像
**Post  api/entrytask/user/update_profile_pic**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| username | string | 用户名<br />长度限制：[4,13] |
| profilePic | file | 头像 |
| sessionID | string | 从cookie中获取 |

**返回值**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| profilePath | string | 图片路径 |

```json
{
    "errCode":"7",
    "errMsg":"success",
    "data":{
        "profilePath":"test-2022-08-01.jpg"
    }
}
{
    "errCode":"2",
    "errMsg":"InvalidParamsError",
    "data":""
}
```
<a name="nfsSb"></a>
### 更新昵称
**Post  api/entrytask/user/update_nickname**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| nickname | string | 更新的昵称<br />长度限制：[1,8] |
| sessionID | string | 从cookie中获取 |

**返回值**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| nickname | string | 昵称 |

```json
{
    "errCode":"7",
    "errMsg":"success",
    "data":"xxxxx"
}
{
    "errCode":"0",
    "errMsg":"ServerError",
    "data":""
}
```
<a name="DB2p9"></a>
## **五、存储设计**
<a name="RteyI"></a>
### MySql
**user表**

| **字段名称** | **字段类型** | **字段注释** |
| --- | --- | --- |
| id | bigint | 主键id |
| gmt_create | date_time | 创建时间 |
| gmt_modified | date_time | 修改时间 |
| username | varchar(64) | 用户名 |
| nickname | varchar(64) | 昵称 |
| password | varchar(64) | 密码 |
| salt | char(4) | 生成密码用到的随机值 |
| profile_path | varchar(128) | 图片存储路径 |

**primary key**: id<br />**unique key**：username
<a name="sujMW"></a>
### <br />Redis
**sessionID和username映射** ：string
```
key：sessionID 
value：username
```
**user信息缓存**：hash
```
key：username
value：[username:xq nickname:xxx profilePath:xxxxxxxxx.jpg]
```
过期时间设置：30分钟

<a name="pDzKc"></a>
## **六、外部依赖与限制**
无

<a name="lS885"></a>
## **七、部署方案与环境要求**


_<br />
<a name="l1avW"></a>
## **八、SLA**

- 200固定用户 qps大于3000    压测结果均值10000左右
```shell
wrk -c200 -t8 -d120s -s benchmark/getUserInfo.lua  -H "Cookie: sessionID=120f6f4a-4e92-44eb-89db-e39bb68e16ea" --latency http://localhost:9090/api/entrytask/user/get_user_info 
```
![固定200.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659062099644-3f58f7c7-bfeb-41fc-ad4b-92cd035abbf6.png#clientId=ua583142f-62f6-4&crop=0&crop=0&crop=1&crop=1&from=drop&id=uc741e62c&margin=%5Bobject%20Object%5D&name=%E5%9B%BA%E5%AE%9A200.png&originHeight=710&originWidth=3322&originalType=binary&ratio=1&rotation=0&showTitle=false&size=173304&status=done&style=none&taskId=u5eab2f6f-ebd7-42d9-8ce2-26efea41044&title=)

- 2000固定用户 qps大于1500    压测结果均值9000左右   
```shell
wrk -c2000 -t8 -d60s -s benchmark/getUserInfo.lua -H "Cookie: sessionID=6e87c1c2-a160-4ef4-9e6c-64f8ed459a33"" --latency http://localhost:9090/api/entrytask/user/get_user_info
```
![固定2000.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659062106812-b399a249-f2e2-47d2-942b-cc8d898eeed1.png#clientId=ua583142f-62f6-4&crop=0&crop=0&crop=1&crop=1&from=drop&id=uddcfd8a8&margin=%5Bobject%20Object%5D&name=%E5%9B%BA%E5%AE%9A2000.png&originHeight=752&originWidth=3342&originalType=binary&ratio=1&rotation=0&showTitle=false&size=188561&status=done&style=none&taskId=udc65dcaa-60c8-41c0-ab57-b40d4f71886&title=)

- 200随机用户 qps大于1000    压测结果均值6000左右 
```shell
wrk -c200 -t8 -d120s -s benchmark/signIn.lua --latency http://localhost:9090/api/entrytask/user/signin
```
![随机200.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659062116708-7a593e4f-881e-4738-88f3-09e17bf4faec.png#clientId=ua583142f-62f6-4&crop=0&crop=0&crop=1&crop=1&from=drop&id=ud0eda86e&margin=%5Bobject%20Object%5D&name=%E9%9A%8F%E6%9C%BA200.png&originHeight=758&originWidth=2218&originalType=binary&ratio=1&rotation=0&showTitle=false&size=159282&status=done&style=none&taskId=u7972e1da-b627-4df5-b9fa-1700d51ea25&title=)

- 2000随机用户 qps大于800    压测结果均值4000左右
```shell
wrk -c2000 -t8 -d120s -s benchmark/signIn.lua --latency http://localhost:9090/api/entrytask/user/signin
```
![随机2000.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659062123812-fc6b2933-7363-4304-9692-bb60dbc6c5c7.png#clientId=ua583142f-62f6-4&crop=0&crop=0&crop=1&crop=1&from=drop&id=u5336db1c&margin=%5Bobject%20Object%5D&name=%E9%9A%8F%E6%9C%BA2000.png&originHeight=754&originWidth=2164&originalType=binary&ratio=1&rotation=0&showTitle=false&size=159583&status=done&style=none&taskId=u74133068-2e2c-40eb-b266-f39ea34c6c9&title=)
<a name="tmhcV"></a>
## **九、遗留问题与风险预估**
rpc超时处理

_<br />
<a name="AvVRu"></a>
## **十、附录**
无
