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
	_ datasource.DataSource              = &GroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupsDataSource{}
)

func NewGroupsDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

// GroupsDataSource defines the data source implementation.
type GroupsDataSource struct {
	client *hundApiV1.Client
}

// GroupsDataSourceModel describes the data source data model.
type GroupsDataSourceModel struct {
	Groups []models.GroupModel `tfsdk:"groups"`
}

func (d *GroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *GroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Groups data source",

		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
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
						"collapsed": schema.BoolAttribute{
							Computed: true,
						},
						"position": schema.Int64Attribute{
							Computed: true,
						},
						"components": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *GroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	limit := 100
	rsp, err := d.client.GetAllGroups(ctx, &hundApiV1.GetAllGroupsParams{
		Limit: &limit,
	}, hundApiV1.Unexpand("data.components"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Hund Groups",
			err.Error(),
		)
		return
	}

	groups, err := hundApiV1.ParseGetAllGroupsResponse(rsp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Hund Groups",
			err.Error(),
		)
		return
	}

	if groups.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Failed response code from Hund API",
			"Received a non-200 status code: "+fmt.Sprint(groups.StatusCode()),
		)
		return
	}

	// data.Groups = []models.IssueTemplateModel{}

	for _, group := range groups.HALJSON200.Data {
		groupModel, diag := models.ToGroupModel(group)
		resp.Diagnostics.Append(diag...)

		data.Groups = append(data.Groups, groupModel)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
