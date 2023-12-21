package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
	"github.com/hundio/terraform-provider-hund/internal/planmodifiers"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IssueUpdateResource{}
var _ resource.ResourceWithConfigure = &IssueUpdateResource{}
var _ resource.ResourceWithImportState = &IssueUpdateResource{}
var _ resource.ResourceWithConfigValidators = &IssueUpdateResource{}
var _ resource.ResourceWithModifyPlan = &IssueUpdateResource{}

func NewIssueUpdateResource() resource.Resource {
	return &IssueUpdateResource{}
}

// IssueUpdateResource defines the resource implementation.
type IssueUpdateResource struct {
	client *hundApiV1.Client
}

// IssueUpdateResourceModel describes the resource data model.
type IssueUpdateResourceModel models.UpdateModel

func (r *IssueUpdateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_update"
}

func (r *IssueUpdateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Issue Updates describe a particular phase in the evolution of an issue. Updates have their own bodies and label, and can also change the current state override of the Issue. Updates are also responsible for resolving/reopening Issues, as well as adding addendums/postmortems to the end of an Issue.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("Update"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"archive_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "When true, this Update will not be destroyed from your status page if the resource is destroyed in your Terraform configuration. This option is **recommended** for maintaining a history on your status page of past Issues.",
				Optional:            true,
			},
			"issue_id": schema.StringAttribute{
				MarkdownDescription: "The Issue that this Update pertains to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: createdAtFieldMarkdownDescription("Update"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: updatedAtFieldMarkdownDescription("Update"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"effective": schema.BoolAttribute{
				MarkdownDescription: "When true, denotes that this Update is the latest update on this Issue (hence, the \"effective\" Update according to `effective_after`).",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"reopening": schema.BoolAttribute{
				MarkdownDescription: "Whether this Update reopened the Issue if it was already resolved in an Update before this one.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"effective_after": schema.StringAttribute{
				MarkdownDescription: "The time after which this Update is considered the latest Update on its Issue, until the `effective_after` time of the Update succeeding this one, if one exists.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The body text of this Update in raw markdown."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body_html": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("An HTML rendered view of the markdown in `body`."),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("body")),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("body_translations")),
				},
			},
			"body_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The body text of this Update in raw markdown."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"body_html_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("An HTML rendered view of the markdown in `body`."),
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("body")),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("body_translations")),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "The label applied to this update, as well as the issue at large when this Update is the *latest* Update in the Issue. The label can be thought of as the \"state\" of the Issue as of this Update (e.g. \"Problem Identified\", \"Monitoring\", \"Resolved\").",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state_override": schema.Int64Attribute{
				MarkdownDescription: "The integer state which overrides the state of affected Components in `component`. A value of `null` indicates no override is present.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 1),
				},
			},
			"template": issueTemplateApplicationSchema(),
		},
	}
}

func (r *IssueUpdateResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("body"),
			path.MatchRoot("body_translations"),
			path.MatchRoot("template"),
		),
	}
}

func (r *IssueUpdateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		var data IssueUpdateResourceModel

		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

		if !data.ArchiveOnDestroy.ValueBool() {
			if data.Reopening.ValueBool() || data.Label.ValueString() == string(hundApiV1.UPDATELABELResolved) {
				resp.Diagnostics.AddWarning(
					"Cannot destroy Reopening/Resolving Issue Updates",
					"Updates that affect the standing/resolved state of an Issue cannot be deleted. "+
						"Applying this resource destruction will only remove this Update from "+
						"the Terraform state. This Update will only be destroyed if the Issue "+
						"itself is also destroyed.",
				)
			}

			resp.Diagnostics.Append(IssueAndUpdateDestructionWarning())
		}

		return
	}

	if req.State.Raw.IsNull() {
		return
	}

	var plan, state, config IssueUpdateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	planmodifiers.PlanModifyIssueTemplateUpdateApplication(ctx, path.Root("template"), state.Template, plan.Template, resp)

	var templateState, templatePlan attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("template"), &templatePlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("template"), &templateState)...)

	templateChanged := !templateState.Equal(templatePlan)

	if templateChanged {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("label"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body_html"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body_html_translations"), types.MapUnknown(types.StringType))...)
	} else if state.Label.IsNull() && config.Label.IsNull() && plan.Label.IsUnknown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("label"), types.StringNull())...)
	}

	if templateChanged || !state.Body.Equal(plan.Body) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body_translations"), types.MapUnknown(types.StringType))...)
	}

	if templateChanged || !state.BodyTranslations.Equal(plan.BodyTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body"), types.StringUnknown())...)
	}

	if !plan.EffectiveAfter.Equal(state.EffectiveAfter) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("effective"), types.BoolUnknown())...)
	}

	if !req.Plan.Raw.Equal(req.State.Raw) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated_at"), types.StringUnknown())...)
	}
}

func (r *IssueUpdateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IssueUpdateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, config IssueUpdateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	effectiveAfter, err := hundApiV1.ToIntTimestampPtr(config.EffectiveAfter.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.Append(models.TimestampError(err))
		return
	}

	form := hundApiV1.UpdateFormCreate{
		EffectiveAfter: effectiveAfter,
		StateOverride:  hundApiV1.DblPtr((*hundApiV1.IntegerState)(hundApiV1.ToIntPtr(data.StateOverride.ValueInt64Pointer()))),
	}

	if data.Template != nil {
		templateForm := hundApiV1.IssueTemplateApplicationFormCreate{}

		prepareIssueTemplateApplicationUpdate(ctx, *data.Template, &templateForm, &resp.Diagnostics)

		if resp.Diagnostics.HasError() {
			return
		}

		opaqueForm := hundApiV1.UpdateFormCreate_Template{}
		err = opaqueForm.FromUpdateFormCreateTemplate1(templateForm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Template conversion error",
				"Got error encoding IssueTemplateApplication: "+err.Error(),
			)
			return
		}

		form.Template = &opaqueForm
	} else {
		body, err := hundApiV1.ToI18nStringPtr(data.Body, data.BodyTranslations)
		if err != nil {
			resp.Diagnostics.Append(models.I18nStringError(err))
			return
		}

		form.Body = hundApiV1.DblPtr(body)

		form.Label = hundApiV1.DblPtr((*hundApiV1.UPDATELABEL)(data.Label.ValueStringPointer()))
	}

	rsp, err := r.client.CreateAUpdate(ctx, data.IssueId.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Issue Update",
			err.Error(),
		)
		return
	}

	update, err := hundApiV1.ParseCreateAUpdateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	if update.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(update.StatusCode())+
				"\nError: "+string(update.Body),
		)
		return
	}

	apiUpdate, err := hundApiV1.UnwrapUpdateExpansionary(*update.HALJSON201)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	newState, diag := models.ToUpdateModel(ctx, *apiUpdate)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.ArchiveOnDestroy = data.ArchiveOnDestroy

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueUpdateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAUpdate(ctx, data.IssueId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Issue Update",
			err.Error(),
		)
		return
	}

	update, err := hundApiV1.ParseRetrieveAUpdateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	if update.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(update.StatusCode())+
				"\nError: "+string(update.Body),
		)
		return
	}

	apiUpdate, err := hundApiV1.UnwrapUpdateExpansionary(*update.HALJSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	newState, diag := models.ToUpdateModel(ctx, *apiUpdate)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.ArchiveOnDestroy = data.ArchiveOnDestroy

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, config IssueUpdateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	effectiveAfter, err := hundApiV1.ToIntTimestampPtr(config.EffectiveAfter.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.Append(models.TimestampError(err))
		return
	}

	form := hundApiV1.UpdateFormUpdate{
		EffectiveAfter: effectiveAfter,
		StateOverride:  hundApiV1.DblPtr((*hundApiV1.IntegerState)(hundApiV1.ToIntPtr(data.StateOverride.ValueInt64Pointer()))),
	}

	if data.Template != nil {
		templateForm := hundApiV1.IssueTemplateApplicationFormCreate{}

		prepareIssueTemplateApplicationUpdate(ctx, *data.Template, &templateForm, &resp.Diagnostics)

		if resp.Diagnostics.HasError() {
			return
		}

		opaqueForm := hundApiV1.UpdateFormUpdate_Template{}
		err = opaqueForm.FromUpdateFormUpdateTemplate1(templateForm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Template conversion error",
				"Got error encoding IssueTemplateApplication: "+err.Error(),
			)
			return
		}

		form.Template = &opaqueForm
	} else {
		body, err := hundApiV1.ToI18nStringPtr(data.Body, data.BodyTranslations)
		if err != nil {
			resp.Diagnostics.Append(models.I18nStringError(err))
			return
		}

		form.Body = hundApiV1.DblPtr(body)
		form.Label = hundApiV1.DblPtr((*hundApiV1.UPDATELABEL)(data.Label.ValueStringPointer()))
	}

	rsp, err := r.client.ReviseAnUpdate(ctx, data.IssueId.ValueString(), data.Id.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Issue Update",
			err.Error(),
		)
		return
	}

	update, err := hundApiV1.ParseReviseAnUpdateResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	if update.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(update.StatusCode())+
				"\nError: "+string(update.Body),
		)
		return
	}

	apiUpdate, err := hundApiV1.UnwrapUpdateExpansionary(*update.HALJSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Update",
			err.Error(),
		)
		return
	}

	newState, diag := models.ToUpdateModel(ctx, *apiUpdate)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.ArchiveOnDestroy = data.ArchiveOnDestroy

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueUpdateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueUpdateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ArchiveOnDestroy.ValueBool() {
		return
	}

	if data.Reopening.ValueBool() || data.Label.ValueString() == string(hundApiV1.UPDATELABELResolved) {
		return
	}

	rsp, err := r.client.DeleteAUpdate(ctx, data.IssueId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund Issue Update",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		update, err := hundApiV1.ParseDeleteAUpdateResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(update.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *IssueUpdateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	issueId, id, ok := strings.Cut(req.ID, "/")

	if !ok {
		resp.Diagnostics.AddError(
			"Could not parse Hund Issue Update Import ID",
			"To properly import an Issue Update, please provide both the Issue ID, a"+
				" forward slash, and finally the Issue Update ID: <ISSUE_ID>/<UPDATE_ID>",
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("issue_id"), issueId)...)
}

func prepareIssueTemplateApplicationUpdate(ctx context.Context, model models.IssueTemplateApplicationUpdateModel, form *hundApiV1.IssueTemplateApplicationFormCreate, diags *diag.Diagnostics) {
	variables, diag0 := model.Variables.ApiValue()
	diags.Append(diag0...)
	if diags.HasError() {
		return
	}

	body, err := hundApiV1.ToI18nStringPtr(model.Body, model.BodyTranslations)
	if err != nil {
		diags.Append(models.I18nStringError(err))
		return
	}

	form.IssueTemplate = model.IssueTemplateId.ValueString()
	form.Body = body
	form.Variables = &variables

	if !model.Label.IsUnknown() {
		labelPtr := hundApiV1.MaybeISSUETEMPLATELABEL(model.Label.ValueStringPointer())
		form.Label = hundApiV1.DblPtr(labelPtr)
	}

	if !model.Schema.IsUnknown() {
		schema := map[string]models.IssueTemplateVariableModel{}
		diags.Append(model.Schema.ElementsAs(ctx, &schema, false)...)
		schemaForm := models.ToIssueTemplateVariablesForm(schema)

		form.Schema = &schemaForm
	}
}
