# beewolf
small crawler powered by go!!

### 结构?

ship包是一个简单的爬虫框架，spiders包含了熊猫TV、战旗TV、全民TV、虎牙TV、斗鱼TV(未实现)这几个主流平台的爬虫


### 如何运行?

```go run claw.go
...
2016/06/09 19:16:10 [熊猫TV]: Pool of urls is empty...
2016/06/09 19:16:10 [全民TV]: Pool of urls is empty...
2016/06/09 19:16:10 [战旗TV]: Pool of urls is empty...
[战旗TV]: 当前在线5949087人
2016/06/09 19:16:10 [战旗TV]: Finish after work...
[熊猫TV]: 当前在线4479101人
[全民TV]: 当前在线6213468人
2016/06/09 19:16:10 [熊猫TV]: Finish after work...
2016/06/09 19:16:10 [全民TV]: Finish after work...
2016/06/09 19:16:10 HTTP GET http://www.huya.com/index.php?m=Live&do=ajaxAllLiveByPage&page=168&pageNum=1
2016/06/09 19:16:10 [虎牙TV]: Pool of urls is empty...
[虎牙TV]: 当前在线384049人
2016/06/09 19:16:10 [虎牙TV]: Finish after work...
```

### 接下来?

斗鱼TV的爬虫还没有完成.

### 问题?

chan的超时机制和goroutine切换引起的冲突问题，本意是想用selet配合一个超时chan来实现信道读取的超时退出。但是多个goroutine的自动切换导致爬虫会被强制超时退出。
