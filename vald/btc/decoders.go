package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	vault "github.com/scalarorg/bitcoin-vault/ffi/go-vault"
	chainUtils "github.com/scalarorg/bitcoin-vault/go-utils/chain"
	evmUtils "github.com/scalarorg/bitcoin-vault/go-utils/evm"
)

func DecodeEventContractCall(tx *rpc.BTCTransaction, evmConfigs map[int64]evmTypes.EVMConfig) (evmTypes.EventContractCall, error) {
	log.Infof("Decoding BTC transaction %+v\n", tx)

	// Decode the hex string into bytes
	txRaw, err := hex.DecodeString(tx.Data.Hex)
	if err != nil {
		log.Errorf("Failed to decode hex string %v", err)
		return evmTypes.EventContractCall{}, err
	}
	// Parse the transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)
	err = msgTx.Deserialize(bytes.NewReader(txRaw))
	if err != nil {
		log.Errorf("Failed to parse transaction %v", err)
		return evmTypes.EventContractCall{}, err
	}

	if len(msgTx.TxOut) < 3 {
		return evmTypes.EventContractCall{}, fmt.Errorf("btcLocking tx must have at least 3 outputs")
	}

	embeddedDataTxOut := msgTx.TxOut[1]
	if embeddedDataTxOut == nil || embeddedDataTxOut.PkScript == nil || embeddedDataTxOut.PkScript[0] != txscript.OP_RETURN {
		return evmTypes.EventContractCall{}, fmt.Errorf("transaction does not have expected payload op return output")
	}

	output, err := vault.ParseVaultEmbeddedData(embeddedDataTxOut.PkScript)
	if err != nil || output == nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot parse payload op return data: %w", err)
	}

	var senderBytes [20]byte
	copy(senderBytes[:], output.DestinationRecipientAddress[:])

	var txIdBytes [32]byte
	txId := msgTx.TxID()
	txBytes, err := hex.DecodeString(txId)
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("failed to decode tx id: %w", err)
	}
	copy(txIdBytes[:], txBytes)

	var mintingAmount int64 = msgTx.TxOut[0].Value

	_, payloadHash, err := evmUtils.CalculateStakingPayloadHash(senderBytes, mintingAmount, txIdBytes)
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("failed to get payload hash: %w", err)
	}

	parsedDestinationChain := chainUtils.NewDestinationChainFromBytes(output.DestinationChain)
	if parsedDestinationChain == nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("failed to parse destination chain")
	}

	destinationChain := nexus.ChainName(evmConfigs[int64(parsedDestinationChain.ChainID)].Name)
	contractAddress := hex.EncodeToString(output.DestinationContractAddress[:])

	log.Debugf("Encoded Data: %v\n", map[string]interface{}{
		"sender":        output.DestinationRecipientAddress,
		"mintingAmount": mintingAmount,
		"txId":          txId,
		"payloadHash":   payloadHash,
	})

	log.Infof("Encoded BTC info to EVM Call: %v\n", evmTypes.EventContractCall{
		Sender:           evmTypes.Address(output.DestinationRecipientAddress),
		DestinationChain: destinationChain, // not used
		ContractAddress:  contractAddress,
		PayloadHash:      evmTypes.Hash(common.BytesToHash(payloadHash)),
	})

	return evmTypes.EventContractCall{
		Sender:           evmTypes.Address(output.DestinationRecipientAddress),
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      evmTypes.Hash(common.BytesToHash(payloadHash)),
	}, nil
}
