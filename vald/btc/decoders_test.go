package btc_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"

	coreUtils "github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/vald/btc"
)

func TestDecodeEventContractCall(t *testing.T) {

	txHex := "02000000000101b32997c3254860bd022b6a6ec90005af7977cc7a67be6bba1fcc2378093a71210200000000fdffffff03a086010000000000225120dade785d43c753bcc8c66f21fef05643ebb4d9812aa60782c7440189255bbb4b00000000000000003d6a013504531801040100080000000000aa36a714b91e3a8ef862567026d6f376c9f3d6b814ca43371424a1db57fa3ecafcbad91d6ef068439aceeae09044ee8e000000000016001450dceca158a9c872eb405d52293d351110572c9e0247304402206bedb73e1bb7c8fb3d536a132051dc65edf641fe61bd7395a20d25d38cd1e59902207c49c52f0c704546868355e526fe99325fa03deda139f97ebc987931e61ab6d70121022ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350a00000000"

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

	t.Logf("Chain ID: %d", chainId)

	payloadHash, err := btc.GetPayloadHash(payload.DestinationRecipientAddr, msgTx.TxOut[0].Value, 1730982570)
	if err != nil {
		t.Error("cannot get payload hash: ", err)
	}

	t.Logf("Payload hash: %x", payloadHash)
	payloadHex := hex.EncodeToString(payloadHash[:])
	t.Logf("Payload hex: %s", payloadHex)

	if payloadHex != "2ec538c2c7579a8ed73b43d6424b19df7beecd8ae883a6776df439f4b1281bc3" {
		t.Errorf("payload hash does not match")
	}

}
