package tax

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// keeper of the stake store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// Make Tax Keeper
func NewTaxKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

const (
	// default paramspace for params keeper
	DefaultParamspace = "tax_rate"
)

// make tax rate store key from denom
func RateStoreKey(Denom string) []byte {
	return append([]byte("tax_rate:"), []byte(Denom)...)
}

//______________________________________________________________________

// get the tax rate of deom
func (k Keeper) GetTaxRate(ctx sdk.Context, Denom string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(RateStoreKey(Denom))
	if b == nil {
		return sdk.ZeroDec()
	}

	var taxRate Rate
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, taxRate)
	v, err := sdk.NewDecFromStr(taxRate.Rate)
	if err == nil {
		return sdk.ZeroDec()
	}
	return v
}

// set the tax rate
func (k Keeper) SetTaxRate(ctx sdk.Context, taxRate Rate) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(taxRate)
	store.Set(RateStoreKey(taxRate.Denom), b)
}

// get all tax rates
func (k Keeper) GetTaxRates(ctx sdk.Context, Denom string) (taxRate Rate) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(RateStoreKey(Denom))
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &taxRate)
	return taxRate
}

// iterate all tax
func (k Keeper) IterateTaxRates(ctx sdk.Context, process func(Rate) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte("tax_rate:"))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		b := store.Get(val)
		rate := Rate{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rate)
		if process(rate) {
			return
		}
		iter.Next()
	}
}

