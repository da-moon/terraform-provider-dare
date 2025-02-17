package main

import (
	plugin "github.com/hashicorp/terraform-plugin-sdk/plugin"
	terraform "github.com/hashicorp/terraform-plugin-sdk/terraform"

	provider "github.com/da-moon/terraform-provider-dare/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return provider.Provider()
		},
	})
}
