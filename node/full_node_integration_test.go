package node

import (
	"context"
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	rollabci "github.com/rollkit/rollkit/abci/types"
	"github.com/rollkit/rollkit/config"
	"github.com/rollkit/rollkit/da"
	mockda "github.com/rollkit/rollkit/da/mock"
	"github.com/rollkit/rollkit/mocks"
	"github.com/rollkit/rollkit/p2p"
	rollproxy "github.com/rollkit/rollkit/proxy"
	"github.com/rollkit/rollkit/store"
	"github.com/rollkit/rollkit/types"
)

func TestAggregatorMode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	app := &mocks.Application{}
	app.On("InitChain", mock.Anything).Return(abci.ResponseInitChain{})
	app.On("CheckTx", mock.Anything).Return(abci.ResponseCheckTx{})
	app.On("BeginBlock", mock.Anything).Return(abci.ResponseBeginBlock{})
	app.On("DeliverTx", mock.Anything).Return(abci.ResponseDeliverTx{})
	app.On("EndBlock", mock.Anything).Return(abci.ResponseEndBlock{})
	app.On("Commit", mock.Anything).Return(abci.ResponseCommit{})
	app.On("GetAppHash", mock.Anything).Return(rollabci.ResponseGetAppHash{})

	key, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	anotherKey, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	genesisValidators, signingKey := getGenesisValidatorSetWithSigner(1)
	blockManagerConfig := config.BlockManagerConfig{
		BlockTime:   1 * time.Second,
		NamespaceID: types.NamespaceID{1, 2, 3, 4, 5, 6, 7, 8},
	}
	node, err := newFullNode(context.Background(), config.NodeConfig{DALayer: "mock", Aggregator: true, BlockManagerConfig: blockManagerConfig}, key, signingKey, rollproxy.NewLocalClientCreator(app), &tmtypes.GenesisDoc{ChainID: "test", Validators: genesisValidators}, log.TestingLogger())
	require.NoError(err)
	require.NotNil(node)

	assert.False(node.IsRunning())

	err = node.Start()
	assert.NoError(err)
	defer func() {
		err := node.Stop()
		assert.NoError(err)
	}()
	assert.True(node.IsRunning())

	pid, err := peer.IDFromPrivateKey(anotherKey)
	require.NoError(err)
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				node.incomingTxCh <- &p2p.GossipMessage{Data: []byte(time.Now().String()), From: pid}
				time.Sleep(time.Duration(mrand.Uint32()%20) * time.Millisecond) //nolint:gosec
			}
		}
	}()
	time.Sleep(3 * time.Second)
	cancel()
}

// TestTxGossipingAndAggregation setups a network of nodes, with single aggregator and multiple producers.
// Nodes should gossip transactions and aggregator node should produce blocks.
func TestTxGossipingAndAggregation(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	clientNodes := 4
	nodes, apps := createAndStartNodes(clientNodes, false, t)
	aggApp := apps[0]
	apps = apps[1:]

	aggApp.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
	aggApp.AssertExpectations(t)

	for i, app := range apps {
		app.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
		app.AssertExpectations(t)

		// assert that we have most of the blocks from aggregator
		beginCnt := 0
		endCnt := 0
		commitCnt := 0
		for _, call := range app.Calls {
			switch call.Method {
			case "BeginBlock":
				beginCnt++
			case "EndBlock":
				endCnt++
			case "Commit":
				commitCnt++
			}
		}
		aggregatorHeight := nodes[0].Store.Height()
		adjustedHeight := int(aggregatorHeight - 3) // 3 is completely arbitrary
		assert.GreaterOrEqual(beginCnt, adjustedHeight)
		assert.GreaterOrEqual(endCnt, adjustedHeight)
		assert.GreaterOrEqual(commitCnt, adjustedHeight)

		// assert that all blocks known to node are same as produced by aggregator
		for h := uint64(1); h <= nodes[i].Store.Height(); h++ {
			aggBlock, err := nodes[0].Store.LoadBlock(h)
			require.NoError(err)
			nodeBlock, err := nodes[i].Store.LoadBlock(h)
			require.NoError(err)
			assert.Equal(aggBlock, nodeBlock, fmt.Sprintf("height: %d", h))
		}
	}
}

func TestLazyAggregator(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	app := &mocks.Application{}
	app.On("InitChain", mock.Anything).Return(abci.ResponseInitChain{})
	app.On("CheckTx", mock.Anything).Return(abci.ResponseCheckTx{})
	app.On("BeginBlock", mock.Anything).Return(abci.ResponseBeginBlock{})
	app.On("DeliverTx", mock.Anything).Return(abci.ResponseDeliverTx{})
	app.On("EndBlock", mock.Anything).Return(abci.ResponseEndBlock{})
	app.On("Commit", mock.Anything).Return(abci.ResponseCommit{})
	app.On("GetAppHash", mock.Anything).Return(rollabci.ResponseGetAppHash{})

	key, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	genesisValidators, signingKey := getGenesisValidatorSetWithSigner(1)
	blockManagerConfig := config.BlockManagerConfig{
		BlockTime:   1 * time.Second,
		NamespaceID: types.NamespaceID{1, 2, 3, 4, 5, 6, 7, 8},
	}

	node, err := NewNode(context.Background(), config.NodeConfig{
		DALayer:            "mock",
		Aggregator:         true,
		BlockManagerConfig: blockManagerConfig,
		LazyAggregator:     true,
	}, key, signingKey, rollproxy.NewLocalClientCreator(app), &tmtypes.GenesisDoc{ChainID: "test", Validators: genesisValidators}, log.TestingLogger())
	assert.False(node.IsRunning())
	assert.NoError(err)
	err = node.Start()
	assert.NoError(err)
	defer func() {
		err := node.Stop()
		assert.NoError(err)
	}()
	assert.True(node.IsRunning())

	require.NoError(err)

	// delay to ensure it waits for block 1 to be built
	time.Sleep(1 * time.Second)

	client := node.GetClient()

	_, err = client.BroadcastTxCommit(context.Background(), []byte{0, 0, 0, 1})
	assert.NoError(err)
	time.Sleep(2 * time.Second)
	assert.Equal(node.(*FullNode).Store.Height(), uint64(2))

	_, err = client.BroadcastTxCommit(context.Background(), []byte{0, 0, 0, 2})
	assert.NoError(err)
	time.Sleep(2 * time.Second)
	assert.Equal(node.(*FullNode).Store.Height(), uint64(3))

	_, err = client.BroadcastTxCommit(context.Background(), []byte{0, 0, 0, 3})
	assert.NoError(err)
	time.Sleep(2 * time.Second)
	assert.Equal(node.(*FullNode).Store.Height(), uint64(4))

}

func TestHeaderExchange(t *testing.T) {
	testSingleAggreatorSingleFullNode(t)
	testSingleAggreatorTwoFullNode(t)
	testSingleAggreatorSingleFullNodeTrustedHash(t)
	testSingleAggreatorSingleFullNodeSingleLightNode(t)
}

func testSingleAggreatorSingleFullNode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var wg sync.WaitGroup
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	clientNodes := 1
	nodes, _ := createNodes(aggCtx, ctx, clientNodes+1, false, &wg, t)

	node1 := nodes[0]
	node2 := nodes[1]

	require.NoError(node1.Start())
	time.Sleep(2 * time.Second) // wait for more than 1 blocktime for syncer to work
	require.NoError(node2.Start())

	time.Sleep(3 * time.Second)

	n1h := node1.hExService.headerStore.Height()
	aggCancel()
	require.NoError(node1.Stop())

	time.Sleep(3 * time.Second)

	n2h := node2.hExService.headerStore.Height()
	cancel()
	require.NoError(node2.Stop())

	assert.Equal(n1h, n2h, "heights must match")
}

func testSingleAggreatorTwoFullNode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var wg sync.WaitGroup
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	clientNodes := 2
	nodes, _ := createNodes(aggCtx, ctx, clientNodes+1, false, &wg, t)

	node1 := nodes[0]
	node2 := nodes[1]
	node3 := nodes[2]

	require.NoError(node1.Start())
	time.Sleep(2 * time.Second) // wait for more than 1 blocktime for syncer to work
	require.NoError(node2.Start())
	require.NoError(node3.Start())

	time.Sleep(3 * time.Second)

	n1h := node1.hExService.headerStore.Height()
	aggCancel()
	require.NoError(node1.Stop())

	time.Sleep(3 * time.Second)

	n2h := node2.hExService.headerStore.Height()
	cancel()
	require.NoError(node2.Stop())

	n3h := node3.hExService.headerStore.Height()
	require.NoError(node3.Stop())

	assert.Equal(n1h, n2h, "heights must match")
	assert.Equal(n1h, n3h, "heights must match")
}

func testSingleAggreatorSingleFullNodeTrustedHash(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var wg sync.WaitGroup
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	clientNodes := 1
	nodes, _ := createNodes(aggCtx, ctx, clientNodes+1, false, &wg, t)

	node1 := nodes[0]
	node2 := nodes[1]

	require.NoError(node1.Start())
	time.Sleep(2 * time.Second) // wait for more than 1 blocktime for syncer to work
	// Get the trusted hash from node1 and pass it to node2 config
	trustedHash, err := node1.hExService.headerStore.GetByHeight(aggCtx, 1)
	require.NoError(err)
	node2.conf.TrustedHash = trustedHash.Hash().String()
	require.NoError(node2.Start())

	time.Sleep(3 * time.Second)

	n1h := node1.hExService.headerStore.Height()
	aggCancel()
	require.NoError(node1.Stop())

	time.Sleep(3 * time.Second)

	n2h := node2.hExService.headerStore.Height()
	cancel()
	require.NoError(node2.Stop())

	assert.Equal(n1h, n2h, "heights must match")
}

func testSingleAggreatorSingleFullNodeSingleLightNode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var wg sync.WaitGroup
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	clientNodes := 2
	nodes, _ := createNodes(aggCtx, ctx, clientNodes+1, false, &wg, t)

	node1 := nodes[0]
	node2 := nodes[1]
	node3 := nodes[2]

	require.NoError(node1.Start())
	time.Sleep(2 * time.Second) // wait for more than 1 blocktime for syncer to work
	require.NoError(node2.Start())

	node3.conf.Light = true
	require.NoError(node3.Start())

	time.Sleep(3 * time.Second)

	n1h := node1.hExService.headerStore.Height()
	aggCancel()
	require.NoError(node1.Stop())

	time.Sleep(3 * time.Second)

	n2h := node2.hExService.headerStore.Height()
	n3h := node3.hExService.headerStore.Height()
	cancel()
	require.NoError(node2.Stop())
	require.NoError(node3.Stop())

	assert.Equal(n1h, n2h, "heights must match")
	assert.Equal(n1h, n3h, "heights must match")
}

// TODO: rewrite this integration test to accommodate gossip/halting mechanism of full nodes after fraud proof generation (#693)
// TestFraudProofTrigger setups a network of nodes, with single malicious aggregator and multiple producers.
// Aggregator node should produce malicious blocks, nodes should detect fraud, and generate fraud proofs
/* func TestFraudProofTrigger(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	clientNodes := 4
	nodes, apps := createAndStartNodes(clientNodes, true, t)
	aggApp := apps[0]
	apps = apps[1:]

	aggApp.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
	aggApp.AssertExpectations(t)

	for i, app := range apps {
		//app.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
		app.AssertExpectations(t)

		// assert that we have most of the blocks from aggregator
		beginCnt := 0
		endCnt := 0
		commitCnt := 0
		generateFraudProofCnt := 0
		for _, call := range app.Calls {
			switch call.Method {
			case "BeginBlock":
				beginCnt++
			case "EndBlock":
				endCnt++
			case "Commit":
				commitCnt++
			case "GenerateFraudProof":
				generateFraudProofCnt++
			}
		}
		aggregatorHeight := nodes[0].Store.Height()
		adjustedHeight := int(aggregatorHeight - 3) // 3 is completely arbitrary
		assert.GreaterOrEqual(beginCnt, adjustedHeight)
		assert.GreaterOrEqual(endCnt, adjustedHeight)
		assert.GreaterOrEqual(commitCnt, adjustedHeight)

		// GenerateFraudProof should have been called on each call to
		// BeginBlock, DeliverTx, and EndBlock so the sum of their counts
		// should be equal.
		// Note: The value of clientNodes represents number of calls to DeliverTx
		assert.Equal(beginCnt+clientNodes+endCnt, generateFraudProofCnt)

		// assert that all blocks known to node are same as produced by aggregator
		for h := uint64(1); h <= nodes[i].Store.Height(); h++ {
			nodeBlock, err := nodes[i].Store.LoadBlock(h)
			require.NoError(err)
			aggBlock, err := nodes[0].Store.LoadBlock(h)
			require.NoError(err)
			assert.Equal(aggBlock, nodeBlock)
		}
	}
} */

// Creates a starts the given number of client nodes along with an aggregator node. Uses the given flag to decide whether to have the aggregator produce malicious blocks.
func createAndStartNodes(clientNodes int, isMalicious bool, t *testing.T) ([]*FullNode, []*mocks.Application) {
	var wg sync.WaitGroup
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	nodes, apps := createNodes(aggCtx, ctx, clientNodes+1, isMalicious, &wg, t)
	startNodes(nodes, &wg, t)
	aggCancel()
	time.Sleep(100 * time.Millisecond)
	for _, n := range nodes {
		require.NoError(t, n.Stop())
	}
	cancel()
	time.Sleep(100 * time.Millisecond)
	return nodes, apps
}

// Starts the given nodes using the given wait group to synchronize them
// and wait for them to gossip transactions
func startNodes(nodes []*FullNode, wg *sync.WaitGroup, t *testing.T) {
	numNodes := len(nodes)
	wg.Add((numNodes) * (numNodes - 1))

	// Wait for aggregator node to publish the first block for full nodes to initialize header exchange service
	require.NoError(t, nodes[0].Start())
	time.Sleep(1 * time.Second)
	for i := 1; i < len(nodes); i++ {
		require.NoError(t, nodes[i].Start())
	}

	// wait for nodes to start up and establish connections; 1 second ensures that test pass even on CI.
	time.Sleep(1 * time.Second)

	for i := 1; i < len(nodes); i++ {
		data := strconv.Itoa(i) + time.Now().String()
		require.NoError(t, nodes[i].P2P.GossipTx(context.TODO(), []byte(data)))
	}

	timeout := time.NewTimer(time.Second * 30)
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		wg.Wait()
	}()
	select {
	case <-doneChan:
	case <-timeout.C:
		t.FailNow()
	}
}

// Creates the given number of nodes the given nodes using the given wait group to synchornize them
func createNodes(aggCtx, ctx context.Context, num int, isMalicious bool, wg *sync.WaitGroup, t *testing.T) ([]*FullNode, []*mocks.Application) {
	t.Helper()

	if aggCtx == nil {
		aggCtx = context.Background()
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// create keys first, as they are required for P2P connections
	keys := make([]crypto.PrivKey, num)
	for i := 0; i < num; i++ {
		keys[i], _, _ = crypto.GenerateEd25519Key(rand.Reader)
	}

	nodes := make([]*FullNode, num)
	apps := make([]*mocks.Application, num)
	dalc := &mockda.DataAvailabilityLayerClient{}
	ds, _ := store.NewDefaultInMemoryKVStore()
	_ = dalc.Init([8]byte{}, nil, ds, log.TestingLogger())
	_ = dalc.Start()
	nodes[0], apps[0] = createNode(aggCtx, 0, isMalicious, true, dalc, keys, wg, t)
	for i := 1; i < num; i++ {
		nodes[i], apps[i] = createNode(ctx, i, isMalicious, false, dalc, keys, wg, t)
	}

	return nodes, apps
}

func createNode(ctx context.Context, n int, isMalicious bool, aggregator bool, dalc da.DataAvailabilityLayerClient, keys []crypto.PrivKey, wg *sync.WaitGroup, t *testing.T) (*FullNode, *mocks.Application) {
	t.Helper()
	require := require.New(t)
	// nodes will listen on consecutive ports on local interface
	// random connections to other nodes will be added
	startPort := 10000
	p2pConfig := config.P2PConfig{
		ListenAddress: "/ip4/127.0.0.1/tcp/" + strconv.Itoa(startPort+n),
	}
	bmConfig := config.BlockManagerConfig{
		DABlockTime: 100 * time.Millisecond,
		BlockTime:   1 * time.Second, // blocks must be at least 1 sec apart for adjacent headers to get verified correctly
		NamespaceID: types.NamespaceID{8, 7, 6, 5, 4, 3, 2, 1},
		FraudProofs: true,
	}
	for i := 0; i < len(keys); i++ {
		if i == n {
			continue
		}
		r := i
		id, err := peer.IDFromPrivateKey(keys[r])
		require.NoError(err)
		p2pConfig.Seeds += "/ip4/127.0.0.1/tcp/" + strconv.Itoa(startPort+r) + "/p2p/" + id.Pretty() + ","
	}
	p2pConfig.Seeds = strings.TrimSuffix(p2pConfig.Seeds, ",")

	app := &mocks.Application{}
	app.On("InitChain", mock.Anything).Return(abci.ResponseInitChain{})
	app.On("CheckTx", mock.Anything).Return(abci.ResponseCheckTx{})
	app.On("BeginBlock", mock.Anything).Return(abci.ResponseBeginBlock{})
	app.On("EndBlock", mock.Anything).Return(abci.ResponseEndBlock{})
	app.On("Commit", mock.Anything).Return(abci.ResponseCommit{})
	maliciousAppHash := []byte{9, 8, 7, 6}
	nonMaliciousAppHash := []byte{1, 2, 3, 4}
	if isMalicious && aggregator {
		app.On("GetAppHash", mock.Anything).Return(rollabci.ResponseGetAppHash{AppHash: maliciousAppHash})
	} else {
		app.On("GetAppHash", mock.Anything).Return(rollabci.ResponseGetAppHash{AppHash: nonMaliciousAppHash})
	}

	if isMalicious && !aggregator {
		app.On("GenerateFraudProof", mock.Anything).Return(rollabci.ResponseGenerateFraudProof{FraudProof: &rollabci.FraudProof{}})
	}
	app.On("DeliverTx", mock.Anything).Return(abci.ResponseDeliverTx{}).Run(func(args mock.Arguments) {
		wg.Done()
	})
	if ctx == nil {
		ctx = context.Background()
	}

	genesisValidators, signingKey := getGenesisValidatorSetWithSigner(1)
	node, err := newFullNode(
		ctx,
		config.NodeConfig{
			P2P:                p2pConfig,
			DALayer:            "mock",
			Aggregator:         aggregator,
			BlockManagerConfig: bmConfig,
		},
		keys[n],
		signingKey,
		rollproxy.NewLocalClientCreator(app),
		&tmtypes.GenesisDoc{ChainID: "test", Validators: genesisValidators},
		log.TestingLogger().With("node", n))
	require.NoError(err)
	require.NotNil(node)

	// use same, common DALC, so nodes can share data
	node.dalc = dalc
	node.blockManager.SetDALC(dalc)

	return node, app
}
