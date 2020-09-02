你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
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

