# hook api

zlmediakit程序配套hook api程序

# zlmediakit hook api

## 鉴权类

文件鉴权、rtsp鉴权、shell鉴权、流播放鉴权等
```
POST /api/v1/auth/http_file
POST /api/v1/auth/rtsp_auth
POST /api/v1/auth/rtsp_play
POST /api/v1/auth/shell
POST /api/v1/auth/stream_play
POST /api/v1/auth/stream_publish

```
## 系统类
流量报告、心跳报告、rtsp关闭报告、系统启动报告
```
POST /api/v1/service/flow_report
POST /api/v1/service/keepalive_report
POST /api/v1/service/rtp_close_report
POST /api/v1/service/rtp_timeout_report
POST /api/v1/service/startup_report
```

## 流事件类
要播放的流不存在、流无人观看
```
POST /api/v1/stream/change
POST /api/v1/stream/none_reader
POST /api/v1/stream/not_found
```

## 录像类
mp4、ts录像完成事件
```
POST /api/v1/record/mp4_finish
POST /api/v1/record/ts_finish
```

# 设备管理api
对外提供设备管理restfull api, 提供设备的增加、修改、删除、查询等，方便把录像设备使用sqlite管理，在播放流不存在的时候，通过查询设备管理表，把对应的设备的流接入到zlmediakit中.

## api
```
查询所有设备信息，密码加密
GET /api/v1/device/
增加设备信息
POST /api/v1/device/
修改设备信息
PUT /api/v1/device/
查询指定设备信息
GET /api/v1/device/{id}/
删除设备
DELETE /api/v1/device/{id}/
```

# 文件组织格式

```
root
    api                          api接入目录
                    auth.go      鉴权类
                    common.go    公共类定义了统一返回对象
                    device.go    设备管理类
                    record.go    录像类
                    service.go   系统相关的逻辑
                    stream.go    流相关逻辑
    library                      类库,主要包括logger、config、database、utils等
                config           实现配置文件解析
                database         sqlite数据库实现
                logger           日志
                net              网络相关
                os               环境变量读取、解压缩等
                utils            序号的生成、json string
    message                      请求消息 定义
    routes                       路由定义
                     routes.go   路由定义规则
    service                      复杂逻辑实现
                     device_service.go  设备管理dao实现
    main.go 程序入口
```

# 测试步骤

## url地址
```
使用vlc测试播放地址:
播放地址: rtsp://localhost/appname/streamid

打开一个设备流url:
http://localhost:8080/api/v1/stream/open/streamid

关闭一个设备流url:
http://localhost:8080/api/v1/stream/close/streamid

```
## 添加视频设备
通过api添加视频设备信息
url: http://localhost:8080/api/v1/device/
method: POST
data:{"streamId":"1000000004","username":"admin","password":"******","hostname":"172.17.18.233","appName":"app","vhostName":""}

## 播放视频
根据上一步骤中添加设备的streamId、app信息,可以使用vlc验证是否可以播放
播放url地址: rtsp://localhost/appname/streamid
