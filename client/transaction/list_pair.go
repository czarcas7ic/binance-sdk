package transaction

import (
	"gitlab.com/thorchain/binance-sdk/types/msg"
	"gitlab.com/thorchain/binance-sdk/types/tx"
)

type ListPairResult struct {
	tx.TxCommitResult
}

func (c *client) ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool, options ...Option) (*ListPairResult, error) {
	fromAddr := c.keyManager.GetAddr()

	burnMsg := msg.NewDexListMsg(fromAddr, proposalId, baseAssetSymbol, quoteAssetSymbol, initPrice)
	commit, err := c.broadcastMsg(burnMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &ListPairResult{*commit}, nil

}
