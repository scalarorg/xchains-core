package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"
)

func DecodeEventContractCall(tx *rpc.BTCTransaction, evmConfigs map[int64]evmTypes.EVMConfig) (evmTypes.EventContractCall, error) {
	// TODO_SCALAR: Parse the tx data to extract the sender, destination chain, and contract address
	// sender := types.Address(common.BytesToAddress([]byte(tx.Data.Hash)))
	// destinationChain := nexus.ChainName("ethereum-sepolia")
	// contractAddress := "0xe432150cce91c13a887f7D836923d5597adD8E31"

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
	payload, err := NewOpReturnData(msgTx.TxOut)
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot parse payload op return data: %w", err)
	}

	if payload == nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("transaction does not have expected payload op return output")
	}

	chainId, err := utils.BytesToInt64BigEndian(payload.DestinationChainID[:])
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot parse chain id: %w", err)
	}

	sender := payload.DestinationRecipientAddr
	mintingAmount := msgTx.TxOut[0].Value
	blockTime := tx.Data.Blocktime

	payloadHash, err := GetPayloadHash(sender, mintingAmount, blockTime)
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("failed to get payload hash: %w", err)
	}

	destinationChain := nexus.ChainName(evmConfigs[chainId].Name)
	contractAddress := hex.EncodeToString(payload.DestinationContractAddr[:])

	log.Debugf("Encoded Data: %v\n", map[string]interface{}{
		"sender":        sender,
		"mintingAmount": mintingAmount,
		"blockTime":     blockTime,
		"payloadHash":   payloadHash,
	})

	log.Infof("Encoded BTC info to EVM Call: %v\n", evmTypes.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain, // not used
		ContractAddress:  contractAddress,
		PayloadHash:      payloadHash,
	})

	return evmTypes.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
		ContractAddress:  contractAddress,
		PayloadHash:      payloadHash,
	}, nil
}
