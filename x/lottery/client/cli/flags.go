package cli

import flag "github.com/spf13/pflag"

const (
	FlagAddress = "address"
	FlagAmount  = "amount"
	FlagFrom    = "from"
	FlagPubKey  = "pubkey"
)

// FlagSetAddress Returns the flagset for External Address related operations.
func FlagSetAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagAddress, "", "The external path containing the path for external key address to check your node license validator")
	return fs
}
