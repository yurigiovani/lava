package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	authKeeper stakingtypes.AccountKeeper
	bankKeeper stakingtypes.BankKeeper
}

// NewKeeper creates a new lottery Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// EnterLottery receive types.MsgEnterLottery and set a new entry for current lottery
func (k Keeper) EnterLottery(ctx sdk.Context, msg *types.MsgEnterLottery) error {
	ctx.Logger().With("address", msg.Address).Info("started to enter on lottery", msg)
	defer ctx.Logger().Info("finished to enter on lottery", msg)

	if k.hasEntryLottery(ctx, msg.Address) {
		err := errors.New("address is already participating from lottery")
		k.Logger(ctx).Error(err.Error())

		return err
	}

	if err := k.addEntry(ctx, 1, msg); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not set lottery: %s", err))

		return err
	}

	k.incrementCounter(ctx)

	return nil
}

func (k Keeper) GetLastLottery(ctx sdk.Context) types.MsgEnterLotteryList {
	k.Logger(ctx).Info("getting last lottery")
	return types.MsgEnterLotteryList{}
}

// GetCounter method to get counter state from current lottery
func (k Keeper) GetCounter(ctx sdk.Context) int64 {
	bcounter := k.store(ctx).Get(types.KeyLotteryCounter)
	var counter int64 = 1

	if bcounter != nil {
		counter = int64(bcounter[0])
	}

	return counter
}

// DrawLottery method to draw lottery and choose the winner of the lottery
func (k Keeper) DrawLottery(ctx sdk.Context) (sdk.Address, error) {
	msgLotteryList, err := k.getCurrentLottery(ctx)

	if err != nil {
		return nil, err
	}

	winnerIndex := chooseWinner(ctx.BlockHeader().DataHash, len(msgLotteryList))
	winner := msgLotteryList[winnerIndex].Address

	return sdk.AccAddress(winner), nil
}

func chooseWinner(dataHash []byte, numTransactions int) uint16 {
	hashResult := crypto.Sha256(dataHash)
	index := (binary.BigEndian.Uint16(hashResult) ^ 0xFFFF) % uint16(numTransactions)

	return index
}
