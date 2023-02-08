package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// NewLotteryCmd returns a root CLI command handler for all x/staking transaction commands.
func NewLotteryCmd() *cobra.Command {
	lotteryTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	lotteryTxCmd.AddCommand(
		NewEnterLotteryCmd(),
	)

	return lotteryTxCmd
}

func NewEnterLotteryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enter [address] [amount] [tokenDenom]",
		Short: "Enter the lottery",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1] + " " + args[2])

			if err != nil {
				return err
			}

			txResponse, err := sendEnterLottery(clientCtx, cmd.Flags(), address, amount, args[2])

			if err != nil {
				return err
			}

			txRaw := []byte(txResponse.RawLog)

			return clientCtx.PrintRaw(txRaw)
		},
	}
	return cmd
}

func sendEnterLottery(clientCtx client.Context, flags *flag.FlagSet, address sdk.AccAddress, amount sdk.Coin, tokenDenom string) (sdk.TxResponse, error) {
	msg := types.NewMsgEnterLottery(address, amount, tokenDenom)

	txFactory := tx.NewFactoryCLI(clientCtx, flags)
	txFactory = txFactory.WithAccountNumber(0).WithSequence(0)

	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	// set the new appened msgs into builder
	txBuilder.SetMsgs(msg)

	// set the memo,fees,feeGranter,feePayer from cmd flags
	txBuilder.SetMemo(txFactory.Memo())
	txBuilder.SetFeeAmount(txFactory.Fees())
	txBuilder.SetFeeGranter(clientCtx.FeeGranter)
	txBuilder.SetFeePayer(clientCtx.FeePayer)

	// set the gasLimit
	txBuilder.SetGasLimit(txFactory.Gas())

	sign(clientCtx, txBuilder, txFactory, "from")

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())

	res, err := clientCtx.BroadcastTx(txBytes)

	return *res, err
}

func sign(clientCtx client.Context, txBuilder client.TxBuilder, txFactory tx.Factory, from string) error {
	_, fromName, _, err := client.GetFromFields(clientCtx, txFactory.Keybase(), from)
	if err != nil {
		return fmt.Errorf("error getting account from keybase: %w", err)
	}

	if err = authclient.SignTx(txFactory, clientCtx, fromName, txBuilder, true, true); err != nil {
		return err
	}

	return nil
}
