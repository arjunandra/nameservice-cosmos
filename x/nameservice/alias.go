package nameservice

import (
	"github.com/arjunandra/namservice-cosmos/x/nameservice/internal/keeper"
	"github.com/arjunandra/namservice-cosmos/x/nameservice/internal/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QueryParams       = types.QueryParams
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgSetName 		= types.NewMsgSetName
	NewMsgBuyName 		= types.NewMsgBuyName
	NewMsgDeleteName 	= types.NewMsgDeleteName
	newWhoIs			= types.newWhoIs
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	// TODO: Fill out function aliases

	// variable aliases
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper       	= keeper.Keeper
	GenesisState 	= types.GenesisState
	Params       	= types.Params
	MsgSetName	 	= types.MsgSetName
	MsgBuyName 	 	= types.MsgBuyName
	MsgDeleteName	= types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames	= types.QueryResNames
	whoIs			= types.whoIs
)
