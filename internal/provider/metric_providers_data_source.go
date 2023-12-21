package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &MetricProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &MetricProvidersDataSource{}
)

func NewMetricProvidersDataSource() datasource.DataSource {
	return &MetricProvidersDataSource{}
}

// MetricProvidersDataSource defines the data source implementation.
type MetricProvidersDataSource struct {
	client *hundApiV1.Client
}

// MetricProvidersDataSourceModel describes the data source data model.
type MetricProvidersDataSourceModel struct {
	Watchdog types.String `tfsdk:"watchdog"`
	Default  types.Bool   `tfsdk:"default"`

	MetricProviders []models.MetricProviderModel `tfsdk:"metric_providers"`
}

func (d *MetricProvidersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_providers"
}

func (d *MetricProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "MetricProviders data source",

		Attributes: map[string]schema.Attribute{
			"watchdog": schema.StringAttribute{
				MarkdownDescription: "ObjectId for a particular Watchdog to retrieve MetricProviders on.",
				Optional:            true,
			},
			"default": schema.BoolAttribute{
				MarkdownDescription: "When true, returns only MetricProviders for which `default` is true (i.e. returns MetricProviders that are considered the \"default\" for their respective Watchdogs). When used in conjunction with the `watchdog` parameter, returns the *single* default MetricProvider of that Watchdog, if it exists.",
				Optional:            true,
			},
			"metric_providers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"watchdog": schema.StringAttribute{
							Computed: true,
						},
						"default": schema.BoolAttribute{
							Computed: true,
						},
						"instances": schema.MapNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"slug": schema.StringAttribute{
										Computed: true,
									},
									"definition_slug": schema.StringAttribute{
										Computed: true,
									},
									"enabled": schema.BoolAttribute{
										Computed: true,
									},
									"top_level_enabled": schema.BoolAttribute{
										Computed: true,
									},
									"title": schema.StringAttribute{
										Computed: true,
									},
									"title_translations": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"x_title": schema.StringAttribute{
										Computed: true,
									},
									"x_title_translations": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"y_title": schema.StringAttribute{
										Computed: true,
									},
									"y_title_translations": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"x_type": schema.StringAttribute{
										Computed: true,
									},
									"y_type": schema.StringAttribute{
										Computed: true,
									},
									"y_supremum": schema.Float64Attribute{
										Computed: true,
									},
									"plot_type": schema.StringAttribute{
										Computed: true,
									},
									"interpolation": schema.StringAttribute{
										Computed: true,
									},
									"aggregation": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"service": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"builtin": schema.SingleNestedAttribute{
									Computed:   true,
									Attributes: map[string]schema.Attribute{},
								},
								"updown": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"monitor_api_key": schema.StringAttribute{
											Computed:  true,
											Sensitive: true,
										},
										"monitor_token": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"pingdom": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"api_token": schema.StringAttribute{
											Computed:  true,
											Sensitive: true,
										},
										"check_id": schema.StringAttribute{
											Computed: true,
										},
										"check_type": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"uptimerobot": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"monitor_api_key": schema.StringAttribute{
											Computed:  true,
											Sensitive: true,
										},
									},
								},
								"webhook": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"webhook_key": schema.StringAttribute{
											Computed:  true,
											Sensitive: true,
										},
									},
								},
								"icmp": nativeIcmpServiceDataSourceSchema(),
								"http": nativeHttpServiceDataSourceSchema(),
								"dns":  nativeDnsServiceDataSourceSchema(),
								"tcp":  nativeTcpServiceDataSourceSchema(),
								"udp":  nativeUdpServiceDataSourceSchema(),
							},
						},
					},
				},
			},
		},
	}
}

func (d *MetricProvidersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *MetricProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MetricProvidersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	limit := 100
	rsp, err := d.client.GetAllMetricProviders(ctx, &hundApiV1.GetAllMetricProvidersParams{
		Watchdog: data.Watchdog.ValueStringPointer(),
		Default:  data.Default.ValueBoolPointer(),

		Limit: &limit,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund MetricProviders",
			err.Error(),
		)
		return
	}

	metricProviders, err := hundApiV1.ParseGetAllMetricProvidersResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund MetricProviders",
			err.Error(),
		)
		return
	}

	if metricProviders.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(metricProviders.StatusCode()),
		)
		return
	}

	// data.MetricProviders = []models.IssueTemplateModel{}

	for _, metricProvider := range metricProviders.HALJSON200.Data {
		metricProviderModel, diag := models.ToMetricProviderModel(metricProvider)
		resp.Diagnostics.Append(diag...)

		data.MetricProviders = append(data.MetricProviders, metricProviderModel)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
