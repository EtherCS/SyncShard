# SyncShard
SyncShard is an availability-priority synchronization protocol for a blockchain sharding system

## test

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

### client parameter description

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