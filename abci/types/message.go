package types

import "github.com/tendermint/tendermint/abci/types"

func ToRequestEcho(message string) *Request {
	return &Request{
		Value: &Request_Echo{&types.RequestEcho{Message: message}},
	}
}

func ToRequestFlush() *Request {
	return &Request{
		Value: &Request_Flush{&types.RequestFlush{}},
	}
}

func ToRequestInfo(req types.RequestInfo) *Request {
	return &Request{
		Value: &Request_Info{&req},
	}
}

func ToRequestSetOption(req types.RequestSetOption) *Request {
	return &Request{
		Value: &Request_SetOption{&req},
	}
}

func ToRequestDeliverTx(req types.RequestDeliverTx) *Request {
	return &Request{
		Value: &Request_DeliverTx{&req},
	}
}

func ToRequestCheckTx(req types.RequestCheckTx) *Request {
	return &Request{
		Value: &Request_CheckTx{&req},
	}
}

func ToRequestCommit() *Request {
	return &Request{
		Value: &Request_Commit{&types.RequestCommit{}},
	}
}

func ToRequestQuery(req types.RequestQuery) *Request {
	return &Request{
		Value: &Request_Query{&req},
	}
}

func ToRequestInitChain(req types.RequestInitChain) *Request {
	return &Request{
		Value: &Request_InitChain{&req},
	}
}

func ToRequestBeginBlock(req types.RequestBeginBlock) *Request {
	return &Request{
		Value: &Request_BeginBlock{&req},
	}
}

func ToRequestEndBlock(req types.RequestEndBlock) *Request {
	return &Request{
		Value: &Request_EndBlock{&req},
	}
}

func ToRequestListSnapshots(req types.RequestListSnapshots) *Request {
	return &Request{
		Value: &Request_ListSnapshots{&req},
	}
}

func ToRequestOfferSnapshot(req types.RequestOfferSnapshot) *Request {
	return &Request{
		Value: &Request_OfferSnapshot{&req},
	}
}

func ToRequestLoadSnapshotChunk(req types.RequestLoadSnapshotChunk) *Request {
	return &Request{
		Value: &Request_LoadSnapshotChunk{&req},
	}
}

func ToRequestApplySnapshotChunk(req types.RequestApplySnapshotChunk) *Request {
	return &Request{
		Value: &Request_ApplySnapshotChunk{&req},
	}
}

func ToRequestGetAppHash(req RequestGetAppHash) *Request {
	return &Request{
		Value: &Request_GetAppHash{&req},
	}
}

func ToRequestGenerateFraudProof(req RequestGenerateFraudProof) *Request {
	return &Request{
		Value: &Request_GenerateFraudProof{&req},
	}
}

func ToRequestVerifyFraudProof(req RequestVerifyFraudProof) *Request {
	return &Request{
		Value: &Request_VerifyFraudProof{&req},
	}
}

// ------

func ToResponseException(errStr string) *Response {
	return &Response{
		Value: &Response_Exception{&types.ResponseException{Error: errStr}},
	}
}

func ToResponseEcho(message string) *Response {
	return &Response{
		Value: &Response_Echo{&types.ResponseEcho{Message: message}},
	}
}

func ToResponseFlush() *Response {
	return &Response{
		Value: &Response_Flush{&types.ResponseFlush{}},
	}
}

func ToResponseInfo(res types.ResponseInfo) *Response {
	return &Response{
		Value: &Response_Info{&res},
	}
}

func ToResponseSetOption(res types.ResponseSetOption) *Response {
	return &Response{
		Value: &Response_SetOption{&res},
	}
}

func ToResponseDeliverTx(res types.ResponseDeliverTx) *Response {
	return &Response{
		Value: &Response_DeliverTx{&res},
	}
}

func ToResponseCheckTx(res types.ResponseCheckTx) *Response {
	return &Response{
		Value: &Response_CheckTx{&res},
	}
}

func ToResponseCommit(res types.ResponseCommit) *Response {
	return &Response{
		Value: &Response_Commit{&res},
	}
}

func ToResponseQuery(res types.ResponseQuery) *Response {
	return &Response{
		Value: &Response_Query{&res},
	}
}

func ToResponseInitChain(res types.ResponseInitChain) *Response {
	return &Response{
		Value: &Response_InitChain{&res},
	}
}

func ToResponseBeginBlock(res types.ResponseBeginBlock) *Response {
	return &Response{
		Value: &Response_BeginBlock{&res},
	}
}

func ToResponseEndBlock(res types.ResponseEndBlock) *Response {
	return &Response{
		Value: &Response_EndBlock{&res},
	}
}

func ToResponseListSnapshots(res types.ResponseListSnapshots) *Response {
	return &Response{
		Value: &Response_ListSnapshots{&res},
	}
}

func ToResponseOfferSnapshot(res types.ResponseOfferSnapshot) *Response {
	return &Response{
		Value: &Response_OfferSnapshot{&res},
	}
}

func ToResponseLoadSnapshotChunk(res types.ResponseLoadSnapshotChunk) *Response {
	return &Response{
		Value: &Response_LoadSnapshotChunk{&res},
	}
}

func ToResponseApplySnapshotChunk(res types.ResponseApplySnapshotChunk) *Response {
	return &Response{
		Value: &Response_ApplySnapshotChunk{&res},
	}
}

func ToResponseGetAppHash(res ResponseGetAppHash) *Response {
	return &Response{
		Value: &Response_GetAppHash{&res},
	}
}

func ToResponseGenerateFraudProof(res ResponseGenerateFraudProof) *Response {
	return &Response{
		Value: &Response_GenerateFraudProof{&res},
	}
}

func ToResponseVerifyFraudProof(res ResponseVerifyFraudProof) *Response {
	return &Response{
		Value: &Response_VerifyFraudProof{&res},
	}
}
