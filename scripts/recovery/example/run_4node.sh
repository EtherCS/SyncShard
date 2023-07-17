#GOPATH=/home/gudale/go
GOSRC=$GOPATH/src
TEST_SCENE="example"
TM_HOME="$HOME/.example"
WORKSPACE="$GOSRC/github.com/EtherCS/SyncShard"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=90


rm -rf $TM_HOME/*
mkdir -p $TM_HOME
mkdir -p $LOG_DIR

cp -r $WORKSPACE/configs/recovery/4nodes/* $TM_HOME
echo "configs generated"

pkill -9 example
./build/node -home $TM_HOME/node0 &> $LOG_DIR/node0.log &
sleep 2
./build/node -home $TM_HOME/node1 &> $LOG_DIR/node1.log &
sleep 2
./build/node -home $TM_HOME/node2 &> $LOG_DIR/node2.log &
sleep 2
./build/node -home $TM_HOME/node3 &> $LOG_DIR/node3.log &
echo "recovery testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION
pkill -9 node
pkill -9 client
pkill -9 user
echo "all done"

#curl -s 'localhost:20057/broadcast_tx_commit?tx="abcd"'