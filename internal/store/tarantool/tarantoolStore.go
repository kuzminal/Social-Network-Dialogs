package tarantool

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"sync"
)

type Tarantool struct {
	conn *tarantool.Connection
}

var (
	tarantoolInstance *Tarantool
	tarantoolOnce     sync.Once
)

func NewTarantoolMaster(host string, port string, user string, password string) (*Tarantool, error) {
	return getNewTarantoolInstance(host, port, user, password)
}

func NewTarantoolSlave(host string, port string, user string, password string) (*Tarantool, error) {
	return getNewTarantoolInstance(host, port, user, password)
}

func getNewTarantoolInstance(host string, port string, user string, password string) (*Tarantool, error) {
	if len(host) == 0 {
		host = "localhost"
	}
	if len(port) == 0 {
		port = "3301"
	}
	tarantoolOnce.Do(func() {
		opts := tarantool.Opts{
			User: user,
			Pass: password,
			/*Timeout:       2500 * time.Millisecond,
			Reconnect:     1 * time.Second,
			MaxReconnects: 3,*/
		}
		connStr := fmt.Sprintf("%s:%s", host, port)
		conn, err := tarantool.Connect(connStr, opts)
		if err != nil {
			fmt.Println("Connection refused:", err)
		}
		tarantoolInstance = &Tarantool{conn: conn}
	})
	return tarantoolInstance, nil
}

func (t *Tarantool) Ping(ctx context.Context) error {
	_, err := t.conn.Ping()
	return err
}

func (t *Tarantool) Close() error {
	return t.conn.Close()
}
