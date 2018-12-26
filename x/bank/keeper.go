package bank

import (
	"encoding/hex"
	"fmt"

	"github.com/PhenixChain/PhenixChain/codec"
	sdk "github.com/PhenixChain/PhenixChain/types"
	"github.com/PhenixChain/PhenixChain/x/auth"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

//-----------------------------------------------------------------------------
// Keeper

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper

	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error)
}

// BaseKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseKeeper struct {
	BaseSendKeeper

	ak auth.AccountKeeper

	tk TxKeeper
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(ak auth.AccountKeeper, tk TxKeeper) BaseKeeper {
	return BaseKeeper{
		BaseSendKeeper: NewBaseSendKeeper(ak, tk),
		ak:             ak,
		tk:             tk,
	}
}

// SetCoins sets the coins at the addr.
func (keeper BaseKeeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return setCoins(ctx, keeper.ak, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper BaseKeeper) SubtractCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return subtractCoins(ctx, keeper.ak, keeper.tk, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseKeeper) AddCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Tags, sdk.Error) {

	return addCoins(ctx, keeper.ak, keeper.tk, addr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(
	ctx sdk.Context, inputs []Input, outputs []Output,
) (sdk.Tags, sdk.Error) {

	return inputOutputCoins(ctx, keeper, inputs, outputs)
}

//-----------------------------------------------------------------------------
// Send Keeper

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// SendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak auth.AccountKeeper

	tk TxKeeper
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak auth.AccountKeeper, tk TxKeeper) BaseSendKeeper {
	return BaseSendKeeper{
		BaseViewKeeper: NewBaseViewKeeper(ak),
		ak:             ak,
		tk:             tk,
	}
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	return sendCoins(ctx, keeper, fromAddr, toAddr, amt)
}

//-----------------------------------------------------------------------------
// View Keeper

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak auth.AccountKeeper
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak auth.AccountKeeper) BaseViewKeeper {
	return BaseViewKeeper{
		ak: ak,
	}
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.ak, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.ak, addr, amt)
}

//-----------------------------------------------------------------------------

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func subtractCoins(ctx sdk.Context, am auth.AccountKeeper, tk TxKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	oldCoins := getCoins(ctx, am, addr)
	newCoins, hasNeg := oldCoins.SafeMinus(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}

	err := setCoins(ctx, am, addr, newCoins)
	// nikolas
	storeAddressTxHash(ctx, tk, addr)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, tk TxKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	// nikolas
	storeAddressTxHash(ctx, tk, addr)
	tags := sdk.NewTags("recipient", []byte(addr.String()))
	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
// NOTE: Make sure to revert state changes from tx on error
func sendCoins(ctx sdk.Context, bsk BaseSendKeeper, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	_, subTags, err := subtractCoins(ctx, bsk.ak, bsk.tk, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, bsk.ak, bsk.tk, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, bk BaseKeeper, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, bk.ak, bk.tk, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, bk.ak, bk.tk, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}

//#############################################################################################
//<!-- < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < << < < ☺
//v                  ✰  实现通过地址查找交易记录 ✰
//v
//v    保存address对应的txhash，每个地址最多保存最近的300个txhash
//v                                                        --- nikolas
//☺ > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > >  -->

func storeAddressTxHash(ctx sdk.Context, tk TxKeeper, addr sdk.AccAddress) {
	key := append([]byte{0x01}, addr.Bytes()...)
	store := ctx.KVStore(tk.key)

	// 取出历史数据
	getTxs := []Tx{}
	txHash := store.Get(key)
	if txHash != nil {
		if err := tk.cdc.UnmarshalJSON(txHash, &getTxs); err != nil {
			panic(err)
		}
	}

	hexStr := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	tx := []Tx{Tx{hexStr}}
	txs := []Tx{}

	// 新数据插在历史数据前面
	txs = append(tx, getTxs...)
	if len(txs) > 310 {
		// 数组长度最大暂定300
		txs = txs[:300]
	}
	ay, err := tk.cdc.MarshalJSON(txs)
	if err != nil {
		panic(err)
	}

	store.Set(key, ay)
}

type Tx struct {
	ID string `json:"tx"`
}

type TxKeeper struct {
	key sdk.StoreKey
	cdc *codec.Codec
}

func NewTxKeeper(cdc *codec.Codec, key sdk.StoreKey) TxKeeper {
	return TxKeeper{
		key: key,
		cdc: cdc,
	}
}

//#############################################################################################
