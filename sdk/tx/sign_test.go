package tx

import (
	"reflect"
	"testing"

	"./txmsg"
)

func TestSign(t *testing.T) {

	priv, acc := PrivAndAddr()
	newOrderMsg := txmsg.NewNewOrderMsg(
		acc,
		txmsg.GenerateOrderID(1, acc),
		txmsg.OrderSide.BUY,
		"BNB_NNB",
		100000000,
		500000000,
	)

	signMsg := StdSignMsg{
		ChainID:       "bnc-chain-1",
		AccountNumber: 100,
		Sequence:      1,
		Memo:          "",
		Fee:           NewStdFee(5000, Coin{Denom: "BNB", Amount: 100000000}),
		Msgs:          []txmsg.Msg{newOrderMsg},
	}

	tx := &Tx{}
	hexTx, err := tx.Sign(priv.Bytes(), signMsg)

	// fmt.Println("stdTx: ", string(stdTx))

	if err != nil {
		t.Errorf("tx.Sign() failed, expected signed tx but got error: %v", err)
	}

	if len(hexTx) == 0 {
		t.Errorf("tx.Sign() failed, expected signed tx but got empty data: %v", hexTx)
	}

	bTx := DecodeHex(hexTx)

	var stdTx StdTx
	Cdc.UnmarshalBinaryLengthPrefixed(bTx, &stdTx)

	if !reflect.DeepEqual(stdTx.GetMsgs(), signMsg.Msgs) {
		t.Errorf("tx.Sign() decode failed, expected decoded msgs: %v to equal encoded msgss: %v", stdTx.GetMsgs(), signMsg.Msgs)
	}
}