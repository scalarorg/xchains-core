package btc_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/btc-vault/btcvault"

	coreUtils "github.com/axelarnetwork/axelar-core/utils"
)

func TestDecodeEventContractCall(t *testing.T) {

	txHex := "02000000000101502f447e05fcbabe0c8d3d401db31cb1a2f74c66d530ff593b4df83adc0465ca0000000000fdffffff04f82a0000000000002251205cd309728cdf3d5e8f723530c7429c61855337866ee8e180c5e1e3f6a4ac76610000000000000000476a450102030400d529c49a30b8dfd061b3458893abe74345651adf0c723d1061aa945fcb0b5ff2332343ee3d4d01132920ab977f59f767b1aefb83d208e3e1a3a349b949ccb66a00000000000000003a6a380000000000aa36a7130c4810d57140e1e62967cbf742caeae91b6ece768e8de8cf0c7747d41f75f83c914a19c5921cf30000000000002af8a3eba704000000001600145c19d6bd3f9a24eeb584c9ef0f9a8d6c2348e65402463043021f7ee5a635e7d968f722b871b2aaaa802f3015ee303cae4ea8f7fc540ef646170220269cb89a422a893aeb754c78837d52a9f5aef7de6f411e0be01d352d34cc222b012103d529c49a30b8dfd061b3458893abe74345651adf0c723d1061aa945fcb0b5ff200000000"

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
	payloadData, err := btcvault.NewPayloadOpReturnDataFromTxOutput(msgTx.TxOut[2])
	if err != nil {
		t.Error("cannot parse payload op return data: ", err)
	}

	chainId, err := coreUtils.BytesToInt64BigEndian(payloadData.ChainID)
	if err != nil {
		t.Error("cannot parse chain id: ", err)
	}

	t.Log("Chain ID: ", chainId)
}
