package btc_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"

	coreUtils "github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/vald/btc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
)

func TestDecodeEventContractCall(t *testing.T) {

	txHex := "020000000001016f6651122da3f1fee9af65734274f5dbbc58d8b21d8e12a6510d0e12d40af1030200000000fdffffff03a086010000000000225120dade785d43c753bcc8c66f21fef05643ebb4d9812aa60782c7440189255bbb4b00000000000000003d6a013504531801040100080000000000aa36a714b91e3a8ef862567026d6f376c9f3d6b814ca43371424a1db57fa3ecafcbad91d6ef068439aceeae090c53d8f000000000016001450dceca158a9c872eb405d52293d351110572c9e02483045022100cd1c70983bf4991b8a9adc1ee5f3b6b67cdcd1c524241154cf7601dd01301347022022001f080250c4cc65b643438108d0724ed3cf20dd7f7ff6ad61fa0cb76f6a170121022ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350a00000000"

	// Decode the hex string into bytes
	txRaw, err := hex.DecodeString(txHex)
	if err != nil {
		log.Errorf("Failed to decode hex string %v", err)
	}

	t.Logf("Decoding BTC transaction %+v\n", txRaw)

	// Parse the transaction
	msgTx := wire.NewMsgTx(2)
	err = msgTx.Deserialize(bytes.NewReader(txRaw))
	if err != nil {
		log.Errorf("Failed to parse transaction %v", err)
	}
	if len(msgTx.TxOut) < 3 {
		t.Errorf("btcLocking tx must have at least 3 outputs")
	}

	payload, err := btc.NewOpReturnData(msgTx.TxOut)
	if err != nil {
		t.Error("cannot parse payload op return data: ", err)
	}

	t.Logf("Payload tag: %x", payload.Tag)
	t.Logf("Payload version: %x", payload.Version)
	t.Logf("Payload destination chain id: %x", payload.DestinationChainID)
	t.Logf("Payload destination recipient addr: %x", payload.DestinationRecipientAddr)
	t.Logf("Payload destination contract addr: %x", payload.DestinationContractAddr)

	chainId, err := coreUtils.BytesToInt64BigEndian(payload.DestinationChainID[:])
	if err != nil {
		t.Error("cannot parse chain id: ", err)
	}

	mintingAmount, err := btc.GetMintingAmount(msgTx.TxOut[0])
	if err != nil {
		t.Error("cannot get minting amount: ", err)
	}

	abi_minting_payload := types.Hash(common.BytesToHash(mintingAmount[:]))

	t.Logf("Minting amount: %x", abi_minting_payload)

	t.Log("Chain ID: ", chainId)
}
