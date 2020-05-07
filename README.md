# Cosmos SDK DAPP Development Flow

This guide will take us through the fundamentals required to build an end to end application using Cosmos SDK.

We will be building the tutorial application 'nameservice', a simple name bidding application where users can bid for and sell names.

I've framed this tutorial to be a **23 step** process, where each step involves creating/modifying a seperate .go file.

### Prerequisites:

>  - golang > 1.13.0 installed
>  - Scaffold Tool (git@github.com:cosmos/scaffold.git)

### Tree

This is the completed application tree that we will eventually be left with.

Each of the individual .go files are the leaf nodes in the tree that we will be working with each step.

```
./nameservice
├── Makefile
├── Makefile.ledger
├── app.go
├── cmd
│   ├── acli
│   │   └── main.go
│   └── aud
│       └── main.go
├── go.mod
├── go.sum
└── x
    └── nameservice
        ├── alias.go
        ├── client
        │   ├── cli
        │   │   ├── query.go
        │   │   └── tx.go
        │   └── rest
        │       ├── query.go
        │       ├── rest.go
        │       └── tx.go
        ├── genesis.go
        ├── handler.go
        ├── internal
        │   ├── keeper
        │   │   ├── keeper.go
        │   │   └── querier.go
        │   └── types
        │       ├── codec.go
        │       ├── errors.go
        │       ├── expected_keepers.go
        │       ├── key.go
        │       ├── msgs.go
        │       ├── querier.go
        │       └── types.go
        └── module.go


```
### Initialize Application
\
We initialize the application with the help of the scaffold tool through the command below.
```
scaffold app lvl-1 [user] [repo]
```
Then we need to create the modules for the application to access.

We do this in the /x/ folder with the scaffold tool, using the command
```
scaffold module [user] [repo] nameservice
```

## Step 1 - types.go

**Create** the types.go file in the ``` ./x/nameservice/internal/types/types.go ``` directory.

**Initialize:**

```
package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)
```

**Define Structures & Variables**
``` 
// Initial Name Value
var minNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type WhoIs struct {
	Value string			`json:"value"`
	Owner sdk.AccAddress 	`json:"owner"`
	Price sdk.Coins			`json:"price"`
}
```

**Define Constructors & Print Functions**
```
// whoIs Constructor
func NewWhoIs() WhoIs {
	return WhoIs {
		Price: minNamePrice,
	}
}

// whoIs Print Function
func (w WhoIs) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s\n Value: %s\n Price: %s`, w.Owner, w.Value, w.Price))
}
```

## Step 2 - key.go

**Define Module Name, Key Values, & Routes**
```
package types

const (
	// ModuleName is the name of the module
	ModuleName = "nameservice"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)
```

## Step 3 - errors.go

**Define Custom Errors Required For The Application**
```
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNameDoesNotExist = sdkerrors.Register(ModuleName, 1, "Name Doesn't Exist")
)
```

## Step 4 - expected_keepers.go

Create interfaces of what you expect the other keepers to have to be able to use this module.

```
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) 
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

## Step 5 - keeper.go

**Import Libraries & Declare Keeper Structure & Constructor**

```
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper	types.BankKeeper
	storeKey	sdk.StoreKey
	cdc 		*codec.Codec
}

// Keeper Constructor
func NewKeeper(coinkeeper types.BankKeeper, storekey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper {
		CoinKeeper: coinkeeper,
		storeKey: storekey,
		cdc: cdc,
	}
}
```

**Declare Getters, Setters, & Deletes For Structures (declared in types.go)**

```
// whoIs Getter & Setter
func (k Keeper) SetWhoIs(ctx sdk.Context, name string, w types.WhoIs) {     

	// No Owner
	if w.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(w))
}

func (k Keeper) GetWhoIs(ctx sdk.Context, name string) types.WhoIs {
	store := ctx.KVStore(k.storeKey)

	// No whoIs
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhoIs()
	}

	bz := store.Get([]byte(name))

	var whoIs types.WhoIs

	k.cdc.MustUnmarshalBinaryBare(bz, &whoIs)
	return whoIs
}

func (k Keeper) DeleteWhoIs(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}
```

**Declare Getters, Setters, & Iterators For Structures' Attributes (declared in types.go)**

```
// Name Getter & Setter & Bool & Iterator
func (k Keeper) GetName(ctx sdk.Context, name string) string {
	return k.GetWhoIs(ctx, name).Value
}

func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhoIs(ctx, name)
	whois.Value = value 
	k.SetWhoIs(ctx, name, whois)
}

func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// Owner Getter, Setter, & Bool

func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhoIs(ctx, name).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhoIs(ctx, name)
	whois.Owner = owner 
	k.SetWhoIs(ctx, name, whois)
}

func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhoIs(ctx, name).Owner.Empty()
}

// Price Getter & Setter

func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhoIs(ctx, name).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhoIs(ctx, name)
	whois.Price = price
	k.SetWhoIs(ctx, name, whois)
}

```

## Step 6 - msg.go

**Message Structure Declarations*
```
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Structure Declarations

type MsgSetName struct {
	Name string				`json:"name"`
	Value string			`json:"value"`	
	Owner sdk.AccAddress	`json:"owner"`	
}

type MsgBuyName struct {
	Name string				`json:"name"`
	Bid sdk.Coins			`json:"bid"`
	Buyer sdk.AccAddress	`json:"buyer"`
}

type MsgDeleteName struct {
	Name string				`json:"name"`
	Owner sdk.AccAddress	`json:"owner"`
}
```

**Message Constructors**
```
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName {
		Name: name,
		Value: value,
		Owner: owner,
	}
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName {
		Name: name,
		Bid: bid,
		Buyer: buyer,
	}
}

func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName {
		Name: name,
		Owner: owner,
	}
}
```

**Message Route Declarations**
```
func (msg MsgSetName) Route() string { return RouterKey }
func (msg MsgBuyName) Route() string { return RouterKey }
func (msg MsgDeleteName) Route() string { return RouterKey }
```

**Message Type Declarations**
```
func (msg MsgSetName) Type() string { return "set_name" }
func (msg MsgBuyName) Type() string {return "buy_name"}
func (msg MsgDeleteName) Type() string { return "delete_name" }
```

**Stateless Checks**
```
func (msg MsgSetName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg.Owner.String())
	}

	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Value cannot be empty")
	}

	return nil
}

func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	if !msg.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}

	return nil
}

func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	return nil
}
```

**Message Sign Bytes Getter**
```
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) 
}
```

**Message Signers Getter**
```
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
```

## Step 7 - handler.go

**Initialization & Handler Constructor**

```
package nameservice

import (
	"fmt"

	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the nameservice type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, k, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, k, msg)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
```

**Declare Each Individual Handle Function**

```
// Handler Functions

func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) (*sdk.Result, error){
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.SetName(ctx, msg.Name, msg.Value) 
	return &sdk.Result{}, nil
}

func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) (*sdk.Result, error) {

	// Check If Current Price > Bid
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid Didn't Surpass Current Price")
	}

	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		
		// Error Occurred
		if err != nil {
			return nil, err
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)

		// Error Occured
		if err != nil {
			return nil, err
		}
	}

	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg MsgDeleteName) (*sdk.Result, error){
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}

	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.DeleteWhoIs(ctx, msg.Name)
	return &sdk.Result{}, nil
}
```

## Step 7 - /internal/types/querier.go

**Declare Querier Structures**

```
type QueryResResolve struct {
	Value string `json:"value"`
}

type QueryResNames []string
```

**Implement fmt.Stringer For Structures**

```
func (r QueryResResolve) String() string {
	return r.Value
}

func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}
```

## Step 8 - /internal/keeper/querier.go

**Add Query End-Points & Routes In NewQuerier**

```
// Query End-Points
const (
	QueryResolve = "resolve"
	QueryWhoIs = "whois"
	QueryNames = "names"
)

// NewQuerier creates a new querier for naeservice clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, k)
		case QueryWhoIs:
			return queryWhoIs(ctx, path[1:], req, k)
		case QueryNames:
			return queryNames(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}
```

**Define Input Parameters & Responses For Each Query**

```
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	value := keeper.GetName(ctx, path[0])

	if len(value) == 0 {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Couldn't Resolve Name")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResResolve{Value: value})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryWhoIs(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	whois := keeper.GetWhoIs(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}	

func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var namesList types.QueryResNames

	iterator := keeper.GetNamesIterator(ctx)

	for; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
```

## Step 9 - alias.go

**Declare Constants Exposed From Internal Package**

```
const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
)
```

**Declare Function Aliases**

```
var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgSetName 		= types.NewMsgSetName
	NewMsgBuyName 		= types.NewMsgBuyName
	NewMsgDeleteName 	= types.NewMsgDeleteName
	NewWhoIs			= types.NewWhoIs
	RegisterCodec       = types.RegisterCodec
)
```

**Declare Variable Aliases**

```
var (
	ModuleCdc     = types.ModuleCdc
)
```

**Declare Required Structures (from keeper & types)**

```
type (
	Keeper       	= keeper.Keeper
	Params       	= types.Params
	MsgSetName	 	= types.MsgSetName
	MsgBuyName 	 	= types.MsgBuyName
	MsgDeleteName	= types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames	= types.QueryResNames
	whoIs			= types.WhoIs
)
```

## Step 10 - codec.go

**Register Module Messages**

```
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
	cdc.RegisterConcrete(MsgDeleteName{}, "nameservice/DeleteName", nil)
}
```

## Step 11 - /client/cli/query.go

**Add Query Commands To ```GetQueryCmd()```**

```
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
```
**Define ```cobra.Command```s For Each Module's Added Querier Command**

```
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
```

## Step 12 - /client/cli/tx.go

**Add Transaction Commands To ```GetTxCmd()```**

```
// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	nameserviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceTxCmd.AddCommand(flags.PostCommands(

		// Added Transaction Commands
		
		GetCmdBuyName(cdc),
		GetCmdSetName(cdc),
		GetCmdDeleteName(cdc),
	)...)

	return nameserviceTxCmd
}
```

**Define ```cobra.Command```s For Each Module's Added Transaction Command**

```
func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-name [name] [amount]",
		Short: "Buy New name / Bid For Existing Name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBuyName(args[0], coins, cliCtx.GetFromAddress())

			// State-less Checks
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdSetName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-name [name] [value]",
		Short: "Set Value For Your Name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// if err := cliCtx.EnsureAccountExists(); err != nil {
			// 	return err
			// }

			msg := types.NewMsgSetName(args[0], args[1], cliCtx.GetFromAddress())

			// State-less Checks
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdDeleteName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-name [name]",
		Short: "Delete The Name That You Own (along with it's associated fields)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDeleteName(args[0], cliCtx.GetFromAddress())

			// State-less Checks
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
```

## Step 13 - /client/rest/rest.go

**Define REST Client Interface (Declare All Handler Routes From query.go & tx.go)**

```
// RegisterRoutes registers nameservice-related REST handlers to a router

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	// registerQueryRoutes(cliCtx, r)
	// registerTxRoutes(cliCtx, r)
	
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), namesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoIsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), deleteNameHandler(cliCtx)).Methods("DELETE")
}
```

## Step 14 - /client/rest/query.go

Create the file if it doesn't already exist.

**Define Handlers For Query Commands**

The query commands were previously defined in /client/cli/query.go.

```
func resolveNameHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func whoIsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func namesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
```

## Step 15 - /client/rest/tx.go

Create the file if it doesn't already exist.

**Declare Request Structures Of Transaction Commands**

The transaction commands were previously defined in /client/cli/query.go.

```
type buyNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Amount  string       `json:"amount"`
	Buyer   string       `json:"buyer"`
}

type setNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Value   string       `json:"value"`
	Owner   string       `json:"owner"`
}

type deleteNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Owner   string       `json:"owner"`
}
```

**Define Handlers For Transaction Commands**

```
func buyNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyNameReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed To Parse Request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		// State-less Checks
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Retrieve Account
		addr, err := sdk.AccAddressFromBech32(req.Buyer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Send Coins
		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgBuyName(req.Name, coins, addr)

		// State-less Checks
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func setNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setNameReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Retrieve Address
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgSetName(req.Name, req.Value, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func deleteNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deleteNameReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}	

		// Retrieve Address
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Create Message
		msg := types.NewMsgDeleteName(req.Name, addr)
		err = msg.ValidateBasic()
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate Response
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
```

## Step 16 - module.go

**Add Keepers That The Module Depends On In AppModule Structure & NewAppModule()**

```
// AppModule implements an application module for the nameservice module.
type AppModule struct {
	AppModuleBasic

	keeper        Keeper
	
	// Added Keepers
	
	bankKeeper bank.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k Keeper, /*TODO: Add Keepers that your application depends on*/ bankKeeper bank.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		keeper:              k,
		
		// Added Keepers

		bankKeeper: bankKeeper,
	}
}
```

## Step 17 - x/[name]/genesis.go

**Define ```GenesisState```'s Structure, Constructor, & Default State**

```
type GenesisState struct {
	whoIsRecords []types.WhoIs `json:"whois_records"`
}

// Constructor 

func NewGenesisState(whoIsRecords []types.WhoIs) GenesisState {
	return GenesisState{whoIsRecords: whoIsRecords}
}

// Default 

func DefaultGenesisState() GenesisState {
	return GenesisState{
		whoIsRecords: []types.WhoIs{},
	}
}
```

**Define State-less Checks In ```ValidateGenesis()```**

```
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
```

**Define Logic For Genesis Initialization**

```
// InitGenesis initialize default parameters
// and the keeper's address to pubkey map

func InitGenesis(ctx sdk.Context, k Keeper, /* TODO: Define what keepers the module needs */ genState GenesisState) []abci.ValidatorUpdate {
	
	// Fetch & Iterate Through Names' whoIs

	for _, whoIs := range genState.whoIsRecords {
		// Assign whoIs Structures
		k.SetWhoIs(ctx, whoIs.Value, whoIs)
	}
	return []abci.ValidatorUpdate{}
}
```

**Define Logic For Exporting State**

```
// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis

func ExportGenesis(ctx sdk.Context, k Keeper) (GenesisState) {
	
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
```

## Step 18 - /internal/types/genesis.go

**Copy ```GenesisState```'s Structure, Constructor, Default State, & State-less Validation From ```x/[name]/genesis.go```**

```
// GenesisState - all nameservice state that must be provided at genesis
type GenesisState struct {
	whoIsRecords []WhoIs	`json:"whois_records"`
}


// NewGenesisState creates a new GenesisState object
func NewGenesisState( /* TODO: Fill out with what is needed for genesis state */) GenesisState {
	return GenesisState{
		whoIsRecords: nil,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
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
```

## Step 19 - app.go

**Import Newly Defined Modules**

```
import "github.com/[username]/[repo_name]/x/[repo_name]"
```

**Add Defined Modules To ```ModuleBasics```**

```
// ModuleBasics The module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.

ModuleBasics = module.NewBasicManager(
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	distr.AppModuleBasic{},
	params.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	
	// Added Modules

	nameservice.AppModule{},
)
```

**Add Defined Keeper To ```NewApp```'s Structure**

```
// NewApp extended ABCI application
type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper

	// Added Keeper 
	nsKeeper       nameservice.Keeper
}
```

**Initialize Added Keepers In ```NewInitApp()```**

```
app.nsKeeper = nameservice.NewKeeper(
		app.bankKeeper,
		keys[nameservice.StoreKey],
		app.cdc,
	)
```

**Include Added Modules in ```app.mm```'s Definition & ```.SetOrderInitGenesis()```**

```
// NOTE: Any module instantiated in the module manager that is later modified
// must be passed by reference here.

app.mm = module.NewManager(
	genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
	auth.NewAppModule(app.accountKeeper),
	bank.NewAppModule(app.bankKeeper, app.accountKeeper),
	supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
	distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
	slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
	staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
	slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),

	// Added Module
	nameservice.NewAppModule(app.nsKeeper, app.bankKeeper),
)
```

```
// Sets the order of Genesis - Order matters, genutil is to always come last
// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.

app.mm.SetOrderInitGenesis(
	genaccounts.ModuleName,
	distr.ModuleName,
	staking.ModuleName,
	auth.ModuleName,
	bank.ModuleName,
	slashing.ModuleName,

	// Added Module
	nameservice.ModuleName,
	
	genutil.ModuleName,
)
```

**Define ```ExportAppStateAndValidators()``` To Help Bootstrap The Initial State For The Application**

```
func (app *NewApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	
		// as if they could withdraw from the start of the next block
		ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	
		genState := app.mm.ExportGenesis(ctx)

		appState, err = codec.MarshalJSONIndent(app.cdc, genState)
		if err != nil {
			return nil, nil, err
		}
	
		validators = staking.WriteValidators(ctx, app.stakingKeeper)
	
		return appState, validators, nil
}
```

## Step 20 - /cmd/aud/main.go

**Import App**

```
import app "github.com/[username]/[repo_name]/app"
```

## Step 21 - /cmd/acli/main.go

**Import App**

```
import app "github.com/[username]/[repo_name]/app"
```

## Step 22 - go.mod & Makefile

**Add Module Path**

```
module github.com/[username]/[repo_name]
```

**Create A New File ```Makefile.ledger``` To Include Ledger Nano S Support**

```
LEDGER_ENABLED ?= true

build_tags =
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif
```

**Include ```Makefile.ledger``` In The Makefile**

```
PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NameService \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=aud \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=acli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

# Include Ledger File
include Makefile.ledger

all: install

install: go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/aud
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/acli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
```

## Step 23 - Build The Application!

```
# Install the app into your $GOBIN
make install
```

```
# Now you should be able to run the following commands:
aud help
acli help
```

If aud & acli commands cannot be found, refresh ```.bashrc``` file and try again:

```
source ~/.bashrc
```


### Created By Arjun Andra


