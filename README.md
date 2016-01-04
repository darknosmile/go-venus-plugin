#go-venus-plug
go-venus-plugin是开源服务框架[Venus](http://wiki.hexnova.com/display/Venus/HOME)的[GO](http://golang.org/)语言客户端。

### Documentation
* [Veuns通讯协议](http://wiki.hexnova.com/pages/viewpage.action?pageId=1507425)
* [Venus协议以及交互流程](http://wiki.hexnova.com/pages/viewpage.action?pageId=622616)

### How do I use it?

```go
	venus := Conncect("10.32.172.124:16800")
	venus.AuthByDummy("venus")
	var userData models.FindNameData
	userData.UserName = "cscadmin"
	body, err := json.Marshal(userData)
	if err != nil {
		panic(err.Error())
	}
	venus.Request("permissionServiceOP.findUserByUserName", "1", string(body))
```



