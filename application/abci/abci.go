package abci

import (
	"encoding/binary"

	// "log"
	// "time"

	// "log"
	"strconv"

	syncNode "github.com/EtherCS/SyncShard/application/validator"
	syntypes "github.com/EtherCS/SyncShard/types"
	abcicode "github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

var (
	// stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion uint64 = 0x1
)

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

var _ abcitypes.Application = (*syncApplication)(nil)

type syncApplication struct {
	abcitypes.BaseApplication
	Node *syncNode.ValidatorInterface
}

func NewsyncApplication(node *syncNode.ValidatorInterface) *syncApplication {
	return &syncApplication{
		Node: node,
	}
}

func (syncApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (syncApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *syncApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	return abcitypes.ResponseCheckTx{Code: abcicode.CodeTypeOK, GasWanted: 1}
}

func (app *syncApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

func (app *syncApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	_, tx_json := syncNode.ResolveTx(req.Tx)
	var err1, err2 error
	var events []abcitypes.Event
	var event_type string
	new_tx := syntypes.TransactionType{
		From_shard: tx_json.From_shard,
		To_shard:   tx_json.To_shard,
		Tx_type:    tx_json.Tx_type,
		From:       tx_json.From,
		To:         tx_json.To,
		Value:      tx_json.Value,
		Data:       tx_json.Data,
	}
	if tx_json.Tx_type == syncNode.IntraShard_TX {
		event_type = "intra-shard transaction"
		err1 = app.Node.BCState.Database.Set(prefixKey(tx_json.From), []byte("0"))
		err2 = app.Node.BCState.Database.Set(prefixKey(tx_json.To), []byte("0"))
		app.Node.BCState.Size++
		// count frequncy
		value_from, ok_from := app.Node.KeyFrequency[string(tx_json.From)]
		if ok_from { // already exists
			app.Node.KeyFrequency[string(tx_json.From)] = value_from + 1
		} else { // first appear
			app.Node.KeyFrequency[string(tx_json.From)] = 0
		}
		value_to, ok_to := app.Node.KeyFrequency[string(tx_json.To)]
		if ok_to { // already exists
			app.Node.KeyFrequency[string(tx_json.To)] = value_to + 1
		} else { // first appear
			app.Node.KeyFrequency[string(tx_json.To)] = 0
		}
	} else if tx_json.Tx_type == syncNode.InterShard_TX_Verify {
		event_type = "inter-shard verification transaction"
		err1 = app.Node.BCState.Database.Set(prefixKey(tx_json.From), []byte("lock"))
		value_from, ok_from := app.Node.KeyFrequency[string(tx_json.From)]
		if ok_from { // already exists
			app.Node.KeyFrequency[string(tx_json.From)] = value_from + 1
		} else { // first appear
			app.Node.KeyFrequency[string(tx_json.From)] = 0
		}
		new_tx.Tx_type = syncNode.InterShard_TX_Execute
		_, exec_tx := syncNode.Deserilization(new_tx)
		if app.Node.Leader {
			go app.Node.DeliverExecutionTx(exec_tx, new_tx.To_shard)
		}

	} else if tx_json.Tx_type == syncNode.InterShard_TX_Execute {
		event_type = "inter-shard execution transaction"
		err2 = app.Node.BCState.Database.Set(prefixKey(tx_json.To), []byte("lock"))
		value_to, ok_to := app.Node.KeyFrequency[string(tx_json.To)]
		if ok_to { // already exists
			app.Node.KeyFrequency[string(tx_json.To)] = value_to + 1
		} else { // first appear
			app.Node.KeyFrequency[string(tx_json.To)] = 0
		}
		new_tx.Tx_type = syncNode.InterShard_TX_Commit
		_, commit_tx := syncNode.Deserilization(new_tx)
		if app.Node.Leader {
			go app.Node.DeliverCommitTx(commit_tx, new_tx.From_shard)
		}
	} else if tx_json.Tx_type == syncNode.InterShard_TX_Commit {
		event_type = "inter-shard commit transaction"
		err1 = app.Node.BCState.Database.Set(prefixKey(tx_json.From), []byte("0"))
		new_tx.Tx_type = syncNode.InterShard_TX_Update
		_, update_tx := syncNode.Deserilization(new_tx)
		if app.Node.Leader {
			go app.Node.DeliverUpdateTx(update_tx, new_tx.To_shard)
		}
	} else if tx_json.Tx_type == syncNode.InterShard_TX_Update {
		event_type = "inter-shard update transaction"
		err2 = app.Node.BCState.Database.Set(prefixKey(tx_json.To), []byte("0"))
		app.Node.BCState.Size++
	} else if tx_json.Tx_type == syncNode.Synchronization_TX {
		// print simulation latency
		syntypes.PrintSyncLatency(app.Node.KeyFrequency, int(app.Node.ShowKeyNum), int(app.Node.BCState.Height+1))
		// print frequency information
		syntypes.PrintKeyFrequency(app.Node.KeyFrequency, int(app.Node.ShowKeyNum), int(app.Node.BCState.Height+1))
	}
	if err1 != nil || err2 != nil {
		panic(err1)
	}
	events = []abcitypes.Event{
		{
			Type: event_type,
			Attributes: []abcitypes.EventAttribute{
				{Key: "from", Value: string(tx_json.From), Index: true},
				{Key: "to", Value: string(tx_json.To), Index: true},
				{Key: "value", Value: strconv.Itoa(int(tx_json.Value)), Index: true},
				{Key: "data", Value: string(tx_json.Data), Index: true},
			},
		},
	}
	return abcitypes.ResponseDeliverTx{Code: abcicode.CodeTypeOK, Events: events}
}

func (app *syncApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *syncApplication) Commit() abcitypes.ResponseCommit {
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.Node.BCState.Size)
	app.Node.BCState.AppHash = appHash
	app.Node.BCState.Height++
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *syncApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	if reqQuery.Prove {
		value, err := app.Node.BCState.Database.Get(prefixKey(reqQuery.Data))
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
		resQuery.Height = app.Node.BCState.Height

		return
	}

	resQuery.Key = reqQuery.Data
	value, err := app.Node.BCState.Database.Get(prefixKey(reqQuery.Data))
	if err != nil {
		panic(err)
	}
	if value == nil {
		resQuery.Log = "does not exist"
	} else {
		resQuery.Log = "exists"
	}
	resQuery.Value = value
	resQuery.Height = app.Node.BCState.Height

	return resQuery
}
