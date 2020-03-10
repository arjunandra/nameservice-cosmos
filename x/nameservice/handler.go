package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the nameservice type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handler Functions

func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) (*sdk.Result, error){
	if !msg.Owner.Equals(keeper.getOwner()) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.setName(ctx, msg.Name, msg.Value) 
	return &sdk.Result{}, nil
}

func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) (*sdk.Result, error) {

	// Check If Current Price > Bid
	if keeper.getPrice(ctx, msg.Name).isAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid Didn't Surpass Current Price")
	}

	if keeper.hasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.getOwner(ctx, msg.Name), msg.Bid)
		
		// Error Occurred
		if err != nil {
			return nil, err
		}
	} else {
		err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)

		// Error Occured
		if err != nil {
			return nil, err
		}
	}

	keeper.setOwner(ctx, msg.Name, msg.Buyer)
	keeper.setPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

func handleMsgDeleteName(ctx sdk.Context, keeper Keeper) (*sdk.Result, error){
	if !keeper.isNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}

	if !msg.Owner.Equals(keeper.getOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.deleteWhoIs(ctx, msg.Name)
	return &sdk.Result{}, nil
}

