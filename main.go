package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/skysoft-atm/terraform-provider-elastic/elastic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: elastic.Provider,
	})
}
