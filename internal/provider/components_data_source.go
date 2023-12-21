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
	_ datasource.DataSource              = &ComponentsDataSource{}
	_ datasource.DataSourceWithConfigure = &ComponentsDataSource{}
)

func NewComponentsDataSource() datasource.DataSource {
	return &ComponentsDataSource{}
}

// ComponentsDataSource defines the data source implementation.
type ComponentsDataSource struct {
	client *hundApiV1.Client
}

// ComponentsDataSourceModel describes the data source data model.
type ComponentsDataSourceModel struct {
	Group types.String `tfsdk:"group"`
	Issue types.String `tfsdk:"issue"`

	Components []models.ComponentModel `tfsdk:"components"`
}

func (d *ComponentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_components"
}

func (d *ComponentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Components data source",

		Attributes: map[string]schema.Attribute{
			"group": schema.StringAttribute{
				MarkdownDescription: "Return the Components for the provided Group ObjectId.",
				Optional:            true,
			},
			"issue": schema.StringAttribute{
				MarkdownDescription: "Return the Components for the provided Issue ObjectId.",
				Optional:            true,
			},
			"components": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"group": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"name_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"description_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"description_html": schema.StringAttribute{
							Computed: true,
						},
						"description_html_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"exclude_from_global_history": schema.BoolAttribute{
							Computed: true,
						},
						"exclude_from_global_uptime": schema.BoolAttribute{
							Computed: true,
						},
						"last_event_at": schema.StringAttribute{
							Computed: true,
						},
						"percent_uptime": schema.Float64Attribute{
							Computed: true,
						},
						"watchdog": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"high_frequency": schema.BoolAttribute{
									Computed: true,
								},
								"latest_status": schema.StringAttribute{
									Computed: true,
								},
								"service": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"manual": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"state": schema.Int64Attribute{
													Computed: true,
												},
											},
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
												"unconfirmed_is_down": schema.BoolAttribute{
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
												"unconfirmed_is_down": schema.BoolAttribute{
													Computed: true,
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
												"deadman": schema.BoolAttribute{
													Computed: true,
												},
												"consecutive_checks": schema.Int64Attribute{
													Computed: true,
												},
												"reporting_interval": schema.Int64Attribute{
													Computed: true,
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
			},
		},
	}
}

func (d *ComponentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ComponentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ComponentsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	limit := 100
	rsp, err := d.client.GetAllComponents(ctx, &hundApiV1.GetAllComponentsParams{
		Group: data.Group.ValueStringPointer(),
		Issue: data.Issue.ValueStringPointer(),

		Limit: &limit,
	}, hundApiV1.Expand("data.watchdog"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Components",
			err.Error(),
		)
		return
	}

	components, err := hundApiV1.ParseGetAllComponentsResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Components",
			err.Error(),
		)
		return
	}

	if components.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(components.StatusCode()),
		)
		return
	}

	// data.Components = []models.IssueTemplateModel{}

	for _, component := range components.HALJSON200.Data {
		componentModel, diag := models.ToComponentModel(ctx, component)
		resp.Diagnostics.Append(diag...)

		data.Components = append(data.Components, componentModel)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
