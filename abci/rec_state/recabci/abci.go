package recabci

import (
	"fmt"
	"log"
	"time"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"

	recstate "github.com/EtherCS/SyncShard/abci/rec_state/recovery"
)

//---------------------------------------------------

// declare RecoveryApplication is an implementation of Application Interface
var _ types.Application = (*RecoveryApplication)(nil)

type RecoveryApplication struct {
	types.BaseApplication

	recState     *recstate.RecoveryState
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)
	txToRemove   map[string]struct{}
}

func NewRecoveryApplication(db *recstate.RecoveryState) *RecoveryApplication {
	// state := recstate.LoadState(db)
	return &RecoveryApplication{
		recState:     db,
		RetainBlocks: 0,
	}
}

func (app *RecoveryApplication) State() *recstate.RecoveryState {
	if app.recState == nil {
		panic("the state of application is nil")
	}
	return app.recState
}

func (app *RecoveryApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.recState.Size()),
		Version:          version.ABCIVersion,
		AppVersion:       recstate.ProtocolVersion,
		LastBlockHeight:  app.recState.Height(),
		LastBlockAppHash: app.recState.AppHash,
	}
}

func (app *RecoveryApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	log.Println("abci: begin a block")
	fmt.Println("Ether recovery test: start block at time", time.Now())
	app.txToRemove = map[string]struct{}{}
	app.recState.BeginBlock(req.Header.Height)
	return types.ResponseBeginBlock{}
}

// tx is either "key=value" or just arbitrary bytes
func (app *RecoveryApplication) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	// log.Println("abci: deliver a transaction")
	key := app.recState.ExecuteTx(req.Tx)
	events := []types.Event{
		{
			Type: "app",
			Attributes: []types.EventAttribute{
				{Key: "creator", Value: "Cosmoshi Netowoko", Index: true},
				{Key: "key", Value: key, Index: true},
				{Key: "index_key", Value: "index is working", Index: true},
				{Key: "noindex_key", Value: "index is working", Index: false},
			},
		},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

func (app *RecoveryApplication) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	if req.Type == types.CheckTxType_Recheck {
		if _, ok := app.txToRemove[string(req.Tx)]; ok {
			return types.ResponseCheckTx{Code: 5, GasWanted: 1}
		}
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *RecoveryApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	log.Println("abci: end a block")
	return types.ResponseEndBlock{}
}

func (app *RecoveryApplication) Commit() types.ResponseCommit {
	log.Println("abci: commit a block")
	// Using a memdb - just return the big endian size of the db
	appHash := app.recState.Commit()
	// app.recState.EndAndCommitBlock()
	resp := types.ResponseCommit{Data: appHash}
	if app.RetainBlocks > 0 && app.recState.Height() >= app.RetainBlocks {
		resp.RetainHeight = app.recState.Height() - app.RetainBlocks + 1
	}
	fmt.Println("Ether recovery test: commit block at time", time.Now())
	return resp
}

// Returns an associated value or nil if missing.
func (app *RecoveryApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	if reqQuery.Prove {
		value, err := app.recState.KVdb.Get(recstate.PrefixKey(reqQuery.Data))
		if err != nil {
			panic(err)
		}
		if value == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		resQuery.Height = app.recState.Height()

		return
	}

	resQuery.Key = reqQuery.Data
	value, err := app.recState.KVdb.Get(recstate.PrefixKey(reqQuery.Data))
	if err != nil {
		panic(err)
	}
	if value == nil {
		resQuery.Log = "does not exist"
	} else {
		resQuery.Log = "exists"
	}
	resQuery.Value = value
	resQuery.Height = app.recState.Height()

	return resQuery
}

// func (app *RecoveryApplication) BeginRollback(req types.RequestBeginRollback) types.ResponseBeginRollback {
// 	log.Println("abci: begin rollback, checkpoint is ", app.recState.CheckpointHeight())
// 	app.recState.Rollback(req.From)
// 	return types.ResponseBeginRollback{CheckpointHeight: app.recState.CheckpointHeight()}
// }

// func (app *RecoveryApplication) ReplayRollbackTx(req types.RequestReplayRollbackTx) types.ResponseReplayRollbackTx {
// 	app.recState.ExecuteTx(req.Tx)
// 	return types.ResponseReplayRollbackTx{}
// }

// func (app *RecoveryApplication) EndRollback(req types.RequestEndRollback) types.ResponseEndRollback {
// 	// app.recState.Commit()
// 	log.Println("abci: end rollback")
// 	return types.ResponseEndRollback{}
// }
