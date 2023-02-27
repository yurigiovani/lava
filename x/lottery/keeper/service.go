package keeper

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/tendermint/go-amino"
)

// store returns a lottery store.
func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// addEntry method to add some msg into a desired lottery by id
func (k Keeper) addEntry(ctx sdk.Context, id int64, msg *types.MsgEnterLottery) error {
	entries, err := k.getCurrentLottery(ctx)

	if err != nil {
		return err
	}

	entries = append(entries, msg)
	bentries, err := amino.MarshalBinaryLengthPrefixed(entries)

	if err != nil {
		return err
	}

	k.store(ctx).Set(types.GetLotteryEntriesKey(id), bentries)

	return nil
}

func (k Keeper) getCurrentLottery(ctx sdk.Context) (types.MsgEnterLotteryList, error) {
	var currentID = k.getCurrentLotteryID(ctx)
	var entries types.MsgEnterLotteryList
	lotteryStore := k.store(ctx)
	lotteries := lotteryStore.Get(types.GetLotteryEntriesKey(currentID))

	if len(lotteries) <= 0 {
		return entries, nil
	}

	if err := amino.UnmarshalBinaryLengthPrefixed(lotteries, &entries); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error to get lottery: %s", err))
		return nil, err
	}

	return entries, nil
}

func (k Keeper) hasEntryLottery(ctx sdk.Context, address string) bool {
	var currentID = k.getCurrentLotteryID(ctx)
	var msgs types.MsgEnterLotteryList
	lotteryStore := k.store(ctx)
	iterator := sdk.KVStorePrefixIterator(lotteryStore, types.GetLotteryEntriesKey(currentID))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		err := amino.UnmarshalBinaryLengthPrefixed(lotteryStore.Get(iterator.Key()), &msgs)

		if err != nil {
			fmt.Println("error to amino.UnmarshalBinaryLengthPrefixed", err)
			return false
		}
	}

	for _, msg := range msgs {
		if msg.Address != address {
			continue
		}

		return true
	}

	return false
}

func (k Keeper) incrementCounter(ctx sdk.Context) {
	counter := k.GetCounter(ctx)
	k.store(ctx).Set(types.KeyLotteryCounter, []byte{byte(counter)})
}

func (k Keeper) getCurrentLotteryID(ctx sdk.Context) int64 {
	bid := k.store(ctx).Get(types.KeyLotteryCurrentyID)
	var id int64 = 1

	if bid != nil {
		id = int64(binary.BigEndian.Uint64(bid))
	}

	return id
}
