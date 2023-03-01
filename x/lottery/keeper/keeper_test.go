package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	testutilbank "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/lottery/keeper"
	"github.com/cosmos/cosmos-sdk/x/lottery/testutil"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnterLotterySuccess(t *testing.T) {
	app, ctx := createTestApp(t)

	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)
	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	require.NoError(t, k.EnterLottery(ctx, &msg1))
	require.NoError(t, k.EnterLottery(ctx, &msg2))

	bal := testutil.GetModuleBalance(bankKeeper.(bankkeeper.Keeper), accountKeeper, ctx, types.ModuleName, "stake")

	require.Equal(t, math.NewInt(27), bal.Amount)

	balAddr0 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], "stake")
	balAddr1 := testutil.GetAccountBalance(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], "stake")

	require.Equal(t, math.NewInt(39), balAddr0.Amount)
	require.Equal(t, math.NewInt(34), balAddr1.Amount)
}

func TestEnterLotteryErrorParticipatingTwice(t *testing.T) {
	app, ctx := createTestApp(t)

	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)
	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	require.NoError(t, k.EnterLottery(ctx, &msg1))
	require.NoError(t, k.EnterLottery(ctx, &msg2))
	require.Error(t, k.EnterLottery(ctx, &msg2))
}

func TestEnterLotteryErrorBetLessThanAllowed(t *testing.T) {
	app, ctx := createTestApp(t)

	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)
	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(500))

	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(50))))
	testutilbank.FundAccount(bankKeeper.(bankkeeper.Keeper), ctx, addrs[1], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(3))))

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(5)), sdk.NewCoin("stake", math.NewInt(1)))

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)

	require.NoError(t, k.EnterLottery(ctx, &msg1))
	require.Error(t, k.EnterLottery(ctx, &msg2))
}
