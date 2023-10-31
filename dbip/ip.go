package dbip

import (
	"github.com/sshterm/ipdb/ip"
	_ "embed"
)

//go:embed ip.db
var db []byte

func NewIP() *ip.IP {
	return ip.NewDBByte(db)
}
