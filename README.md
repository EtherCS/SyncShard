# SyncShard
SyncShard is an availability-priority synchronization protocol for a blockchain sharding system

## test

- run consensus nodes
```
./scripts/run_test.sh -n 2 -m 2 -p 10057 -i "127.0.0.1" -s "20057,21057" -x "127.0.0.2,127.0.0.3"
```
- run client
```
./build/user -shards 2 -beaconport 10057 -beaconip "127.0.0.1" -shardports "20057,21057" -shardips "127.0.0.1,127.0.0.1" -batch 10 -ratio 0.0
```