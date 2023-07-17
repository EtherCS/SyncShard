package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"
	// "time"
)

var shardNum, reqDuration, concurrentNum uint
var shardPorts, shardIps string

// group 1: ./build/user -parallel 200 -shardports "20057,20157,20257" -shardips "18.188.221.188,3.145.146.41,13.59.36.127"
// two groups: ./build/user -parallel 200 -shardports "20057,20157,20257,21057,21157,21357" -shardips "18.188.221.188,3.145.146.41,13.59.36.127,18.188.221.188,3.145.146.41,13.59.36.127"
// ./build/user -parallel 200 -shardports "20057,20157,20257,21057,21157,21357" -shardips "18.188.221.188,3.145.146.41,13.59.36.127,3.139.65.78,3.19.238.75,18.221.138.47"
func init() {
	flag.UintVar(&shardNum, "shards", 2, "the number of shards")
	flag.UintVar(&reqDuration, "duration", 120, "duration of sending request")
	flag.StringVar(&shardPorts, "shardports", "20057,20157,20257", "shards chain port")
	flag.StringVar(&shardIps, "shardips", "127.0.0.1,127.0.0.1,127.0.0.1", "shards chain ip")
	flag.UintVar(&concurrentNum, "parallel", 100, "concurrent number for sending requests")
}

func main() {
	flag.Parse()
	shard_ports_temp := []byte(shardPorts)
	shard_ports := bytes.Split(shard_ports_temp, []byte(","))
	shard_ips_temp := []byte(shardIps)
	shard_ips := bytes.Split(shard_ips_temp, []byte(","))
	var ports_value64 []uint64
	for _, shard_port := range shard_ports {
		temp_port, _ := strconv.ParseUint(string(shard_port), 10, 64)
		ports_value64 = append(ports_value64, temp_port)
	}

	for p := 0; p < int(concurrentNum); p++ {
		for i, _ := range shard_ports {
			go send_request(uint(ports_value64[i]), string(shard_ips[i]))
		}
		// time.Sleep(time.Duration(requestRate) * time.Millisecond)
	}
	time.Sleep(time.Duration(reqDuration) * time.Second)
	// for i := 0; i < 100; i++ {
	// 	// for true {
	// 	request1 := fmt.Sprintf("http://127.0.0.1:20057/broadcast_tx_commit?tx=\"abcd%v\"", get_rand(math.MaxInt64))
	// 	request2 := fmt.Sprintf("http://127.0.0.1:21057/broadcast_tx_commit?tx=\"abcd%v\"", get_rand(math.MaxInt64))
	// 	go send_request("127.0.0.1")
	// 	go http.Get(request2)
	// 	// time.Sleep(100 * time.Millisecond)
	// }
}

func get_rand(upperBond int64) string {
	maxInt := new(big.Int).SetInt64(upperBond)
	i, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		fmt.Printf("Can't generate random value: %v, %v", i, err)
	}
	outputRand := fmt.Sprintf("%v", i)
	return outputRand
}

func send_request(port uint, ip string) {
	for {
		http.Get(fmt.Sprintf("http://%v:%v/broadcast_tx_commit?tx=\"abcd%v\"", ip, port, get_rand(math.MaxInt64)))
	}
}
