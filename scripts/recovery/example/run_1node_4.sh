#GOPATH=/home/gudale/go
GOSRC=$GOPATH/src
TEST_SCENE="example"
TM_HOME="$HOME/.example"
WORKSPACE="$GOSRC/github.com/EtherCS/SyncShard"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=90

NODE_ID=0
while getopts ":n:h:" opt
do 
    case $opt in
    n) # 
        echo "node is is $OPTARG"
        NODE_ID=$OPTARG
        ;;  
    h)
        echo "config home is $OPTARG"
        CONFIG_HOME=$OPTARG
        ;;  
    ?)  
        echo "unknown: $OPTARG"
        ;;
    esac
done

rm -rf $TM_HOME/*
mkdir -p $TM_HOME
mkdir -p $LOG_DIR

cp -r $WORKSPACE/configs/recovery/4nodes/* $TM_HOME
echo "configs generated"

pkill -9 recnode
./build/recnode -home $TM_HOME/node$NODE_ID &> $LOG_DIR/node$NODE_ID.log &
sleep 2

echo "recovery testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION
pkill -9 recnode

echo "all done"

#curl -s 'localhost:20057/broadcast_tx_commit?tx="abcd"'