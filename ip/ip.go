package ip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
)

type IP struct {
	reader  io.ReaderAt
	version int64
	ipV4    int64
	ipV6    int64
	country map[byte]string
}

func NewDB(reader io.ReaderAt) *IP {
	ip := &IP{reader: reader}
	ip.init()
	return ip
}
func NewDBByte(data []byte) *IP {
	ip := &IP{reader: bytes.NewReader(data)}
	ip.init()
	return ip
}
func (i *IP) init() {
	i.initIndex()
	i.initCountry()
}
func (i *IP) Version() int64 {
	return i.version
}
func (i *IP) CountryCN(ip net.IP) bool {
	return i.Country(ip) == "CN"
}
func (i *IP) Country(ip net.IP) string {
	if ip.To4() != nil {
		return i.country4(ip)
	} else if ip.To16() != nil {
		return i.country16(ip)
	}
	return ""
}
func (i *IP) country4(ip net.IP) (country string) {
	return i.lookup(ip, 12, i.ipV4+14, 6, 4)
}
func (i *IP) country16(ip net.IP) (country string) {
	return i.lookup(ip, i.ipV4+12, i.ipV4+i.ipV6+12, 18, 16)
}
func (i *IP) initIndex() {
	if i.version == 0 {
		_index := make([]byte, 12)
		_, err := i.reader.ReadAt(_index, 0)
		if err == nil {
			i.version = new(big.Int).SetBytes(_index[:4]).Int64()
			i.ipV4 = new(big.Int).SetBytes(_index[4:8]).Int64()
			i.ipV6 = new(big.Int).SetBytes(_index[8:12]).Int64()
		}
	}
}
func (i *IP) lookup(ip net.IP, off, max, size int64, len int) (country string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			if off >= max {
				break
			}
			data := make([]byte, size)
			_, err := i.reader.ReadAt(data, off)
			if err != nil {
				break
			}
			off += size
			if len == 16 {
				if data[0] != ip.To16()[0] {
					continue
				}
			} else if len == 4 {
				if data[0] != ip.To4()[0] {
					continue
				}
			}
			_, i2, err := net.ParseCIDR(fmt.Sprintf("%s/%d", net.IP(data[:len]), data[len]))
			if err != nil {
				break
			}
			if i2.Contains(ip) {
				if v, ok := i.country[data[len+1]]; ok {
					country = v
				}
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
	return
}
func (i *IP) initCountry() {
	if len(i.country) == 0 {
		i.country = make(map[byte]string)
		_index := make([]byte, 4)
		off := i.ipV4 + i.ipV6 + 12
		_, err := i.reader.ReadAt(_index, off)
		if err == nil {
			size := new(big.Int).SetBytes(_index[:4]).Int64()
			data := make([]byte, int(size))
			_, err := i.reader.ReadAt(data, off+4)
			if err == nil {
				_ = json.Unmarshal(data, &i.country)
			}
		}
	}
}
