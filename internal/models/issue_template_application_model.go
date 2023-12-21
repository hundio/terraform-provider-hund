package models

import (
	"context"
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

type IssueTemplateApplicationUpdateModel struct {
	Id               types.String                           `tfsdk:"id"`
	IssueTemplateId  types.String                           `tfsdk:"issue_template_id"`
	Body             types.String                           `tfsdk:"body"`
	BodyTranslations types.Map                              `tfsdk:"body_translations"`
	Label            types.String                           `tfsdk:"label"`
	Schema           types.Map                              `tfsdk:"schema"`
	Variables        IssueTemplateVariablesApplicationModel `tfsdk:"variables"`
}

type IssueTemplateApplicationIssueModel struct {
	Id               types.String                           `tfsdk:"id"`
	IssueTemplateId  types.String                           `tfsdk:"issue_template_id"`
	Body             types.String                           `tfsdk:"body"`
	BodyTranslations types.Map                              `tfsdk:"body_translations"`
	Label            types.String                           `tfsdk:"label"`
	Schema           types.Map                              `tfsdk:"schema"`
	Variables        IssueTemplateVariablesApplicationModel `tfsdk:"variables"`

	Title             types.String `tfsdk:"title"`
	TitleTranslations types.Map    `tfsdk:"title_translations"`
}

type IssueTemplateVariableApplicationModel struct {
	String     types.String `tfsdk:"string"`
	I18nString types.Map    `tfsdk:"i18n_string"`
	Datetime   types.String `tfsdk:"datetime"`
	Number     types.Number `tfsdk:"number"`
}

type IssueTemplateVariablesApplicationModel map[string]IssueTemplateVariableApplicationModel

func (m IssueTemplateVariablesApplicationModel) ApiValue() (hundApiV1.IssueTemplateVariablesApplication, diag.Diagnostics) {
	var diags diag.Diagnostics

	variables := hundApiV1.IssueTemplateVariablesApplication{}

	for k, variable := range m {
		apiVar, err := variable.ApiValue()
		if err != nil {
			diags.AddError(
				"Template Variable Application conversion error",
				"Got error encoding IssueTemplateVariablesApplication: "+err.Error(),
			)
			return variables, diags
		}

		variables[k] = *apiVar
	}

	return variables, diags
}

func (m IssueTemplateVariableApplicationModel) ApiValue() (*hundApiV1.IssueTemplateVariableApplication, error) {
	var opaque hundApiV1.IssueTemplateVariableApplication

	if !m.String.IsNull() {
		err := opaque.FromIssueTemplateVariableApplication2(m.String.ValueString())
		return &opaque, err
	}

	if !m.I18nString.IsNull() {
		i18n := map[string]string{}

		for k, v := range m.I18nString.Elements() {
			i18n[k] = v.(types.String).ValueString()
		}

		err := opaque.FromIssueTemplateVariableApplication3(i18n)

		return &opaque, err
	}

	if !m.Datetime.IsNull() {
		ts, err := hundApiV1.ToIntTimestamp(m.Datetime.ValueString())
		if err != nil {
			return nil, err
		}

		err = opaque.FromIssueTemplateVariableApplication1(ts)

		return &opaque, err
	}

	if !m.Number.IsNull() {
		num, _ := m.Number.ValueBigFloat().Float64()

		err := opaque.FromIssueTemplateVariableApplication0(num)

		return &opaque, err
	}

	return nil, errors.New("could not find non-null field")
}

func ToIssueTemplateVariableApplicationModel(variable hundApiV1.IssueTemplateVariableApplication, schema IssueTemplateVariableModel) (IssueTemplateVariableApplicationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := IssueTemplateVariableApplicationModel{
		I18nString: types.MapNull(types.StringType),
	}

	switch schema.Type.ValueString() {
	case "string":
		str, err := variable.AsIssueTemplateVariableApplication2()

		if err != nil {
			diags.AddError(
				"Could not parse IssueTemplate Variable Application",
				"Failed to parse variable as string",
			)
			return model, diags
		}

		model.String = types.StringValue(str)
	case "datetime":
		num, err := variable.AsIssueTemplateVariableApplication1()

		if err != nil {
			diags.AddError(
				"Could not parse IssueTemplate Variable Application",
				"Failed to parse variable as datetime",
			)
			return model, diags
		}

		model.Datetime = types.StringValue(hundApiV1.ToStringTimestamp(num))
	case "i18n-string":
		i18n, err := variable.AsIssueTemplateVariableApplication3()

		if err != nil {
			diags.AddError(
				"Could not parse IssueTemplate Variable Application",
				"Failed to parse variable as i18n-string",
			)
			return model, diags
		}

		i18nModel := map[string]attr.Value{}
		for k, v := range i18n {
			i18nModel[k] = types.StringValue(v)
		}

		i18nMap, diag0 := types.MapValue(types.StringType, i18nModel)
		diags.Append(diag0...)

		model.I18nString = i18nMap
	case "number":
		num, err := variable.AsIssueTemplateVariableApplication0()

		if err != nil {
			diags.AddError(
				"Could not parse IssueTemplate Variable Application",
				"Failed to parse variable as number",
			)
			return model, diags
		}

		model.Number = types.NumberValue(big.NewFloat(num))
	default:
		diags.AddError(
			"Could not parse IssueTemplate Variable Application",
			"Unsupported variable schema type: "+schema.Type.ValueString(),
		)
	}

	return model, diags
}

func ToIssueTemplateVariablesApplicationModel(variables hundApiV1.IssueTemplateVariablesApplication, schema map[string]IssueTemplateVariableModel) (map[string]IssueTemplateVariableApplicationModel, diag.Diagnostics) {
	model := map[string]IssueTemplateVariableApplicationModel{}
	diags := diag.Diagnostics{}

	for name, variable := range variables {
		s, ok := schema[name]
		if !ok {
			diags.AddError(
				"Could not parse IssueTemplate Variable Application",
				"Variable `"+name+"`: missing schema",
			)

			continue
		}

		var diag0 diag.Diagnostics
		model[name], diag0 = ToIssueTemplateVariableApplicationModel(variable, s)
		diags.Append(diag0...)
	}

	return model, diags
}

func ToIssueTemplateApplicationUpdateModel(ctx context.Context, app hundApiV1.IssueTemplateApplication) (*IssueTemplateApplicationUpdateModel, diag.Diagnostics) {
	model := IssueTemplateApplicationUpdateModel{
		Id:              types.StringValue(app.Id),
		IssueTemplateId: types.StringValue(app.IssueTemplate),
		Label:           types.StringPointerValue((*string)(app.Label)),
	}

	if app.Body != nil {
		body, bodyTranslations, diag := hundApiV1.FromI18nString(*app.Body)
		if diag.HasError() {
			return nil, diag
		}

		model.Body = *body
		model.BodyTranslations = *bodyTranslations
	}

	var diags diag.Diagnostics

	schemaModel := ToIssueTemplateVariables(app.Schema)

	model.Schema, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"type":     types.StringType,
		"required": types.BoolType,
	}}, schemaModel)

	if diags.HasError() {
		return &model, diags
	}

	model.Variables, diags = ToIssueTemplateVariablesApplicationModel(app.Variables, schemaModel)

	return &model, diags
}

func ToIssueTemplateApplicationIssueModel(ctx context.Context, app hundApiV1.IssueTemplateApplicationIssue) (*IssueTemplateApplicationIssueModel, diag.Diagnostics) {
	model := IssueTemplateApplicationIssueModel{
		Id:              types.StringValue(app.Id),
		IssueTemplateId: types.StringValue(app.IssueTemplate),
		Label:           types.StringPointerValue((*string)(app.Label)),
	}

	if app.Title != nil {
		title, titleTranslations, diag := hundApiV1.FromI18nString(*app.Title)
		if diag.HasError() {
			diag.Append(diag...)
			return nil, diag
		}

		model.Title = *title
		model.TitleTranslations = *titleTranslations
	}

	if app.Body != nil {
		body, bodyTranslations, diag := hundApiV1.FromI18nString(*app.Body)
		if diag.HasError() {
			return nil, diag
		}

		model.Body = *body
		model.BodyTranslations = *bodyTranslations
	}

	var diags diag.Diagnostics

	schemaModel := ToIssueTemplateVariables(app.Schema)

	model.Schema, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"type":     types.StringType,
		"required": types.BoolType,
	}}, schemaModel)

	if diags.HasError() {
		return &model, diags
	}

	model.Variables, diags = ToIssueTemplateVariablesApplicationModel(app.Variables, schemaModel)

	return &model, diags
}
