package nameservice

import (
	"fmt"

	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	whoIsRecords []types.WhoIs `json:"whois_records"`
}

func NewGenesisState(whoIsRecords []types.WhoIs) GenesisState {
	return GenesisState{whoIsRecords: whoIsRecords}
}

func ValidateGenesis(genState GenesisState) error {

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
		whoIsRecords: []types.WhoIs{},
	}
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, /* TODO: Define what keepers the module needs */ genState GenesisState) []abci.ValidatorUpdate {
	// TODO: Define logic for when you would like to initalize a new genesis

	// Fetch & Iterate Through Names' whoIs
	for _, whoIs := range genState.whoIsRecords {
		// Assign whoIs Structures
		k.SetWhoIs(ctx, whoIs.Value, whoIs)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (GenesisState) {
	// TODO: Define logic for exporting state
	var names []types.WhoIs

	// Retrieve All The Names
	iterator := k.GetNamesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {

		// Get Key (Name)
		key := string(iterator.Key())

		// Get whoIs Of Name
		whois := k.GetWhoIs(ctx, key)

		// Append To Names List
		names = append(names, whois)
	}

	return GenesisState{whoIsRecords: names}
}
