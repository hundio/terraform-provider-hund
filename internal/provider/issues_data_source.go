package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &IssuesDataSource{}
	_ datasource.DataSourceWithConfigure = &IssuesDataSource{}
)

func NewIssuesDataSource() datasource.DataSource {
	return &IssuesDataSource{}
}

// IssuesDataSource defines the data source implementation.
type IssuesDataSource struct {
	client *hundApiV1.Client
}

type IssuesDataSourceModel struct {
	Upcoming types.Bool `tfsdk:"upcoming"`
	Standing types.Bool `tfsdk:"standing"`
	Resolved types.Bool `tfsdk:"resolved"`

	Components types.Set `tfsdk:"components"`

	Issues []models.IssueModel `tfsdk:"issues"`
}

func (d *IssuesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issues"
}

func (d *IssuesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Issues data source",

		Attributes: map[string]schema.Attribute{
			"upcoming": schema.BoolAttribute{
				MarkdownDescription: "When true, returns only upcoming scheduled Issues.",
				Optional:            true,
			},
			"resolved": schema.BoolAttribute{
				MarkdownDescription: "When true, returns only resolved Issues.",
				Optional:            true,
			},
			"standing": schema.BoolAttribute{
				MarkdownDescription: "When true, returns only ongoing Issues.",
				Optional:            true,
			},
			"components": schema.SetAttribute{
				MarkdownDescription: "One or more Components to return Issues for.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"issues": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"archive_on_destroy": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"began_at": schema.StringAttribute{
							Computed: true,
						},
						"ended_at": schema.StringAttribute{
							Computed: true,
						},
						"cancelled_at": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"body": schema.StringAttribute{
							Computed: true,
						},
						"body_html": schema.StringAttribute{
							Computed: true,
						},
						"title_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"body_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"body_html_translations": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"label":                schema.StringAttribute{Computed: true},
						"specialization":       schema.StringAttribute{Computed: true},
						"duration":             schema.Int64Attribute{Computed: true},
						"open_graph_image_url": schema.StringAttribute{Computed: true},
						"priority":             schema.Int64Attribute{Computed: true},
						"state_override":       schema.Int64Attribute{Computed: true},
						"resolved":             schema.BoolAttribute{Computed: true},
						"retrospective":        schema.BoolAttribute{Computed: true},
						"scheduled":            schema.BoolAttribute{Computed: true},
						"standing":             schema.BoolAttribute{Computed: true},
						"component_ids": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"schedule": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id":                    schema.StringAttribute{Computed: true},
								"started":               schema.BoolAttribute{Computed: true},
								"ended":                 schema.BoolAttribute{Computed: true},
								"notified":              schema.BoolAttribute{Computed: true},
								"starts_at":             schema.StringAttribute{Computed: true},
								"ends_at":               schema.StringAttribute{Computed: true},
								"notify_subscribers_at": schema.StringAttribute{Computed: true},
							},
						},
						"template": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"issue_template_id": schema.StringAttribute{
									Computed: true,
								},
								"title": schema.StringAttribute{
									Computed: true,
								},
								"body": schema.StringAttribute{
									Computed: true,
								},
								"title_translations": schema.MapAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"body_translations": schema.MapAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"label": schema.StringAttribute{
									Computed: true,
								},
								"schema": schema.MapNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Computed: true,
											},
											"required": schema.BoolAttribute{
												Computed: true,
											},
										},
									},
								},
								"variables": schema.MapNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"string":      schema.StringAttribute{Computed: true},
											"number":      schema.NumberAttribute{Computed: true},
											"i18n_string": schema.MapAttribute{Computed: true, ElementType: types.StringType},
											"datetime":    schema.StringAttribute{Computed: true},
										},
									},
								},
							},
						},
						"updates": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":                 schema.StringAttribute{Computed: true},
									"issue_id":           schema.StringAttribute{Computed: true},
									"archive_on_destroy": schema.BoolAttribute{Computed: true},
									"created_at":         schema.StringAttribute{Computed: true},
									"updated_at":         schema.StringAttribute{Computed: true},
									"effective":          schema.BoolAttribute{Computed: true},
									"reopening":          schema.BoolAttribute{Computed: true},
									"effective_after":    schema.StringAttribute{Computed: true},
									"body": schema.StringAttribute{
										Computed: true,
									},
									"body_html": schema.StringAttribute{
										Computed: true,
									},
									"body_translations": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"body_html_translations": schema.MapAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"label":          schema.StringAttribute{Computed: true},
									"state_override": schema.Int64Attribute{Computed: true},
									"template": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed: true,
											},
											"issue_template_id": schema.StringAttribute{
												Computed: true,
											},
											"body": schema.StringAttribute{
												Computed: true,
											},
											"body_translations": schema.MapAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"label": schema.StringAttribute{
												Computed: true,
											},
											"schema": schema.MapNestedAttribute{
												Computed: true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"type": schema.StringAttribute{
															Computed: true,
														},
														"required": schema.BoolAttribute{
															Computed: true,
														},
													},
												},
											},
											"variables": schema.MapNestedAttribute{
												Computed: true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"string":      schema.StringAttribute{Computed: true},
														"number":      schema.NumberAttribute{Computed: true},
														"i18n_string": schema.MapAttribute{Computed: true, ElementType: types.StringType},
														"datetime":    schema.StringAttribute{Computed: true},
													},
												},
											},
										},
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

func (d *IssuesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IssuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IssuesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var components *[]string
	if !data.Components.IsNull() {
		comps := []string{}

		for _, v := range data.Components.Elements() {
			str, ok := v.(types.String)
			if !ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("components"),
					"Cannot interpret component ID as string",
					"Failed to cast a given component ID to string type.",
				)

				return
			}

			comps = append(comps, str.ValueString())
		}

		components = &comps
	}

	limit := 100
	rsp, err := d.client.GetAllIssues(ctx, &hundApiV1.GetAllIssuesParams{
		Upcoming: data.Upcoming.ValueBoolPointer(),
		Standing: data.Standing.ValueBoolPointer(),
		Resolved: data.Resolved.ValueBoolPointer(),

		Components: components,

		Limit: &limit,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Issues",
			err.Error(),
		)
		return
	}

	issues, err := hundApiV1.ParseGetAllIssuesResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issues",
			err.Error(),
		)
		return
	}

	if issues.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(issues.StatusCode()),
		)
		return
	}

	data.Issues = []models.IssueModel{}

	for _, issue := range issues.HALJSON200.Data {
		issuesState, diag := models.ToIssueModel(ctx, issue)
		resp.Diagnostics.Append(diag...)

		data.Issues = append(data.Issues, issuesState)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
