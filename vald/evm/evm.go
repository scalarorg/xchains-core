package evm

import (
	"context"
	goerrors "errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"

	"github.com/axelarnetwork/axelar-core/sdk-utils/broadcast"
	"github.com/axelarnetwork/axelar-core/utils/errors"
	btc "github.com/axelarnetwork/axelar-core/vald/btc"
	"github.com/axelarnetwork/axelar-core/vald/evm/rpc"
	btcTypes "github.com/axelarnetwork/axelar-core/x/btc/types"
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/axelarnetwork/utils/log"
	"github.com/axelarnetwork/utils/monads/results"
	"github.com/axelarnetwork/utils/slices"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
)

// ErrNotFinalized is returned when a transaction is not finalized
var ErrNotFinalized = goerrors.New("not finalized")

// ErrTxFailed is returned when a transaction has failed
var ErrTxFailed = goerrors.New("transaction failed")

// Mgr manages all communication with Ethereum
type Mgr struct {
	rpcs                      map[string]rpc.Client
	broadcaster               broadcast.Broadcaster
	validator                 sdk.ValAddress
	proxy                     sdk.AccAddress
	latestFinalizedBlockCache LatestFinalizedBlockCache
	//TODO: btcMgr is a map of chainID to btc.Mgr -> refactor ChainID as type
	btcMgrs map[string]*btc.Mgr
	//TODO: evmConfigs is a map of chainID to EVMConfig -> refactor ChainID as type
	evmConfigs map[int64]evmTypes.EVMConfig
}

// NewMgr returns a new Mgr instance
func NewMgr(
	clientCtx sdkClient.Context,
	rpcs map[string]rpc.Client,
	broadcaster broadcast.Broadcaster,
	valAddr sdk.ValAddress,
	caches LatestFinalizedBlockCache,
	evmConfigs []evmTypes.EVMConfig,
	btcConfigs []btcTypes.BTCConfig) *Mgr {

	btcMgrs := make(map[string]*btc.Mgr)

	for _, cfg := range btcConfigs {
		_, ok := btcMgrs[cfg.ChainID]
		if ok {
			err := sdkerrors.Wrap(goerrors.New("duplicate chain id"), fmt.Sprintf("duplicate chain id %s", cfg.ChainID))
			log.Error(err.Error())
			panic(err)
		}

		mgr, err := btc.NewMgr(cfg, clientCtx, broadcaster, valAddr, clientCtx.FromAddress)
		if err != nil {
			err = sdkerrors.Wrap(err, "failed to create an RPC connection for Btc chain. Verify your BTC Connection config.")
			log.Error(err.Error())
			panic(err)
		}

		btcMgrs[cfg.ChainID] = mgr
	}

	selfConfigs := make(map[int64]evmTypes.EVMConfig)

	for _, cfg := range evmConfigs {
		_, ok := selfConfigs[cfg.ChainID]
		if ok {
			err := sdkerrors.Wrap(goerrors.New("duplicate chain id"), fmt.Sprintf("duplicate chain id %d", cfg.ChainID))
			log.Error(err.Error())
			panic(err)
		}

		selfConfigs[cfg.ChainID] = cfg
	}

	return &Mgr{
		rpcs:                      rpcs,
		broadcaster:               broadcaster,
		validator:                 valAddr,
		proxy:                     clientCtx.FromAddress,
		latestFinalizedBlockCache: caches,
		evmConfigs:                selfConfigs,
		btcMgrs:                   btcMgrs,
	}
}

func (mgr Mgr) GetBtcMgr(chainID string) (*btc.Mgr, error) {
	btcMgr, ok := mgr.btcMgrs[chainID]
	if !ok {
		return nil, sdkerrors.Wrap(goerrors.New("chain id not found"), fmt.Sprintf("chain id %s not found", chainID))
	}

	return btcMgr, nil
}

func (mgr Mgr) logger(keyvals ...any) log.Logger {
	keyvals = append([]any{"listener", "evm"}, keyvals...)
	return log.WithKeyVals(keyvals...)
}

// ProcessNewChain notifies the operator that vald needs to be restarted/udpated for a new chain
func (mgr Mgr) ProcessNewChain(event *evmTypes.ChainAdded) (err error) {
	mgr.logger().Info(fmt.Sprintf("VALD needs to be updated and restarted for new chain %s", event.Chain.String()))
	return nil
}

func (mgr Mgr) isFinalized(chain nexus.ChainName, txReceipt geth.Receipt, confHeight uint64) (bool, error) {
	client, ok := mgr.rpcs[strings.ToLower(chain.String())]
	if !ok {
		return false, fmt.Errorf("rpc client not found for chain %s", chain.String())
	}

	if mgr.latestFinalizedBlockCache.Get(chain).Cmp(txReceipt.BlockNumber) >= 0 {
		return true, nil
	}

	latestFinalizedBlockNumber, err := client.LatestFinalizedBlockNumber(context.Background(), confHeight)
	if err != nil {
		return false, err
	}

	mgr.latestFinalizedBlockCache.Set(chain, latestFinalizedBlockNumber)

	if latestFinalizedBlockNumber.Cmp(txReceipt.BlockNumber) < 0 {
		return false, nil
	}

	return true, nil
}

// GetTxReceiptIfFinalized retrieves receipt for provided transaction ID.
//
// # Result is
//
// - Ok(receipt) if the transaction is finalized and successful
//
// - Err(ethereum.NotFound) if the transaction is not found
//
// - Err(ErrTxFailed) if the transaction is finalized but failed
//
// - Err(ErrNotFinalized) if the transaction is not finalized
//
// - Err(err) otherwise
func (mgr Mgr) GetTxReceiptIfFinalized(chain nexus.ChainName, txID common.Hash, confHeight uint64) (results.Result[geth.Receipt], error) {
	txReceipts, err := mgr.GetTxReceiptsIfFinalized(chain, []common.Hash{txID}, confHeight)
	if err != nil {
		return results.Result[geth.Receipt]{}, err
	}

	return txReceipts[0], err
}

// GetTxReceiptsIfFinalized retrieves receipts for provided transaction IDs.
//
// # Individual result is
//
// - Ok(receipt) if the transaction is finalized and successful
//
// - Err(ethereum.NotFound) if the transaction is not found
//
// - Err(ErrTxFailed) if the transaction is finalized but failed
//
// - Err(ErrNotFinalized) if the transaction is not finalized
//
// - Err(err) otherwise
func (mgr Mgr) GetTxReceiptsIfFinalized(chain nexus.ChainName, txIDs []common.Hash, confHeight uint64) ([]results.Result[geth.Receipt], error) {
	client, ok := mgr.rpcs[strings.ToLower(chain.String())]
	if !ok {
		return nil, fmt.Errorf("rpc client not found for chain %s", chain.String())
	}

	receipts, err := client.TransactionReceipts(context.Background(), txIDs)
	if err != nil {
		return slices.Map(txIDs, func(_ common.Hash) results.Result[geth.Receipt] {
			return results.FromErr[geth.Receipt](
				sdkerrors.Wrapf(
					errors.With(err, "chain", chain.String(), "tx_ids", txIDs),
					"cannot get transaction receipts"),
			)
		}), nil
	}

	return slices.Map(receipts, func(receipt rpc.TxReceiptResult) results.Result[geth.Receipt] {
		return results.Pipe(results.Result[geth.Receipt](receipt), func(receipt geth.Receipt) results.Result[geth.Receipt] {

			isFinalized, err := mgr.isFinalized(chain, receipt, confHeight)
			if err != nil {
				return results.FromErr[geth.Receipt](sdkerrors.Wrapf(errors.With(err, "chain", chain.String()),
					"cannot determine if the transaction %s is finalized", receipt.TxHash.Hex()),
				)
			}

			if !isFinalized {
				return results.FromErr[geth.Receipt](ErrNotFinalized)
			}

			if receipt.Status != geth.ReceiptStatusSuccessful {
				return results.FromErr[geth.Receipt](ErrTxFailed)
			}

			return results.FromOk[geth.Receipt](receipt)
		})
	}), nil
}

// isParticipantOf checks if the validator is in the poll participants list
func (mgr Mgr) isParticipantOf(participants []sdk.ValAddress) bool {
	return slices.Any(participants, func(v sdk.ValAddress) bool { return v.Equals(mgr.validator) })
}
