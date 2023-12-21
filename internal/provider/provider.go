package provider

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	hundApiV1 "github.com/hundio/terraform-provider-hund/internal/hund_api_v1"
)

// Ensure HundProvider satisfies various provider interfaces.
var _ provider.Provider = &HundProvider{}

// HundProvider defines the provider implementation.
type HundProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HundProviderModel describes the provider data model.
type HundProviderModel struct {
	Domain types.String `tfsdk:"domain"`
	Key    types.String `tfsdk:"key"`
}

func (p *HundProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hund"
	resp.Version = p.version
}

func (p *HundProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The [Hund](https://hund.io) provider offers several resources and data sources to provision and query various objects on a Hund hosted status page.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				MarkdownDescription: "The [domain](https://hund.io/help/api#section/Base-URL) at which to call the Hund API. Usually, this should be the domain of your status page.",
				Optional:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The [Hund API key](https://hund.io/help/api#section/Authentication) used to authenticate with the API.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *HundProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HundProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Domain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("domain"),
			"Unknown Hund API Domain",
			"The provider cannot create the Hund API client as there is an unknown configuration value for the Hund API domain. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HUND_DOMAIN environment variable.",
		)
	}

	if data.Key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown Hund API Key",
			"The provider cannot create the Hund API client as there is an unknown configuration value for the Hund API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HUND_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	domain := os.Getenv("HUND_DOMAIN")
	key := os.Getenv("HUND_KEY")

	if !data.Domain.IsNull() {
		domain = data.Domain.ValueString()
	}

	if !data.Key.IsNull() {
		key = data.Key.ValueString()
	}

	if domain == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("domain"),
			"Missing Hund API Domain",
			"The provider cannot create the Hund API client as there is a missing or empty value for the Hund API domain. "+
				"Set the domain value in the configuration or use the HUND_DOMAIN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Missing Hund API Key",
			"The provider cannot create the Hund API client as there is a missing or empty value for the Hund API key. "+
				"Set the key value in the configuration or use the HUND_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := newProviderClient(p.version, domain, key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hund API Client",
			"An unexpected error occurred when creating the Hund API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Hund Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *HundProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
		NewGroupComponentOrderingResource,
		NewComponentResource,
		NewMetricProviderResource,
		NewIssueResource,
		NewIssueUpdateResource,
		NewIssueTemplateResource,
	}
}

func (p *HundProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGroupsDataSource,
		NewComponentsDataSource,
		NewMetricProvidersDataSource,
		NewIssuesDataSource,
		NewIssueTemplatesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HundProvider{
			version: version,
		}
	}
}

func newProviderClient(version string, domain string, key string) (*hundApiV1.Client, error) {
	security, err := hundApiV1.WithSecurity(key)
	if err != nil {
		return nil, err
	}

	var endpoint string

	if strings.HasSuffix(domain, ".localhost") {
		endpoint = "http://" + domain + ":3000/api/v1"
	} else {
		endpoint = "https://" + domain + "/api/v1"
	}

	options := func(client *hundApiV1.Client) error {
		client.RequestEditors = append(client.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Hund-Version", "2021-09-01")
			req.Header.Add("Accept-Language", "*")
			req.Header.Set("User-Agent", "terraform-provider-hund/"+version)

			return nil
		})

		retryableClient := retryablehttp.NewClient()

		client.Client = retryableClient.StandardClient()

		return nil
	}

	return hundApiV1.NewClient(endpoint, security, options)
}
