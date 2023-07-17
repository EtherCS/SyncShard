#GOPATH=/home/gudale/go
source ~/.bashrc
GOSRC=$GOPATH/src
TEST_SCENE="example"
TM_HOME="$HOME/.example"
WORKSPACE="$GOSRC/github.com/EtherCS/SyncShard"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR_ROOT="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
LOG_DIR1="$LOG_DIR_ROOT/group1"
LOG_DIR2="$LOG_DIR_ROOT/group2"
TM_HOME_GROUP1="$TM_HOME/group1"
TM_HOME_GROUP2="$TM_HOME/group2"
DURATION=120

echo "group1 data in $LOG_DIR1"
echo "group2 data in $LOG_DIR2"

rm -rf $TM_HOME/*
mkdir -p $TM_HOME_GROUP1
mkdir -p $TM_HOME_GROUP2
mkdir -p $LOG_DIR1
mkdir -p $LOG_DIR2

cp -r $WORKSPACE/configs/recovery/byzantine-single/group1/* $TM_HOME_GROUP1
cp -r $WORKSPACE/configs/recovery/byzantine-single/group2/* $TM_HOME_GROUP2
echo "configs generated"

echo "recovery testnet launched"
echo "running for ${DURATION}s..."

pkill -9 recnode
pkill -9 monitor
pkill -9 user
# run nodes
echo "run tendermint nodes"
./build/recnode -home $TM_HOME/group1/node0 &> $LOG_DIR1/node0.log &
./build/recnode -home $TM_HOME/group2/node0 &> $LOG_DIR2/node0.log &
sleep 2
./build/recnode -home $TM_HOME/group1/node1 &> $LOG_DIR1/node1.log &
./build/recnode -home $TM_HOME/group2/node1 &> $LOG_DIR2/node1.log &
sleep 2
./build/recnode -home $TM_HOME/group1/node2 &> $LOG_DIR1/node2.log &
./build/recnode -home $TM_HOME/group2/node2 &> $LOG_DIR2/node2.log &
sleep 2

# run user for sending transactions
echo "run user, sending transactions"
./build/user &
sleep 20

# run client to get block, and trigger rollback, default 5
echo "run client, recovery detection"
./build/monitor -block_height 12 &


sleep $DURATION
pkill -9 recnode
pkill -9 monitor
pkill -9 user
echo "all done"

#curl -s 'localhost:20057/broadcast_tx_commit?tx="abcd"'