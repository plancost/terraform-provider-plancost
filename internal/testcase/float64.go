package testcase

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ knownvalue.Check = float64Exact{}

type float64Exact struct {
	value float64
}

// CheckValue determines whether the passed value is of type float64, and
// contains a matching float64 value.
func (v float64Exact) CheckValue(other any) error {
	num, ok := other.(float64)

	if !ok {
		return fmt.Errorf("expected json.Number value for Float64Exact check, got: %T", other)
	}

	if num != v.value {
		return fmt.Errorf("expected value %s for Float64Exact check, got: %s", v.String(), strconv.FormatFloat(num, 'f', -1, 64))
	}

	return nil
}

// String returns the string representation of the float64 value.
func (v float64Exact) String() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

// Float64Exact returns a Check for asserting equality between the
// supplied float64 and the value passed to the CheckValue method.
func Float64Exact(value float64) float64Exact {
	return float64Exact{
		value: value,
	}
}

var _ knownvalue.Check = float64Approx{}

type float64Approx struct {
	value     float64
	tolerance float64
}

// CheckValue determines whether the passed value is of type float64, and
// contains a matching float64 value within the tolerance.
func (v float64Approx) CheckValue(other any) error {
	num, ok := other.(float64)

	if !ok {
		return fmt.Errorf("expected float64 value for Float64Approx check, got: %T", other)
	}

	diff := math.Abs(num - v.value)
	if diff > v.tolerance {
		return fmt.Errorf("expected value %s (tolerance %s) for Float64Approx check, got: %s",
			strconv.FormatFloat(v.value, 'f', -1, 64),
			strconv.FormatFloat(v.tolerance, 'f', -1, 64),
			strconv.FormatFloat(num, 'f', -1, 64))
	}

	return nil
}

// String returns the string representation of the float64 value.
func (v float64Approx) String() string {
	return fmt.Sprintf("%s (tolerance %s)", strconv.FormatFloat(v.value, 'f', -1, 64), strconv.FormatFloat(v.tolerance, 'f', -1, 64))
}

// Float64Approx returns a Check for asserting equality between the
// supplied float64 and the value passed to the CheckValue method within a tolerance.
func Float64Approx(value float64, tolerance float64) float64Approx {
	return float64Approx{
		value:     value,
		tolerance: tolerance,
	}
}
