package utils

import (
	"encoding/json"
	"io/ioutil"

	"zinx/ziface"
)

// 全局参数的结构体， 供其他模块使用
// 参数通过zinx.json 配置

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer //  当前zinx 全局的Server 对象
	Host      string         //  当前服务器主机监听的IP
	TcpPort   int            //  当前服务器主机监听的端口
	Name      string         //  当前服务器的名称

	// Zinx
	Version        string // 当前Zinx 的版本号
	MaxConn        int    // 允许的最大连接数
	MaxPackageSize uint32 //  数据包的最大值
}

// 从 zinx.json 加载用户的配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 将 json 数据解析到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 定一个全局的 GlobalObj
var GlobalObject *GlobalObj

// init GlobalObject
func init() {

	// 配置文件没有加载，有一组默认的初始值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        9999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1024,
		MaxPackageSize: 4096,
	}

	// 读取用户配置
	GlobalObject.Reload()
}
