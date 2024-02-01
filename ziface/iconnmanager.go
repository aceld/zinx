// @Title iconnmanager.go
// @Description Connection management related operations, including adding, removing, getting a connection object by a connection ID, methods to get the current number of connections and clear all connections.
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
IConnManager Connection Management Abstract Layer
*/
type IConnManager interface {
	Add(IConnection)                                                        // Add connection
	Remove(IConnection)                                                     // Remove connection
	Get(uint64) (IConnection, error)                                        // Get a connection by ConnID
	Get2(string) (IConnection, error)                                       // Get a connection by string ConnID
	Len() int                                                               // Get current number of connections
	ClearConn()                                                             // Remove and stop all connections
	GetAllConnID() []uint64                                                 // Get all connection IDs
	GetAllConnIdStr() []string                                              // Get all string connection IDs
	Range(func(uint64, IConnection, interface{}) error, interface{}) error  // Traverse all connections
	Range2(func(string, IConnection, interface{}) error, interface{}) error // Traverse all connections 2
}
