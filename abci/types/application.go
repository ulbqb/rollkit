package types

import (
	"github.com/tendermint/tendermint/abci/types"
)

type Application interface {
	types.Application

	// Get appHash
	GetAppHash(RequestGetAppHash) ResponseGetAppHash
	// Generate Fraud Proof
	GenerateFraudProof(RequestGenerateFraudProof) ResponseGenerateFraudProof
	// Verifies a Fraud Proof
	VerifyFraudProof(RequestVerifyFraudProof) ResponseVerifyFraudProof
}

//-------------------------------------------------------
// BaseApplication is a base form of Application

var _ Application = (*BaseApplication)(nil)

type BaseApplication struct {
}

func NewBaseApplication() *BaseApplication {
	return &BaseApplication{}
}

func (BaseApplication) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{}
}

func (BaseApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (BaseApplication) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

func (BaseApplication) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

func (BaseApplication) Commit() types.ResponseCommit {
	return types.ResponseCommit{}
}

func (BaseApplication) Query(req types.RequestQuery) types.ResponseQuery {
	return types.ResponseQuery{Code: types.CodeTypeOK}
}

func (BaseApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	return types.ResponseInitChain{}
}

func (BaseApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	return types.ResponseBeginBlock{}
}

func (BaseApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{}
}

func (BaseApplication) ListSnapshots(req types.RequestListSnapshots) types.ResponseListSnapshots {
	return types.ResponseListSnapshots{}
}

func (BaseApplication) OfferSnapshot(req types.RequestOfferSnapshot) types.ResponseOfferSnapshot {
	return types.ResponseOfferSnapshot{}
}

func (BaseApplication) LoadSnapshotChunk(req types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk {
	return types.ResponseLoadSnapshotChunk{}
}

func (BaseApplication) ApplySnapshotChunk(req types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk {
	return types.ResponseApplySnapshotChunk{}
}

func (BaseApplication) GetAppHash(req RequestGetAppHash) ResponseGetAppHash {
	return ResponseGetAppHash{}
}

func (BaseApplication) GenerateFraudProof(req RequestGenerateFraudProof) ResponseGenerateFraudProof {
	return ResponseGenerateFraudProof{}
}

func (BaseApplication) VerifyFraudProof(req RequestVerifyFraudProof) ResponseVerifyFraudProof {
	return ResponseVerifyFraudProof{}
}
