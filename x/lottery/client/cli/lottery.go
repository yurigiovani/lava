package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// NewLotteryCmd returns a root CLI command handler for all x/staking transaction commands.
func NewLotteryCmd() *cobra.Command {
	lotteryTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Lottery subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	lotteryTxCmd.AddCommand(
		NewEnterLotteryCmd(),
	)

	flags.AddTxFlagsToCmd(lotteryTxCmd)

	return lotteryTxCmd
}

func NewEnterLotteryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enter",
		Short: "Enter the lottery",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}

			err = sendEnterLottery(clientCtx, cmd)

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(FlagBet, "", "The bet that will send to lottery (min 1stake)")
	cmd.Flags().String(FlagFrom, "", "The address that will enter on lottery")
	cmd.Flags().String(FlagPubKey, "", "The pubkey from address")

	_ = cmd.MarkFlagRequired(FlagBet)
	_ = cmd.MarkFlagRequired(FlagFrom)
	_ = cmd.MarkFlagRequired(FlagPubKey)

	return cmd
}

// senEnterLottery is to build the message and send the broadcast the transaction
func sendEnterLottery(clientCtx client.Context, cmd *cobra.Command) error {
	txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
		WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

	txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())

	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
}

// newBuildCreateValidatorMsg to build the message types.MsgEnterLottery
func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgEnterLottery, error) {
	fBet, _ := fs.GetString(FlagBet)
	bet, err := sdk.ParseCoinNormalized(fBet)

	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(FlagPubKey)

	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	lotteryFee := sdk.NewCoin("stake", math.NewInt(types.LotteryFee))

	msg := types.NewMsgEnterLottery(
		valAddr.String(), bet, lotteryFee,
	)

	if err != nil {
		return txf, nil, err
	}

	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, &msg, nil
}
