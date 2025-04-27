package firewall

import (
	"context"

	"github.com/floshodan/hrobot-go/hrobot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

const (
	DataSourceType                     = "hrobot_firewall"
	DataSourceTypeTemplate             = "hrobot_firewall_template"
	DataSourceFirewallTemplateListType = "hrobot_firewall_templates"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHrobotFirewallRead,
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"whitelist_hos": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dst_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tcp_flags": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func validateOneServerIdentifier(v any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	m := v.(map[string]interface{})
	serverIp := m["server_ip"]
	serverNumber := m["server_number"]

	if serverIp == nil && serverNumber == nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Both server_ip and server_number are nil",
			Detail:   "One of server_ip or server_number must be provided",
		}
		diags = append(diags, diag)
	} else if serverIp != nil && serverNumber != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Both server_ip and server_number are set",
			Detail:   "Only one of server_ip or server_number can be provided",
		}
		diags = append(diags, diag)
	}
	return diags
}

func dataSourceHrobotFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hrobot.Client)

	if id, ok := d.GetOk("server_ip"); ok {
		i, _, err := client.Firewall.GetFirewallByIP(ctx, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if i == nil {
			return diag.Errorf("no firewall found with id %d", id)
		}
		setFirewallSchema(d, i)
		return nil
	}

	return nil
}
