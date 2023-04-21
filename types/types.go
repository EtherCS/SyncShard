package types

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strconv"
)

const (
	Addr_Length uint8 = 6
	Data_Length uint8 = 4
	Key_Num     uint8 = 100
)

type SyncAddress struct {
	Ip   net.IP
	Port uint16
}

type TransactionType struct {
	From_shard uint8  // the sender's shard ID
	To_shard   uint8  // the receiver's shard ID
	Tx_type    uint8  // transaction request or synchronization request
	From       []byte // the address of the transaction's sender
	To         []byte // the address of the transaction's receiver
	Value      uint32
	Data       []byte // the data used to execute the transaction
	Nonce      uint32
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

type SortMap struct {
	Key   string
	Value int
}

type SortMapList []SortMap

func (p SortMapList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p SortMapList) Len() int {
	return len(p)
}
func (p SortMapList) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

func sortMapByValue(m map[string]int) SortMapList {
	p := make(SortMapList, len(m))
	i := 0
	for k, v := range m {
		p[i] = SortMap{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

func PrintKeyFrequency(m map[string]int, showNum int, height int) {
	var sortMap SortMapList = sortMapByValue(m)
	// var count int = 0
	fmt.Println("Key Ratio")
	var ratio float64 = 0.0
	var total_num int = 0
	for i := 0; i < len(sortMap); i++ {
		if i >= showNum {
			break
		}
		total_num += sortMap[i].Value
	}
	for i := 0; i < len(sortMap); i++ {
		if i >= showNum {
			fmt.Printf("Total ratio is %f \n", ratio)
			break
		}
		fmt.Printf("%s %f \n", sortMap[i].Key, float64(sortMap[i].Value)/float64(total_num))
		if i <= showNum/2 {
			ratio += float64(sortMap[i].Value) / float64(total_num)
		}
	}
}

func PrintSyncLatency(m map[string]int, showNum int, blockNum int, blockSize int) {
	var sortMap SortMapList = sortMapByValue(m)
	var ratio float64 = 0.0
	var total_num int = 0
	for i := 0; i < len(sortMap); i++ {
		if i >= showNum {
			break
		}
		total_num += sortMap[i].Value
	}
	for i := 0; i < len(sortMap); i++ {
		if i >= showNum {
			break
		}
		if i <= showNum/2 {
			ratio += float64(sortMap[i].Value) / float64(total_num)
		}
	}
	fmt.Printf("Synchronization latency with SyncShard: %f s \n", 1*ratio+(1.0-ratio)*float64(blockNum)*(float64(blockSize)/100.0))
	fmt.Printf("Synchronization latency with traditional method: %f s\n", float64(blockNum)*(float64(blockSize)/100.0))
}
