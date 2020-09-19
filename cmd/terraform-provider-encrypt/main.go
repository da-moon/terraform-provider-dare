package main

import (
	provider "github.com/da-moon/terraform-provider-dare/pkg/provider"
	plugin "github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
