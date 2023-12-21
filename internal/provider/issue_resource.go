package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
var _ resource.Resource = &IssueResource{}
var _ resource.ResourceWithConfigure = &IssueResource{}
var _ resource.ResourceWithImportState = &IssueResource{}
var _ resource.ResourceWithConfigValidators = &IssueResource{}
var _ resource.ResourceWithModifyPlan = &IssueResource{}

func NewIssueResource() resource.Resource {
	return &IssueResource{}
}

// IssueResource defines the resource implementation.
type IssueResource struct {
	client *hundApiV1.Client
}

// IssueResourceModel describes the resource data model.
type IssueResourceModel models.IssueModel

func (r *IssueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue"
}

func (r *IssueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The Issue resource represents an evolving incident in time. Issues have updates, which describe the evolution of an issue, often up to resolution. Issues may also set a schedule, which allows automatically starting and ending the issue.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("Issue"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"archive_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "When true, this Issue will not be destroyed from your status page if the resource is destroyed in your Terraform configuration. This option is **recommended** for maintaining a history on your status page of past Issues.",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: createdAtFieldMarkdownDescription("Issue"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: updatedAtFieldMarkdownDescription("Issue"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"began_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp at which this Issue began affecting its given Components.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					planmodifiers.UnknownOnDependentFieldChange(path.Root("schedule").AtName("starts_at")),
				},
			},
			"ended_at": schema.StringAttribute{
				MarkdownDescription: "The UNIX timestamp at which this Issue stopped affecting its given Components. This field is `null` if it has not ended yet.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.NullDefault(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cancelled_at": schema.StringAttribute{
				MarkdownDescription: "The time at which this Issue was cancelled. This field is `null` if the Issue has not been cancelled or is not scheduled.",
				Computed:            true,
				Default:             planmodifiers.NullDefault(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The title of the Issue."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The initial body text of the issue in raw markdown."),
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
			"title_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The title of the Issue."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"body_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The initial body text of the issue in raw markdown."),
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
				MarkdownDescription: "The initial label applied to the issue. The \"current\" label of the entire issue may be updated by the labels of Issue Updates, though this must be taken from the latest Update in `updates`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"specialization": schema.StringAttribute{
				MarkdownDescription: "Whether this Issue has special abilities or connotations. `general` is the default behavior, indicating no specialization. Other values include `maintenance`, which indicates an Issue shows affected components as \"under maintenance,\" and `informational`, which indicates that the Issue is an informational bulletin.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"duration": schema.Int64Attribute{
				MarkdownDescription: "The effective duration of this Issue in seconds. That is, the total amount of time for which this Issue affects its Components. Thus, this field only accumulates while the Issue is ongoing/open.\n\n-> This value is zero for cancelled and informational Issues. For scheduled Issues, this field will remain zero until the Issue begins according to the Schedule.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"open_graph_image_url": schema.StringAttribute{
				MarkdownDescription: "The URL to an image which will be displayed alongside this issue when shared on social media websites.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The integer priority of the Issue. Priority pertains to how notifications are\ntriggered for this Issue: -1 indicates **low priority**, meaning no\nnotifications whatsoever will be triggered for this issue; 0 indicates\n**normal priority**, which is the default behavior; and, 1 indicates\n**high priority**, meaning all subscriptions across all notifiers will receive\nnotifications for this Issue regardless of their notification preferences.\n",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Default: int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(-1, 1),
				},
			},
			"state_override": schema.Int64Attribute{
				MarkdownDescription: "The integer state which overrides the state of affected Components in\n`component`. A value of `null` indicates no override is present.\n",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(-1, 1),
				},
			},
			"resolved": schema.BoolAttribute{
				MarkdownDescription: "Whether this Issue is currently resolved, thus no longer affecting its given\nComponents.\n",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"retrospective": schema.BoolAttribute{
				MarkdownDescription: "Whether this Issue is retrospective; that is, the Issue was created both resolved\n*and* backdated.\n",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"scheduled": schema.BoolAttribute{
				MarkdownDescription: "Whether this Issue has a Schedule.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"standing": schema.BoolAttribute{
				MarkdownDescription: "Whether this Issue is currently active and affecting its given Components.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"component_ids": schema.SetAttribute{
				MarkdownDescription: "The Components IDs affected by this Issue.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"schedule": schema.SingleNestedAttribute{
				MarkdownDescription: "An object detailing the Schedule of this issue if it is scheduled. This field is `null` if the Issue is not scheduled.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: idFieldMarkdownDescription("Schedule"),
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"started": schema.BoolAttribute{
						MarkdownDescription: "Whether this scheduled Issue has started.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"ended": schema.BoolAttribute{
						MarkdownDescription: "Whether this scheduled Issue has ended.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"notified": schema.BoolAttribute{
						MarkdownDescription: "Whether this scheduled Issue has fired an `issue_upcoming` notification.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"starts_at": schema.StringAttribute{
						MarkdownDescription: "The time at which this scheduled Issue will begin.",
						Required:            true,
					},
					"ends_at": schema.StringAttribute{
						MarkdownDescription: "The time at which this scheduled Issue will end.",
						Required:            true,
					},
					"notify_subscribers_at": schema.StringAttribute{
						MarkdownDescription: "The time at which this scheduled Issue will fire a \"heads-up\" `issue_upcoming` notification, informing subscribers that the Issue will begin soon. This field is `null` if the Issue will not be sending an `issue_upcoming` notification.\n\n-> This field cannot be changed once the Issue has emitted an `issue_upcoming` notification. Consider creating a new Update if you'd like to remind subscribers additional times apart from this automated notification.",
						Optional:            true,
						Computed:            true,
						Default:             planmodifiers.NullDefault(),
					},
				},
			},
			"template": issueTemplateApplicationIssueSchema(),
			"updates": schema.ListNestedAttribute{
				MarkdownDescription: "An optional list of Updates to create when initially creating the Issue. When creating a sequence of Updates, ensure that their `effective_after` timestamps do not encroach upon one another, or an error will occur.\n\n~> This field is primarily meant for assisting with the creation of retrospective Issues, rather than creating new Updates as they arise. Please use the full hund_issue_update resource instead, rather than configuring this field directly.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					models.UpdateList{},
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: idFieldMarkdownDescription("Update"),
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"archive_on_destroy": schema.BoolAttribute{
							MarkdownDescription: "This field is unused when embedded in Issues.",
							Computed:            true,
						},
						"issue_id": schema.StringAttribute{
							MarkdownDescription: "The Issue that this Update pertains to.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
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
						},
						"effective": schema.BoolAttribute{
							MarkdownDescription: "When true, denotes that this Update is the latest update on this Issue (hence, the \"effective\" Update according to `effective_after`).",
							Computed:            true,
						},
						"reopening": schema.BoolAttribute{
							MarkdownDescription: "Whether this Update reopened the Issue if it was already resolved in an Update before this one.",
							Computed:            true,
						},
						"effective_after": schema.StringAttribute{
							MarkdownDescription: "The time after which this Update is considered the latest Update on its Issue, until the `effective_after` time of the Update succeeding this one, if one exists.",
							Computed:            true,
							Optional:            true,
						},
						"body": schema.StringAttribute{
							MarkdownDescription: translationOriginalFieldMarkdownDescription("The body text of this Update in raw markdown."),
							Optional:            true,
							Computed:            true,
						},
						"body_html": schema.StringAttribute{
							MarkdownDescription: translationOriginalFieldMarkdownDescription("An HTML rendered view of the markdown in `body`."),
							Computed:            true,
						},
						"body_translations": schema.MapAttribute{
							MarkdownDescription: translationFieldMarkdownDescription("The body text of this Update in raw markdown."),
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
						},
						"body_html_translations": schema.MapAttribute{
							MarkdownDescription: translationFieldMarkdownDescription("An HTML rendered view of the markdown in `body`."),
							Computed:            true,
							ElementType:         types.StringType,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The label applied to this update, as well as the issue at large when this Update is the *latest* Update in the Issue. The label can be thought of as the \"state\" of the Issue as of this Update (e.g. \"Problem Identified\", \"Monitoring\", \"Resolved\").",
							Optional:            true,
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
				},
			},
		},
	}
}

func (r *IssueResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("body"),
			path.MatchRoot("body_translations"),
			path.MatchRoot("template"),
		),
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("title"),
			path.MatchRoot("title_translations"),
			path.MatchRoot("template"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("began_at"),
			path.MatchRoot("schedule"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("ended_at"),
			path.MatchRoot("updates"),
		),
	}
}

func (r *IssueResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		var data IssueResourceModel

		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

		if !data.ArchiveOnDestroy.ValueBool() {
			resp.Diagnostics.Append(IssueAndUpdateDestructionWarning())
		}

		return
	}

	if req.State.Raw.IsNull() {
		return
	}

	var plan, state, config IssueResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	planmodifiers.PlanModifyIssueTemplateIssueApplication(ctx, path.Root("template"), state.Template, plan.Template, resp)

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

	if templateChanged || !state.Title.Equal(plan.Title) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("title_translations"), types.MapUnknown(types.StringType))...)
	}

	if templateChanged || !state.TitleTranslations.Equal(plan.TitleTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("title"), types.StringUnknown())...)
	}

	if templateChanged || !state.Body.Equal(plan.Body) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body_translations"), types.MapUnknown(types.StringType))...)
	}

	if templateChanged || !state.BodyTranslations.Equal(plan.BodyTranslations) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("body"), types.StringUnknown())...)
	}

	if !resp.Plan.Raw.Equal(req.State.Raw) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("updated_at"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("duration"), types.Int64Unknown())...)
	}
}

func (r *IssueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IssueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueResourceModel

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

	var began *int64
	if !data.BeganAt.IsUnknown() {
		began, err = hundApiV1.ToIntTimestampPtr(data.BeganAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}
	}

	var ended *int64
	if !data.EndedAt.IsUnknown() {
		ended, err = hundApiV1.ToIntTimestampPtr(data.EndedAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}
	}

	var label **hundApiV1.ISSUELABEL
	if !data.Label.IsUnknown() {
		label = hundApiV1.DblPtr((*hundApiV1.ISSUELABEL)(data.Label.ValueStringPointer()))
	}

	form := hundApiV1.IssueFormCreate{
		BeganAt:           began,
		EndedAt:           ended,
		Title:             title,
		Body:              body,
		Components:        hundApiV1.ToStringList(data.ComponentIds),
		Label:             label,
		OpenGraphImageUrl: hundApiV1.DblPtr(data.OpenGraphImageUrl.ValueStringPointer()),
		Priority:          (*hundApiV1.IssueFormCreatePriority)(hundApiV1.ToIntPtr(data.Priority.ValueInt64Pointer())),
		StateOverride:     hundApiV1.DblPtr((*hundApiV1.IntegerState)(hundApiV1.ToIntPtr(data.StateOverride.ValueInt64Pointer()))),
		Schedule:          nil,
	}

	if data.Schedule != nil {
		starts, err := hundApiV1.ToIntTimestamp(data.Schedule.StartsAt.ValueString())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		ends, err := hundApiV1.ToIntTimestamp(data.Schedule.EndsAt.ValueString())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		notify, err := hundApiV1.ToIntTimestampPtr(data.Schedule.NotifySubscribersAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		form.Schedule = &hundApiV1.ScheduleFormCreate{
			StartsAt:            starts,
			EndsAt:              ends,
			NotifySubscribersAt: notify,
		}
	}

	if data.Template != nil {
		templateForm := hundApiV1.IssueTemplateApplicationIssueFormCreate{}

		prepareIssueTemplateApplication(ctx, *data.Template, &templateForm, &resp.Diagnostics)

		if resp.Diagnostics.HasError() {
			return
		}

		opaqueForm := hundApiV1.IssueFormCreate_Template{}
		err = opaqueForm.FromIssueFormCreateTemplate1(templateForm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Template conversion error",
				"Got error encoding IssueTemplateApplication: "+err.Error(),
			)
			return
		}

		form.Template = &opaqueForm
	}

	rsp, err := r.client.CreateAIssue(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Issue",
			err.Error(),
		)
		return
	}

	issue, err := hundApiV1.ParseCreateAIssueResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue",
			err.Error(),
		)
		return
	}

	if issue.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-201 status code: "+fmt.Sprint(issue.StatusCode())+
				"\nError: "+string(issue.Body),
		)
		return
	}

	newState, diag := models.ToIssueModel(ctx, *issue.HALJSON201)
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

func (r *IssueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAIssue(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Issue",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	issue, err := hundApiV1.ParseRetrieveAIssueResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue",
			err.Error(),
		)
		return
	}

	if issue.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(issue.StatusCode())+
				"\nError: "+string(issue.Body),
		)
		return
	}

	newState, diag := models.ToIssueModel(ctx, *issue.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.ArchiveOnDestroy = data.ArchiveOnDestroy

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, config IssueResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

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

	began, err := hundApiV1.ToIntTimestampPtr(config.BeganAt.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.Append(models.TimestampError(err))
		return
	}

	components := hundApiV1.ToStringList(data.ComponentIds)

	form := hundApiV1.IssueFormUpdate{
		BeganAt:           began,
		Title:             title,
		Body:              body,
		Components:        &components,
		Label:             hundApiV1.DblPtr((*hundApiV1.ISSUELABEL)(data.Label.ValueStringPointer())),
		OpenGraphImageUrl: hundApiV1.DblPtr(data.OpenGraphImageUrl.ValueStringPointer()),
		Priority:          (*hundApiV1.ISSUEPRIORITY)(hundApiV1.ToIntPtr(data.Priority.ValueInt64Pointer())),
		StateOverride:     hundApiV1.DblPtr((*hundApiV1.IntegerState)(hundApiV1.ToIntPtr(data.StateOverride.ValueInt64Pointer()))),
		Schedule:          nil,
	}

	if data.Schedule != nil {
		starts, err := hundApiV1.ToIntTimestampPtr(data.Schedule.StartsAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		ends, err := hundApiV1.ToIntTimestampPtr(data.Schedule.EndsAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		notify, err := hundApiV1.ToIntTimestampPtr(data.Schedule.NotifySubscribersAt.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.Append(models.TimestampError(err))
			return
		}

		form.Schedule = &hundApiV1.ScheduleFormUpdate{
			StartsAt:            starts,
			EndsAt:              ends,
			NotifySubscribersAt: notify,
		}
	}

	if data.Template != nil {
		templateForm := hundApiV1.IssueTemplateApplicationIssueFormCreate{}

		prepareIssueTemplateApplication(ctx, *data.Template, &templateForm, &resp.Diagnostics)

		if resp.Diagnostics.HasError() {
			return
		}

		opaqueForm := hundApiV1.IssueFormUpdate_Template{}
		err = opaqueForm.FromIssueFormUpdateTemplate1(templateForm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Template conversion error",
				"Got error encoding IssueTemplateApplication: "+err.Error(),
			)
			return
		}

		form.Template = &opaqueForm
	}

	rsp, err := r.client.ReviseAnIssue(ctx, data.Id.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Issue",
			err.Error(),
		)
		return
	}

	issue, err := hundApiV1.ParseReviseAnIssueResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue",
			err.Error(),
		)
		return
	}

	if issue.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(issue.StatusCode())+
				"\nError: "+string(issue.Body),
		)
		return
	}

	newState, diag := models.ToIssueModel(ctx, *issue.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.ArchiveOnDestroy = data.ArchiveOnDestroy

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *IssueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ArchiveOnDestroy.ValueBool() {
		return
	}

	rsp, err := r.client.DeleteAIssue(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund Issue",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		issue, err := hundApiV1.ParseDeleteAIssueResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(issue.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *IssueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func prepareIssueTemplateApplication(ctx context.Context, model models.IssueTemplateApplicationIssueModel, form *hundApiV1.IssueTemplateApplicationIssueFormCreate, diags *diag.Diagnostics) {
	variables, diag0 := model.Variables.ApiValue()
	diags.Append(diag0...)
	if diags.HasError() {
		return
	}

	title, err := hundApiV1.ToI18nStringPtr(model.Title, model.TitleTranslations)
	if err != nil {
		diags.Append(models.I18nStringError(err))
		return
	}

	body, err := hundApiV1.ToI18nStringPtr(model.Body, model.BodyTranslations)
	if err != nil {
		diags.Append(models.I18nStringError(err))
		return
	}

	form.IssueTemplate = model.IssueTemplateId.ValueString()
	form.Title = title
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
