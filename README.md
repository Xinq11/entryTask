## 一、背景及目的
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
_给出整体系统架构图，或者与当前方案相关的部分架构图，包括主要模块、数据流、上下文关系。_<br />_如果是部分变动，应使用文字描述或者图形标识突出有变化的地方。_<br />_使用表格或者子章节的形式，描述主要模块和接口的功能；_
<a name="vvESq"></a>
### 系统架构图
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1658470040023-de40819d-5d6a-4ec6-ab0e-bdb1c83f1099.jpeg)<br />加UI层<br />加登出功能
<a name="ORSG3"></a>
### 页面转换
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1658470042186-4ddec177-e2c3-4c2e-bbba-5849662dd96e.jpeg)
<a name="r57di"></a>
## **三、核心逻辑详细设计**
_可包括：_<br />_1、模块间交互的时序图、状态图等；_<br />_2、复杂事务处理，一致性方案、重试与修复机制等；_<br />_3、数据采集、计算逻辑；_<br />_4、批量任务调度、依赖关系等；_<br />_5、涉及到信息和数据安全风险，如SQL注入、伪造和篡改请求等，及其应对的校验、授权机制；_<br />更新缓存
<a name="iJrpz"></a>
### 注册流程
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657620689598-bd7ffb30-2c44-4087-bb4b-b944932e0fdc.jpeg)
<a name="HQIEL"></a>
### 登录流程
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657643220261-0ad518a4-ece0-4f5c-b3ba-6a9d9a246a9b.jpeg)
<a name="XvW7B"></a>
### 查看用户信息流程
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657643171883-e96aea0b-ae63-4a01-a29a-5f32ca710313.jpeg)
<a name="cchy9"></a>
### 修改头像流程
在controller层处理图片，保存好后把路径通过Rpc传输<br />![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657643222190-dba3d362-8bf3-4c0a-9cfa-ab067984de8f.jpeg)
<a name="ebBA2"></a>
### 修改昵称流程
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657643224359-71c31497-98b2-4ea6-a73c-79a791d2ec4f.jpeg)
<a name="iFQbi"></a>
### 鉴权

- 用户登录后校验用户名密码是否正确，正确则生成sessionID(全局唯一)，并将sessionID - username缓存到redis中(设置过期时间)
- 将生成的sessionID返回给客户端(set-cookie)，后续每次请求在cookie中携带sessionID
- 验证sessionID是否过期，如果过期则需要重新登录

**sessionID生成**<br />Google uuid
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
![](https://cdn.nlark.com/yuque/0/2022/jpeg/21719644/1657608067652-c80d8ed2-9ad6-4a80-9dff-fc2282ddf80c.jpeg)
<a name="umi8r"></a>
#### 序列化/反序列化方式
基准测试

| 序列化方式 | <br /> |  |
| --- | --- | --- |
| json | <br /> | <br /> |
| gob |  |  |
| gogo-protobuf |  |  |

<a name="COzMk"></a>
#### 自定义传输协议
固定报头长度，把数据长度写入报头。
<a name="W5L88"></a>
## 四、接口设计
_给出对外接口描述，包括名称、地址或URL、类型（GRPC/HTTPS等）、协议内容（PB/JSON等）、各参数（类型、描述、限制）等；_<br />_对外接口需要给出鉴权机制和其他安全性考虑；_
<a name="lws9T"></a>
### 错误码
| **err_code** | **err_msg** | **注释** |
| --- | --- | --- |
| 0 | success | 成功 |
| 1 | InvalidParamsError | 参数错误 |
| 2 | PasswordError | 密码错误 |
| 3 | UserNotExistError | 用户不存在 |
| 4 | UserExistedError | 用户已存在 |
| 5 | InvalidSessionError | 过期Session |
| 6 | ServerError | 服务端错误 |

<a name="pfuyO"></a>
### 注册
**Post  api/entrytask/user/signup**<br />**入参：**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| username | string | 用户名<br />长度限制：0<len<=8 | 否 |
| password | string | 密码<br />长度限制：8<=len<=16<br />格式限制：包含字母和数字 | 否 |

```json
// success
{
    "err_code":"0",
    "err_msg":"success",
    "data": ""
}
// fail
{
    "err_code":"1",
    "err_msg":"InvalidParamsError",
    "data":""
}
```
<a name="JyYo5"></a>
### 登录
**Post  api/entrytask/user/signin**<br />**入参**

| **字段名称** | **类型** | **注释** | **是否可空** |
| --- | --- | --- | --- |
| username | string | 用户名<br />长度限制：0<len<=8 | 否 |
| password | string | 密码<br />长度限制：8<=len<=16<br />格式限制：包含字母和数字 | 否 |

**返回值**

| **字段名称** | **类型** | **注释** | **是否可空** |
| --- | --- | --- | --- |
| sessionID | string | 设置在set-cookie中返回 | <br /> |

```json
// success
{
    "err_code":"0",
    "err_msg":"success",
    "data":""
}
// fail
{
    "err_code":"2",
    "err_msg":"PasswordError",
    "data":""
}
```
<a name="hRbLv"></a>
### 查看用户信息
**GET  api/entrytask/user/get_user_info **<br />**入参**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| sessionID | string | 从cookie中获取 | <br /> |

**返回值**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| username | string | 用户名 | 否 |
| nickname | string | 昵称 | 是 |
| profilePath | string | 图片路径 | 是 |

```json
// success
{
  "err_code":"0",
  "err_msg":"success",
  "data":{
    "username":"xq",
    "nickname":"nick",
    "profilePath":"xxxxxxx",
  }
}
// fail
{
    "err_code":"5",
    "err_msg":"InvalidSessionError",
    "data":""
}
```
<a name="EvlbF"></a>
### 更新头像
**Post  api/entrytask/user/update_profile_pic**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| username | string | 用户名 | 否 |
| profilePic | file | 图片<br />大小限制：2Mb | 否 |
| sessionID | string | 从cookie中获取 | <br /> |

**返回值**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| profilePath | string | 图片路径 | 是 |

```json
// success
{
    "err_code":"0",
    "err_msg":"success",
    "data":{
        "profilePath":"xxxxxxxxxxxxxx"
    }
}
// fail
{
    "err_code":"1",
    "err_msg":"InvalidParamsError",
    "data":""
}
```
<a name="nfsSb"></a>
### 更新昵称
**Post  api/entrytask/user/update_nickname**<br />**入参**

| **字段名称** | **字段类型** | **字段注释** | **是否可空** |
| --- | --- | --- | --- |
| nickname | string | 更新后的昵称<br />长度限制：0<len<=16 | 否 |
| sessionID | string | 从cookie中获取 | <br /> |

```json
// success
{
    "err_code":"0",
    "err_msg":"success",
    "data":""
}
// fail
{
    "err_code":"5",
    "err_msg":"ServerError",
    "data":""
}
```
<a name="DB2p9"></a>
## **五、存储设计**
_可包括：_<br />_1、数据库表定义、字段定义、索引、主/备库读写等；""_<br />_2、缓存KV设计、加载/更新/失效逻辑等；_<br />_3、文件存储介质、目录组织形式、数据格式、索引、保存周期与清理机制等；_<br />_4、消息队列等中间件选型；_

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
_可包括：_<br />_1、对公司内部跨产品系统或者外部的依赖，如SPM、CoreServer、银行、供应商等。_<br />_2、受限于外部依赖的现有条件产生的限制，比如接口时延、TPS等；_

<a name="lS885"></a>
## **七、部署方案与环境要求**
_可包括：_<br />_1、配置初始化、更改、下发、推送等；_<br />_2、各种存储容量的预估、需要扩容的实例、备库、账号要求等；_<br />_3、接入层、逻辑层实例数；_<br />_4、VIP、域名、防火墙等特殊要求；_<br />_<br />
<a name="l1avW"></a>
## **八、SLA**
_可包括：_<br />_1、可支持的存储容量；_<br />_2、关键接口可支撑的TPS、QPS、时延等；_<br />_3、系统可用性保证，如故障时间、恢复时间、Crash率等；_

不用全部压测， 压测登录或者查看用户信息界面接口_go-wrk -t=8 -c=100 -n=200 -H = "sessionID=470157e3-2109-4f93-b534-2a9867cb9b90" "[http://localhost:9090/api/entrytask/user/get_user_info](http://localhost:9090/api/entrytask/user/get_user_info)"<br />�go-torch -u "[http://localhost:9090/api/entrytask/user/get_user_info](http://localhost:9090/api/entrytask/user/get_user_info)" -t 30

<a name="tmhcV"></a>
## **九、遗留问题与风险预估**
_可包括：_<br />_1、本次方案受限于时间、人力、外部因素等原因，未充分设计或者实现，可能带来的影响，以及下阶段的改进计划。_<br />_2、受限于外部依赖限制、硬件资源、网络等不可控条件，存在的运行和运营风险。_<br />_<br />
<a name="AvVRu"></a>
## **十、附录**
_可包括一些附带的内容，例如引用的文档链接、提供的操作手册附件等。_
