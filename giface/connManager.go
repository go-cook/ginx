package giface

type ConnManager interface {
	Add(conn Conn)                   // 添加连接
	Remove(conn Conn)                //删除连接
	Get(connId uint32) (Conn, error) // 利用ConnId获取连接
	Count() int                      //获取当前连接数
	ClearConn()                      //删除所有连接
}
