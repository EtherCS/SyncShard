DURATION=120
GROUP_ID=1
NODE_ID=0
CONFIG_HOME="./configs/recovery/byzantine/group1"

while getopts ":g:n:h:" opt
do 
    case $opt in
    g) # 
        echo "group id is $OPTARG"
        GROUP_ID=$OPTARG
        ;;
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

#GOPATH=/home/gudale/go
GOSRC=$GOPATH/src
TEST_SCENE="example"
TM_HOME="$HOME/.example"
WORKSPACE="$GOSRC/github.com/EtherCS/recovery"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR_ROOT="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE-group$GROUP_ID"

TM_HOME_GROUP="$TM_HOME/group$GROUP_ID"
# LOG_DIR_GROUP="$LOG_DIR_ROOT/group$GROUP_ID"

rm -rf $TM_HOME/*
mkdir -p $TM_HOME_GROUP
mkdir -p $LOG_DIR_ROOT
# mkdir -p $LOG_DIR_GROUP

cp -r $CONFIG_HOME/* $TM_HOME
echo "configs generated"

echo "recovery testnet launched"
echo "running for ${DURATION}s..."

pkill -9 recnode

# run nodes
echo "run group$GROUP_ID node$NODE_ID"
./build/recnode -home $TM_HOME/node$NODE_ID &> $LOG_DIR_ROOT/node$NODE_ID.log &

sleep $DURATION
pkill -9 recnode
echo "all done"

#curl -s 'localhost:20057/broadcast_tx_commit?tx="abcd"'