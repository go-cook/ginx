package giface

// 封装路由
type Router interface {
	PreHandle(req Request)  // 业务前
	Handle(req Request)     // 业务中
	PostHandle(req Request) // 业务后
}
