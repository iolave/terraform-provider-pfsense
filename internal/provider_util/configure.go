package providerutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

func UnknownProviderValue(value string) (string, string) {
	return fmt.Sprintf("Unknown pfSense %s", value),
		fmt.Sprintf("The provider cannot create the pfSense client as there is an unknown configuration value for the %s. ", value) +
			"Either target apply the source of the value first, set the value statically in the configuration."
}

func UnexpectedConfigureType(value string, providerData any) (string, string) {
	return fmt.Sprintf("Unexpected %s Configure Type", value),
		fmt.Sprintf("Expected *pfsense.Client, got: %T. Please report this issue to the provider developers.", providerData)
}

func ConfigureDataSourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (*pfsense.Client, bool) {
	if req.ProviderData == nil {
		return nil, false
	}

	client, ok := req.ProviderData.(*pfsense.Client)

	if !ok {
		summary, detail := UnexpectedConfigureType("Data Source", req.ProviderData)
		resp.Diagnostics.AddError(summary, detail)
	}

	return client, ok
}

func ConfigureResourceClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse) (*pfsense.Client, bool) {
	if req.ProviderData == nil {
		return nil, false
	}

	client, ok := req.ProviderData.(*pfsense.Client)

	if !ok {
		summary, detail := UnexpectedConfigureType("Resource", req.ProviderData)
		resp.Diagnostics.AddError(summary, detail)
	}

	return client, ok
}

func AddError(diag *diag.Diagnostics, summary string, err error) bool {
	if err != nil {
		diag.AddError(summary, fmt.Sprintf("unexpected error: %v", err))
		return true
	}
	return false
}
