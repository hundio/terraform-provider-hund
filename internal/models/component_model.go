package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type ComponentModel struct {
	Id                          types.String   `tfsdk:"id"`
	CreatedAt                   types.String   `tfsdk:"created_at"`
	UpdatedAt                   types.String   `tfsdk:"updated_at"`
	Name                        types.String   `tfsdk:"name"`
	NameTranslations            types.Map      `tfsdk:"name_translations"`
	Description                 types.String   `tfsdk:"description"`
	DescriptionTranslations     types.Map      `tfsdk:"description_translations"`
	DescriptionHtml             types.String   `tfsdk:"description_html"`
	DescriptionHtmlTranslations types.Map      `tfsdk:"description_html_translations"`
	ExcludeFromGlobalHistory    types.Bool     `tfsdk:"exclude_from_global_history"`
	ExcludeFromGlobalUptime     types.Bool     `tfsdk:"exclude_from_global_uptime"`
	Group                       types.String   `tfsdk:"group"`
	LastEventAt                 types.String   `tfsdk:"last_event_at"`
	PercentUptime               types.Float64  `tfsdk:"percent_uptime"`
	Watchdog                    *WatchdogModel `tfsdk:"watchdog"`
}

func ToComponentModel(ctx context.Context, comp hundApiV1.ComponentExpansionary) (ComponentModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := ComponentModel{
		Id:                       types.StringValue(comp.Id),
		CreatedAt:                types.StringValue(hundApiV1.ToStringTimestamp(comp.CreatedAt)),
		UpdatedAt:                types.StringValue(hundApiV1.ToStringTimestamp(comp.UpdatedAt)),
		ExcludeFromGlobalHistory: types.BoolValue(comp.ExcludeFromGlobalHistory),
		ExcludeFromGlobalUptime:  types.BoolValue(comp.ExcludeFromGlobalUptime),
		LastEventAt:              types.StringPointerValue(hundApiV1.ToStringTimestampPtr(comp.LastEventAt)),
		PercentUptime:            types.Float64Value(float64(comp.PercentUptime)),
	}

	name, nameMap, diag0 := hundApiV1.FromI18nString(comp.Name)
	diags.Append(diag0...)

	if comp.Description != nil {
		description, descriptionMap, diag0 := hundApiV1.FromI18nString(*comp.Description)
		diags.Append(diag0...)

		if diags.HasError() {
			return model, diags
		}

		model.Description = *description
		model.DescriptionTranslations = *descriptionMap
	}

	descriptionHtml, descriptionHtmlMap, diag0 := hundApiV1.FromI18nString(comp.DescriptionHtml)
	diags.Append(diag0...)

	if diags.HasError() {
		return model, diags
	}

	model.Name = *name
	model.NameTranslations = *nameMap

	model.DescriptionHtml = *descriptionHtml
	model.DescriptionHtmlTranslations = *descriptionHtmlMap

	group, err := comp.Group.AsComponentExpansionaryGroup0()
	if err != nil {
		diags.AddError(
			"Could not parse Component",
			"Failed to parse Group as string: "+err.Error(),
		)
	}

	model.Group = types.StringValue(group)

	watchdog, err := comp.Watchdog.AsComponentExpansionaryWatchdog1()
	if err != nil {
		diags.AddError(
			"Could not parse Component",
			"Failed to parse Component Watchdog: "+err.Error(),
		)
	}

	watchdogModel, diag0 := ToWatchdogModel(watchdog)
	diags.Append(diag0...)

	model.Watchdog = &watchdogModel

	return model, diags
}
