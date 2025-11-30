package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

// func IssueTemplateVariableAttrTypes() map[string]attr.Type {
// 	return map[string]attr.Type{
// 		"type":     types.StringType,
// 		"required": types.BoolType,
// 	}
// }

// IssueTemplateModel describes the data source data model.
type IssueTemplateModel struct {
	Id                types.String                          `tfsdk:"id"`
	CreatedAt         types.String                          `tfsdk:"created_at"`
	UpdatedAt         types.String                          `tfsdk:"updated_at"`
	Name              types.String                          `tfsdk:"name"`
	Title             types.String                          `tfsdk:"title"`
	Body              types.String                          `tfsdk:"body"`
	TitleTranslations types.Map                             `tfsdk:"title_translations"`
	BodyTranslations  types.Map                             `tfsdk:"body_translations"`
	Kind              types.String                          `tfsdk:"kind"`
	Label             types.String                          `tfsdk:"label"`
	Variables         map[string]IssueTemplateVariableModel `tfsdk:"variables"`
}

type IssueTemplateVariableModel struct {
	Type     types.String `tfsdk:"type"`
	Required types.Bool   `tfsdk:"required"`
}

func ToIssueTemplateVariableModel(variable hundApiV1.IssueTemplateVariable) IssueTemplateVariableModel {
	return IssueTemplateVariableModel{
		Type:     types.StringValue(string(variable.Type)),
		Required: types.BoolValue(variable.Required),
	}
}

func ToIssueTemplateVariables(variables hundApiV1.IssueTemplateVariables) map[string]IssueTemplateVariableModel {
	model := map[string]IssueTemplateVariableModel{}

	for name, variable := range variables {
		model[name] = ToIssueTemplateVariableModel(variable)
	}

	return model
}

func ToIssueTemplateModel(template hundApiV1.IssueTemplate) (IssueTemplateModel, diag.Diagnostics) {
	diag := diag.Diagnostics{}

	model := IssueTemplateModel{
		Id:        types.StringValue(template.Id),
		Name:      types.StringValue(template.Name),
		Kind:      types.StringValue(string(template.Kind)),
		CreatedAt: types.StringValue(hundApiV1.ToStringTimestamp(template.CreatedAt)),
		UpdatedAt: types.StringValue(hundApiV1.ToStringTimestamp(template.UpdatedAt)),
		Label:     types.StringPointerValue((*string)(template.Label)),
	}

	if template.Title != nil {
		title, titleMap, diag0 := hundApiV1.FromI18nString(*template.Title)
		diag.Append(diag0...)

		model.Title = *title
		model.TitleTranslations = *titleMap
	} else {
		model.TitleTranslations = types.MapNull(types.StringType)
	}

	if template.Body != nil {
		body, bodyMap, diag0 := hundApiV1.FromI18nString(*template.Body)
		diag.Append(diag0...)

		model.Body = *body
		model.BodyTranslations = *bodyMap

		if diag.HasError() {
			return model, diag
		}
	} else {
		model.BodyTranslations = types.MapNull(types.StringType)
	}

	model.Variables = ToIssueTemplateVariables(template.Variables)

	return model, diag
}

func ToIssueTemplateVariablesForm(model map[string]IssueTemplateVariableModel) hundApiV1.IssueTemplateVariablesForm {
	form := hundApiV1.IssueTemplateVariablesForm{}

	for k, itvm := range model {
		form[k] = hundApiV1.IssueTemplateVariableForm{
			Type:     (*hundApiV1.IssueTemplateVariableFormType)(itvm.Type.ValueStringPointer()),
			Required: itvm.Required.ValueBoolPointer(),
		}
	}

	return form
}
