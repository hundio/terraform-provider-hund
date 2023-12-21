package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
	"github.com/hundio/terraform-provider-hund/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IssueTemplateResource{}
var _ resource.ResourceWithConfigure = &IssueTemplateResource{}
var _ resource.ResourceWithImportState = &IssueTemplateResource{}
var _ resource.ResourceWithConfigValidators = &IssueTemplateResource{}
var _ resource.ResourceWithModifyPlan = &IssueTemplateResource{}

func NewIssueTemplateResource() resource.Resource {
	return &IssueTemplateResource{}
}

// IssueTemplateResource defines the resource implementation.
type IssueTemplateResource struct {
	client *hundApiV1.Client
}

// IssueTemplateResourceModel describes the resource data model.
type IssueTemplateResourceModel models.IssueTemplateModel

func (r *IssueTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_template"
}

func (r *IssueTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource representing templates for creating Issues, as well as Issue Updates.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("IssueTemplate"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "The \"kind\" of this IssueTemplate. This field can be either `issue` or `update`, depending on whether this IssueTemplate can be applied to an Issue or Issue Update, respectively.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(hundApiV1.ISSUETEMPLATEKINDIssue),
						string(hundApiV1.ISSUETEMPLATEKINDUpdate),
					),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "An internal name for identifying this IssueTemplate.",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: createdAtFieldMarkdownDescription("IssueTemplate"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: updatedAtFieldMarkdownDescription("IssueTemplate"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("When `kind` is `issue`, then the applied Issue will take on this title. This field supports [Liquid templating](https://shopify.github.io/liquid/)."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The body to use for an Issue/Update applied against this template. This field supports [Liquid templating](https://shopify.github.io/liquid/)."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("When `kind` is `issue`, then the applied Issue will take on this title. This field supports [Liquid templating](https://shopify.github.io/liquid/)."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"body_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The body to use for an Issue/Update applied against this template. This field supports [Liquid templating](https://shopify.github.io/liquid/)."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "The label to use for an Issue/Update applied against this template.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"variables": schema.MapNestedAttribute{
				MarkdownDescription: "An object defining a set of typed variables that can be provided in an application of this IssueTemplate. The variables can be accessed from any field in the IssueTemplate supporting Liquid.",
				Optional:            true,
				NestedObject:        issueTemplateVariablesSchema(),
			},
		},
	}
}

func (r *IssueTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hundApiV1.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hundApiV1.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IssueTemplateResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("body"),
			path.MatchRoot("body_translations"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("title"),
			path.MatchRoot("title_translations"),
		),
		validators.IssueTemplateKindFields(),
	}
}

func (r *IssueTemplateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state IssueTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Title.Equal(plan.Title) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("title_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.TitleTranslations.Equal(plan.TitleTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("title"), types.StringUnknown())...)
	}

	if !state.Body.Equal(plan.Body) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.BodyTranslations.Equal(plan.BodyTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body"), types.StringUnknown())...)
	}

	if !req.Plan.Raw.Equal(req.State.Raw) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated_at"), types.StringUnknown())...)
	}
}

func (r *IssueTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	title, err := hundApiV1.ToI18nStringPtr(data.Title, data.TitleTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	body, err := hundApiV1.ToI18nStringPtr(data.Body, data.BodyTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	form := hundApiV1.IssueTemplateFormCreate{
		Name:  data.Name.ValueString(),
		Kind:  hundApiV1.ISSUETEMPLATEKIND(data.Kind.ValueString()),
		Title: title,
		Body:  body,
		Label: (*hundApiV1.ISSUETEMPLATELABEL)(data.Label.ValueStringPointer()),
	}

	if len(data.Variables) > 0 {
		formVariables := models.ToIssueTemplateVariablesForm(data.Variables)

		form.Variables = &formVariables
	}

	rsp, err := r.client.CreateAIssueTemplate(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund IssueTemplate",
			err.Error(),
		)
		return
	}

	template, err := hundApiV1.ParseCreateAIssueTemplateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund IssueTemplate",
			err.Error(),
		)
		return
	}

	if template.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-201 status code: "+fmt.Sprint(template.StatusCode())+
				"\nError: "+string(template.Body),
		)
		return
	}

	newState, diag := models.ToIssueTemplateModel(*template.HALJSON201)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAIssueTemplate(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund IssueTemplate",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	template, err := hundApiV1.ParseRetrieveAIssueTemplateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund IssueTemplate",
			err.Error(),
		)
		return
	}

	if template.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(template.StatusCode())+
				"\nError: "+string(template.Body),
		)
		return
	}

	newState, diag := models.ToIssueTemplateModel(*template.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IssueTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	title, err := hundApiV1.ToI18nStringPtr(data.Title, data.TitleTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	body, err := hundApiV1.ToI18nStringPtr(data.Body, data.BodyTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	form := hundApiV1.IssueTemplateFormUpdate{
		Name:  data.Name.ValueStringPointer(),
		Title: title,
		Body:  body,
		Label: (*hundApiV1.ISSUETEMPLATELABEL)(data.Label.ValueStringPointer()),
	}

	if len(data.Variables) > 0 {
		formVariables := hundApiV1.IssueTemplateVariablesForm{}

		for k, itvm := range data.Variables {
			formVariables[k] = hundApiV1.IssueTemplateVariableForm{
				Type:     (*hundApiV1.IssueTemplateVariableFormType)(itvm.Type.ValueStringPointer()),
				Required: itvm.Required.ValueBoolPointer(),
			}
		}

		form.Variables = &formVariables
	}

	rsp, err := r.client.UpdateAIssueTemplate(ctx, data.Id.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Issue",
			err.Error(),
		)
		return
	}

	template, err := hundApiV1.ParseUpdateAIssueTemplateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue",
			err.Error(),
		)
		return
	}

	if template.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(template.StatusCode())+
				"\nError: "+string(template.Body),
		)
		return
	}

	newState, diag := models.ToIssueTemplateModel(*template.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.DeleteAIssueTemplate(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund Issue",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		template, err := hundApiV1.ParseDeleteAIssueTemplateResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(template.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *IssueTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
