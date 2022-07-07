package giface

// 定义服务接口
type Server interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行服务
	Server()
}
