package lazyledger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/lazyledger/lazyledger-app/app/params"
	apptypes "github.com/lazyledger/lazyledger-app/x/lazyledgerapp/types"

	"github.com/lazyledger/optimint/da"
	"github.com/lazyledger/optimint/log"
	"github.com/lazyledger/optimint/types"
)

type Config struct {
	// PayForMessage related params
	NamespaceID []byte
	PubKey      []byte
	BaseRateMax uint64
	TipRateMax  uint64
	From        string

	// RPC related params
	Address string
	ChainID string
	Timeout time.Duration

	// keyring related params
	AppName string
	Backend string
	RootDir string
}

type LazyLedger struct {
	config Config
	logger log.Logger

	keyring keyring.Keyring

	rpcClient *grpc.ClientConn
}

var _ da.DataAvailabilityLayerClient = &LazyLedger{}

// Init is called once to allow DA client to read configuration and initialize resources.
func (ll *LazyLedger) Init(config []byte, logger log.Logger) error {
	ll.logger = logger
	err := toml.Unmarshal(config, &ll.config)
	if err != nil {
		return err
	}
	var userInput io.Reader
	// TODO(tzdybal): this means interactive reading from stdin - shouldn't we replace this somehow?
	userInput = os.Stdin
	ll.keyring, err = keyring.New(ll.config.AppName, ll.config.Backend, ll.config.RootDir, userInput)
	return err
}

func (ll *LazyLedger) Start() (err error) {
	ll.rpcClient, err = grpc.Dial(ll.config.Address, grpc.WithInsecure())
	return
}

func (ll *LazyLedger) Stop() error {
	return ll.rpcClient.Close()
}

// SubmitBlock submits the passed in block to the DA layer.
// This should create a transaction which (potentially)
// triggers a state transition in the DA layer.
func (ll *LazyLedger) SubmitBlock(block *types.Block) da.ResultSubmitBlock {
	msg, err := ll.preparePayForMessage(block)
	if err != nil {
		return da.ResultSubmitBlock{Code: da.StatusError, Message: err.Error()}
	}

	err = ll.callRPC(msg)
	if err != nil {
		return da.ResultSubmitBlock{Code: da.StatusError, Message: err.Error()}
	}

	return da.ResultSubmitBlock{Code: da.StatusSuccess}
}

func (ll *LazyLedger) preparePayForMessage(block *types.Block) (*apptypes.MsgWirePayForMessage, error) {
	// TODO(tzdybal): serialize block
	var message []byte
	message, err := block.Serialize()
	if err != nil {
		return nil, err
	}

	// create PayForMessage message
	msg, err := apptypes.NewMsgWirePayForMessage(
		ll.config.NamespaceID,
		message,
		ll.config.PubKey,
		&apptypes.TransactionFee{
			BaseRateMax: ll.config.BaseRateMax,
			TipRateMax:  ll.config.TipRateMax,
		},
		apptypes.SquareSize,
	)
	if err != nil {
		return nil, err
	}

	// sign the PayForMessage's ShareCommitments
	err = msg.SignShareCommitments(ll.config.From, ll.keyring)
	if err != nil {
		return nil, err
	}

	// run message checks
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ll *LazyLedger) sign(msg *apptypes.MsgWirePayForMessage) (*tx.BroadcastTxRequest, error) {
	encCfg := params.MakeEncodingConfig()

	// Create a new TxBuilder.
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	txBuilder = setConfigs(txBuilder)

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, err
	}

	info, err := ll.keyring.Key(ll.config.From)
	if err != nil {
		return nil, err
	}

	sigV2 := signing.SignatureV2{
		PubKey: info.GetPubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: 0, // we need to find this by querying the node
	}

	txBuilder.SetSignatures(sigV2)

	signerData := authsigning.SignerData{
		ChainID:       ll.config.ChainID,
		AccountNumber: 0,
		Sequence:      0,
	}

	// Generate the bytes to be signed.
	bytesToSign, err := encCfg.TxConfig.SignModeHandler().GetSignBytes(
		signing.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return nil, err
	}

	// Sign those bytes
	sigBytes, _, err := ll.keyring.Sign(ll.config.From, bytesToSign)
	if err != nil {
		return nil, err
	}

	// Construct the SignatureV2 struct
	sig := signing.SignatureV2{
		PubKey: info.GetPubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: sigBytes,
		},
		Sequence: 0,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	// Generated Protobuf-encoded bytes.
	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	}, nil
}

func setConfigs(builder client.TxBuilder) client.TxBuilder {
	coin := sdk.Coin{
		Denom:  "token",
		Amount: sdk.NewInt(10),
	}
	builder.SetGasLimit(100000)
	builder.SetFeeAmount(sdk.NewCoins(coin))
	return builder
}

func (ll *LazyLedger) callRPC(msg *apptypes.MsgWirePayForMessage) error {
	txReq, err := ll.sign(msg)
	if err != nil {
		return err
	}

	txClient := tx.NewServiceClient(ll.rpcClient)

	resp, err := txClient.BroadcastTx(context.Background(), txReq)
	fmt.Println("tzdybal:", resp)
	if err != nil {
		return err
	}

	return nil
}