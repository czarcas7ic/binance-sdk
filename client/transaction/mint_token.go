package transaction

import (
	"fmt"

	"gitlab.com/thorchain/binance-sdk/types/msg"
	"gitlab.com/thorchain/binance-sdk/types/tx"
)

type MintTokenResult struct {
	tx.TxCommitResult
}

func (c *client) MintToken(symbol string, amount int64, sync bool, options ...Option) (*MintTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can'c be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	mintMsg := msg.NewMintMsg(
		fromAddr,
		symbol,
		amount,
	)
	commit, err := c.broadcastMsg(mintMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &MintTokenResult{*commit}, nil

}
