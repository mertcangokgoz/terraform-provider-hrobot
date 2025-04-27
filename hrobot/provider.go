package hrobot

import (
	"errors"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/floshodan/terraform-provider-hrobot/internal/boot"
	"github.com/floshodan/terraform-provider-hrobot/internal/firewall"
	"github.com/floshodan/terraform-provider-hrobot/internal/sshkey"
	"github.com/floshodan/terraform-provider-hrobot/internal/vswitch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HROBOT_TOKEN", nil),
				Description: "The API token to access the Hetzner Robot Interface. Has the following Format: username:password ",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					token := val.(string)
					if token == "" {
						errs = append(errs, errors.New("entered token is invalid"))
					}
					return
				},
				Sensitive: true,
			}},
		ResourcesMap: map[string]*schema.Resource{
			//"hrobot_server_list": server.ResourceEvent(),
			//server.ResourceType: server.Resource(),
			sshkey.ResourceType:   sshkey.Resource(),
			firewall.ResourceType: firewall.Resource(),
			boot.ResourceType:     boot.Resource(),
			vswitch.ResourceType:  vswitch.Resource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			sshkey.DataSourceType:      sshkey.DataSource(),
			sshkey.DataSourceListType:  sshkey.DataSourceList(),
			firewall.DataSourceType:    firewall.DataSource(),
			boot.DataSourceType:        boot.DataSourceList(),
			vswitch.DataSourceType:     vswitch.DataSource(),
			vswitch.DataSourceListType: vswitch.DataSourceList(),
			//rescue.DataSourceType:     rescue.DataSource(),
		},
		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// TODO
	opts := []hrobot.ClientOption{
		hrobot.WithToken(d.Get("token").(string)),
	}

	//username := d.Get("username").(string)
	//password := d.Get("password").(string)
	return hrobot.NewClient(opts...), nil
}
