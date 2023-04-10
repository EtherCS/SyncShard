source ~/.bashrc
GOSRC=$GOPATH/src
ROOT=$GOSRC/github.com/EtherCS/SyncShard

mkdir -p build


go build -o build/appnode $ROOT/cmd/app
go build -o build/sync $ROOT/cmd/sync
go build -o build/user $ROOT/cmd/user

chmod +x build/*
