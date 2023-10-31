package cndb

import (
	"github.com/sshterm/ipdb/db/ip"
	_ "embed"
)

//go:embed ip.db
var db []byte

func NewIP() *ip.IP {
	return ip.NewDBByte(db)
}
