package gnet

type BaseRouter struct{}

func (r *BaseRouter) PreHandle(Request)  {}
func (r *BaseRouter) Handle(Request)     {}
func (r *BaseRouter) PostHandle(Request) {}
