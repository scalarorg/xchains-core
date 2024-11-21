package btc_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/bitcoin-vault/ffi/go-vault"
	evmUtils "github.com/scalarorg/bitcoin-vault/go-utils/evm"
)

// // go test -timeout 10m -run ^TestDecodeEventContractCall$ github.com/axelarnetwork/axelar-core/vald/btc -v -count=1
// func TestDecodeEventContractCall(t *testing.T) {

// 	txHex := "020000000001011713e20bd169b9fe7afd16831989b4a893945150c40f252047cf58b7acaffcfa0200000000fdffffff03102700000000000022512067bff357780a93826a444646aec681c4ff1f4316244478c0d611f91a75c93b8a0000000000000000416a3f5343414c41526c69676874000100030000000000aa36a7b91e3a8ef862567026d6f376c9f3d6b814ca433724a1db57fa3ecafcbad91d6ef068439aceeae090c3a97f000000000016001450dceca158a9c872eb405d52293d351110572c9e02483045022100dccb6ce51e325bb96c1fa0dde646eafab89cca6e55d0b3f8c2fe021da91f5ba6022066715e44928d25f3869d8c424743d94f64867aeb883e95863bd2b1d9df85d2250121022ae31ea8709aeda8194ba3e2f7e7e95e680e8b65135c8983c0a298d17bc5350a00000000"

// 	// Decode the hex string into bytes
// 	txRaw, err := hex.DecodeString(txHex)
// 	if err != nil {
// 		log.Errorf("Failed to decode hex string %v", err)
// 	}

// 	t.Logf("Decoding BTC transaction %+v\n", txRaw)

// 	// Parse the transaction
// 	msgTx := wire.NewMsgTx(2)
// 	err = msgTx.Deserialize(bytes.NewReader(txRaw))
// 	if err != nil {
// 		log.Errorf("Failed to parse transaction %v", err)
// 	}
// 	if len(msgTx.TxOut) < 3 {
// 		t.Errorf("btcLocking tx must have at least 3 outputs")
// 	}

// 	output, err := goutils.ParseVaultEmbeddedData(msgTx.TxOut[1].PkScript)
// 	if err != nil {
// 		t.Error("cannot parse payload op return data: ", err)
// 	}

// 	txid, err := hex.DecodeString(msgTx.TxID())
// 	if err != nil {
// 		t.Error("cannot decode txid: ", err)
// 	}
// 	t.Logf("BTC TxID: %x", txid)
// 	t.Logf("BTC TxHash: %x", msgTx.TxHash())
// 	t.Logf("Payload tag: %x", output.Tag)
// 	t.Logf("Payload version: %x", output.Version)
// 	t.Logf("Payload destination chain id: %x", output.DestinationChainID)
// 	t.Logf("Payload destination recipient addr: %x", output.DestinationRecipientAddress)
// 	t.Logf("Payload destination contract addr: %x", output.DestinationContractAddress)
// 	t.Logf("Output: %+v\n", output)

// 	var senderBytes [20]byte
// 	copy(senderBytes[:], output.DestinationRecipientAddress[:])

// 	var txIdBytes [32]byte
// 	txId := msgTx.TxID()
// 	txBytes, err := hex.DecodeString(txId)
// 	if err != nil {
// 		t.Error("failed to decode tx id: ", err)
// 	}
// 	copy(txIdBytes[:], txBytes)

// 	var mintingAmount int64 = msgTx.TxOut[0].Value

// 	_, payloadHash, err := goutils.CalculateStakingPayloadHash(senderBytes, mintingAmount, txIdBytes)
// 	if err != nil {
// 		t.Error("failed to get payload hash: ", err)
// 	}

// 	t.Logf("Payload hash: %x", payloadHash)
// }

// Note: To run this test, must build bitcoin-vault-ffi first then copy to the lib folder
// cp -p ../../bitcoin-vault/target/release/libbitcoin_vault_ffi.* ./lib

// CGO_LDFLAGS="-L./lib -lbitcoin_vault_ffi" CGO_CFLAGS="-I./lib" go test -timeout 10m -run ^TestDecodeEventContractCall$ github.com/axelarnetwork/axelar-core/vald/btc -v -count=1
func TestDecodeEventContractCall(t *testing.T) {

	txHex := "02000000000101bb520c6dc5043249b56157246cf381d0ab8abcf6fae82f2c18f757a9d02e616e0200000000fdffffff030f270000000000002251207f815abf6dfd78423a708aa8db1c2c906eecac910c035132d342e4988a37b8d500000000000000003c6a3a5343414c415200014003a736aa00000000001f98c06d8734d5a9ff0b53e3294626e62e4d232c130c4810d57140e1e62967cbf742caeae91b6ece1aa2052a0100000022512095033d48b6029174ed3ba21390756c56e90c41eeeef41c172c81d1d09a167cda014097f6cab5f3e5c2a8e65120dd9f1582d2c5fce543728381cbf014ceb64bc79fb2420673c60540ad0040f44d4e9872ed47c43155e7edd4f6e9555ee87c37d2c10d00000000"

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

	output, err := vault.ParseVaultEmbeddedData(msgTx.TxOut[1].PkScript)
	if err != nil {
		t.Error("cannot parse payload op return data: ", err)
	}

	txid, err := hex.DecodeString(msgTx.TxID())
	if err != nil {
		t.Error("cannot decode txid: ", err)
	}

	t.Logf("BTC TxID: %x", txid)
	t.Logf("BTC TxHash: %x", msgTx.TxHash())

	t.Logf("Payload tag: %x", output.Tag)
	t.Logf("Payload version: %x", output.Version)
	t.Logf("Payload network id: %x", output.NetworkID)
	t.Logf("Payload flags: %x", output.Flags)
	t.Logf("Payload service tag: %x", output.ServiceTag)
	t.Logf("Payload have only covenants: %v", output.HaveOnlyCovenants)
	t.Logf("Payload covenant quorum: %x", output.CovenantQuorum)
	t.Logf("Payload destination chain id: %x", output.DestinationChain)
	t.Logf("Payload destination contract addr: %x", output.DestinationContractAddress)
	t.Logf("Payload destination recipient addr: %x", output.DestinationRecipientAddress)

	t.Logf("Output: %+v\n", output)

	var senderBytes [20]byte
	copy(senderBytes[:], output.DestinationRecipientAddress[:])

	var txIdBytes [32]byte
	txId := msgTx.TxID()
	txBytes, err := hex.DecodeString(txId)
	if err != nil {
		t.Error("failed to decode tx id: ", err)
	}
	copy(txIdBytes[:], txBytes)

	var mintingAmount int64 = msgTx.TxOut[0].Value

	_, payloadHash, err := evmUtils.CalculateStakingPayloadHash(senderBytes, mintingAmount, txIdBytes)
	if err != nil {
		t.Error("failed to get payload hash: ", err)
	}

	t.Logf("Payload hash: %x", payloadHash)
}
