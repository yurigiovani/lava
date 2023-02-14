package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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

	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagFrom)

	return cmd
}

func sendEnterLottery(clientCtx client.Context, txFactory tx.Factory, msg *types.MsgEnterLottery) error {
	return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txFactory, msg)
}

func sign(clientCtx client.Context, txBuilder client.TxBuilder, txFactory tx.Factory, from string) error {
	_, fromName, _, err := client.GetFromFields(clientCtx, txFactory.Keybase(), from)
	if err != nil {
		return fmt.Errorf("error getting account from keybase: %w", err)
	}

	if err = authclient.SignTx(txFactory, clientCtx, fromName, txBuilder, false, false); err != nil {
		return err
	}

	return nil
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgEnterLottery, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()

	//pkStr, err := fs.GetString(cli.FlagPubKey)
	//if err != nil {
	//	return txf, nil, err
	//}

	pkStr := "{\"@type\":\"/cosmos.crypto.ed25519.PubKey\",\"key\":\"E2Q//ZPEwt092mrYNZkY3UpN8Ioi1P0BKhhm3gFmmHg=\"}"

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	//moniker, _ := fs.GetString(FlagMoniker)
	//identity, _ := fs.GetString(FlagIdentity)
	//website, _ := fs.GetString(FlagWebsite)
	//security, _ := fs.GetString(FlagSecurityContact)
	//details, _ := fs.GetString(FlagDetails)
	//description := types.NewDescription(
	//	moniker,
	//	identity,
	//	website,
	//	security,
	//	details,
	//)

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
