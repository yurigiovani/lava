package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper creates a new lottery Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
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

func (k Keeper) ResetCounter(ctx sdk.Context) {
	k.store(ctx).Set(types.KeyLotteryCounter, []byte{0})
}

// DrawLottery method to draw lottery and choose the winner of the lottery
func (k Keeper) DrawLottery(ctx sdk.Context) (*types.MsgEnterLottery, error) {
	msgLotteryList, err := k.getCurrentLottery(ctx)

	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not getCurrentLottery: %s", err))
		return nil, err
	}

	winnerMsg := k.chooseWinner(ctx, msgLotteryList)
	return winnerMsg, nil
}

func (k Keeper) chooseWinner(ctx sdk.Context, msgs types.MsgEnterLotteryList) *types.MsgEnterLottery {
	var dataHash []byte

	for _, msg := range msgs {
		bmsg, err := msg.Marshal()

		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("some error while Marshal a msg: %s", err))
			continue
		}

		dataHash = append(dataHash, bmsg...)
	}

	hashResult := crypto.Sha256(dataHash)
	index := (binary.BigEndian.Uint16(hashResult) ^ 0xFFFF) % uint16(len(msgs))

	return msgs[index]
}

func (k Keeper) Payout(ctx sdk.Context, winner types.MsgEnterLottery) (bool, error) {
	msgs, _ := k.getCurrentLottery(ctx)

	isPayable, amount := k.calculatePayout(ctx, winner, msgs)

	if isPayable == false {
		return false, nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(winner.Bet.Denom, amount))

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winner.GetSigners()[0], coins); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not send module from Lottery pool to account %s: %s", winner.Address, err))

		return false, err
	}

	return true, nil
}
