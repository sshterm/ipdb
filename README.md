# ipdb
IP数据库，极速返回 IP 所在国家或地区,支持 IP4和 IPV6
数据库 约2.6 MB

```go
package main

import (
	"github.com/sshterm/ipdb/db"
	"log"
	"net"
)

func main() {
	data := db.NewIP()
	log.Println(data.Version()) //1698710400
	log.Println(data.Country(net.ParseIP("120.231.109.110"))) //CN
	log.Println(data.Country(net.ParseIP("2409:8a55:f2fb:ff2:42:c0ff:fea8:a0a"))) //CN
}

```
