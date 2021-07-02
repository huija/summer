<div align=center><img width="800" height=480" src="https://huija.github.io/images/summer.png"/></div>

# summer
![GitHub](https://img.shields.io/github/license/huija/summer)
![Build](https://www.travis-ci.com/huija/summer.svg?branch=main)

Blooming like summer flowers
> config your own web server

- http: github.com/gin-gonic/gin
- log: go.uber.org/zap & gopkg.in/natefinch/lumberjack.v2
- mysql: github.com/go-sql-driver/mysql
- redis: github.com/go-redis/redis
- mongo: go.mongodb.org/mongo-driver/mongo

> 后续也可以集成其余框架的一些轮子, 通过配置项直接选择

## create a basic repo

- create your repo
```shell script
cd /your/aim/path
mkdir repo
cd repo
go mod init github.com/huija/repo
mkdir conf
touch main.go conf/config.yaml
```

- `main.go`
```go
package main

import (
	"github.com/huija/summer"
)

func main() {
	_ = summer.Bloom()
}
```

- cli exec
```go
go run main.go
```

- open in broswer: http://localhost:8080,  `404` because there was no route

## add a router

`summer` use gin as the default web framework, use`AddStage` to append router init func.

> All of summer's inner objects will init in `summer.Bloom()`, it's a pipeline init.
  So, do not use them before `summer.Bloom()` , the true way is `AddStage`.
> the name prefix is necessary, such as `summer.Gin+summer.Split`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huija/summer"
	"github.com/huija/summer/logs"
	"net/http"
)

func main() {
	summer.AddStage(summer.Gin+summer.Split+"router", router)
	err := summer.Bloom()
	if err != nil {
		logs.SugaredLogger.Panic(err)
	}
}

func router() error {
	summer.GinSrv.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})
	return nil
}
```
> 这样启动后, 可以访问`http://localhost:9090`, 就可以得到json返回

> summer内置了一些配置GinSrv的方法, 可以选择性打开，只需要在配置文件配置好就行

```yaml
srv:
  cors: true #允许跨域 
  trace: true #trace日志
  pprof: true #pprof路由
```

## simple config

`summer`通过配置来直接配置初始化内容, 减少代码的重复书写.

- `conf.yaml`的内容改成如下
```yaml
srv:
  name: repo
  tag: v0.0.1
  host: localhost
  author: huija
  listen: 0.0.0.0
  port: 8080
logs:
  # -1:debug, 0:info, 1: warn, 2: error, 3: panic, 4: fatal
  - level: -1
    type: file
    file:
      path: ./conf/test.log
      localzone: true
  - level: -1
    type: console
```
- 启动后的端口就变成了8080, 并且日志会同时输出到终端和文件.
- 终端还是普通的日志格式, 文件变成了文件
- 自定义部分转后面高级: TODO

## dbs config
### redis config
先以redis为例, 在dbs键下添加配置(单机, 哨兵, 集群同配置即可)
```yaml
dbs:
  redis:
    addrs: [127.0.0.1:6379]
    maxpoolsize: 100
    minpoolsize: 10
```
然后就可以直接通过dbs.RedisDB来访问redis了(比如将如下代码加在之前的`router`中)
```go
	summer.GinSrv.GET("/redis/ping", func(c *gin.Context) {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		set := dbs.RedisDB.Set(ctx, "ping", "pong", 0)
		logs.SugaredLogger.Debug(set)
		c.JSON(http.StatusOK, gin.H{
			"msg": dbs.RedisDB.Get(ctx, "ping").String(),
		})
	})
```
启动
```shell script
go run main.go
```
访问
```shell script
curl -i http://localhost:9090/redis/ping
```
> 需要注意的是, 所有这类全局的连接池, 都是在`summer.Bloom()`后才会初始化的, 所以不能在`Bloom()`之前使用
> 如果仍然报了空指针异常, 估计是配置的字符串不对, 可以根据日志查一下
> 会打印出实际使用的配置文件的总体
> `=>`字符也会描述程序启动的内部流程

### mysql config
mysql配置如下:
```yaml
dbs:
  mysql:
    schema: root:123456@tcp(127.0.0.1:3306)/repo?charset=utf8mb4&parseTime=True&loc=Local
    maxpoolsize: 100
    minpoolsize: 10
```
访问方式
```go
	summer.GinSrv.GET("/mysql/ping", func(c *gin.Context) {
		err := dbs.MysqlDB.Ping()
		if err != nil {
			logs.SugaredLogger.Error(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
```
### mongodb config
mongodb配置如下, 不同部署模式(单机, 复制集, 分片集群)都只需要改schema即可
```yaml
dbs:
  mongo:
    schema: mongodb://127.0.0.1:27017/repo
    maxpoolsize: 100
    minpoolsize: 10
```
调用方式
```go
	summer.GinSrv.GET("/mongo/ping", func(c *gin.Context) {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		err := dbs.MongoDB.Ping(ctx, readpref.Primary())
		if err != nil {
			logs.SugaredLogger.Error(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
```

## Advanced
TODO

### RegisterClose
- 这个注册close方法, 可以实现优雅关闭
- 常在stage方法内调用, 保证与rootPipe相反的顺序进行关闭
```go
summer.RegisterClose(name string,fn func() error)
```

*TODO*
