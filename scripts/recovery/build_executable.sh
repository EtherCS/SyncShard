# GOPATH=/Users/gudale/go
source ~/.bashrc
GOSRC=$GOPATH/src
ROOT=$GOSRC/github.com/EtherCS/SyncShard

mkdir -p build
go build -o build/recnode $ROOT/abci/cmd/recnode/example    # compiler node
# go build -o build/monitor $ROOT/cmd/recover/monitor_client # client for detection
go build -o build/user $ROOT/cmd/recover/user_client    # user for sending tx
go build -o build/latency $ROOT/cmd/recover/latency_client  # client for latency test

chmod +x build/*