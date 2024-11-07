package btc

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// OP_RETURN
// OP_PUSHBYTES_1 35
// OP_PUSHBYTES_4 53180104
// OP_PUSHBYTES_1 00
// OP_PUSHBYTES_8 0000000000aa36a7
// OP_PUSHBYTES_20 b91e3a8ef862567026d6f376c9f3d6b814ca4337
// OP_PUSHBYTES_20 24a1db57fa3ecafcbad91d6ef068439aceeae090

type OpReturnData struct {
	Tag                      [4]byte
	Version                  byte
	DestinationChainID       [8]byte
	DestinationRecipientAddr [20]byte
	DestinationContractAddr  [20]byte
}

func NewOpReturnData(txOuts []*wire.TxOut) (*OpReturnData, error) {

	selectedTxOut := &wire.TxOut{}

	// find the op return tx out
	for _, txOut := range txOuts {
		if txOut.PkScript[0] == txscript.OP_RETURN {
			selectedTxOut = txOut
			break
		}
	}

	if selectedTxOut == nil {
		return nil, errors.New("no op return tx out found")
	}

	data := selectedTxOut.PkScript[1:]

	if len(data) < 59 { // Adjusted for OP codes and lengths
		return nil, fmt.Errorf("data too short: got %d bytes, expected at least 59", len(data))
	}

	var opReturn OpReturnData

	// Parse the OP_RETURN data considering opcode lengths
	// Skip `OP_RETURN` (1 byte) + `OP_PUSHBYTES_1` (1 byte)
	idx := 2

	// Skip over `OP_PUSHBYTES_4` (1 byte) and map 4 bytes to Tag
	copy(opReturn.Tag[:], data[idx+1:idx+5])
	idx += 5 // Move index by 4 bytes for Tag + 1 for opcode

	// Skip over `OP_PUSHBYTES_1` (1 byte) and map 1 byte to Version
	opReturn.Version = data[idx+1]
	idx += 2 // Move index by 1 byte for Version + 1 for opcode

	// Skip over `OP_PUSHBYTES_8` (1 byte) and map 8 bytes to DestinationChainID
	copy(opReturn.DestinationChainID[:], data[idx+1:idx+9])
	idx += 9 // Move index by 8 bytes for ChainID + 1 for opcode

	// Skip over `OP_PUSHBYTES_20` (1 byte) and map 20 bytes to DestinationRecipientAddr
	copy(opReturn.DestinationRecipientAddr[:], data[idx+1:idx+21])
	idx += 21 // Move index by 20 bytes for RecipientAddr + 1 for opcode

	// Skip over `OP_PUSHBYTES_20` (1 byte) and map 20 bytes to DestinationContractAddr
	copy(opReturn.DestinationContractAddr[:], data[idx+1:idx+21])
	idx += 21 // Move index by 20 bytes for ContractAddr + 1 for opcode

	return &opReturn, nil
}

func GetMintingAmount(txOut *wire.TxOut) ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(txOut.Value))
	return buf, nil
}
