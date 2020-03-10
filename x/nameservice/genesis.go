package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	whoIsRecords []whoIs `json:"whois_records"`
}

func NewGenesisState(whoIsRecords []whoIs) GenesisState {
	return GenesisState{whoIsRecords: whoIsRecords}
}

func validateGensis(genState GenesisState) error {

	// Fetch & Iterate Through Names' whoIs
	for _, whoIs := range genState.whoIsRecords {

		if whoIs.Owner == nil {
			return fmt.Errorf("Invalid whoIsRecord: %s (Value) - Missing Owner", whoIs.Value)
		}

		if whoIs.Value == "" {
			return fmt.Errorf("Invalid whoIsRecord: %s (Owner) - Missing Value", whoIs.Owner)
		}

		if whoIs.Price == nil {
			return fmt.Errorf("Invalid whoIsRecord: %s (Value) - Missing Price", whoIs.Value)
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		whoIsRecords: []whoIs{},
	}
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, /* TODO: Define what keepers the module needs */, genState GenesisState) []abci.ValidatorUpdate {
	// TODO: Define logic for when you would like to initalize a new genesis

	// Fetch & Iterate Through Names' whoIs
	for _, whoIs := range genState.whoIsRecords {
		// Assign whoIs Structures
		keeper.setWhoIs(ctx, whoIs.Value, whoIs)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (GenesisState) {
	// TODO: Define logic for exporting state
	var names []whoIs

	// Retrieve All The Names
	iterator := k.getNamesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {

		// Get Key (Name)
		key := string(iterator.Key())

		// Get whoIs Of Name
		whois := k.getWhoIs(ctx, key)

		// Append To Names List
		names = append(names, whois)
	}

	return GenesisState{whoIsRecords: names}
}
