package mock

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	libbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/rpc/client"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"gitlab.com/thorchain/binance-sdk/client/rpc"
)

// ABCIApp will send all abci related request to the named app,
// so you can test app behavior from a client without needing
// an entire tendermint node
type ABCIApp struct {
	App abci.Application
}

var (
	_ rpc.ABCIClient = ABCIApp{}
	_ rpc.ABCIClient = ABCIMock{}
	_ rpc.ABCIClient = (*ABCIRecorder)(nil)
)

func (a ABCIApp) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	resp, err := a.App.Info(ctx, proxy.RequestInfo)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultABCIInfo{Response: *resp}, nil
}

func (a ABCIApp) ABCIQuery(ctx context.Context, path string, data libbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return a.ABCIQueryWithOptions(ctx, path, data, client.DefaultABCIQueryOptions)
}

func (a ABCIApp) ABCIQueryWithOptions(ctx context.Context, path string, data libbytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	q, err := a.App.Query(ctx, &abci.RequestQuery{
		Data:   data,
		Path:   path,
		Height: opts.Height,
		Prove:  opts.Prove,
	})
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultABCIQuery{Response: *q}, nil
}

// NOTE: Caller should call a.App.Commit() separately,
// this function does not actually wait for a commit.
// TODO: Make it wait for a commit and set res.Height appropriately.
func (a ABCIApp) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return &ctypes.ResultBroadcastTxCommit{}, nil
}

func (a ABCIApp) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	c, err := a.App.CheckTx(ctx, &abci.RequestCheckTx{Tx: tx})
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultBroadcastTx{Code: c.Code, Data: c.Data, Log: c.Log, Hash: tx.Hash()}, nil
}

func (a ABCIApp) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	c, err := a.App.CheckTx(ctx, &abci.RequestCheckTx{Tx: tx})
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultBroadcastTx{Code: c.Code, Data: c.Data, Log: c.Log, Hash: tx.Hash()}, nil
}

// ABCIMock will send all abci related request to the named app,
// so you can test app behavior from a client without needing
// an entire tendermint node
type ABCIMock struct {
	Info            Call
	Query           Call
	BroadcastCommit Call
	Broadcast       Call
}

func (m ABCIMock) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	res, err := m.Info.GetResponse(nil)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultABCIInfo{Response: res.(abci.ResponseInfo)}, nil
}

func (m ABCIMock) ABCIQuery(ctx context.Context, path string, data libbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return m.ABCIQueryWithOptions(ctx, path, data, client.DefaultABCIQueryOptions)
}

func (m ABCIMock) ABCIQueryWithOptions(ctx context.Context, path string, data libbytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	res, err := m.Query.GetResponse(QueryArgs{path, data, opts.Height, opts.Prove})
	if err != nil {
		return nil, err
	}
	resQuery := res.(abci.ResponseQuery)
	return &ctypes.ResultABCIQuery{Response: resQuery}, nil
}

func (m ABCIMock) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return &ctypes.ResultBroadcastTxCommit{}, nil
}

func (m ABCIMock) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := m.Broadcast.GetResponse(tx)
	if err != nil {
		return nil, err
	}
	return res.(*ctypes.ResultBroadcastTx), nil
}

func (m ABCIMock) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := m.Broadcast.GetResponse(tx)
	if err != nil {
		return nil, err
	}
	return res.(*ctypes.ResultBroadcastTx), nil
}

// ABCIRecorder can wrap another type (ABCIApp, ABCIMock, or Client)
// and record all ABCI related calls.
type ABCIRecorder struct {
	Client rpc.ABCIClient
	Calls  []Call
}

func NewABCIRecorder(client rpc.ABCIClient) *ABCIRecorder {
	return &ABCIRecorder{
		Client: client,
		Calls:  []Call{},
	}
}

type QueryArgs struct {
	Path   string
	Data   libbytes.HexBytes
	Height int64
	Prove  bool
}

func (r *ABCIRecorder) addCall(call Call) {
	r.Calls = append(r.Calls, call)
}

func (r *ABCIRecorder) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	res, err := r.Client.ABCIInfo(ctx)
	r.addCall(Call{
		Name:     "abci_info",
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *ABCIRecorder) ABCIQuery(ctx context.Context, path string, data libbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return r.ABCIQueryWithOptions(ctx, path, data, client.DefaultABCIQueryOptions)
}

func (r *ABCIRecorder) ABCIQueryWithOptions(ctx context.Context, path string, data libbytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	res, err := r.Client.ABCIQueryWithOptions(ctx, path, data, opts)
	r.addCall(Call{
		Name:     "abci_query",
		Args:     QueryArgs{path, data, opts.Height, opts.Prove},
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *ABCIRecorder) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return &ctypes.ResultBroadcastTxCommit{}, nil
}

func (r *ABCIRecorder) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := r.Client.BroadcastTxAsync(ctx, tx)
	r.addCall(Call{
		Name:     "broadcast_tx_async",
		Args:     tx,
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *ABCIRecorder) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := r.Client.BroadcastTxSync(ctx, tx)
	r.addCall(Call{
		Name:     "broadcast_tx_sync",
		Args:     tx,
		Response: res,
		Error:    err,
	})
	return res, err
}
