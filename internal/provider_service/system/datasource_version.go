package system

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providerutil "github.com/marshallford/terraform-provider-pfsense/internal/provider_util"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = &VersionDataSource{}
	_ datasource.DataSourceWithConfigure = &VersionDataSource{}
)

func NewVersionDataSource() datasource.DataSource {
	return &VersionDataSource{}
}

type VersionDataSource struct {
	client *pfsense.Client
}

type VersionDataSourceModel struct {
	Current types.String `tfsdk:"current"`
	Latest  types.String `tfsdk:"latest"`
}

func (d *VersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_version", req.ProviderTypeName)
}

func (d *VersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves current and latest system version.",
		Attributes: map[string]schema.Attribute{
			"current": schema.StringAttribute{
				Description: "Current pfSense system version.",
				Computed:    true,
			},
			"latest": schema.StringAttribute{
				Description: "Latest pfSense system version.",
				Computed:    true,
			},
		},
	}
}

func (d *VersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := providerutil.ConfigureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *VersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VersionDataSourceModel

	version, err := d.client.GetSystemVersion(ctx)
	if providerutil.AddError(&resp.Diagnostics, "Unable to get system version", err) {
		return
	}

	data.Current = types.StringValue(version.Current)
	data.Latest = types.StringValue(version.Latest)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
