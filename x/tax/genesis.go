package tax

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Rate - store tax rate of one coin
type Rate struct {
	Denom  string  	`json:"denom"`
	Rate   string 	`json:"rate"`
}

// Rates - store tax rates of all coins
type Rates []Rate

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Rates Rates `json:"rates"` // taxrate object
}

func NewGenesisState(taxRates Rates) GenesisState {
	return GenesisState{
		Rates: taxRates,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Rates: Rates{
			Rate{"terra", "0.001"},
		},
	}
}

// new tax genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, taxRate := range data.Rates {
		keeper.SetTaxRate(ctx, taxRate)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var taxRates Rates
	keeper.IterateTaxRates(ctx, func(taxRate Rate) (stop bool) {
		taxRates = append(taxRates, taxRate)
		return false
	})

	return NewGenesisState(taxRates)
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	return nil
}
