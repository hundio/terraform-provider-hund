package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExactlyOneNonNullAttribute() exactlyOneNonNullAttribute {
	return exactlyOneNonNullAttribute{}
}

type exactlyOneNonNullAttribute struct{}

var _ validator.Object = &exactlyOneNonNullAttribute{}

func (v exactlyOneNonNullAttribute) Description(ctx context.Context) string {
	return "Validate object has a single non-null attribute."
}

func (v exactlyOneNonNullAttribute) MarkdownDescription(ctx context.Context) string {
	return "Validate object has a single non-null attribute."
}

func (v exactlyOneNonNullAttribute) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	nonNull := []string{}

	for name, val := range req.ConfigValue.Attributes() {
		if !val.IsNull() {
			nonNull = append(nonNull, name)
		}
	}

	nonNullCount := len(nonNull)

	if nonNullCount != 1 {
		var moreOrLess string

		if nonNullCount < 1 {
			moreOrLess = "Less"
		} else {
			moreOrLess = "More"
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Expected a single non-null Object field.",
			moreOrLess+" than one field in this Object is Non-null. There are "+
				fmt.Sprintf("%d: %v", len(nonNull), nonNull),
		)
	}
}
