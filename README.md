# SyncShard
SyncShard is an availability-priority synchronization protocol for a blockchain sharding system

## Architecture
![alt text](/figure/architecture.png)
## Download dependency
```
go mod tidy
```
## Compile
```
make build
```
## Test
## node parameter description

| node parameters | description        | default               |
|-----------------|--------------------|-----------------------|
| n               | shard number       | 2                     |
| m               | shard size         | 2                     |
| p               | beacon chain port  | 10057                 |
| i               | beacon chain ip    | "127.0.0.1"           |
| s               | shard chain points | "20057,21057"         |
| x               | shard chain ips    | "127.0.0.1,127.0.0.1" |
| k               | the number of key  | 100 |

### transaction client parameter description

| client parameters | description               | default               |
|-------------------|---------------------------|-----------------------|
| shards            | shard number              | 2                     |
| beaconport        | beacon chain port         | 10057                 |
| beaconip          | beacon chain ip           | "127.0.0.1"           |
| shardports        | shard chain points        | "20057,21057"         |
| shardips          | shard chain ips           | "127.0.0.2,127.0.0.3" |
| batch             | batch size per request    | 10                    |
| ratio             | cross-shard txs ratio     | 0.8                   |
| parallel          | concurrent request number | 100                   |
| duration          | execution duration, s     | 120                   |

### synchronization client parameter description

| client parameters | description               | default               |
|-------------------|---------------------------|-----------------------|
| shards            | shard number              | 2                     |
| beaconport        | beacon chain port         | 10057                 |
| beaconip          | beacon chain ip           | "127.0.0.1"           |
| shardports        | shard chain points        | "20057,21057"         |
| shardips          | shard chain ips           | "127.0.0.2,127.0.0.3" |
| blocksize         | block size                | 100                    |


### 1. Single machine test

- run consensus nodes
```
./scripts/run_test.sh -n 2 -m 2 -p 10057 -i "127.0.0.1" -s "20057,21057" -x "127.0.0.2,127.0.0.3" -k 100
```
- run client
```
./build/user -shards 2 -beaconport 10057 -beaconip "127.0.0.1" -shardports "20057,21057" -shardips "127.0.0.1,127.0.0.1" -batch 10 -ratio 0.0
```

- run synchronization request
```
./build/sync -shardport "20057" -shardip "127.0.0.1"
```

After testing finishes, you can get the log files of each node in: ./tmplog/app-$TEST-TIME

### 2. Multiple machines test (different IPs)
**Very simple! Just replace "127.0.0.1" with the machine IPs, and run the script on different machines respectively. For example, assume two machines have #IP1 and #IP2**
- run consensus nodes on two machines respectively
```
./scripts/run_test.sh -n 2 -m 2 -p 10057 -i "127.0.0.1" -s "20057,21057" -x "#IP1,#IP2" -k 100
```
- run client
```
./build/user -shards 2 -beaconport 10057 -beaconip "127.0.0.1" -shardports "20057,21057" -shardips "#IP1,#IP2" -batch 10 -ratio 0.0
```

- run synchronization request
```
./build/sync -shardport "20057" -shardip "#IP1"
```

After testing finishes, you can get the log files of each node in: ./tmplog/app-$TEST-TIME