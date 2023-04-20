package abcicli

import (
	types "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/service"
	tmsync "github.com/tendermint/tendermint/libs/sync"

	rolltypes "github.com/rollkit/rollkit/abci/types"
)

type localClient struct {
	service.BaseService

	mtx *tmsync.Mutex
	rolltypes.Application
	Callback
}

func NewLocalClient(mtx *tmsync.Mutex, app rolltypes.Application) Client {
	if mtx == nil {
		mtx = new(tmsync.Mutex)
	}
	cli := &localClient{
		mtx:         mtx,
		Application: app,
	}
	cli.BaseService = *service.NewBaseService(nil, "localClient", cli)
	return cli
}

func (app *localClient) SetResponseCallback(cb Callback) {
	app.mtx.Lock()
	app.Callback = cb
	app.mtx.Unlock()
}

// TODO: change types.Application to include Error()?
func (app *localClient) Error() error {
	return nil
}

func (app *localClient) FlushAsync() *ReqRes {
	// Do nothing
	return newLocalReqRes(rolltypes.ToRequestFlush(), nil)
}

func (app *localClient) EchoAsync(msg string) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.callback(
		rolltypes.ToRequestEcho(msg),
		rolltypes.ToResponseEcho(msg),
	)
}

func (app *localClient) InfoAsync(req types.RequestInfo) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Info(req)
	return app.callback(
		rolltypes.ToRequestInfo(req),
		rolltypes.ToResponseInfo(res),
	)
}

func (app *localClient) SetOptionAsync(req types.RequestSetOption) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.SetOption(req)
	return app.callback(
		rolltypes.ToRequestSetOption(req),
		rolltypes.ToResponseSetOption(res),
	)
}

func (app *localClient) DeliverTxAsync(params types.RequestDeliverTx) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(params)
	return app.callback(
		rolltypes.ToRequestDeliverTx(params),
		rolltypes.ToResponseDeliverTx(res),
	)
}

func (app *localClient) CheckTxAsync(req types.RequestCheckTx) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return app.callback(
		rolltypes.ToRequestCheckTx(req),
		rolltypes.ToResponseCheckTx(res),
	)
}

func (app *localClient) QueryAsync(req types.RequestQuery) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Query(req)
	return app.callback(
		rolltypes.ToRequestQuery(req),
		rolltypes.ToResponseQuery(res),
	)
}

func (app *localClient) CommitAsync() *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return app.callback(
		rolltypes.ToRequestCommit(),
		rolltypes.ToResponseCommit(res),
	)
}

func (app *localClient) InitChainAsync(req types.RequestInitChain) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return app.callback(
		rolltypes.ToRequestInitChain(req),
		rolltypes.ToResponseInitChain(res),
	)
}

func (app *localClient) BeginBlockAsync(req types.RequestBeginBlock) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return app.callback(
		rolltypes.ToRequestBeginBlock(req),
		rolltypes.ToResponseBeginBlock(res),
	)
}

func (app *localClient) EndBlockAsync(req types.RequestEndBlock) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return app.callback(
		rolltypes.ToRequestEndBlock(req),
		rolltypes.ToResponseEndBlock(res),
	)
}

func (app *localClient) ListSnapshotsAsync(req types.RequestListSnapshots) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ListSnapshots(req)
	return app.callback(
		rolltypes.ToRequestListSnapshots(req),
		rolltypes.ToResponseListSnapshots(res),
	)
}

func (app *localClient) OfferSnapshotAsync(req types.RequestOfferSnapshot) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.OfferSnapshot(req)
	return app.callback(
		rolltypes.ToRequestOfferSnapshot(req),
		rolltypes.ToResponseOfferSnapshot(res),
	)
}

func (app *localClient) LoadSnapshotChunkAsync(req types.RequestLoadSnapshotChunk) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.LoadSnapshotChunk(req)
	return app.callback(
		rolltypes.ToRequestLoadSnapshotChunk(req),
		rolltypes.ToResponseLoadSnapshotChunk(res),
	)
}

func (app *localClient) ApplySnapshotChunkAsync(req types.RequestApplySnapshotChunk) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ApplySnapshotChunk(req)
	return app.callback(
		rolltypes.ToRequestApplySnapshotChunk(req),
		rolltypes.ToResponseApplySnapshotChunk(res),
	)
}

func (app *localClient) GetAppHashAsync(req rolltypes.RequestGetAppHash) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GetAppHash(req)
	return app.callback(
		rolltypes.ToRequestGetAppHash(req),
		rolltypes.ToResponseGetAppHash(res),
	)
}

func (app *localClient) GenerateFraudProofAsync(req rolltypes.RequestGenerateFraudProof) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GenerateFraudProof(req)
	return app.callback(
		rolltypes.ToRequestGenerateFraudProof(req),
		rolltypes.ToResponseGenerateFraudProof(res),
	)
}

func (app *localClient) VerifyFraudProofAsync(req rolltypes.RequestVerifyFraudProof) *ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.VerifyFraudProof(req)
	return app.callback(
		rolltypes.ToRequestVerifyFraudProof(req),
		rolltypes.ToResponseVerifyFraudProof(res),
	)
}

//-------------------------------------------------------

func (app *localClient) FlushSync() error {
	return nil
}

func (app *localClient) EchoSync(msg string) (*types.ResponseEcho, error) {
	return &types.ResponseEcho{Message: msg}, nil
}

func (app *localClient) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Info(req)
	return &res, nil
}

func (app *localClient) SetOptionSync(req types.RequestSetOption) (*types.ResponseSetOption, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.SetOption(req)
	return &res, nil
}

func (app *localClient) DeliverTxSync(req types.RequestDeliverTx) (*types.ResponseDeliverTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(req)
	return &res, nil
}

func (app *localClient) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return &res, nil
}

func (app *localClient) QuerySync(req types.RequestQuery) (*types.ResponseQuery, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Query(req)
	return &res, nil
}

func (app *localClient) CommitSync() (*types.ResponseCommit, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return &res, nil
}

func (app *localClient) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return &res, nil
}

func (app *localClient) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return &res, nil
}

func (app *localClient) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return &res, nil
}

func (app *localClient) ListSnapshotsSync(req types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ListSnapshots(req)
	return &res, nil
}

func (app *localClient) OfferSnapshotSync(req types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.OfferSnapshot(req)
	return &res, nil
}

func (app *localClient) LoadSnapshotChunkSync(
	req types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.LoadSnapshotChunk(req)
	return &res, nil
}

func (app *localClient) ApplySnapshotChunkSync(
	req types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ApplySnapshotChunk(req)
	return &res, nil
}

func (app *localClient) GetAppHashSync(
	req rolltypes.RequestGetAppHash) (*rolltypes.ResponseGetAppHash, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GetAppHash(req)
	return &res, nil
}

func (app *localClient) GenerateFraudProofSync(
	req rolltypes.RequestGenerateFraudProof) (*rolltypes.ResponseGenerateFraudProof, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.GenerateFraudProof(req)
	return &res, nil
}

func (app *localClient) VerifyFraudProofSync(
	req rolltypes.RequestVerifyFraudProof) (*rolltypes.ResponseVerifyFraudProof, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.VerifyFraudProof(req)
	return &res, nil
}

func (app *localClient) callback(req *rolltypes.Request, res *rolltypes.Response) *ReqRes {
	app.Callback(req, res)
	rr := newLocalReqRes(req, res)
	rr.callbackInvoked = true
	return rr
}

func newLocalReqRes(req *rolltypes.Request, res *rolltypes.Response) *ReqRes {
	reqRes := NewReqRes(req)
	reqRes.Response = res
	return reqRes
}
