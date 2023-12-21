package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &IssueTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &IssueTemplatesDataSource{}
)

func NewIssueTemplatesDataSource() datasource.DataSource {
	return &IssueTemplatesDataSource{}
}

// IssueTemplatesDataSource defines the data source implementation.
type IssueTemplatesDataSource struct {
	client *hundApiV1.Client
}

// IssueTemplatesDataSourceModel describes the data source data model.
type IssueTemplatesDataSourceModel struct {
	Kind types.String `tfsdk:"kind"`

	IssueTemplates []models.IssueTemplateModel `tfsdk:"issue_templates"`
}

func (d *IssueTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_templates"
}

func (d *IssueTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "IssueTemplates data source",

		Attributes: map[string]schema.Attribute{
			"kind": schema.StringAttribute{
				MarkdownDescription: "Return only IssueTemplates for the given kind. Either `issue` or `update`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("issue", "update"),
				},
			},
			"issue_templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"kind": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
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
						"variables": schema.MapNestedAttribute{
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
					},
				},
			},
		},
	}
}

func (d *IssueTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IssueTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IssueTemplatesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	limit := 100
	rsp, err := d.client.GetAllIssueTemplates(ctx, &hundApiV1.GetAllIssueTemplatesParams{
		Kind: (*hundApiV1.GetAllIssueTemplatesParamsKind)(data.Kind.ValueStringPointer()),

		Limit: &limit,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Issue Templates",
			err.Error(),
		)
		return
	}

	templates, err := hundApiV1.ParseGetAllIssueTemplatesResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Issue Templates",
			err.Error(),
		)
		return
	}

	if templates.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(templates.StatusCode()),
		)
		return
	}

	// data.IssueTemplates = []models.IssueTemplateModel{}

	for _, template := range templates.HALJSON200.Data {
		templateState, diag := models.ToIssueTemplateModel(template)
		resp.Diagnostics.Append(diag...)

		data.IssueTemplates = append(data.IssueTemplates, templateState)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
