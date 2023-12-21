package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
	"github.com/hundio/terraform-provider-hund/internal/planmodifiers"
	"github.com/hundio/terraform-provider-hund/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ComponentResource{}
var _ resource.ResourceWithConfigure = &ComponentResource{}
var _ resource.ResourceWithImportState = &ComponentResource{}
var _ resource.ResourceWithConfigValidators = &ComponentResource{}
var _ resource.ResourceWithModifyPlan = &ComponentResource{}

func NewComponentResource() resource.Resource {
	return &ComponentResource{}
}

// ComponentResource defines the resource implementation.
type ComponentResource struct {
	client *hundApiV1.Client
}

// ComponentResourceModel describes the resource data model.
type ComponentResourceModel models.ComponentModel

func (r *ComponentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

func (r *ComponentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "An object representing a specific part of a service, which is potentially subject to downtime.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ObjectId of this Component.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp at which this Component was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp at which this Component was last updated.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The ID of the Group that this Component belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("The name of this Component."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("The name of this Component."),
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: translationOriginalFieldMarkdownDescription("A description of this Component, potentially with markdown formatting."),
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description_translations": schema.MapAttribute{
				MarkdownDescription: translationFieldMarkdownDescription("A description of this Component, potentially with markdown formatting."),
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
			"exclude_from_global_history": schema.BoolAttribute{
				MarkdownDescription: "Exclude this Component's uptime percentage from being factored into the global percent uptime calculation.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"exclude_from_global_uptime": schema.BoolAttribute{
				MarkdownDescription: "Exclude this Component from appearing in the global history.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"last_event_at": schema.StringAttribute{
				MarkdownDescription: "A timestamp at which the last event for this Component occurred. This includes automated status changes, as well as issue creation and update.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"percent_uptime": schema.Float64Attribute{
				MarkdownDescription: "The rolling 30-day percent uptime of this Component.",
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"watchdog": schema.SingleNestedAttribute{
				MarkdownDescription: "The Watchdog that supplies the current status for this Component.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The ObjectId of this Watchdog.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"high_frequency": schema.BoolAttribute{
						MarkdownDescription: "When true, this Watchdog will run every 30 seconds, instead of the standard 1 minute.\n\n-> You are billed extra for each high frequency Watchdog. Please see our [pricing page](https://hund.io/pricing) for more details.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"latest_status": schema.StringAttribute{
						MarkdownDescription: "The ObjectId of the latest Status object generated by this Watchdog. When `null`, this Watchdog is still pending initial status.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"service": schema.SingleNestedAttribute{
						MarkdownDescription: "The service configuration for this Watchdog, which describes how the Watchdog determines current status.",
						Required:            true,
						Validators: []validator.Object{
							validators.ExactlyOneNonNullAttribute(),
						},
						Attributes: map[string]schema.Attribute{
							"manual": schema.SingleNestedAttribute{
								MarkdownDescription: "A manually updated Watchdog.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"state": schema.Int64Attribute{
										MarkdownDescription: "An integer denoting operational state (1 => operational, 0 => degraded, -1 => outage).",
										Optional:            true,
										Computed:            true,
										Default:             int64default.StaticInt64(1),
										Validators: []validator.Int64{
											int64validator.Between(-1, 1),
										},
									},
								},
							},
							"updown": schema.SingleNestedAttribute{
								MarkdownDescription: "An [Updown.io](https://updown.io) service.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"monitor_api_key": schema.StringAttribute{
										MarkdownDescription: "An Updown.io monitor API key. This API key can be read-only.",
										Required:            true,
										Sensitive:           true,
									},
									"monitor_token": schema.StringAttribute{
										MarkdownDescription: "An Updown.io monitor token to retrieve status from.",
										Required:            true,
									},
								},
							},
							"pingdom": schema.SingleNestedAttribute{
								MarkdownDescription: "A [pingdom](https://www.pingdom.com) service.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"api_token": schema.StringAttribute{
										MarkdownDescription: "The Pingdom API v3 key.",
										Required:            true,
										Sensitive:           true,
									},
									"check_id": schema.StringAttribute{
										MarkdownDescription: "The ID of the check to pull status from on Pingdom.",
										Required:            true,
									},
									"check_type": schema.StringAttribute{
										MarkdownDescription: "The type of the Pingdom check. `check` denotes a normal Pingdom uptime check, and `transactional` denotes a Pingdom TMS check.",
										Optional:            true,
										Computed:            true,
										Default:             stringdefault.StaticString("check"),
										Validators: []validator.String{
											stringvalidator.OneOf(
												string(hundApiV1.PINGDOMCHECKTYPECheck),
												string(hundApiV1.PINGDOMCHECKTYPETransactional),
											),
										},
									},
									"unconfirmed_is_down": schema.BoolAttribute{
										MarkdownDescription: "When true, triggers Watchdog outage when Pingdom reports a yet unconfirmed outage.",
										Optional:            true,
										Computed:            true,
										Default:             booldefault.StaticBool(false),
									},
								},
							},
							"uptimerobot": schema.SingleNestedAttribute{
								MarkdownDescription: "An [Uptime Robot](https://uptimerobot.com) service.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"monitor_api_key": schema.StringAttribute{
										MarkdownDescription: "An Uptime Robot monitor API key to retrieve status from.",
										Required:            true,
										Sensitive:           true,
									},
									"unconfirmed_is_down": schema.BoolAttribute{
										MarkdownDescription: "When true, triggers Watchdog outage when UptimeRobot reports a yet unconfirmed outage.",
										Optional:            true,
										Computed:            true,
										Default:             booldefault.StaticBool(false),
									},
								},
							},
							"webhook": schema.SingleNestedAttribute{
								MarkdownDescription: "A [webhook](https://hund.io/help/integrations/webhooks) service.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"webhook_key": schema.StringAttribute{
										MarkdownDescription: "The key to use for this webhook, expected in request headers.",
										Optional:            true,
										Computed:            true,
										Sensitive:           true,
									},
									"deadman": schema.BoolAttribute{
										MarkdownDescription: "When true, turns on a \"Dead Man's Switch\" for the Watchdog, according to the" +
											" configuration set by `reporting_interval` and `consecutive_checks`. The Watchdog" +
											" will trigger an \"outage\" state if the webhook does not receive a call after" +
											" the configured number of consecutive checks (according to the reporting interval)." +
											" This switch can be useful when a lack of webhook reporting from the specific" +
											" component should be taken to mean that the component itself is down.,",
										Optional: true,
										Computed: true,
										Default:  booldefault.StaticBool(false),
									},
									"consecutive_checks": schema.Int64Attribute{
										MarkdownDescription: "This property is only required when `deadman: true`. This property configures" +
											" how many checks (i.e. the number of times `reporting_interval` elapses) must" +
											" fail (i.e. no status reported to the webhook) before triggering the \"Dead Man's" +
											" Switch.\"" +
											" For example, if `deadman: true` and `reporting_interval: 60`, then a setting" +
											" of `consecutive_checks: 5` would cause the Watchdog to wait for 5 consecutive" +
											" minutes to receive a webhook call before triggering outage. Since the count is" +
											" consecutive, it is reset whenever a new webhook call comes through to the Watchdog.",
										Optional: true,
									},
									"reporting_interval": schema.Int64Attribute{
										MarkdownDescription: "This property is only required when `deadman: true`. This property configures how often (in seconds) that you expect to POST status to the webhook.",
										Optional:            true,
									},
								},
							},
							"icmp": nativeIcmpServiceSchema(),
							"http": nativeHttpServiceSchema(),
							"dns":  nativeDnsServiceSchema(),
							"tcp":  nativeTcpServiceSchema(),
							"udp":  nativeUdpServiceSchema(),
						},
					},
				},
			},
		},
	}
}

func (r *ComponentResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
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

func (r *ComponentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, config, state ComponentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	planmodifiers.WatchdogComputeHighFrequency(ctx, path.Root("watchdog"), plan.Watchdog, config.Watchdog, resp)

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
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("last_event_at"), types.StringUnknown())...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("percent_uptime"), types.Float64Unknown())...)
	}
}

func (r *ComponentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ComponentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ComponentResourceModel

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

	serviceForm, err := data.Watchdog.Service.ApiCreateForm()
	if err != nil {
		resp.Diagnostics.Append(WatchdogServiceError(err))
	}

	form := hundApiV1.ComponentFormCreate{
		Name:                     name,
		Group:                    data.Group.ValueString(),
		Description:              &description,
		ExcludeFromGlobalHistory: data.ExcludeFromGlobalHistory.ValueBoolPointer(),
		ExcludeFromGlobalUptime:  data.ExcludeFromGlobalUptime.ValueBoolPointer(),
		Watchdog: hundApiV1.WatchdogFormCreate{
			HighFrequency: data.Watchdog.HighFrequency.ValueBoolPointer(),
			Service:       serviceForm,
		},
	}

	rsp, err := r.client.CreateAComponent(ctx, form, hundApiV1.Expand("watchdog"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund Component",
			err.Error(),
		)
		return
	}

	component, err := hundApiV1.ParseCreateAComponentResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Component",
			err.Error(),
		)
		return
	}

	if component.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-201 status code: "+fmt.Sprint(component.StatusCode())+
				"\nError: "+string(component.Body),
		)
		return
	}

	newState, diag := models.ToComponentModel(ctx, *component.HALJSON201)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.Watchdog.Service.ReplaceSensitiveAttributes(data.Watchdog.Service)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ComponentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ComponentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAComponent(ctx, data.Id.ValueString(), hundApiV1.Expand("watchdog"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Component",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	component, err := hundApiV1.ParseRetrieveAComponentResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Component",
			err.Error(),
		)
		return
	}

	if component.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(component.StatusCode())+
				"\nError: "+string(component.Body),
		)
		return
	}

	newState, diag := models.ToComponentModel(ctx, *component.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Watchdog != nil {
		newState.Watchdog.Service.ReplaceSensitiveAttributes(data.Watchdog.Service)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ComponentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ComponentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	var watchdogPlan, watchdogState types.Object

	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, path.Root("watchdog"), &watchdogPlan)...,
	)
	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, path.Root("watchdog"), &watchdogState)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	var watchdogForm *hundApiV1.WatchdogFormUpdate

	servicePlan := watchdogPlan.Attributes()["service"].(types.Object)
	serviceState := watchdogState.Attributes()["service"].(types.Object)

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

	if watchdogServiceTypeChanged(servicePlan, serviceState) {
		keepOriginalMetrics := false

		serviceForm, err := data.Watchdog.Service.ApiCreateForm()
		if err != nil {
			resp.Diagnostics.Append(WatchdogServiceError(err))
			return
		}

		conversionForm := hundApiV1.WatchdogFormConvert{
			KeepOriginalDefaultMetricProvider: &keepOriginalMetrics,
			HighFrequency:                     data.Watchdog.HighFrequency.ValueBoolPointer(),
			Service:                           serviceForm,
		}

		rsp, err := r.client.ConvertAWatchdogsServiceType(ctx, data.Watchdog.Id.ValueString(), conversionForm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Hund Watchdog Service",
				err.Error(),
			)
			return
		}

		watchdog, err := hundApiV1.ParseConvertAWatchdogsServiceTypeResponse(rsp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Parse Hund Watchdog",
				err.Error(),
			)
			return
		}

		if watchdog.StatusCode() != 200 {
			resp.Diagnostics.AddError(
				"Failed response code from Hund API",
				"Received a non-200 status code: "+fmt.Sprint(watchdog.StatusCode())+
					"\nError: "+string(watchdog.Body),
			)
			return
		}
	} else if !watchdogPlan.Equal(watchdogState) {
		serviceForm, err := data.Watchdog.Service.ApiUpdateForm()
		if err != nil {
			resp.Diagnostics.Append(WatchdogServiceError(err))
			return
		}

		watchdogForm = &hundApiV1.WatchdogFormUpdate{
			HighFrequency: data.Watchdog.HighFrequency.ValueBoolPointer(),
			Service:       serviceForm,
		}
	}

	form := hundApiV1.ComponentFormUpdate{
		Name:                     &name,
		Description:              &description,
		Group:                    data.Group.ValueStringPointer(),
		ExcludeFromGlobalHistory: data.ExcludeFromGlobalHistory.ValueBoolPointer(),
		ExcludeFromGlobalUptime:  data.ExcludeFromGlobalUptime.ValueBoolPointer(),
		Watchdog:                 watchdogForm,
	}

	rsp, err := r.client.UpdateAComponent(ctx, data.Id.ValueString(), form, hundApiV1.Expand("watchdog"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Hund Component",
			err.Error(),
		)
		return
	}

	component, err := hundApiV1.ParseUpdateAComponentResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Component",
			err.Error(),
		)
		return
	}

	if component.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(component.StatusCode())+
				"\nError: "+string(component.Body),
		)
		return
	}

	newState, diag := models.ToComponentModel(ctx, *component.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.Watchdog.Service.ReplaceSensitiveAttributes(data.Watchdog.Service)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ComponentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ComponentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.DeleteAComponent(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund Component",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		component, err := hundApiV1.ParseDeleteAComponentResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(component.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *ComponentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
