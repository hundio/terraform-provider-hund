package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupComponentOrderingResource{}
var _ resource.ResourceWithConfigure = &GroupComponentOrderingResource{}
var _ resource.ResourceWithImportState = &GroupComponentOrderingResource{}
var _ resource.ResourceWithConfigValidators = &GroupComponentOrderingResource{}

// var _ resource.ResourceWithModifyPlan = &GroupComponentOrderingResource{}

func NewGroupComponentOrderingResource() resource.Resource {
	return &GroupComponentOrderingResource{}
}

// GroupComponentOrderingResource defines the resource implementation.
type GroupComponentOrderingResource struct {
	client *hundApiV1.Client
}

// GroupComponentOrderingResourceModel describes the resource data model.
type GroupComponentOrderingResourceModel models.GroupComponentOrderingModel

func (r *GroupComponentOrderingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_component_ordering"
}

func (r *GroupComponentOrderingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource representing the ordering of Components within a Group",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The Group whose ordering is managed by this resource.",
				Required:            true,
			},
			"components": schema.ListAttribute{
				MarkdownDescription: "The list of Component IDs in this Group, listed in the exact order they will appear under the Group.\n\n~> This list **must not** omit nor add any Components not already in the referenced Group, or an error will occur. This resource is **only** for managing an order for the Components of a Group.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
		},
	}
}

func (r *GroupComponentOrderingResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}

func (r *GroupComponentOrderingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupComponentOrderingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupComponentOrderingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState := r.commitGroupComponentOrdering(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *GroupComponentOrderingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupComponentOrderingResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAGroup(ctx, data.Group.ValueString(), hundApiV1.Unexpand("components"))
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

	newState, diag := models.ToGroupComponentOrderingModel(*group.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *GroupComponentOrderingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupComponentOrderingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState := r.commitGroupComponentOrdering(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *GroupComponentOrderingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "deleting group_component_ordering")
}

func (r *GroupComponentOrderingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("group"), req, resp)
}

func (r *GroupComponentOrderingResource) commitGroupComponentOrdering(ctx context.Context, data GroupComponentOrderingResourceModel, diags *diag.Diagnostics) *models.GroupComponentOrderingModel {
	ordering := []string{}

	for _, sv := range data.Components {
		ordering = append(ordering, sv.ValueString())
	}

	rsp, err := r.client.ReorderAGroupsComponents(ctx, data.Group.ValueString(), ordering, hundApiV1.Unexpand("components"))
	if err != nil {
		diags.AddError(
			"Unable to Reorder Hund Group's Components",
			err.Error(),
		)
		return nil
	}

	if rsp.StatusCode == 404 {
		diags.AddAttributeError(
			path.Root("group"),
			"Group ID Does not Exist",
			"A group_component_ordering requires an existing Group. The given ID was not found.",
		)
		return nil
	}

	group, err := hundApiV1.ParseReorderAGroupsComponentsResponse(rsp)
	if err != nil {
		diags.AddError(
			"Unable to Parse Group Reorder Error",
			err.Error(),
		)
		return nil
	}

	if group.StatusCode() != 200 {
		diags.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(group.StatusCode())+
				"\nError: "+string(group.Body),
		)
		return nil
	}

	newState, diag := models.ToGroupComponentOrderingModel(*group.HALJSON200)
	diags.Append(diag...)

	return &newState
}
