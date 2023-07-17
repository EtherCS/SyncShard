package recovery

import (
	"bytes"
	"encoding/binary"
	"log"

	dbm "github.com/tendermint/tm-db"
)

// implement the following interfaces:
// rollbackState

var (
	stateLatestKey     = []byte("stateLatestKey")
	stateCheckpointKey = []byte("stateCheckpointKey")
	kvPairPrefixKey    = []byte("kvPairKey:")

	ProtocolVersion uint64 = 0x1
)

const (
	valSetCheckpointInterval = 10 // how often to update a checkpoint
	valBufferBlockNum        = 10 // equals to valSetCheckpointInterval
)

type TxBlock struct {
	BlockHeight int64
	Txs         [][]byte
}

type RecoveryState struct {
	KVdb              dbm.DB
	OldKVdb           dbm.DB
	BufferBlocks      []TxBlock
	VoteBlock         TxBlock // VoteBlock is a temp block used for constructing BufferBlocks
	size              int64   // number of transactions
	height            int64
	AppHash           []byte
	checkpointHeight  int64 // set 10 blocks as a checkpoint
	checkpointSize    int64 // number of transactions until checkpoint height
	checkpointAppHase []byte
	StoreOptions
}

type StoreOptions struct {
	// DiscardABCIResponses determines whether or not the store
	// retains all ABCIResponses. If DiscardABCiResponses is enabled,
	// the store will maintain only the response object from the latest
	// height.
	DiscardABCIResponses bool
}

func NewRecoveryState(name string, dir string) *RecoveryState {
	var recoveryState RecoveryState
	var err, err2 error
	recoveryState.KVdb, err = dbm.NewDB(name, dbm.GoLevelDBBackend, dir+"/latest")
	if err != nil {
		log.Fatalf("Create latest database error: %v", err)
	}
	recoveryState.OldKVdb, err2 = dbm.NewDB(name, dbm.GoLevelDBBackend, dir+"/old")
	if err2 != nil {
		log.Fatalf("Create old database error: %v", err)
	}
	recoveryState.size = 0
	recoveryState.height = 0
	recoveryState.checkpointHeight = 0
	recoveryState.checkpointSize = 0
	// CommitCheckpointState(&recoveryState)
	return &recoveryState
}

func (state *RecoveryState) Size() int64 {
	return state.size
}
func (state *RecoveryState) Height() int64 {
	return state.height
}
func (state *RecoveryState) CheckpointHeight() int64 {
	return state.checkpointHeight
}
func (state *RecoveryState) CheckpointSize() int64 {
	return state.checkpointSize
}
func (state *RecoveryState) BeginBlock(blockHeight int64) {
	state.VoteBlock = TxBlock{BlockHeight: blockHeight}
}

func (state *RecoveryState) ExecuteTx(tx []byte) string {
	var key, value string

	parts := bytes.Split(tx, []byte("="))
	if len(parts) == 2 {
		key, value = string(parts[0]), string(parts[1])
	} else {
		key, value = string(tx), string(tx)
	}

	err := state.KVdb.Set(PrefixKey([]byte(key)), []byte(value))
	if err != nil {
		panic(err)
	}
	state.size++

	// add tx to VoteBlock
	state.VoteBlock.Txs = append(state.VoteBlock.Txs, tx)
	return key
}

func (state *RecoveryState) EndAndCommitBlock() {
	state.BufferBlocks = append(state.BufferBlocks, state.VoteBlock)
	state.VoteBlock.Txs = state.VoteBlock.Txs[:0]
}

func (state *RecoveryState) Commit() []byte {
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, state.size)
	state.AppHash = appHash
	state.height++
	state.EndAndCommitBlock()
	if len(state.BufferBlocks) == valBufferBlockNum {
		state.checkpointHeight = state.height
		state.checkpointSize = state.size
		state.checkpointAppHase = state.AppHash
		state.UpdateCheckpointState()
	}
	return appHash
}

func (state *RecoveryState) UpdateCheckpointState() {
	for _, block := range state.BufferBlocks {
		for _, tx := range block.Txs {
			var key, value string
			parts := bytes.Split(tx, []byte("="))
			if len(parts) == 2 {
				key, value = string(parts[0]), string(parts[1])
			} else {
				key, value = string(tx), string(tx)
			}

			state.OldKVdb.Set(PrefixKey([]byte(key)), []byte(value))
		}
	}
	state.BufferBlocks = state.BufferBlocks[:0] // clear buffer blocks
}

// rollback to the checkpoint state
func (state *RecoveryState) Rollback(height int64) {
	if state.checkpointHeight <= height {
		// copy(state.KVdb,oldState.KVdb)
		// TODO
		state.KVdb = state.OldKVdb
		state.size = state.checkpointSize
		state.height = state.checkpointHeight
		state.AppHash = state.checkpointAppHase
	}
	// TODO: the height of rollback is less than the checkpoint
}

func PrefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}
