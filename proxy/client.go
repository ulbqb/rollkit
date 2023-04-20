package proxy

import (
	tmsync "github.com/tendermint/tendermint/libs/sync"

	rollabcicli "github.com/rollkit/rollkit/abci/client"
	rolltypes "github.com/rollkit/rollkit/abci/types"
)

type ClientCreator interface {
	// NewABCIClient returns a new ABCI client.
	NewABCIClient() (rollabcicli.Client, error)
}

type localClientCreator struct {
	mtx *tmsync.Mutex
	app rolltypes.Application
}

// NewLocalClientCreator returns a ClientCreator for the given app,
// which will be running locally.
func NewLocalClientCreator(app rolltypes.Application) ClientCreator {
	return &localClientCreator{
		mtx: new(tmsync.Mutex),
		app: app,
	}
}

func (l *localClientCreator) NewABCIClient() (rollabcicli.Client, error) {
	return rollabcicli.NewLocalClient(l.mtx, l.app), nil
}
