package conn_pool

import (
	"database/sql"
)

// A database connection
type Conn interface {
	Query(string) (*sql.Rows, error)
}

// Use this struct for storing pool related information
type Pool struct {
}

// NewPool creates a new connection pool with connections to a database.
// It takes as input a function that actually establishes connections with
// the underlying database.
// The returned pool should have one connection ready and available for use
// Establishing a connection to a database might fail, in which case NewPool
// propagates the error from the underlying driver to its caller.
func NewPool(newConnection func() (Conn, error)) (*Pool, error) {
	return &Pool{}, nil
}

// SetMaxConnections sets the maximum number of connections that a pool can
// maintain. If it sets to a number m that is smaller than the number of currently
// open connections, SetMaxConnections should block (i.e., not return) until the
// number of open connections drops below m.
func (p *Pool) SetMaxConnections(m int) {
}

// GetMaxConnections returns the maximum number of connections that the pool can
// maintain.
func (p *Pool) GetMaxConnections() int {
	return 0
}

// Open returns a connection from the connection pool. It should only return a
// connection when the number of open connections is less than the maximum pool size.
// Open should block (i.e., not return) when the number of open connections
// exceeds the pool's maximum.
// This function needs to be safe for concurrent use. When it is called by
// multiple goroutines, it should return unique connections to each caller.
func (p *Pool) Open() Conn {
	return nil
}

// Close returns a connection back to the connection pool without actually closing
// it with the underlying database.
// When a connection is closed, using it has undefined behavior (i.e., applications
// should not do it).
// Close needs to return a connection to the pool it was originally allocated in.
// If the connection does not belong to the pool (i.e., c was not allocated to p
// originally), Close is a noop (i.e., no effects).
// Closing a connection that's not open is a noop.
func (p *Pool) Close(c Conn) {
}
