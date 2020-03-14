package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group nameservice queries under a subcommand
	nameserviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceQueryCmd.AddCommand(
		flags.GetCommands(
			
			// Added Query Commands

			GetCmdGetName(queryRoute, cdc),
			GetCmdWhoIs(queryRoute, cdc),
			GetCmdNames(queryRoute, cdc),
		)...,
	)

	return nameserviceQueryCmd
}

// Define cobra.Commands For Each Module's Added Querier Command

func GetCmdGetName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use: "get [name]",
		Short: "get name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, name), nil)
			
			if err != nil {
				fmt.Sprintf("Couldn't get name - %s \n", name)
				return nil
			}

			var output types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &output)
			return cliCtx.PrintOutput(output)
		},
	}
}

func GetCmdWhoIs(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use: "whois [name]",
		Short: "Query whois info of name",
		Args: cobra.ExactArgs(1),
		RunE: func (cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, name), nil)
		
			if err != nil {
				fmt.Sprintf("Couldn't Resolve whoIs - %s \n", name)
				return nil
			}

			var output types.WhoIs
			cdc.MustUnmarshalJSON(res, &output)
			return cliCtx.PrintOutput(output)
		},
	}
}

func GetCmdNames(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use: "names",
		Short: "names",
		Args: cobra.ExactArgs(1),
		RunE: func (cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names", queryRoute), nil)
		
			if err != nil {
				fmt.Sprintf("Couldn't Get Query Names \n")
				return nil
			}

			var output types.QueryResNames
			cdc.MustUnmarshalJSON(res, &output)
			return cliCtx.PrintOutput(output)
		},
	}
}