package gnet

type BaseRouter struct{}

func (r *BaseRouter) PreHandle(req Request)  {}
func (r *BaseRouter) Handle(req Request)     {}
func (r *BaseRouter) PostHandle(req Request) {}
