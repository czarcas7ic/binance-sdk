package tx

import (
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/crypto/sr25519"
	"github.com/tendermint/go-amino"

	"gitlab.com/thorchain/binance-sdk/types/msg"
)

// cdc global variable
var Cdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)
	RegisterAmino(cdc)
	msg.RegisterCodec(cdc)
}

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)
	cdc.RegisterConcrete(sr25519.PrivKey{},
		sr25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&ed25519.PrivKey{},
		ed25519.PrivKeyName, nil)
	cdc.RegisterConcrete(&secp256k1.PrivKey{},
		secp256k1.PrivKeyName, nil)
}

func init() {
	RegisterCodec(Cdc)
}
