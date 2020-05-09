package vm

import (
	"fmt"
	"math/big"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/netcloth/netcloth-chain/app/v0/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	GasLimit  uint64
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	StateDB   *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return st.StateDB.GetBalance(acc).Cmp(amount) >= 0
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) {
	st.StateDB.SubBalance(from, amount)
	st.StateDB.AddBalance(to, amount)
}

func (st StateTransition) GetHashFn(header abci.Header) func() sdk.Hash {
	return func() sdk.Hash {
		var res = sdk.Hash{}
		blockID := header.GetLastBlockId()
		res.SetBytes(blockID.GetHash())
		return res
	}
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context, vmParams *types.Params) (*big.Int, *sdk.Result, error) {
	ctx = ctx.WithLogger(ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName)))
	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHashFn(ctx.BlockHeader()),

		Origin: st.Sender,

		CoinBase:    ctx.BlockHeader().ProposerAddress, // validator consensus address, not account address
		Time:        sdk.NewInt(int64(ctx.BlockHeader().Time.Unix())).BigInt(),
		GasLimit:    st.GasLimit,
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	// This gas meter is set up to consume gas from gaskv during evm execution and be ignored
	currentGasMeter := ctx.GasMeter()
	csdb := st.StateDB.WithContext(ctx.WithGasMeter(sdk.NewInfiniteGasMeter())).WithTxHash(tmhash.Sum(ctx.TxBytes()))
	// Clear cache of accounts to handle changes outside of the EVM
	csdb.UpdateAccounts()

	cfg := Config{OpConstGasConfig: &vmParams.VMOpGasParams, CommonGasConfig: &vmParams.VMCommonGasParams}
	evm := NewEVM(evmCtx, csdb, cfg)

	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		vmerr       error
	)

	if st.Recipient.Empty() {
		ret, addr, leftOverGas, vmerr = evm.Create(st.Sender, st.Payload, st.GasLimit, st.Amount.BigInt())
		ctx.Logger().Info(fmt.Sprintf("create contract, consumed gas = %v , leftOverGas = %v, vm err = %v ", st.GasLimit-leftOverGas, leftOverGas, vmerr))
	} else {
		ret, leftOverGas, vmerr = evm.Call(st.Sender, st.Recipient, st.Payload, st.GasLimit, st.Amount.BigInt())

		if vmerr == ErrExecutionReverted && len(ret) > 4 {
			ctx.Logger().Info(fmt.Sprintf("VM revert error, reason provided by the contract: %v", string(ret[4:])))
		}

		ctx.Logger().Info(fmt.Sprintf("call contract, ret = %x, consumed gas = %v , leftOverGas = %v, vm err = %v", ret, st.GasLimit-leftOverGas, leftOverGas, vmerr))
	}

	vmGasUsed := st.GasLimit - leftOverGas

	if vmerr != nil {
		return nil, &sdk.Result{Data: ret, GasUsed: ctx.GasMeter().GasConsumed() + vmGasUsed}, vmerr
	}

	if ctx.GasMeter().GasConsumed()+vmGasUsed > ctx.GasMeter().Limit() {
		// vm rum out of gas
		ctx.Logger().Info("VM run out of gas")
		return nil, &sdk.Result{Data: ret, GasUsed: ctx.GasMeter().GasConsumed() + vmGasUsed}, ErrOutOfGas
	}

	// comsume vm gas
	ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(vmGasUsed, "EVM execution consumption")

	st.StateDB.Finalise(true)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeNewContract,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
		),
	})

	return nil, &sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}, nil
}

func DoStateTransition(ctx sdk.Context, msg types.MsgContract, k Keeper, readonly bool) (*big.Int, *sdk.Result, error) {
	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.To,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		StateDB:   k.StateDB.WithContext(ctx),
	}

	if readonly {
		st.StateDB = types.NewStateDB(k.StateDB).WithContext(ctx)
		st.GasLimit = DefaultGasLimit
	}

	params := k.GetParams(ctx)
	gasLimit := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	st.GasLimit = gasLimit
	return st.TransitionCSDB(ctx, &params)
}
