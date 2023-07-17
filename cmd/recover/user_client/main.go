package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// var groupBlocks1, groupBlocks2 []types.Block
var groupBlocks1, groupBlocks2 tmbytes.HexBytes = nil, nil
var strBlockHash1, strBlockHash2 string = "", ""

// var (
// 	DETECT_HEIGHT int64 = 5
// )

var block_height int64
var shardPorts, shardIps string

// var isLeader bool

func init() {
	flag.Int64Var(&block_height, "block_height", 5, "block height for recovery detection")
	flag.StringVar(&shardPorts, "shardports", "20351,21351", "shards chain port")
	flag.StringVar(&shardIps, "shardips", "127.0.0.1,127.0.0.1", "shards chain ip")
}

func main() {
	flag.Parse()
	shard_ports_temp := []byte(shardPorts)
	shard_ports := bytes.Split(shard_ports_temp, []byte(","))
	shard_ips_temp := []byte(shardIps)
	shard_ips := bytes.Split(shard_ips_temp, []byte(","))
	// var ports_value64 []uint64
	// for _, shard_port := range shard_ports {
	// 	temp_port, _ := strconv.ParseUint(string(shard_port), 10, 64)
	// 	ports_value64 = append(ports_value64, temp_port)
	// }
	fmt.Println("Recovery height is", block_height)
	fmt.Println("connecting to node", string(shard_ips[0])+":"+string(shard_ports[0]))
	fmt.Println("connecting to node", string(shard_ips[1])+":"+string(shard_ports[1]))
	conn1, err1 := net.Dial("tcp", string(shard_ips[0])+":"+string(shard_ports[0]))
	conn2, err2 := net.Dial("tcp", string(shard_ips[1])+":"+string(shard_ports[1]))
	if err1 != nil || err2 != nil {
		fmt.Println("conn server failed")
		return
	}
	defer conn1.Close()
	defer conn2.Close()
	done1 := make(chan string)
	done2 := make(chan string)

	// wait 20s for requesting blocks
	// time.Sleep(20000 * time.Microsecond)

	go handleWriteGroupOne(conn1, done1)
	go handleReadGroupOne(conn1, done1)
	go handleWriteGroupTwo(conn2, done2)
	go handleReadGroupTwo(conn2, done2)

	// catch the signal
	fmt.Println(<-done1)
	fmt.Println(<-done2)
}

// Request format: 5 is the inconsistent height
// Get block: Rollback=0;Height=5;
// Send inconsistent evidence: Rollback=1;Height=5;blockA.Hash()=blockB.Hash()
// Note: blockA.Hash() means its chain is longer than blockB.Hash(), and don't need to rollback
func handleWriteGroupOne(conn net.Conn, done chan string) {
	for {
		_, e := conn.Write([]byte("Rollback=0;Height=" + strconv.Itoa(int(block_height)) + ";\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		time.Sleep(600000 * time.Millisecond)
	}
	done <- "Sent to group 1"
}
func handleReadGroupOne(conn net.Conn, done chan string) {
	for {
		// receive blocks
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		recv := string(buf[:reqLen])
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			break
		}
		// TODO: detection and construct evidence, write to send to node
		fmt.Println("block from group1 is: " + string(buf[:reqLen-1]))
		hashes := bytes.Split([]byte(recv), []byte(";"))
		if len(hashes) != 0 {
			strBlockHash1 = string(hashes[0])
			for strBlockHash2 == "" {
				// break
				// fmt.Println("Ether: wait for group2's block")
			}
			// fmt.Println("Ether: now, block1 is", strBlockHash1, "block2 is", strBlockHash2)
			if strBlockHash1 != strBlockHash2 {
				_, e := conn.Write([]byte("Rollback=1;Height=" + strconv.Itoa(int(block_height)) + ";" + strBlockHash1 + "=" + strBlockHash2 + ";\n"))
				if e != nil {
					fmt.Println("Error to send message because of ", e.Error())
					break
				}
				time.Sleep(600000 * time.Millisecond)
			}
		} else {
			// fmt.Println("Ether: receive a invalid format hash message", recv)
			continue
		}
	}

	done <- "Read from group 1"
}

func handleWriteGroupTwo(conn net.Conn, done chan string) {
	for {
		_, e := conn.Write([]byte("Rollback=0;Height=" + strconv.Itoa(int(block_height)) + ";\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		time.Sleep(600000 * time.Millisecond)
	}
	done <- "Sent to group 2"
}
func handleReadGroupTwo(conn net.Conn, done chan string) {
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		recv := string(buf[:reqLen])
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			break
		}
		// TODO: detection and construct evidence, write to send to node
		fmt.Println("block from group2 is: " + string(buf[:reqLen-1]))
		hashes := bytes.Split([]byte(recv), []byte(";"))
		if len(hashes) != 0 {
			strBlockHash2 = string(hashes[0])
			for strBlockHash1 == "" {
				// break
				// fmt.Println("Ether: wait for group1's block")
			}
			// fmt.Println("Ether: now, block1 is", strBlockHash1, "block2 is", strBlockHash2)
			if strBlockHash1 != strBlockHash2 {
				_, e := conn.Write([]byte("Rollback=1;Height=" + strconv.Itoa(int(block_height)) + ";" + strBlockHash1 + "=" + strBlockHash2 + ";\n"))
				if e != nil {
					fmt.Println("Error to send message because of ", e.Error())
					break
				}
				time.Sleep(600000 * time.Millisecond)
			}
		} else {
			// fmt.Println("Ether: receive a invalid format hash message", recv)
			continue
		}
	}
	done <- "Read from group 2"
}
func detectInconsistency(b1 tmbytes.HexBytes, b2 tmbytes.HexBytes) bool {
	return bytes.Equal(b1, b2)
}
