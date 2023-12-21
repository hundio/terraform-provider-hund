package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type GroupModel struct {
	Id                          types.String `tfsdk:"id"`
	CreatedAt                   types.String `tfsdk:"created_at"`
	UpdatedAt                   types.String `tfsdk:"updated_at"`
	Name                        types.String `tfsdk:"name"`
	NameTranslations            types.Map    `tfsdk:"name_translations"`
	Description                 types.String `tfsdk:"description"`
	DescriptionTranslations     types.Map    `tfsdk:"description_translations"`
	DescriptionHtml             types.String `tfsdk:"description_html"`
	DescriptionHtmlTranslations types.Map    `tfsdk:"description_html_translations"`
	Collapsed                   types.Bool   `tfsdk:"collapsed"`
	Position                    types.Int64  `tfsdk:"position"`
	Components                  types.List   `tfsdk:"components"`
}

func ToGroupModel(group hundApiV1.Group) (GroupModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := GroupModel{
		Id:        types.StringValue(group.Id),
		CreatedAt: types.StringValue(hundApiV1.ToStringTimestamp(group.CreatedAt)),
		UpdatedAt: types.StringValue(hundApiV1.ToStringTimestamp(group.UpdatedAt)),
		Collapsed: types.BoolValue(group.Collapsed),
		Position:  types.Int64Value(int64(group.Position)),
	}

	components, err := group.Components.AsGroupComponents1()
	if err != nil {
		diags.AddError(
			"Unable to parse Group Components",
			"Expected []string. Error: "+err.Error(),
		)
		return model, diags
	}

	modelComponents := []attr.Value{}
	for _, v := range components {
		modelComponents = append(modelComponents, types.StringValue(v))
	}

	var diag0 diag.Diagnostics
	model.Components, diag0 = types.ListValue(types.StringType, modelComponents)
	diags.Append(diag0...)

	name, nameMap, diag0 := hundApiV1.FromI18nString(group.Name)
	diags.Append(diag0...)

	if group.Description != nil {
		description, descriptionMap, diag0 := hundApiV1.FromI18nString(*group.Description)
		diags.Append(diag0...)

		if diags.HasError() {
			return model, diags
		}

		model.Description = *description
		model.DescriptionTranslations = *descriptionMap
	}

	descriptionHtml, descriptionHtmlMap, diag0 := hundApiV1.FromI18nString(group.DescriptionHtml)
	diags.Append(diag0...)

	if diags.HasError() {
		return model, diags
	}

	model.Name = *name
	model.NameTranslations = *nameMap

	model.DescriptionHtml = *descriptionHtml
	model.DescriptionHtmlTranslations = *descriptionHtmlMap

	return model, diags
}

type GroupComponentOrderingModel struct {
	Id    types.String `tfsdk:"id"`
	Group types.String `tfsdk:"group"`

	Components []types.String `tfsdk:"components"`
}

func ToGroupComponentOrderingModel(group hundApiV1.Group) (GroupComponentOrderingModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := GroupComponentOrderingModel{
		Id:    types.StringValue(group.Id),
		Group: types.StringValue(group.Id),
	}

	components, err := group.Components.AsGroupComponents1()
	if err != nil {
		diags.AddError(
			"Unable to parse Group Components",
			"Expected []string. Error: "+err.Error(),
		)
		return model, diags
	}

	for _, v := range components {
		model.Components = append(model.Components, types.StringValue(v))
	}

	return model, diags
}
