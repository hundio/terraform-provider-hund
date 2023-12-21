package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type unknownOnDependentFieldChange struct {
	dependentField path.Path
}

func UnknownOnDependentFieldChange(field path.Path) unknownOnDependentFieldChange {
	return unknownOnDependentFieldChange{dependentField: field}
}

var _ planmodifier.String = unknownOnDependentFieldChange{}
var _ planmodifier.Map = unknownOnDependentFieldChange{}

func (ud unknownOnDependentFieldChange) Description(ctx context.Context) string {
	return "Replace string values with unknown when the given field is changed in the plan."
}

func (ud unknownOnDependentFieldChange) MarkdownDescription(ctx context.Context) string {
	return "Replace string values with unknown when the given field is changed in the plan."
}

func (ud unknownOnDependentFieldChange) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var fieldState, fieldPlan attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, ud.dependentField, &fieldPlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, ud.dependentField, &fieldState)...)

	if !fieldPlan.IsUnknown() && !fieldState.Equal(fieldPlan) {
		tflog.Debug(ctx, "changing value to unknown", map[string]interface{}{
			"path":           req.Path.String(),
			"dependent":      ud.dependentField.String(),
			"dependentState": fieldState.String(),
			"dependentPlan":  fieldPlan.String(),
		})

		resp.PlanValue = types.StringUnknown()
	}
}

func (ud unknownOnDependentFieldChange) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	var fieldState, fieldPlan attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, ud.dependentField, &fieldPlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, ud.dependentField, &fieldState)...)

	if !fieldPlan.IsUnknown() && !fieldState.Equal(fieldPlan) {
		tflog.Debug(ctx, "changing value to unknown", map[string]interface{}{
			"path":           req.Path.String(),
			"dependent":      ud.dependentField.String(),
			"dependentState": fieldState.String(),
			"dependentPlan":  fieldPlan.String(),
		})

		resp.PlanValue = types.MapUnknown(req.PlanValue.ElementType(ctx))
	}
}
