package btc

import (
	goerrors "errors"

	"github.com/axelarnetwork/axelar-core/sdk-utils/broadcast"
	"github.com/axelarnetwork/axelar-core/utils/errors"
	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	btctypes "github.com/axelarnetwork/axelar-core/x/btc/types"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/axelarnetwork/utils/slices"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ErrNotFinalized is returned when a transaction is not finalized
var ErrNotFinalized = goerrors.New("not finalized")

// ErrTxFailed is returned when a transaction has failed
var ErrTxFailed = goerrors.New("transaction failed")

// Mgr manages all communication with Bitcoin
type Mgr struct {
	rpc         rpc.Client
	broadcaster broadcast.Broadcaster
	validator   sdk.ValAddress
	proxy       sdk.AccAddress
	config      btctypes.BTCConfig

	// TODO: evmConfigs is a map of chainID EVMConfig -> refactor ChainID as type
	evmConfigs map[int64]types.EVMConfig
}

// NewMgr returns a new Mgr instance
func NewMgr(btcConfig btctypes.BTCConfig, cliCtx sdkClient.Context, broadcaster broadcast.Broadcaster, valAddr sdk.ValAddress, proxy sdk.AccAddress) (*Mgr, error) {

	btcClientLogger := log.WithKeyVals("chain", btcConfig.Name, "url", btcConfig.Host)
	rpcConfig := btcConfig.GetRPCConfig()
	if rpcConfig.Host == "" {
		return nil, sdkerrors.New("btc", 1, "btc rpc host is empty") // nolint: goerr113
	}
	btcClient, err := rpc.NewClient(rpcConfig, btcClientLogger)
	if err != nil {
		err = sdkerrors.Wrap(err, "failed to create a BTC RPC client")
		log.Error(err.Error())
		return nil, err
	}
	return &Mgr{
		config:      btcConfig,
		broadcaster: broadcaster,
		validator:   valAddr,
		proxy:       proxy,
		rpc:         btcClient,

		// FIXME: evmConfigs is not initialized
		evmConfigs: make(map[int64]types.EVMConfig),
	}, nil
}

func (mgr Mgr) logger(keyvals ...any) log.Logger {
	keyvals = append([]any{"listener", "btc"}, keyvals...)
	return log.WithKeyVals(keyvals...)
}

// isParticipantOf checks if the validator is in the poll participants list
func (mgr Mgr) isParticipantOf(participants []sdk.ValAddress) bool {
	return slices.Any(participants, func(v sdk.ValAddress) bool { return v.Equals(mgr.validator) })
}

func (mgr Mgr) isFinalized(tx rpc.BTCTransaction, confHeight int64) (bool, error) {
	if confHeight < 0 {
		return false, goerrors.New("ConfHeight cannot less than zero")
	}

	if tx.Data.Confirmations < uint64(confHeight) {
		return false, nil
	}

	return true, nil
}

func (mgr Mgr) GetTxIfFinalized(txID types.Hash, confHeight uint64) (results.Result[rpc.BTCTransaction], error) {
	txs, err := mgr.GetTxsIfFinalized([]types.Hash{txID}, confHeight)
	if err != nil {
		return results.Result[rpc.BTCTransaction]{}, err
	}

	return txs[0], err
}

func (mgr Mgr) GetTxsIfFinalized(txIDs []types.Hash, confHeight uint64) ([]results.Result[rpc.BTCTransaction], error) {
	txs, err := mgr.rpc.GetTransactions(txIDs)

	if err != nil {
		return slices.Map(txIDs, func(_ types.Hash) results.Result[rpc.BTCTransaction] {
			return results.FromErr[rpc.BTCTransaction](
				sdkerrors.Wrapf(
					errors.With(err, "chain", mgr.config.Name, "tx_ids", txIDs),
					"cannot get transaction receipts"),
			)
		}), nil
	}

	return slices.Map(txs, func(tx rpc.TxResult) results.Result[rpc.BTCTransaction] {
		return results.Pipe(results.Result[rpc.BTCTransaction](tx), func(tx rpc.BTCTransaction) results.Result[rpc.BTCTransaction] {

			isFinalized, err := mgr.isFinalized(tx, int64(confHeight))
			if err != nil {
				return results.FromErr[rpc.BTCTransaction](sdkerrors.Wrapf(errors.With(err, "chain", mgr.config.Name),
					"cannot determine if the transaction %s is finalized", tx.Data.Txid),
				)
			}

			if !isFinalized {
				return results.FromErr[rpc.BTCTransaction](ErrNotFinalized)
			}

			// TODO_SCALAR: Check if bitcoin transactions can fail after being added to mempool,
			// and if so, check if we need to check the status of the transaction

			// if receipt.Status != geth.ReceiptStatusSuccessful {
			// 	return results.FromErr[geth.Receipt](ErrTxFailed)
			// }

			return results.FromOk[rpc.BTCTransaction](tx)
		})

	}), nil
}
