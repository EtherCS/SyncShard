package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"
	// "time"
)

// ./build/latency -shardport "20057" -shardip "18.188.221.188"

var shardPort, shardIp, info string

func init() {
	flag.StringVar(&shardPort, "shardport", "20057", "shards chain port")
	flag.StringVar(&shardIp, "shardip", "127.0.0.1", "shards chain ip")
	flag.StringVar(&info, "info", "T-REC: performance, c=10", "info")
}

func main() {
	flag.Parse()
	fmt.Println("Test info:", info)
	start := time.Now()
	fmt.Println("start in:", start)
	request1 := fmt.Sprintf("http://%v:%v/broadcast_tx_commit?tx=\"abcd%v\"", shardIp, shardPort, get_rand(math.MaxInt64))
	res, _ := http.Get(request1)
	elapsed := time.Since(start)
	fmt.Println("receive", res)
	fmt.Println("confirmation latency is:", elapsed)
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
