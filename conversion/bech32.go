package conversion

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SetupBech32Prefix() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("inv", "invpub")
	config.SetBech32PrefixForValidator("invval", "invvalpub")
	config.SetBech32PrefixForConsensusNode("invvalcons", "invvalconspub")
}
