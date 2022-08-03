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
│   ├── mySqlPool.go
│   └── redisPool.go
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
│   ├── pool
│   │   └── tcpPool.go
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
![](https://cdn.nlark.com/yuque/__puml/c14c9d1f515c3e1e55480660d79edd88.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5cbmNvbnRyb2xsZXIgLT4gc2VydmljZTogcnBjIHJlcXVlc3RcbmFjdGl2YXRlIHNlcnZpY2Vcblxuc2VydmljZSAtPiBteXNxbDog5p-l55yL55So5oi35piv5ZCm5a2Y5ZyoXG5hY3RpdmF0ZSBteXNxbFxuXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nlKjmiLfkv6Hmga9cbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2Vcbm5vdGUgcmlnaHRcbueUqOaIt-W3suWtmOWcqFxuZW5kIG5vdGVcblxucm5vdGUgb3ZlciBzZXJ2aWNlXG7lr4bnoIHliqDlr4ZcbmVuZHJub3RlXG5cbnNlcnZpY2UgLT4gbXlzcWw6IOWtmOWCqOeUqOaIt-S_oeaBr1xuYWN0aXZhdGUgbXlzcWxcbm15c3FsIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSBteXNxbFxuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvYzE0YzlkMWY1MTVjM2UxZTU1NDgwNjYwZDc5ZWRkODguc3ZnIiwiaWQiOiJLSzlLUiIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="HQIEL"></a>
### 登录流程
![](https://cdn.nlark.com/yuque/__puml/82eacce7974f2d297c77a4e587eae9b3.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5cbmNvbnRyb2xsZXIgLT4gc2VydmljZTogcnBjIHJlcXVlc3RcbmFjdGl2YXRlIHNlcnZpY2Vcblxuc2VydmljZSAtPiBteXNxbDog5p-l55yL55So5oi35piv5ZCm5a2Y5ZyoXG5hY3RpdmF0ZSBteXNxbFxuXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nlKjmiLfkv6Hmga9cbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2Vcbm5vdGUgcmlnaHRcbueUqOaIt-S4jeWtmOWcqFxuZW5kIG5vdGVcblxucm5vdGUgb3ZlciBzZXJ2aWNlXG7lr4bnoIHmoKHpqoxcbmVuZHJub3RlXG5ybm90ZSBvdmVyIHNlcnZpY2VcbueUn-aIkHNlc3Npb25JRFxuZW5kcm5vdGVcbnNlcnZpY2UgLT4gcmVkaXM6IOe8k-WtmHNlc3Npb25JROWSjOeUqOaIt-S_oeaBr1xuYWN0aXZhdGUgcmVkaXNcbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvODJlYWNjZTc5NzRmMmQyOTdjNzdhNGU1ODdlYWU5YjMuc3ZnIiwiaWQiOiJkSENpOCIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="boHMv"></a>
### 登出
![](https://cdn.nlark.com/yuque/__puml/7492e33a7a93aef0cdfd33fc810bc127.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5cbmNvbnRyb2xsZXIgLT4gc2VydmljZTogcnBjIHJlcXVlc3RcbmFjdGl2YXRlIHNlcnZpY2Vcblxuc2VydmljZSAtPiByZWRpczog5Yig6Zmkc2Vzc2lvbklEXG5hY3RpdmF0ZSByZWRpc1xuXG5yZWRpcyAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgcmVkaXNcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2VcbmRlYWN0aXZhdGUgc2VydmljZVxuXG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWkhOeQhue7k-aenFxuZW5kcm5vdGVcblxuY29udHJvbGxlciAtPiB3ZWI6IGh0dHAgcmVzcG9uc2VcblxuXG5cbkBlbmR1bWwiLCJ1cmwiOiJodHRwczovL2Nkbi5ubGFyay5jb20veXVxdWUvX19wdW1sLzc0OTJlMzNhN2E5M2FlZjBjZGZkMzNmYzgxMGJjMTI3LnN2ZyIsImlkIjoiWVllcFciLCJtYXJnaW4iOnsidG9wIjp0cnVlLCJib3R0b20iOnRydWV9LCJjYXJkIjoiZGlhZ3JhbSJ9)<a name="XvW7B"></a>
### 查看用户信息流程
![](https://cdn.nlark.com/yuque/__puml/abb92c2dd0ae336a382e5a48fa56773b.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5cbmNvbnRyb2xsZXIgLT4gc2VydmljZTogcnBjIHJlcXVlc3RcbmFjdGl2YXRlIHNlcnZpY2Vcblxuc2VydmljZSAtPiByZWRpczog6aqM6K-Bc2Vzc2lvbklEXG5hY3RpdmF0ZSByZWRpc1xuXG5yZWRpcyAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgcmVkaXNcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2Vcbm5vdGUgcmlnaHRcbnNlc3Npb25JROS4jeWtmOWcqFxuZW5kIG5vdGVcblxuc2VydmljZSAtPiByZWRpczog5p-l6K-i57yT5a2YXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gc2VydmljZTog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHJlZGlzXG5cbnNlcnZpY2UgLT4gY29udHJvbGxlcjogcnBjIHJlc3BvbnNlXG5ub3RlIHJpZ2h0XG7mnInnvJPlrZjliJnov5Tlm57mlbDmja5cbmVuZCBub3RlXG5cbnNlcnZpY2UgLT4gbXlzcWw6IOafpeivoueUqOaIt-S_oeaBr1xuYWN0aXZhdGUgbXlzcWxcbm15c3FsIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSBteXNxbFxuXG5zZXJ2aWNlIC0-IHJlZGlzOiDnvJPlrZjnlKjmiLfkv6Hmga9cbmFjdGl2YXRlIHJlZGlzXG5yZWRpcyAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgcmVkaXNcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2VcbmRlYWN0aXZhdGUgc2VydmljZVxuXG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuWkhOeQhue7k-aenFxuZW5kcm5vdGVcblxuY29udHJvbGxlciAtPiB3ZWI6IGh0dHAgcmVzcG9uc2VcblxuXG5cbkBlbmR1bWwiLCJ1cmwiOiJodHRwczovL2Nkbi5ubGFyay5jb20veXVxdWUvX19wdW1sL2FiYjkyYzJkZDBhZTMzNmEzODJlNWE0OGZhNTY3NzNiLnN2ZyIsImlkIjoiRFJPM0IiLCJtYXJnaW4iOnsidG9wIjp0cnVlLCJib3R0b20iOnRydWV9LCJjYXJkIjoiZGlhZ3JhbSJ9)<a name="cchy9"></a>
### 修改头像流程
![](https://cdn.nlark.com/yuque/__puml/b59cd676e46ab602bd6d7393f6ef4c44.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5ybm90ZSBvdmVyIGNvbnRyb2xsZXJcbuS_neWtmOWbvueJh1xuZW5kcm5vdGVcbmFjdGl2YXRlIGNvbnRyb2xsZXJcblxuY29udHJvbGxlciAtPiBzZXJ2aWNlOiBycGMgcmVxdWVzdFxuYWN0aXZhdGUgc2VydmljZVxuXG5zZXJ2aWNlIC0-IHJlZGlzOiDpqozor4FzZXNzaW9uSURcbmFjdGl2YXRlIHJlZGlzXG5cbnJlZGlzIC0-IHNlcnZpY2U6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxubm90ZSByaWdodFxuc2Vzc2lvbklE5LiN5a2Y5ZyoXG5lbmQgbm90ZVxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmm7TmlrDnlKjmiLfkv6Hmga9cbmFjdGl2YXRlIG15c3FsXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiByZWRpczog5Yig6Zmk55So5oi35L-h5oGv57yT5a2YXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gcmVkaXM6IOWksei0pemHjeivlVxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvYjU5Y2Q2NzZlNDZhYjYwMmJkNmQ3MzkzZjZlZjRjNDQuc3ZnIiwiaWQiOiJSQklIQyIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="ebBA2"></a>
### 修改昵称流程
![](https://cdn.nlark.com/yuque/__puml/7bc07e5d028dd20f38aeb370e0ff6dbd.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcIua1j-iniOWZqFwiIGFzIHdlYlxucGFydGljaXBhbnQgXCJIVFRQIFNlcnZlclwiIGFzIGNvbnRyb2xsZXJcbnBhcnRpY2lwYW50IFwiVENQIFNlcnZlclwiIGFzIHNlcnZpY2VcbnBhcnRpY2lwYW50IFwiUmVkaXNcIiBhcyByZWRpc1xucGFydGljaXBhbnQgXCJNeVNxbFwiIGFzIG15c3FsXG5cbmFjdGl2YXRlIHdlYlxud2ViIC0-IGNvbnRyb2xsZXI6IGh0dHAgcmVxdWVzdFxucm5vdGUgb3ZlciBjb250cm9sbGVyXG7lj4LmlbDmoKHpqoxcbmVuZHJub3RlXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5cbmNvbnRyb2xsZXIgLT4gc2VydmljZTogcnBjIHJlcXVlc3RcbmFjdGl2YXRlIHNlcnZpY2Vcblxuc2VydmljZSAtPiByZWRpczog6aqM6K-Bc2Vzc2lvbklEXG5hY3RpdmF0ZSByZWRpc1xuXG5yZWRpcyAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgcmVkaXNcblxuc2VydmljZSAtPiBjb250cm9sbGVyOiBycGMgcmVzcG9uc2Vcbm5vdGUgcmlnaHRcbnNlc3Npb25JROS4jeWtmOWcqFxuZW5kIG5vdGVcblxuXG5zZXJ2aWNlIC0-IG15c3FsOiDmm7TmlrDnlKjmiLfkv6Hmga9cbmFjdGl2YXRlIG15c3FsXG5teXNxbCAtPiBzZXJ2aWNlOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgbXlzcWxcblxuc2VydmljZSAtPiByZWRpczog5Yig6Zmk55So5oi35L-h5oGv57yT5a2YXG5hY3RpdmF0ZSByZWRpc1xucmVkaXMgLT4gcmVkaXM6IOWksei0pemHjeivlVxuZGVhY3RpdmF0ZSByZWRpc1xuXG5zZXJ2aWNlIC0-IGNvbnRyb2xsZXI6IHJwYyByZXNwb25zZVxuZGVhY3RpdmF0ZSBzZXJ2aWNlXG5cbnJub3RlIG92ZXIgY29udHJvbGxlclxu5aSE55CG57uT5p6cXG5lbmRybm90ZVxuXG5jb250cm9sbGVyIC0-IHdlYjogaHR0cCByZXNwb25zZVxuXG5cblxuQGVuZHVtbCIsInVybCI6Imh0dHBzOi8vY2RuLm5sYXJrLmNvbS95dXF1ZS9fX3B1bWwvN2JjMDdlNWQwMjhkZDIwZjM4YWViMzcwZTBmZjZkYmQuc3ZnIiwiaWQiOiJyNElNbyIsIm1hcmdpbiI6eyJ0b3AiOnRydWUsImJvdHRvbSI6dHJ1ZX0sImNhcmQiOiJkaWFncmFtIn0=)<a name="iFQbi"></a>
### 鉴权

- 用户登录后校验用户名密码是否正确，正确则生成sessionID(全局唯一)，并将sessionID - username缓存到redis中(设置过期时间)
- 将生成的sessionID返回给客户端(set-cookie)，后续每次请求在cookie中携带sessionID
- 验证sessionID是否过期，如果过期则需要重新登录
<a name="M5x4B"></a>
### sessionID生成

1. 根据随机数生成UUID
1. 将第一步生成的UUID和username进行md5计算生成最终的sessionID
<a name="kjyBR"></a>
### 密码加密

- **salt加密**：密码先进行一次 MD5（或其它哈希算法）加密；将得到的 MD5 值前后加上随机串，再进行一次 MD5 加密
<a name="VUHOM"></a>
### RPC
<a name="PawGC"></a>
#### 模块划分
**client**：RPC客户端，用于发起rpc调用<br />**server**：RPC服务端，监听端口，接收rpc调用，并根据服务名称选择服务进行处理<br />**service**：RPC具体服务，根据rpc调用请求的方法名称，利用反射调用本地方法得到结果
<a name="SJToI"></a>
#### 流程图
![](https://cdn.nlark.com/yuque/__puml/f980229e16185219e154575ec149d32a.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuYXV0b251bWJlclxuXG5wYXJ0aWNpcGFudCBcImNvbnRyb2xsZXJcIiBhcyBjb250cm9sbGVyXG5wYXJ0aWNpcGFudCBcIlJQQyBDbGllbnRcIiBhcyBjbGllbnRcbnBhcnRpY2lwYW50IFwiUlBDIFNlcnZlclwiIGFzIHNlcnZlclxucGFydGljaXBhbnQgXCJSUEMgU2VydmljZVwiIGFzIHN2Y1xucGFydGljaXBhbnQgXCJzZXJ2aWNlXCIgYXMgc2VydmljZVxuXG5jbGllbnQgLT4gY2xpZW50OiBNYWtlQ2xpZW50XG5ybm90ZSBvdmVyIGNsaWVudFxuTWFrZUNsaWVudFxu5Yib5bu6dGNw6L-e5o6l5rGgXG5lbmRybm90ZVxuXG5cbnNlcnZlciAtPiAgc2VydmVyOiBNYWtlU2VydmVyXG5cblxucm5vdGUgb3ZlciBzZXJ2ZXJcbk1ha2VTZXJ2ZXJcbuWIm-W7unNlcnZpY2VOYW1l5ZKMc2VydmljZeaYoOWwhOihqFxuZW5kcm5vdGVcblxuc3ZjIC0-IHN2YzogTWFrZVNlcnZpY2VcbnJub3RlIG92ZXIgc3ZjXG5NYWtlU2VydmljZVxu5Yib5bu6bWV0aG9kTmFtZeWSjG1ldGhvZOaYoOWwhOihqFxuZW5kcm5vdGVcblxuc2VydmVyIC0-IHN2YzogcmVnaXN0ZXJcbnJub3RlIG92ZXIgc3ZjXG7lsIZzZXJ2aWNl5re75Yqg5Yiwc2VydmVy55qE5pig5bCE6KGo5LitXG5lbmRybm90ZVxuXG5cbnNlcnZlciAtPiBzZXJ2ZXI6IEFjY2VwdFxcbuebkeWQrOerr-WPo1xuYWN0aXZhdGUgc2VydmVyXG5cblxuXG5jb250cm9sbGVyIC0-IGNsaWVudDog6LCD55SoY2xpZW50LmNhbGwg5Y-R6LW3cnBj6LCD55SoXG5hY3RpdmF0ZSBjb250cm9sbGVyXG5hY3RpdmF0ZSBjbGllbnRcblxucm5vdGUgb3ZlciBjbGllbnRcbuiOt-WPlnRjcOi_nuaOpVxuZW5kcm5vdGVcblxucm5vdGUgb3ZlciBjbGllbnRcbuW6j-WIl-WMllxuZW5kcm5vdGVcblxuY2xpZW50IDwtPiBzZXJ2ZXI6IHRjcCBjb25uZWN0IC4uLi4uXG5cbnJub3RlIG92ZXIgc2VydmVyXG7lj43luo_liJfljJZcbmVuZHJub3RlXG5cbnJub3RlIG92ZXIgc2VydmVyXG7liKTmlq1tYXDkuK3mmK_lkKbljIXlkKtcbuivt-axgueahHNlcnZpY2VOYW1lXG5lbmRybm90ZVxuXG5zZXJ2ZXIgLT4gc3ZjOiDpgInmi6lzZXJ2aWNlXG5hY3RpdmF0ZSBzdmNcbnJub3RlIG92ZXIgc3ZjXG7liKTmlq1tYXDkuK3mmK_lkKbljIXlkKtcbuivt-axgueahG1ldGhvZE5hbWVcbmVuZHJub3RlXG5cbnN2YyAtPiBzZXJ2aWNlOiDlj43lsITosIPnlKjmnKzlnLDmlrnms5VcbmFjdGl2YXRlIHNlcnZpY2VcbnNlcnZpY2UgLT4gc3ZjOiDov5Tlm57nu5PmnpxcbmRlYWN0aXZhdGUgc2VydmljZVxuc3ZjIC0-IHNlcnZlcjog6L-U5Zue57uT5p6cXG5kZWFjdGl2YXRlIHN2Y1xuXG5ybm90ZSBvdmVyIHNlcnZlclxu5bqP5YiX5YyWXG5lbmRybm90ZVxuXG5zZXJ2ZXIgPC0-IGNsaWVudDogdGNwIGNvbm5lY3QgLi4uLi5cblxuXG5ybm90ZSBvdmVyIGNsaWVudFxu5Y-N5bqP5YiX5YyWXG5lbmRybm90ZVxuXG5ybm90ZSBvdmVyIGNsaWVudFxu6YeK5pS-dGNw6L-e5o6lXG5lbmRybm90ZVxuXG5jbGllbnQtPmNvbnRyb2xsZXI6IOi_lOWbnue7k-aenFxuZGVhY3RpdmF0ZSBjbGllbnRcbmRlYWN0aXZhdGUgY29udHJvbGxlclxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC9mOTgwMjI5ZTE2MTg1MjE5ZTE1NDU3NWVjMTQ5ZDMyYS5zdmciLCJpZCI6IlZFd0I5IiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="bVxEr"></a>
#### tcp连接池
**获取连接**<br />
![](https://cdn.nlark.com/yuque/__puml/1fc4777a6adf8c73c8b1dc79f866c198.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuc3RhcnRcblxuOmdldCBjb25uZWN0aW9uO1xuXG5pZiAo5a2Y5Zyo56m66Zey6L-e5o6lKSB0aGVuICh0cnVlKVxuICA6cmV0dXJuO1xuXHRlbmRcbmVsc2UgKGZhbHNlKVxuICBpZiAo5b2T5YmN6L-e5o6l5pWw5aSn5LqO5pyA5aSn6L-e5o6l5pWwKSB0aGVuICh0cnVlKVxuXHRcdFx0OumYu-Whnjtcblx0XHRcdGVuZFxuXHRlbHNlIChmYWxzZSlcblx0XHRcdDrliJvlu7rmlrDov57mjqU7XG5cdFx0XHRlbmRcblx0ZW5kaWZcbmVuZGlmXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC8xZmM0Nzc3YTZhZGY4YzczYzhiMWRjNzlmODY2YzE5OC5zdmciLCJpZCI6IlFsNHlMIiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)
<br />**释放连接**<br />
![](https://cdn.nlark.com/yuque/__puml/e0b597b88a46de9e97bc85e9e565c699.svg#lake_card_v2=eyJ0eXBlIjoicHVtbCIsImNvZGUiOiJAc3RhcnR1bWxcblxuc3RhcnRcblxuOnJlbGVhc2UgY29ubmVjdGlvbjtcblxuaWYgKOWtmOWcqOmYu-WhnuetieW-hei_nuaOpemHiuaUvikgdGhlbiAodHJ1ZSlcbiAgOui_nuaOpeWkjeeUqDtcbiAgZW5kXG5lbHNlIChmYWxzZSlcbiAgaWYgKOW9k-WJjei_nuaOpeaVsOWkp-S6juacgOWkp-i_nuaOpeaVsHx856m66Zey6L-e5o6l5rGg5bey5ruhKSB0aGVuICh0cnVlKVxuXHRcdDrlhbPpl63ov57mjqU7XG5cdFx0ZW5kXG5cdGVsc2UgKGZhbHNlKVxuXHRcdDrmlL7lm57nqbrpl7Lov57mjqXmsaA7XG4gIGVuZFxuZW5kaWZcblxuXG5AZW5kdW1sIiwidXJsIjoiaHR0cHM6Ly9jZG4ubmxhcmsuY29tL3l1cXVlL19fcHVtbC9lMGI1OTdiODhhNDZkZTllOTdiYzg1ZTllNTY1YzY5OS5zdmciLCJpZCI6IldTRFFhIiwibWFyZ2luIjp7InRvcCI6dHJ1ZSwiYm90dG9tIjp0cnVlfSwiY2FyZCI6ImRpYWdyYW0ifQ==)<a name="Eaot8"></a>
#### 传输协议
设置固定字节长度来存放数据长度
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
<a name="hRbLv"></a>
### 查看用户信息
**GET  api/entrytask/user/get_user_info **<br />**入参**

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
<a name="WkIJq"></a>
### Redis
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
web Server for Chrome插件：用于在本地电脑启动一个临时文件服务器，进行图片访问
<a name="lS885"></a>
## **七、部署方案与环境要求**

- Golang 版本：1.12.7
- MySql版本：5.7.10
- Redis版本：7.0.3
- MySql 初始用户数据量：<br />![截屏2022-08-03 下午2.43.37.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659509084455-c3b6d6d9-c964-496e-8b40-9040d4a3e099.png#clientId=u8a701fd8-f2d7-4&crop=0&crop=0&crop=1&crop=1&from=ui&id=ufba0beb5&margin=%5Bobject%20Object%5D&name=%E6%88%AA%E5%B1%8F2022-08-03%20%E4%B8%8B%E5%8D%882.43.37.png&originHeight=122&originWidth=340&originalType=binary&ratio=1&rotation=0&showTitle=false&size=52322&status=done&style=none&taskId=u5d160efe-761d-4d2f-be0c-0be0ea82f5a&title=)
<a name="eCmSM"></a>
### 启动tcpServer
```shell
cd /Users/qi.xin/Projects/EntryTask/cmd/tcpServer
go build tcpServer.go
./tcpServer
```
<a name="EdZqH"></a>
### 启动httpServer
```shell
cd /Users/qi.xin/Projects/EntryTask/cmd/httpServer
go build httpServer.go
./httpServer
```
<a name="l1avW"></a>
## **八、SLA**

- 200固定用户 qps大于3000    压测结果均值10000左右
```shell
wrk -c200 -t8 -d120s -s benchmark/getUserInfo.lua  -H "Cookie: sessionID=5fb4c437-8310-35fb-848c-8985461f762f" --latency http://localhost:9090/api/entrytask/user/get_user_info 
```
![固定 200.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659458763234-2ff09473-6b58-4732-9e91-a6bbfee105f1.png#clientId=u3b873f51-20af-4&crop=0&crop=0&crop=1&crop=1&from=ui&id=u086d2221&margin=%5Bobject%20Object%5D&name=%E5%9B%BA%E5%AE%9A%20200.png&originHeight=760&originWidth=3336&originalType=binary&ratio=1&rotation=0&showTitle=false&size=187098&status=done&style=none&taskId=ucd81811d-f636-4f81-81d9-874aef24439&title=)

- 2000固定用户 qps大于1500    压测结果均值7000左右   
```shell
wrk -c2000 -t8 -d120s -s benchmark/getUserInfo.lua -H "Cookie: sessionID=a73fd784-c132-3963-84bf-ba586b98a6a9" --latency http://localhost:9090/api/entrytask/user/get_user_info
```
![固定 2000.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659458774991-58f41dfc-f8b7-4e9a-b178-88b671e5ec01.png#clientId=u3b873f51-20af-4&crop=0&crop=0&crop=1&crop=1&from=ui&id=u9210e18c&margin=%5Bobject%20Object%5D&name=%E5%9B%BA%E5%AE%9A%202000.png&originHeight=756&originWidth=3304&originalType=binary&ratio=1&rotation=0&showTitle=false&size=184884&status=done&style=none&taskId=u6dcb4255-40c9-4d8c-b525-85b9da68999&title=)

- 200随机用户 qps大于1000    压测结果均值6000左右 
```shell
wrk -c200 -t8 -d120s -s benchmark/signIn.lua --latency http://localhost:9090/api/entrytask/user/signin
```
![随机 200.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659521017147-1c9348e6-a136-4195-8ecb-802bb7b115d8.png#clientId=ubff15b93-6bb8-4&crop=0&crop=0&crop=1&crop=1&from=ui&id=u527d4cad&margin=%5Bobject%20Object%5D&name=%E9%9A%8F%E6%9C%BA%20200.png&originHeight=750&originWidth=2158&originalType=binary&ratio=1&rotation=0&showTitle=false&size=157505&status=done&style=none&taskId=u77ca23b6-172b-4403-adb4-950de917097&title=)

- 2000随机用户 qps大于800    压测结果均值5000左右
```shell
wrk -c2000 -t8 -d120s -s benchmark/signIn.lua --latency http://localhost:9090/api/entrytask/user/signin
```
![随机 2000.png](https://cdn.nlark.com/yuque/0/2022/png/21719644/1659521060519-ef9918e2-a0ab-43dc-870e-376dd1e9e06e.png#clientId=ubff15b93-6bb8-4&crop=0&crop=0&crop=1&crop=1&from=ui&id=uaa6f4ea6&margin=%5Bobject%20Object%5D&name=%E9%9A%8F%E6%9C%BA%202000.png&originHeight=760&originWidth=2174&originalType=binary&ratio=1&rotation=0&showTitle=false&size=158299&status=done&style=none&taskId=u47fd536c-8842-44c3-a735-ef55f892ff2&title=)
<a name="tmhcV"></a>
## **九、遗留问题与风险预估**
接口设计存在冗余，可以简化<br />rpc动态代理，超时处理

<a name="AvVRu"></a>
## **十、附录**
无

