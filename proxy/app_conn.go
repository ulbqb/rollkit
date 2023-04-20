package proxy

import (
	"github.com/tendermint/tendermint/abci/types"

	rollabcicli "github.com/rollkit/rollkit/abci/client"
	rolltypes "github.com/rollkit/rollkit/abci/types"
)

type AppConnConsensus interface {
	SetResponseCallback(rollabcicli.Callback)
	Error() error

	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)

	BeginBlockSync(types.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	DeliverTxAsync(types.RequestDeliverTx) *rollabcicli.ReqRes
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
	CommitSync() (*types.ResponseCommit, error)

	GetAppHashSync(rolltypes.RequestGetAppHash) (*rolltypes.ResponseGetAppHash, error)
	GenerateFraudProofSync(rolltypes.RequestGenerateFraudProof) (*rolltypes.ResponseGenerateFraudProof, error)
	VerifyFraudProofSync(rolltypes.RequestVerifyFraudProof) (*rolltypes.ResponseVerifyFraudProof, error)
}

type AppConnMempool interface {
	SetResponseCallback(rollabcicli.Callback)
	Error() error

	CheckTxAsync(types.RequestCheckTx) *rollabcicli.ReqRes
	CheckTxSync(types.RequestCheckTx) (*types.ResponseCheckTx, error)

	FlushAsync() *rollabcicli.ReqRes
	FlushSync() error
}

type AppConnQuery interface {
	Error() error

	EchoSync(string) (*types.ResponseEcho, error)
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
	QuerySync(types.RequestQuery) (*types.ResponseQuery, error)

	//	SetOptionSync(key string, value string) (res types.Result)
}

type AppConnSnapshot interface {
	Error() error

	ListSnapshotsSync(types.RequestListSnapshots) (*types.ResponseListSnapshots, error)
	OfferSnapshotSync(types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error)
	LoadSnapshotChunkSync(types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error)
	ApplySnapshotChunkSync(types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error)
}

type appConnConsensus struct {
	appConn rollabcicli.Client
}

func (app *appConnConsensus) SetResponseCallback(cb rollabcicli.Callback) {
	app.appConn.SetResponseCallback(cb)
}

func (app *appConnConsensus) Error() error {
	return app.appConn.Error()
}

func (app *appConnConsensus) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	return app.appConn.InitChainSync(req)
}

func (app *appConnConsensus) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	return app.appConn.BeginBlockSync(req)
}

func (app *appConnConsensus) DeliverTxAsync(req types.RequestDeliverTx) *rollabcicli.ReqRes {
	return app.appConn.DeliverTxAsync(req)
}

func (app *appConnConsensus) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	return app.appConn.EndBlockSync(req)
}

func (app *appConnConsensus) CommitSync() (*types.ResponseCommit, error) {
	return app.appConn.CommitSync()
}

func (app *appConnConsensus) GetAppHashSync(req rolltypes.RequestGetAppHash) (*rolltypes.ResponseGetAppHash, error) {
	return app.appConn.GetAppHashSync(req)
}

func (app *appConnConsensus) GenerateFraudProofSync(
	req rolltypes.RequestGenerateFraudProof,
) (*rolltypes.ResponseGenerateFraudProof, error) {
	return app.appConn.GenerateFraudProofSync(req)
}

func (app *appConnConsensus) VerifyFraudProofSync(
	req rolltypes.RequestVerifyFraudProof,
) (*rolltypes.ResponseVerifyFraudProof, error) {
	return app.appConn.VerifyFraudProofSync(req)
}

func NewAppConnConsensus(appConn rollabcicli.Client) AppConnConsensus {
	return &appConnConsensus{
		appConn: appConn,
	}
}

//------------------------------------------------
// Implements AppConnMempool (subset of abcicli.Client)

type appConnMempool struct {
	appConn rollabcicli.Client
}

func NewAppConnMempool(appConn rollabcicli.Client) AppConnMempool {
	return &appConnMempool{
		appConn: appConn,
	}
}

func (app *appConnMempool) SetResponseCallback(cb rollabcicli.Callback) {
	app.appConn.SetResponseCallback(cb)
}

func (app *appConnMempool) Error() error {
	return app.appConn.Error()
}

func (app *appConnMempool) FlushAsync() *rollabcicli.ReqRes {
	return app.appConn.FlushAsync()
}

func (app *appConnMempool) FlushSync() error {
	return app.appConn.FlushSync()
}

func (app *appConnMempool) CheckTxAsync(req types.RequestCheckTx) *rollabcicli.ReqRes {
	return app.appConn.CheckTxAsync(req)
}

func (app *appConnMempool) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	return app.appConn.CheckTxSync(req)
}

//------------------------------------------------
// Implements AppConnQuery (subset of abcicli.Client)

type appConnQuery struct {
	appConn rollabcicli.Client
}

func NewAppConnQuery(appConn rollabcicli.Client) AppConnQuery {
	return &appConnQuery{
		appConn: appConn,
	}
}

func (app *appConnQuery) Error() error {
	return app.appConn.Error()
}

func (app *appConnQuery) EchoSync(msg string) (*types.ResponseEcho, error) {
	return app.appConn.EchoSync(msg)
}

func (app *appConnQuery) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	return app.appConn.InfoSync(req)
}

func (app *appConnQuery) QuerySync(reqQuery types.RequestQuery) (*types.ResponseQuery, error) {
	return app.appConn.QuerySync(reqQuery)
}

//------------------------------------------------
// Implements AppConnSnapshot (subset of abcicli.Client)

type appConnSnapshot struct {
	appConn rollabcicli.Client
}

func NewAppConnSnapshot(appConn rollabcicli.Client) AppConnSnapshot {
	return &appConnSnapshot{
		appConn: appConn,
	}
}

func (app *appConnSnapshot) Error() error {
	return app.appConn.Error()
}

func (app *appConnSnapshot) ListSnapshotsSync(req types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	return app.appConn.ListSnapshotsSync(req)
}

func (app *appConnSnapshot) OfferSnapshotSync(req types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	return app.appConn.OfferSnapshotSync(req)
}

func (app *appConnSnapshot) LoadSnapshotChunkSync(
	req types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	return app.appConn.LoadSnapshotChunkSync(req)
}

func (app *appConnSnapshot) ApplySnapshotChunkSync(
	req types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	return app.appConn.ApplySnapshotChunkSync(req)
}
