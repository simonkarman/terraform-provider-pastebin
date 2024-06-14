package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/simonkarman/pastebin-client-go"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &pastebinProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pastebinProvider{
			version: version,
		}
	}
}

// pastebinProvider is the provider implementation.
type pastebinProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *pastebinProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pastebin"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
type pastebinProviderModel struct {
	Host    types.String `tfsdk:"host"`
	DevKey  types.String `tfsdk:"dev_key"`
	UserKey types.String `tfsdk:"user_key"`
}

func (p *pastebinProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"dev_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"user_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a pastebin API client for data sources and resources.
func (p *pastebinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config pastebinProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown PasteBin API Host",
			"The provider cannot create the PasteBin API client as there is an unknown configuration value for the PasteBin API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PASTEBIN_HOST environment variable.",
		)
	}

	if config.DevKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dev_key"),
			"Unknown PasteBin API Dev Key",
			"The provider cannot create the PasteBin API client as there is an unknown configuration value for the PasteBin API dev key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PASTEBIN_DEV_KEY environment variable.",
		)
	}

	if config.UserKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_key"),
			"Unknown PasteBin API User Key",
			"The provider cannot create the PasteBin API client as there is an unknown configuration value for the PasteBin API user key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PASTEBIN_USER_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("PASTEBIN_HOST")
	devKey := os.Getenv("PASTEBIN_DEV_KEY")
	userKey := os.Getenv("PASTEBIN_USER_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.DevKey.IsNull() {
		devKey = config.DevKey.ValueString()
	}

	if !config.UserKey.IsNull() {
		userKey = config.UserKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		host = "https://pastebin.com"
	}
	hostUrl, err := url.Parse(host)
	if err != nil || hostUrl == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Invalid PasteBin API Host",
			"The provider cannot create the PasteBin API client as the provided host is not a valid URL. "+
				"Ensure the host value or the PASTEBIN_HOST environment variable is set to a valid url. "+
				"If you leave the host value empty, the default 'https://pastebin.com' url is used.",
		)
	}

	if devKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("dev_key"),
			"Missing PasteBin API Dev Key",
			"The provider cannot create the PasteBin API client as there is a missing or empty value for the PasteBin API dev key. "+
				"Set the dev_key value in the configuration or use the PASTEBIN_DEV_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if userKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_key"),
			"Missing PasteBin API User Key",
			"The provider cannot create the PasteBin API client as there is a missing or empty value for the PasteBin API user key. "+
				"Set the dev_key value in the configuration or use the PASTEBIN_USER_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new PasteBin client using the configuration values
	client := pastebin.New(*hostUrl, devKey, userKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create PasteBin API Client",
			"An unexpected error occurred when creating the PasteBin API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"PasteBin Client Error: "+err.Error(),
		)
		return
	}

	// Make the PasteBin client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *pastebinProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *pastebinProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
