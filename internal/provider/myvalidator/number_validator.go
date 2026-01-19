package myvalidator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Custom validator for number range
type numberRangeValidator struct {
	min *big.Float
	max *big.Float
}

func (v numberRangeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("number must be between %s and %s", v.min.String(), v.max.String())
}

func (v numberRangeValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("number must be between %s and %s", v.min.String(), v.max.String())
}

func (v numberRangeValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	val := req.ConfigValue.ValueBigFloat()
	if val.Cmp(v.min) < 0 || val.Cmp(v.max) > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Number Range",
			fmt.Sprintf("The number must be between %s and %s", v.min.String(), v.max.String()),
		)
	}
}

func NumberBetween(minValue, maxValue float64) validator.Number {
	return numberRangeValidator{
		min: big.NewFloat(minValue),
		max: big.NewFloat(maxValue),
	}
}
