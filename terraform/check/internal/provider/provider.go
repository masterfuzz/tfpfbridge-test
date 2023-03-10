package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider satisfies various provider interfaces.
var _ provider.Provider = &CheckProvider{}

// CheckProvider defines the provider implementation.
type CheckProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Schema implements provider.Provider
func (*CheckProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

// CheckProviderModel describes the provider data model.
type CheckProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *CheckProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "check"
	resp.Version = p.version
}

func (p *CheckProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// 	var data CheckProviderModel

// 	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// Configuration values are now available.
// 	// if data.Endpoint.IsNull() { /* ... */ }

// 	// Example client configuration for data sources and resources
// 	client := http.DefaultClient
// 	resp.DataSourceData = client
// 	resp.ResourceData = client
// }

func (p *CheckProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHttpHealthResource,
		NewLocalCommandResource,
	}
}

func (p *CheckProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CheckProvider{
			version: version,
		}
	}
}
