package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	testutilbank "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/keeper"
	"github.com/cosmos/cosmos-sdk/x/lottery/testutil"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func createTestApp(t *testing.T) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	return app, ctx
}

func createExpectedKeepers(ctx sdk.Context, app *simapp.SimApp) (types.AccountKeeper, types.BankKeeper) {
	maccPerms := simapp.GetMaccPerms()
	appCodec := simapp.MakeTestEncodingConfig().Codec
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms, sdk.Bech32MainPrefix,
	)
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, app.GetKey(types.StoreKey), accountKeeper, app.GetSubspace(banktypes.ModuleName), nil)

	return accountKeeper, bankKeeper
}

func TestCalculatePayoutLowestBet(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))

	msgs := types.MsgEnterLotteryList{}

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))
	msgs = append(msgs, &msg1, &msg2)

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])

	require.Nil(t, err1)
	require.Nil(t, err2)

	isPayable, err := k.Payout(ctx, *msgs[0])

	require.Nil(t, err)
	require.False(t, isPayable, "this msg must not be paid")

	balAddr0 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], "stake")
	balAddr1 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], "stake")

	require.Equal(t, math.NewInt(39), balAddr0.Amount)
	require.Equal(t, math.NewInt(34), balAddr1.Amount)
}

func TestCalculatePayoutLowestHighestBet(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))

	msgs := types.MsgEnterLotteryList{}

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))

	msgs = append(msgs, &msg1, &msg2)

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])

	require.Nil(t, err1)
	require.Nil(t, err2)

	isPayable, err := k.Payout(ctx, *msgs[1])

	require.Nil(t, err)
	require.True(t, isPayable, "this msg must be paid")

	balAddr0 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], "stake")
	balAddr1 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], "stake")

	require.Equal(t, math.NewInt(39), balAddr0.Amount)
	require.Equal(t, math.NewInt(61), balAddr1.Amount)
}

func TestCalculatePayoutBetKeepingFeeOnLotteryPool(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 3, math.NewInt(0))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[2], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))

	msgs := types.MsgEnterLotteryList{}

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))
	msg3 := types.NewMsgEnterLottery(addrs[2].String(), sdk.NewCoin("stake", math.NewInt(25)), sdk.NewCoin("stake", math.NewInt(1)))

	msgs = append(msgs, &msg1, &msg2, &msg3)

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])
	err3 := k.EnterLottery(ctx, msgs[2])

	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Nil(t, err3)

	isPayable, err := k.Payout(ctx, *msgs[1])

	require.Nil(t, err)
	require.True(t, isPayable, "this msg must be paid")

	balAddr0 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], "stake")
	balAddr1 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], "stake")
	balAddr2 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[2], "stake")

	require.Equal(t, math.NewInt(39), balAddr0.Amount)
	require.Equal(t, math.NewInt(84), balAddr1.Amount)
	require.Equal(t, math.NewInt(24), balAddr2.Amount)
}
