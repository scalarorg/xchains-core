package btc

import (
	goerrors "errors"

	"github.com/axelarnetwork/axelar-core/sdk-utils/broadcast"
	"github.com/axelarnetwork/axelar-core/utils/errors"
	"github.com/axelarnetwork/axelar-core/vald/btc/rpc"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/axelarnetwork/utils/slices"
	"github.com/btcsuite/btcd/rpcclient"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CHAIN_BITCOIN = "bitcoin"
)

// ErrNotFinalized is returned when a transaction is not finalized
var ErrNotFinalized = goerrors.New("not finalized")

// ErrTxFailed is returned when a transaction has failed
var ErrTxFailed = goerrors.New("transaction failed")

// Mgr manages all communication with Bitcoin
type Mgr struct {
	rpc    rpc.Client
	config rpcclient.ConnConfig
}

// NewMgr returns a new Mgr instance
func NewMgr(btcConfig rpcclient.ConnConfig, cliCtx sdkClient.Context, b broadcast.Broadcaster, valAddr sdk.ValAddress) (*Mgr, error) {

	btcClientLogger := log.WithKeyVals("chain", CHAIN_BITCOIN, "url", btcConfig.Host)
	btcClient, err := rpc.NewClient(&btcConfig, btcClientLogger)
	if err != nil {
		err = sdkerrors.Wrap(err, "failed to create a BTC RPC client")
		log.Error(err.Error())
		return nil, err
	}
	return &Mgr{
		config: btcConfig,
		rpc:    btcClient,
	}, nil
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
					errors.With(err, "chain", CHAIN_BITCOIN, "tx_ids", txIDs),
					"cannot get transaction receipts"),
			)
		}), nil
	}

	return slices.Map(txs, func(tx rpc.TxResult) results.Result[rpc.BTCTransaction] {
		return results.Pipe(results.Result[rpc.BTCTransaction](tx), func(tx rpc.BTCTransaction) results.Result[rpc.BTCTransaction] {

			isFinalized, err := mgr.isFinalized(tx, int64(confHeight))
			if err != nil {
				return results.FromErr[rpc.BTCTransaction](sdkerrors.Wrapf(errors.With(err, "chain", CHAIN_BITCOIN),
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