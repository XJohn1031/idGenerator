#ID Generator

用go实现的id生成器, 支持每秒qps:131072, 超出需要等待下一秒


### 依赖
#### mysql(或zk,redis等)
需要使用mysql来保证多台机器获取到的workId不同
当然, 如果是单点, 那随意设置workId

### 使用介绍

#### init.sql
mysql的建表语句

#### config.json
在connect_info上添上连接mysql的信息

#### 用法
```go
    ring := NewRing(1<<17, GetWorkId("config.json"))
    id, err := ring.Take()
```
这里GetWorkId默认使用的是mysql, 你也可以使用zk等自定义的workId, 传入即可

