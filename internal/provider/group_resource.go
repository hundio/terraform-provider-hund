package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
	"github.com/hundio/terraform-provider-hund/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupResource{}
var _ resource.ResourceWithConfigure = &GroupResource{}
var _ resource.ResourceWithImportState = &GroupResource{}
var _ resource.ResourceWithConfigValidators = &GroupResource{}
var _ resource.ResourceWithModifyPlan = &GroupResource{}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

// GroupResource defines the resource implementation.
type GroupResource struct {
	client *hundApiV1.Client
}

// GroupResourceModel describes the resource data model.
type GroupResourceModel models.GroupModel

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource representing a logical grouping of Components.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("Group"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: createdAtFieldMarkdownDescription("Group"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: updatedAtFieldMarkdownDescription("Group"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The name of this Group."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The name of this Group."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("A description of this Group, potentially with markdown formatting."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("A description of this Group, potentially with markdown formatting."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"description_html": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("An HTML rendering of the markdown-formatted `description`."),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("description")),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("description_translations")),
				},
			},
			"description_html_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("An HTML rendering of the markdown-formatted `description`."),
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("description")),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("description_translations")),
				},
			},
			"collapsed": schema.BoolAttribute{
				MarkdownDescription: "Whether or not this group is displayed collapsed by default on the status page.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"position": schema.Int64Attribute{
				MarkdownDescription: "An integer representing the position of this Group. Groups are displayed on the status page in ascending order according to this value.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"components": schema.ListAttribute{
				MarkdownDescription: "A list of the Component IDs contained within this Group.",
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *GroupResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("name"),
			path.MatchRoot("name_translations"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("description"),
			path.MatchRoot("description_translations"),
		),
	}
}

func (r *GroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Name.Equal(plan.Name) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("name_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.NameTranslations.Equal(plan.NameTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("name"), types.StringUnknown())...)
	}

	if !state.Description.Equal(plan.Description) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("description_translations"), types.MapUnknown(types.StringType))...)
	}

	if !state.DescriptionTranslations.Equal(plan.DescriptionTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("description"), types.StringUnknown())...)
	}

	if !req.Plan.Raw.Equal(req.State.Raw) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated_at"), types.StringUnknown())...)
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name, err := hundApiV1.ToI18nString(data.Name, data.NameTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	description, err := hundApiV1.ToI18nStringPtr(data.Description, data.DescriptionTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	form := hundApiV1.GroupFormCreate{
		Description: &description,
		Name:        name,
		Collapsed:   data.Collapsed.ValueBoolPointer(),
		Position:    hundApiV1.ToIntPtr(data.Position.ValueInt64Pointer()),
	}

	rsp, err := r.client.CreateAGroup(ctx, form, hundApiV1.Unexpand("components"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Group",
			err.Error(),
		)
		return
	}

	group, err := hundApiV1.ParseCreateAGroupResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Group",
			err.Error(),
		)
		return
	}

	if group.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-201 status code: "+fmt.Sprint(group.StatusCode())+
				"\nError: "+string(group.Body),
		)
		return
	}

	newState, diag := models.ToGroupModel(*group.HALJSON201)
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

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAGroup(ctx, data.Id.ValueString(), hundApiV1.Unexpand("components"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Group",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	group, err := hundApiV1.ParseRetrieveAGroupResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Group",
			err.Error(),
		)
		return
	}

	if group.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(group.StatusCode())+
				"\nError: "+string(group.Body),
		)
		return
	}

	newState, diag := models.ToGroupModel(*group.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name, err := hundApiV1.ToI18nString(data.Name, data.NameTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	description, err := hundApiV1.ToI18nStringPtr(data.Description, data.DescriptionTranslations)
	if err != nil {
		resp.Diagnostics.Append(models.I18nStringError(err))
		return
	}

	form := hundApiV1.GroupFormUpdate{
		Name:        &name,
		Description: &description,
		Collapsed:   data.Collapsed.ValueBoolPointer(),
		Position:    hundApiV1.ToIntPtr(data.Position.ValueInt64Pointer()),
	}

	rsp, err := r.client.UpdateAGroup(ctx, data.Id.ValueString(), form, hundApiV1.Unexpand("components"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Hund Group",
			err.Error(),
		)
		return
	}

	group, err := hundApiV1.ParseUpdateAGroupResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Group",
			err.Error(),
		)
		return
	}

	if group.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(group.StatusCode())+
				"\nError: "+string(group.Body),
		)
		return
	}

	newState, diag := models.ToGroupModel(*group.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// if len(data.Components.Elements()) > 0 {
	// 	resp.Diagnostics.AddError(
	// 		"Cannot delete a non-empty Group",
	// 		"This Group cannot be deleted unless its Components are moved to another"+
	// 			" Group, or otherwise deleted.",
	// 	)
	// 	return
	// }

	rsp, err := r.client.DeleteAGroup(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund Group",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		group, err := hundApiV1.ParseDeleteAGroupResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(group.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
