package types

import (
	"fmt"
)

// GenesisState - all nameservice state that must be provided at genesis
type GenesisState struct {
	// TODO: Fill out what is needed by the module for genesis
	whoIsRecords []WhoIs	`json:"whois_records"`
}


// NewGenesisState creates a new GenesisState object
func NewGenesisState( /* TODO: Fill out with what is needed for genesis state */) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state

		whoIsRecords: nil,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state, these values will be initialized but empty
	
		whoIsRecords: []WhoIs{},
	}
}

// ValidateGenesis validates the nameservice genesis parameters
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