package common

import (
	"net"
)

type P2pclient struct {
	Name string
	Conn net.Conn
}
