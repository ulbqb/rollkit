package node

import (
	"context"

	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/rollkit/rollkit/config"
	rollproxy "github.com/rollkit/rollkit/proxy"
)

type Node interface {
	Start() error
	GetClient() rpcclient.Client
	Stop() error
	IsRunning() bool
}

// Add Defaults?

func NewNode(
	ctx context.Context,
	conf config.NodeConfig,
	p2pKey crypto.PrivKey,
	signingKey crypto.PrivKey,
	appClient rollproxy.ClientCreator,
	genesis *tmtypes.GenesisDoc,
	logger log.Logger,
) (Node, error) {
	if !conf.Light {
		return newFullNode(
			ctx,
			conf,
			p2pKey,
			signingKey,
			appClient,
			genesis,
			logger,
		)
	} else {
		return newLightNode(
			ctx,
			conf,
			p2pKey,
			appClient,
			genesis,
			logger,
		)
	}
}
