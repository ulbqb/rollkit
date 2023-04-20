package abcicli

import (
	"sync"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/service"
	tmsync "github.com/tendermint/tendermint/libs/sync"

	rolltypes "github.com/rollkit/rollkit/abci/types"
)

type Client interface {
	service.Service

	SetResponseCallback(Callback)
	Error() error

	FlushAsync() *ReqRes
	EchoAsync(msg string) *ReqRes
	InfoAsync(types.RequestInfo) *ReqRes
	SetOptionAsync(types.RequestSetOption) *ReqRes
	DeliverTxAsync(types.RequestDeliverTx) *ReqRes
	CheckTxAsync(types.RequestCheckTx) *ReqRes
	QueryAsync(types.RequestQuery) *ReqRes
	CommitAsync() *ReqRes
	InitChainAsync(types.RequestInitChain) *ReqRes
	BeginBlockAsync(types.RequestBeginBlock) *ReqRes
	EndBlockAsync(types.RequestEndBlock) *ReqRes
	ListSnapshotsAsync(types.RequestListSnapshots) *ReqRes
	OfferSnapshotAsync(types.RequestOfferSnapshot) *ReqRes
	LoadSnapshotChunkAsync(types.RequestLoadSnapshotChunk) *ReqRes
	ApplySnapshotChunkAsync(types.RequestApplySnapshotChunk) *ReqRes

	FlushSync() error
	EchoSync(msg string) (*types.ResponseEcho, error)
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
	SetOptionSync(types.RequestSetOption) (*types.ResponseSetOption, error)
	DeliverTxSync(types.RequestDeliverTx) (*types.ResponseDeliverTx, error)
	CheckTxSync(types.RequestCheckTx) (*types.ResponseCheckTx, error)
	QuerySync(types.RequestQuery) (*types.ResponseQuery, error)
	CommitSync() (*types.ResponseCommit, error)
	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)
	BeginBlockSync(types.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
	ListSnapshotsSync(types.RequestListSnapshots) (*types.ResponseListSnapshots, error)
	OfferSnapshotSync(types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error)
	LoadSnapshotChunkSync(types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error)
	ApplySnapshotChunkSync(types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error)

	GetAppHashSync(rolltypes.RequestGetAppHash) (*rolltypes.ResponseGetAppHash, error)
	GenerateFraudProofSync(rolltypes.RequestGenerateFraudProof) (*rolltypes.ResponseGenerateFraudProof, error)
	VerifyFraudProofSync(rolltypes.RequestVerifyFraudProof) (*rolltypes.ResponseVerifyFraudProof, error)
}

type Callback func(*rolltypes.Request, *rolltypes.Response)

type ReqRes struct {
	*rolltypes.Request
	*sync.WaitGroup
	*rolltypes.Response // Not set atomically, so be sure to use WaitGroup.

	mtx tmsync.Mutex

	// callbackInvoked as a variable to track if the callback was already
	// invoked during the regular execution of the request. This variable
	// allows clients to set the callback simultaneously without potentially
	// invoking the callback twice by accident, once when 'SetCallback' is
	// called and once during the normal request.
	callbackInvoked bool
	cb              func(*rolltypes.Response) // A single callback that may be set.
}

func NewReqRes(req *rolltypes.Request) *ReqRes {
	return &ReqRes{
		Request:   req,
		WaitGroup: waitGroup1(),
		Response:  nil,

		callbackInvoked: false,
		cb:              nil,
	}
}

// Sets sets the callback. If reqRes is already done, it will call the cb
// immediately. Note, reqRes.cb should not change if reqRes.done and only one
// callback is supported.
func (r *ReqRes) SetCallback(cb func(res *rolltypes.Response)) {
	r.mtx.Lock()

	if r.callbackInvoked {
		r.mtx.Unlock()
		cb(r.Response)
		return
	}

	r.cb = cb
	r.mtx.Unlock()
}

// InvokeCallback invokes a thread-safe execution of the configured callback
// if non-nil.
func (r *ReqRes) InvokeCallback() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.cb != nil {
		r.cb(r.Response)
	}
	r.callbackInvoked = true
}

// GetCallback returns the configured callback of the ReqRes object which may be
// nil. Note, it is not safe to concurrently call this in cases where it is
// marked done and SetCallback is called before calling GetCallback as that
// will invoke the callback twice and create a potential race condition.
//
// ref: https://github.com/tendermint/tendermint/issues/5439
func (r *ReqRes) GetCallback() func(*rolltypes.Response) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.cb
}

func waitGroup1() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	return
}
