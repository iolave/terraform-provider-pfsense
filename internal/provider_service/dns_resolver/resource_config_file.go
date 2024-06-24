package dnsresolver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providerutil "github.com/marshallford/terraform-provider-pfsense/internal/provider_util"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var _ resource.Resource = &ConfigFileResource{}
var _ resource.ResourceWithImportState = &ConfigFileResource{}

func NewConfigFileResource() resource.Resource {
	return &ConfigFileResource{}
}

type ConfigFileResource struct {
	client *pfsense.Client
}

type ConfigFileResourceModel struct {
	Name    types.String `tfsdk:"name"`
	Content types.String `tfsdk:"content"`
	Apply   types.Bool   `tfsdk:"apply"`
}

func (r *ConfigFileResourceModel) SetFromValue(ctx context.Context, configFile *pfsense.ConfigFile) diag.Diagnostics {
	r.Name = types.StringValue(configFile.Name)
	r.Content = types.StringValue(configFile.Content)

	return nil
}

func (r ConfigFileResourceModel) Value(ctx context.Context) (*pfsense.ConfigFile, diag.Diagnostics) {
	var configFile pfsense.ConfigFile
	var err error
	var diags diag.Diagnostics

	err = configFile.SetName(r.Name.ValueString())
	if err != nil {
		diags.AddAttributeError(
			path.Root("name"),
			"Name cannot be parsed",
			err.Error(),
		)
	}

	err = configFile.SetContent(r.Content.ValueString())
	if err != nil {
		diags.AddAttributeError(
			path.Root("content"),
			"Content cannot be parsed",
			err.Error(),
		)
	}

	return &configFile, diags
}

func (r *ConfigFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_configfile", req.ProviderTypeName)
}

func (r *ConfigFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "DNS resolver (Unbound) config file. Prerequisite: Must add the directive 'include-toplevel: /var/unbound/conf.d/*' to the DNS resolver custom options input. Use with caution, content is not checked/validated.",
		MarkdownDescription: "DNS resolver (Unbound) [config file](https://man.freebsd.org/cgi/man.cgi?unbound.conf). **Prerequisite**: Must add the directive `include-toplevel: /var/unbound/conf.d/*` to the DNS resolver custom options input. **Use with caution**, content is not checked/validated.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of config file.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description:         "Contents of file. Must specify Unbound clause(s). Comments start with '#' and last to the end of line.",
				MarkdownDescription: "Contents of file. Must specify Unbound clause(s). Comments start with `#` and last to the end of line.",
				Required:            true,
			},
			"apply": schema.BoolAttribute{
				Description:         "Apply change, defaults to 'true'.",
				MarkdownDescription: "Apply change, defaults to `true`.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *ConfigFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := providerutil.ConfigureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *ConfigFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ConfigFileResourceModel
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configFileReq, d := data.Value(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	configFile, err := r.client.CreateDNSResolverConfigFile(ctx, *configFileReq)
	if providerutil.AddError(&resp.Diagnostics, "Error creating config file", err) {
		return
	}

	diags = data.SetFromValue(ctx, configFile)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		if providerutil.AddError(&resp.Diagnostics, "Error applying config file", err) {
			return
		}
	}
}

func (r *ConfigFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ConfigFileResourceModel
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configFile, err := r.client.GetDNSResolverConfigFile(ctx, data.Name.ValueString())
	if providerutil.AddError(&resp.Diagnostics, "Error reading config file", err) {
		return
	}

	diags = data.SetFromValue(ctx, configFile)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ConfigFileResourceModel
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configFileReq, d := data.Value(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configFile, err := r.client.UpdateDNSResolverConfigFile(ctx, *configFileReq)
	if providerutil.AddError(&resp.Diagnostics, "Error updating config file", err) {
		return
	}

	diags = data.SetFromValue(ctx, configFile)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		if providerutil.AddError(&resp.Diagnostics, "Error applying config file", err) {
			return
		}
	}
}

func (r *ConfigFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ConfigFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDNSResolverConfigFile(ctx, data.Name.ValueString())
	if providerutil.AddError(&resp.Diagnostics, "Error deleting config file", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		if providerutil.AddError(&resp.Diagnostics, "Error applying config file", err) {
			return
		}
	}
}

func (r *ConfigFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
