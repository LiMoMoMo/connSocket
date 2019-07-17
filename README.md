## connC
连接客户端。
### 接口
1. NewConnC(port string, addr string)
返回客户端实例
2. Write(Report)
写消息到服务端
3. GetCmdChan() chan Command
获取接收消息的管道
## connS
连接服务端。
### 接口
1. NewConnS(port string, addr string)
返回服务端实例
2. Start()
开始服务
3. WriteTo(id string, Command)
指定给哪个客户端写消息
4. GetRepChan() chan Report
获取接收消息的管道
## models
仿照`beego`框架的模型, 将所有数据集合到一个目录下。
1. Report
客户端上报到服务端的消息类型
2. Command
服务端发送指定到客户端的消息类型
### 发送消息
根消息使用字段`Content`(interface{})类型接收所有类型结构体
```go
// client
addr := models.Report{
	Type: models.Type_AddressInfo,
	Content: models.AddressInfo{
		ID:        "qwerty",
		Addresses: []string{"0.0.0.0", "1.1.1.1"},
	},
}
connc.Write(&addr)
```
### 接收消息
根消息使用字段`Content`(interface{})类型返回所有消息
```go
// client
go func() {
	for {
		select {
		case cmd := <-connc.GetCmdChan():
			fmt.Println("Receive Msg", cmd.Content.(models.AddressInfo))
		case <-ctx.Done():
			return
		}
	}
}()
```
## example
使用例子

## feature
1. 解决`TCP`流式传输中的粘包/半包的问题;
2. 提供简单的接口；