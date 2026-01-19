package myvalidator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Custom validator for regex patterns
type regexValidator struct{}

func (v regexValidator) Description(ctx context.Context) string {
	return "string must be a valid regular expression"
}

func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return "string must be a valid regular expression"
}

func (v regexValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, err := regexp.Compile(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Regular Expression",
			fmt.Sprintf("The string must be a valid regular expression: %s", err),
		)
	}
}

func ValidRegex() validator.String {
	return regexValidator{}
}
