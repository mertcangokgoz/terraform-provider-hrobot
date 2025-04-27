package main

import (
	"github.com/floshodan/terraform-provider-hrobot/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hrobot.Provider})
}
