package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type nullDefault struct{}

func NullDefault() nullDefault { return nullDefault{} }

var _ defaults.String = nullDefault{}
var _ planmodifier.String = nullDefault{}

func (nd nullDefault) Description(ctx context.Context) string {
	return "Null default value for Timestamps."
}

func (nd nullDefault) MarkdownDescription(ctx context.Context) string {
	return "Null default value for Timestamps."
}

func (nd nullDefault) DefaultString(ctx context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringNull()
}

func (nd nullDefault) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	if !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringNull()
		return
	}
}
