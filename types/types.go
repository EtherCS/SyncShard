package types

import (
	"bytes"
	"net"
	"strconv"
)

type SyncAddress struct {
	Ip   net.IP
	Port uint16
}

type TransactionType struct {
	From_shard uint8 // the sender's shard
	To_shard   uint8 // the receiver's shard
	Tx_type    uint8
	From       []byte
	To         []byte
	Value      uint32
	Data       []byte
	Nonce      uint32 // TODO: enable contineous tx requests by setting vary nonce
}

func BytesToIp(bt []byte) net.IP {
	ip_parts := bytes.Split(bt, []byte("."))
	temp_p1, _ := strconv.ParseUint(string(ip_parts[0]), 10, 64)
	p1 := uint8(temp_p1)
	temp_p2, _ := strconv.ParseUint(string(ip_parts[1]), 10, 64)
	p2 := uint8(temp_p2)
	temp_p3, _ := strconv.ParseUint(string(ip_parts[2]), 10, 64)
	p3 := uint8(temp_p3)
	temp_p4, _ := strconv.ParseUint(string(ip_parts[3]), 10, 64)
	p4 := uint8(temp_p4)
	return []byte{p1, p2, p3, p4}
}
