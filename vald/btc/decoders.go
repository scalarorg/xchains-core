package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/btc/types"
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/axelarnetwork/utils/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	payloadData, err := NewOpReturnData(msgTx.TxOut)
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot parse payload op return data: %w", err)
	}

	if payloadData == nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("transaction does not have expected payload op return output")
	}

	chainId, err := utils.BytesToInt64BigEndian(payloadData.DestinationChainID[:])
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot parse chain id: %w", err)
	}

	sender := evmTypes.Address(common.BytesToAddress(payloadData.DestinationRecipientAddr[:]))
	// Find and Get the chain name
	// numberChainID := binary.BigEndian.Uint64(payloadData.ChainID)
	// destinationChain := nexus.ChainName(strconv.FormatUint(numberChainID, 10))
	// Todo: Scalar hardcoded chain name for now
	destinationChain := nexus.ChainName(evmConfigs[chainId].Name)
	// Get the contract address
	contractAddress := hex.EncodeToString(payloadData.DestinationContractAddr[:])
	// need "0x"?

	mintingAmount, err := GetMintingAmount(msgTx.TxOut[0])
	if err != nil {
		return evmTypes.EventContractCall{}, fmt.Errorf("cannot get minting amount: %w", err)
	}

	abi_minting_payload := types.Hash(common.BytesToHash(mintingAmount[:]))
	abi_address_payload := types.Hash(common.BytesToHash(payloadData.DestinationRecipientAddr[:]))

	abi_payload := append(abi_address_payload[:], abi_minting_payload[:]...)

	btcTxBlockTimeToBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(btcTxBlockTimeToBytes, uint64(tx.Data.Blocktime))

	btcTxBlockTime := types.Hash(common.BytesToHash(btcTxBlockTimeToBytes))

	payloadHash := evmTypes.Hash(common.BytesToHash(crypto.Keccak256(abi_payload, btcTxBlockTime[:])))

	log.Infof("Encoded BTC info to EVM Call: %v\n", evmTypes.EventContractCall{
		Sender:           sender,
		DestinationChain: destinationChain,
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
