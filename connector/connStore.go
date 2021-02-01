package connector

var connStore = make(map[string]*Connection)

// SaveConnection 保存连接
func SaveConnection(connection *Connection) {
	connStore[connection.uid] = connection
}

// GetConnection 获取连接
func GetConnection(uid string) *Connection {
	connection, ok := connStore[uid]
	if ok {
		return connection
	}
	return nil
}

// DelConnection 删除连接
func DelConnection(uid string) {
	delete(connStore, uid)
}
