package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net/http"
)

var shardPort, shardIp, beaconPort, beaconIp, blockSize string

func init() {
	flag.StringVar(&shardPort, "shardport", "20057", "shards chain port")
	flag.StringVar(&shardIp, "shardip", "127.0.0.1", "shards chain ip")
	flag.StringVar(&beaconPort, "beaconport", "10057", "beacon chain port")
	flag.StringVar(&beaconIp, "beaconip", "127.0.0.1", "beacon chain ip")
	flag.StringVar(&blockSize, "blocksize", "100", "beacon chain ip")
}

func main() {
	flag.Parse()
	http.Get(fmt.Sprintf("http://%v:%v/broadcast_tx_commit?tx=\"fromid=%v,toid=%v,type=%v,from=SYNC,to=WXYZ,value=%v,data=NONE,nonce=%v\"", shardIp, shardPort, 0, 1, 5, blockSize, get_rand(math.MaxInt32)))
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
