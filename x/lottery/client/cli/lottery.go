package cli

import (
	"fmt"
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

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())

			if err != nil {
				return err
			}

			err = sendEnterLottery(clientCtx, txf, msg)

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(FlagAmount, "", "The amount that will send to lottery")
	cmd.Flags().String(FlagFrom, "", "The address that will enter on lottery")
	cmd.Flags().String(FlagPubKey, "", "The pubkey from address")

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagFrom)
	_ = cmd.MarkFlagRequired(FlagPubKey)

	return cmd
}

func sendEnterLottery(clientCtx client.Context, txFactory tx.Factory, msg *types.MsgEnterLottery) error {
	return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txFactory, msg)
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgEnterLottery, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()

	if err != nil {
		fmt.Println("error while GetAccount", err)
		return txf, nil, err
	}

	pkStr, err := fs.GetString(FlagPubKey)

	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	msg := types.NewMsgEnterLottery(
		valAddr, amount.Amount.Int64(), amount.Denom,
	)

	if err != nil {
		return txf, nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, &msg, nil
}
