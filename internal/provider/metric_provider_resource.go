package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
	"github.com/hundio/terraform-provider-hund/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MetricProviderResource{}
var _ resource.ResourceWithConfigure = &MetricProviderResource{}
var _ resource.ResourceWithImportState = &MetricProviderResource{}
var _ resource.ResourceWithConfigValidators = &MetricProviderResource{}
var _ resource.ResourceWithModifyPlan = &MetricProviderResource{}

func NewMetricProviderResource() resource.Resource {
	return &MetricProviderResource{}
}

// MetricProviderResource defines the resource implementation.
type MetricProviderResource struct {
	client *hundApiV1.Client
}

// MetricProviderResourceModel describes the resource data model.
type MetricProviderResourceModel models.MetricProviderModel

func (r *MetricProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_provider"
}

func (r *MetricProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "MetricProviders gather metrics from a configured service for viewing on the Status Page.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: idFieldMarkdownDescription("MetricProvider"),
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"watchdog": schema.StringAttribute{
				MarkdownDescription: "The Watchdog that owns this MetricProvider.",
				Required:            true,
			},
			"default": schema.BoolAttribute{
				MarkdownDescription: "When true, denotes that this MetricProvider is the default MetricProvider of the Watchdog. This implies that they share the same service configuration, which the MetricProvider inherits from the Watchdog. This MetricProvider is created **automatically**, depending on the Watchdog, and cannot be deleted without also deleting the Watchdog.\n\n~> Default MetricProviders cannot be created directly, and must be imported to be managed by Terraform. Deleting a default MetricProvider from your Terraform configuration will only remove the resource from the state.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"instances": schema.MapNestedAttribute{
				MarkdownDescription: "A Map of MetricInstances, which describe each Metric that the MetricProvider provides. The keys of this Map define the slugs of each provided metric.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: idFieldMarkdownDescription("MetricInstance"),
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "A string that uniquely identifies this MetricInstance by referencing the `definition_slug` and MetricProvider `id`.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"definition_slug": schema.StringAttribute{
							MarkdownDescription: "A descriptive string that identifies the metric definition this instance derives from (e.g. `http.tcp_connection_time`, `apdex`, etc.).",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether or not to show this metric on the Component that uses it (through the Watchdog).",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"top_level_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether or not to show this metric on the status page home.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"title": schema.StringAttribute{
							MarkdownDescription: translationOriginalFieldMarkdownDescription("The title of the metric, displayed above its graph on the status page."),
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("title_translations")),
							},
						},
						"title_translations": schema.MapAttribute{
							MarkdownDescription: translationFieldMarkdownDescription("The title of the metric, displayed above its graph on the status page."),
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							PlanModifiers: []planmodifier.Map{
								mapplanmodifier.UseStateForUnknown(),
							},
						},
						"x_title": schema.StringAttribute{
							MarkdownDescription: translationOriginalFieldMarkdownDescription("The title of the x-axis of this metric."),
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("x_title_translations")),
							},
						},
						"x_title_translations": schema.MapAttribute{
							MarkdownDescription: translationFieldMarkdownDescription("The title of the x-axis of this metric."),
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							PlanModifiers: []planmodifier.Map{
								mapplanmodifier.UseStateForUnknown(),
							},
						},
						"y_title": schema.StringAttribute{
							MarkdownDescription: translationOriginalFieldMarkdownDescription("The title of the y-axis of this metric."),
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("y_title_translations")),
							},
						},
						"y_title_translations": schema.MapAttribute{
							MarkdownDescription: translationFieldMarkdownDescription("The title of the y-axis of this metric."),
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
							PlanModifiers: []planmodifier.Map{
								mapplanmodifier.UseStateForUnknown(),
							},
						},
						"x_type": schema.StringAttribute{
							MarkdownDescription: "The type of quantity represented by the x-axis. One of `time` or `measure`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("time", "measure"),
							},
						},
						"y_type": schema.StringAttribute{
							MarkdownDescription: "The type of quantity represented by the y-axis. One of `time` or `measure`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("time", "measure"),
							},
						},
						"y_supremum": schema.Float64Attribute{
							MarkdownDescription: "The least upper bound to display the y-axis on. The metric will always display up to at least this value on the y-axis regardless of the graphed data. If the graph exceeds this value, then the bound will be raised as much as necessary to accommodate the data.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Float64{
								float64planmodifier.UseStateForUnknown(),
							},
						},
						"plot_type": schema.StringAttribute{
							MarkdownDescription: "The kind of visualization to display the metric with. One of `line` or `bar`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("line", "bar"),
							},
						},
						"interpolation": schema.StringAttribute{
							MarkdownDescription: "The kind of interpolation to use between points displayed in the graph (line plots only). One of `linear`, `step`, `basis`, `bundle`, or `cardinal`.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("linear", "step", "basis", "bundle", "cardinal"),
							},
						},
						"aggregation": schema.StringAttribute{
							MarkdownDescription: "The kind of aggregation method to use in case multiple displayed data points share the same time-axis value (depending on the axis configured for time, by default x).\n\n-> this field does not have any effect on the underlying data; it is purely cosmetic, and applied only when viewing the data on the status page.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("sum", "average", "first", "last", "max", "min"),
							},
						},
					},
				},
			},
			"service": schema.SingleNestedAttribute{
				MarkdownDescription: "The service configuration for this MetricProvider, which describes how the given `instances` are provided.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Object{
					validators.ExactlyOneNonNullAttribute(),
				},
				Attributes: map[string]schema.Attribute{
					"builtin": schema.SingleNestedAttribute{
						MarkdownDescription: "The builtin Hund metric provider, which provides metrics based on recorded uptime and incidents.",
						Optional:            true,
						Attributes:          map[string]schema.Attribute{},
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
						},
					},
					"webhook": schema.SingleNestedAttribute{
						MarkdownDescription: "A [webhook](https://hund.io/help/documentation/incoming-webhook-metrics) service.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"webhook_key": schema.StringAttribute{
								MarkdownDescription: "The key to use for this webhook, expected in request headers.",
								Optional:            true,
								Computed:            true,
								Sensitive:           true,
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
	}
}

func (r *MetricProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		validators.MetricProviderInstancesMatchService(),
	}
}

func (r *MetricProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		var instances types.Map

		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("instances"), &instances)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if instances.IsUnknown() {
			var service types.Object
			resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("service"), &service)...)

			serviceType := models.MetricProviderServiceType(service)

			knownInstances, ok := models.MetricProviderServiceInstances()[serviceType]
			if knownInstances == nil {
				knownInstances = []string{}
			}
			if !ok {
				resp.Diagnostics.Append(MetricProviderServiceError(errors.New("unknown service type: " + serviceType)))
			}

			attrSchema, diag0 := req.Plan.Schema.AttributeAtPath(ctx, path.Root("instances"))
			resp.Diagnostics.Append(diag0...)

			instanceType, ok := attrSchema.GetType().(types.MapType)
			if !ok {
				resp.Diagnostics.AddError(
					"Unexpected schema type",
					"MetricProvider `instances` was not MapType",
				)
				return
			}

			instanceMap := map[string]models.MetricInstanceModel{}
			for _, v := range knownInstances {
				instanceMap[v] = models.MetricInstanceModelUnknown()
			}

			value, diag0 := types.MapValueFrom(ctx, instanceType.ElementType(), instanceMap)
			resp.Diagnostics.Append(diag0...)

			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("instances"), &value)...)
		}

		return
	}

	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state MetricProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Default.Equal(state.Default) {
		resp.RequiresReplace.Append(path.Root("default"))
	}

	if plan.Service.ServiceType() != state.Service.ServiceType() {
		resp.RequiresReplace.Append(path.Root("service"))
	}
}

func (r *MetricProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MetricProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, config MetricProviderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Default.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot create Default MetricProvider",
			"MetricProviders that are created automatically from an existing Watchdog"+
				" cannot be created directly in Terraform. Please import this resource instead:\n"+
				"\tterraform import <ADDRESS> default/"+data.Watchdog.ValueString(),
		)
		return
	}

	serviceForm, err := data.Service.ApiCreateForm()
	if err != nil {
		resp.Diagnostics.Append(MetricProviderServiceError(err))
		return
	}

	form := hundApiV1.MetricProviderFormCreate{
		Service:   serviceForm,
		Watchdog:  data.Watchdog.ValueString(),
		Instances: []hundApiV1.MetricInstanceFormCreate{},
	}

	// Only create instances given by practitioner config
	for slug, mim := range config.Instances {
		instanceForm, diag0 := mim.ApiCreateForm(slug)
		resp.Diagnostics.Append(diag0...)

		form.Instances = append(form.Instances, instanceForm)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.CreateAMetricProvider(ctx, form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund MetricProvider",
			err.Error(),
		)
		return
	}

	metric_provider, err := hundApiV1.ParseCreateAMetricProviderResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund MetricProvider",
			err.Error(),
		)
		return
	}

	if metric_provider.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-201 status code: "+fmt.Sprint(metric_provider.StatusCode())+
				"\nError: "+string(metric_provider.Body),
		)
		return
	}

	newState, diag := models.ToMetricProviderModel(*metric_provider.HALJSON201)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.Service.ReplaceSensitiveAttributes(*data.Service)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *MetricProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MetricProviderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rsp, err := r.client.RetrieveAMetricProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund MetricProvider",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	metric_provider, err := hundApiV1.ParseRetrieveAMetricProviderResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund MetricProvider",
			err.Error(),
		)
		return
	}

	if metric_provider.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(metric_provider.StatusCode())+
				"\nError: "+string(metric_provider.Body),
		)
		return
	}

	newState, diag := models.ToMetricProviderModel(*metric_provider.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Service != nil {
		newState.Service.ReplaceSensitiveAttributes(*data.Service)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *MetricProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, config, state MetricProviderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	var servicePlan, serviceState types.Object
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("service"), &servicePlan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("service"), &serviceState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	form := hundApiV1.MetricProviderFormUpdate{}

	if !servicePlan.Equal(serviceState) {
		serviceForm, err := data.Service.ApiUpdateForm()
		if err != nil {
			resp.Diagnostics.Append(MetricProviderServiceError(err))
			return
		}

		form.Service = &serviceForm
	}

	instances := []hundApiV1.MetricProviderFormUpdate_Instances_Item{}

	instanceSlugs := map[string]any{}

	for slug, mim := range config.Instances {
		if mim.Id.IsUnknown() {
			// Create MetricInstance
			instanceForm, diag0 := mim.ApiCreateForm(slug)
			resp.Diagnostics.Append(diag0...)

			item := hundApiV1.MetricProviderFormUpdate_Instances_Item{}
			err := item.FromMetricProviderFormUpdateInstances0(instanceForm)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error encoding MetricInstance Create Form",
					"Error: "+err.Error(),
				)
				return
			}

			instances = append(instances, item)
		} else {
			// Record this ID as existent, to prevent its removal below
			instanceSlugs[slug] = nil

			// Update MetricInstance
			instanceForm, diag0 := mim.ApiUpdateForm(slug)
			resp.Diagnostics.Append(diag0...)

			item := hundApiV1.MetricProviderFormUpdate_Instances_Item{}
			err := item.FromMetricProviderFormUpdateInstances1(instanceForm)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error encoding MetricInstance Update Form",
					"Error: "+err.Error(),
				)
				return
			}

			instances = append(instances, item)
		}
	}

	for slug := range state.Instances {
		_, ok := instanceSlugs[slug]

		if !ok {
			deleted := true

			instanceDeletion := hundApiV1.MetricInstanceFormEmbeddedUpdate{
				DefinitionSlug: slug,
				Deleted:        &deleted,
			}

			item := hundApiV1.MetricProviderFormUpdate_Instances_Item{}
			err := item.FromMetricProviderFormUpdateInstances1(instanceDeletion)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error encoding MetricInstance Deletion Form",
					"Error: "+err.Error(),
				)
				return
			}

			instances = append(instances, item)
		}
	}

	form.Instances = &instances

	rsp, err := r.client.UpdateAMetricProvider(ctx, data.Id.ValueString(), form)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Hund MetricProvider",
			err.Error(),
		)
		return
	}

	metric_provider, err := hundApiV1.ParseUpdateAMetricProviderResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund MetricProvider",
			err.Error(),
		)
		return
	}

	if metric_provider.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(metric_provider.StatusCode())+
				"\nError: "+string(metric_provider.Body),
		)
		return
	}

	newState, diag := models.ToMetricProviderModel(*metric_provider.HALJSON200)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState.Service.ReplaceSensitiveAttributes(*data.Service)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *MetricProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MetricProviderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Default.ValueBool() {
		return
	}

	rsp, err := r.client.DeleteAMetricProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Hund MetricProvider",
			err.Error(),
		)
		return
	}

	if rsp.StatusCode != 204 && rsp.StatusCode != 404 {
		summary := "Received a non-200 status code: " + fmt.Sprint(rsp.StatusCode)

		metric_provider, err := hundApiV1.ParseDeleteAMetricProviderResponse(rsp)

		if err == nil {
			summary = summary + "\nError: " + string(metric_provider.Body)
		}

		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			summary,
		)
		return
	}
}

func (r *MetricProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	watchdogId, viaWatchdog := strings.CutPrefix(req.ID, "default/")

	if viaWatchdog {
		model := r.findDefaultMetricProvider(ctx, watchdogId, &resp.Diagnostics)

		if model == nil {
			return
		}

		resp.State.SetAttribute(ctx, path.Root("id"), model.Id)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *MetricProviderResource) findDefaultMetricProvider(ctx context.Context, watchdogId string, diags *diag.Diagnostics) *MetricProviderResourceModel {
	rsp, err := r.client.GetAllMetricProviders(ctx, &hundApiV1.GetAllMetricProvidersParams{
		Watchdog: &watchdogId,
		Default:  hundApiV1.Ptr(true),
	})
	if err != nil {
		diags.AddError(
			"Unable to find default Hund MetricProvider",
			err.Error(),
		)
		return nil
	}

	metric_providers, err := hundApiV1.ParseGetAllMetricProvidersResponse(rsp)
	if err != nil {
		diags.AddError(
			"Unable to Parse Hund MetricProviders",
			err.Error(),
		)
		return nil
	}

	if metric_providers.StatusCode() != 200 {
		diags.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(metric_providers.StatusCode())+
				"\nError: "+string(metric_providers.Body),
		)
		return nil
	}

	if len(metric_providers.HALJSON200.Data) != 1 {
		diags.AddError(
			"Could not find default MetricProvider",
			fmt.Sprintf("It does not appear that this Watchdog (id=%[1]q) has a default"+
				" MetricProvider associated with it. Please set the `default` attribute to"+
				" false, or remove the attribute.", watchdogId),
		)

		return nil
	}

	model, diag := models.ToMetricProviderModel(metric_providers.HALJSON200.Data[0])
	diags.Append(diag...)

	return (*MetricProviderResourceModel)(&model)
}
